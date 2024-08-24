package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("request creation error: %s", err)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 399 {
		return "", fmt.Errorf("invalid status (%d): %s", res.StatusCode, res.Status)
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("invalid content type: %s", contentType)
	}

	rawHTML, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response body: %s", err)
	}

	return string(rawHTML), nil
}
