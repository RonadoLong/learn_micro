package main

import (
	"fmt"
	"sync"
)

// var mu sync.Mutex
// var chain string

// func main() {
// 	chain = "main"
// 	A()
// 	fmt.Println(chain)
// }
// func A() {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	chain = chain + " --> A"
// 	B()
// }
// func B() {
// 	chain = chain + " --> B"
// 	C()
// }
// func C() {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	chain = chain + " --> C"
// }

// var mu sync.RWMutex
// var count int

// func main() {
// 	go A()
// 	time.Sleep(2 * time.Second)
// 	mu.Lock()
// 	defer mu.Unlock()
// 	count++
// 	fmt.Println(count)
// }
// func A() {
// 	mu.RLock()
// 	defer mu.RUnlock()
// 	B()
// }
// func B() {
// 	time.Sleep(5 * time.Second)
// 	C()
// }
// func C() {
// 	mu.RLock()
// 	defer mu.RUnlock()
// }

// func main() {
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		time.Sleep(time.Second)
// 		wg.Done()
// 		wg.Add(1)
// 	}()
// 	wg.Wait()
// }

// type MyMutex struct {
// 	count int
// 	sync.Mutex
// }

// func main() {
// 	var mu MyMutex
// 	mu.Lock()
// 	var mu2 = mu
// 	mu.count++
// 	mu.Unlock()
// 	mu2.Lock()
// 	mu2.count++
// 	mu2.Unlock()
// 	fmt.Println(mu.count, mu2.count)
// }

// func main() {
// 	var wg sync.WaitGroup
// 	wg.Add(2)
// 	var ints = make([]int, 0, 1000)
// 	go func() {
// 		for i := 0; i < 1000; i++ {
// 			ints = append(ints, i)
// 		}
// 		wg.Done()
// 	}()
// 	go func() {
// 		for i := 0; i < 1000; i++ {
// 			ints = append(ints, i)
// 		}
// 		wg.Done()
// 	}()
// 	wg.Wait()
// 	fmt.Println(len(ints))
// }

// func main() {
	var m sync.Map
// 	m.LoadOrStore("a", 1)
// 	m.Delete("a")
// 	fmt.Println(m.l)
// }


func main() {
	var ch chan int
	var count int
	go func() {
		ch <- 1
	}()
	go func() {
		count++
		close(ch)
	}()
	<-ch
	fmt.Println(count)
}
