package main

import (
	"github.com/simeonmiteff/nzcovid19cases"
)

func main() {
	rawCases, err := nzcovid19cases.Scrape()
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
}