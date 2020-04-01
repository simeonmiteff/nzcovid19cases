package nzcovid19cases

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type AgeRange struct {
	Valid               bool
	OlderOrEqualToAge   int
	YoungerThanAge		int
}

type TravelRelated struct {
	Valid               bool
	Value               bool
}

type TravelDate struct {
	Valid               bool
	Value               time.Time
}


type NormalisedCase struct {
	CaseNumber          int
	ReportedDate		time.Time
	LocationName        string
	Age                 AgeRange
	Gender              string
	IsTravelRelated		TravelRelated
	DepartureDate		TravelDate
	ArrivalDate			TravelDate
	LastCityBeforeNZ	string
	FlightNumber		string
	CaseType			string
}

var yesNoLookup = map[string]TravelRelated{
	"yes": {Valid: true, Value: true},
	"y": {Valid: true, Value: true},
	"no": {Valid: true, Value: false},
	"n": {Valid: true, Value: false},
	"": {Valid: false, Value: false},
	"unknown": {Valid: false, Value: false},
}

var ageLookup = map[string]AgeRange{
	"<1":         {Valid: true, OlderOrEqualToAge: 0, YoungerThanAge: 1},
	"1 to 4":     {Valid: true, OlderOrEqualToAge: 1, YoungerThanAge: 4},
	"5 to 9":     {Valid: true, OlderOrEqualToAge: 1, YoungerThanAge: 4},
	"10 to 14":   {Valid: true, OlderOrEqualToAge: 10, YoungerThanAge: 14},
	"15 to 19":   {Valid: true, OlderOrEqualToAge: 15, YoungerThanAge: 19},
	"20 to 29":   {Valid: true, OlderOrEqualToAge: 20, YoungerThanAge: 29},
	"20 to 29":   {Valid: true, OlderOrEqualToAge: 20, YoungerThanAge: 29}, // FUUUUU
	"30 to 39":   {Valid: true, OlderOrEqualToAge: 30, YoungerThanAge: 39},
	"40 to 49":   {Valid: true, OlderOrEqualToAge: 40, YoungerThanAge: 49},
	"50 to 59":   {Valid: true, OlderOrEqualToAge: 50, YoungerThanAge: 59},
	"50 to 59":   {Valid: true, OlderOrEqualToAge: 50, YoungerThanAge: 59}, // FfffUUUUUU!
	"60 to 69":   {Valid: true, OlderOrEqualToAge: 60, YoungerThanAge: 69},
	"60 to 69":   {Valid: true, OlderOrEqualToAge: 60, YoungerThanAge: 69}, // Unicode spaces
	"70+":        {Valid: true, OlderOrEqualToAge: 70, YoungerThanAge: 110},
	"Unknown":    {Valid: false, OlderOrEqualToAge: 0, YoungerThanAge: 0},
	"":          {Valid: false, OlderOrEqualToAge: 0, YoungerThanAge: 0},
}

var genderLookup = map[string]string{
	"F":      			"Female",
	"Female": 			"Female",
	"M":      			"Male",
	"Male":   			"Male",
	"":      			"Unknown or undisclosed",
	"Not provided":     "Unknown or undisclosed",
}

var levelLookup = map[int]string {
	1:"Prepare",
	2:"Reduce",
	3:"Restrict",
	4:"Eliminate",
}

var locationNames = map[string]string{
	"Capital & Coast": 		"Capital and Coast",
	"Hawkes’s Bay": 		"Hawke's Bay",
	"Hawke’s Bay": 			"Hawke's Bay",
	"Hawkes Bay": 			"Hawke's Bay",
	"Nelson-Marlborough":	"Nelson Marlborough",
	"Southern DHB":			"Southern",
}

var validDHBs = map[string]bool{
	"Auckland":true,
	"Bay of Plenty":true,
	"Canterbury":true,
	"Capital and Coast":true,
	"Counties Manukau":true,
	"Hawke's Bay":true,
	"Hutt Valley":true,
	"Lakes":true,
	"MidCentral":true,
	"Nelson Marlborough":true,
	"Northland":true,
	"South Canterbury":true,
	"Southern":true,
	"Tairawhiti":true,
	"Taranaki":true,
	"Waikato":true,
	"Wairarapa":true,
	"Waitemata":true,
	"West Coast":true,
	"Whanganui":true,
}

var ValidDHBsList =[]string{
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

func (n *NormalisedCase) FromRaw(r *RawCase) error {
	ageRange, ok := ageLookup[strings.TrimSpace(r.Age)]
	if !ok {
		exactAge, err := strconv.Atoi(strings.TrimSpace(r.Age))
		if err == nil {
			ageRange = AgeRange{Valid: true, OlderOrEqualToAge: exactAge, YoungerThanAge: exactAge}
		} else {
			return fmt.Errorf("age string \"%v\" not found in lookup table", r.Age)
		}
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
		return fmt.Errorf("failed to load timezone: %w", tz)
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
		n.DepartureDate = TravelDate{Valid:true, Value:d}
	}

	if strings.TrimSpace(r.ArrivalDate) != "" {
		d, err := time.ParseInLocation(TimeFormat, r.ArrivalDate, tz)
		if err != nil {
			return fmt.Errorf("problem parsing arrival date (%v): %w", r.ArrivalDate, err)
		}
		n.ArrivalDate = TravelDate{Valid:true, Value:d}
	}

	//geometry, ok := locations[n.LocationName]
	//if !ok {
	//	return fmt.Errorf("unknown location: \"%v\"", n.LocationName)
	//}
	//n.LocationCentrePoint = geometry

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
