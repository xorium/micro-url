package repo

import (
	"context"
	"fmt"
)

var (
	ErrNotExists    = fmt.Errorf("value for a such key doesn't exist")
	ErrAlreadySaved = fmt.Errorf("value for a such key has already saved")
)

type URLRepo interface {
	PutURL(ctx context.Context, shortID string, longURL string) error
	GetURL(ctx context.Context, shortID string) (longURL string, err error)
	// IncrementLatestCounterValue returns the last value of global counter for all
	// services named svcName incremented by delta.
	IncrementLatestCounterValue(ctx context.Context, svcName string, delta uint64) (newValue uint64, err error)
}
