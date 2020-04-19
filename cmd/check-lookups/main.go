package main

import (
	"github.com/simeonmiteff/nzcovid19cases"
)

//nolint:funlen
func main() {
	rawCases, err := nzcovid19cases.ScrapeCases()
	if err != nil {
		panic(err)
	}

	caseStats, err := nzcovid19cases.ScrapeCaseStats()
	if err != nil {
		panic(err)
	}

	_, err = nzcovid19cases.RenderCaseStats(caseStats, "json")
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

	_, err = nzcovid19cases.RenderLocations(locations, "json")
	if err != nil {
		panic(err)
	}

	levelInt, levelString, err := nzcovid19cases.ScrapeLevel()
	if err != nil {
		panic(err)
	}

	_, err = nzcovid19cases.RenderLevels(levelInt, levelString, "json")
	if err != nil {
		panic(err)
	}

	//gS, gR, err := nzcovid19cases.ScrapeGrants()
	//if err != nil {
	//	panic(err)
	//}

	//_, err = nzcovid19cases.RenderGrants(gS, gR, "json")
	//if err != nil {
	//	panic(err)
	//}

	clusters, err := nzcovid19cases.ScrapeClusters()
	if err != nil {
		panic(err)
	}

	_, err = nzcovid19cases.RenderClusters(clusters, "json")
	if err != nil {
		panic(err)
	}
}
