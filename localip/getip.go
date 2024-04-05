// Package localip provides functionality to retrieve local IP addresses.
package localip

import (
	"net"
)

// IPAddrInfo holds information about an IP address
// including its string representation and version (IPv4 or IPv6).
type IPAddrInfo struct {
	Address string
	Version string
}

// GetLocalIPAddr allows to get local IPv4 or IPv6 address
// GetLocalIPAddr() -> IPAddrInfo{"10.100.1.10", "IPv4"}
// getLocalIPAddr() -> IPAddrInfo{"fe80::2fb1:53c4:d88e:241", "IPv6"}
func GetLocalIPAddr() (ips []IPAddrInfo, err error) {
	var ipAddresses []IPAddrInfo

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)

			if !ok || ipnet.IP.IsLoopback() {
				continue
			}

			if ipnet.IP.To4() != nil {
				ipAddresses = append(ipAddresses, IPAddrInfo{Address: ipnet.IP.String(), Version: "IPv4"})
			} else if ipnet.IP.To16() != nil {
				ipAddresses = append(ipAddresses, IPAddrInfo{Address: ipnet.IP.String(), Version: "IPv6"})
			}
		}

	}
	return ipAddresses, nil
}
