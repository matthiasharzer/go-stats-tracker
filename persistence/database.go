package persistence

type SheetContext struct {
	GoogleRefreshToken  string
	TargetSpreadsheetID string
}

type Database interface {
	Lookup(accessKey string) (*SheetContext, error)
	Save(accessKey string, sheetContext SheetContext) error
}
