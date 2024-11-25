package main

import (
	"fmt"
	"os"
	"time"
)

func recordTime(funcName string) func() {
	start := time.Now()
	fmt.Printf("enter %s\n", funcName)
	return func() {
		fmt.Printf("exit %s, spend %v\n", funcName, time.Since(start))
	}
}

func readFromFile(filename string) error {
	defer recordTime("readFromFile")() // defer延迟调用
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close() // 函数结束时关闭文件
	return nil
}

func main() {
	err := readFromFile("test.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("main end")
}
