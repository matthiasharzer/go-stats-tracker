package analyzer

type PlayerStats struct {
	Level         int
	TotalLevelXP  int64
	GainedLevelXP int64
}

type Analyzer interface {
	ExtractPlayerStats() (PlayerStats, error)
	IsPlayerPage() (bool, error)
}
