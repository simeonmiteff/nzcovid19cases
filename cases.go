package nzcovid19cases

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

type RawCase struct {
	ReportedDate     string
	Case             int
	Location         string
	Age              string
	Gender           string
	TravelRelated    string
	LastCityBeforeNZ string
	FlightNumber     string
	DepartureDate    string
	ArrivalDate      string
	CaseType         string
}

type CaseStatsResponse struct {
	ConfirmedCasesTotal     int
	ConfirmedCasesNew24h    int
	ProbableCasesTotal      int
	ProbableCasesNew24h     int
	RecoveredCasesTotal     int
	RecoveredCasesNew24h    int
	HospitalisedCasesTotal  int
	HospitalisedCasesNew24h int
	DeathCasesTotal         int
	DeathCasesNew24h        int
}

func parseRow(cols []soup.Root) RawCase {
	var c RawCase
	c.ReportedDate = cols[0].Text()
	c.Gender = cols[1].Text()
	c.Age = cols[2].Text()
	c.Location = cols[3].Text()
	// Handle case with erroneous colspan
	_, sigh := cols[3].Attrs()["colspan"]
	if sigh {
		return c
	}

	c.TravelRelated = cols[4].Text()
	c.LastCityBeforeNZ = cols[5].Text()
	c.FlightNumber = cols[6].Text()
	c.ArrivalDate = cols[7].Text()
	c.DepartureDate = cols[8].Text()

	return c
}

const (
	CSVRenderType  = "csv"
	JSONRenderType = "json"
)

func parseStat(stat soup.Root) (int, int, error) {
	tds := stat.FindAll("td")
	if len(tds) != 2 { //nolint:gomnd
		return 0, 0, fmt.Errorf("expected two columns")
	}

	numString := strings.ReplaceAll(tds[0].Text(), ",", "")

	num, err := strconv.Atoi(numString)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert %v to number: %w", tds[1].Text(), err)
	}

	if strings.TrimSpace(tds[1].Text()) == "" {
		return num, 0, nil
	}

	num24hString := strings.ReplaceAll(tds[1].Text(), ",", "")

	num24h, err := strconv.Atoi(num24hString)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert %v to number: %w", tds[1].Text(), err)
	}

	return num, num24h, nil
}

func ScrapeCases() ([]*RawCase, error) {
	resp, err := soup.Get("https://www.health.govt.nz/our-work/diseases-and-conditions/" +
		"covid-19-novel-coronavirus/covid-19-current-cases/covid-19-current-cases-details")
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(resp)

	tables := doc.FindAll("table", "class", "table-style-two")

	var offset = 0

	if len(tables) > 2 { //nolint:gomnd
		offset = 1
	}

	rows := tables[offset].FindAll("tr")

	var cases []*RawCase //nolint:prealloc

	// Note: slice starting at 1, skipping the header
	for i, row := range rows[1:] {
		cols := row.FindAll("td")

		if len(cols) < 8 { //nolint:gomnd
			return nil, fmt.Errorf("table 1 row has %v columns, not at least 8", len(cols))
		}

		c := parseRow(cols)

		c.CaseType = "confirmed"
		c.Case = len(rows) - i - 1

		cases = append(cases, &c)
	}

	rows = tables[offset+1].FindAll("tr")

	// Note: slice starting at 1, skipping the header
	for i, row := range rows[1:] {
		cols := row.FindAll("td")

		if len(cols) != 9 { //nolint:gomnd
			return nil, fmt.Errorf("table 1 row has %v columns, not 9", len(cols))
		}

		c := parseRow(cols)

		c.CaseType = "probable"
		c.Case = len(rows) - i - 1

		cases = append(cases, &c)
	}

	return cases, nil
}

func RenderCases(normCases []*NormalisedCase, viewType string) (string, error) {
	validViewTypes := map[string]bool{
		CSVRenderType:  true,
		JSONRenderType: true,
	}
	if !validViewTypes[viewType] {
		return "", InvalidUsageError{fmt.Sprintf("unknown view type: %v", viewType)}
	}

	var sb strings.Builder

	switch viewType {
	case CSVRenderType:
		sb.WriteString(
			`"CaseNumber",` +
				`"ReportedDate",` +
				`"LocationName",` +
				`"AgeValid",` +
				`"OlderOrEqualToAge",` +
				`"YoungerThanAge",` +
				`"Gender",` +
				`"IsTravelRelated",` +
				`"DepartureDateValid",` +
				`"DepartureDate",` +
				`"ArrivalDateValid",` +
				`"ArrivalDate",` +
				`"LastCityBeforeNZ",` +
				`"FlightNumber",` +
				`"CaseType"`,
		)
		sb.WriteRune('\n')

		for _, c := range normCases {
			sb.WriteString(fmt.Sprintf(
				`%v, "%v", "%v", "%v", "%v", "%v", %v, %v, "%v", "%v", "%v", "%v", "%v", "%v", "%v", "%v"`,
				c.CaseNumber, c.ReportedDate, c.LocationName, c.Age.Valid, c.Age.OlderOrEqualToAge,
				c.Age.YoungerThanAge, c.Gender, c.IsTravelRelated.Valid, c.IsTravelRelated.Value,
				c.DepartureDate.Valid, c.DepartureDate.Value, c.ArrivalDate.Valid, c.ArrivalDate.Value,
				c.LastCityBeforeNZ, c.FlightNumber, c.CaseType),
			)
			sb.WriteRune('\n')
		}

	case JSONRenderType:
		b, err := json.MarshalIndent(normCases, "", "  ")
		if err != nil {
			return "", err
		}

		sb.Write(b)
	}

	return sb.String(), nil
}

func ScrapeCaseStats() (CaseStatsResponse, error) {
	var cS CaseStatsResponse

	resp, err := soup.Get("https://www.health.govt.nz/our-work/diseases-and-conditions/" +
		"covid-19-novel-coronavirus/covid-19-current-cases")
	if err != nil {
		return cS, err
	}

	doc := soup.HTMLParse(resp)

	tables := doc.FindAll("table")
	stats := tables[0].FindAll("tr")

	if len(stats) != 7 { //nolint:gomnd
		return cS, fmt.Errorf("stats table has %v TR, not 7", len(stats))
	}

	cS.ConfirmedCasesTotal, cS.ConfirmedCasesNew24h, err = parseStat(stats[1])
	if err != nil {
		return cS, fmt.Errorf("problem parsing confirmed cases %w", err)
	}

	cS.ProbableCasesTotal, cS.ProbableCasesNew24h, err = parseStat(stats[2])
	if err != nil {
		return cS, fmt.Errorf("problem parsing probable cases %w", err)
	}

	cS.HospitalisedCasesTotal, cS.HospitalisedCasesNew24h, err = parseStat(stats[4]) //parseStatHosp(stats[4])
	if err != nil {
		return cS, fmt.Errorf("problem parsing hospitalised cases %w", err)
	}

	cS.RecoveredCasesTotal, cS.RecoveredCasesNew24h, err = parseStat(stats[5])
	if err != nil {
		return cS, fmt.Errorf("problem parsing recovered cases %w", err)
	}

	cS.DeathCasesTotal, cS.DeathCasesNew24h, err = parseStat(stats[6])
	if err != nil {
		return cS, fmt.Errorf("problem parsing dead cases %w", err)
	}

	return cS, nil
}

func RenderCaseStats(cS CaseStatsResponse, viewType string) (string, error) {
	var sb strings.Builder

	if viewType != JSONRenderType {
		return "", InvalidUsageError{fmt.Sprintf("unknown view type: %v", viewType)}
	}

	b, err := json.MarshalIndent(cS, "", "  ")
	if err != nil {
		return "", err
	}

	sb.Write(b)

	return sb.String(), nil
}
