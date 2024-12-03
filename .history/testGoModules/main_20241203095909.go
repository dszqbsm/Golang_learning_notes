package main

import "fmt"

type Worker struct {
	Name string
	Age  int
	Salary
}

type Salary struct {
	Money int
}

func (s Salary) func1() {
	fmt.Println("Salary func1")
}
func (s Salary) func2() {
	fmt.Println("Salary func2")
}

func main() {
	s := Salary{Money: 1000}
	w := Worker{Name: "zhangsan", Age: 20, Salary: s}
	w.func1()
	w.func2()
}
