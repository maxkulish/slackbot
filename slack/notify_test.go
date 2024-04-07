package slack

import (
	"testing"

	"github.com/maxkulish/slackbot/localip"
)

func TestPrepareIPList(t *testing.T) {
	singleIP := []localip.IPAddrInfo{
		{
			Address: "192.168.1.1",
			Version: "IPv4",
		},
	}

	manyIP4s := []localip.IPAddrInfo{
		{
			Address: "192.168.1.1",
			Version: "IPv4",
		},
		{
			Address: "127.0.0.1",
			Version: "IPv4",
		},
		{
			Address: "14.89.76.251",
			Version: "IPv4",
		},
	}

	mixedIPs := []localip.IPAddrInfo{
		{
			Address: "192.168.1.1",
			Version: "IPv4",
		},
		{
			Address: "fe80::8811:86eb:a566:5033",
			Version: "IPv6",
		},
	}

	cases := []struct {
		desc string
		ips  []localip.IPAddrInfo
		want string
	}{
		{"empty IP list", []localip.IPAddrInfo{}, "`unknown`"},
		{"single IP", singleIP, "`192.168.1.1`"},
		{"multiple IPs, IPv4 and IPv6 mix", mixedIPs, "`192.168.1.1`"},
		{"multiple IPv4 IPs", manyIP4s, "`192.168.1.1`, `127.0.0.1`, `14.89.76.251`"},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			got := PrepareIPList(c.ips)
			if got != c.want {
				t.Errorf("PrepareIPList(%v) == %q, want %q", c.ips, got, c.want)
			}
		})
	}
}
