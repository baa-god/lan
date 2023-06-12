package ip2region

import (
	"fmt"
	"github.com/baa-god/lan/lan"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"strings"
)

var (
	searcher *xdb.Searcher
)

type Region struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Region  string `json:"region"`
	Prov    string `json:"prov"`
	City    string `json:"city"`
	Err     string `json:"err"`
}

func (r Region) IsIntranet() bool {
	return r.City == "内网IP"
}

func (r Region) Valid() bool {
	return r.Country != "" && r.Prov != "" && r.City != ""
}

func Search(ip string) (r *Region) {
	r = &Region{IP: ip}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("err:", err)
		}
	}()

	region, err := searcher.SearchByStr(ip)
	x := strings.Split(region, "|")

	for i := 0; i < len(x); i++ {
		if x[i] == "0" || x[i] == "内网IP" {
			x[i] = ""
		}
	}

	r.Country = x[0]
	r.Region = x[1]
	r.Prov = strings.TrimSuffix(x[2], "省")
	r.City = strings.TrimSuffix(x[3], "市")
	r.Err = lan.ErrOr(err)

	return
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
