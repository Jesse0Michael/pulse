package service

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/jesse0michael/testhelpers/pkg/testserver"
)

func TestOpenAI_Summarize(t *testing.T) {
	resp, _ := os.ReadFile("testdata/openai_chatcompletion.json")
	tests := []struct {
		name    string
		server  *testserver.Server
		want    string
		wantErr bool
	}{
		{
			name: "list user activity",
			server: testserver.New(
				testserver.Handler{Path: "/chat/completions", Status: http.StatusOK, Response: resp},
			),
			want:    `Hello there, how may I assist you today?`,
			wantErr: false,
		},
		{
			name: "failed list user activity",
			server: testserver.New(
				testserver.Handler{Path: "/chat/completions", Status: http.StatusTooManyRequests, Response: []byte(`{"message": "Too Many Requests"}`)},
			),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewOpenAI(OpenAIConfig{URL: tt.server.URL})
			got, err := o.Summarize(context.TODO(), "test-content")
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenAI.Summarize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OpenAI.Summarize() = %v, want %v", got, tt.want)
			}
		})
	}
}
