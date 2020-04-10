package nzcovid19cases

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Location struct {
	LocationName string
	CaseCount int
}

func BuildLocations(normCases []*NormalisedCase) map[string]*Location {
	locations := make(map[string]*Location)

	for _, c := range normCases {
		loc, ok := locations[c.LocationName]
		if !ok {
			locations[c.LocationName] = &Location{
				LocationName: c.LocationName,
				CaseCount:    1,
			}

			continue
		}
		loc.CaseCount++
	}

	return locations
}

func RenderLocations(locations map[string]*Location, viewType string) (string, error) {
	validViewTypes := map[string]bool{
		"csv":  true,
		"json": true,
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
	}

	return sb.String(), nil
}
