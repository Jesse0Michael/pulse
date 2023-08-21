package service

import "context"

type Config struct {
}

type Pulser struct {
	cfg Config
}

func NewPulser(cfg Config) *Pulser {
	return &Pulser{
		cfg: cfg,
	}
}

func (p *Pulser) Summary(ctx context.Context, req SummaryRequest) (*Summary, error) {
	return &Summary{
		Summary: "test",
	}, nil
}
