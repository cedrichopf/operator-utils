package network

import (
	"errors"
	"fmt"
	"net/url"
)

// HostFromURL returns the host of a given URL or returns an error if the host cannot be
// extracted from the URL
func HostFromURL(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("unable to get host from url: %v", err)
	}
	host := u.Hostname()
	if host == "" {
		return "", errors.New("error while getting hostname from url, got empty hostname")
	}
	return host, nil
}
