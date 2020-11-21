package localip

import (
	"net"
)

// getLocalIPAddr allows to get local IPv4 or IPv6 address
// getLocalIPAddr(v4) -> "10.100.1.10"
// getLocalIPAddr(v6) -> "fe80::2fb1:53c4:d88e:241"
func getLocalIPAddr(version string) (ips []string, err error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return []string{}, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return []string{}, err
		}

		for _, addr := range addrs {
			ipnet, status := addr.(*net.IPNet)

			if status == true {
				switch version {
				case "v4":
					// prevent 127 and check IP4.
					if !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
						ips = append(ips, ipnet.IP.String())
					}
				case "v6":
					// prevent ::1 and check IP6.
					if !ipnet.IP.IsLoopback() && ipnet.IP.To16() == nil {
						ips = append(ips, ipnet.IP.String())
					}
				}
			}
		}
	}

	return

}

func GetIPv4() (ip []string, err error) {
	return getLocalIPAddr("v4")
}

func GetIPv6() (ip []string, err error) {
	return getLocalIPAddr("v6")
}
