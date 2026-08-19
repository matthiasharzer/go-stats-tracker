package spreadsheet

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/matthiasharzer/go-stats-tracker/analyzer"
	"github.com/matthiasharzer/go-stats-tracker/persistence"
	"github.com/matthiasharzer/go-stats-tracker/utils/timeutils"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const TargetRange = "Sheet1!A3:C1000"

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
	date  time.Time
	level int
	xp    int64
}

func findRow(rows []row, date time.Time) *row {
	for _, row := range rows {
		if row.date.Equal(date) {
			return &row
		}
	}
	return nil
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
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, TargetRange).ValueRenderOption("UNFORMATTED_VALUE").Do()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve spreadsheet data: %w", err)
	}

	var rows []row
	for _, rowValues := range resp.Values {
		if len(rowValues) < 3 {
			break
		}
		dateSerial, ok := rowValues[0].(float64)
		if !ok {
			return nil, fmt.Errorf("failed to parse date")
		}
		date := ParseSheetsSerial(dateSerial)
		level := rowValues[1].(float64)
		xp := rowValues[2].(float64)
		rows = append(rows, row{date: date, level: int(level), xp: int64(xp)})
	}
	return rows, nil
}

//func addRow(srv *sheets.Service, spreadsheetID string, date time.Time, xp int64)

func setRow(srv *sheets.Service, spreadsheetID string, date time.Time, level int, xp int64, setLevel bool) error {

}

func (s *Service) AppendStats(userContext persistence.UserContext, stats analyzer.PlayerStats) error {
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

	today := timeutils.TodayDate().Add(24 * time.Hour)
	yesterday := today.Add(-1 * 24 * time.Hour)

	todaysRow := findRow(rows, today)
	yesterdaysRow := findRow(rows, yesterday)

	if todaysRow {

	}

	_ = todaysRow

	//srv.Spreadsheets.Values.Get(userContext.TargetSpreadsheetID, "Sheet1!A3:B1000")

	//timestamp := time.Now().Format(time.RFC3339)

	return nil

}
