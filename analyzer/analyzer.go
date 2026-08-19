package analyzer

type Analyzer interface {
	ExtractPlayerXP(imageData []byte) (int64, error)
	IsPlayerProfile(imageData []byte) (bool, error)
}
