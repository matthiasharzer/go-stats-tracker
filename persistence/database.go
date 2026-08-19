package persistence

type UserContext struct {
	GoogleRefreshToken  string
	TargetSpreadsheetID string
}

type Database interface {
	Lookup(key string) *UserContext
	Save(key string, userContext UserContext) error
}
