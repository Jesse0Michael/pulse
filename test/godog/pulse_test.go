package godog

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

func (s *Steps) iRunThePulseGithubCommandOnUser(ctx context.Context, user string) (context.Context, error) {
	cmd := exec.Command( // nolint:gosec
		s.cfg.PulseCLIPath,
		"github",
		user,
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("GITHUB_URL=%s/", s.client.URL()), fmt.Sprintf("OPENAI_URL=%s", s.client.URL()))

	err := cmd.Run()
	output := stdout.String()
	fmt.Println(stderr.String())
	fmt.Println(output)
	ctx = context.WithValue(ctx, outputContextKey, output)

	return ctx, err
}

func (s *Steps) thePulseOutputShouldEqual(ctx context.Context, data string) error {
	output, _ := ctx.Value(outputContextKey).(string)

	e, err := os.ReadFile(data)
	if err != nil {
		return err
	}

	if output != string(e) {
		return fmt.Errorf("expected output to equal:\n%s\nrecieved:\n%s", string(e), output)
	}

	return nil
}
