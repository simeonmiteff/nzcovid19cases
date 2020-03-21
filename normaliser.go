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
	"Dunedin":           geojson.NewPointGeometry([]float64{170.5028, -45.8788}),
	"Hawkes Bay":        geojson.NewPointGeometry([]float64{176.7416, -39.1090}),
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
}

var ageLookup = map[string]AgeRange{
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
	n.LocationName = r.Location
	n.TravelDetailsUnstructured = r.TravelDetails

	return nil
}
