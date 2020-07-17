package cidrmerge

import (
	"fmt"
	"net"

	"github.com/armon/go-radix"
)

func Merge(ipNets []*net.IPNet) []*net.IPNet {
	if len(ipNets) <= 1 {
		return ipNets
	}
	t := radix.New()
	for i, ipNet := range ipNets {
		t.Insert(binprefix(ipNet), ipNets[i])
	}
	for {
		done := true
		for _, ipNet := range ipNets {
			superIPNet := supernet(ipNet)
			if superIPNet == nil {
				continue
			}
			prefix := binprefix(superIPNet)
			found := 0
			t.WalkPrefix(prefix, func(s string, v interface{}) bool {
				found++
				return found == 2
			})
			if found == 2 {
				t.DeletePrefix(prefix)
				t.Insert(prefix, superIPNet)
				done = false
			}
		}
		if done {
			break
		}
		ipNets = make([]*net.IPNet, 0, t.Len())
		t.WalkPrefix("", func(s string, v interface{}) bool {
			ipNets = append(ipNets, v.(*net.IPNet))
			return false
		})
	}
	return ipNets
}

func binprefix(ipNet *net.IPNet) string {
	var (
		s  string
		ip = ipNet.IP
	)
	if ip4 := ip.To4(); ip4 != nil {
		ip = ip4
	}
	for _, b := range ip {
		s += fmt.Sprintf("%08b", b)
	}
	ones, _ := ipNet.Mask.Size()
	return s[0:ones]
}

func supernet(ipNet *net.IPNet) *net.IPNet {
	ones, bits := ipNet.Mask.Size()
	if ones == 0 {
		return nil
	}
	mask := net.CIDRMask(ones-1, bits)
	return &net.IPNet{
		IP:   ipNet.IP.Mask(mask),
		Mask: mask,
	}
}
