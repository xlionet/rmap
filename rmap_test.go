package rmap

import (
	"testing"
	"time"
)


func Test_Rmap(t *testing.T){
	r :=New(2)
	go func(){
		tick := time.NewTicker(time.Second)

		for range tick.C{
			r.Clean()
		}
	}()

	for i:=0; i< 100000000000000; i++{
		r.Set(i, i)
	}
	t.Log("down")
}