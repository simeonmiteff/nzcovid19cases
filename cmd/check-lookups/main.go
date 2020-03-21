package main

import (
	"github.com/simeonmiteff/nzcovid19cases"
)

func main() {
	_, err := nzcovid19cases.RenderCases("csv")
	if err != nil {
		panic(err)
	}
}