package persistence

type UserContext struct {
	GoogleRefreshToken  string
	TargetSpreadsheetID string
}

type Database interface {
	Lookup(userID string) (*UserContext, error)
	Save(userID string, userContext UserContext) error
	Exists(userID string) (bool, error)
}
