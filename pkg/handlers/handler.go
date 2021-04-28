package handlers

import (
	"net/http"
	"sync/atomic"

	"github.com/gorilla/mux"
)

func Router(isReady *atomic.Value) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthz)
	r.HandleFunc("/readyz", readyz(isReady))
	return r
}

// healthz is a liveness probe.
func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// readyz is a readiness probe.
func readyz(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady != nil {
			v := isReady.Load()
			if v != nil {
				f := v.(bool)
				if f {
					w.WriteHeader(http.StatusOK)
					return
				}
			}

		}
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
}
