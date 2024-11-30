package main

import (
	"fmt"
)

type Student1 struct {
	ID int
	Name string
}

type Student2 struct {
	ID int `json:"stu_id"`
	Name string `json:"stu_name"`
}

s1 := Student1{
	ID: 1,
	Name: "zhangsan",
}

s2 := Student2{
	ID: 2,
	Name: "zhangssi",
}

func mian() {
	// 没有设置标签的，在序列化时默认使用结构体字段名称
	b1, err := json.Marshal(s1)
	fmt.Println(b1)		// json:{"ID":1, "Name":"zhangsan"}
	// 设置标签的，在序列化时使用标签指定的字段名称
	b2, err := json.Marshal(s2)
	fmt.Println(b2)		// json:{"ID":2, "stu_name":"zhangssi"}

}


