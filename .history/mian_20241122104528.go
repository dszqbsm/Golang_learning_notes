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
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}

	// 省略初始语句
	i := 10
	for ; i < 20; i++ {
		fmt.Println(i)
	}

	// 省略初始语句和结束语句，实现while循环
	for i < 30 {
		fmt.Println(i)
		i++
	}

	// 使用break、goto、return、panic语句强制退出for循环，或者通过continue跳出本次循环进入下一次循环
	for i := 0; i < 10; i++ {
		if i == 5 {
			break
		}
		if i == 3 {
			continue
		}
		if i == 7 {
			panic("panic")
		}

	}

}
