package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCrawler(t *testing.T) {
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

	c := newCrawler(ts.URL, log.New(os.Stdout, "", log.LstdFlags))
	c.crawlPage(ts.URL)
	c.wg.Wait()

	if hasRequestedExternalURL {
		t.Errorf("did not expect external url to be crawled")
	}

	if len(c.pages) != 2 {
		t.Errorf("unexpected numbers of pages crawled. expected: 1. got: %d. pages: %v", len(c.pages), c.pages)
	}
}
