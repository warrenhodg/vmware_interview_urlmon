package urlchecker

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

type UrlChecker struct {
	logger  *zap.Logger
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	options *Options
	c       chan string
}

func New(logger *zap.Logger, options *Options) *UrlChecker {
	return &UrlChecker{
		logger:  logger.With(zap.String("module", "urlchecker")),
		options: options,
	}
}

func (u *UrlChecker) Run(ctx context.Context) {
	setupMetrics()

	ctx, cancel := context.WithCancel(ctx)
	u.cancel = cancel
	u.c = make(chan string, u.options.workers)

	for i := 0; i < u.options.workers; i++ {
		u.wg.Add(1)
		go u.work(ctx)
	}

	u.wg.Add(1)
	go u.produce(ctx)
}

func (u *UrlChecker) produce(ctx context.Context) {
	u.logger.Info("producer startup")
	defer u.logger.Info("producer shutdown")
	defer u.wg.Done()

	ticker := time.NewTicker(u.options.checkPeriod)
	defer ticker.Stop()

	defer close(u.c)

	// Loop until context is done (cancel func called)
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			u.logger.Info("producer starting monitoring cycle")

			for _, url := range u.options.urls {
				u.c <- url
			}
		}
	}
}

func (u *UrlChecker) check(url string) {
	up := false
	start := time.Now()
	defer func() {
		observe(url, up, time.Since(start))
	}()

	u.logger.With(zap.String("url", url)).Info("checking url")
	time.Sleep(time.Millisecond * 300)
}

func (u *UrlChecker) work(ctx context.Context) {
	u.logger.Info("worker startup")
	defer u.logger.Info("worker shutdown")
	defer u.wg.Done()

	// Loop will end when channel is closed
	for url := range u.c {
		u.check(url)
	}
}

// drainChan prevents a deadlock in the producer logic
func (u *UrlChecker) drainChan() {
	for range u.c {
	}
}

// Shutdown gracefully
func (u *UrlChecker) Shutdown() {
	u.cancel()
	u.drainChan()
	u.wg.Wait()
}
