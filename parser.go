package main

import (
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}

	var urls []string
	var f func(node *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					rawURL := strings.TrimSpace(a.Val)
					if rawURL == "" {
						continue
					}

					linkURL, err := url.Parse(rawURL)
					if err != nil {
						continue
					}
					finalUrl := a.Val

					if linkURL.Hostname() == "" {
						finalUrl = rawBaseURL + "/" + strings.TrimLeft(finalUrl, "/")
					}

					urls = append(urls, finalUrl)
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
