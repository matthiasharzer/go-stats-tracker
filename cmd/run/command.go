package run

import (
	"fmt"
	"net/http"
	"os"

	"github.com/matthiasharzer/go-stats-tracker/analyzer/tesseract"
	"github.com/matthiasharzer/go-stats-tracker/logging"
	"github.com/matthiasharzer/go-stats-tracker/persistence/sqlite"
	"github.com/matthiasharzer/go-stats-tracker/tracker"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func getGoogleOauthConfig() (*oauth2.Config, error) {
	redirectURL := os.Getenv("REDIRECT_URL")
	if redirectURL == "" {
		return nil, fmt.Errorf("REDIRECT_URL is not set")
	}

	clientID := os.Getenv("CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("CLIENT_ID is not set")
	}

	clientSecret := os.Getenv("CLIENT_SECRET")
	if clientSecret == "" {
		return nil, fmt.Errorf("CLIENT_SECRET is not set")
	}

	return &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/drive.file"},
		Endpoint:     google.Endpoint,
	}, nil
}

func getFilePickerCredentials() (apiKey string, appId string, err error) {
	apiKey = os.Getenv("FILE_PICKER_API_KEY")
	if apiKey == "" {
		return "", "", fmt.Errorf("FILE_PICKER_API_KEY is not set")
	}
	appId = os.Getenv("FILE_PICKER_APP_ID")
	if appId == "" {
		return "", "", fmt.Errorf("FILE_PICKER_APP_ID is not set")
	}

	return apiKey, appId, nil
}

var httpPort int
var httpHost string
var dbFile string

func init() {
	Command.Flags().IntVarP(&httpPort, "port", "p", 4000, "The HTTP server port to listen on")
	Command.Flags().StringVarP(&httpHost, "host", "", "", "The HTTP server host (default: all interfaces)")
	Command.Flags().StringVarP(&dbFile, "database-file", "d", "data/db.sqlite", "The database file to use")
}

var Command = &cobra.Command{
	Use:   "run",
	Short: "run the stats tracker server",
	RunE: func(cmd *cobra.Command, args []string) error {
		oauthConfig, err := getGoogleOauthConfig()
		if err != nil {
			return fmt.Errorf("failed to get oauth config: %w", err)
		}

		filePickerDeveloperKey, filePickerAppId, err := getFilePickerCredentials()
		if err != nil {
			return fmt.Errorf("failed to get file picker API key: %w", err)
		}

		database, err := sqlite.NewDatabase(dbFile)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}

		statsTracker := tracker.NewStatsTracker(tesseract.NewAnalyzer, database, *oauthConfig)

		mux := getMux(*oauthConfig, database, statsTracker, filePickerDeveloperKey, filePickerAppId)

		addr := fmt.Sprintf("%s:%d", httpHost, httpPort)
		logging.Info("starting go-stats-tracker server", "host", httpHost, "port", httpPort)
		err = http.ListenAndServe(
			addr,
			mux,
		)

		return fmt.Errorf("failed to start server: %w", err)
	},
}
