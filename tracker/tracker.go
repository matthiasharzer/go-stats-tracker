package tracker

type StatsTracker struct {
}

func NewStatsTracker() *StatsTracker {
	return &StatsTracker{}
}

func (s *StatsTracker) Submit(imageData []byte) error {

}
