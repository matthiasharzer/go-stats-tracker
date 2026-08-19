package tracker

import (
	"fmt"

	"github.com/matthiasharzer/go-stats-tracker/analyzer"
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

func (s *StatsTracker) Submit(userID string, imageData []byte) error {

	stats, err := s.analyzer.ExtractPlayerStats(imageData)
	if err != nil {
		return fmt.Errorf("failed to extract player XP: %w", err)
	}

	userContext, err := s.database.Lookup(userID)
	if err != nil {
		return fmt.Errorf("failed to lookup user: %w", err)
	}
	if userContext == nil {
		return fmt.Errorf("user not found")
	}

	s.sheetService.AppendStats(*userContext, stats)

	fmt.Printf("Level %d\n", stats.Level)
	fmt.Printf("Total XP: %d\n", stats.TotalLevelXP)
	fmt.Printf("Gained XP: %d\n", stats.GainedLevelXP)

	return nil
}
