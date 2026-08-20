package register

import (
	"net/http"

	"github.com/matthiasharzer/go-stats-tracker/api"
	"github.com/matthiasharzer/go-stats-tracker/logging"
	"golang.org/x/oauth2"
)

func Handler(sharedState *api.SharedState, oauth oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sheetID := r.URL.Query().Get("sheet_id")
		if sheetID == "" {
			http.Error(w, "Missing sheet_id", http.StatusBadRequest)
			return
		}

		stateID := sharedState.NewStateID(sheetID)
		authURL := oauth.AuthCodeURL(stateID, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

		logging.Info("requesting access", "sheet_id", sheetID)

		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	}
}
