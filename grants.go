package nzcovid19cases

import (
	"encoding/json"
	"fmt"
	"strings"
)

type GrantsSummary struct {
	Clients int
	Grants int
	SumGrantAmount int
}

type GrantRegions struct {
	Auckland int
	EastCoast int
	BayOfPlenty int
	Northland int
	Wellington int
	Nelson int
	Canterbury int
	Southern int
	Other int
	Total int
}

func RenderGrants(gS GrantsSummary, gR GrantRegions, viewType string) (string, error) {
	var sb strings.Builder
	if viewType != "json" {
		return "", InvalidUsageError{fmt.Sprintf("unknown view type: %v", viewType)}
	}
	b, err := json.MarshalIndent(gS, "", "  ")
	if err != nil {
		return "", err
	}
	sb.Write(b)
	b, err = json.MarshalIndent(gR, "", "  ")
	if err != nil {
		return "", err
	}
	sb.Write(b)
	return sb.String(), nil
}