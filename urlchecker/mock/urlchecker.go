package mock

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Response stores the expected response
// for a single url's check
type Response struct {
	Up       bool
	Duration time.Duration
}

// Request is used to keep track of requests made
type Request struct {
	Time time.Time
	URL  string
}

// URLChecker implements the Checker interface
// and should be used for tests
type URLChecker struct {
	Logger *zap.Logger

	// Wait being set to true causes the Check function
	// call to actually take the specified duration
	Wait bool

	// DefaultResponse is used for when the url is
	// not present in the Responses map
	DefaultResponse Response

	// Responses stores expected responses to various
	// url requests
	Responses map[string]Response

	// Requests
	Requests []*Request

	// Mutex allows -race tests to pass
	Mutex sync.Mutex
}

// Check simulates checking a url, and can optionally
// wait the specified duration. If the response behaviour
// is not present in the Responses map, it uses the
// DefaultResponse
func (u *URLChecker) Check(ctx context.Context, url string) (up bool, duration time.Duration) {
	var r Response
	found := false

	req := &Request{
		Time: time.Now(),
		URL:  url,
	}

	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	u.Requests = append(u.Requests, req)

	if u.Logger != nil {
		u.Logger.With(zap.String("url", url)).Info("Checking")
	}

	if u.Responses != nil {
		r, found = u.Responses[url]
	}

	if !found {
		r = u.DefaultResponse
	}

	if u.Wait {
		time.Sleep(r.Duration)
	}

	return r.Up, r.Duration
}

// RequestCount is a goroutine-safe
// way to determine how many requests have
// been made
func (u *URLChecker) RequestCount() int {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	return len(u.Requests)

}
