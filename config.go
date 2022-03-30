package main

import (
	"time"
)

// Config defines application level configs that
// then are passed to various modules' initializations
// that make up the application
type Config struct {
	listenAddr string

	// urls is the list of urls to check
	urls []string

	// checkPeriod is how frequently to check the urls
	checkPeriod time.Duration

	// workers is how many concurrent urls to check
	workers int

	// shutdownDuration is the time provided to wait for existing connections to close
	shutdownDuration time.Duration

	// shutdownReadyFailDuration is the duration for which the server will fail its readiness
	// checks (forcing upstream traffic redirection away) prior to actually shutting down
	shutdownReadyFailDuration time.Duration
}

// GetConfig returns the application configuration
func GetConfig() *Config {
	return &Config{
		listenAddr: ":8080",
		urls: []string{
			"http://httpstat.us/503",
			"http://httpstat.us/200",
		},
		checkPeriod:               time.Second * 5,
		workers:                   2,
		shutdownDuration:          time.Second * 1,
		shutdownReadyFailDuration: time.Second * 1,
	}
}
