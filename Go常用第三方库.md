# Go常用第三方库

整理我在写代码过程中常用的标准库，以及一些常用的方法，并用例子展现使用说明，作为自己的标准库使用手册

## gin框架

gin框架基于httprouter开发，是十分流行的轻量级Web框架

```go
// 在项目根目录下
go mod init gin_demo
go get -u github.com/gin-gonic/gin      // 下载gin框架

package main
import (
    "github.com/gin-gonic/gin"
)
func main() {
    // 创建一个默认的路由引擎
    r := gin.Default()
    // GET：请求方式；/hello：请求的路径
    // 当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
    r.GET("/hello", func(c *gin.Context) {
        // c.JSON：返回JSON格式的数据
        c.JSON(200, gin.H{
            "message": "Hello world",
        })
    })
    // 启动HTTP服务，默认在0.0.0.0:8080启动服务
    r.Run()
}
```

1. 路由

路由分组通常用作划分业务逻辑或者API版本中

- `r.GET("/index", func(c *gin.Context) {...})`: 匹配/index路径的GET请求
- `r.POST("/login", func(c *gin.Context) {...})`: 匹配/login路径的POST请求
- `r.Any("/test", func(c *gin.Context) {...})`: 匹配/test路径的所有请求类型
- `r.NoRoute(func(c *gin.Context) { c.HTML(http.StatusNotFound, "views/404.html", nil )})`: 匹配路径不存在的请求

```go
// {}包裹共同URL前缀路由组
func main() {
    r := gin.Default()
    userGroup := r.Group("/user")       // 用户业务组
    {
        userGroup.GET("/index", func(c *gin.Context) {...})     // /user/index
        userGroup.POST("/login", func(c *gin.Context) {...})    // /user/login
    }
    shopGroup := r.Group("/shop")       // 商城业务组
    {
        shopGroup.GET("/index", func(c *gin.Context) {...})     // /shop/index
        shopGroup.POST("/cart", func(c *gin.Context) {...})     // /shop/cart
        // 嵌套路由组
        xx := shopGroup.Group("/xx")
        {
            xx.GET("/oo", func(c *gin.Context) {...})               // /shop/xx/oo
        }
    }
    r.Run()
}
```

2. 获取参数

- 获取URL中?后面携带的参数，如/user/search?username=zhangsan&address=beijing

```go
func main() {
    r := gin.Default()
    r.GET("/user/search", func(c *gin.Context) {
        username := c.DefaultQuery("username", "zhangsan")      // 默认值为zhangsan
        address := c.Query("address")
        // 输出JSON结果给调用方
        c.JSON(http.StatusOK, gin.H{
            "message":  "ok",
            "username": username,
            "address":  address,
        })
    })
    r.Run()
}
```

- 获取form表单提交数据，如/user/search发送一个POST请求，获取请求数据

```go
func main() {
    r := gin.Default() 
    r.POST("/user/search", func(c *gin.Context) {
        username := c.DefaultPostForm("username", "zhangsan")       // DefaultPostForm取不到值时会返回指定的默认值
        address := c.PostForm("address")
        // 输出JSON结果给调用方
        c.JSON(http.StatusOK, gin.H{
            "message":  "ok",
            "username": username,
            "address":  address,
        })
    })
    r.Run(":8080")
}
```

- 获取json方式提交的数据，如向/json发送了一个POST请求，获取请求数据

```go
r.POST("/json", func(c *gin.Context) {
    b, _ := c.GetRawData()  // 从c.Request.Body读取请求数据
    var m map[string]interface{}
    // 反序列化
    _ = json.Unmarshal(b, &m)
    c.JSON(http.StatusOK, m)
})
```

- 获取通过URL路径传递的参数数据，如/user/search/zhangsan/beijing

```go
func main() {
    r := gin.Default() 
    r.GET("/user/search/:username/:address", func(c *gin.Context) {
        username := c.Param("username")
        address := c.Param("address")
        // 输出JSON结果给调用方
        c.JSON(http.StatusOK, gin.H{
            "message":  "ok",
            "username": username,
            "address":  address,
        })
    })
    r.Run(":8080")
}
```

- .ShouldBing()能够基于请求自动提取JSON、form表单和QueryString类型的数据，并把值绑定到指定的结构体对象

更方便地获取请求相关参数的方法，基于请求的Content-Type识别请求数据类型并利用反射机制自动提取请求中的数据并存入结构体中

ShouldBind会按照下面的顺序解析请求中的数据完成绑定

（1）如果是GET请求，则只会使用Form绑定引擎（query）
（2）如果是POST请求，则会先检查content-tyep是否为JSON或XML，再使用Form（form-data）绑定引擎

```go
type LOgin struct {
    User string `form:"user" json:"user" binding:"required"`
    Password string `form:"password" json:"password" binding:"required"`
}
func main() {
    r := gin.Default()
    // 绑定JSON的示例({"user": "zhangsan", "password": "123456"})
    r.POST("/loginJSON", func(c *gin.Context) {
        var login Login
        if err := c.ShouldBind(&login); err == nil {
            fmt.Printf("login info:%#v\n", login)
            c.JSON(http.StatusOK, gin.H{
                "user":     login.User,
                "password": login.Password,
            })
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
    })
    // 绑定form表单示例(user=zhangsan&password=123456)
    r.POST("/loginForm", func(c *gin.Context) {
        var login Login
        // ShouldBind()会根据请求的Content-Type自动选择绑定器
        if err := c.ShouldBind(&login); err == nil {
            fmt.Printf("login info:%#v\n", login)
            c.JSON(http.StatusOK, gin.H{
                "user":     login.User,
                "password": login.Password,
            })
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
    })
    // 绑定QueryString示例(/loginQuery?user=zhangsan&password=123456)
    r.GET("/loginQuery", func(c *gin.Context) {
        var login Login
        // ShouldBind()会根据请求的Content-Type自动选择绑定器 
        if err := c.ShouldBind(&login); err == nil {
            fmt.Printf("login info:%#v\n", login)
            c.JSON(http.StatusOK, gin.H{
                "user":     login.User,
                "password": login.Password,
            })
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
    })
    r.Run(":8080")
}
```

3. 上传文件

```html
<!-- 上传文件前端代码 -->
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <title>上文文件示例</title>
</head>
<body>
    <form action="/upload" method="post" enctype="multipart/form-data">      // action属性指定了表单提交的目标URL，method属性指定了表单提交的方式，enctype属性指定了表单提交时使用的编码类型
        <input type="file" name="f1">                                        // input标签type属性为file，提供了一个文件选择按钮，用户可以从本地文件系统中选择文件上传，name属性定义了文件在表单数据中的名称，后端用这个名称检索文件
        <input type="submit" value="上传">                                   // input标签type属性为submit，提供了一个提交按钮，用户点击后会将表单数据提交到服务器
    </form>
</body>
</html>
```

```go
func main() {
    r := gin.Default()
    // 处理multipart forms提交文件时默认的内存限制是32MB
    // 可以通过以下的方式修改
    // r.MaxMultipartMemory = 8 << 20 // 8MB
    r.POST("/upload", func(c *gin.Context) {

        // 单个文件
        // 获取上传的文件
        file, err := c.FormFile("f1")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": err.Error(),
            })
            return
        }
        log.Println(file.Filename)
        dst := fmt.Sprintf("C:/tmp/%s", file.Filename)
        // 上传文件到指定的目录
        c.SaveUploadedFile(file, dst)

        // 多个文件
        // 获取上传的文件
        form, _ := c.MultipartForm()        // 这个方法用于解析请求体中的multipart表单数据，该方法返回两个值，第一个值是*http.MultipartForm类型的指针，包含了解析后的表单数据，包括文件和普通字段，第二个是error类型的值
        files := form.File["file"]          // 从MultipartForm结构体中获取名为"file"的文件列表，form.File是一个映射(map[string][]*multipart.FileHeader)他的键是表单字段名，值是*multipart.FileHeader类型的切片，每个*multipart.FileHeader代表一个上传的文件
        for index, file := range files {
            log.Println(file.Filename)
            dst := fmt.Sprintf("C:/tmp/%s_%d", file.Filename, index)
            // 上传文件到指定的目录
            c.SaveUploadedFile(file, dst)
        }

        c.JSON(http.StatusOK, gin.H{
            "message": fmt.Sprintf("'%s' uploaded!", file.Filename),
        })
    })
    r.Run(":8080")
}
```

4. gin渲染响应数据

根据请求中的数据类型来返回对应类型的响应数据

- HTML格式响应数据

虽然浏览器能够渲染HTML页面，但当涉及到动态内容或服务端逻辑时，需要后端渲染HTML模板，例如后端可以根据代码逻辑执行用户认证，数据库查询等操作，然后将结果与HTML模板结合，生成最终的HTML页面响应给用户，也就是说，实际上不是后端渲染HTML页面，本质上还是浏览器渲染HTML页面，后端只是根据用户提交的数据，动态对HTML模板进行修改再返回，这样能够做到动态内容响应

在templates文件夹中存放模板文件

/templates/posts/index.html文件内容如下：

```html
{{define "posts/index.html"}}
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>posts/index</title>
</head>
<body>
    {{.title}}
</body>
</html>
{{end}}
```

/templates/users/index.html文件内容如下：

```html
{{define "users/index.html"}}
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>users/index</title>
</head>
<body>
    {{.title}}
</body>
</html>
{{end}}
```

gin框架中使用LoadHTMLGlob()或者LoadHTMLFiles()方法渲染HTML模板

```go
func main() {
    r := gin.Default()
    // 加载模板文件
    r.LoadHTMLGlob("templates/**/*")                                                // 加载templates目录及其子目录下的所有文件
    // r.LoadHTMLFiles("templates/posts/index.html", "templates/users/index.html")  // 指定具体的模板文件
    r.GET("/posts/index", func(c *gin.Context) {
        // 根据文件名渲染
        c.HTML(http.StatusOK, "posts/index.html", gin.H{        // c.HTML用来渲染模板，第一个参数是http响应状态码，第二个参数表示模板文件的名称，必须与加载的匹配，第三个参数是一个映射，包含要传递给模板的变量
            "title": "posts/index",                             // 模板接收title变量，在HTML模板中，使用该变量的值替换模板中{{.title}}的位置，从而动态生成内容    
        })
    })
    r.GET("/users/index", func(c *gin.Context) {
        // 根据文件名渲染
        c.HTML(http.StatusOK, "users/index.html", gin.H{
            "title": "users/index",
        })
    })
    r.Run(":8080")
}
```

- 静态文件处理

当渲染的HTML文件引用了静态文件时，需要按照以下方法在渲染页面前调用gin.Static方法

```go
func main() {
    r := gin.Default()
    r.Static("/static", "./static")     // 将指定的文件系统目录作为静态文件服务的根目录，这里代表所有/static/开头的URL都会映射到./static目录下的文件
    r.LoadHTMLGlob("templates/**/*")
    ...
    r.Run(":8080")
}
```

- JSON格式响应数据

```go
func main() {
    r := gin.Default() 
    // gin.H是map[string]interface{}的缩写
    r.GET("/soneJSON", func(c *gin.Context) {
        // 方式一：自己拼接JSON
        c.JSON(http.StatusOK, gin.H{
            "message": "hello world!",
        })

        // 方式二：使用结构体
        var msg struct {
            Name    string `json:"user"`
            Message string
            Number  int
        }
        msg.Name = "zhangsan"
        msg.Message = "hello world!"
        msg.Number = 123
        c.JSON(http.StatusOK, msg)
    })
    r.Run(":8080")
}
```

- protobuf格式响应数据

Protocol Buffers是一种轻量高效的结构化数据存储格式，可以用于结构化数据的序列化，反序列化和传输

Protocol Buffers支持多种编程语言，使得使用不同语言编写的系统之间可以轻松地进行数据交换，并且效率比XML和JSON高，因为他不需要解析复杂的文本结构，而是直接操作二进制数据

关于Protocol Buffers这种数据传输格式，还不够熟悉

```go
r.GET("/someProtoBuf", func(c *gin.Context) {
    resp := []int64{int64(1), int64(2)}
    label := "test"
    data := &protoexample.Test{
        Label: &label,
        Resp:  resp,
    }
    // 注意：数据在想用中变为二进制
    // 输出被protoexample.Test protobuf 序列化了的数据
    // data结构体将被序列化为二进制格式的protobuf数据，然后写入到响应中
    c.ProtoBuf(http.StatusOK, data)
})
```

5. 重定向

支持将请求重定向到内部网址或外部网址

```go
r.GET("/test", func(c *gin.Context) {
    // c.Redirect用于发送一个http重定向响应到客户端
    c.Redirect(http.StatusMovedPermanently, "https://www.liwenzhou.com/")       // http.StatusMovedPermanently是重定向状态码301，第二个参数是重定向的目标URL
})
```

6. 路由重定向

可以将当前请求交给其他路由处理函数处理

```go
r.GET("/test", func(c *gin.Context) {
    // 指定重定向的URL
    c.Request.URL.Path = "/test2"
    r.HandlContext(c)               // 将/test请求交给/test2的路由处理函数
})
r.GET("/test2", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"hello world!"})
})
```

7. 中间件

在处理请求的过程中加入自己的钩子（Hook）函数，这个钩子函数就叫中间件，适合处理一些公共的业务逻辑，例如登录认证、权限校验、数据分页、记录日志、耗时统计等

`type HandleFunc func(*Context)`：中间件必须是gin.HandlerFunc类型，接收一个*gin.Context参数

gin框架中常见的中间件示例：

- 记录接口耗时的中间件

```go
func StatCost() gin.HandlerFunc {
    return func(c *gin.Context) {       // c是*gin.Context类型的指针，代表gin的上下文，包含了请求和响应的信息
        start := time.Now()
        c.Set("name", "zhangsan") // 在上下文中设置值，后续的处理函数能够取到该值

        c.Next() // 调用该请求的剩余处理程序，告诉gin框架继续执行处理链中的下一个处理器

        // 计算耗时
        cost := time.Since(start)
        log.Println(cost)
    }
}
```

- 记录响应体的中间件

```go
type bodyLogWriter struct {
    gin.ResponseWriter      // 嵌入gin框架的ResponseWriter
    body *bytes.Buffer      // 用于记录的response
}
// Write 写入响应体数据
func (w bodyLogWriter) Write(b []byte) (int, error) {
    w.body.Write(b)                     // 将数据写入到缓存中
    return w.ResponseWriter.Write(b)    // 调用原始的ResponseWriter的Write方法，将数据写入响应体
}
// ginBodyLogMiddleware 一个记录返回给客户端响应体的中间件
func ginBodyLogMiddleware(c *gin.Context) {
    blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
    c.Writer = blw      // 使用自定义的类型替换默认的Writer，从而实现记录响应体
    c.Next()
    fmt.Println("Response body: " + blw.body.String())
}
```

- 跨域中间件

使用第三方库gin-contrib/cors库，该库支持各种常用的配置项，以下中间件需要注册在业务处理函数前面

```go
package main
import (
    "time"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)
func main() {
    r := gin.Default()
    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{"URL_ADDRESS"},      // 允许跨域发来请求的网站
        AllowMethods: []string{"PUT", "PATCH"},     // 允许使用的请求方法
        AllowHeaders: []string{"Origin"},           // 允许使用的请求头
        ExposeHeaders: []string{"Content-Length"},  // 允许获取的响应头
        AllowCredentials: true,                     // 是否允许cookie
        AllowOriginFunc: func(origin string) bool { // 自定义过滤源站的方法
            return origin == "URL_ADDRESS"
        },
        MaxAge: 12 * time.Hour,
    }))
    r.Run(":8080")
}

// 使用默认配置，允许所有的跨域请求
func main() {
    r := gin.Default()
    // 写法1
    r.User(cors.Default())
    // 写法2
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    r.Use(cors.New(config))
    r.Run(":8080")
}
```

- 为路由添加中间件

为全局路由注册中间件

```go
func main() {
    // 新建一个没有任何默认中间件的路由
    r := gin.New()
    // 注册一个全局中间件
    r.Use(StatCost())
    r.GET("/test", func(c *gin.Context) {
        name := c.MustGet("name").(string)      // 从上下文取值
        log.Println(name)
        c.JSON(http.StatusOk, gin.H{
            "message": "hello world",
        })
    })
    r.Run(":8080")
}
```

为某个路由单独注册中间件

```go
r.GET("/test2", StatCost(), func(c *gin.Context) {
    name := c.MustGet("name").(string)
    log.Println(name)
    c.JSON(http.StatusOk, gin.H{
        "message": "hello world",
    })
})
```

为路由组注册中间件

```go
// 写法1
shopGroup := r.Group("/sho", StatCost())
{
    shoGroup.GET("/index", func(c *gin.Context) {...})
}

// 写法2
shopGroup := r.Group("/shop")
shopGroup.Use(StatCost())
{
    shopGroup.GET("/index", func(c *gin.Context) {...})
}
```

`gin.Default()`默认使用了Logger和Recovery中间件，并且即使配置了GIN_MODE=release，Logger中间件也会将日志写入gin.DefaultWriter；Recovery中间件会recover任何panic，如果有panic，则返回500错误码

可以使用gin.New()创建一个没有任何默认中间件的路由

当在中间件中启动新的goroutine时，为了并发安全，在创建的goroutine中不能使用原始的上下文(c *gin.Context)，必须使用其只读副本(c.Copy())

8. 运行多个服务

```go
package main
import (
    "log"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "golang.org/x/sync/errgroup"
)
var (
    g errgroup.Group        // 用于管理并发执行的组
)
func route01() http.Handler {
    e := gin.New()
    e.Use(gin.Recovery())
    e.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "code": http.StatusOK,
            "error": "Welcome server 01",
        })
    })
    return e
}
func route02() http.Handler {
    e := gin.New()
    e.Use(gin.Recovery())
    e.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "code": http.StatusOK,
            "error": "Welcome server 02",
        })
    })
    return e
}
func main() {
    server01 := &http.Server{
        Addr: ":8080",
        Handler: route01(),
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    server02 := &http.Server{
        Addr: ":8081",
        Handler: route02(),
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    // 借助errgroup.Group或者自行开启两个goroutine分别启动两个服务
    // Go方法会并发执行传入的函数，并返回一个错误
    g.Go(func() error {
        return server01.ListenAndServe()
    })
    g.Go(func() error {
        return server02.ListenAndServe()
    })
    if err := g.Wait(); err != nil {        // 会等待所有由g.Go启动的goroutine完成，如果有任何一个goroutine返回错误，Wait会立即返回该错误信息，并使用log.Fatal记录错误并停止程序
        log.Fatal(err)
    }
}
```














