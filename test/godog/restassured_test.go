package godog

import (
	"context"
	"os"

	"github.com/jesse0michael/go-rest-assured/v4/pkg/assured"
)

func (s *Steps) iHaveACleanRestAssuredEnvironment(ctx context.Context) error {
	return s.client.ClearAll()
}

func (s *Steps) restAssuredReturns(ctx context.Context, data, method, path string) error {
	b, err := os.ReadFile(data)
	if err != nil {
		return err
	}

	return s.client.Given(assured.Call{
		Method:   method,
		Path:     path,
		Response: b,
	})
}
