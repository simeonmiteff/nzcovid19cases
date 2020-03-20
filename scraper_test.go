package nzcovid19cases

import "testing"

func TestScrape(t *testing.T) {
	r, err := Scrape()
	if err != nil {
		t.Error(err)
	}
	// Would be super happy if this test fails because the page has fewer than 10 cases on it
	if len(r) < 10 {
		t.Errorf("only %v cases returned, is the scraper broken, or are we out of the woods?", len(r))
	}
}