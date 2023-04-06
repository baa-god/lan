package main

import (
	"fmt"
	"github.com/baa-god/lan/ip2region"
)

func main() {
	ip2region.InitIP2Region("ip2region.xdb")

	r := ip2region.MustSearch("192.168.1.237")
	fmt.Printf("%+v\n", r)
}
