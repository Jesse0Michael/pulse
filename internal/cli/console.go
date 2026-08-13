// Package cli provides a styled bubbletea console launcher for ADK agents.
package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

// Launcher is an ADK SubLauncher that runs an interactive bubbletea TUI.
type Launcher struct {
	AgentName string
	ModelName string
}

// New returns a new console SubLauncher.
func New() *Launcher { return &Launcher{} }

func (*Launcher) Keyword() string           { return "console" }
func (*Launcher) SimpleDescription() string { return "interactive TUI chat" }
func (*Launcher) CommandLineSyntax() string { return "" }

func (*Launcher) Parse(args []string) ([]string, error) { return args, nil }

func (l *Launcher) Run(ctx context.Context, cfg *adk.Config) error {
	const userID, appName = "console_user", "console_app"

	sess := cfg.SessionService
	if sess == nil {
		sess = session.InMemoryService()
	}
	resp, err := sess.Create(ctx, &session.CreateRequest{AppName: appName, UserID: userID})
	if err != nil {
		return fmt.Errorf("creating session: %w", err)
	}

	r, err := runner.New(runner.Config{
		AppName:         appName,
		Agent:           cfg.AgentLoader.RootAgent(),
		SessionService:  sess,
		ArtifactService: cfg.ArtifactService,
	})
	if err != nil {
		return fmt.Errorf("creating runner: %w", err)
	}

	m := newModel(ctx, r, userID, resp.Session.ID(), l.AgentName, l.ModelName)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithContext(ctx), tea.WithInputTTY())
	_, err = p.Run()
	return err
}

// styles.
var (
	barStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	userStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("250")) // light grey
	inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))  // white
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
	hintStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Italic(true)
)

// streaming bridge messages.
type streamChunk struct{ text string }
type streamDone struct{}
type streamErr struct{ err error }
type focusMsg struct{}

// slashCommand describes an in-console command (typed as "/name ...").
type slashCommand struct {
	Name string
	Desc string
}

var slashCommands = []slashCommand{
	{Name: "skill", Desc: "run a skill"},
}

type model struct {
	ctx       context.Context
	runner    *runner.Runner
	sessionID string
	userID    string

	agentName string
	modelName string

	vp    viewport.Model
	input textarea.Model

	turns     []string // finalized rendered turns
	streaming string   // partial AI text in flight
	width     int
	height    int
	busy      bool
	ready     bool
	events    chan tea.Msg
	cancel    context.CancelFunc

	md      *glamour.TermRenderer // refreshed on resize
	mdWidth int                   // width md was built for

	skills []Skill
}

func newModel(ctx context.Context, r *runner.Runner, userID, sessionID, agentName, modelName string) *model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Prompt = inputStyle.Render("> ")
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.SetHeight(1)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.Blur()

	return &model{
		ctx:       ctx,
		runner:    r,
		sessionID: sessionID,
		userID:    userID,
		agentName: agentName,
		modelName: modelName,
		input:     ta,
		events:    make(chan tea.Msg, 64),
		skills:    LoadSkills(),
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.waitForEvent())
}

func (m *model) waitForEvent() tea.Cmd {
	return func() tea.Msg {
		select {
		case msg := <-m.events:
			return msg
		case <-m.ctx.Done():
			return nil
		}
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		wasReady := m.ready
		m.width = msg.Width
		m.height = msg.Height
		m.relayout()
		m.vp.SetContent(m.renderHistory())
		m.vp.GotoBottom()
		if !wasReady {
			return m, tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg { return focusMsg{} })
		}
		return m, nil

	case focusMsg:
		m.input.Reset()
		m.input.Focus()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit
		case "esc":
			if m.busy && m.cancel != nil {
				m.cancel()
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			if m.busy {
				return m, nil
			}
			text := strings.TrimSpace(m.input.Value())
			if text == "" {
				return m, nil
			}
			payload, display, cmdErr := m.expandCommand(text)
			if cmdErr != nil {
				m.input.Reset()
				m.turns = append(m.turns, errorStyle.Render("error: "+cmdErr.Error()))
				m.vp.SetContent(m.renderHistory())
				m.vp.GotoBottom()
				return m, nil
			}
			m.input.Reset()
			m.turns = append(m.turns, renderUser(display, m.contentWidth()))
			m.vp.SetContent(m.renderHistory())
			m.vp.GotoBottom()
			m.busy = true
			return m, m.startTurn(payload)
		}

	case streamChunk:
		m.streaming += msg.text
		m.vp.SetContent(m.renderHistoryWithStream())
		m.vp.GotoBottom()
		return m, m.waitForEvent()

	case streamDone:
		if m.streaming != "" {
			m.turns = append(m.turns, m.renderAIFinal(m.streaming))
			m.streaming = ""
		}
		m.busy = false
		m.cancel = nil
		m.vp.SetContent(m.renderHistory())
		m.vp.GotoBottom()
		return m, m.waitForEvent()

	case streamErr:
		m.turns = append(m.turns, errorStyle.Render("error: "+msg.err.Error()))
		m.streaming = ""
		m.busy = false
		m.cancel = nil
		m.vp.SetContent(m.renderHistory())
		m.vp.GotoBottom()
		return m, m.waitForEvent()
	}

	var cmds []tea.Cmd
	var c tea.Cmd
	m.input, c = m.input.Update(msg)
	cmds = append(cmds, c)
	m.vp, c = m.vp.Update(msg)
	cmds = append(cmds, c)
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if !m.ready {
		return ""
	}
	topBar := m.topBar()
	topSep := barStyle.Render(strings.Repeat("─", m.width))
	botSep := barStyle.Render(strings.Repeat("─", m.width))
	hint := ""
	if m.busy {
		hint = hintStyle.Render(" thinking… (esc to cancel)")
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		topBar,
		topSep,
		m.vp.View(),
		botSep+hint,
		m.input.View(),
		m.helperLine(),
	)
}

func (m *model) relayout() {
	const (
		topBarH = 1
		topSepH = 1
		botSepH = 1
		inputH  = 1
		helperH = 1
		minVPH  = 3
	)
	vpH := max(m.height-topBarH-topSepH-botSepH-inputH-helperH, minVPH)
	if !m.ready {
		m.vp = viewport.New(m.width, vpH)
		m.ready = true
	} else {
		m.vp.Width = m.width
		m.vp.Height = vpH
	}
	m.input.SetWidth(m.width)
	m.input.SetHeight(1)

	if w := m.contentWidth(); w != m.mdWidth {
		if r, err := glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(w*9/10),
		); err == nil {
			m.md = r
			m.mdWidth = w
		}
	}
}

func (m *model) contentWidth() int {
	if m.width <= 0 {
		return 80
	}
	return m.width
}

func (m *model) topBar() string {
	var parts []string
	if m.agentName != "" {
		parts = append(parts, "agent: "+m.agentName)
	}
	if m.modelName != "" {
		parts = append(parts, "model: "+m.modelName)
	}
	info := strings.Join(parts, "  │  ")
	return barStyle.Render(info)
}

func (m *model) renderHistory() string {
	return strings.Join(m.turns, "\n")
}

func (m *model) renderHistoryWithStream() string {
	if m.streaming == "" {
		return m.renderHistory()
	}
	tail := renderAI(m.streaming, m.contentWidth())
	if len(m.turns) == 0 {
		return tail
	}
	return m.renderHistory() + "\n" + tail
}

func renderUser(text string, w int) string {
	bubbleW := w * 7 / 10
	if bubbleW < 10 {
		bubbleW = w
	}
	bubble := userStyle.Width(bubbleW).Align(lipgloss.Right).Render(text)
	return lipgloss.PlaceHorizontal(w, lipgloss.Right, bubble)
}

func renderAI(text string, w int) string {
	bubbleW := w * 9 / 10
	if bubbleW < 10 {
		bubbleW = w
	}
	return lipgloss.NewStyle().Width(bubbleW).Render(text)
}

// renderAIFinal renders the finalized AI turn through glamour for markdown
// styling. Falls back to renderAI if the renderer isn't ready or errors.
func (m *model) renderAIFinal(text string) string {
	if m.md == nil {
		return renderAI(text, m.contentWidth())
	}
	out, err := m.md.Render(text)
	if err != nil {
		return renderAI(text, m.contentWidth())
	}
	return strings.TrimRight(out, "\n")
}

// expandCommand resolves slash commands like "/skill weekly-summary" into the
// payload sent to the agent and the text shown in the user bubble. Non-slash
// input passes through unchanged.
func (m *model) expandCommand(text string) (payload, display string, err error) {
	if !strings.HasPrefix(text, "/") {
		return text, text, nil
	}
	parts := strings.SplitN(strings.TrimPrefix(text, "/"), " ", 2)
	cmd := parts[0]
	arg := ""
	if len(parts) > 1 {
		arg = strings.TrimSpace(parts[1])
	}
	switch cmd {
	case "skill":
		if arg == "" {
			return "", "", errors.New("usage: /skill <name>")
		}
		var found *Skill
		for i := range m.skills {
			if m.skills[i].Name == arg {
				found = &m.skills[i]
				break
			}
		}
		if found == nil {
			return "", "", fmt.Errorf("skill %q not found", arg)
		}
		content, cerr := found.Content()
		if cerr != nil {
			return "", "", cerr
		}
		return "Execute the following skill:\n\n" + content, "/skill " + arg, nil
	}
	return "", "", fmt.Errorf("unknown command: /%s", cmd)
}

// helperLine returns the contextual hint rendered under the input. Empty
// input or non-slash input yields an empty line so the layout stays stable.
func (m *model) helperLine() string {
	val := strings.TrimLeft(m.input.Value(), " \t")
	if !strings.HasPrefix(val, "/") {
		return ""
	}
	rest := strings.TrimPrefix(val, "/")
	// If a space appears, the user has moved past the command name and is
	// typing an argument — show completions for that argument.
	if idx := strings.IndexAny(rest, " \n"); idx >= 0 {
		cmd := rest[:idx]
		arg := strings.TrimSpace(rest[idx:])
		switch cmd {
		case "skill":
			var names []string
			for _, s := range m.skills {
				if arg == "" || strings.HasPrefix(s.Name, arg) {
					names = append(names, s.Name)
				}
			}
			if len(names) == 0 {
				return hintStyle.Render("  no skills match")
			}
			return hintStyle.Render("  skills: " + strings.Join(names, ", "))
		}
		return hintStyle.Render("  unknown command")
	}
	// No space yet — list commands matching the partial name.
	var matches []string
	for _, c := range slashCommands {
		if strings.HasPrefix(c.Name, rest) {
			matches = append(matches, "/"+c.Name+" — "+c.Desc)
		}
	}
	if len(matches) == 0 {
		return hintStyle.Render("  no commands match")
	}
	return hintStyle.Render("  " + strings.Join(matches, "   "))
}

// startTurn launches a goroutine that ranges over the runner's event stream
// and posts streaming messages back to the bubbletea program via m.events.
func (m *model) startTurn(text string) tea.Cmd {
	msg := genai.NewContentFromText(text, genai.RoleUser)
	turnCtx, cancel := context.WithCancel(m.ctx)
	m.cancel = cancel

	go func() {
		defer cancel()
		var prev string
		for ev, err := range m.runner.Run(turnCtx, m.userID, m.sessionID, msg, agent.RunConfig{
			StreamingMode: agent.StreamingModeSSE,
		}) {
			if err != nil {
				m.events <- streamErr{err: err}
				return
			}
			if ev == nil || ev.LLMResponse.Content == nil {
				continue
			}
			var s string
			var sSb418 strings.Builder
			for _, p := range ev.LLMResponse.Content.Parts {
				sSb418.WriteString(p.Text)
			}
			s += sSb418.String()
			if s == "" {
				continue
			}
			// Mirror ADK's console: skip a duplicate final-response event
			// when its text matches what we already streamed.
			if ev.IsFinalResponse() && s == prev {
				continue
			}
			m.events <- streamChunk{text: s}
			prev += s
		}
		m.events <- streamDone{}
	}()
	return m.waitForEvent()
}
