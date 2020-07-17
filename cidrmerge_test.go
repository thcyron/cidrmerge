package cidrmerge_test

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/thcyron/cidrmerge"
)

func mustParseIPNet(s string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	if ipnet == nil {
		panic("no *net.IPNet")
	}
	return ipnet
}

func TestMerge(t *testing.T) {
	testCases := map[string]struct {
		Nets   []*net.IPNet
		Merged []*net.IPNet
	}{
		"ipv4": {
			Nets: []*net.IPNet{
				mustParseIPNet("192.168.0.0/24"),
				mustParseIPNet("192.168.1.0/24"),
				mustParseIPNet("192.168.2.0/23"),
				mustParseIPNet("192.168.8.0/23"),
			},
			Merged: []*net.IPNet{
				mustParseIPNet("192.168.0.0/22"),
				mustParseIPNet("192.168.8.0/23"),
			},
		},
		"ipv6": {
			Nets: []*net.IPNet{
				mustParseIPNet("fe80:1234:5678:9abc::/64"),
				mustParseIPNet("fe80:1234:5678:9abd::/64"),
				mustParseIPNet("fe80:1234:5678:ffff::/64"),
			},
			Merged: []*net.IPNet{
				mustParseIPNet("fe80:1234:5678:9abc::/63"),
				mustParseIPNet("fe80:1234:5678:ffff::/64"),
			},
		},
		"ipv4 and ipv6": {
			Nets: []*net.IPNet{
				mustParseIPNet("192.168.0.0/24"),
				mustParseIPNet("192.168.1.0/24"),
				mustParseIPNet("fe80:1234:5678:9abc::/64"),
				mustParseIPNet("fe80:1234:5678:9abd::/64"),
			},
			Merged: []*net.IPNet{
				mustParseIPNet("192.168.0.0/23"),
				mustParseIPNet("fe80:1234:5678:9abc::/63"),
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			merged := cidrmerge.Merge(testCase.Nets)
			if !reflect.DeepEqual(merged, testCase.Merged) {
				t.Errorf("unexpected result: %v", merged)
			}
		})
	}
}

func Example() {
	_, ipNet1, _ := net.ParseCIDR("192.168.0.0/24")
	_, ipNet2, _ := net.ParseCIDR("192.168.1.0/24")
	merged := cidrmerge.Merge([]*net.IPNet{ipNet1, ipNet2})
	fmt.Println(merged)
	// Output: [192.168.0.0/23]
}
