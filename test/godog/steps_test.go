package godog

import (
	"github.com/jesse0michael/go-rest-assured/v4/pkg/assured"
)

type contextKey string

const (
	outputContextKey = contextKey("output")
)

type Config struct {
	PulseCLIPath string `envconfig:"PULSE_CLI_PATH" default:"../../bin/pulse"`
}

type Steps struct {
	cfg    Config
	client *assured.Client
}

func NewSteps(cfg Config) *Steps {
	return &Steps{
		cfg:    cfg,
		client: assured.NewClientServe(),
	}
}
