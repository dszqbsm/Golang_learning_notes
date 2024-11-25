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
	m := map[string]string{
		"zhangsan": "zhangsan",
		"lisi":     "lisi",
	}
	fmt.Println(len(m))
	// 添加元素
	m["zhangsi"] = "zhangsi"
	// 访问元素
	fmt.Println(m["zhangsan"])
	fmt.Println(len(m))
}
