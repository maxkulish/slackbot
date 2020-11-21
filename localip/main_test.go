package localip

import (
	"testing"
)

func TestGetIPv4(t *testing.T) {
	f := func(version, ipv4 string) {
		t.Helper()
		ips, err := GetIPv4()
		if err != nil {
			t.Fatal("error", err)
		}

		for _, ip := range ips {
			if ip == ipv4 {
				t.Fatalf("unexpected resutl for GetIPv4. got: %s; want: %s", ip, ipv4)
			}
		}

	}

	f("v4", "127.0.0.1")
	f("v4", "0.0.0.0")
}

func BenchmarkGetIPv6(b *testing.B) {

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = GetIPv4()
	}
}
