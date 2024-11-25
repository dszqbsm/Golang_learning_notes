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
	// 值为切片类型的映射
	// 键为字符串，值为字符串切片的映射
	m := make(map[string][]stirng, 3)
	key := "zhangsan"
	value, ok := m[key]
	if !ok {
		fmt.Println("not found")
	}
	// 找到的话，值就是一个切片
	value = append(value, "lisi")
	m[key] = value

}
