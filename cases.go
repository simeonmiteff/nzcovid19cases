package nzcovid19cases

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	geojson "github.com/paulmach/go.geojson"
	"os"
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

type CaseStats struct {
	Confirmed				int
	Recovered				int
	CommunityTransmission	int
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

	stats := doc.Find("div", "property", "content:encoded").FindAll("li")

	for i, stat := range stats {
		parts := strings.Split(stat.Text(), " ")
		fmt.Println(i, parts[len(parts)-1])
	}

	os.Exit(0)


	return cases, nil
}

func RenderCases(normCases []*NormalisedCase, viewType string) (string, error) {
	validViewTypes := map[string]bool{
		"csv":     true,
		"json":    true,
		"geojson": true,
	}
	if !validViewTypes[viewType] {
		return "", InvalidUsageError{fmt.Sprintf("unknown view type: %v", viewType)}
	}

	var sb strings.Builder
	switch viewType {
	case "csv":
		sb.WriteString(`"CaseNumber", "LocationName", "AgeValid", "OlderOrEqualToAge", "YoungerOrEqualToAge"` +
			`,"Gender", "TravelDetailsUnstructured", "LocationCentrePointLongitude",` +
			`"LocationCentrePointLatitude"`)
		sb.WriteRune('\n')
		for _, c := range normCases {
			sb.WriteString(fmt.Sprintf(`%v, "%v", "%v", %v, %v, "%v", "%v", %v, %v`,
				c.CaseNumber, c.LocationName, c.Age.Valid, c.Age.OlderOrEqualToAge,
				c.Age.YoungerOrEqualToAge, c.Gender, c.TravelDetailsUnstructured,
				c.LocationCentrePoint.Point[0], c.LocationCentrePoint.Point[1]))
			sb.WriteRune('\n')
		}
	case "json":
		b, err := json.MarshalIndent(normCases, "", "  ")
		if err != nil {
			return "", err
		}
		sb.Write(b)
	case "geojson":
		fc := geojson.NewFeatureCollection()
		for _, c := range normCases {
			var feature geojson.Feature
			feature.Geometry = c.LocationCentrePoint
			feature.SetProperty("LocationName", c.LocationName)
			feature.SetProperty("CaseNumber", c.CaseNumber)
			feature.SetProperty("AgeValid", c.Age.Valid)
			feature.SetProperty("OlderOrEqualToAge", c.Age.OlderOrEqualToAge)
			feature.SetProperty("YoungerOrEqualToAge", c.Age.YoungerOrEqualToAge)
			feature.SetProperty("Gender", c.Gender)
			feature.SetProperty("Travel details", c.TravelDetailsUnstructured)
			fc.AddFeature(&feature)
		}
		b, err := fc.MarshalJSON()
		if err != nil {
			return "", err
		}
		sb.Write(b)
		sb.WriteRune('\n')
	}
	return sb.String(), nil
}
