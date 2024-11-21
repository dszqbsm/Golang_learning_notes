package main

import (
	"fmt"
	"strconv"
	"unsafe"
)

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func main() {
	var a = 1
	c := float32(a)
	// 整数转换成字符串
	i := 123
	// 使用fmt.Sprintf返回一个格式化的字符串
	s1 := fmt.Sprintf("%d", i)
	// 使用strconv.Itoa返回一个字符串
	s2 := strconv.Itoa(i)
	// 使用strconv.FormatInt返回一个字符串
	s3 := strconv.FormatInt(int64(i), 10)
	// 使用strconv.FormatUint返回一个字符串
	s4 := strconv.FormatUint(uint64(i), 10)

	// 字符串转换成整数
	s := "123"
	// 使用strconv.Atoi返回一个整数
	i1, _ := strconv.Atoi(s)
	// 使用strconv.ParseInt返回一个整数
	i2, _ := strconv.ParseInt(s, 10, 64)
	// 使用strconv.ParseUint返回一个无符号整数
	i3, _ := strconv.ParseUint(s, 10, 64)

	// 字符串转换成浮点数
	s = "123.456"
	// 使用strconv.ParseFloat返回一个浮点数
	f, _ := strconv.ParseFloat(s, 64)
}
