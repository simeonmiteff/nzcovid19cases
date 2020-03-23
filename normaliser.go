package nzcovid19cases

import (
	"fmt"
	geojson "github.com/paulmach/go.geojson"
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
	LocationCentrePoint *geojson.Geometry
	//LocationShape *geojson.Feature
	Age                       AgeRange
	Gender                    string
	TravelDetailsUnstructured string
}

var locations = map[string]*geojson.Geometry{
	"Auckland":          geojson.NewPointGeometry([]float64{174.7633, -36.8485}),
	"Canterbury":        geojson.NewPointGeometry([]float64{171.1637, -43.7542}),
	"Dundedin":          geojson.NewPointGeometry([]float64{170.5028, -45.8788}), // Typo
	"Dunedin":           geojson.NewPointGeometry([]float64{170.5028, -45.8788}),
	"Hawkes Bay":        geojson.NewPointGeometry([]float64{176.7416, -39.1090}), // Typo
	"Invercargill":      geojson.NewPointGeometry([]float64{168.3538, -46.4132}),
	"Northland":         geojson.NewPointGeometry([]float64{173.7624, -35.5795}),
	"Queenstown":        geojson.NewPointGeometry([]float64{168.6626, -45.0312}),
	"Rotorua":           geojson.NewPointGeometry([]float64{176.2497, -38.1368}),
	"Southern DHB":      geojson.NewPointGeometry([]float64{170.5086, -45.8694}), // Dunedin coordinates?
	"Taranaki":          geojson.NewPointGeometry([]float64{174.4383, -39.3538}),
	"Waikato":           geojson.NewPointGeometry([]float64{175.1894, -37.4558}),
	"Wellington":        geojson.NewPointGeometry([]float64{174.7762, -41.2865}),
	"Nelson":            geojson.NewPointGeometry([]float64{173.2840, -41.2706}),
	"Manawatu":          geojson.NewPointGeometry([]float64{175.4376, -39.7273}),
	"Taupo":             geojson.NewPointGeometry([]float64{176.0702, -38.6857}),
	"Wellington Region": geojson.NewPointGeometry([]float64{175.4376, -41.0299}),
	"Otago":             geojson.NewPointGeometry([]float64{170.1548, -45.4791}),
	"Hamilton":          geojson.NewPointGeometry([]float64{175.2793, -37.7870}),
	"Bay of Plenty":     geojson.NewPointGeometry([]float64{177.1423, -37.6893}),
	"Coramandel":        geojson.NewPointGeometry([]float64{175.4981, -36.7613}), // Typo
	"Wairarapa":         geojson.NewPointGeometry([]float64{175.6574, -40.9511}), // Masterton coordinates
	"Marlborough":       geojson.NewPointGeometry([]float64{173.4217, -41.5727}),
	"Tasman":            geojson.NewPointGeometry([]float64{172.7347, -41.2122}),
}

var locationNames = map[string]string{
	"Dundedin":          "Dunedin",
	"Hawkes Bay":        "Hawke's Bay",
	"Coramandel":        "Coromandel",
}

var ageLookup = map[string]AgeRange{
	"Child": {Valid: true, OlderOrEqualToAge: 0, YoungerOrEqualToAge: 12},
	"Teens": {Valid: true, OlderOrEqualToAge: 13, YoungerOrEqualToAge: 19}, // Does the MOH use 13-19?
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
	"":    {Valid: false, OlderOrEqualToAge: 0, YoungerOrEqualToAge: 0},
}

var genderLookup = map[string]string{
	"F":      "Female",
	"Female": "Female",
	"M":      "Male",
	"Male":   "Male",
	"":      "Unknown or undisclosed",
}

func (n *NormalisedCase) FromRaw(r *RawCase) error {
	ageRange, ok := ageLookup[strings.TrimSpace(r.Age)]
	if !ok {
		return fmt.Errorf("age string \"%v\" not found in lookup table", r.Age)
	}
	n.Age = ageRange
	gender, ok := genderLookup[strings.TrimSpace(r.Gender)]
	if !ok {
		return fmt.Errorf("gender string \"%v\" not found in lookup table", r.Gender)
	}
	geometry, ok := locations[strings.TrimSpace(r.Location)]
	if !ok {
		return fmt.Errorf("unknown location: \"%v\"", r.Location)
	}
	n.LocationCentrePoint = geometry

	n.Gender = gender
	n.CaseNumber = r.Case
	correctedName, ok := locationNames[r.Location]
	if ok {
		n.LocationName = correctedName
	} else {
		n.LocationName = r.Location
	}
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
