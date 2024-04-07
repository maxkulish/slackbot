package slack

import (
	"fmt"
	"strings"
	"testing"
	"time"

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

// TestPrepareMessage tests the PrepareMessage function for correctness.
func TestPrepareMessage(t *testing.T) {
	// Mock inputs
	hostname := "testHost"
	message := "Test message"
	ips := []localip.IPAddrInfo{
		{Address: "192.168.1.1", Version: "IPv4"},
		{Address: "10.0.0.1", Version: "IPv4"},
	}

	// Call the function to test
	result := PrepareMessage(hostname, message, ips)

	// Check the result
	if result.Text != "Test message" {
		t.Errorf("Unexpected result text: got %v want %v", result.Text, "Test message")
	}

	// Checking for the presence of hostname and IP list in the resulting message
	date := time.Now().Format("2006-01-02 15:04:05")
	expectedHostnameAndDate := fmt.Sprintf(":calendar: *%s*  |  :computer: %s", date, hostname)
	if !strings.Contains(result.Blocks[0].Elements[0].Text, expectedHostnameAndDate) {
		t.Errorf("Hostname and date not found in message blocks")
	}

	expectedIPList := ":information_source: *IPv4* `192.168.1.1`, `10.0.0.1`"
	if !strings.Contains(result.Blocks[1].Text.Text, expectedIPList) {
		t.Errorf("IP list not correctly formatted in message blocks")
	}

	// Check for custom message content
	if !strings.Contains(result.Blocks[3].Text.Text, message) {
		t.Errorf("Custom message not found or not correctly formatted in message blocks")
	}
}
