package main

import "fmt"

type Worker struct {
	Name string
	Age  int
	Salary
}

func (w Worker) func1() {
	fmt.Println("Worker func1")
}
func (w Worker) func2() {
	fmt.Println("Worker func2")
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
	// 若没有重写func1()、func2()方法，会调用Salary的func1()、func2()方法
	w.func1()
	w.func2()
	w.Salary.func1()
	w.Salary.func2()
}
