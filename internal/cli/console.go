// Package cli provides a styled bubbletea console launcher for ADK agents.
package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

// Launcher is an ADK SubLauncher that runs an interactive bubbletea TUI.
type Launcher struct{}

// New returns a new console SubLauncher.
func New() launcher.SubLauncher { return &Launcher{} }

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

	m := newModel(ctx, r, userID, resp.Session.ID())
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithContext(ctx))
	_, err = p.Run()
	return err
}

// styles
var (
	barStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	userStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("250")) // light grey
	inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))  // white
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
	hintStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Italic(true)
)

// streaming bridge messages
type streamChunk struct{ text string }
type streamDone struct{}
type streamErr struct{ err error }

type model struct {
	ctx       context.Context
	runner    *runner.Runner
	sessionID string
	userID    string

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
}

func newModel(ctx context.Context, r *runner.Runner, userID, sessionID string) *model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Prompt = inputStyle.Render("> ")
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.SetHeight(3)
	ta.Focus()

	return &model{
		ctx:       ctx,
		runner:    r,
		sessionID: sessionID,
		userID:    userID,
		input:     ta,
		events:    make(chan tea.Msg, 64),
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.waitForEvent())
}

func (m *model) waitForEvent() tea.Cmd {
	return func() tea.Msg { return <-m.events }
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.relayout()
		m.vp.SetContent(m.renderHistory())
		m.vp.GotoBottom()
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
			m.input.Reset()
			m.turns = append(m.turns, renderUser(text, m.contentWidth()))
			m.vp.SetContent(m.renderHistory())
			m.vp.GotoBottom()
			m.busy = true
			return m, m.startTurn(text)
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
	bar := barStyle.Render(strings.Repeat("─", m.width))
	hint := ""
	if m.busy {
		hint = hintStyle.Render(" thinking… (esc to cancel)")
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		m.vp.View(),
		bar+hint,
		m.input.View(),
	)
}

func (m *model) relayout() {
	const (
		barH      = 1
		inputH    = 3
		minVPH    = 3
	)
	vpH := m.height - barH - inputH
	if vpH < minVPH {
		vpH = minVPH
	}
	if !m.ready {
		m.vp = viewport.New(m.width, vpH)
		m.ready = true
	} else {
		m.vp.Width = m.width
		m.vp.Height = vpH
	}
	m.input.SetWidth(m.width)
	m.input.SetHeight(inputH)

	if w := m.contentWidth(); w != m.mdWidth {
		if r, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
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
			for _, p := range ev.LLMResponse.Content.Parts {
				s += p.Text
			}
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
