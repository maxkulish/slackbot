// Package localip provides functionality to retrieve local IP addresses.
package localip

import (
	"bytes"
	"io"
	"net"
	"net/http"
)

// IPAddrInfo holds information about an IP address
// including its string representation and version (IPv4 or IPv6).
type IPAddrInfo struct {
	Address string
	Version string
	Local   bool
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
				ipAddresses = append(ipAddresses, IPAddrInfo{
					Address: ipnet.IP.String(),
					Version: "IPv4",
					Local:   true,
				})
			} else if ipnet.IP.To16() != nil {
				ipAddresses = append(ipAddresses, IPAddrInfo{
					Address: ipnet.IP.String(),
					Version: "IPv6",
					Local:   true,
				})
			}
		}

	}
	return ipAddresses, nil
}

// GetPublicIPAddr sends a request to checkip.amazonaws.com and returns the public IP address as a string.
func GetPublicIPAddr() (IPAddrInfo, error) {
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return IPAddrInfo{}, err
	}
	defer resp.Body.Close() // Ensure we close the response body

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return IPAddrInfo{}, err
	}

	netIP := net.ParseIP(string(bytes.TrimSpace(body)))
	if netIP == nil {
		return IPAddrInfo{}, err
	}

	ipType := "IPv6"
	if netIP.To4() != nil {
		ipType = "IPv4"
	}

	return IPAddrInfo{
		Address: netIP.String(),
		Version: ipType,
		Local:   false,
	}, nil // Return the IPAddrInfo struct and error as nil
}
