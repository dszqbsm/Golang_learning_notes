package main

import (
	"fmt"
	"reflect"
	"errors"
)

type Student struct {
	Name string `info:"name"`
	Age int `info:"age"`
}

// 加载数据至变量v
func LoadInfo(s string, v interface{}) (err error) {
	// 确保传入的v是结构体指针
	tInfo := reflect.TypeOf(v)
	if tInfo.Kind() != reflect.Ptr {
		err = errors.New("Please pass into a struct ptr")
		return 
	}
	if tInfo.Elem().Kind() != reflect.Struct {
		err = errors.New("Please pass into a struct ptr")
		return 
	}

	vIndo := reflect.ValueOf(v)
	// 按行分隔
	list :=
}
