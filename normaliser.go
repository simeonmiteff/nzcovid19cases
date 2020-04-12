package nzcovid19cases

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AgeRange struct {
	Valid             bool
	OlderOrEqualToAge int
	YoungerThanAge    int
}

type TravelRelated struct {
	Valid bool
	Value bool
}

type TravelDate struct {
	Valid bool
	Value time.Time
}

type NormalisedCase struct {
	CaseNumber       int
	ReportedDate     time.Time
	LocationName     string
	Age              AgeRange
	Gender           string
	IsTravelRelated  TravelRelated
	DepartureDate    TravelDate
	ArrivalDate      TravelDate
	LastCityBeforeNZ string
	FlightNumber     string
	CaseType         string
}

var yesNoLookup = map[string]TravelRelated{ //nolint:gochecknoglobals
	"yes":     {Valid: true, Value: true},
	"y":       {Valid: true, Value: true},
	"no":      {Valid: true, Value: false},
	"n":       {Valid: true, Value: false},
	"":        {Valid: false, Value: false},
	"unknown": {Valid: false, Value: false},
}

var reAgeRange = regexp.MustCompile(`(\d+)`) //nolint:gochecknoglobals

var ageLookup = map[string]AgeRange{ //nolint:gochecknoglobals
	"<1":      {Valid: true, OlderOrEqualToAge: 0, YoungerThanAge: 1},
	"70+":     {Valid: true, OlderOrEqualToAge: 70, YoungerThanAge: 110}, //nolint:gomnd
	"Unknown": {Valid: false, OlderOrEqualToAge: 0, YoungerThanAge: 0},
	"":        {Valid: false, OlderOrEqualToAge: 0, YoungerThanAge: 0},
}

var genderLookup = map[string]string{ //nolint:gochecknoglobals
	"F":            "Female",
	"Female":       "Female",
	"M":            "Male",
	"Male":         "Male",
	"":             "Unknown or undisclosed",
	"Not provided": "Unknown or undisclosed",
}

var levelLookup = map[int]string{ //nolint:gochecknoglobals
	1: "Prepare",
	2: "Reduce",
	3: "Restrict",
	4: "Eliminate",
}

var locationNames = map[string]string{ //nolint:gochecknoglobals
	"Capital & Coast":    "Capital and Coast",
	"Hawkes’s Bay":       "Hawke's Bay",
	"Hawke’s Bay":        "Hawke's Bay",
	"Hawkes Bay":         "Hawke's Bay",
	"Nelson-Marlborough": "Nelson Marlborough",
	"Southern DHB":       "Southern",
}

var validDHBs = map[string]bool{ //nolint:gochecknoglobals
	"Auckland":           true,
	"Bay of Plenty":      true,
	"Canterbury":         true,
	"Capital and Coast":  true,
	"Counties Manukau":   true,
	"Hawke's Bay":        true,
	"Hutt Valley":        true,
	"Lakes":              true,
	"MidCentral":         true,
	"Nelson Marlborough": true,
	"Northland":          true,
	"South Canterbury":   true,
	"Southern":           true,
	"Tairawhiti":         true,
	"Taranaki":           true,
	"Waikato":            true,
	"Wairarapa":          true,
	"Waitemata":          true,
	"West Coast":         true,
	"Whanganui":          true,
}

var ValidDHBsList = []string{ //nolint:gochecknoglobals
	"Auckland",
	"Bay of Plenty",
	"Canterbury",
	"Capital and Coast",
	"Counties Manukau",
	"Hawke's Bay",
	"Hutt Valley",
	"Lakes",
	"MidCentral",
	"Nelson Marlborough",
	"Northland",
	"South Canterbury",
	"Southern",
	"Tairawhiti",
	"Taranaki",
	"Waikato",
	"Wairarapa",
	"Waitemata",
	"West Coast",
	"Whanganui",
}

const TimeFormat = "2/01/2006"

//nolint:funlen
func (n *NormalisedCase) FromRaw(r *RawCase) error {
	age := strings.TrimSpace(r.Age)

	var ageRange AgeRange

	var ok bool

	matches := reAgeRange.FindAllString(age, 2)

	if len(matches) != 2 { //nolint:gomnd
		ageRange, ok = ageLookup[age]
		if !ok {
			exactAge, err := strconv.Atoi(age)
			if err == nil {
				ageRange = AgeRange{Valid: true, OlderOrEqualToAge: exactAge, YoungerThanAge: exactAge}
			} else {
				return fmt.Errorf("age string \"%v\" not found in lookup table", age)
			}
		}
	} else {
		num1, err := strconv.Atoi(matches[0])
		if err != nil {
			return fmt.Errorf("failed to convert %v to number: %w", matches[0], err)
		}
		num2, err := strconv.Atoi(matches[1])
		if err != nil {
			return fmt.Errorf("failed to convert %v to number: %w", matches[1], err)
		}
		ageRange = AgeRange{Valid: true, OlderOrEqualToAge: num1, YoungerThanAge: num2}
	}

	n.Age = ageRange

	gender, ok := genderLookup[strings.TrimSpace(r.Gender)]
	if !ok {
		return fmt.Errorf("gender string \"%v\" not found in lookup table", r.Gender)
	}

	noSpaces := strings.TrimSpace(r.Location)

	correctedName, ok := locationNames[noSpaces]
	if ok {
		_, ok = validDHBs[correctedName]
		if !ok {
			return fmt.Errorf("DHB name \"%v\" not found in lookup table", correctedName)
		}

		n.LocationName = correctedName
	} else {
		n.LocationName = noSpaces
	}

	yesNo, ok := yesNoLookup[strings.TrimSpace(strings.ToLower(r.TravelRelated))]
	if !ok {
		return fmt.Errorf("travel related string \"%v\" not found in lookup table", r.TravelRelated)
	}

	tz, err := time.LoadLocation("Pacific/Auckland")
	if err != nil {
		return fmt.Errorf("failed to load timezone: %w", err)
	}

	reportedDate, err := time.ParseInLocation(TimeFormat, r.ReportedDate, tz)
	if err != nil {
		return fmt.Errorf("problem parsing reported date (%v): %w", r.ReportedDate, err)
	}

	if strings.TrimSpace(r.DepartureDate) != "" {
		d, err := time.ParseInLocation(TimeFormat, r.DepartureDate, tz)
		if err != nil {
			return fmt.Errorf("problem parsing departure date (%v): %w", r.DepartureDate, err)
		}

		n.DepartureDate = TravelDate{Valid: true, Value: d}
	}

	if strings.TrimSpace(r.ArrivalDate) != "" {
		d, err := time.ParseInLocation(TimeFormat, r.ArrivalDate, tz)
		if err != nil {
			return fmt.Errorf("problem parsing arrival date (%v): %w", r.ArrivalDate, err)
		}

		n.ArrivalDate = TravelDate{Valid: true, Value: d}
	}

	n.ReportedDate = reportedDate
	n.Gender = gender
	n.CaseNumber = r.Case
	n.CaseType = r.CaseType
	n.IsTravelRelated = yesNo
	n.LastCityBeforeNZ = r.LastCityBeforeNZ
	n.FlightNumber = r.FlightNumber

	return nil
}

func NormaliseCases(rawCases []*RawCase) ([]*NormalisedCase, error) {
	normCases := make([]*NormalisedCase, len(rawCases))

	for i, cp := range rawCases {
		var normCase NormalisedCase

		err := normCase.FromRaw(cp)
		if err != nil {
			return nil, fmt.Errorf("problem normalising case from line %v: %w", i, err)
		}

		normCases[i] = &normCase
	}

	return normCases, nil
}
