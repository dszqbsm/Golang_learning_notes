package main

import (
	"fmt"
)

func funA() {
	fmt.Println("funA")
}

func funB() {
	defer func() {
		err := recover()
		// 如果程序出现了panic，那么可以通过recover恢复
		if err != nil {
			fmt.Printf("recover from panic: %v\n", err)
		}
	}()

	panic("funB panic")
}

func funC() {
	fmt.Println("funC")
}

func main() {
	funA()
	funB()
	funC()
}
