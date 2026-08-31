package tracker

import (
	"fmt"
	"time"

	"github.com/matthiasharzer/go-stats-tracker/analyzer"
	"github.com/matthiasharzer/go-stats-tracker/logging"
	"github.com/matthiasharzer/go-stats-tracker/persistence"
	"github.com/matthiasharzer/go-stats-tracker/tracker/spreadsheet"
	"golang.org/x/oauth2"
)

type CreateAnalyzerFunc = func(imageData []byte) (analyzer.Analyzer, error)

type StatsTracker struct {
	createAnalyzer CreateAnalyzerFunc
	database       persistence.Database
	sheetService   *spreadsheet.Service
}

func NewStatsTracker(createAnalyzer CreateAnalyzerFunc, database persistence.Database, oauth oauth2.Config) *StatsTracker {
	return &StatsTracker{
		createAnalyzer: createAnalyzer,
		database:       database,
		sheetService:   spreadsheet.NewService(oauth),
	}
}

func (s *StatsTracker) Submit(userContext persistence.SheetContext, imageData []byte, date time.Time) error {
	statAnalyzer, err := s.createAnalyzer(imageData)
	if err != nil {
		logging.Error("failed to create analyzer", "err", err)
		return fmt.Errorf("unable to create analyzer: %w", err)
	}

	isPlayerPage, err := statAnalyzer.IsPlayerPage()
	if err != nil {
		return fmt.Errorf("unable to check if image is a player page: %w", err)
	}
	if !isPlayerPage {
		return fmt.Errorf("image is not a player page")
	}

	stats, err := statAnalyzer.ExtractPlayerStats()
	if err != nil {
		logging.Error("failed to extract player stats", "err", err)
		return fmt.Errorf("failed to extract player stats: %w", err)
	}

	logging.Info("extracted player stats", "level", stats.Level, "gainedLevelXP", stats.GainedLevelXP, "totalLevelXP", stats.TotalLevelXP, "date", date.Format("2006-01-02"))

	err = s.sheetService.AppendStats(userContext, stats, date)
	if err != nil {
		return fmt.Errorf("failed to append stats to spreadsheet: %w", err)
	}
	return nil
}
