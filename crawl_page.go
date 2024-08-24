package main

import (
	"log"
	"strings"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	if !strings.HasPrefix(rawCurrentURL, rawBaseURL) {
		log.Printf("external url, skipping %q", rawCurrentURL)
		return
	}

	currentNormalized, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("error normalizing url. err: %s\n", err)
		return
	}

	if _, exists := pages[currentNormalized]; exists {
		log.Printf("incrementing %q", currentNormalized)
		pages[currentNormalized]++
		return
	}

	log.Printf("querying %q", rawCurrentURL)

	rawHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Printf("could get HTML for %q. err: %s\n", rawCurrentURL, err)
		return
	}

	pages[currentNormalized] = 1

	urls, err := getURLsFromHTML(rawHTML, rawBaseURL)
	if err != nil {
		log.Printf("error extracting urls. err: %s\n", err)
		return
	}

	log.Printf("Found %d URLs on %q", len(urls), rawCurrentURL)

	for _, u := range urls {
		crawlPage(rawBaseURL, u, pages)
	}
}
