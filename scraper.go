package nzcovid19cases

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"regexp"
	"strconv"
	"strings"
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
	resp, err := soup.Get("https://covid19.govt.nz/government-actions/covid-19-alert-level")
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

func ScrapeGrants() (GrantsSummary, GrantRegions, error) {
	var gS GrantsSummary
	var gR GrantRegions

	resp, err := soup.Get("https://www.msd.govt.nz/about-msd-and-our-work/newsroom/2020/covid-19/covid-19-data.html")
	if err != nil {
		return gS, gR, err
	}
	doc := soup.HTMLParse(resp)
	div := doc.Find("div", "id", "content")
	if div.Error != nil {
		return gS, gR, fmt.Errorf("could not find content div")
	}

	tables := div.FindAll("table")
	if len(tables) != 2 {
		return gS, gR, fmt.Errorf("content div has fewer than 2 (%v) tables", len(tables))
	}

	tds := tables[0].FindAll("td")

	gS.Clients, err = strconv.Atoi(strings.ReplaceAll(tds[0].Text(), ",", ""))
	if err != nil {
		return gS, gR,fmt.Errorf("could not convert client count (%v) to int", tds[0].Text())
	}
	gS.Grants, err = strconv.Atoi(strings.ReplaceAll(tds[1].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert grant count (%v) to int", tds[1].Text())
	}
	gS.SumGrantAmount, err = strconv.Atoi(strings.ReplaceAll(tds[2].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert grant total amount (%v) to int", tds[2].Text())
	}

	tds = tables[1].FindAll("td")
	gR.Auckland, err = strconv.Atoi(strings.ReplaceAll(tds[0].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Auckland region client count (%v) to int", tds[0].Text())
	}
	gR.EastCoast, err = strconv.Atoi(strings.ReplaceAll(tds[1].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert EastCoast region client count (%v) to int", tds[1].Text())
	}
	gR.BayOfPlenty, err = strconv.Atoi(strings.ReplaceAll(tds[2].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert BayOfPlenty region client count (%v) to int", tds[2].Text())
	}
	gR.Northland, err = strconv.Atoi(strings.ReplaceAll(tds[3].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Northland region client count (%v) to int", tds[3].Text())
	}
	gR.Wellington, err = strconv.Atoi(strings.ReplaceAll(tds[4].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Wellington region client count (%v) to int", tds[4].Text())
	}
	gR.Nelson, err = strconv.Atoi(strings.ReplaceAll(tds[5].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Nelson region client count (%v) to int", tds[5].Text())
	}
	gR.Canterbury, err = strconv.Atoi(strings.ReplaceAll(tds[6].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Canterbury region client count (%v) to int", tds[6].Text())
	}
	gR.Southern, err = strconv.Atoi(strings.ReplaceAll(tds[7].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Southern region client count (%v) to int", tds[7].Text())
	}
	gR.Other, err = strconv.Atoi(strings.ReplaceAll(tds[8].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Other region client count (%v) to int", tds[8].Text())
	}
	gR.Total, err = strconv.Atoi(strings.ReplaceAll(tds[9].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Total region client count (%v) to int", tds[9].Text())
	}

	return gS, gR, nil
}