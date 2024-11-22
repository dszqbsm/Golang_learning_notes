package main

import (
	"fmt"
	"unsafe"
)

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func main() {
	// 使用goto跳出双层for循环
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j == 20 {
				// 跳转到退出标签
				goto breakTag
			}
			fmt.Println(i, j)
		}
	}

breakTag:
	fmt.Println("breakTag")
}
