package run

import (
	"net/http"

	"github.com/matthiasharzer/go-stats-tracker/api"
	"github.com/matthiasharzer/go-stats-tracker/api/v1/callback"
	"github.com/matthiasharzer/go-stats-tracker/api/v1/ingest"
	"github.com/matthiasharzer/go-stats-tracker/api/v1/register"
	"github.com/matthiasharzer/go-stats-tracker/persistence"
	"github.com/matthiasharzer/go-stats-tracker/tracker"
	"golang.org/x/oauth2"
)

func getMux(oauth oauth2.Config, database persistence.Database, tracker *tracker.StatsTracker) *http.ServeMux {
	sharedState := api.NewSharedState()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /api/v1/register", register.Handler(sharedState, oauth))
	mux.HandleFunc("GET /api/v1/callback", callback.Handler(sharedState, oauth, database))
	mux.HandleFunc("POST /api/v1/ingest/{userID}", ingest.Handler(tracker))

	return mux
}
