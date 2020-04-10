package nzcovid19cases

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"regexp"
	"strconv"
)

type AlertLevel struct {
	Level     int
	LevelName string
}

var re = regexp.MustCompile(`(?m)Level (\d)`)

func ScrapeLevel() (int, string, error) {
	resp, err := soup.Get("https://covid19.govt.nz/")
	if err != nil {
		return 0, "", err
	}

	doc := soup.HTMLParse(resp)

	h3 := doc.Find("h3")
	if h3.Error != nil {
		return 0, "", fmt.Errorf("could not find h3")
	}

	matches := re.FindStringSubmatch(h3.Text())
	if len(matches) != 2 {
		return 0, "", fmt.Errorf("expected two match elements")
	}

	levelString := matches[1]

	levelInt, err := strconv.Atoi(levelString)
	if err != nil {
		return 0, "", fmt.Errorf("could not convert level (%v) to int: %w", levelString, err)
	}

	levelName, ok := levelLookup[levelInt]
	if !ok {
		return 0, "", fmt.Errorf("could not look up level name from level (%v) : %w", levelInt, err)
	}

	return levelInt, levelName, nil
}

func RenderLevels(levelInt int, levelString, viewType string) (string, error) {
	if viewType != "json" {
		return "", InvalidUsageError{fmt.Sprintf("unknown view type: %v", viewType)}
	}

	level := AlertLevel{
		Level:     levelInt,
		LevelName: levelString,
	}

	b, err := json.MarshalIndent(level, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
