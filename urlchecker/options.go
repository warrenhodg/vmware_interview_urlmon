package urlchecker

import (
	"time"
)

type Options struct {
	// urls is the list of urls to check
	urls []string

	// checkPeriod is how frequently to check the urls
	checkPeriod time.Duration

	// workers is how many concurrent urls to check
	workers int

	// observe specifies whether prometheus observations
	// will occur (disabled for testing)
	observe bool
}

// DefaultOptions returns a baseline
// set of options that can be customised
func DefaultOptions() *Options {
	return &Options{
		urls:        []string{},
		checkPeriod: time.Second,
		workers:     1,
		observe:     true,
	}
}

// Urls specifies the set of urls to check
func (o *Options) Urls(v []string) *Options {
	o.urls = v
	return o
}

// CheckPeriod specifies how often the set of urls
// will be checked
func (o *Options) CheckPeriod(v time.Duration) *Options {
	o.checkPeriod = v
	return o
}

// Workers specifies how many workers will run
func (o *Options) Workers(v int) *Options {
	o.workers = v
	return o
}

// Observe specifies whether prometheus observations
// will occur (disabled for testing)
func (o *Options) Observe(v bool) *Options {
	o.observe = v
	return o
}
