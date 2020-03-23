package main

import (
	"github.com/simeonmiteff/nzcovid19cases"
)

func main() {
	rawCases, err := nzcovid19cases.ScrapeCases()
	if err != nil {
		panic(err)
	}
	normCases, err := nzcovid19cases.NormaliseCases(rawCases)
	if err != nil {
		panic(err)
	}
	_, err = nzcovid19cases.RenderCases(normCases, "csv")
	if err != nil {
		panic(err)
	}
	locations := nzcovid19cases.BuildLocations(normCases)
	_, err = nzcovid19cases.RenderLocations(locations, "csv")
	if err != nil {
		panic(err)
	}
}