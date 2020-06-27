package network

import (
	"sync/atomic"
	"time"
)

var atomicNetId int64 = time.Now().UnixNano() / 1e6

func getAndIncrement() int64 {
	return atomic.AddInt64(&atomicNetId, 1)
}
