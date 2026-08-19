package tracker

import (
	"fmt"

	"github.com/matthiasharzer/go-stats-tracker/analyzer"
)

type StatsTracker struct {
	analyzer analyzer.Analyzer
}

func NewStatsTracker(analyzer analyzer.Analyzer) *StatsTracker {
	return &StatsTracker{
		analyzer: analyzer,
	}
}

func (s *StatsTracker) Submit(imageData []byte) error {

	stats, err := s.analyzer.ExtractPlayerStats(imageData)
	if err != nil {
		return fmt.Errorf("failed to extract player XP: %w", err)
	}

	fmt.Printf("Level %d\n", stats.Level)
	fmt.Printf("Total XP: %d\n", stats.TotalLevelXP)
	fmt.Printf("Gained XP: %d\n", stats.GainedLevelXP)

	return nil
}
