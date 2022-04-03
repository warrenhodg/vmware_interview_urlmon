package urlchecker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Errors
var (
	ErrNilOptions    = fmt.Errorf("system options cannot be nil")
	ErrTooFewWorkers = fmt.Errorf("system workers must be a positive integer")
)

// System represents an entire checking system,
// including producer and consumers.
// It also manages the lifecycle of the
// goroutines.
type System struct {
	logger  *zap.Logger
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	options *Options
	c       chan string
	checker Checker
}

// New instantiates a System
func New(logger *zap.Logger, options *Options, checker Checker) (*System, error) {
	if logger != nil {
		logger = logger.With(zap.String("module", "urlchecker"))
	}

	if options == nil {
		return nil, ErrNilOptions
	}

	if options.workers <= 0 {
		return nil, ErrTooFewWorkers
	}

	return &System{
		logger:  logger,
		options: options,
		checker: checker,
	}, nil
}

// Run sets up the prometheus metrics, and starts
// the relevant producer and worker goroutines
// The goroutines can be terminated later by calling
// Shutdown
func (s *System) Run(ctx context.Context) {
	if s.options.observe {
		setupMetrics()
	}

	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.c = make(chan string, s.options.workers)

	for i := 0; i < s.options.workers; i++ {
		s.wg.Add(1)
		go s.work(ctx)
	}

	s.wg.Add(1)
	go s.produce(ctx)
}

// produce enqueues urls on the provided channel
// using the timing defined in the options until its
// context is done, then it closes the channel
func (s *System) produce(ctx context.Context) {
	if s.logger != nil {
		s.logger.Info("producer startup")
		defer s.logger.Info("producer shutdown")
	}
	defer s.wg.Done()

	ticker := time.NewTicker(s.options.checkPeriod)
	defer ticker.Stop()

	defer close(s.c)

	// Loop until context is done (cancel func called)
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if s.logger != nil {
				s.logger.Info("producer starting monitoring cycle")
			}

			for _, url := range s.options.urls {
				s.c <- url
			}
		}
	}
}

// work pulls urls off the channel, and checks
// them, saving the results to metrics.
// It terminates when the channel is closed.
// It is called as a goroutine by Run, and
// may have several instances running concurrently.
func (s *System) work(ctx context.Context) {
	if s.logger != nil {
		s.logger.Info("worker startup")
		defer s.logger.Info("worker shutdown")
	}
	defer s.wg.Done()

	// Loop will end when channel is closed
	for url := range s.c {
		up, duration := s.checker.Check(ctx, url)
		if s.options.observe {
			observe(url, up, duration)
		}
	}
}

// drainChan prevents a deadlock in the producer logic
func (s *System) drainChan() {
	for range s.c {
	}
}

// Shutdown gracefully by:
// 1. terminating the context, causing
// 2. the producer to terminate, causing
// 3. the channel to be closed, causing
// 4. the workers to stop
func (s *System) Shutdown() {
	s.cancel()
	s.drainChan()
	s.wg.Wait()
}
