package nzcovid19cases

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"regexp"
	"strconv"
)

type RawCase struct {
	Case          int
	Location      string
	Age           string
	Gender        string
	TravelDetails string
}

func parseRow(cols []soup.Root) (RawCase, error) {
	var c RawCase
	caseNum, err := strconv.Atoi(cols[0].Text())
	if err != nil {
		return c, fmt.Errorf("failed to convert %v to case number: %w", cols[0].Text(), err)
	}
	c.Case = caseNum
	c.Location = cols[1].Text()
	c.Age = cols[2].Text()
	c.Gender = cols[3].Text()
	c.TravelDetails = cols[4].Text()
	return c, nil
}

func ScrapeCases() ([]*RawCase, error) {
	resp, err := soup.Get("https://www.health.govt.nz/our-work/diseases-and-conditions/covid-19-novel-coronavirus/covid-19-current-cases")
	if err != nil {
		return nil, err
	}
	doc := soup.HTMLParse(resp)
	rows := doc.Find("table", "class", "table-style-two").FindAll("tr")
	var cases []*RawCase

	// Note: slice starting at 1, skipping the header
	for i, row := range rows[1:] {
		cols := row.FindAll("td")
		// This deals with the colspan=5 row that appeared
		if len(cols) == 1 {
			continue
		}
		if len(cols) != 5 {
			return cases, fmt.Errorf("row has %v columns, not 5", len(cols))
		}
		c, err := parseRow(cols)
		if err != nil {
			return nil, fmt.Errorf("problem parsing row %v from html table: %w", i, err)
		} else {
			cases = append(cases, &c)
		}
	}
	return cases, nil
}

var re = regexp.MustCompile(`(?m)Level (\d)`)

func ScrapeLevel() (int, string, error) {
	resp, err := soup.Get("https://covid19.govt.nz/government-actions/covid-19-alert-system")
	if err != nil {
		return 0, "", err
	}
	doc := soup.HTMLParse(resp)
	div := doc.Find("div", "class", "hero-statement")
	if div.Error != nil {
		return 0, "", fmt.Errorf("could not find div")
	}

	matches := re.FindStringSubmatch(div.Text())
	if len(matches) != 2 {
		return 0, "", fmt.Errorf("expected two match elements")
	}

	levelString := matches[1]

	levelInt, err := strconv.Atoi(levelString)
	if err != nil {
		return 0, "", fmt.Errorf("could not convert level (%v) to int: %w", levelString, err)
	}

	levelName, ok := levelLoopup[levelInt]
	if !ok {
		return 0, "", fmt.Errorf("could not look up level name from level (%v) : %w", levelInt, err)
	}

	return levelInt, levelName, nil
}