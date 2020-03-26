package nzcovid19cases

import (
	"fmt"
	"strconv"
	"strings"
)

type AgeRange struct {
	Valid               bool
	OlderOrEqualToAge   int
	YoungerOrEqualToAge int
}

type NormalisedCase struct {
	CaseNumber          int
	LocationName        string
	Age                       AgeRange
	Gender                    string
	TravelDetailsUnstructured string
	CaseType				  string
}

var ageLookup = map[string]AgeRange{
	"Child": {Valid: true, OlderOrEqualToAge: 0, YoungerOrEqualToAge: 12},
	"Teens": {Valid: true, OlderOrEqualToAge: 13, YoungerOrEqualToAge: 19}, // Does the MOH use 13-19?
	"Teen":  {Valid: true, OlderOrEqualToAge: 13, YoungerOrEqualToAge: 19}, // Does the MOH use 13-19?
	"20s":   {Valid: true, OlderOrEqualToAge: 20, YoungerOrEqualToAge: 29},
	"30s":   {Valid: true, OlderOrEqualToAge: 30, YoungerOrEqualToAge: 39},
	"40s":   {Valid: true, OlderOrEqualToAge: 40, YoungerOrEqualToAge: 49},
	"50s":   {Valid: true, OlderOrEqualToAge: 50, YoungerOrEqualToAge: 59},
	"60s":   {Valid: true, OlderOrEqualToAge: 60, YoungerOrEqualToAge: 69},
	"70s":   {Valid: true, OlderOrEqualToAge: 70, YoungerOrEqualToAge: 79},
	// Not seen in the data (yet)
	"80s":  {Valid: true, OlderOrEqualToAge: 80, YoungerOrEqualToAge: 89},
	"90s":  {Valid: true, OlderOrEqualToAge: 90, YoungerOrEqualToAge: 99},
	"100s": {Valid: true, OlderOrEqualToAge: 100, YoungerOrEqualToAge: 109},
	"Unknown":    {Valid: false, OlderOrEqualToAge: 0, YoungerOrEqualToAge: 0},
}

var genderLookup = map[string]string{
	"F":      			"Female",
	"Female": 			"Female",
	"M":      			"Male",
	"Male":   			"Male",
	"":      			"Unknown or undisclosed",
	"Not provided":     "Unknown or undisclosed",
}

var levelLoopup = map[int]string {
	1:"Prepare",
	2:"Reduce",
	3:"Restrict",
	4:"Eliminate",
}

func (n *NormalisedCase) FromRaw(r *RawCase) error {
	ageRange, ok := ageLookup[strings.TrimSpace(r.Age)]
	if !ok {
		exactAge, err := strconv.Atoi(r.Age)
		if err == nil {
			ageRange = AgeRange{Valid: true, OlderOrEqualToAge: exactAge, YoungerOrEqualToAge: exactAge}
		} else {
			return fmt.Errorf("age string \"%v\" not found in lookup table", r.Age)
		}
	}
	n.Age = ageRange
	gender, ok := genderLookup[strings.TrimSpace(r.Gender)]
	if !ok {
		return fmt.Errorf("gender string \"%v\" not found in lookup table", r.Gender)
	}

	//noSpaces := strings.TrimSpace(r.Location)
	//correctedName, ok := locationNames[noSpaces]
	//if ok {
	//	n.LocationName = correctedName
	//} else {
	//	n.LocationName = noSpaces
	//}
	//
	//geometry, ok := locations[n.LocationName]
	//if !ok {
	//	return fmt.Errorf("unknown location: \"%v\"", n.LocationName)
	//}
	//n.LocationCentrePoint = geometry

	n.LocationName = r.Location

	n.Gender = gender
	n.CaseNumber = r.Case
	n.CaseType = r.CaseType

	n.TravelDetailsUnstructured = r.TravelDetails

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
