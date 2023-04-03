package wood

import (
	"os"
	"testing"
)

func TestFunc(t *testing.T) {
	logger := New(os.Stdout)
	logger.Trace("xxx")
	Trace("ddd")
}
