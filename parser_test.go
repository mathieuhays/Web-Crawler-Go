package main

import (
	"reflect"
	"testing"
)

const htmlOne = `
<html><body>
	<a href="/path/one">
		<span>Boot.dev</span>
	</a>
	<p>
		<a href="https://example.com">Other url</a>
	</p>
	<a name="am_an_anchor">should be skipped</a>
	<div>
		<div>
			<a id="non_absolute_relative" href="path">Hey</a>
		</div>
	</div>
	<a href="">empty link</a>
</body></html>`

const htmlTwo = `<test <a href="https://example.com">hey</a`

const htmlThree = `<div><a href="https://test.com">test</a></div>`

func TestGetURLsFromHTML(t *testing.T) {
	t.Run("test valid links", func(t *testing.T) {
		urls, err := getURLsFromHTML(htmlOne, "https://boot.dev")
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		expectedURLs := []string{"https://boot.dev/path/one", "https://example.com", "https://boot.dev/path"}

		if !reflect.DeepEqual(expectedURLs, urls) {
			t.Errorf("Extracted URL not as expected. got %v -- expected: %v", urls, expectedURLs)
		}
	})

	t.Run("malformed html", func(t *testing.T) {
		urls, err := getURLsFromHTML(htmlTwo, "https://boot.dev")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if len(urls) != 0 {
			t.Errorf("unexpected amount of urls found. Expected: 0. Found: %v", urls)
		}
	})

	t.Run("minimal html", func(t *testing.T) {
		urls, err := getURLsFromHTML(htmlThree, "https://boot.dev")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		expectedURLs := []string{"https://test.com"}
		if !reflect.DeepEqual(expectedURLs, urls) {
			t.Errorf("Extracted URL not as expected. got %v -- expected: %v", urls, expectedURLs)
		}
	})
}
