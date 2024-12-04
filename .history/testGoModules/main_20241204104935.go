package main

import (
	"fmt"
	"reflect"
)

func reflectSetValue1(x interface{}) {
	v := reflect.ValueOf(x)
	if v.kind() == reflect.Int64 {
		fmt.Println("x is int64")
		v.SetInt(1000)
	}
}
