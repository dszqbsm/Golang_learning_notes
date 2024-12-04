package main

import (
	"fmt"
	"sync"
)

// 声明全局等待组变量
var wg sync.WaitGroup

func hello(i int) {
	defer wg.Done() // goroutine结束就-1
	fmt.Println("Hello World", i)
}

func main() {
	for i := 0; i < 10; i++ {
		wg.Add(1) // 启动一个goroutine就+1
		go hello(i)
	}
	wg.Wait() // 等待所有goroutine结束
}
