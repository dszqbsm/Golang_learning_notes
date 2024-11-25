package main

func main() {
	var a []int
	s = make([]int, 1) // make函数的返回值就是一个切片，分配一个内存
	s[0] = 100
	s[1] = 200

	var m map[string]int
	m = make(map[string]int, 1) // make函数的返回值就是一个映射，分配两个内存
	m["zhangsan"] = 100

	m["zhangsi"] = 100
}
