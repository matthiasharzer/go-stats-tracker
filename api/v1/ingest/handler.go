package ingest

import (
	"io"
	"net/http"

	"github.com/docker/go-units"
	"github.com/matthiasharzer/go-stats-tracker/logging"
	"github.com/matthiasharzer/go-stats-tracker/persistence"
	"github.com/matthiasharzer/go-stats-tracker/tracker"
	"github.com/matthiasharzer/go-stats-tracker/utils/timeutils"
)

const fileSizeLimit = 10 * units.MiB

func Handler(database persistence.Database, tracker *tracker.StatsTracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			logging.Debug("method not allowed")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.PathValue("userID")
		if userID == "" {
			logging.Debug("missing userID in request path")
			http.Error(w, "userID is required", http.StatusBadRequest)
			return
		}

		userContext, err := database.Lookup(userID)
		if err != nil {
			logging.Debug("failed to lookup user", "error", err, "userID", userID)
			http.Error(w, "failed to lookup user", http.StatusInternalServerError)
			return
		}
		if userContext == nil {
			logging.Debug("user not found", "userID", userID)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		if r.Body == nil {
			logging.Debug("missing request body in request", "userID", userID)
			http.Error(w, "missing request body", http.StatusBadRequest)
			return
		}
		defer func() {
			_ = r.Body.Close()
		}()

		limitedReader := http.MaxBytesReader(w, r.Body, int64(fileSizeLimit))
		fileContent, err := io.ReadAll(limitedReader)
		if err != nil {
			logging.Warn("failed to read request body", "error", err, "userID", userID)
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}

		today := timeutils.TodayDate()

		err = tracker.Submit(*userContext, fileContent, today)
		if err != nil {
			logging.Warn("failed to submit stats", "error", err, "userID", userID)
			http.Error(w, "failed to submit stats", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("OK"))
		if err != nil {
			logging.Error("failed to write response", "error", err, "userID", userID)
		}
	}
}
