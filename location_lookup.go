package nzcovid19cases

import (
	"fmt"
	gj "github.com/paulmach/go.geojson"
)

var locations = map[string]*gj.Geometry{
	"Auckland":     gj.NewPointGeometry([]float64{174.7633, -36.8485}),
	"Canterbury":   gj.NewPointGeometry([]float64{171.1637, -43.7542}),
	"Dunedin":      gj.NewPointGeometry([]float64{170.5028, -45.8788}),
	"Hawkes Bay":   gj.NewPointGeometry([]float64{176.7416, -39.1090}),
	"Invercargill": gj.NewPointGeometry([]float64{168.3538, -46.4132}),
	"Northland":    gj.NewPointGeometry([]float64{173.7624, -35.5795}),
	"Queenstown":   gj.NewPointGeometry([]float64{168.6626, -45.0312}),
	"Rotorua":      gj.NewPointGeometry([]float64{176.2497, -38.1368}),
	"Southern DHB": gj.NewPointGeometry([]float64{170.5086, -45.8694}), // Dunedin coordinates?
	"Taranaki":     gj.NewPointGeometry([]float64{174.4383, -39.3538}),
	"Waikato":      gj.NewPointGeometry([]float64{175.1894, -37.4558}),
	"Wellington":   gj.NewPointGeometry([]float64{174.7762, -41.2865}),
}

func GetLocation(location string) (*gj.Geometry, error) {
	geometry, ok := locations[location]
	if !ok {
		return nil, fmt.Errorf("unknown location: %v", geometry)
	}
	return geometry, nil
}