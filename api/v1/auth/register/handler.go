package register

import (
	"net/http"

	"github.com/matthiasharzer/go-stats-tracker/api/v1/auth"
	"github.com/matthiasharzer/go-stats-tracker/logging"
	"github.com/matthiasharzer/go-stats-tracker/utils/stringutils"
	"golang.org/x/oauth2"
)

func Handler(sharedState *auth.SharedState, oauth oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessKey := stringutils.RandomString(48)

		stateID := sharedState.NewStateID(accessKey)
		authURL := oauth.AuthCodeURL(stateID, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

		logging.Info("requesting access", "state_id", stateID)

		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	}
}
