package main

import (
	"fmt"
	"io"
	"sort"
)

type resource struct {
	URL   string
	count int
}

type report struct {
	baseURL   string
	resources []resource
}

func newReport(pages map[string]int, baseURL string) *report {
	resources := make([]resource, len(pages))
	i := 0

	for u, count := range pages {
		resources[i] = resource{
			URL:   u,
			count: count,
		}
		i++
	}

	return &report{
		baseURL:   baseURL,
		resources: resources,
	}
}

func (r *report) sort() {
	sort.Slice(r.resources, func(i, j int) bool {
		return r.resources[i].count > r.resources[j].count
	})
}

func (r *report) print(out io.Writer) {
	_, _ = fmt.Fprintf(out, `
=============================
  REPORT for %s
=============================
`, r.baseURL)

	for _, res := range r.resources {
		_, _ = fmt.Fprintf(out, "Found %d internal links to %s\n", res.count, res.URL)
	}
}
