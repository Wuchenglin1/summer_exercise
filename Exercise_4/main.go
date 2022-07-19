package main

import (
	"fmt"
	"sync"
	"time"
)

type TSlice struct {
	i    []interface{}
	size int        //大小
	lock sync.Mutex //加锁，并发安全
}

func main() {
	begin := time.Now()
	var wg sync.WaitGroup
	var t TSlice

	//放入、取出 10w 个数测试
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func(n int) {
			t.push(uint(n))
			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("before : queue len = ", t.size)
	l := t.size
	for i := 0; i < l; i++ {
		fmt.Println(t.pop())
	}
	fmt.Println("after : queue len = ", t.size)
	fmt.Println("spend time : ", time.Now().Sub(begin))

	//一边拿一边取 20w 个数 性能测试:
	ch := make(chan interface{}, 200000)

	for i := 0; i < 10; i++ {
		go func(j int) {
			for k := j * 20000; k < 20000*(j+1); k++ {
				t.push(k)
			}

		}(i)
	}

	for i := 0; i < 200000; i++ {
		select {
		case ch <- t.pop():
			fmt.Println(<-ch)
		}
	}

	fmt.Println(time.Now().Sub(begin))
}
