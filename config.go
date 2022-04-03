package main

import (
	"time"

	"github.com/alecthomas/kong"
)

// Config defines application level configs that
// then are passed to various modules' initializations
// that make up the application
type Config struct {
	ListenAddr string `name:"listen-addr" help:"Address on which to listen" short:"l" default:":8080"`

	// URLs is the list of urls to check
	URLs []string `name:"urls" help:"List of urls to check" short:"u" default:"http://httpstat.us/503,http://httpstat.us/200" required:""`

	// CheckPeriod is how frequently to check the urls
	CheckPeriod time.Duration `name:"check-period" help:"How often to perform checks" short:"p" default:"5s"`

	// Workers is how many concurrent urls to check
	Workers int `name:"workers" help:"How many checks can be performed concurrently" short:"w" default:"2"`

	// ShutdownDuration is the time provided to wait for existing connections to close
	ShutdownDuration time.Duration `name:"shutdown-duration" help:"Time allowed for existing connections to close when shutting down" default:"1s"`

	// ShutdownReadyFailDuration is the duration for which the server will fail its readiness
	// checks (forcing upstream traffic redirection away) prior to actually shutting down
	ShutdownReadyFailDuration time.Duration `name:"shutdown-ready-fail-duration" help:"Fail readiness check for this long before actually shutting down" default:"1s"`
}

// GetConfig returns the application configuration
func GetConfig() (*Config, error) {
	cfg := &Config{}
	ctx := kong.Parse(cfg)
	return cfg, ctx.Error
}
