package urlchecker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/warrenhodg/urlmon/urlchecker/mock"
)

func errorsMatch(err, pattern error) bool {
	if pattern == nil {
		return err == nil
	}

	return errors.Is(err, pattern)
}

func TestSystem(t *testing.T) {
	testCases := []struct {
		description   string
		urls          []string
		checkPeriod   time.Duration
		checkDuration time.Duration
		totalTime     time.Duration
		workers       int
		count         int
		err           error
	}{
		{
			description:   "1 url every 10ms for 55ms",
			urls:          []string{"http://httpstat.us/200"},
			checkPeriod:   time.Millisecond * 10,
			checkDuration: time.Millisecond * 2,
			totalTime:     time.Millisecond * 55,
			workers:       1,
			count:         5,
			err:           nil,
		},
		{
			description:   "too many workers",
			urls:          []string{"http://httpstat.us/200"},
			checkPeriod:   time.Millisecond * 10,
			checkDuration: time.Millisecond * 2,
			totalTime:     time.Millisecond * 55,
			workers:       5,
			count:         5,
			err:           nil,
		},
		{
			description:   "too many workers",
			urls:          []string{"http://httpstat.us/200"},
			checkPeriod:   time.Millisecond * 10,
			checkDuration: time.Millisecond * 2,
			totalTime:     time.Millisecond * 55,
			workers:       0,
			count:         0,
			err:           ErrTooFewWorkers,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			m := &mock.URLChecker{
				Wait: true,
				DefaultResponse: mock.Response{
					Duration: testCase.checkDuration,
				},
			}

			options := DefaultOptions().
				Urls(testCase.urls).
				Workers(testCase.workers).
				CheckPeriod(testCase.checkPeriod).
				Observe(false)
			s, err := New(nil, options, m)
			if !errorsMatch(err, testCase.err) {
				t.Fatalf("unexpected error. expected: %v. received: %v", testCase.err, err)
				return
			}

			if s == nil {
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), testCase.totalTime)
			defer cancel()

			s.Run(ctx)
			<-ctx.Done()

			requests := m.RequestCount()
			if requests != testCase.count {
				t.Fatalf("incorrect number of tests performed: %d vs %d", requests, testCase.count)
				return
			}
		})
	}
}
