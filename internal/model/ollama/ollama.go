// Package ollama implements the model.LLM interface for Ollama's OpenAI-compatible API.
package ollama

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

type ollamaModel struct {
	client    *openai.Client
	modelName string
}

// NewModel returns a [model.LLM] backed by Ollama's OpenAI-compatible API.
//
// baseURL should point at the Ollama server (e.g. "http://localhost:11434").
// modelName is the Ollama model tag (e.g. "llama3.2", "qwen2.5:7b").
func NewModel(baseURL, modelName string) model.LLM {
	cfg := openai.DefaultConfig("ollama") // Ollama ignores the token
	cfg.BaseURL = strings.TrimRight(baseURL, "/") + "/v1"
	return &ollamaModel{
		client:    openai.NewClientWithConfig(cfg),
		modelName: modelName,
	}
}

func (m *ollamaModel) Name() string {
	return m.modelName
}

// GenerateContent implements model.LLM.
func (m *ollamaModel) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	msgs := contentsToMessages(req)
	if stream {
		return m.generateStream(ctx, msgs)
	}
	return func(yield func(*model.LLMResponse, error) bool) {
		resp, err := m.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:    m.modelName,
			Messages: msgs,
		})
		if err != nil {
			yield(nil, fmt.Errorf("ollama: %w", err))
			return
		}
		if len(resp.Choices) == 0 {
			yield(nil, fmt.Errorf("ollama: empty response"))
			return
		}
		yield(&model.LLMResponse{
			Content:      genai.NewContentFromText(resp.Choices[0].Message.Content, "model"),
			TurnComplete: true,
			FinishReason: genai.FinishReasonStop,
		}, nil)
	}
}

func (m *ollamaModel) generateStream(ctx context.Context, msgs []openai.ChatCompletionMessage) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		stream, err := m.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
			Model:    m.modelName,
			Messages: msgs,
			Stream:   true,
		})
		if err != nil {
			yield(nil, fmt.Errorf("ollama stream: %w", err))
			return
		}
		defer stream.Close()

		// Accumulate the full response so we can yield a final non-partial event.
		// ADK requires the last yielded event to have Partial=false (TurnComplete).
		var accumulated strings.Builder
		for {
			chunk, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				yield(nil, fmt.Errorf("ollama stream recv: %w", err))
				return
			}
			if len(chunk.Choices) == 0 {
				continue
			}
			delta := chunk.Choices[0].Delta.Content
			if delta == "" {
				continue
			}
			accumulated.WriteString(delta)
			// Yield each partial delta for live display in the UI.
			if !yield(&model.LLMResponse{
				Content: genai.NewContentFromText(delta, "model"),
				Partial: true,
			}, nil) {
				return
			}
		}

		// Yield the final aggregated response (non-partial). This is what ADK
		// uses as the last event and stores in the conversation history.
		if accumulated.Len() > 0 {
			yield(&model.LLMResponse{
				Content:      genai.NewContentFromText(accumulated.String(), "model"),
				TurnComplete: true,
				FinishReason: genai.FinishReasonStop,
			}, nil)
		}
	}
}

// contentsToMessages converts an ADK LLMRequest into OpenAI chat messages.
func contentsToMessages(req *model.LLMRequest) []openai.ChatCompletionMessage {
	var msgs []openai.ChatCompletionMessage

	// System instruction from the agent's Instruction field
	if req.Config != nil && req.Config.SystemInstruction != nil {
		if text := partsToText(req.Config.SystemInstruction.Parts); text != "" {
			msgs = append(msgs, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: text,
			})
		}
	}

	for _, c := range req.Contents {
		if c == nil {
			continue
		}
		text := partsToText(c.Parts)
		if text == "" {
			continue
		}
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    genaiRoleToOpenAI(c.Role),
			Content: text,
		})
	}
	return msgs
}

func genaiRoleToOpenAI(role string) string {
	switch role {
	case "model":
		return openai.ChatMessageRoleAssistant
	case "system":
		return openai.ChatMessageRoleSystem
	default:
		return openai.ChatMessageRoleUser
	}
}

func partsToText(parts []*genai.Part) string {
	var sb strings.Builder
	for _, p := range parts {
		if p != nil {
			sb.WriteString(p.Text)
		}
	}
	return sb.String()
}
