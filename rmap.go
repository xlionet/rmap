package rmap

import (
	"fmt"
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

var lock sync.RWMutex

type RMap struct {
	buckets   []rmap
	curIndex  int32 //
	lastIndex int32
}

func New(bucketSize int) *RMap {
	bu := newBucket()
	return &RMap{
		buckets:   []rmap{bu},
		lastIndex: -1,
	}
}

func (r *RMap) Set(k, v int) {
	// index := atomic.LoadInt32(&r.curIndex)
	if r.curIndex != r.lastIndex {
		fmt.Println("set last index", r.lastIndex, "new index", r.curIndex)
	}

	// fmt.Println("cur index", r.lastIndex)
	lock.RLock()
	r.buckets[r.curIndex].m[k] = v
	r.lastIndex = r.curIndex
	lock.RUnlock()

}

// 支持异步清理
func (r *RMap) Clean() {
	if r.curIndex < 0 || int(r.curIndex) >= len(r.buckets) {
		return
	}

	oldIndex := r.curIndex

	//开始clean的时候，先释放这个map
	lock.Lock()
	r.sw()
	lock.Unlock()

	fmt.Println("===>start clean index ", oldIndex, "cur index", r.curIndex)

	for k := range r.buckets[oldIndex].m {
		r.buckets[oldIndex].m[k] = 0
	}
	r.buckets[oldIndex].free = true // 改用atomic 包
	fmt.Println("===>end clean index ", oldIndex, "cur index", r.curIndex)
}

// 切换找到一个可用的bucket,如果没有可用的，就创建一个
func (r *RMap) sw() {
	fmt.Println("++++++>switch from index :", r.curIndex)
	for i := int32(1); i < int32(len(r.buckets)); i++ {
		index := (r.curIndex + i) % int32(len(r.buckets))
		if r.buckets[index].free {
			r.buckets[index].free = false
			// atomic.LoadInt32(addr * int32)

			// atomic.StoreInt32(&r.curIndex, int32(index))
			r.curIndex = index

			fmt.Println("++++++>switch new index :", r.curIndex, "len=", len(r.buckets))
			return
		}
	}

	// r.curIndex = -1 //没有找到，触发重新add
	bucket := newBucket()
	r.buckets = append(r.buckets, bucket)

	// atomic.StoreInt32(&r.curIndex, int32(len(r.buckets)-1))
	r.curIndex = int32(len(r.buckets) - 1)

	fmt.Println("scale bucket", r.curIndex)

	fmt.Println("++++++>switch to index :", r.curIndex)
}
