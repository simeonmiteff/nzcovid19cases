package nzcovid19cases

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"strconv"
	"strings"
)

type GrantsSummary struct {
	Clients        int
	Grants         int
	SumGrantAmount int
}

type GrantsRegions struct {
	Auckland    int
	EastCoast   int
	BayOfPlenty int
	Northland   int
	Wellington  int
	Nelson      int
	Canterbury  int
	Southern    int
	Other       int
	Total       int
}

type Grants struct {
	Summary GrantsSummary
	Regions GrantsRegions
}

func ScrapeGrants() (GrantsSummary, GrantsRegions, error) {
	var gS GrantsSummary
	var gR GrantsRegions

	resp, err := soup.Get("https://www.msd.govt.nz/about-msd-and-our-work/newsroom/2020/covid-19/covid-19-data.html")
	if err != nil {
		return gS, gR, err
	}
	doc := soup.HTMLParse(resp)
	div := doc.Find("div", "id", "content")
	if div.Error != nil {
		return gS, gR, fmt.Errorf("could not find content div")
	}

	tables := div.FindAll("table")
	if len(tables) != 2 {
		return gS, gR, fmt.Errorf("content div has fewer than 2 (%v) tables", len(tables))
	}

	tds := tables[0].FindAll("td")

	gS.Clients, err = strconv.Atoi(strings.ReplaceAll(tds[0].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert client count (%v) to int", tds[0].Text())
	}
	gS.Grants, err = strconv.Atoi(strings.ReplaceAll(tds[1].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert grant count (%v) to int", tds[1].Text())
	}
	gS.SumGrantAmount, err = strconv.Atoi(strings.ReplaceAll(tds[2].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert grant total amount (%v) to int", tds[2].Text())
	}

	tds = tables[1].FindAll("td")
	gR.Auckland, err = strconv.Atoi(strings.ReplaceAll(tds[0].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Auckland region client count (%v) to int", tds[0].Text())
	}
	gR.EastCoast, err = strconv.Atoi(strings.ReplaceAll(tds[1].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert EastCoast region client count (%v) to int", tds[1].Text())
	}
	gR.BayOfPlenty, err = strconv.Atoi(strings.ReplaceAll(tds[2].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert BayOfPlenty region client count (%v) to int", tds[2].Text())
	}
	gR.Northland, err = strconv.Atoi(strings.ReplaceAll(tds[3].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Northland region client count (%v) to int", tds[3].Text())
	}
	gR.Wellington, err = strconv.Atoi(strings.ReplaceAll(tds[4].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Wellington region client count (%v) to int", tds[4].Text())
	}
	gR.Nelson, err = strconv.Atoi(strings.ReplaceAll(tds[5].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Nelson region client count (%v) to int", tds[5].Text())
	}
	gR.Canterbury, err = strconv.Atoi(strings.ReplaceAll(tds[6].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Canterbury region client count (%v) to int", tds[6].Text())
	}
	gR.Southern, err = strconv.Atoi(strings.ReplaceAll(tds[7].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Southern region client count (%v) to int", tds[7].Text())
	}
	gR.Other, err = strconv.Atoi(strings.ReplaceAll(tds[8].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Other region client count (%v) to int", tds[8].Text())
	}
	gR.Total, err = strconv.Atoi(strings.ReplaceAll(tds[9].Text(), ",", ""))
	if err != nil {
		return gS, gR, fmt.Errorf("could not convert Total region client count (%v) to int", tds[9].Text())
	}

	return gS, gR, nil
}

func RenderGrants(gS GrantsSummary, gR GrantsRegions, viewType string) (string, error) {
	var sb strings.Builder
	if viewType != "json" {
		return "", InvalidUsageError{fmt.Sprintf("unknown view type: %v", viewType)}
	}
	b, err := json.MarshalIndent(Grants{
		Summary: gS,
		Regions: gR,
	}, "", "  ")
	if err != nil {
		return "", err
	}
	sb.Write(b)
	return sb.String(), nil
}
