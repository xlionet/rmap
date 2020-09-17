package rmap

import (
	"log"
	"testing"
	"time"
)

func Test_Rmap(t *testing.T) {
	r := New(2)
	go func() {
		tick := time.NewTicker(time.Second)

		for range tick.C {
			r.Clean()
			// r.sw()
			log.Printf("index %d len %d", r.curIndex, len(r.buckets[r.curIndex].m))
		}
	}()

	for i := 0; i < 10000000; i++ {
		r.Set(i, i)
	}
	t.Log("down")
}
