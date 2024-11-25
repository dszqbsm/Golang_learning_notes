package main

import (
	"unsafe"
)

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func main() {
	// 默认初始值为nil
	m1 := map[string]string

	m1["zhangsan"] = "zhangsan" // 报错

	m2 := map[string]string{
		"zhangsan": "zhangsan",
		"lisi":     "lisi",
	}

	// 默认初始值不为nil
	m3 := make(map[string]string)

	m3["zhangsan"] = "zhangsan" // 正常
}
