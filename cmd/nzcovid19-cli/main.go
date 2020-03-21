package main

import (
	"fmt"
	"github.com/simeonmiteff/nzcovid19cases"
	"os"
	"strings"
)

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `

Usage: %v <action>
	Where <action> is one of:
		- cases/json
		- cases/csv
		- cases/geojson
`, os.Args[0])
	os.Exit(1)
}

func abort(err error, exitCode int) {
	_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(exitCode)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	arg1 := os.Args[1]
	parts := strings.SplitN(arg1, "/", 2)
	if len(parts) != 2 {
		usage()
	}
	dataType := parts[0]
	viewType := parts[1]

	validDataTypes := map[string]bool{
		"cases": true,
	}

	if !validDataTypes[dataType] {
		_, _ = fmt.Fprintf(os.Stderr, "Unknown data type specified\n")
		usage()
	}

	rawCases, err := nzcovid19cases.Scrape()
	if err != nil {
		abort(err, 2)
	}

	normCases, err := nzcovid19cases.NormaliseCases(rawCases)
	if err != nil {
		abort(err, 3)
	}

	var result string
	switch dataType {
	case "cases":
		result, err = nzcovid19cases.RenderCases(normCases, viewType)
	}

	if err != nil {
		invalidUsageError, ok := err.(nzcovid19cases.InvalidUsageError)
		if ok {
			_, _ = fmt.Fprint(os.Stderr, invalidUsageError)
			usage()
		}
		abort(err, 100)
	}

	fmt.Print(result)

}
