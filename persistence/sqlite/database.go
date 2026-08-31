package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/matthiasharzer/go-stats-tracker/persistence"
	"github.com/matthiasharzer/go-stats-tracker/utils/funcutils"
	_ "modernc.org/sqlite"
)

type Database struct {
	file string
}

func NewDatabase(dbFile string) (persistence.Database, error) {
	db := &Database{
		file: dbFile,
	}
	err := db.createTable()
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	return db, nil
}

func (d *Database) connect() (*sql.DB, error) {
	return sql.Open("sqlite", d.file)
}

func (d *Database) createTable() error {
	db, err := d.connect()
	if err != nil {
		return err
	}
	defer funcutils.LogError(db.Close, "failed to close database connection")

	createTableSQL := `CREATE TABLE IF NOT EXISTS sheet_context (
		accessKey TEXT PRIMARY KEY,
		googleRefreshToken TEXT NOT NULL,
		targetSpreadsheetID TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) Lookup(accessKey string) (*persistence.SheetContext, error) {
	db, err := d.connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer funcutils.LogError(db.Close, "failed to close database connection")

	selectSQL := `SELECT googleRefreshToken, targetSpreadsheetID FROM sheet_context WHERE accessKey = ?`
	row := db.QueryRow(selectSQL, accessKey)

	var googleRefreshToken, targetSpreadsheetID string
	err = row.Scan(&googleRefreshToken, &targetSpreadsheetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &persistence.SheetContext{
		GoogleRefreshToken:  googleRefreshToken,
		TargetSpreadsheetID: targetSpreadsheetID,
	}, nil
}
func (d *Database) Save(accessKey string, sheetContext persistence.SheetContext) error {
	db, err := d.connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer funcutils.LogError(db.Close, "failed to close database connection")

	insertSQL := `INSERT INTO sheet_context (accessKey, googleRefreshToken, targetSpreadsheetID) VALUES (?, ?, ?)`
	_, err = db.Exec(insertSQL, accessKey, sheetContext.GoogleRefreshToken, sheetContext.TargetSpreadsheetID)
	if err != nil {
		return fmt.Errorf("failed to insert into sheet_context: %w", err)
	}

	return nil
}

func (d *Database) Exists(accessKey string) (bool, error) {
	db, err := d.connect()
	if err != nil {
		return false, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer funcutils.LogError(db.Close, "failed to close database connection")

	selectSQL := `SELECT COUNT(*) FROM sheet_context WHERE accessKey = ?`
	row := db.QueryRow(selectSQL, accessKey)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to scan row: %w", err)
	}

	return count > 0, nil
}
