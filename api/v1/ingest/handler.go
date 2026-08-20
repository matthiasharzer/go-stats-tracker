package ingest

import (
	"io"
	"net/http"

	"github.com/docker/go-units"
	"github.com/matthiasharzer/go-stats-tracker/tracker"
	"github.com/matthiasharzer/go-stats-tracker/utils/timeutils"
)

const fileSizeLimit = 10 * units.MiB

func Handler(tracker *tracker.StatsTracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.PathValue("userID")
		if userID == "" {
			http.Error(w, "userID is required", http.StatusBadRequest)
			return
		}

		if r.Body == nil {
			http.Error(w, "missing request body", http.StatusBadRequest)
			return
		}
		defer func() {
			_ = r.Body.Close()
		}()

		limitedReader := http.MaxBytesReader(w, r.Body, int64(fileSizeLimit))
		fileContent, err := io.ReadAll(limitedReader)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}

		today := timeutils.TodayDate()

		err = tracker.Submit(userID, fileContent, today)
		if err != nil {
			http.Error(w, "failed to submit stats", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
