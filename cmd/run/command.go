package run

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/matthiasharzer/go-stats-tracker/analyzer/tesseract"
	"github.com/matthiasharzer/go-stats-tracker/logging"
	"github.com/matthiasharzer/go-stats-tracker/persistence/inmemory"
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
		// This MUST match exactly what you added in the GCP Console
		Scopes:   []string{"https://www.googleapis.com/auth/spreadsheets"},
		Endpoint: google.Endpoint,
	}, nil
}

var httpPort int
var httpHost string

func init() {
	Command.Flags().IntVarP(&httpPort, "port", "p", 4000, "The HTTP server port to listen on")
	Command.Flags().StringVarP(&httpHost, "host", "", "", "The HTTP server host (default: all interfaces)")
}

var Command = &cobra.Command{
	Use:   "run",
	Short: "run the stats tracker server",
	RunE: func(cmd *cobra.Command, args []string) error {
		analyzer, err := tesseract.NewAnalyzer()
		if err != nil {
			return fmt.Errorf("failed to create analyzer: %w", err)
		}

		statsTracker := tracker.NewStatsTracker(analyzer)

		oauthConfig, err := getGoogleOauthConfig()
		if err != nil {
			return fmt.Errorf("failed to get oauth config: %w", err)
		}

		database := inmemory.NewDatabase()

		mux := getMux(*oauthConfig, database)

		file, err := os.Open("test.png")
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		imageBytes, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		err = statsTracker.Submit(imageBytes)
		if err != nil {
			return fmt.Errorf("failed to submit image: %w", err)
		}

		addr := fmt.Sprintf("%s:%d", httpHost, httpPort)
		logging.Info("starting go-stats-tracker server", "host", httpHost, "port", httpPort)
		err = http.ListenAndServe(
			addr,
			mux,
		)

		return fmt.Errorf("failed to start server: %w", err)
	},
}
