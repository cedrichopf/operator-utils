package network

import (
	"errors"
	"fmt"
	"net/url"
)

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
