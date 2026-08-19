package analyzer

type PlayerStats struct {
	Level         int
	TotalLevelXP  int64
	GainedLevelXP int64
}

type Analyzer interface {
	ExtractPlayerStats(imageData []byte) (PlayerStats, error)
}
