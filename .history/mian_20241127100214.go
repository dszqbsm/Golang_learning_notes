package main

import (
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"time"
)

type User struct {
	Name string
	Age int
	Birthday time.Time
}

func main() {
	db, err := gorm.Open(										
		mysql.Open("user:password@tcp(127.0.0.1:3306)/hello")	
	)

	var users []user
	err = db.Select("id", "name").Find(&users, 1).Error

	// 操作数据库
	db.AutoMigrate(&Product{})
	db.Migrator().CreateTable(&Product{})
	
	// 创建
	user := User{
		Name: "zhangsan",
		Age:  18,
		Birthday: time.Now(),
	}
}