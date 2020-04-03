# NZ COVID-19 cases scraper

# Overview

This code is intended to scrape the following sources of COVID-19 data in New Zealand, and render the data in various formats suitable for mapping, visualisation and analysis:
 - Ministry Of Health COVID-19 [COVID-19 case page](https://www.health.govt.nz/our-work/diseases-and-conditions/covid-19-novel-coronavirus/covid-19-current-cases)
 - Ministry Of Social Development [COVID-19 hardship grants page](https://www.msd.govt.nz/about-msd-and-our-work/newsroom/2020/covid-19/covid-19-data.html)
 - The government COVID-19 [alert level page](https://covid19.govt.nz/government-actions/covid-19-alert-level/) 

Use this with caution - the NZ government may change their pages and break the scraper at any time.

This code is used as the core of an API service I'm running: https://nzcovid19api.xerra.nz/

Courtesy of @gizmoguy, the metrics exported are scraped by a Prometheus server, and visualised on a 
[Grafana dashboard](https://nzcovid19.grafana.sla.ac/d/r4XZV79Wz/new-zealand-covid-19-tracker?orgId=1):

![Screenshot of COVID-19 Grafana dashboard at https://nzcovid19.grafana.sla.ac/d/r4XZV79Wz/new-zealand-covid-19-tracker?orgId=1](doc/grafana_screenshot.png "Screenshot of COVID-19 Grafana dashboard")

## Building

Building requires a go 1.13+ toolchain.

`./build.sh`

## Usage

For now there is a CLI tool:

```
cmd/nzcovid19-cli$ ./nzcovid19-cli 

Usage: ./cmd/nzcovid19-cli/nzcovid19-cli <action>
	Where <action> is one of:
		- cases/json
		- cases/csv
		- locations/json
		- locations/csv
		- alertlevel/json
		- grants/json
		- casestats/json
		- clusters/json
		- clusters/csv
```

## Code license

This code is published under the [MIT license](LICENSE.txt).

## Data copyright

The data processed by this tool is published under:
 - The [Ministry Of Health's copyright](https://www.health.govt.nz/about-site/copyright) and [Ministry Of Social Development's copyright](https://www.msd.govt.nz/about-msd-and-our-work/tools/copyright-statement.html), which at the time
of writing is Creative Commons Attribution 4.0 International Licence (with some exceptions for MOH).
 - The [Crown Copyright](https://www.iponz.govt.nz/about-ip/copyright/crown-copyright/).