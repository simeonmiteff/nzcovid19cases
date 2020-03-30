package nzcovid19cases

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"strconv"
	"strings"
)

type Cluster struct {
	Name  string
	Cases int
}

func ScrapeClusters() ([]*Cluster, error) {
	var clusters []*Cluster
	resp, err := soup.Get("https://www.health.govt.nz/our-work/diseases-and-conditions/covid-19-novel-coronavirus/covid-19-current-situation/covid-19-current-cases/covid-19-clusters")
	if err != nil {
		return clusters, err
	}
	doc := soup.HTMLParse(resp)

	tables := doc.FindAll("table")
	if len(tables) < 1 {
		return clusters, fmt.Errorf("page must have at least one table")
	}

	trs := tables[0].FindAll("tr")
	clusters = make([]*Cluster, len(trs)-1)

	for i, tr := range trs[1:] {
		var c Cluster
		tds := tr.FindAll("td")
		c.Name = strings.TrimSpace(tds[0].Text())
		c.Cases, err = strconv.Atoi(strings.TrimSpace(tds[1].Text()))
		if err != nil {
			return clusters,fmt.Errorf("on row %v could not convert case count (%v) to int", i, tds[1].Text())
		}
		clusters[i] = &c
	}

	return clusters, nil
}

func RenderClusters(clusters []*Cluster, viewType string) (string, error) {
	var sb strings.Builder
	validViewTypes := map[string]bool{
		"csv":     true,
		"json":    true,
	}
	if !validViewTypes[viewType] {
		return "", InvalidUsageError{fmt.Sprintf("unknown view type: %v", viewType)}
	}

	switch viewType {
	case "json":
		b, err := json.MarshalIndent(clusters, "", "  ")
		if err != nil {
			return "", err
		}
		sb.Write(b)
	case "csv":
		sb.WriteString(`"Name", "Cases"`)
		sb.WriteRune('\n')
		for _, c := range clusters {
			sb.WriteString(fmt.Sprintf(`"%v", %v`, c.Name, c.Cases))
			sb.WriteRune('\n')
		}
	}
	return sb.String(), nil
}