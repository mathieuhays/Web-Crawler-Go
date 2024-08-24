package main

import (
	"reflect"
	"testing"
)

func TestReportSort(t *testing.T) {
	pages := map[string]int{
		"test.com/about": 12,
		"test.com/blog":  4,
		"test.com":       24,
	}

	r := newReport(pages, "test.com")
	r.sort()

	expected := []resource{
		{
			URL:   "test.com",
			count: 24,
		},
		{
			URL:   "test.com/about",
			count: 12,
		},
		{
			URL:   "test.com/blog",
			count: 4,
		},
	}

	if !reflect.DeepEqual(r.resources, expected) {
		t.Errorf("invalid resources. expected resources sorted DESC. got: %v", r.resources)
	}
}
