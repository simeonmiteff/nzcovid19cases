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

	var err error
	var result string
	switch dataType {
	case "cases":
		result, err = nzcovid19cases.RenderCases(viewType)
	default:
		usage()
	}

	if err != nil {
		invalidUsageError, ok := err.(nzcovid19cases.InvalidUsageError)
		if ok {
			_, _ = fmt.Fprint(os.Stderr, invalidUsageError)
			usage()
		}
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	fmt.Print(result)

}