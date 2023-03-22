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

func IsIntranet(ip string) bool {
	return MustSearch(ip).City == "内网IP"
}

func InitIP2Region(dbPath string) (err error) {
	// 从 dbPath 加载 VectorIndex 缓存，把下述 vIndex 变量全局到内存里面。
	vIndex, err := xdb.LoadVectorIndexFromFile(dbPath)
	if err != nil {
		fmt.Printf("failed to load vector index from `%s`: %s\n", dbPath, err)
		return
	}

	// 用全局的 vIndex 创建带 VectorIndex 缓存的查询对象
	searcher, err = xdb.NewWithVectorIndex(dbPath, vIndex)
	if err != nil {
		fmt.Printf("failed to create searcher with vector index: %s\n", err)
		return
	}

	return
}
