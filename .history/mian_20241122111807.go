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
				continue breakLabel //
			}
		}
		fmt.Println(i)
	}
}
