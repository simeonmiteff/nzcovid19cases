package nzcovid19cases

import (
	"encoding/json"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"log"
	"strings"
)

func RenderCases(viewType string) (string, error) {
	var sb strings.Builder
	switch viewType {
	case "csv":
		cases, err := Scrape()
		if err != nil {
			log.Fatal(err)
		}
		sb.WriteString(`"Case", "Location", "Age", "Gender", "Travel details"`)
		for _, cp := range cases {
			c := *cp
			sb.WriteString(fmt.Sprintf(`"%v", "%v", "%v", "%v", "%v"`, c.Case, c.Location, c.Age, c.Gender, c.TravelDetails))
			sb.WriteRune('\n')
		}
	case "json":
		cases, err := Scrape()
		if err != nil {
			return "", err
		}
		b, err := json.MarshalIndent(cases, "", "  ")
		if err != nil {
			return "", err
		}
		sb.Write(b)
	case "geojson":
		cases, err := Scrape()
		if err != nil {
			return "", err
		}
		fc := geojson.NewFeatureCollection()
		for i, cp := range cases {
			c := *cp
			loc, err := GetLocation(c.Location)
			if err != nil {
				return "", fmt.Errorf("problem getting location for table line %v: %w", i, err)
			}
			var feature geojson.Feature
			feature.Geometry = loc
			feature.SetProperty("Location", c.Location)
			feature.SetProperty("Case", c.Case)
			feature.SetProperty("Age", c.Age)
			feature.SetProperty("Gender", c.Gender)
			feature.SetProperty("Travel details", c.TravelDetails)
			fc.AddFeature(&feature)
		}
		b, err := fc.MarshalJSON()
		if err != nil {
			return "", err
		}
		sb.Write(b)
		sb.WriteRune('\n')
	default:
		return "", InvalidUsageError{fmt.Sprintf("unknown view type: %v", viewType)}
	}
	return sb.String(), nil
}