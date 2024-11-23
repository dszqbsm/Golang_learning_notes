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
	a := [5]int{1, 2, 3, 4, 5}

	s := a[1:3:4] // 额外指定max，控制切片容量

	fmt.Printf("s: %v len: %v cap: %v\n", s, len(s), cap(s)) // [2, 3] len: 2 cap: 4

	s1 := a[1:3]

	fmt.Printf("s1: %v len: %v cap: %v\n", s1, len(s1), cap(s1)) // [2, 3] len: 2 cap: 4

}
