package urlchecker

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// URLChecker actually connects to the specified url
// when performing its check
type URLChecker struct {
	client *http.Client
	logger *zap.Logger
}

// NewURLChecker instantiates a new URLChecker
func NewURLChecker(logger *zap.Logger) *URLChecker {
	if logger != nil {
		logger = logger.With(zap.String("module", "URLChecker"))
	}

	return &URLChecker{
		client: &http.Client{},
		logger: logger,
	}
}

// Check actually reads the specified url with a GET method, and returns
// up=true if the response code is 200, otherwise false
func (u *URLChecker) Check(ctx context.Context, url string) (up bool, duration time.Duration) {
	defer func() {
		if u.logger != nil {
			u.logger.With(
				zap.Bool("up", up),
				zap.Duration("duration", duration),
				zap.String("url", url),
			).Info("checked url")
		}
	}()

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, time.Since(start)
	}

	res, err := u.client.Do(req)
	if err != nil {
		return false, time.Since(start)
	}
	defer res.Body.Close()

	return res.StatusCode == http.StatusOK, time.Since(start)
}
