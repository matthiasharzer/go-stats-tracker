package spreadsheet

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/matthiasharzer/go-stats-tracker/analyzer"
	"github.com/matthiasharzer/go-stats-tracker/logging"
	"github.com/matthiasharzer/go-stats-tracker/persistence"
	"github.com/matthiasharzer/go-stats-tracker/utils/sheetutils"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type cell[T any] struct {
	column string
	row    int
	value  T
}

func (c cell[T]) format() string {
	return fmt.Sprintf("%s%d", c.column, c.row)
}

type cellRange struct {
	from cell[any]
	to   cell[any]
}

func (c cellRange) formatString() string {
	return fmt.Sprintf("%s:%s", c.from.format(), c.to.format())
}

var TargetRange = cellRange{
	from: cell[any]{column: "A", row: 3},
	to:   cell[any]{column: "C", row: 1000},
}

func ParseSheetsSerial(serial float64) time.Time {
	baseDate := time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)

	days := math.Floor(serial)
	fraction := serial - days

	dayInNanoseconds := float64(24 * time.Hour)
	nanoseconds := math.Round(fraction * dayInNanoseconds)

	result := baseDate.AddDate(0, 0, int(days)).Add(time.Duration(nanoseconds))
	return result
}

type row struct {
	date  cell[time.Time]
	level cell[int]
	xp    cell[int64]
}

func findRow(rows []row, date time.Time) *row {
	for _, row := range rows {
		if row.date.value.Equal(date) {
			return &row
		}
	}
	return nil
}

func parseRow(rowValues []any, currentRowNumber int) (row, error) {
	if len(rowValues) != 3 {
		return row{}, fmt.Errorf("invalid row length: expected 3, got %d", len(rowValues))
	}
	dateSerial, ok := rowValues[0].(float64)
	if !ok {
		return row{}, fmt.Errorf("failed to parse date")
	}
	date := ParseSheetsSerial(dateSerial)
	level := rowValues[1].(float64)
	xp := rowValues[2].(float64)
	return row{
		date: cell[time.Time]{
			column: TargetRange.from.column,
			row:    currentRowNumber,
			value:  date,
		},
		level: cell[int]{
			column: sheetutils.IncrementColumnN(TargetRange.from.column, 1),
			row:    currentRowNumber,
			value:  int(level),
		},
		xp: cell[int64]{
			column: sheetutils.IncrementColumnN(TargetRange.from.column, 2),
			row:    currentRowNumber,
			value:  int64(xp),
		},
	}, nil
}

type Service struct {
	oauth oauth2.Config
}

func NewService(oauth oauth2.Config) *Service {
	return &Service{
		oauth: oauth,
	}
}

func (s *Service) parseExistingRows(srv *sheets.Service, spreadsheetID string) ([]row, error) {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, TargetRange.formatString()).ValueRenderOption("UNFORMATTED_VALUE").Do()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve spreadsheet data: %w", err)
	}

	currentRowNumber := TargetRange.from.row
	var rows []row
	for _, rowValues := range resp.Values {
		if len(rowValues) < 3 {
			break
		}
		parsedRow, err := parseRow(rowValues, currentRowNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to parse row: %w", err)
		}
		rows = append(rows, parsedRow)
		currentRowNumber++
	}
	return rows, nil
}

//func addRow(srv *sheets.Service, spreadsheetID string, date time.Time, xp int64)

func (s *Service) setRow(srv *sheets.Service, spreadsheetID string, dateCell cell[time.Time], level int, xp int64) error {
	targetRange := cellRange{
		from: cell[any]{
			column: dateCell.column,
			row:    dateCell.row,
		},
		to: cell[any]{
			column: sheetutils.IncrementColumnN(TargetRange.from.column, 2),
			row:    dateCell.row,
		},
	}

	dateFormatted := dateCell.value.Format("2006-01-02")

	var vr sheets.ValueRange
	row := []any{dateFormatted, level, xp}
	vr.Values = append(vr.Values, row)

	_, err := srv.Spreadsheets.Values.Update(spreadsheetID, targetRange.formatString(), &vr).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		return fmt.Errorf("failed to append row: %w", err)
	}
	return nil
}

func (s *Service) resolveDateCell(rows []row, date time.Time) (cell[time.Time], error) {
	if len(rows) == 0 {
		return cell[time.Time]{
			column: TargetRange.from.column,
			row:    TargetRange.from.row,
			value:  date,
		}, nil
	}
	latestRow := rows[len(rows)-1]
	existingRow := findRow(rows, date)

	if existingRow != nil {
		logging.Info("found existing row", "date", date.Format("2006-01-02"), "row", existingRow.date.row)
		return existingRow.date, nil
	}

	if date.Before(latestRow.date.value) {
		return cell[time.Time]{}, fmt.Errorf("cannot insert out of order date")
	}

	logging.Info("no existing row found, creating new one", "date", date.Format("2006-01-02"))
	return cell[time.Time]{
		column: latestRow.date.column,
		row:    latestRow.date.row + 1,
		value:  date,
	}, nil
}

func (s *Service) AppendStats(userContext persistence.UserContext, stats analyzer.PlayerStats, date time.Time) error {
	ctx := context.Background()

	token := &oauth2.Token{
		RefreshToken: userContext.GoogleRefreshToken,
		TokenType:    "Bearer",
	}

	client := s.oauth.Client(ctx, token)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	rows, err := s.parseExistingRows(srv, userContext.TargetSpreadsheetID)
	if err != nil {
		return fmt.Errorf("failed to parse existing rows: %w", err)
	}
	if len(rows) == 0 {
		return fmt.Errorf("empty rows no supported yet")
	}

	targetDateCell, err := s.resolveDateCell(rows, date)
	if err != nil {
		return fmt.Errorf("failed to resolve date cell: %w", err)
	}

	logging.Info("writing row", "spreadsheetID", userContext.TargetSpreadsheetID, "date", targetDateCell.value.Format("2006-01-02"), "dateCell", targetDateCell.format(), "level", stats.Level, "xp", stats.GainedLevelXP)
	err = s.setRow(srv, userContext.TargetSpreadsheetID, targetDateCell, stats.Level, stats.GainedLevelXP)
	if err != nil {
		return fmt.Errorf("failed to set row: %w", err)
	}

	return nil
}
