// service holds the actual business logic (no HTTP logic)
package service

import (
	"context"
	"sync/atomic"
	"time"
)

// Processor is what handlers call; interface so we can mock in tests
type Processor interface {
	Process(ctx context.Context, input ProcessInput) (ProcessResult, error)
}

// ProcessInput is the payload we get from the HTTP body
type ProcessInput struct {
	Payload string `json:"payload"`
	Echo    bool   `json:"echo,omitempty"`
}

// ProcessResult is what we return as JSON
type ProcessResult struct {
	Result    string `json:"result"`
	Processed int64  `json:"processed_at_unix"`
}

// ProcessService does the work: prefix payload or echo it, set timestamp
type ProcessService struct {
	counter atomic.Int64
}

// NewProcessService returns a new processor
func NewProcessService() *ProcessService {
	return &ProcessService{}
}

// Process runs the logic
func (s *ProcessService) Process(ctx context.Context, input ProcessInput) (ProcessResult, error) {
	select {
	case <-ctx.Done():
		return ProcessResult{}, ctx.Err()
	default:
	}
	result := input.Payload
	if !input.Echo {
		result = "processed:" + result
	}
	ts := time.Now().Unix()
	s.counter.Add(1)
	return ProcessResult{
		Result:    result,
		Processed: ts,
	}, nil
}
