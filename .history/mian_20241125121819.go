package main

import (
	"fmt"
	"sort"
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
	fmt.Println(len(m)) // 2
	// 添加元素
	m["zhangsi"] = "zhangsi"
	// 访问元素
	fmt.Println(m["zhangsan"])
	fmt.Println(len(m)) // 3

	// 删除元素
	delete(m, "zhangsi")

	// 判断键值对是否存在
	value, ok := m["zhangsi"]

	// 按顺序遍历映射
	// 先取出所有key存入切片中
	keys := make([]string, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}
	// 对切片进行排序
	sort.Strings(keys)
	// 按顺序遍历映射
	for _, key := range keys {
		fmt.Println(key, m[key])
	}
}
