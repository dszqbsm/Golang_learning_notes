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
