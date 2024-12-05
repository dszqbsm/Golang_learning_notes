package main

import (
	"fmt"
	"sync"
)

// demo1 通道误用导致的bug
func demo1() {
	wg := sync.WaitGroup{}
	ch := make(chan int, 10)
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)

	wg.Add(3)
	for j := 0; j < 3; j++ {
		go func() {
			for {
				task := <-ch
				fmt.Println(task)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
