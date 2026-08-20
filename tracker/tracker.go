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

type StatsTracker struct {
	analyzer     analyzer.Analyzer
	database     persistence.Database
	sheetService *spreadsheet.Service
}

func NewStatsTracker(analyzer analyzer.Analyzer, database persistence.Database, oauth oauth2.Config) *StatsTracker {
	return &StatsTracker{
		analyzer:     analyzer,
		database:     database,
		sheetService: spreadsheet.NewService(oauth),
	}
}

func (s *StatsTracker) Submit(userID string, imageData []byte, date time.Time) error {
	stats, err := s.analyzer.ExtractPlayerStats(imageData)
	if err != nil {
		return fmt.Errorf("failed to extract player XP: %w", err)
	}

	logging.Info("extracted player stats", "userID", userID, "level", stats.Level, "gainedLevelXP", stats.GainedLevelXP, "totalLevelXP", stats.TotalLevelXP, "date", date.Format("2006-01-02"))

	userContext, err := s.database.Lookup(userID)
	if err != nil {
		return fmt.Errorf("failed to lookup user: %w", err)
	}
	if userContext == nil {
		return fmt.Errorf("user not found")
	}

	err = s.sheetService.AppendStats(*userContext, stats, date)
	if err != nil {
		return fmt.Errorf("failed to append stats to spreadsheet: %w", err)
	}
	return nil
}
