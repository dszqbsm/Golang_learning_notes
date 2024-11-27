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
		mysql.Open("user:password@tcp(127.0.0.1:3306)/hello"),
	)

	if !err {
		// 处理打开连接的错误
	}

	var users []User
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
	result := db.Create(&user)
	
	// user.ID 			// 返回主键 last insert id
	// result.Error		// 返回 error
	// result.RowsAffected	// 返回受影响的行数

	// 批量创建
	var users = []User{
		{Name: "zhangsan", Age: 18, Birthday: time.Now()},
		{Name: "lisi", Age: 20, Birthday: time.Now()},
	}
	db.Create(&users)
	db.CreateInBatches(&users, 100)

	for _, user := range users {
		user.ID		// 1, 2, 3	 
	}

	// 读取
	var product Product
	db.First(&product, 1)					// 查询id为1的product
	db.First(&product, "code = ?", "L1212")	// 查询code为L1212的product

	result := db.Find(&users, []int{1, 2, 3})
	result.RowsAffected						// 返回找到的记录数
	errors.Is(result.Error, gorm.ErrRecordNotFound)	// 判断是否找不到记录

	// 更新某个字段
	db.Model(&product).Update("Price", 2000)
	db.Model(&product).UpdateColumn("Price", 2000)

	// 更新多个字段
	db.Model(&Product{}).Where("price < ?", 2000).Updates(map[string]interface{}{"Price": 2000})

	// 删除 - 删除product
	db.Delete(&product)
}

// 方法1：使用数据库约束自动删除
type User struct {
	ID uint
	Name string
	Account Account				`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreditCards []CreditCard	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Orders []Order				`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// 需要使用GORM Migrate数据库迁移数据库外键才行
db.AutoMigrate(&User{})
// 如果未启用软删除，在删除User时会自动删除其依赖
db.Delete(&User{})

// 方法2：使用Select实现级联删除，不依赖数据库约束及软删除
// 删除user时，也删除user的account
db.Select("Account").Delete(&User)

// 删除user时，也删除user的Orders、CreditCards记录
db.Select("Orders", "CreditCards").Delete(&User)

// 删除user时，也删除user的Orders、CreditCards记录，也删除订单的BillingAddress
db.Select("Orders", "CreditCards", "BillingAddress").Delete(&User)

// 删除user时，也删除用户及其依赖的所有has one/many、many2many记录
db.Select(clause.Associations).Delete(&User)


