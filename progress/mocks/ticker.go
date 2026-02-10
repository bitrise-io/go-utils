package mocks

import "time"

type mockTicker struct {
	C chan time.Time
}

func NewTicker() mockTicker {
	return mockTicker{
		C: make(chan time.Time), // Buffered channel to prevent blocking
	}
}

// Chan ...
func (t mockTicker) Chan() <-chan time.Time {
	return t.C
}

// Stop ...
func (t mockTicker) Stop() {
}

func (t mockTicker) DoTicks(n int) {
	for i := 0; i < n; i++ {
		t.C <- time.Now()
	}
}
