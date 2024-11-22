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
breakLabel:
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j == 2 {
				fmt.Println(j)
				continue breakLabel //跳出内层循环，继续执行外层循环，即只会打印十个2
			}
		}
		fmt.Println(i)
	}
}
