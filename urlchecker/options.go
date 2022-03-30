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
}

func DefaultOptions() *Options {
	return &Options{
		urls:        []string{},
		checkPeriod: time.Second,
		workers:     1,
	}
}

func (o *Options) Urls(v []string) *Options {
	o.urls = v
	return o
}

func (o *Options) CheckPeriod(v time.Duration) *Options {
	o.checkPeriod = v
	return o
}

func (o *Options) Workers(v int) *Options {
	o.workers = v
	return o
}
