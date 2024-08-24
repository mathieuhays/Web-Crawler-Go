package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/url"
	"strings"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return []string{}, fmt.Errorf("could not parse base url: %s", err)
	}

	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, fmt.Errorf("could not parse HTML: %s", err)
	}

	var urls []string
	var f func(node *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					rawHref := strings.TrimSpace(a.Val)
					if rawHref == "" {
						continue
					}

					href, err := url.Parse(rawHref)
					if err != nil {
						log.Printf("could not parse href '%v': %v\n", a.Val, err)
						continue
					}

					resolvedURL := baseURL.ResolveReference(href)
					urls = append(urls, resolvedURL.String())
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return urls, nil
}
