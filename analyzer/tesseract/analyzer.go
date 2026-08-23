package tesseract

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/otiai10/gosseract/v2"

	"github.com/matthiasharzer/go-stats-tracker/analyzer"
	"github.com/matthiasharzer/go-stats-tracker/logging"
)

// Matches `<level> ??? LEVEL` where ??? may be anything and contain up to 2 line breaks, since the level number is
// usually not on the same line as the LEVEL text
var levelRegex = regexp.MustCompile(`(?P<level>\d+)(.*\s?){0,2}LEVEL`)

// Matches `xxx.xxx.xxx / yyy.yyy.yyy`, where the left side of / is the current gained xp for the level and the right
// side is the total XP required to level up. Includes tolerance for optional whitespace around the /
var xpRegex = regexp.MustCompile(`(?P<gained>(\d+\.?)+)\s*/\s*(?P<total>(\d+\.?)+)`)

var userPageRegexes []*regexp.Regexp

func init() {
	for _, loc := range localizations {
		regexStr := fmt.Sprintf("(?i)%s\\s+%s\\s+%s", loc.You, loc.Friends, loc.Groups)
		userPageRegexes = append(userPageRegexes, regexp.MustCompile(regexStr))
	}
}

func parseLevel(levelStr string) (int64, error) {
	levelNormalized := strings.ReplaceAll(levelStr, ".", "")
	level, err := strconv.ParseInt(levelNormalized, 10, 64)
	if err != nil {
		return 0, errors.New("failed to convert player level to int")
	}
	return level, nil
}

type Analyzer struct {
	imageData []byte
	text      string
}

func newClient() (*gosseract.Client, func()) {
	client := gosseract.NewClient()
	return client, func() {
		err := client.Close()
		if err != nil {
			logging.Warn("Failed to close tesseract client", "err", err)
		}
	}
}

func NewAnalyzer(imageData []byte) (analyzer.Analyzer, error) {
	client, closeClient := newClient()
	defer closeClient()

	err := client.SetImageFromBytes(imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to set image from bytes: %w", err)
	}

	text, err := client.Text()
	if err != nil {
		return nil, fmt.Errorf("failed to extract text: %w", err)
	}

	return &Analyzer{
		imageData: imageData,
		text:      text,
	}, nil
}

func (a *Analyzer) extractPlayerLevel() (int, error) {
	matches := levelRegex.FindStringSubmatch(a.text)
	if len(matches) < 2 {
		return 0, errors.New("could not extract player level")
	}

	level, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, errors.New("failed to convert player level to int")
	}

	return level, nil
}

func (a *Analyzer) extractPlayerXP() (gained int64, total int64, err error) {
	matches := xpRegex.FindStringSubmatch(a.text)

	if len(matches) < 3 {
		return 0, 0, errors.New("could not extract player xp")
	}

	gained, err = parseLevel(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf("could not extract player gained xp: %w", err)
	}

	total, err = parseLevel(matches[3])
	if err != nil {
		return 0, 0, fmt.Errorf("could not extract player total xp: %w", err)
	}

	return gained, total, nil
}

func (a *Analyzer) ExtractPlayerStats() (analyzer.PlayerStats, error) {
	level, err := a.extractPlayerLevel()
	if err != nil {
		logging.Warn("failed to extract player level")
	}

	gained, total, err := a.extractPlayerXP()
	if err != nil {
		return analyzer.PlayerStats{}, fmt.Errorf("failed to extract player xp: %w", err)
	}

	return analyzer.PlayerStats{
		Level:         level,
		TotalLevelXP:  total,
		GainedLevelXP: gained,
	}, nil
}

func (a *Analyzer) IsPlayerPage() (bool, error) {
	for _, regex := range userPageRegexes {
		if regex.MatchString(a.text) {
			return true, nil
		}
	}
	return false, nil
}
