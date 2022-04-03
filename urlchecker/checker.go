package urlchecker

import (
	"context"
	"time"
)

// Checker interface represents something that
// can check a url, returning whether it is up
// or down, and how long it took
type Checker interface {
	Check(ctx context.Context, url string) (up bool, duration time.Duration)
}
