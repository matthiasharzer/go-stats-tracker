package sqlite

import (
	"database/sql"
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
	defer funcutils.IgnoreError(db.Close)

	createTableSQL := `CREATE TABLE IF NOT EXISTS user_context (
		userID TEXT PRIMARY KEY,
		googleRefreshToken TEXT NOT NULL,
		targetSpreadsheetID TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) Lookup(userID string) (*persistence.UserContext, error) {
	db, err := d.connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer funcutils.IgnoreError(db.Close)

	selectSQL := `SELECT googleRefreshToken, targetSpreadsheetID FROM user_context WHERE userID = ?`
	row := db.QueryRow(selectSQL, userID)

	var googleRefreshToken, targetSpreadsheetID string
	err = row.Scan(&googleRefreshToken, &targetSpreadsheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &persistence.UserContext{
		GoogleRefreshToken:  googleRefreshToken,
		TargetSpreadsheetID: targetSpreadsheetID,
	}, nil
}
func (d *Database) Save(userID string, userContext persistence.UserContext) error {
	db, err := d.connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer funcutils.IgnoreError(db.Close)

	insertSQL := `INSERT INTO user_context (userID, googleRefreshToken, targetSpreadsheetID) VALUES (?, ?, ?)`
	_, err = db.Exec(insertSQL, userID, userContext.GoogleRefreshToken, userContext.TargetSpreadsheetID)
	if err != nil {
		return fmt.Errorf("failed to insert into user_context: %w", err)
	}

	return nil
}

func (d *Database) Exists(userID string) (bool, error) {
	db, err := d.connect()
	if err != nil {
		return false, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer funcutils.IgnoreError(db.Close)

	selectSQL := `SELECT COUNT(*) FROM user_context WHERE userID = ?`
	row := db.QueryRow(selectSQL, userID)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to scan row: %w", err)
	}

	return count > 0, nil
}
