package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 空指针
	var a *int
	va := reflect.ValueOf(a)
	fmt.Println(va.IsNil())   // true
	fmt.Println(va.IsValid()) // true

	// nil值
	v := reflect.ValueOf(nil)
	fmt.Println(v.IsValid()) // false
	fmt.Println(v.IsNil())   // false
}
