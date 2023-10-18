package godog

import (
	"os"

	"github.com/jesse0michael/go-rest-assured/v4/pkg/assured"
)

func (s *Steps) iHaveACleanRestAssuredEnvironment() error {
	return s.client.ClearAll()
}

func (s *Steps) restAssuredReturns(data, method, path string) error {
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
