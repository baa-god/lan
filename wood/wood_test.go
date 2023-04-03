package wood

import (
	"os"
	"testing"
)

func TestFunc(t *testing.T) {
	Trace("xxx")
	Debug("xxx")
	//wood.Error("xxx")

	l := New(os.Stdout)
	l.Debug("ssss1")

	Debug("xxx")
	l.Panic("xxxx1")
}
