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

	s := a[1:3] // [2, 3]

	s1 := s[3:4] // 切片的切片，high可以到切片的容量

	fmt.Printf("s1: %v len: %v cap: %v\n", s1, len(s1), cap(s1)) // s1: [5] len: 1 cap: 1
}
