package nzcovid19cases

import (
	"fmt"
	"strconv"
	"strings"
)

type AgeRange struct {
	Valid               bool
	OlderOrEqualToAge   int
	YoungerThanAge		int
}

type NormalisedCase struct {
	CaseNumber          int
	LocationName        string
	Age                 AgeRange
	Gender              string
	CaseType			string
}

var ageLookup = map[string]AgeRange{
	"<1":         {Valid: true, OlderOrEqualToAge: 0, YoungerThanAge: 1},
	"1 to 4":     {Valid: true, OlderOrEqualToAge: 1, YoungerThanAge: 4},
	"5 to 9":     {Valid: true, OlderOrEqualToAge: 1, YoungerThanAge: 4},
	"10 to 14":   {Valid: true, OlderOrEqualToAge: 10, YoungerThanAge: 14},
	"15 to 19":   {Valid: true, OlderOrEqualToAge: 15, YoungerThanAge: 19},
	"20 to 29":   {Valid: true, OlderOrEqualToAge: 20, YoungerThanAge: 29},
	"30 to 39":   {Valid: true, OlderOrEqualToAge: 30, YoungerThanAge: 39},
	"40 to 49":   {Valid: true, OlderOrEqualToAge: 40, YoungerThanAge: 49},
	"50 to 59":   {Valid: true, OlderOrEqualToAge: 50, YoungerThanAge: 59},
	"60 to 69":   {Valid: true, OlderOrEqualToAge: 60, YoungerThanAge: 69},
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
	"Nelson Marlborough":	"Nelson-Marlborough",
	"Southern DHB":			"Southern",
}

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
		n.LocationName = correctedName
	} else {
		n.LocationName = noSpaces
	}

	//geometry, ok := locations[n.LocationName]
	//if !ok {
	//	return fmt.Errorf("unknown location: \"%v\"", n.LocationName)
	//}
	//n.LocationCentrePoint = geometry

	n.Gender = gender
	n.CaseNumber = r.Case
	n.CaseType = r.CaseType

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
