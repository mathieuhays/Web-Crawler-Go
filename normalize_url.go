package main

import (
	"errors"
	"net/url"
	"strings"
)

var errInvalidURL = errors.New("invalid URL")

func normalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", errInvalidURL
	}

	if len(u.Hostname()) == 0 {
		return "", errInvalidURL
	}

	path := strings.TrimRight(u.Path, "/")

	return strings.ToLower(u.Hostname() + path), nil
}
