package pkg

import (
	"sync"
	"time"
	log "unknwon.dev/clog/v2"
)

type LimiterServer struct {
	Interval time.Duration
	MaxCount int
	Lock     sync.Mutex
	ReqCount int
}

func NewLimiterServer(interval time.Duration, maxCount int) *LimiterServer {
	limiter := &LimiterServer{
		Interval: interval,
		MaxCount: maxCount,
	}

	go func() {
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			limiter.Lock.Lock()
			log.Info("Reset LimiterServer Count...")

			limiter.ReqCount = 0
			limiter.Lock.Unlock()
		}
	}()

	return limiter
}

func (limiter *LimiterServer) Increase() {
	limiter.Lock.Lock()
	defer limiter.Lock.Unlock()

	limiter.ReqCount += 1
}

func (limiter *LimiterServer) IsAvailable() bool {
	limiter.Lock.Lock()
	defer limiter.Lock.Unlock()

	return limiter.ReqCount < limiter.MaxCount
}
