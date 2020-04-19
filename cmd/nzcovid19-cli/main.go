package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/simeonmiteff/nzcovid19cases"
)

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `

Usage: %v <action>
	Where <action> is one of:
		- cases/json
		- cases/csv
		- locations/json
		- locations/csv
		- alertlevel/json
		- casestats/json
		- clusters/json
		- clusters/csv
`, os.Args[0])

	os.Exit(1)
}

func abort(err error, exitCode int) {
	_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(exitCode)
}

func getCases() ([]*nzcovid19cases.NormalisedCase, error) {
	rawCases, err := nzcovid19cases.ScrapeCases()
	if err != nil {
		return nil, err
	}

	normCases, err := nzcovid19cases.NormaliseCases(rawCases)
	if err != nil {
		return nil, err
	}

	return normCases, nil
}

func checkTypes(dataType string) {
	validDataTypes := map[string]bool{
		"cases":      true,
		"casestats":  true,
		"locations":  true,
		"alertlevel": true,
		//"grants":     true,
		"clusters":   true,
	}

	if !validDataTypes[dataType] {
		_, _ = fmt.Fprintf(os.Stderr, "Unknown data type specified\n")

		usage()
	}
}

func getArgs() (string, string) {
	if len(os.Args) < 2 { //nolint:gomnd
		usage()
	}

	arg1 := os.Args[1]
	parts := strings.SplitN(arg1, "/", 2) //nolint:gomnd

	if len(parts) != 2 { //nolint:gomnd
		usage()
	}

	return parts[0], parts[1]
}

//nolint:funlen
func main() {
	dataType, viewType := getArgs()

	checkTypes(dataType)

	var result string

	var renderErr error

	switch dataType {
	case "cases":
		normCases, err := getCases()
		if err != nil {
			abort(err, 4)
		}

		result, renderErr = nzcovid19cases.RenderCases(normCases, viewType)
	case "locations":
		normCases, err := getCases()
		if err != nil {
			abort(err, 4)
		}

		locations := nzcovid19cases.BuildLocations(normCases)
		result, renderErr = nzcovid19cases.RenderLocations(locations, viewType)
	case "alertlevel":
		levelInt, levelString, err := nzcovid19cases.ScrapeLevel()
		if err != nil {
			abort(err, 5)
		}

		result, renderErr = nzcovid19cases.RenderLevels(levelInt, levelString, viewType)
	//case "grants":
	//	gS, gR, err := nzcovid19cases.ScrapeGrants()
	//	if err != nil {
	//		abort(err, 6)
	//	}
	//
	//	result, renderErr = nzcovid19cases.RenderGrants(gS, gR, viewType)
	case "casestats":
		caseStats, err := nzcovid19cases.ScrapeCaseStats()
		if err != nil {
			abort(err, 3)
		}

		result, renderErr = nzcovid19cases.RenderCaseStats(caseStats, "json")
	case "clusters":
		clusters, err := nzcovid19cases.ScrapeClusters()
		if err != nil {
			abort(err, 3)
		}

		result, renderErr = nzcovid19cases.RenderClusters(clusters, viewType)
	}

	if renderErr != nil {
		invalidUsageError, ok := renderErr.(nzcovid19cases.InvalidUsageError)
		if ok {
			_, _ = fmt.Fprint(os.Stderr, invalidUsageError)

			usage()
		}

		abort(renderErr, 100)
	}

	fmt.Print(result)
}
