package progress

import (
	"testing"
	"time"
)

func TestSimpleProgress(t *testing.T) {
	startTime := time.Now()

	SimpleProgress(".", 2*time.Second, func() {
		time.Sleep(5)
	})

	duration := time.Now().Sub(startTime)
	if duration >= time.Duration(6)*time.Second {
		t.Fatalf("Should take no more than 6 sec, but got: %s", duration)
	}
	if duration < time.Duration(4)*time.Second {
		t.Fatalf("Should take at least 4 sec, but got: %s", duration)
	}
}
