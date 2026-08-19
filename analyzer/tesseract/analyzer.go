package tesseract

import (
	"fmt"

	"github.com/otiai10/gosseract/v2"

	"github.com/matthiasharzer/go-stats-tracker/analyzer"
	"github.com/matthiasharzer/go-stats-tracker/logging"
)

type Analyzer struct {
}

func newClient() (*gosseract.Client, func()) {
	client := gosseract.NewClient()
	return client, func() {
		err := client.Close()
		if err != nil {
			logging.Warn("Failed to close tesseract client", "err", err)
		}
	}
}

func NewAnalyzer() (analyzer.Analyzer, error) {
	return &Analyzer{}, nil
}

func (a *Analyzer) ExtractPlayerXP(imageData []byte) (int64, error) {
	client, closeClient := newClient()
	defer closeClient()

	err := client.SetImageFromBytes(imageData)
	if err != nil {
		return 0, fmt.Errorf("failed to set image from bytes: %w", err)
	}

	text, err := client.Text()
	if err != nil {
		return 0, fmt.Errorf("failed to extract text: %w", err)
	}

	fmt.Print("Extracted text:\n", text)

	return 0, nil
}

func (a *Analyzer) IsPlayerProfile(imageData []byte) (bool, error) {
	return true, nil
}
