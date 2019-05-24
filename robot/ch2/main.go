package main

import (
	"fmt"
	"math/rand"
	//	"runtime"
	"sync"
	"time"
)

func main() {
	st := time.Now()
	wg := &sync.WaitGroup{}

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			// 防止被编译器优化，随机化初始值
			x := int(rand.Int31())
			for i := 0; i < 5000000; i++ {
				// 防止被编译器优化，随机化操作
				if (x+i)%2 == 0 {
					x += i
				} else {
					x -= i
				}
 
			}
			fmt.Println(x)
			wg.Done()
		}()
	}
 
	wg.Wait()
 
	fmt.Println("time used: ", time.Now().Sub(st))
}



