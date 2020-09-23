package rmap

import (
	"sync"
	// "sync/atomic"
)

type rmap struct {
	m    map[int]int
	free bool
}

func newBucket() rmap {
	return rmap{m: make(map[int]int), free: false}
}

type RMap struct {
	buckets  []rmap
	curIndex int32 //
	mu       *sync.RWMutex
}

func New(bucketSize int) *RMap {
	bu := newBucket()
	return &RMap{
		buckets: []rmap{bu},
		mu:      &sync.RWMutex{},
	}
}

func (r *RMap) Set(k, v int) {
	r.mu.RLock()
	r.buckets[r.curIndex].m[k] = v
	r.mu.RUnlock()

}

// 支持异步清理
func (r *RMap) Clean() {
	if r.curIndex < 0 || int(r.curIndex) >= len(r.buckets) {
		return
	}

	oldIndex := r.curIndex

	r.mu.Lock()
	r.sw()
	r.mu.Unlock()

	for k := range r.buckets[oldIndex].m {
		r.buckets[oldIndex].m[k] = 0
	}
	r.buckets[oldIndex].free = true
}

// 切换找到一个可用的bucket,如果没有可用的，就创建一个
func (r *RMap) sw() {
	for i := int32(1); i < int32(len(r.buckets)); i++ {
		index := (r.curIndex + i) % int32(len(r.buckets))

		if r.buckets[index].free {
			r.buckets[index].free = false
			r.curIndex = index
			return
		}
	}

	bucket := newBucket()
	r.buckets = append(r.buckets, bucket)

	r.curIndex = int32(len(r.buckets) - 1)

}
