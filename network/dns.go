package network

import "net"

func HostIsResolvable(host string) bool {
	_, err := net.LookupIP(host)
	return err == nil
}
