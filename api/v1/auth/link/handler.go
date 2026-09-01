package link

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/docker/go-units"
	"github.com/matthiasharzer/go-stats-tracker/api/v1/auth"
	"github.com/matthiasharzer/go-stats-tracker/persistence"
	"github.com/matthiasharzer/go-stats-tracker/utils/funcutils"
)

const maxBodySize = 1 * units.MiB

func Handler(sharedState *auth.SharedState, database persistence.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer funcutils.LogError(r.Body.Close, "failed to request body")
		limitedReader := io.LimitReader(r.Body, maxBodySize)
		body, err := io.ReadAll(limitedReader)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		var requestBody RequestBody
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		if requestBody.State == "" || requestBody.SheetID == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		authFlowState := sharedState.PopState(requestBody.State)
		if authFlowState == nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = database.Save(authFlowState.UserAccessKey, persistence.SheetContext{
			GoogleRefreshToken:  authFlowState.RefreshToken,
			TargetSpreadsheetID: requestBody.SheetID,
		})
		if err != nil {
			http.Error(w, "Failed to save sheet context", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
