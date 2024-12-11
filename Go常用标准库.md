# Go常用标准库

整理我在写代码过程中常用的标准库，以及一些常用的方法，并用例子展现使用说明，作为自己的标准库使用手册

## fmt包

包括向外输出内容和获取输入内容两大部分

### 向外输出内容

1. Print系列函数

会将内容输出到系统的标准输出

- `func Print(a ...interface{}) (n int, err error)`: 直接输出内容

- `func Printf(format string, a ...interface{}) (n int, err error)`: 支持格式化输出字符串

- `func Println(a ...interface{}) (n int, err error)`: 在输出内容的结尾添加一个换行符

> `a ...interface{}`是变长参数的表示方式，这里表示接收零个或多个空接口类型

2. Fprint系列函数

会将内容输出到一个io.Writer接口类型的变量w中，不仅仅是标准输出和文件，只要满足io.Writer接口的类型（例如网络I/O等）都支持写入

- `func Fprint(w io.Writer, a ...interface{}) (n int, err error)`：直接输出内容

- `func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)`：支持格式化输出字符串

- `func Fprintln(w io.Writer, a ...interface{}) (n int, err error)`：在输出内容的结尾添加一个换行符

```go
// 区别于Print系列函数
// 也可以输出到标准输出（os.Stdout）
func fprintlnDemo() {
    str := "12345"
    fmt.Fprinitln(os.Stdout, str)
}

// fprintfDemo 将格式化后的字符串写入xx.txt文件
func fprintfDemo() {
    name := "zhangsan"
    fileObj, _ := os.OpenFile("./xx.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    fmt.Fprintf(fileObj, "name: %s", name) 
}
// 编译执行后，会在当前目录下创建一个名为xx.txt的文件
```

> `os.O_CREATE|os.O_WRONLY|os.O_APPEND`：这是一个标志位的组合，用于指定文件打开的方式，os.O_CREATE表示如果文件不存在，则创建文件；os.O_WRONLY表示以只写方式打开文件；os.O_APPEND表示如果文件存在，则写入的数据会被追加到文件末尾；0644这是一个八进制数，用于指定文件的权限；6表示文件所有者有读写权限；4表示组用户有读权限；4表示其他用户有读权限

3. Sprint系列函数

Sprint系列函数会把传入的数据生成并返回一个字符串

- `func Sprint(a ...interface{}) string`：直接返回内容

- `func Sprintf(format string, a ...interface{}) string`：支持格式化输出字符串

- `func Sprintln(a ...interface{}) string`：在输出内容的结尾添加一个换行符

```go
func sprintDemo() {
    name := "zhangsan"
    age := 18
    s := fmt.Sprintf("name: %s, age: %d", name, age)
    fmt.Println(s)      // name: zhangsan, age: 18
}
```

4. Errorf函数

Errorf函数根据format参数生成格式化字符串并返回一个包含该字符串的error，通常使用这种方式来自定义error，如`err := fmt.Errorf("无效的id")`

`func Errorf(format string, a ...interface{}) error`

还可以使用格式化动词来生成一个可以包含指定error的新error

```go
func errorDemo() {
    e := errors.New("连接失败")                 // 原始错误
    err := fmt.Errorf("查询失败， err: %w", e)  // 生成一个包含原始error的新error
}
```

### 格式化占位符

`fmt.*printf`系列函数都支持format格式化参数

1. 通用占位符

| 占位符 | 说明 |
| :---: | :---: |
| %v | 值的默认格式表示 |
| %+v | 类似%v，但输出结构体时会添加字段名 |
| %#v | 值的Go语法表示 |
| %T | 值的类型 |
| %% | 百分号 |

```go
fmt.Printf("%v\n", 100)     // 100
fmt.Printf("%v\n", false)   // false
o := struct{
    name string
}{"zhangsan"}
fmt.Printf("%v\n", o)       // {zhangsan}
fmt.Printf("%#v\n", o)      // struct { name string }{ name:"zhangsan" }
fmt.Printf("%T\n", o)       // struct
fmt.Printf("100%%\n")       // 100%
```

2. 布尔型

布尔型占位符主要指`%t`，表示true或false

3. 整型

| 占位符 | 说明 |
| :---: | :---: |
| %b | 二进制表示 |
| %c | 该值对应的unicode码值 |
| %d | 10进制表示 |
| %o | 8进制表示 |
| %x | 16进制表示，使用a-f |
| %X | 16进制表示，使用A-F |
| %U | unicode格式表示 |
| %q | 该值对应单引号括起来的Go语法字符字面值，必要时会采用安全的转移表示 |

```go
n := 65
fmt.Printf("%b\n", n)   // 1000001
fmt.Printf("%c\n", n)   // A
fmt.Printf("%d\n", n)   // 65
fmt.Printf("%o\n", n)   // 101
fmt.Printf("%x\n", n)   // 41
fmt.Printf("%X\n", n)   // 41
```

4. 浮点数与复数

| 占位符 | 说明 |
| :---: | :---: |
| %b | 无小数部分，二进制指数的科学计数法，如-123456p-78 |
| %e | 科学计数法，如-1234.456e+78 |
| %E | 科学计数法，如-1234.456E+78 |
| %f | 有小数部分但无指数部分，如123.456 |
| %F | 等价于%f |
| %g | 根据情况选择%e或%f格式（以获得更简洁、准确的输出） |
| %G | 根据情况选择%E或%F格式（以获得更简洁、准确的输出） |

```go
f := 12.34
fmt.Printf("%b\n", f)   // 694680242521899p-49
fmt.Printf("%e\n", f)   // 1.234000e+01
fmt.Printf("%E\n", f)   // 1.234000E+01
fmt.Printf("%f\n", f)   // 12.340000
fmt.Printf("%F\n", f)   // 12.340000
fmt.Printf("%g\n", f)   // 12.34
fmt.Printf("%G\n", f)   // 12.34
```

5. 字符串和[]byte

| 占位符 | 说明 |
| :---: | :--- :|
| %s | 直接输出字符串或者[]byte |
| %q | 该值对应双引号括起来的Go语法字符串字面值，必要时会采用安全的转义表示 |
| %x | 每个字节用两字符十六进制数表示（使用a-f） |
| %X | 每个字节用两字符十六进制数表示（使用A-F） |

```go
s := "zhangsan"
fmt.Printf("%s\n", s)   // zhangsan   
fmt.Printf("%q\n", s)   // "zhangsan"
fmt.Printf("%x\n", s)   // 7a68616e6773616e
fmt.Printf("%X\n", s)   // E68616E6773616E
```

6. 指针

指针占位符主要指`%p`，代表表示为十六进制，并加上前导的0x

```go
a := 10
fmt.Printf("%p\n", &a)      // 0xc0000140b8
fmt.Printf("%#p\n", &a)     // c0000140b8
```

7. 宽度标识符

通过一个紧跟在百分号后面的十进制数指定，若未指定则在表示值时非必要不填充

精度通过点号后面的十进制数指定，若未指定则使用默认精度

| 占位符 | 说明 |
| :---: | :---: |
| %f | 默认宽度，默认精度 |
| %9f | 宽度9，默认精度 |
| %.2f | 默认宽度，精度2 |
| %9.2f | 宽度9，精度2 |
| %9.f | 宽度9，精度0 |

```go
n := 12.34
fmt.Printf("%f\n", n)       // 12.340000
fmt.Printf("%9f\n", n)      // 12.340000
fmt.Printf("%.2f\n", n)     // 12.34
fmt.Printf("%9.2f\n", n)    //     12.34
fmt.Printf("%9.f\n", n)     //        12
```

8. 其他flag

| 占位符 | 说明 |
| :---: | :---: |
| '+' | 总是输出数值的正负号，对%q（%+q）会生成全部是ASCII字符的输出（通过转义） |
| '-' | 在输出右边而不是左边填充空白（即从默认的右对齐切换为左对齐） |
| '#' | 八进制数前加0（%#o），十六进制数前加0x（%#x）或0X（%#X），指针类型去掉前面的0x（%#p），对%q（%#q）、%U（%#U）会输出空格和单引号括起来的Go字面值 |
| '' | 对于数值，在正数前面加空格、在负数前面加负号；对于字符串，在采用%x或%X时会在输出的个字节之间加空格 |
| '0' | 使用0而不是空格填充，对于数值类型，会把填充的0放在正负号后面 |

```go
s := "Go"
fmt.Printf("%s\n", s)           // Go
fmt.Printf("%5s\n", s)          //      Go
fmt.Printf("%-5s\n", s)         // Go     
fmt.Printf("%5.7s\n", s)        //      Go
fmt.Printf("%-5.7s\n", s)       // Go
fmt.Printf("%5.2s\n", s)        // Go
fmt.Printf("%05s\n", s)         //      Go
// 数字类型
i := -10
fmt.Printf("%d\n", i)           // -10
fmt.Printf(" %d\n", i)          //  -10
fmt.Printf("%5d\n", i)          //      -10
fmt.Printf("%-5d\n", i)         // -10
fmt.Printf("%05d\n", i)         // -0010
f := 12.34
fmt.Printf("%f\n", f)           // 12.340000
fmt.Printf("% f\n", f)          //  12.340000
fmt.Printf("%f\n", f)           // 12.340000
fmt.Printf("%-f\n", f)          // 12.340000
```

### 获取输入内容

以下三个函数可以在程序运行过程中从标准输入获取用户的输入

1. fmt.Scan函数

`func Scan(a ...interface{}) (n int, err error)`: 从标准输入扫描文本，读取由空白符分隔的值保存到传递的参数中，换行符视为空白符，返回成功扫描的数据个数和遇到的任何错误

```go
func scanDemo() {
    var (
        job string
        num int
        skip bool
    )
    fmt.Scan(&job, &num, &skip)     // 获取输入，输入的内容按空格分隔
    fmt.Printf("获取的输入内容 job: %s, num: %d, skip: %t\n", job, num, skip)
}
```

2. fmt.Scanf函数

`func Scanf(format string, a ...interface{}) (n int, err error)`：从标准输入扫描文本，根据format参数指定的格式去读取由空白符分隔的值保存到传递的参数中，换行符视为空白符，返回成功扫描的数据个数和遇到的任何错误

只有按照format格式输入的数据才会被扫描并赋值给对应变量，否则所有变量均为默认值

```go
func scanfDemo() {
    var (
        job string
        num int
        skip bool
    )
    fmt.Scanf("1:%s 2:%d 3:%t", &job, &num, &skip)
    fmt.Printf("获取的输入内容 job: %s, num: %d, skip: %t\n", job, num, skip)
}
```

3. fmt.Scanln函数

`func Scanln(a ...interface{}) (n int, err error)`：从标准输入扫描文本，在遇到换行时才停止扫描，最后一个数据后面必须有换行或者达到结束位置，返回成功扫描的数据个数和遇到的任何错误

```go
func scanlnDemo() {
    var (
        job string
        num int
        skip bool
    )
    fmt.Scanln(&job, &num, &skip)
    fmt.Printf("获取的输入内容 job: %s, num: %d, skip: %t\n", job, num, skip)
}
```

4. bufio包

当输入的内容可能包含空格时，要想完整获取输入的内容可以使用bufio包

```go
func bufioDemo() {
    reader := bufio.NewReader(os.Stdin)     // 从标准输入生成读对象
    fmt.Print("请输入内容：")
    text, _ := reader.ReadString('\n')      // 读取内容，直到遇到换行符
    text = strings.TrimSpace(text)          // 去掉首尾多余空格
    fmt.Printf("获取的输入内容：%#v\n", text)
}
```

5. Fscan系列函数

以下几个函数功能类似于`fmt.Scan`、`fmt.Scanf`、`fmt.Scanln`，只不过它们不是从标准输入中读取数据，而是从io.Reader中读取数据

- `func Fscan(r io.Reader, a ...interface{}) (n int, err error)`
- `func Fscanf(r io.Reader, format string, a ...interface{}) (n int, err error)`
- `func Fscanln(r io.Reader, a ...interface{}) (n int, err error)`

6. Sscan系列函数

以下几个函数功能类似于`fmt.Scan`、`fmt.Scanf`、`fmt.Scanln`，只不过它们不是从标准输入中读取数据，而是从指定字符串中读取数据

- `func Sscan(str string, a ...interface{}) (n int, err error)`
- `func Sscanf(str string, format string, a ...interface{}) (n int, err error)`
- `func Sscanln(str string, a ...interface{}) (n int, err error)`

### 获取命令行参数

通过`os.Args`可以获取到命令行参数，他是一个存储命令行参数的字符串切片，第一个元素是执行文件的名称

```go
func main() {
    // os.Args是一个[]string
    if len(os.Args) > 0 {
        for index, arg := range os.Args {
            fmt.Printf("args[%d]=%v\n", index, arg)
        }
    }
}
/*
> ./args_demo a b c d
args[0]=./args_demo
args[1]=a
args[2]=b
args[3]=c
args[4]=d
*/
```

## flag包

flag包用于实现命令行参数的解析


1. 定义命令行参数

对于需要在执行程序时通过命令行指定运行时所需的参数，这种场景，有以下两种常用的定义命令行flag参数的方法，并且在定义号命令行flag参数后，需要调用`flag.Parse`对命令行参数进行解析

- `flag.Type(flag名, 默认值, 帮助信息) *Type`

```go
// 此时job、num、skip、delay均为指针类型
job := flag.String("job", "work", "任务名称")
num := flag.Int("num", 1, "次数")
skip := flag.Bool("skip", false, "是否跳过失败任务")
delay := flag.Duration("d", 0, "任务间隔时间")
```

- `flag.TypeVar(Type指针, flag名, 默认值, 帮助信息)`

```go
var (
    job string
    num int
    skip bool
    delay time.Duration
)
flag.StringVar(&job, "job", "work", "任务名称")
flag.IntVar(&num, "num", 1, "次数")
flag.BoolVar(&skip, "skip", false, "是否跳过失败任务")
flag.DurationVar(&delay, "d", 0, "任务间隔时间")
```

2. 命令行参数格式

必须使用等号的方式指定布尔类型的参数，并且flag包解析参数时会在第一个非flag参数（单个"-"不是flag参数）之前停止，或者在终止符"-"之后停止

- `-flag xxx`：使用空格，一个-符号
- `--flag xxx`：使用空格，两个-符号
- `-flag=xxx`：使用等号，一个-符号
- `--flag=xxx`：使用等号，两个-符号

```go
// 一些其他函数用来获取命令行参数的其他信息
flag.Args()     // 以[]string类型返回命令行参数后的其他参数
flag.NArg()     // 返回命令行参数后的其他参数数量
flag.NFlag()    // 返回使用的命令行参数数量
```

3. flag包完整示例

```go
func main() {
    var (
        job string
        num int
        skip bool
        delay time.Duration
    )
    flag.StringVar(&job, "job", "work", "任务名称")
    flag.IntVar(&num, "num", 1, "次数")
    flag.BoolVar(&skip, "skip", false, "是否跳过失败任务")
    flag.DurationVar(&delay, "d", 0, "任务间隔时间")
    // 解析命令行参数
    flag.Parse()
    fmt.Println(job, num, skip, delay)
    // 返回命令行参数后的其他参数
    fmt.Println(flag.Args())
    // 返回命令行参数后的其他参数数量
    fmt.Println(flag.NArg())
    // 返回使用的命令行参数数量
    fmt.Println(flag.NFlag())
}

/*
> ./flag_demo -help
Usage of ./flag_demo:
  -d duration
        任务间隔时间 (default 0s)
  -job string
        任务名称 (default "work")
  -num int
        次数 (default 1)
  -skip
        是否跳过失败任务

// 在执行程序时添加命令行参数
> ./flag_demo -job search --num 10 -skip=true -d=1m
search 10 true 1m0s 
[] 
0 
3

// 在执行程序时使用其他命令行参数
> ./flag_demo test 10
work 10 false 0s            // 命令行没有给定参数，使用默认值
[test 10]
2
0
*/
```

也即是，可以通过flag包来设置命令行的参数设置，从而通过-help能够查看帮助信息，并且可以通过-job、-num、-skip、-d来设置对应的参数值，这些值会赋值到flag包设置的变量中，从而在程序中直接使用，对于未设置的参数，会使用默认值；此外命令行输入的参数都可以通过flag.Args()获取到，因此可以获取到其他命令行参数

## time包

1. 时间类型

```go
func timeDemo() {
    now := time.Now()       // 获取当前时间
    fmt.Printf("current time:%v\n", now)

    year := now.Year()     // 年
    month := now.Month()   // 月
    day := now.Day()       // 日
    hour := now.Hour()     // 小时
    minute := now.Minute() // 分钟
    second := now.Second() // 秒
}
```

2. location和time zone

```go
func timezoneDemo() {
    // 中国没有夏令时，使用固定的8小时的UTC时差
    secondEastOfUTC := int((8 * time.Hour).Seconds())           // 计算了UTC时间与背景时间相差的秒数，8小时转换为秒
    // FixedZone返回始终使用给定区域名称和偏移量的Location
    beijing := time.FixedZone("Beijing Time", secondEastOfUTC)  // 创建一个名为"Beijing Time"的时区，偏移量为secondEastOfUTC，即8小时
    // 根据location获取当前时间
    now := time.Now().In(beijing)                               // 根据时区获取当前时间
    fmt.Println(now)

    // 如果当前系统有时区数据库，则可以使用LoadLocation加载一个位置得到对应的时区
    newYork, err := time.LoadLocation("America/New_York")   // UTC-05:00

    // 创建时间对象需要指定位置，常用的位置是time.Local（当地时间）和time.UTC（UTC时间），以下都表示2009年1月1日23:00:00
    timeInLocal := time.Date(2009, 1, 1, 23, 0, 0, 0, time.Local)       // 系统本地时间
    timeInUTC := time.Date(2009, 1, 1, 23, 0, 0, 0, time.UTC)           // UTC时间
    timeInNewYork := time.Date(2009, 1, 1, 23, 0, 0, 0, newYork)        // 纽约时间
    timeInBeijing := time.Date(2009, 1, 1, 23, 0, 0, 0, beijing)        // 北京时间

    timesAreEqual = timeInUTC.Equal(timeInNewYork)                      // 虽然看似相差8小时，但表示的是同一个时间
}
```

3. Unix Time

表示自1970年1月1日00:00:00UTC至当前时间经过的总秒数

```go
func timestampDemo() {
    now := time.Now()               // 获取当前时间
    timestamp := now.Unix()         // 秒级时间戳
    milli := now.UnixMilli()        // 毫秒级时间戳
    micro := now.UnixMicro()        // 微秒级时间戳
    nanosecond := now.UnixNano()    // 纳秒级时间戳
}
```

```go
func timestamp2Time() {
    // 创建一个名为"Beijing Time"的时区，偏移量为secondEastOfUTC，即8小时
    secondEastOfUTC := int((8 * time.Hour).Seconds())
    beijing := time.FixedZone("Beijing Time", secondEastOfUTC)
    // 根据指定时区获取时间对象，以下时间对象表示时间2009年1月1日23:00:00
    t := time.Date(2009, 1, 1, 23, 0, 0, 0, beijing)
    var (
        sec = t.Unix()          // 秒级时间戳，为int64类型
        msec = t.UnixMilli()    // 毫秒级时间戳，为int64类型
        usec = t.UnixMicro()    // 微秒级时间戳，为int64类型
    )
    // 将以上int64类型的时间戳转换成时间对象
    timeObj := time.Unix(sec, 22)       // 2009-01-01 23:00:00.000000022 +0800 CST
    timeObj = time.UnixMilli(msec)     // 2009-01-01 23:00:00 +0800 CST
    timeObj = time.UnixMicro(usec)     // 2009-01-01 23:00:00 +0800 CST
}
```

4. 时间间隔

`time.Duration`表示两个时间点之间的间隔，以纳秒为单位，最长时间间隔约为290年

```go
// time包中定义的时间间隔类型的常量
const (
    Nanosecond Duration = 1
    Microsecond = 1000 * Nanosecond
    Millisecond = 1000 * Microsecond
    Second = 1000 * Millisecond
    Minute = 60 * Second
    Hour = 60 * Minute
)
// time.Duration表示1纳秒，time.Second表示1秒
```

5. 时间操作

- `func (t Time) Add(d Duration) Time`：返回t+d后的时间
- `func (t Time) Sub(u Time) Duration`：返回t-u的时间间隔，若超过Duration所能表示的最大值/最小值，则返回该最大值/最小值
- `func (t Time) Equal(u Time) bool`：判断两个时间是否相同，会考虑时区的影响，因此不同时区标准的时间也可以正确比较
- `func (t Time) Before(u Time) bool`：判断时间t是否在u之前，会考虑时区的影响，因此不同时区标准的时间也可以正确比较
- `func (t Time) After(u Time) bool`：判断时间t是否在u之后，会考虑时区的影响，因此不同时区标准的时间也可以正确比较

6. 定时器

```go
// 定时器的本质是通道
func tickDemo() {
    ticker := time.Tick(time.Second)        // 定义一个间隔为1s的定时器
    for i := range ticker {
        fmt.Println(i)      // 每秒都会执行的任务
    }
}
```

7. 时间格式化

`time.Format`函数能够将一个时间对象格式化为指定布局的文本表示形式

Go1.20添加了一些常用的格式化布局常量

```go
DateTime = "2006-01-02 15:04:05"
DateOnly = "2006-01-02"
TimeOnly = "15:04:05"

// 时间格式化
func formatDemo() {
    now := time.Now()
    // 使用time包提供的DateTime常量进行格式化
    fmt.Println(now.Format(time.DateTime))
    // 24小时制，小数点后写0，因为有3个0，因此格式化输出的结果也保留3位小数
    fmt.Println(now.Format("2006-01-02 15:04:05.000 Mon Jan"))
    // 12小时制，小数点后写9，省略末尾可能出现的0
    fmt.Println(now.Format("2006-01-02 03:04:05.999 PM Mon Jan"))
    // 只格式化时、分、秒部分
    fmt.Println(now.Format("15:04:05"))
    fmt.Println(now.Format(time.TimeOnly))      // 使用time包的TimeOnly常量
    // 只格式化日期部分
    fmt.Println(now.Format("2006-01-02"))
    fmt.Println(now.Format(time.DateOnly))      // 使用time包的DateOnly常量
}
```

8. 解析字符串格式的时间

使用`time.Parse`和`time.ParseInLocation`两个函数，用于从文本的时间表示中解析出时间对象，第一个参数是用于参考的时间格式，第二个参数是待解析的时间字符串，前者默认使用UTC时间来解析时间字符串，不考虑本地时区的影响，而后者允许指定一个时区，这样解析出来的时间对象会根据提供的时区来设置时间

```go
// 使用parse函数
func parseDemo() {
    // 在没有时区指示符的情况下，time.Parse返回UTC时间
    timeObj, err := time.Parse("2006/01/02 15:04:05", "2016/08/09 21:21:21")
    fmt.Println(timeObj)        // 2016-08-09 21:21:21 +0000 UTC
    // 在有时区指示符的情况下，time.Parse返回对应时区的时间表示
    timeObj, err = time.Parse(time.RFC3339, "2016-08-09T21:21:21+08:00")
    fmt.Println(timeObj)        // 2016-08-09 21:21:21 +0800 CST
}
```

```go
// 使用parseInLocation函数
func parseDemo() {
    now := time.Now()
    // 加载时区
    loc, err := time.LoadLocation("Asia/Shanghai")
    // 按照指定时区和指定格式解析字符串时间
    timeObj, err := time.ParseInLocation("2006/01/02 15:04:05", "2016/08/09 21:21:21", loc)
}
```

## log包

### 默认logger

log包定义了一个logger类型，提供了一些格式化输出日志信息的方法；log包提供了一个预定义的“标准”logger实例，可以直接使用该实例调用Print系列（Print、Printf、Println）函数、Fatal系列（Fatal、Fatalf、Fatalln）函数和Panic系列（Panic、Panicf、Panicln）函数

```go
import (
    "log"
)
func main() {
    log.Println("这是一条普通的日志")
    v := "很普通的"
    log.Printf("这是一条%s日志\n", v)
    log.Fatalln("这是一条会触发fatal的日志")        // 会在输出日志信息后调用os.Exit(1)
    log.Panicln("这是一条会触发panic的日志")        // 会在输出日志信息后panic
}
```

### 自定义logger

默认的logger只会提供日志的日期和时间信息，可以通过`func Flags() int`和`func SetFlags(flag int)`来实现自定义logger，前者返回标准logger的输出配置，后者用来设置标准logger的输出

1. flag选项

```go
// flag选项是一系列定义好的常量，用来控制输出日志信息的细节，但是日志信息输出的顺序是固定的，无法修改
Ldate = 1 << iota       // 日期：2009/01/23
Ltime                   // 时间：01:23:23
Lmicroseconds           // 微妙级别的时间：01:23:23.123123
Llongfile               // 文件全路径名+行号：/a/b/c/d.go:23
Lshortfile              // 文件名+行号：d.go:23（会覆盖Llongfile）
LUTC                    // 使用UTC时间
LstdFlags = Ldate | Ltime // 标准logger的初始值

// 应用示例
log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
log.Println("这是一条很普通的日志")
// 终端输出> 2019/06/19 14:05:17.494943 .../log_demo/main.go:10: 这是一条很普通的日志
```

2. 配置日志前缀

日志前缀支持根据需要为日志设置不同的前缀，方面阅读，`func Prefix() string`和`func SetPrefix(prefix string)`，前者用于查看标准logger的输出前缀，后者用来设置输出前缀

```go
func main() {
    // 自定义日志的配置
    log.SetFlags(log.Llongfile | log.Lmicroseconds | Log.Ldate)
    log.Println("这是一条很普通的日志")
    // 通过SetPrefix自定义前缀
    log.SetPrefix("[order]")
    log.Println("这是一条很普通的日志")
}

// 终端输出> [order] 2019/06/19 14:05:17.494943 .../log_demo/main.go:10: 这是一条很普通的日志
```

3. 配置日志输出位置

`func SetOutput(w io.Writer)`：用于设置日志输出位置，可以将以下代码写入到init函数中，这样就能保证log在程序启动阶段完成配置

```go
import (
    "fmt"
    "log"
    "os"
)

func init() {
    logFile, err := os.OpenFile("./app.log", os.O_CREATE|OS.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        fmt.Println("open log file failed, err:", err)
        return
    }
    log.SetOutput(logFile)      // 设置日志输出位置
    log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}
```

### 创建logger

`func New(out io.Writer, prefix string, flag int) *Logger`：用于创建一个新的logger，其中out参数用于设置日志输出位置，prefix参数用于设置日志前缀，flag参数用于设置日志输出选项

```go
func main() {
    logger := log.New(os.Stdout, "<New>", log.Lshortfile|log.LstdFlags|log.Ltime)
    logger.Println("这是自定义的logger记录的日志")
}
// 终端输出> <New>2019/06/19 14:06:08 main.go:16: 这是自定义的logger记录的日志
```

Go中的log库功能十分有限，一般在实际项目中根据要求编写日志库或者使用第三方日志库，如logrus、zap等

## stronv包

实现了基础数据类型和其字符串表示的相互转换，主要有以下函数：Atoi、Itoa、Parse系列和Format系列

### string与int类型转换

1. Atoi函数

`func Atoi(s string) (int, error)`：Atoi是ParseInt的简写，将字符串类型的证书转换为int类型，若传入的字符串参数无法转换为int类型，则返回错误

```go
func AtoiDemo() {
    str := "45"
    v, err := stronv.Atoi(str)      // 56,nil
    str := "中"
    v, err = stronv.Atoi(str)       // 0,strconv.Atoi: parsing "中": invalid syntax
}
```

2. Itoa函数

- `func Itoa(i int) string`：Itoa是FormatInt的简写，将int类型数据转换为对应的字符串表示

```go
i := 65
v := string(i)      // 使用string(i)返回的是UTF-8编码的表示
v = stronv.Itoa(i)  // "65"
```

### Parse系列函数

将字符串转换为给定类型的值

- `func ParseBool(str string) (value bool, err error)`：返回字符串表示的bool值，只接受1、0、t、f、T、F、true、false、True、False、TRUE、FALSE；传入其他参数会返回错误

- `func ParseInt(s string, base int, bitSize int) (i int64, err error)`：返回字符串表示的整数值，接受正负号，base指定进制（2到36），如果base为0，则会从字符串前置判断，"0x"是16进制，"0"是8进制，否则是10进制；bitSize指定结果必须能赋值（不会溢出）的整数类型，0、8、16、32、64分别代表int、int8、int16、int32、int64；返回的err是*NumError类型的，如果语法有误，err.Error = ErrSyntax；如果结果超出类型范围err.Error = ErrRange

- `func ParseUint(s string, base int, bitSize int) (uint64, error)`：与ParseInt类似，但不接受带政府好的字符串参数，返回字符串表示的无符号整数值

- `func ParseFloat(s string, bitSize int) (float64, error)`：返回字符串表示的浮点数值，bitSize指定了结果必须能赋值（无溢出）的浮点数类型，32是float32（返回值可以不改变精确值的赋值给float32）或64表示float64；返回的err是*NumError类型的，如果语法有误，err.Error = ErrSyntax；如果结果超出类型范围err.Error = ErrRange；如果参数s合乎语法规则，则返回最接近s值的浮点数（使用IEEE754规范舍入）

```go
func parseDemo() {
    b, err := strconv.ParseBool("true")      // true <nil>
    f, err := strconv.ParseFloat("3.1415", 64)      // 3.1415 <nil>
    i, err := strconv.ParseInt("-2", 10, 64)      // -2 <nil>
    u, err := strconv.ParseUint("2", 10, 64)      // 2 <nil>
}
```

### Format系列函数

将给定类型数据格式化为string类型数据

- `func FormatBool(b bool) string`：返回布尔类型b的字符串表示，true或false

- `func FormatInt(i int64, base int) string`：返回整数类型i的base进制的字符串表示，base必须在2~36，例如FormatInt(10, 2)会将二进制的10转换为对应的字符串表示，结果中会使用小写字母a~z表示大于10的数字

- `func FormatUint(i uint64, base int) string`：返回无符号整数类型i的base进制的字符串表示，是FormatInt的无符号整数版本

- `func FormatFloat(f float64, fmt byte, prec, bitSize int) string`：返回浮点数类型f的字符串表示，fmt表示格式，其中'f'(-ddd.dddd)、'b'(-ddddp+-ddd，二进制指数)、'e'(-d.dddde+-dd, 十进制指数)、'E'(-d.ddddE+-dd, 十进制指数)、'g'(指数很大时用'e'格式，否则'f'格式)、'G'(指数很大时用'E'格式，否则用'f'格式)，prec表示精度（排除指数部分），对于'f'、'e'、'E'，他表示小数点后的数字个数，对于'g'、'G'，他控制总的数字个数，若prec为-1，则代表使用最少数量的、但又必须的数字来表示f，bitSize表示f的来源类型（32：float32、64：float64），会据此进行舍入

```go
func foramtDemo() {
    s1 := strconv.FormatBool(true)                  // "true"
    s2 := strconv.FormatFloat(3.1415, 'E', -1, 64)  // ""3.1415E+00"
    s3 := strconv.FormatInt(-10, 2)                 // "-1010"
    s4 := strconv.FormatUint(10, 2)                 // "1010"
}
```

### 其他

1. isPrint函数

`func IsPrint(r rune) bool`：返回一个字符是否可以输出，和unicode.IsPrint一样，参数r只能是字母（广义）、数字、标点符号和ASCII空格

2. CanBackquote

`func CanBackquote(s string) bool`：判断字符串s是否可以不被修改的表示为一个单行的、没有空格和tab之外控制字符的反引号字符串

### 其他Append系列函数、Quote系列函数，可看官方文档

## net/http包

包含了完整的HTTP客户端和服务端的实现

### HTTP客户端









## Context包







# Go常用第三方库

## gin框架



## MySQL


## sqlx


## Redis


## MongoDB


## etcd


## Zap日志库


## Viper


## singleflight包


## Wire


## gRPC























## bufio

当频繁地对少量数据读写时会占用IO，造成性能问题，Golang的`bufio`库使用缓存来一次性进行大块数据的读写，以此降低IO系统调用，提升性能

bufio实现了有缓冲的IO，它包装了一个io.Reader或io.Writer接口对象，闯将另一个也实现了该接口，且同时还提供了缓冲和一些文本IO的帮助函数的对象

bufio包的读写模块提供了针对字节或字符串类型的缓冲机制，因此**很适合用于读写UTF-8编码的文本文件**

### type Reader

```go
type Reader struct {
    buf          []byte
    rd           io.Reader // reader provided by the client
    r, w         int       // buf read and write positions
    err          error
    lastByte     int // last byte read for UnreadByte; -1 means invalid
    lastRuneSize int // size of last rune read for UnreadRune; -1 means invalid
}
```

### func (*Reader) ReadString

`func (b *Reader) ReadString(delim byte) (line string, err error)`

ReadString读取直到第一次遇到delim字节，返回一个**包含已读取的数据和delim字节**的字符串。如果ReadString在读取到delim之前遇到了错误，它会**返回已读取的数据和该错误**（通常是io.EOF）。当且仅当ReadString返回的数据不以delim结尾时，会返回一个非nil的错误。

### func NewReader

`func NewReader(s string) *Reader`

NewReader创建一个从s读取数据的Reader

示例1：

```go
    // 创建一个从os.Stdin读取数据的Reader
    reader := bufio.NewReader(os.Stdin)
    // 读取数据直到遇到换行符为止
    input, err := reader.ReadString('\n')
```

## math/rand



## strings









