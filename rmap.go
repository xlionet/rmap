package rmap

import (
	"fmt"
)

type rmap struct {
	m    map[int]int
	free bool
}

func newBucket() rmap {
	return rmap{m: make(map[int]int), free: false}
}

func (r *rmap) clean() {

}

func (r *rmap) set() {

}

type RMap struct {
	buckets  []rmap
	curIndex int //
}

func New(bucketSize int) *RMap {
	return &RMap{
		buckets:  make([]rmap, 0, bucketSize),
		curIndex: -1,
	}
}

func (r *RMap) Set(k, v int) {
	if r.curIndex < 0 || r.curIndex >= len(r.buckets) {
		bucket := newBucket()
		r.buckets = append(r.buckets, bucket)
		r.curIndex = len(r.buckets) - 1

		fmt.Println("scale bucket", r.curIndex)
	}
	r.buckets[r.curIndex].m[k] = v
}

// 支持异步清理
func (r *RMap) Clean() {
	if r.curIndex < 0 || r.curIndex >= len(r.buckets) {
		return
	}

	oldIndex := r.curIndex

	//开始clean的时候，先释放这个map
	r.sw()

	fmt.Println("clean start", oldIndex, r.curIndex)

	for k := range r.buckets[oldIndex].m {
		r.buckets[oldIndex].m[k] = 0
	}
	r.buckets[oldIndex].free = true // 改用atomic 包
	fmt.Println("clean done", r.curIndex)
}

// 切换找到一个可用的bucket,如果没有可用的，就创建一个
func (r *RMap) sw() {
	fmt.Println("Switch cur index :", r.curIndex)

	for i := 0; i < len(r.buckets); i++ {
		index := (r.curIndex + 1) % len(r.buckets)
		if r.buckets[index].free {
			r.curIndex = index
			fmt.Println("defer Switch cur index :", r.curIndex)
			return
		}
	}

	r.curIndex = -1 //没有找到，触发重新add
	fmt.Println("defer Switch cur index :", r.curIndex)
}
