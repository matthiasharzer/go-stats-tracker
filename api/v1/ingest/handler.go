package ingest

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/docker/go-units"
	"github.com/matthiasharzer/go-stats-tracker/logging"
	"github.com/matthiasharzer/go-stats-tracker/persistence"
	"github.com/matthiasharzer/go-stats-tracker/tracker"
	"github.com/matthiasharzer/go-stats-tracker/utils/timeutils"
)

const fileSizeLimit = 10 * units.MiB

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

func Handler(database persistence.Database, tracker *tracker.StatsTracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sheetID := r.PathValue("sheetID")
		if sheetID == "" {
			logging.Debug("missing sheetID in request path")
			http.Error(w, "sheetID is required", http.StatusBadRequest)
			return
		}

		bearerToken, err := extractBearerToken(r)
		if err != nil {
			logging.Debug("failed to extract bearer token", "error", err, "sheet_id", sheetID)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sheetContext, err := database.Lookup(bearerToken)
		if err != nil {
			logging.Debug("failed to lookup sheet context", "error", err, "sheet_id", sheetID)
			http.Error(w, "failed to lookup sheet context", http.StatusInternalServerError)
			return
		}
		if sheetContext == nil {
			logging.Debug("sheet not found", "userID", sheetID)
			http.Error(w, "sheet not found", http.StatusNotFound)
			return
		}
		if sheetContext.TargetSpreadsheetID != sheetID {
			logging.Debug("sheet id mismatch", "sheet_id", sheetID, "target_spreadsheet_id", sheetContext.TargetSpreadsheetID)
			http.Error(w, "sheet id mismatch", http.StatusUnauthorized)
			return
		}

		if r.Body == nil {
			logging.Debug("missing request body in request", "sheet_id", sheetID)
			http.Error(w, "missing request body", http.StatusBadRequest)
			return
		}
		defer func() {
			_ = r.Body.Close()
		}()

		limitedReader := http.MaxBytesReader(w, r.Body, int64(fileSizeLimit))
		fileContent, err := io.ReadAll(limitedReader)
		if err != nil {
			logging.Warn("failed to read request body", "error", err, "sheet_id", sheetID)
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}

		today := timeutils.TodayDate()

		err = tracker.Submit(*sheetContext, fileContent, today)
		if err != nil {
			logging.Warn("failed to submit stats", "error", err, "sheet_id", sheetID)
			http.Error(w, "failed to submit stats", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("OK"))
		if err != nil {
			logging.Error("failed to write response", "error", err, "userID", sheetID)
		}
	}
}
