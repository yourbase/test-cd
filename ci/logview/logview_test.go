package logview

import (
	"fmt"
	"testing"
)

func TestViewLogLarge(t *testing.T) {
	if testing.Short() {
		return
	}
	// NOTE: The logs from this branch won't be searchable forever. We are
	// currently only searching for 7 days worth of builds.
	out, err := viewLog("241146fcba144d4ba98dca2c3c28c1ae")
	if err != nil {
		t.Fatalf("viewlog test failed: %v", err)
	}
	fmt.Print(out)
}
