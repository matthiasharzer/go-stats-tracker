package run

import (
	"net/http"

	"github.com/matthiasharzer/go-stats-tracker/api/v1/auth"
	"github.com/matthiasharzer/go-stats-tracker/api/v1/auth/callback"
	"github.com/matthiasharzer/go-stats-tracker/api/v1/auth/link"
	"github.com/matthiasharzer/go-stats-tracker/api/v1/auth/register"
	"github.com/matthiasharzer/go-stats-tracker/api/v1/ingest"
	"github.com/matthiasharzer/go-stats-tracker/persistence"
	"github.com/matthiasharzer/go-stats-tracker/tracker"
	"golang.org/x/oauth2"
)

func getMux(oauth oauth2.Config, database persistence.Database, tracker *tracker.StatsTracker, filePickerAPIKey string, filePickerAppID string) *http.ServeMux {
	sharedState := auth.NewSharedState()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /api/v1/register", register.Handler(sharedState, oauth))
	mux.HandleFunc("GET /api/v1/callback", callback.Handler(sharedState, oauth, filePickerAPIKey, filePickerAppID))
	mux.HandleFunc("POST /api/v1/link", link.Handler(sharedState, database))
	mux.HandleFunc("POST /api/v1/ingest/{sheetID}", ingest.Handler(database, tracker))

	return mux
}
