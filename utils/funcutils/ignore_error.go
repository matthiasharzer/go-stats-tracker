package funcutils

import "github.com/matthiasharzer/go-stats-tracker/logging"

func LogError(fn func() error, message string) {
	err := fn()
	if err != nil {
		logging.Error(message, "error", err)
	}
}
