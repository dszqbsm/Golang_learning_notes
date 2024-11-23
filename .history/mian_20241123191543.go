package main

import (
	"fmt"
	"unsafe"
)

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func testArray(a [2]int) {
	a[0] = 100
	fmt.Println(a)
}

func main() {
	a := [2]int{1, 2}
	testArray(a)
	fmt.Println(a)
}
