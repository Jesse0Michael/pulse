package service

import "context"

type Pulser struct {
}

func NewPulser() *Pulser {
	return &Pulser{}
}

func (p *Pulser) Summary(ctx context.Context, req SummaryRequest) (*Summary, error) {
	return &Summary{
		Summary: "test",
	}, nil
}
