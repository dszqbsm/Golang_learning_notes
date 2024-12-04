package main

import (
	"fmt"
	"time"
)

func hello() {
	fmt.Println("Hello World")
}

func main() {
	// 串行执行
	hello()
	fmt.Println("你好")

	// 并发执行，只会给出你好，因为创建goroutine执行hello函数需要时间开销，然而main goroutine是继续执行的，main goroutine结束时所有由main goroutine创建的goroutine都会结束
	go hello()
	fmt.Println("你好")

	// 并发执行，使用time.Sleep等待一秒，你好在前hello在后，因为创建goroutine执行hello函数需要时间开销，然而main goroutine是继续执行的
	go hello()
	fmt.Println("你好")
	time.Sleep(time.Second)
}
