package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Student struct {
	// 不同于json，这里使用的是info作为标签
	Name string `info:"name"`
	Age  int    `info:"age"`
}

func (s Student) Study(title string) {
	fmt.Printf("%s同学正在学习%s\n", s.Name, title)
}

func (s Student) Play(hours int) {
	fmt.Printf("%s同学玩了%d小时\n", s.Name, hours)
}

// 加载数据至变量v
// 思路是：一行一行读取文本文件的内容，并把每行内容都按照等号分割为键值对，然后按键名称去结构体变量中找到对应的结构体字段，并为其赋值
func LoadInfo(s string, v interface{}) (err error) {
	// tInfo为v的类型对象，调用Kind方法获取底层类型
	tInfo := reflect.TypeOf(v)
	if tInfo.Kind() != reflect.Ptr {
		err = errors.New("Please pass into a struct ptr")
		return
	}
	if tInfo.Elem().Kind() != reflect.Struct {
		err = errors.New("Please pass into a struct ptr")
		return
	}

	vInfo := reflect.ValueOf(v)
	// 按行分隔
	list := strings.Split(s, "\n")
	for _, item := range list {
		// 按等号拆分为键值对
		kvList := strings.Split(item, "=")
		if len(kvList) != 2 {
			continue
		}
		fieldName := ""
		key := strings.TrimSpace(kvList[0])
		value := strings.TrimSpace(kvList[1])
		// 遍历结构体字段的tag找到对应的key
		for i := 0; i < tInfo.Elem().NumField(); i++ {
			f := tInfo.Elem().Field(i)
			tagVal := f.Tag.Get("info")
			if tagVal == key {
				fieldName = f.Name
				break
			}
		}
		if len(fieldName) == 0 {
			continue // 找不到跳过
		}
		// 根据找到的结构体字段名称找到结构体的字段
		fv := vInfo.Elem().FieldByName(fieldName)
		switch fv.Type().Kind() {
		case reflect.String:
			fv.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			fv.SetInt(intVal)
		default:
			return fmt.Errorf("unsupport value type: %v", fv.Type().Kind())
		}
	}
	return
}

// Do调用变量v的name方法
// 实现思路：根据入参的方法名name在结构体变量的方法集中查找到对应的方法，然后再调用该方法
func Do(v interface{}, name string, arg interface{}) {
	tInfo := reflect.TypeOf(v)
	vInfo := reflect.ValueOf(v)

	fmt.Println(tInfo.NumMethod())
	m := vInfo.MethodByName(name)
	if !m.IsValid() || m.IsNil() {
		fmt.Printf("%s没有%s方法\n", tInfo.Name(), name)
		return
	}

	// 调用指定方法，通过反射调用方法传递的参数必须是[]reflect.Value类型的
	argVal := reflect.ValueOf(arg)
	m.Call([]reflect.Value{argVal})
}

func main() {
	var stu Student
	// 从文本文件中读取内容
	s, err := os.ReadFile("info.txt")
	if err != nil {
		panic(err)
	}
	err = LoadInfo(string(s), &stu)
	fmt.Printf("stu:%#v err:%v\n", stu, err)
}
