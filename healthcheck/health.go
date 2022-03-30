package healthcheck

import (
	"net/http"
	"sync/atomic"

	"go.uber.org/zap"
)

var (
	// ready is an atomic int used to reflect
	// whether the traffic wishes to receive traffic
	// Non-zero values implies no
	ready uint32

	logger *zap.Logger
)

// Init sets up the health endpoint
func Init(l *zap.Logger) {
	logger = l.With(zap.String("module", "health"))
	http.DefaultServeMux.HandleFunc("/ready", handleReady)
	http.DefaultServeMux.HandleFunc("/live", handleLive)
}

// handleLive always returns 200 OK
func handleLive(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

// handleReady returns whether the app wishes to
// receive traffic
// 200 OK to indicated yes, and
// 500 to indicate no
func handleReady(w http.ResponseWriter, r *http.Request) {
	h := atomic.LoadUint32(&ready)

	if h == 0 {
		_, _ = w.Write([]byte("ok"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// SetReady explicitly sets the health of the system,
// changing what the /health endpoint will return.
// It is used at shutdown to fail the healthchecks
// prior to actually refusing new connections to give the
// load balancer time to reroute traffic
func SetReady(value bool) {
	intValue := uint32(0)
	if !value {
		intValue = 1
	}
	atomic.StoreUint32(&ready, intValue)
	logger.With(zap.Bool("value", value)).Info("/ready result changed")
}
