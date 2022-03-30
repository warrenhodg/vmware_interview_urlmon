package httpserver

import (
	"time"
)

type Options struct {
	listenAddr       string
	shutdownDuration time.Duration
}

func DefaultOptions() *Options {
	return &Options{
		listenAddr:       ":8080",
		shutdownDuration: time.Second * 5,
	}
}

func (o *Options) ListenAddr(v string) *Options {
	o.listenAddr = v
	return o
}
