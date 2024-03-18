package network

import "net"

// HostIsResolvable checks if a given host is resolvable
func HostIsResolvable(host string) bool {
	_, err := net.LookupIP(host)
	return err == nil
}
