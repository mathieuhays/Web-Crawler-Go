package main

import (
	"log"
	"strings"
	"sync"
)

type crawler struct {
	pages              map[string]int
	baseURL            string
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	logger             *log.Logger
	maxPages           int
	totalQueries       int
}

func newCrawler(baseURL string, maxConcurrency, limit int, logger *log.Logger) *crawler {
	return &crawler{
		pages:              map[string]int{},
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		logger:             logger,
		maxPages:           limit,
		totalQueries:       0,
	}
}

func (c *crawler) crawlPage(rawCurrentURL string) {
	if !strings.HasPrefix(rawCurrentURL, c.baseURL) {
		c.logger.Printf("external url, skipping %q\n", rawCurrentURL)
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		c.logger.Printf("error normalizing url: %s\n", err)
		return
	}

	if !c.addPageVisit(normalizedURL) {
		return
	}

	c.wg.Add(1)

	go func() {
		defer c.wg.Done()
		defer func() { <-c.concurrencyControl }()

		c.logger.Printf("starting goroutine: %q\n", rawCurrentURL)

		c.concurrencyControl <- struct{}{}

		c.logger.Printf("executing goroutine: %q\n", rawCurrentURL)

		c.logger.Printf("querying %q\n", rawCurrentURL)

		c.mu.Lock()
		queryCount := c.totalQueries
		c.mu.Unlock()

		if c.maxPages > 0 {
			if queryCount >= c.maxPages {
				c.logger.Println("reached page limit. skipping.")
				return
			}
		}

		c.mu.Lock()
		c.totalQueries++
		c.mu.Unlock()

		rawHTML, err := getHTML(rawCurrentURL)
		if err != nil {
			c.logger.Printf("could get HTML for %q. err: %s\n", rawCurrentURL, err)
			return
		}

		urls, err := getURLsFromHTML(rawHTML, c.baseURL)
		if err != nil {
			c.logger.Printf("error extracting urls: %s\n", err)
			return
		}

		c.logger.Printf("Found %d URLs on %q\n", len(urls), rawCurrentURL)

		for _, u := range urls {
			c.crawlPage(u)
		}
	}()
}

func (c *crawler) addPageVisit(normalizedURL string) (isFirst bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.pages[normalizedURL]; !exists {
		c.pages[normalizedURL] = 1
		return true
	}

	c.logger.Printf("incrementing %q", normalizedURL)
	c.pages[normalizedURL]++
	return false
}
