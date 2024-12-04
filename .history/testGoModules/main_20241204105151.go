package main

import (
	"fmt"
	"reflect"
)

func reflectSetValue1(x interface{}) {
	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Int64 {
		fmt.Println("x is int64")
		v.SetInt(1000) // 修改的是副本，panic
	}
}

func reflectSetValue2(x interface{}) {
	v := reflect.ValueOf(x)
	// 反射中使用Elem()方法可以获取指针对应的值
	if v.Elem().Kind() == reflect.Int64 {
		fmt.Println("x is int64")
		v.Elem().SetInt(1000)
	}
}

func main() {
	var a int64 = 10
	reflectSetValue2(&a)
}
