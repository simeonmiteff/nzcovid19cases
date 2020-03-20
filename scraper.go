package nzcovid19cases

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"strconv"
)

type RawCase struct {
	Case int
	Location string
	Age string
	Gender string
	TravelDetails string
}

func parseRow(root soup.Root) (RawCase, error) {
	cols := root.FindAll("td")
	var c RawCase
	if len(cols) != 5 {
		return c, fmt.Errorf("row has %v columns, not 5", len(cols))
	}
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

func Scrape() ([]*RawCase, error) {
	resp, err := soup.Get("https://www.health.govt.nz/our-work/diseases-and-conditions/covid-19-novel-coronavirus/covid-19-current-cases")
	if err != nil {
		return nil, err
	}
	doc := soup.HTMLParse(resp)
	rows := doc.Find("table", "class", "table-style-two").FindAll("tr")
	cases := make([]*RawCase, len(rows)-1)

	// Note: slice starting at 1, skipping the header
	for i, row := range rows[1:] {
		c, err := parseRow(row)
		if err != nil {
			return nil, fmt.Errorf("problem parsing row %v from html table: %w", i, err)
		}
		cases[i] = &c
	}
	return cases, nil
}