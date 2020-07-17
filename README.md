# cidrmerge

Package cidrmerge merges IPv4 and IPv6 networks.

https://pkg.go.dev/github.com/thcyron/cidrmerge

## Example

```go
package main

import (
    "fmt"
    "net"

    "github.com/thcyron/cidrmerge"
)

func main() {
    _, ipNet1, _ := net.ParseCIDR("192.168.0.0/24")
    _, ipNet2, _ := net.ParseCIDR("192.168.1.0/24")
    merged := cidrmerge.Merge([]*net.IPNet{ipNet1, ipNet2})
    fmt.Println(merged) // Prints: [192.168.0.0/23]
}
```
