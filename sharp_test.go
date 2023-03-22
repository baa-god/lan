package sharp

import (
	"fmt"
	"github.com/baa-god/sharp/ip2region"
	"testing"
)

func TestFunc(t *testing.T) {
	err := ip2region.InitIP2Region("ip2region.xdb")
	fmt.Println("err:", err)

	r := ip2region.MustSearch("192.168.1.241")
	fmt.Printf("r: %+v\n", r)
}
