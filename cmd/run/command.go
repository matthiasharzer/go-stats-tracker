package run

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/matthiasharzer/go-stats-tracker/analyzer/tesseract"
)

var Command = &cobra.Command{
	Use:   "run",
	Short: "run the stats tracker server",
	RunE: func(cmd *cobra.Command, args []string) error {
		analyzer, err := tesseract.NewAnalyzer()
		if err != nil {
			return fmt.Errorf("failed to create analyzer: %w", err)
		}

		file, err := os.Open("test.png")
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		imageBytes, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		_, err = analyzer.ExtractPlayerXP(imageBytes)
		if err != nil {
			return fmt.Errorf("failed to extract player XP: %w", err)
		}

		return nil
	},
}
