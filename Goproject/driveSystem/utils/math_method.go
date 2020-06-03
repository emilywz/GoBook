package utils

import (
	"math/rand"
	"sync"
	"time"
)

var (
	//随机数互斥锁（确保GetRandomInt不能被并发访问）
	randomMutex sync.Mutex
)

/*获取[start,end]范围内的随机数*/
func GetRandomInt(start, end int) int {
	randomMutex.Lock()
	<-time.After(1 * time.Nanosecond)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := start + r.Intn(end-start+1)
	randomMutex.Unlock()
	return n
}
