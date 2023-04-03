package main

import (
	"github.com/baa-god/sharp/wood"
	"os"
)

func main() {
	logger := wood.New(os.Stdout)
	logger.Trace("xxx")
	wood.Trace("ddd")
}
