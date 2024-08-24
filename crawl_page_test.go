package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCrawlPage(t *testing.T) {
	hasRequestedExternalURL := false

	externalTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hasRequestedExternalURL = true
		w.WriteHeader(http.StatusOK)
	}))
	defer externalTs.Close()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		if r.URL.Path == "/sub" {
			_, _ = w.Write([]byte("<div><a href='/'>home</div>"))
			return
		}

		_, _ = w.Write([]byte(fmt.Sprintf("<div><a href='%s'>test</a><a href='/sub'>sub</a></div>", externalTs.URL)))
	}))
	defer ts.Close()

	pages := map[string]int{}
	crawlPage(ts.URL, ts.URL, pages)

	if hasRequestedExternalURL {
		t.Errorf("did not expect external url to be crawled")
	}

	if len(pages) != 2 {
		t.Errorf("unexpected numbers of pages crawled. expected: 1. got: %d. pages: %v", len(pages), pages)
	}
}
