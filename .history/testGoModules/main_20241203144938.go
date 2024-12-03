package main

import "fmt"

type Mover interface {
	Move()
}
type Dog struct {
	Name string
}

func (d *Dog) Move() {
	fmt.Println("dog move")
}

type Car struct {
	Brand string
}

func (c *Car) Move() {
	fmt.Println("car move")
}

// 当一个接口值有多个实际类型需要判断时，可以使用switch
func justifyType(x interface{}) {
	switch v := x.(type) {
	case string:
		fmt.Println("x is string, value is ", v)
	case int:
		fmt.Println("x is int, value is ", v)
	case bool:
		fmt.Println("x is bool, value is ", v)
	default:
		fmt.Println("x is unknown type")
	}
}

func main() {
	var n Mover = &Dog{Name: "dog"}
	// 类型断言语法：第一个参数是n转化为*Dog类型后的变量，第二个参数是一个布尔值，表示断言是否成功
	v, ok := n.(*Dog)
	if ok {
		fmt.Println("类型断言成功")
		v.Name = "富贵" // 变量v的类型是*Dog
	} else {
		fmt.Println("类型断言失败")
	}
}
