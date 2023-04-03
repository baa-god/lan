package main

import (
	"github.com/baa-god/lan/wood"
	"os"
)

func main() {
	wood.Trace("xxx")

	l := wood.New(os.Stdout)
	l.Debug("ssss1")

	wood.Debug("xxx")
	l.Info("xxxx1")
}
