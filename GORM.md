# database/sql包

关系型数据库有很多种类型，Go种database/sql包为访问关系型数据库提供了统一的接口，即通过统一的接口去操作不同的数据库，而对于不同的数据库的具体操作，需要数据库自己实现一个driver

一个常见的问题，数据库卡死了、请求停住了没有响应，可能是因为没有关闭rows，因此初始化了一个rows之后要及时添加defer进行关闭，防御性编程，rows虽然会自己关闭，但是还是要defer，防止在读取过程中出现错误导致的提前返回，从而rows没有关闭导致服务卡死

但是rows关闭其实会出现错误的，但是从rows.Next()关闭rows，这些错误信息会丢失

```go
package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"		// 注册MySQL驱动，使得database/sql包能够使用MySQL数据库，该driver包含了数据库操作的具体实现
)

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/hello")

	rows, err := db.Query("select ic, name from users where id = ?", 1)		// ?是占位符，标识从users表中选择id和name字段，条件是id等于1
	if err != nil {
		// xxx
	}
	defer rows.Close()						// 一定要关闭rows，释放数据库资源，否则会造成资源泄露，rows虽然

	var users []User
	for rows.Next() {						// 遍历查询结果集中的每一行，rows是一个游标，Next不断获取到下一条数据
		var user User
		err := rows.Scan(&user.ID, &user.Name)			// 将当前行的数据扫描到user结构体中

		if err != nil {
			// ...
		}

		users = append(users, user)
	}

	if rows.Err() != nil {						// 只要不是和数据相关的错误，都会通过这里返回
		// ...
	}
}
```

## database/sql包设计原理

以下是database/sql包的设计原理，由Driver支持对不同数据库的连接接口、操作接口，database/sql只暴露给应用程序相同的操作接口，从而实现使用一套统一的接口对不同的数据库进行操作

连接池通过池化技术实现

> 池化操作：提前准备一些资源，在需要时可以重复使用这些预先准备的资源
> 
> - 线程池： 线程池中会启动若干数量的线程，这些线程都处于睡眠状态，当客户端有一个心的请求时，就会环形线程池中的某一个睡眠的线程，让它来处理客户端的请求，处理完之后线程又处于睡眠的状态；能够很高的提升程序的性能
> 
> - 内存池：内存池会预先分配足够大的内存，形成一个初步的内存池，然后每次用户请求内容的适合，就会返回内存池中的一块空闲的内存，并将这块内存的标志置为已使用，当内存使用完毕释放内存的时候，并不是真正地调用free或delete，而是把内存放回内存池的过程，并将标志置为空闲，当应用程序结束后才会将内存池销毁，即将内存池中的每一块内存释放；能够减少内存碎片的发生，提高了内存的使用频率，但是会造成内存的浪费，因为预先分配的内存并不一定会全部被用到
> 
> - 数据库连接池：基本思想是在系统初始化的时候将数据库连接作为对象存储在内存中，当用户需要访问数据库的时候，并非建立一个新的连接，而是从连接池中取出一个已建立的空闲连接对象，使用完毕后也不是将连接关闭，而是将连接放回到连接池中，以供下一个请求访问使用，而这些连接的建立、断开都是由连接池自身来管理的

![设计原理](./images/database.jpg)

常用的连接池配置：

- `func (db *DB) SetConnMaxIdleTime(d time.Duration)`：用于设置数据库连接池中连接的最大空闲时间
- `func (db *DB) SetConnMaxLifeTime(d time.Duration)`：用于设置数据库连接的最大生命周期，即使连接没有被关闭，也会在存活超过这个时间后被关闭和替换
- `func (db *DB) SetMaxIdleConns(n int)`：用于设置数据库连接池中的最大空闲连接数
- `func (db *DB) SetMaxOpenConns(n int)`：用于设置数据库连接池中的最大打开连接数
- `func (db *DB) Status() DBStatus`：用于获取数据库连接池的当前状态，返回一个DBStatus类型的值，包含了连接池的统计信息，如当前打开的连接数、空闲连接数

以下是database/sql包实现sql执行的伪实现过程

```go
for i := 0; i < maxBadConnRetries; i++ {			// maxBadConnRetries默认是两次
    // 从连接池获取连接或通过driver新建连接
    dc, err := db.conn(ctx, strategy)
        // 有空闲连接 -> reuse -> max life time		即复用连接池中的连接
        // 新建连接 -> max open...					即新建连接池外的新连接
    // 将连接放回连接池
    defer dc.db.putConn(dc, err, true)
        // validate Connection有无错误
        // max life time, max idle conns检查
    
    // 连接实现driver.Queryer, driver.Execer等接口
    if err == nil {
        err = dc.ci.Query(sql, args...)
    }

    isBadConn = errors.Is(err, driver.ErrBadConn)
    if !isBadConn {
        break
    }
}
```

> 在插入数据的时候，连接被数据库kill掉，即数据的插入也会进行重试，这样就会导致重复的插入

### database/sql连接接口

database/sql注册全局 driver

```go
// Driver接口
type Driver interface {
	// Open returns a new connection to the database
	Open(name string) (Conn, error)
}

// 注册全局 driver
func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("sql: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	drivers[name] = driver
}
```

业务代码中如何使用driver

```go
import _ "github.com/go-sql-driver/mysql"

func main() {
	db, err := sql.Open("mysql", "gorm:gorm@tcp(localhost:3306)/gorm?charset=utf8&parseTime=True&loc=Local")    // 建立连接
}

//注册 Driver
func init() {
	sql.Register("mysql", &mysql.MySQLDriver{})
}
```

但是会存在一些问题：

- Open中第二个字符串特别长，难以了解是什么意思，有些参数都不能通过参数的方式传进去，也很难做一些密码的转义

- 经常会忘记import driver，因为没有编译检查

而在2017年Go给出了新的连接建立的方法，支持传入interface，基于interface返回一个db，但是推出的太晚，因此很多没有用

```go
package main

import "github.com/go-sql-driver/mysql"         // 有强制的编译检查

type Connector interface {
	Connect(context.Context) (Conn, error)
	Driver() Driver
}

func OpenDB(c driver.Connector) *DB {
	
}

func main() {
	connector, err := mysql.NewConnector(&mysql.Config{
		User: "root",
		Passwd: "123456",
		Net: "tcp",
		Addr: "127.0.0.1:3306",
		DBName: "test",
		ParseTime: true,
	})

	db := sql.OpenDB(connector)
}
```

### database/sql操作接口

1. DB连接的几种类型

- 直接连接 / Conn

每次数据库操作时，都会建立一个新的数据库连接，这种方式在高并发场景下可能会导致数据库服务器压力过大，因为每个操作都需要维护一个独立的连接

- 预编译 / Stmt

先编译SQL语句，然后使用编译后的语句执行操作，能够提高性能，因为编译过程只进行一次，还能预防SQL注入攻击

- 事务 / Tx

指一系列数据库操作的集合，这些操作要么全部成功，要么全部失败，事务具有原子性、一致性、隔离性、持久性的特性，即ACID特性

2. 处理返回数据的几种方式

- Exec / ExecContext -> Result

用于执行不返回结果集的SQL语句，如INSERT、UPDATE、DELETE等，该方式返回一个Result对象，该对象包含了执行后影响的行数和生成的自增ID等信息

- Query / QueryContext -> Rows(Columns)

用于执行返回结果集的SQL查询语句，如SELECT，该方式返回一个Rows结果集，可以通过遍历Rows来获取每一行的数据，COlumns返回查询结果的列名

- QueryRow / QueryRowContext -> Row(Rows 简化) *

与Query类似，但该方法预期查询结果只有一行，如果查询结果有多行，该方法会返回错误，Row是Rows的简化版，用于处理单行结果

database/sql具体是怎么实现去解析数据库的值呢

```go
// driver通过实现这个interface来解析数据库的值
type driver.Rows interface {
	// 返回columns名字
	Columns() []string

	// 实现数据库协议
	// 解析数据到database/sql.Rows.lastcols中
	Next(dest []Value) error

	// 多批数据解析
	HasNextResultSet() bool
	NextResultSet() error

	// ...
}
```

以上主要是从源码角度解析了database/sql包的使用和实现

# GORM基础使用

GORM是一种设计简洁、功能强大、自由扩展的全功能ORM

> ORM是Object-Eelational Mapping（对象关系映射）的缩写，用于在关系数据库和对象程序语言之间转换数据，其核心思想是将数据库中的表映射到程序中的对象，使得开发者可以使用面向对象的方式来操作数据库而不用写复杂的SQL语句

![设计原则](./images/shejiyuanze.jpg)

## GORM基本用法

可以避免需要去import driver，可以避免忘记写defer关闭Rows，代码长度相比前面来说简洁不少

```go
import (
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

func main() {
	// 打开一个到MySQL数据库的连接
	db, err := gorm.Open(										// db是GORM的数据库连接对象，err是可能发生的错误
		mysql.Open("user:password@tcp(127.0.0.1:3306)/hello")	// 包含了用户名、密码、TCP协议、服务器地址、端口和数据库名
	)

	var users []user
	err = db.Select("id", "name").Find(&users, 1).Error			// 执行一个数据库查询操作，指定了查询时只选择id和name两个字段
	// Find(&users, 1)是GORM的方法，用于查询数据库并将结果填充到users切片中，1是查询条件，即查询id等于1的用户
}
```

## GORM基本使用 - CRUD

```go
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
	db.Migrator().CreateTable(&Product{})		// 自动迁移Product模型，即自动创建或修改数据库表结构以匹配Product模型的结构
	
	// 创建
	user := User{								// 直接创建Product表，如果表已存在则不执行任何操作
		Name: "zhangsan",
		Age:  18,
		Birthday: time.Now(),
	}
	result := db.Create(&user)
	// 创建一个新的User对象，并使用db.Create方法将其插入数据库。result对象包含了插入操作的结果
	// user.ID 				// 返回主键 last insert id
	// result.Error			// 返回 error
	// result.RowsAffected	// 返回受影响的行数

	// 批量创建
	var users = []User{
		{Name: "zhangsan", Age: 18, Birthday: time.Now()},
		{Name: "lisi", Age: 20, Birthday: time.Now()},
	}
	db.Create(&users)
	db.CreateInBatches(&users, 100)			// db.CreateInBatches方法也用于批量插入，这里的100参数表示每批处理100条记录

	for _, user := range users {
		fmt.Println(user.ID)				// 遍历users切片，打印每个Users对象的ID，1, 2, 3	 
	}

	// 读取
	var product Product
	db.First(&product, 1)					// 查询id为1的product
	db.First(&product, "code = ?", "L1212")	// 查询code为L1212的product

	result := db.Find(&users, []int{1, 2, 3})			// 使用db.Find方法查询id为1、2、3的User记录，并将结果填充到users切片中
	// result.RowsAffected								// 返回找到的记录数
	// errors.Is(result.Error, gorm.ErrRecordNotFound)	// 判断是否找不到记录

	// 更新某个字段
	db.Model(&product).Update("Price", 2000)
	db.Model(&product).UpdateColumn("Price", 2000)

	// 更新多个字段
	db.Model(&Product{}).Where("price < ?", 2000).Updates(map[string]interface{}{"Price": 2000})		// 更新Product表中所有price小于2000的记录，将它们的Price字段更新为2000

	// 删除 - 删除product
	db.Delete(&product)
}
```

## 模型定义 - 惯例约定

![约定](./images/yueding.jpg)

## GORM关联支持

Go中支持很多种关联

![关联](./images/guanlian.jpg)

### 关联操作 - CRUD

```GO
// 保存用户及其关联
db.Save(&User{
	Name: "zhangsan",
	Languages: []Language{{Name: "zh-CN"}, {Name: "en-US"}},
})

// 关联模式
langAssociation := db.Model(&user).Association("Languages")
// 查询关联
langAssociation.Find(&languages)
// 将汉语，英语添加到用户掌握的语言中
langAssociation.Append([]language{{languageZH, languageEN}})
// 把用户掌握的语言替换成汉语，德语
langAssociation.Replace([]language{languageZH, languageDE})
// 删除用户掌握的两个语言
langAssociation.Delete(languageZH, languageEN)
// 删除用户所有掌握的语言
langAssociation.Clear()
// 返回用户所掌握的语言的数量
langAssociation.Count()

// 批量模式 Append， Replace
var users = []User{user1, user2, user3}
langAssociation := db.Model(&users).Association("Languages")

// 批量模式 Append, Replace, 参数需要与源数据长度相同
// 例如：我们有3个user: 将userA添加到user1的Team
// 将userB添加到user2的Team，将userA、userB、userC添加到user3的Team
db.Model(&users).Association("Team").Append(&userA, &userB, &[]User{userA, userB, userC})	
```

### 关联操作 - Preload / Joins 预加载

这样就可以避免，在查询一个用户的时候，每个用户都要去查一下相关的关联，这样就会产生n+1的SQL操作，可以通过Preload方法或Joins方法来解决，虽然Preload会执行三条SQL，Joins会执行一条SQL，但是不一定三条就慢于一条，因为可能可以使用一些缓存的方式，从而三条的速度也很快，具体使用哪种方式，需要根据业务场景使用

```go
type User struct {
	Orders []Order
	Profile profile
}

// 查询用户的时候并找出其订单，个人信息（1+1条SQL）
db.Preload("Orders").Preload("Profile").Find(&users)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1, 2, 3);		// 一对多
// SELECT * FROM profiles WHERE user_id IN (1, 2, 3);	// 一对一

// 使用Join SQL 加载（单条JOIN SQL）
db.Joins("Company").Joins("Manager").First(&user, 1)
db.Joins("Company", DB,Where(&Company{Alive: true})).Find(&users)

// 预加载全部关联（只加载一级关联）
db.Preload(clause.Associations).Find(&users)

// 多级预加载
db.Preload("Orders.OrderItems.Product").Find(&users)
// 多级预加载 + 预加载全部一级关联
db.Preload("Orders.OrderItems.Product")Preload(clause.Associations).Find(&users)

// 查询用户的时候找出其未取消的订单
db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
db.Preload("Orders", "state = ?", "paid").Preload("Orders.OrderItems").Find(&users)

db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
	return db.Order("orders.amount DESC")
}).Find(&users)
```

### 关联操作 - 级联删除

为了保障进行删除的时候不会导致一些孤儿数据的产生，保证所有数据都是有用

```go
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
```

# GORM设计原理

GORM相当于在database/sql包上再加了一层，用来跟应用程序进行交互

![GORM设计原理](./images/GROM.jpg)

## SQL是怎么生成的

每一个SQL语句，都是由很多个子句生成的，而子句还由好多的表达式构成

![SQL](./images/SQL.jpg)

在GORM的代码中，执行任何一个方法的时候都会产生一个GORM STATEMENT这一个对象，而这个对象由一些Chain method和Finisher Method组成

而一个Chain Method就是用来给GROM STATEMENT添加子句的方法，而这些子句都会用来去生成最终的SQL

而一个Finisher Method则是用来决定GROM STATEMENT最终的类型并执行的方法，从而将子句拼接成一条最终的SQL语句

相当于GORM对SQL语句做了一个仿生的设计，因此可以对SQL进行很好的扩展

![GORM的SQL](./images/GORMSQL.jpg)

比如Where和Limit这两个Chain Method的具体实现，Where通过调用BuildCondition方法生成得到一些conds，再添加到GORM的子句中去，Limit也是类似的，将当前的参数翻译成GORM的子句，再添加到子句中

![GORMAPI](./images/GORMAPI.jpg)

而Finisher Method的实现，会将拿到的参数放到GORM STATEMENT目标值中，对于Find查询，则会调用一个Query的callbacks去执行，callbasks执行返回的是一个processor，然后将当前执行方式支持的子句的形式找出来，再进行翻译，编译成最终的SQL语句，编译好最后交给ConnPoll执行

![GORMFinisher](./images/GORMFinisher.jpg)

这样设计的好处

- 自定义Clause Builder：控制子句生成的SQL

可能不同的数据库或同一数据库不同版本的SQL语句是不同的，此时就可以通过自定义Clause Builder的形式来支持不同的版本，下面的例子中，初始化db的时候会先发一个请求获取版本信息初始化参数，后续再根据参数决定如何生成SQL语句，从而能够支持不同的数据库版本

![自定义Builder](./images/builder.jpg)

- 方便扩展Clause

GORM STATEMENT是由很多子句组成的，GORM在生成SQL的时候会根据当前STATEMENT的类型取出所有支持的子句，把所有子句编译成最终的SQL语句

在编译时，存在一些接口，因此可以在扩展一些编译的接口，因此只需要通过注册一些接口即可，hints就是通过这种方式提供一些支持

下面的第一个例子是在select后面加了一段子句，通过这段子句注册一个查询优化器，没有给出具体的实现，从而可以实现SQL查询优化

第二个第三个例子用于扩展from操作，比如在查询SQL的时候指定一个索引，用来加速操作

后三个例子，自由扩展内容

![扩展Clause](./images/kuozhanziju.jpg)

- 自由选择Clause

实际上不同数据库的子句语法有很大的不同，要在开发前期做技术选型的时候，就要确定好选择哪一个数据库，而不是后续再对该数据库重写一遍接口

GORM通过自定义选择子句的方式对各种数据库添加支持

![选择子句](./images/selectclause.png)

## 插件是怎么工作的

有六种模式，每一种模式都对应一组callbacks

![六种模式](./images/callbacks.jpg)

下面介绍Create这个callbacks，默认有七个callback，第一个callback是在执行create操作的时候，先开启事务，若没有在事务内则会开启当前事务，第二个callback会取出当前Method有哪些before_create的方法并执行，第三个callback会保存一些前缀关联，第四个callback将当前的转换成最终的SQL，然后将SQL发送给database进行执行，再将结果返回到应用程序，最后一个callback会去判断整体的流程是否出错

![插件工作](./images/chajian.jpg)

如何扩展API

![扩展](./images/kuozhan.jpg)

为什么要这样做，就是为了实现灵活定制，和自由扩展，比如实现多租户支持、多数据库或读写分离支持、加解密或混沌工程的支持

多租户支持实现，每个租户都有一个id，希望在查询的时候可以自动的加上这个id，而不是每次查询都要手动加上id，因为手工指定很容易漏掉某些条件，会导致一些事故

下面代码为Query、Update和Delete注册了一个新的callback，用于从当前GORM Statement中取出用户id，然后通过一个where条件更新到GORM STATEMENT中，即后续所有的查询更新删除都会加上这一个过滤条件，而create的时候也要做类似的操作，要从当前context中取出租户id，赋给我们要创建的值，这样就绑定上去，从而我们插入的每一条数据都有一个租户绑定的关系，避免了一些越权等威胁

![多租户实现](./images/duozuhu.jpg)

实现多数据库、读写分离

第一组条件，将db2作为主数据库，db3和db4作为从数据库，还指定了负载策略，主数据库指的就是写的数据库，从数据库就是读的数据库，第一组配置就是全局的默认的配置

第二组条件，只设置了一个从数据库，意味着会使用当前数据库的一个主链接作为写数据库，注册一个db5作为从数据库，并且这一组配置只是应用在User、Address这两个模型里去，即更新这两个模式的时候就会使用这一组条件

第三组条件使用db6和db7作为主数据库，使用db8作为从数据库，只会影响order、Product和secondary

这样就能实现从写数据库查询

![读写分离](./images/kufenli.jpg)

## ConnPool是什么

GORM实际上是和database/sql数据包提供的接口进行交互，但实际实现的时候不是这样的，这里增加了一层新的抽象ConnPool接口，然后由database/sql包去实现这个接口，从而GORM的实际操作就通过这一层ConnPool接口来获得和数据库的交互和实际的执行

![ConnPool](./images/ConnPool.jpg)

这样设计的好处，在读写分离应用场景中，在执行一个create update相关的callback的时候，意识到这是一个写的操作，因此会把写DB的链接替换给ConnPool，从而写相关的操作都会使用写DB进行操作，读也一样，链接替换ConnPool，从而实现使用读数据库进行执行，但是如果当前操作是一个事务，则会继续操作

![ConnPool](./images/ConnPool2.jpg)

基于ConnPool模式可以做什么事情，在开启这个模式之后，所有要执行的SQL都会交给ConnPool进行预编译，GORM会将预编译的语句缓存起来，后面再执行同类型的SQL的时候，就直接使用这些预编译的SQL，修改新的参数进行执行，这样同样的SQL只需要在数据库做一次解析，整体执行数据就会很快

GORM提供了两种预编译模式，全局模式和会话模式

![两种模式](./images/butongmoshi.jpg)

下面介绍一个网络安全的场景，要求数据存储的第一跳，要存到本国的数据库，然后才能存储到别的地方，比如字节有一些海外的业务，存在这样的场景需求

而得益于GORM的扩展架构，可以通过ConnPool扩展能力，写一个业务插件，先把数据发送到该国数据库，再把数据同步到原来的数据库，跟原来的业务代码是透明的关系，只需要在配置的时候设置配置参数就可以完成

![数据存储](./images/shujucunc.jpg)

通过ConnPool接口实现业务需求，还可以做api查询替换，通过封装一个driver，获取某个api服务返回的一些值，而不是需要从数据库中进行查询，也可以通过这个技术去构建interface，开发一个基于缓存的插件，从redis查询相应的数据，在没有查询到缓存的时候再去查询真正的数据库

回到最开始，如何通过一行配置来提升软件的性能，其实MySQL driver还提供了一个interpolateParams参数，默认是false，当该参数是false时候，在执行有参数的SQL时，会先去默认的预编译一下SQL，然后再使用预编译的SQL进行查询，在查询完成之后会将这个预编译的值进行关闭，但是这里并不能加快，因为预编译使用完之后就扔掉了，不能重复使用

这里实际对所有的SQL都会通过这种方式，也就是说执行一个查询请求实际上会发送三条请求，相当于发一个请求就要发三次，这里主要是为了解决一个数据注入的问题，但其实只是在多编码环境中才会出现的问题，对于现在大多数的服务来说都不会有这种问题，因此关掉这个参数能够提升性能

为了能够使用GORM的最佳配置，不用去翻阅文档找到潜在可以优化的配置，字节开发了一个bytedgorm库，将一些不需要用户了解的逻辑封装起来，提供一个最佳的实现，还提供了字节内部的一系列扩展

是如何实现的呢，就和一个Dialector有关，下面例子中，重构标识部分都是一个dialector，他的定义是说GORM怎么去连接数据查询后端的接口，比如postgres定义了怎么去连接一个MySQL以及postgres数据库，第三个定义了怎么连接clickhouse数据库，然后执行查询就可以从clickhouse去操作相关数据；还可以定义一组缓存插件，比如一组callback策略，找不到数据的时候就callback到MySQL数据库中进行查找

![dialector](./images/dialector.jpg)

字节内实现的bytedgorm背后实现的逻辑，dialector在初始化时可以定制一些GORM的功能，bytedgorm还提供了一些Options接口，分成两部分，一部分是在初始化db前回去调用Apply方法，会去修改配置，还提供了另外的方法，可以在初始化的后期调用db所有的api扩展定制功能

![dialector](./images/dialector2.jpg)

# GORM最佳实践

## 数据序列化与SQL表达式

- SQL表达式更新创建

不同数据库的SQL语法不同，为了支持所有的情况，GORM提供了SQL表达式的支持

![SQL表达式](./images/sqlbiaodashi.jpg)

- SQL表达式查询

方法三允许自定义完成子句expression的接口，可以根据环境做一些动态扩展

![SQL表达式](./images/sqlbiaodashi2.jpg)

- 数据序列化

Password方法实现SerializerInterface，通过实现接口，可以把加解密逻辑都写进去，保存时就可以自动的把密码进行加密等操作

![SQL表达式](./images/sqlbiaodashi3.jpg)

## 批量数据操作

- 批量创建/查询

下面代码中，第一个接口会把所有的语句通过一条SQL进行创建，第二个接口会把每100条数据作为一批进行创建

查询可以通过rows，rows其实只是一个迭代器，可以一个一个的迭代取出，从而不用把所有的查询结果都加载到内存中

第二个方法使用FindInBatches方法，可以每100条进行查询，然后将结果交给find方法，进行分治的处理，可以减轻数据库侧的负载

![批量操作](./images/piliangcaozuo.jpg)

- 批量更新

批量更新的接口对于不同的数据库也是不一样的，提供了GORM的OnConflict子句进行支持

![批量更新](./images/pilianggengxin.jpg)

- 批量数据加速操作

GORM为了保持数据的一致性，在create、update以及delete的时候会默认开启一下事务

方法二的话，在注册一个create的callback时候，会将before create和after create的一些方法也执行，但是一般导入数据的时候这并不是非常需要的，可以进行跳过从而提升性能

 第三个方法是预编译的方式

![批量加速](./images/piliangshuju.jpg)

## 代码复用、分库分表、Sharding

- 代码复用

Scopes方法提供代码复用的方法

![代码复用](./images/daimafuyong.jpg)

- 分库分表

常见的数据库优化手段，同一个数据可能存在不同的数据库或者数据表中，可能需要根据一些条件来判断从哪个数据库或者数据表中取数据，可以根据一些规则抽象表的逻辑，再通过Scopes方法复用一下，实现分库分表操作

![分库分表](./images/fenkufenbiao.jpg)

- Sharding支持

因为分库分表大部分都有固定的逻辑，因此可以将其抽象一下，通过Sharding配置

![Sharding](./images/sharding.jpg)

## 混沌工程 / 压测

对于一些特别重要的服务，如钱包，支付系统，中间可能有大量的兼容系统，来校验并告警来保证数据一致性，因此需要人为创造一些错误信息来看系统能不能发现这些问题并解决这些问题


