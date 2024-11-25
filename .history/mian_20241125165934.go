package main

import (
	"errors"
	"fmt"
	"unsafe"
)

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func do(s string) (func(int, int) int, error) {
	switch s {
	case "+":
		return func(i, j int) int { return i + j }, nil
	case "-":
		return func(i, j int) int { return i - j }, nil
	case "*":
		return func(i, j int) int { return i * j }, nil
	default:
		err := errors.New("invalid operator")
		return nil, err
	}
}

func main() {
	f, err := do("+")
	if err != nil {
		fmt.Println(err)
		return
	}
	res := f(1, 2)
	fmt.Println(res)
}
