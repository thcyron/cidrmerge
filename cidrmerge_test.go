package cidrmerge_test

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"testing"

	"github.com/thcyron/cidrmerge"
)

func TestMerge(t *testing.T) {
	testCases := []struct {
		In  []*net.IPNet
		Out []*net.IPNet
	}{
		{
			[]*net.IPNet{},
			[]*net.IPNet{},
		},
		{
			ipnets("192.168.1.0/24"),
			ipnets("192.168.1.0/24"),
		},
		{
			ipnets("192.168.0.0/24", "192.168.1.0/24"),
			ipnets("192.168.1.0/23"),
		},
		{
			ipnets("192.168.1.0/24", "192.168.2.0/24"),
			ipnets("192.168.1.0/24", "192.168.2.0/24"),
		},
		{
			ipnets("192.168.1.1/32", "192.168.1.2/32", "192.168.1.3/32"),
			ipnets("192.168.1.1/32", "192.168.1.2/31"),
		},
		{
			ipnets("fe80:1234:5678:9abc::/64", "fe80:1234:5678:9abd::/64", "fe80:1234:5678:ffff::/64"),
			ipnets("fe80:1234:5678:9abc::/63", "fe80:1234:5678:ffff::/64"),
		},
		{
			ipnets("192.168.0.0/24", "192.168.1.0/24", "fe80:1234:5678:9abc::/64", "fe80:1234:5678:9abd::/64"),
			ipnets("192.168.0.0/23", "fe80:1234:5678:9abc::/63"),
		},
	}
	for i, testCase := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			merged := cidrmerge.Merge(testCase.In)
			if !reflect.DeepEqual(merged, testCase.Out) {
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

func ipnet(s string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	if ipnet == nil {
		panic("no *net.IPNet")
	}
	return ipnet
}

func ipnets(cidrs ...string) []*net.IPNet {
	var ipnets []*net.IPNet
	for _, cidr := range cidrs {
		ipnets = append(ipnets, ipnet(cidr))
	}
	return ipnets
}
