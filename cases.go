package nzcovid19cases

import (
	"encoding/json"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"strings"
)

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
