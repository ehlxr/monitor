package pkg

import (
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	limiter := NewLimiterServer(1*time.Second, 5)

	for {
		if limiter.IsAvailable() {
			limiter.Increase()

			t.Log("hello...", limiter.ReqCount)
		} else {
			return
		}
	}
}
