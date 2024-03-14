package network

import (
	"fmt"
	"net/url"
)

func HostFromURL(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("unable to get host from url: %v", err)
	}
	return u.Hostname(), nil
}
