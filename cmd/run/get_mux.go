package run

import (
	"net/http"

	"github.com/matthiasharzer/go-stats-tracker/api/v1/auth"
	"github.com/matthiasharzer/go-stats-tracker/api/v1/auth/callback"
	"github.com/matthiasharzer/go-stats-tracker/api/v1/auth/link"
	"github.com/matthiasharzer/go-stats-tracker/api/v1/auth/register"
	"github.com/matthiasharzer/go-stats-tracker/api/v1/submit"
	"github.com/matthiasharzer/go-stats-tracker/persistence"
	"github.com/matthiasharzer/go-stats-tracker/tracker"
	"golang.org/x/oauth2"
)

func getMux(oauth oauth2.Config, database persistence.Database, tracker *tracker.StatsTracker, pickerAPIKey string, appID string) *http.ServeMux {
	sharedState := auth.NewSharedState()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /api/v1/auth/register", register.Handler(sharedState, oauth))
	mux.HandleFunc("GET /api/v1/auth/callback", callback.Handler(sharedState, oauth, pickerAPIKey, appID))
	mux.HandleFunc("POST /api/v1/auth/link", link.Handler(sharedState, database))
	mux.HandleFunc("POST /api/v1/submit/{sheetID}", submit.Handler(database, tracker))

	return mux
}
