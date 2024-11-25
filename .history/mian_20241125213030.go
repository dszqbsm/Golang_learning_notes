package main

import "fmt"

func main() {
	var s []int
	s = make([]int, 1) // make函数的返回值就是一个切片，分配一个内存
	s[0] = 100
	s = append(s, 200)
	//s[1] = 200			// 报错，因为只分配了一个内存，不会自动扩容
	fmt.Println(s[0], s[1])

	var m map[string]int
	m = make(map[string]int, 1) // make函数的返回值就是一个映射，分配两个内存
	m["zhangsan"] = 100

	m["zhangsi"] = 100 // 不会报错，映射会自动扩容

	for k, v := range m {
		fmt.Println(k, v)
	}
}
