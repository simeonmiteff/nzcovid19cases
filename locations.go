package nzcovid19cases

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Location struct {
	LocationName        string
	//LocationShape *geojson.Feature
	CaseCount int
}

func BuildLocations(normCases []*NormalisedCase) map[string]*Location{
	locations := make(map[string]*Location)
	for _, c := range normCases {
		loc, ok := locations[c.LocationName]
		if !ok {
			locations[c.LocationName] = &Location{
				LocationName:        c.LocationName,
				CaseCount:           1,
			}
			continue
		}
		loc.CaseCount = loc.CaseCount + 1
	}
	return locations
}

func RenderLocations(locations map[string]*Location, viewType string) (string, error) {
	validViewTypes := map[string]bool{
		"csv":     true,
		"json":    true,
		//"geojson": true,
	}
	if !validViewTypes[viewType] {
		return "", InvalidUsageError{fmt.Sprintf("unknown view type: %v", viewType)}
	}

	var sb strings.Builder
	switch viewType {
	case "csv":
		sb.WriteString(`"LocationName", "CaseCount"`)
		sb.WriteRune('\n')
		for _, l := range locations {
			sb.WriteString(fmt.Sprintf(`"%v", %v`,
				l.LocationName, l.CaseCount))
			sb.WriteRune('\n')
		}
	case "json":
		b, err := json.MarshalIndent(locations, "", "  ")
		if err != nil {
			return "", err
		}
		sb.Write(b)
	//case "geojson":
	//	fc := geojson.NewFeatureCollection()
	//	for _, l := range locations {
	//		var feature geojson.Feature
	//		feature.Geometry = l.LocationCentrePoint
	//		feature.SetProperty("LocationName", l.LocationName)
	//		feature.SetProperty("CaseCount", l.CaseCount)
	//		fc.AddFeature(&feature)
	//	}
	//	b, err := fc.MarshalJSON()
	//	if err != nil {
	//		return "", err
	//	}
	//	sb.Write(b)
	//	sb.WriteRune('\n')
	}
	return sb.String(), nil
}
