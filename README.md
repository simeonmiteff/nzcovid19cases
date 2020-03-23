# NZ COVID-19 cases scraper

This code is intended to scrape the New Zealand Ministry Of Health COVID-19 [case page](https://www.health.govt.nz/our-work/diseases-and-conditions/covid-19-novel-coronavirus/covid-19-current-cases)
and the NZ government COVID-19 alert level from [this page](https://covid19.govt.nz/government-actions/covid-19-alert-level/) and render the data in various formats suitable for mapping, visualisation and analysis.

Use this with caution - the NZ government may change their page and break the scraper at any time.

This code is used as the core of an API service I'm running: https://nzcovid19api.xerra.nz/

## Building

Building requires a go 1.13+ toolchain.

`./build.sh`

## Usage

For now there is a CLI tool:

```
cmd/nzcovid19-cli$ ./nzcovid19-cli 

Usage: ./nzcovid19-cli <action>
        Where <action> is one of:
                - cases/json
                - cases/csv
                - cases/geojson
                - locations/json
                - locations/csv
                - locations/geojson
                - alertlevel/json
```

## Notes about the data

- The GeoJSON renderer looks up a coordinate based on the "Location" column from the table. Currently it has a look-up table that needs to be manually updated when new locations appear. This needs to be replaced with something more reliable.

## Code license

This code is published under the [MIT license](LICENSE.txt).

## Data copyright

The data processed by this tool is published under:
 - The [Ministry Of Health's copyright](https://www.health.govt.nz/about-site/copyright), which at the time
of writing is Creative Commons Attribution 4.0 International Licence, with some exceptions.
 - The [Crown Copyright](https://www.iponz.govt.nz/about-ip/copyright/crown-copyright/).