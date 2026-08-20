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

func (s *StatsTracker) Submit(userContext persistence.UserContext, imageData []byte, date time.Time) error {
	stats, err := s.analyzer.ExtractPlayerStats(imageData)
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
