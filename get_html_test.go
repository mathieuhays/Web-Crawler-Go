package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetHTML(t *testing.T) {
	t.Run("invalid status", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		_, err := getHTML(ts.URL)
		if err == nil {
			t.Fatalf("error expected but none returned")
		}

		if !strings.Contains(err.Error(), "invalid status") {
			t.Errorf("unexpected error. expected to contain 'invalid status'. got: %v", err)
		}
	})

	t.Run("invalid content type", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("{\"not\":\"expected\"}")); err != nil {
				t.Fatalf("unexpected error in request handler: %s", err)
			}
		}))
		defer ts.Close()

		_, err := getHTML(ts.URL)
		if err == nil {
			t.Fatalf("error expected but none returned")
		}

		if !strings.Contains(err.Error(), "invalid content type") {
			t.Errorf("unexpected error. expected to contain 'invalid content type'. got: %v", err)
		}
	})

	t.Run("valid response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("<html><body><h1>website</h1></body></html>")); err != nil {
				t.Fatalf("unexpected error in request handler: %s", err)
			}
		}))
		defer ts.Close()

		rawHTML, err := getHTML(ts.URL)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !strings.Contains(rawHTML, "<html><body>") {
			t.Errorf("didn't get a valid response. expected to contain <html><body>. got: %s", rawHTML)
		}
	})
}
