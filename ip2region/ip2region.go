package ip2region

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"strings"
)

var (
	searcher *xdb.Searcher
)

type Region struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	Region   string `json:"region"`
	Province string `json:"province"`
	City     string `json:"city"`
}

func (r Region) IsIntranet() bool {
	return r.City == "内网IP"
}

func Search(ip string) (r *Region, err error) {
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		fmt.Printf("failed to SearchIP(%s): %s\n", ip, err)
		return
	}

	x := strings.Split(region, "|")
	for i := 0; i < len(x); i++ {
		if x[i] == "0" {
			x[i] = ""
		}
	}

	r = &Region{IP: ip, Country: x[0], Region: x[1], Province: x[2], City: x[3]}
	return
}

func MustSearch(ip string) (r *Region) {
	if r, _ = Search(ip); r == nil {
		r = &Region{}
	}
	return r
}

func InitIP2Region(dbPath string) {
	vIndex, err := xdb.LoadVectorIndexFromFile(dbPath)
	if err != nil {
		panic(err)
	}

	searcher, err = xdb.NewWithVectorIndex(dbPath, vIndex)
	if err != nil {
		panic(err)
	}

	return
}
