package callback

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/matthiasharzer/go-stats-tracker/api"
	"github.com/matthiasharzer/go-stats-tracker/logging"
	"github.com/matthiasharzer/go-stats-tracker/persistence"
	"golang.org/x/oauth2"
)

func Handler(sharedState *api.SharedState, oauth oauth2.Config, database persistence.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		if state == "" {
			http.Error(w, "State is empty", http.StatusBadRequest)
			return
		}

		sheetID := sharedState.PopStateID(state)
		if sheetID == "" {
			http.Error(w, "State is invalid", http.StatusBadRequest)
			return
		}

		code := r.FormValue("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		token, err := oauth.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if token.RefreshToken == "" {
			http.Error(w, "No refresh token returned. Try revoking app permissions and trying again.", http.StatusInternalServerError)
			return
		}

		logging.Info("received callback", "sheet_id", sheetID)

		userID := uuid.NewString()

		err = database.Save(userID, persistence.UserContext{
			GoogleRefreshToken:  token.RefreshToken,
			TargetSpreadsheetID: sheetID,
		})
		if err != nil {
			logging.Error("failed to save tokens to database", "err", err)
			http.Error(w, "Failed to save tokens to database", http.StatusInternalServerError)
			return
		}

		logging.Info("save user context", "user_id", userID, "sheet_id", sheetID)

		response := fmt.Sprintf("Success! Your User ID is: %s\n\nUse this ID to authenticate when uploading screenshots", userID)
		_, err = w.Write([]byte(response))
		if err != nil {
			logging.Error("failed to write response", "err", err)
		}
	}
}
