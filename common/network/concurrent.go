package network

import (
	"sync/atomic"
	"time"
)

var atomicNetId int64 = time.Now().UnixNano() / 1e6

func getAndIncrement(tp int64) int64 {
	return tp<<48 | atomic.AddInt64(&atomicNetId, 1)
}
