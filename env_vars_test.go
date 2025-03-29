package pipe

import (
	"os"
	"testing"
)

func longTestCase(t *testing.T) {
	if os.Getenv("INCLUDE_LONG") != "true" {
		t.Skip()
	}
}
