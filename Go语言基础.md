# Go基础语法

本文记录了笔者学习李文周老师《GO语言之路》的学习笔记，记录Go的基础语法，配合使用例子，加深对Go语言基础语法的熟悉，并作为自己的Go基础语法查询手册

## 包与依赖管理

### 包

包可以用来支持模块化开发，实现代码复用

1. 定义包

可以根据需要创建自定义包，包可以理解为一个存放.go文件的文件夹，这个文件夹下面的所有.go文件都要在第一行添加包声明，表示该文件归属的包

一个文件夹直接包含的文件只能属于一个包

同一个包的文件不能在多个文件夹下

包名为main的包是应用程序的入口包，编译后会得到一个可执行文件，不包含main包的源代码编译则不会得到可执行文件

2. 标识符可见性

只有同一个包内声明的标识符才在同一个命名空间下，包外使用包内标识符，需要添加包名前缀，如`fmt.Println("Hello")`

Go中通过标识符首字母大写来指定该标识符对外可见，对于结构体也一样，只有首字母大写的字段才能导出

3. 包的引入

使用import关键字引入，格式为`import importname "path/to/package"`，其中importname通常省略，默认值为引入包的包名，括号内的表示引入包的路径名称

可以为引入的包指定新包名，使用`_`作为新包名即为匿名引入

> 匿名引入主要是为了满足加载需求，使包中的资源初始化，被匿名引入的包中的`init`函数将被执行并且仅执行一遍

```go
import "fmt"			// 可以单个引入
import (
	"encoding/json"		// 可以批量引入
	"os"
)
import f "fmt"			// 为fmt包指定新包名f
import _ "github.com/go-sql-driver/mysql"		// 匿名引入
```

4. init初始化函数

一个包中可以定义任意的init初始化函数，该函数不会接收任何参数，也没有任何返回值，在代码中也不能主动进行调用，他们只会在程序启动时，按声明的顺序自动串行调用执行

每个包在初始化时都是先执行依赖包中声明的init函数，再执行当前包中声明的init函数，以确保在程序的main函数开始执行时，所有的依赖包都已经初始化完成

每个包的初始化都是从初始化包级别变量开始，即顺序为：先初始化当前包级别的变量，接着按依赖包的导入顺序依次进入依赖包中，先初始化依赖包包级别的变量，（若还有依赖包则继续递归进入初始化）再串行调用依赖包中的init函数，当当前包中所有依赖包都初始化完成后，才串行调用当前包中的init函数，完成初始化，最后执行当前包中的main函数执行程序



```go
// init初始化函数定义
func init() {
	...
}
```

### 依赖管理

早期Go语言的依赖包都需要保存在GOPATH目录下，但是这样不支持版本管理，同一个依赖包只能存在一个版本的代码，但是本地的多个项目可能分别依赖同一个包的不同版本

1. go modules

从Go 1.11版本开始有了依赖管理go modules，Go 1.16版本已经默认开启go modules，常见的命令如下

```go
go mod init				// 初始化项目依赖，生成go.mod文件
go mod download 		// 根据go.mod文件下载依赖
go mod tidy 			// 比对项目文件中引入的依赖与go.mod文件中引入的依赖，并下载缺失的依赖
go mod graph			// 输出依赖关系图
go mod edit 			// 编辑go.mod文件
go mod vendor			// 将项目所有依赖导出到vendor目录
go mod verify			// 校验依赖包是否被篡改过
go mod why 				// 解释为什么需要某个依赖
```

Go语言在go modules的过渡阶段使用环境变量GO111MODULE作为启用go modules功能的开关

- GOPROXY

用于设置模块代理，从代理地址拉去和校验依赖包，能够加快速度，其默认值是https://proxy.golang.org,direct，目前更常用https://goproxy.cn和https://goproxy.io，可以通过下面的方式设置GOPROXY，`go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct`，最后的`direct`用于指示Go回到源地址抓取，类似默认抓取地址

- GOPRIVATE

当项目引用非公开的包（如公司内部git仓库或github私有仓库）时，便无法从代理中拉取依赖包，GOPRIVATE用来告诉Go哪些仓库是私有的，不必通过代理服务器拉去和验证，比如将公司内部的代码仓库设置到GOPRIVATE中`go env -w GOPRIVATE="git.mycompany.com"`，这样就可以正常拉取以`git.mycompany.com`为路径前缀的依赖包

2. 使用go modules引入包

创建一个testGoModules项目，在testGoModules文件夹下执行`go mod init testGoModules`初始化项目，得到go.mod文件，该文件会记录项目使用的第三方依赖包信息，包括包名和版本

```go
module testGoModules	// 用于定义当前项目的导入路径

go 1.23.0				// 标识当前项目使用的Go版本
```

接着在项目目录下创建一个main.go文件，并为testGoModules项目引入第三方包github.com/q1mi/hello来实现功能，需要先将依赖包下载到本地，并在go.mod中记录依赖信息，才能在代码中引入并使用这个包，下载依赖有两种方式

**第一种方式**是使用`go get`命令手动下载`go get -u github.com/q1mi/hello`，默认下载最新版本

也可以指定下载的版本`go get -u github.com/q1mi/hello@v0.1.0`，若依赖包没有发布任何版本，则会拉取最新的提交

也可以通过指定commit hash来下载某个commit对应的代码（一般写hash的前7位即可）`go get github.com/q1mi/hello@2ccfadd`

此时go.mod中就可以看到下载的依赖包及版本信息（按hash下载，go.mod文件中会按默认版本号、最近以此commit时间和hash组成的版本格式）

行尾的`indirect`表示当前程序中所有import语句都没有发现这个包

> 使用go get命令下载一个新的依赖包时一般会额外添加-u参数，用于强制更新现有依赖

```go
module testGoModules

go 1.23.0

require github.com/q1mi/hello v0.1.1 // indirect

module testGoModules

go 1.23.0

require github.com/q1mi/hello v0.1.2-0.20210219092711-2ccfaddad6a3 // indirect
```

**第二种方式**是是直接编辑go.mod文件，写入依赖包和版本信息，如`require github.com/q1mi/hello v0.1.1`或者指定commit进行下载`require github.com/q1mi/hello 2ccfadda`，然后再在项目目录下执行`go mod download`下载依赖包，随后go.mod文件机就会更新对应的版本信息

```go
package main

import (
	"github.com/q1mi/hello"
	"fmt"
)

func main() {
	fmt.Println(hello.SayHi())		// 调用hello包的SayHi函数
}
```

通过以上两种方式下载好依赖包之后，就可以在main.go文件中引入并使用了

在大型项目中，通常会按功能或业务定义多个不同包，如在testGoModules目录下创建一个新的包summer，相当于创建一个文件夹summer，其中的go文件的package都要为summer

此时若想要在当前项目目录下的其他包或者main.go中调用Diving函数，需要以项目的导入路径为前缀导入该包`import "testGoModules/summer"		// 导入当前项目下的包`

```go
package summer

import (
	"fmt"
)

func Diving() {
	fmt.Println("diving")
}
```

若想导入一个没有被发布到其他任何代码仓库的本地包，可以使用replace语句将依赖临时替换成本地包的相对路径，例如在testGoModules的同级目录下有另外一个名为overtime的项目，overtime包只存在于本地，并不能通过网络获取，可以在testGoModules的go.mod文件中正常引入overtime包，然后添加replace语句`replace overtime => ../overtime`，这样就可以在testGoModules项目中引入overtime包了

经常使用replace将项目依赖的某个包替换为其他版本的代码或自己修改的代码包

```go
module testGoModules

go 1.23.0

require github.com/q1mi/hello v0.1.2-0.20210219092711-2ccfaddad6a3 // indirect
require overtime v0.0.0
replace overtime => ./overtime

// testGoModules的main.go文件
package main

import (
	"overtime"

	"github.com/q1mi/hello"
)

func main() {
	hello.SayHi() // 调用hello包的SayHi函数
	overtime.Do()
}
```

通过go modules下载依赖后，项目目录还会生成一个go.sum文件，该文件详细记录了当前项目引入的依赖包的信息及其hash值，用于防止依赖包被非法篡改，通过go.sum机制对依赖包进行校验

下载的依赖都带有版本号以示区别，因此可以在本地保存同一个包的不同版本，依赖都保存在`$GOPATH/pkg/mod`目录下，可以执行`go clean -modcache`命令清除本地缓存的依赖包数据

3. 使用go modules发布包

首先需要创建一个github仓库，并下载到本地`git clone https://github.com/q1mi/hello`，在hello目录下进行初始化`go mod init github.com/q1mi/hello`，然后创建hello.go文件，编写代码实现，然后将改项目的代码push到远程代码仓库，这样就对外发布了一个go包，其他开发者可以通过`github.com/q1mi/hello`引入路径下载并使用，一个完善的包应该使用git tag为代码包打上标签`git tag -a v0.1.0 -m "release version v0.1.0"`

go modules使用语义化版本控制，建议的版本号格式为`v1.2.3`其中1为主版本号，2为此版本号，3为修订号

主版本号：在发布了不兼容的版本迭代时递增
次版本号：在发布了功能性更新时递增
修订号：在发布了bug修复类更新时递增

当要发布新版本时，通常会修改当前包的引入路径，即修改go.mod文件中的module字段`module github.com/q1mi/hello/v2`，然后将修改后的代码打好标签并提交到远程仓库，这样就能在不影响使用旧版本用户的前提下，发布新的版本，对于想使用v2版本代码包的用户，只需要从修改后的引入路径下载即可`go get github.com/q1mi/hello/v2@v2.0.0`，然后进行引用`import "github.com/q1mi/hello/v2"`，这样就可以使用v2版本的代码包了

当我们发布的某个版本存在致命缺陷时，可以在go.mod中使用retract声明废弃，用户使用go get下载v0.1.2版本时就会收到提示，催促其升级到其他版本

```go
module github.com/q1mi/hello
go 1.16
retract v0.1.2
```

## 变量与常量

1.  标识符

Go语言的标识符由字母(A~Z, a~z)、数字(0~9)、下划线(_)组成，区分大小写，不能以数字开头，只能以字母和下划线开头，Go语言中使用**驼峰命名法**命名变量，例如使用`ListenAndServe`和`ReadFromFile`

2. 数据类型

![数据类型](./images/datatype1.jpg)

![数据类型](./images/datatype2.jpg)

```go
func main () {
	// 复数型
	var c1 complex128
	c1 = 3 + 4i

	// 构造复数
	c := complex(1.2, 3.4)
	// 获取实部
	fc1 := real(c)
	// 获取虚部
	fc2 := imag(c)

	// 数字字面量
	v1 := 0b00101101 // 二进制101101，相当于十进制45
	v2 := 0o755 // 八进制755，相当于十进制493
	v3 := 0x2d // 十六进制2d，相当于十进制45

	// 定义一个多行字符串
	str := `
		这是一个多行字符串
		这是第二行
		这是第三行
	`
}
```

- 字符串修改

字符串的底层实现是一个结构体，包含指向底层字节数组的指针和字节数组的长度

```go
type stringStruct struct {
	str unsafe.Pointer
	len int
}
```

要修改字符串，需要先将其转换成`[]rune`或`[]byte`，然后转换成string，无论哪种转换，都会重新分配内存，并复制字节数组

```go
func main () {
	s1 := "big"
	byteS1 := []byte(s1)
	byteS1[0] = 'p'

	s2 := "白萝卜"
	runeS2 := []rune(s2)
	runeS2[0] = '红'
}
```

组成每个字符串的元素是字符，用单引号包裹，Go中字符有两种类型，byte代表ASCII码字符，rune代表UTF-8字符，这两个本质上都是整型，只是为了方便编程，Go语言才给出这两个别名，其中byte是uint8的别名，rune是int32的别名

- 类型转换

数值类型之间可以互相转换，但字符串类型不能与布尔类型、数值类型相互转换

强制类型转换存在溢出情况，当字节长度大的数值转换为字节长度小的数值时，大数值的高位被截位，导致数据精度丢失；浮点数转换为整数时，小数部分会被截掉；可以使用`strconv`包进行转换，该包实现了基本数据类型和其字符串表示的相互转换

```go
func main () {
	var a = 1
	c := float32(a)
	// 整数转换成字符串
	i := 123
	// 使用fmt.Sprintf返回一个格式化的字符串
	s1 := fmt.Sprintf("%d", i)
	// 使用strconv.Itoa返回一个字符串
	s2 := strconv.Itoa(i)
	// 使用strconv.FormatInt返回一个字符串
	s3 := strconv.FormatInt(int64(i), 10)
	// 使用strconv.FormatUint返回一个字符串
	s4 := strconv.FormatUint(uint64(i), 10)
	
	// 字符串转换成整数
	s := "123"
	// 使用strconv.Atoi返回一个整数
	i1, _ := strconv.Atoi(s)
	// 使用strconv.ParseInt返回一个整数
	i2, _ := strconv.ParseInt(s, 10, 64)
	// 使用strconv.ParseUint返回一个无符号整数
	i3, _ := strconv.ParseUint(s, 10, 64)

	// 字符串转换成浮点数
	s = "123.456"
	// 使用strconv.ParseFloat返回一个浮点数
	f, _ := strconv.ParseFloat(s, 64)
}
```



3. 基本使用

- 变量声明

```go
func main() {
	// 变量声明
	var name string
	var age int
	var isOk bool

	// 批量声明
	var (
		name1 string
		age1  int
		isOk1 bool
	)
}
```

- 变量初始化

```go 
func main() {
	// 声明时指定初始值
	var name string = "zhangsan"
	
	// 一次初始化多个变量
	var name, age = "zhangsan", 20

	// 省略类型，编译器自动推导类型
	var name = "zhangsan"
	var age = 20

	// 语法糖，简短声明并初始化方式
	name := "zhangsan"
	age := 20
}
```

- 常量

声明后恒定不变的值

```go
func main() {
	// 常量声明
	const pi = 3.1415

	// 一次性声明多个常量
	const (
		pi = 3.1415
		e = 2.71828
	)

	// 可省略值
	const (
		n1 = 100
		n2		// 100
		n3		// 100
	)

	// 使用预声明标识符
	const (
		leve10 = iota	//0
		leve11 = iota	//1
		// 可省略
		leve12			//2
		// 可跳过
		_
		leve13			//4
		// 可插队
		leve14 = 100	//100
		leve15			//6
		// 一行定义多个
		leve16, leve17 = iota, iota	//7, 8
		leve18, leve19	//8, 9

		// 定义数量级
		_ = iota
		KB = 1 << (10 * iota)	//1 << (10 * 1)
		MB = 1 << (10 * iota)	//1 << (10 * 2)
		GB = 1 << (10 * iota)	//1 << (10 * 3)
		TB = 1 << (10 * iota)	//1 << (10 * 4)
	)
}
```

## 条件语句

Go中规定与if匹配的左括号必须与if在同一行，与else匹配的左括号必须与else在同一行，同时，else必须与上一个if或else if右边的右括号在同一行

if或else if其后的条件语句不需要括号

特殊用法：可以在if表达式之前添加一个执行语句，再根据变量值进行判断

```go
func main() {
	score := 100	// score作用域是整个函数
	if score >= 90 {
		fmt.Println("优秀")
	} else if score >= 80 {
		fmt.Println("良好")
	} else {
		fmt.Println("及格")
	}

	// 特殊用法
	if score := 90; score >= 90 {	// score作用域是if语句块
		fmt.Println("优秀")
	} else if score >= 80 {
		fmt.Println("良好")
	} else {
		fmt.Println("及格")
	}
}
```

## 循环语句

Go中只有for一种循环语句，所有循环类型均可以使用for来实现

for后无需括号

```go
func main() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}

	// 省略初始语句
	i := 10
	for ; i < 20; i++ {
		fmt.Println(i)
	}

	// 省略初始语句和结束语句，实现while循环
	for i < 30 {
		fmt.Println(i)
		i++
	}

	// 使用break、goto、return、panic语句强制退出for循环，或者通过continue跳出本次循环进入下一次循环
	for i := 0; i < 10; i++ {
		if i == 5 {
			break
		}
		if i == 3 {
			continue
		}
		if i == 7 {
			panic("panic")
		}
		if i == 8 {
			goto label
		}
	}
label:
    fmt.Println("hello")

	for {
		fmt.Println("无限循环")
	}
}
```


## switch分支

switch可用于多条件判断，化简条件判断分支过多的情况

Go中的`fallthrough`语法可以在执行完满足条件的case分支后继续执行下一个分支，且`fallthrough`语句只能作为case分支的最后一个语句出现

```go
func main() {
	score := 100
	switch score {
	
	// 可使用多值判断
	case 10, 20, 30:
		fmt.Println("10, 20, 30")
	
	// 使用fallthrough
	case 40:
		fmt.Println("40")
		fallthrough		// 此后的语句都会执行，输出50、60、default
	case 50:
		fmt.Println("50")
		fallthrough
	case 60:
		fmt.Println("60")
	default:
		fmt.Println("default")
	}
	
	// 使用表达式做判断
	switch {
		case score == 100:
			fmt.Println("100")
		case score == 200:
			fmt.Println("200")
		default:
			fmt.Println("default")
	}
}
```

## 标签

Go中标签被定义后必须被使用，标签名允许与变量名重复，但标签名不能重复

- goto

goto语句通过标签在代码间无条件跳转，能快速跳出循环、避免重复退出，能简化一些代码的实现过程

```go
func main() {
	// 使用goto跳出双层for循环
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j == 20 {
				// 跳转到退出标签
				goto breakTag
			}
			fmt.Println(i, j)
		}
	}
	return		// 若不加return，则后面breakTag标签没有goto跳转也会被执行

breakTag:
	fmt.Println("breakTag")
}
```

- break

break语句可以结束for、switch和select语句，但不能用于goto语句，可以在break语句后面添加标签，表示退出某个标签对应的代码块，标签名必须定义在对应的for、switch和select语句块之前

```go
func main() {
breakLabel:
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j == 2 {
				fmt.Println(j)
				break breakLabel // 跳出，不再执行breakLabel下的循环语句，即只会输出2
			}
		}
		fmt.Println(i)
	}
}
```

- continue 

continue语句可以结束当前循环，开始下一次循环，但是仅限在for循环中使用，在continue语句后添加标签，表示继续标签对应的下一次循环

```go
func main() {
breakLabel:
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j == 2 {
				fmt.Println(j)
				continue breakLabel //跳出内层循环，继续执行外层循环，即只会打印十个2
			}
		}
		fmt.Println(i)
	}
}
```

## 数组

数组的大小从声明时就确定了，使用时可以修改数组成员，但是不能修改数组大小，因此平常不常使用数组，而更多使用切片，其长度可以变化

数组是占用一块连续内存的值类型，赋值操作和函数传参会复制整个数组得到一个副本，修改副本的值不会改变原始变量的值

```go
func testArray(a [2]int) {
	a[0] = 100
	fmt.Println(a)		// [100 2]
}

func main() {
	a := [2]int{1, 2}
	testArray(a)
	fmt.Println(a)		// [1 2]
}
```

- 数组的声明

`[3]int`和`[4]int`是不同类型，不能直接进行比较和赋值

```go
func main() {
	// 声明一个数组
	var a [3]int
	var b [4]int

	var n = 5
	// 数组的长度必须是常量或常量表达式
	// var c [n]int		// error

	// 可以通过索引访问数组成员
	fmt.Println(a[1], b[2])

	// 使用内置len函数获取数组中元素数量
	fmt.Println(len(a), len(b))
}
```

- 数组的初始化

```go
func main() {
	var a [3]int					// 初始化为int类型的零值
	
	var b = [3]int{1,2,3}			// 使用指定初始值初始化数组

	var c = [...]int{1,2,3,4,5}		// 编译器根据初始值的数量自行推断数组的长度

	d := [...]int{1: 1, 3: 5}		// 只初始化部分值
}
```

- 数组的遍历

```go
func main() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// for循环遍历
	for i := 0; i < len(array); i++ {
		fmt.Println(array[i])
	}

	// for range遍历
	for index, value := range array {
		fmt.Println(index, value)
	}
}
```

- 多维数组

```go
func main() {
	a := [2][4]int{{1, 2, 3, 4}, {5, 6, 7, 8}}
	b := [2][4]string{
		{"1", "2", "3", "4"},
		{"5", "6", "7", "8"},
	}

	for _, v1 := range b {			// 获取内层数组
		for _, v2 := rangge v1 {	// 遍历内层数组
			fmt.Println(v2)
	}

	// 多维数组只有第一层可以让编译器自动推导数组长度
	c := [...][2]string{
		{"1", "2"},
		{"3", "4"},
		{"5", "6"}
	}
}
```

## 切片

切片是可以存储**相同类型元素**的**长度可变**的序列，是基于数组类型封装的，就像一个长度可变的滑动窗口在底层数组上移动，可以包含数组的全部或部分元素

多个切片可能对应同一个底层数组，切片之间可以重叠，因此修改切片时，有可能影响拥有相同底层数组的其他切片

一个完整切片需要包含三个要素

1. 起始元素的地址：指切片中的底层数组中第一个元素的地址

2. 长度len：指切片中元素的数量，内置len函数获取切片长度

3. 容量cap：指切片中第一个元素到可能达到的最后一个元素的数量，内置cap函数获取切片容量

- 创建切片

1. 切片字面量

```go
func main() {
	a := []int{1, 2, 3}
	b := []string{"1", "2", "3"}
	c := []int{1, 2, 3, 99: 100}	// 使用索引方式指定初始值
	fmt.Println(len(c), cap(c))		// 100 100
}
```

2. 切片表达式

基于数组通过切片表达式得到切片，索引low和high必须满足`0 <= low <= high <= len(array)`，否则会索引越界，左闭右开

由于切片的底层是数组，因此**对数组取切片得到的是数组**，也即是一个切片，因此存在容量cap的概念；而**对字符串取切片得到的是字符串**，没有容量cap的概念

```go
func main() {
	a := [5]int{1, 2, 3, 4, 5}
	s1 := a[1:3]		// [2 3]
	fmt.Printf("s1: %v type: %T len: %v cap: %v\n", s1, s1, len(s1), cap(s1))	// [2 3] []int 2 4，长度为2容量为4

	// 可省略
	a[2:]		// a[2:len(a)]
	a[:3]		// a[0:3]
	a[:]		// a[0:len(a)]

	b := "hello world"
	s2 := b[1:3]		// "el"
	fmt.Printf("s2: %v type: %T len: %v cap: %v\n", s2, s2, len(s2), cap(s2))	// cap(s2)报错，字符串类型没有容量的概念
}
```

**对切片取切片得到的是切片**，此时high的上限是切片的容量cap(a)

```go
func main() {
	a := [5]int{1, 2, 3, 4, 5}

	s := a[1:3] // [2, 3]

	s1 := s[3:4] // 切片的切片，high可以到切片的容量

	fmt.Printf("s1: %v len: %v cap: %v\n", s1, len(s1), cap(s1)) // s1: [5] len: 1 cap: 1
}
```

**完整切片表达式**，除了low和high，增加一个max用来表示结果的容量，此时会将结果切片的容量设置为`max-low`，字符串不支持完整切片表达式，需要满足`0 <= low <= high <= max <= len(array)`，否则会索引越界

```go
func main() {
	a := [5]int{1, 2, 3, 4, 5}

	s := a[1:3:4] // 额外指定max，控制切片容量

	fmt.Printf("s: %v len: %v cap: %v\n", s, len(s), cap(s)) // [2, 3] len: 2 cap: 3

	s2 := s[3:4]

	fmt.Printf("s2: %v len: %v cap: %v\n", s2, len(s2), cap(s2)) // 报错，由于s指定了容量3，此时s已经没有包含5这个元素了

	s1 := a[1:3]

	fmt.Printf("s1: %v len: %v cap: %v\n", s1, len(s1), cap(s1)) // [2, 3] len: 2 cap: 4
}
```

3. make创建切片

用于动态创建一个切片，使用make函数可以将切片所需容量一次申请到位，切片的元素数量小于或等于该切片的容量

```go
func main() {
	s1 := make([]string, 2)			// cap默认与长度相同，len2，cap2
	s2 := make([]string, 2, 10)		// len2, cap10
}
```

- 空切片

nil切片不仅没有大小和容量，也没有指向任何底层的数组

empty切片没有大小和容量，但是指向了一个特殊的内存地址zerobase，即所有0字节分配的基地址

要判断一个切片是否为空，只能用`len(s) == 0`，

切片中的元素不是直接存储的值，因此不允许切片之间直接使用`==`进行比较，切片唯一合法的比较操作是和nil比较

```go
func main() {
	var s1 []int			// nil切片
	s2 := []int{}			// empty切片
	s3 := make([]int, 0)	// empty切片
}
```

- 切片操作

1. 切片遍历

```go
func main() {
	s := []int{1, 2, 3}
	// 索引遍历
	for i := 0; i < len(s); i++ {
		fmt.Println(s[i])
	}

	// range遍历
	for index, value := range s {
		fmt.Println(index, value)
	}
}
```

2. 切片追加

`append`函数支持在切片末尾追加一个，也可以一次性追加多个，也可以追加另一个切片中的元素，但是需要配合`...`符号将被添加的切片展开

Go编译器不允许调用`append`函数后不适用其返回值，因此通常使用原变量接收`append`函数的返回值

```go
func main() {
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5, 6}
	s1 = append(s1, 7)		// [1 2 3 7]
	s2 = append(s2, 8, 10)	// [4 5 6 8 10]
	s1 = append(s1, s2...)	// [1 2 3 7 4 5 6 8 10]

	// append函数中可以直接使用nil slice
	var s []int  			// nil slice
	s = append(s, 1, 2, 3)

	// 使用make时要注意是否初始化时指定初始大小
	s3 := make([]int, 3)
	s3 = append(s3, 1, 2, 3)		// [0 0 0 1 2 3]

	s4 := make([]int, 0, 3)
	s4 = append(s4, 1, 2, 3)		// [1 2 3]
}
```

3. 切片扩容策略

向切片中添加元素时，若底层数组不能容纳新增的元素，切片会自动按照一定的策略进行扩容，此时切片指向的底层数组会更换

- Go1.21.1版本的切片扩容策略

	- 若newLen > doublecap，则newcap = newLen
	- newLen <= doublecap，并且oldCap < 256，则newcap = doublecap
	- oldCap >= 256，则最终容量newcap开始循环增加(newcap + 3 * threshold) / 4，直到newcap >= newLen

4. 切片复制拷贝

简单的复制会导致两个切片共享同一个底层数组，此时对一个切片进行修改会影响另一个切片的内容

```go
func main() {
	s1 := make([]int, 3)
	s2 := s1 				// 将s1赋值给s2，s1和s2共用一个底层数组
	s2[0] = 100				// 修改s2的值，s1的值也会变
	fmt.Println(s1, s2)		// [100 0 0] [100 0 0]
}
```

可以使用`copy`函数将一个切片的数据复制到另一个切片空间，此时两个切片不会共享同一个底层数组，因此可以安全的对一个切片进行修改而不会影响另一个切片

```go
func main() {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := make([]int, 5, 5)

	copy(s2, s1)				// 将切片s1中的元素复制到切片s2中
}
```

5. 切片删除元素

Go语言没办法删除切片中间某个元素并保持剩余元素顺序的办法，只能通过切片操作将删除的元素前后的元素顺序连接起来

```go
func main() {
	s := []int{1, 2, 3}
	// 删除索引1元素
	s1 := append(s[:1], s[2:]...)

	s2 := copy(s[:1], s[2:])
}
```

## 映射map

Go的映射是一种无序的键值对类型，一个映射的所有键都必须是相同的类型，值也必须是相同的类型，但是键和值的类型可以不同

只有能使用`==`操作符进行比较的类型才能作为映射的键，任何类型都可以作为映射的值，不建议使用浮点型作为映射的键

- 创建映射

1. 使用映射字面量的方式创建携带初始值的映射变量

默认初始值为nil，不能直接对nil映射添加元素，因为还没有初始化


2. 使用make函数创建映射

使用make函数创建映射时，不管有没有指定容量，初始值都不再是nil

使用make函数创建映射时，容量参数不是必须的，映射会根据需要自动扩容，但是我们应该在初始化映射时就指定一个合适的容量，这样能有效减少运行时的内存分配，提高代码的执行效率

```go
func main() {
	// 默认初始值为nil
	var m1 map[string]string

	m1["zhangsan"] = "zhangsan" // 报错

	m2 := map[string]string{
		"zhangsan": "zhangsan",
		"lisi":     "lisi",
	}

	// 默认初始值不为nil
	m3 := make(map[string]string)

	m3["zhangsan"] = "zhangsan" // 正常
}
```

- 映射的基本使用

访问、删除、判断键是否存在、遍历

可以直接使用map[key]的方式访问元素，若key不存在，返回值类型对应的零值

可以使用len函数来获取映射中的元素数量

Go内置delete函数用来从映射中删除一组键值对

虽然通过map[key]访问映射中的元素时，若键不在映射中会返回对应类型的零值，但是这样不太友好，Go语言中有专门的写法来判断映射中是否存在某个键`value, ok := map[key]`

同样可以使用range来遍历映射的键和值，但是映射的遍历结果是不固定的，若要按照指定的顺序遍历映射，则可以先将映射的所有键取出放到一个切片中，再对切片进行排序，将遍历该切片的元素作为键去访问映射，间接实现按顺序访问的效果

```go
func main() {
	m := map[string]string{
		"zhangsan": "zhangsan",
		"lisi":     "lisi",
	}
	fmt.Println(len(m)) // 2
	// 添加元素
	m["zhangsi"] = "zhangsi"
	// 访问元素
	fmt.Println(m["zhangsan"])
	fmt.Println(len(m)) // 3

	// 删除元素
	delete(m, "zhangsi")

	// 判断键值对是否存在
	value, ok := m["zhangsi"]

	// 按顺序遍历映射
	// 先取出所有key存入切片中
	keys := make([]string, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}
	// 对切片进行排序
	sort.Strings(keys)
	// 按顺序遍历映射
	for _, key := range keys {
		fmt.Println(key, m[key])
	}
}
```

- 切片作为映射的值

感觉非常适合用来存储一些人员信息，比如一个切片就是一个人员，这个切片中可以存放很多个映射的数据

```go
func main() {
	// 映射类型的切片
	m := make([]map[string]string, 0, 10)

	// 对切片中的映射元素进行初始化
	m[0] = make(map[string]string, 2)
	m[0]["a"] = "a"
	m[0]["b"] = "b"
}
```

- 映射类型的切片

```go
func main() {
	// 值为切片类型的映射
	// 键为字符串，值为字符串切片的映射
	m := make(map[string][]string, 3)
	key := "zhangsan"
	value, ok := m[key]
	if !ok {
		fmt.Println("not found")
	}
	// 找到的话，值就是一个切片
	value = append(value, "lisi")
	m[key] = value
}
```


## 函数

函数使用func关键字声明，`func 函数名 (参数) (返回值)`Go语言支持多个返回值

Go中函数调用过程传参是按值传递，函数接收到的是复制的实参副本，而当传入函数的实参为某些特定类型（指针、切片、映射、函数、通道）时，函数接收到的是实参的引用，此时在函数中修改形参将对实参造成影响

- 函数的一些特殊用法

1. 类型简写

当参数中相邻变量的类型相同时可以省略参数的类型，`func Sum(x, y int) int {}`

2. 可变参数

指函数的参数数量不固定，Go中可变参数通过在形参后添加`...`来标识，通常作为函数的最后一个参数，可变参数本质上是一个切片

```go
// 有多个返回值时，在函数声明中必须用括号将两个返回值的类型括起来
func sum(x int, y ...int) (int, int) {
	summ := 0
	for _, v := range y {		// y是一个切片
		summ += v
	}
	return x, summ
}

res1 := sum(10)
res2 := sum(10, 20)
res3 := sum(10, 20, 30)
fmt.Println(res1, res2, res3)	// 10 30 60
```

3. 命名返回值与特殊返回值

函数定义阶段可以直接给返回值命名，并在函数体中直接使用这些变量

在某些场景中，nil可以作为特殊返回值返回

```go
func sum(x int, y int) (sum int) {
	sum = x + y
	return
}

func someFunc(x string) []int {
	if x == "" {
		return nil	// 没必要返回[]int{}
	}
	return nil
}
```

- 变量作用域

全局变量定义在函数外部，在整个程序运行周期都有效

局部变量又分为：函数内局部变量和语句块局部变量

若局部变量与全局变量重名，则优先访问局部变量（就近原则）

- 函数类型与变量

函数的类型通常称为函数签名，拥有相同参数列表和返回值列表的函数类型相同

```go
// 以下三个函数，函数签名相同
func f1(x int, y int) int {}
func f2(x, y int) int {}
func f3(x, y int) (res int) {}
```

Go语言中函数也可以是一种类型，可以将函数类型的值保存在变量中，函数类型的零值为nil

```go
func sayHello(name string) {
	fmt.Println("Hello,", name)
}

func main() {
	var f func(string)			// 声明函数类型变量f
	f = sayHello				// 将sayHello赋值给f
	f("Go")						// 输出：Hello, Go
}
```

可以使用type关键字为函数类型定义类型别名

如，定义了一个calculation类型，他是`func(int, int) int`类型的别名，只要函数签名满足`func(int, int) int`，那么该函数就可以赋值给calculation类型的变量，并进行调用

```go
type calculation func(int, int) int
func add(x, y int) int {
	return x + y
}
func sub(x, y int) int {
	return x - y
}

func main() {
	var c calculation
	c = add
	fmt.Println(c(1, 2))		// 相当于调用add(1, 2)

	s := sub
	fmt.Println(s(1, 2))		// 相当于调用sub(1, 2)
}
```

- 高阶函数

指以其他函数为参数，或返回一个函数作为结果的函数

1. 函数作为参数

```go
func add(x, y int) int {
	return x + y
}
func sub(x, y int) int {
	return x - y
}

// 将两个整数传入给定op函数进行计算并返回结果
func calc(x, y int, op func(int, int) int) int {
	return op(x, y)
}

func main() {
	addres := calc(10, 20, add)		
	subres := calc(10, 20, sub)
}
```

2. 函数作为返回值

```go
func do(s string) (func(int, int) int, error) {
	switch s {
	case "+":
		return func(i, j int) int { return i + j }, nil
	case "-":
		return func(i, j int) int { return i - j }, nil
	case "*":
		return func(i, j int) int { return i * j }, nil
	default:
		err := errors.New("invalid operator")
		return nil, err
	}
}

func main() {
	f, err := do("+")
	if err != nil {
		fmt.Println(err)
		return
	}
	res := f(1, 2)
	fmt.Println(res)
}
```

- 匿名函数和闭包

匿名函数可以看作没有函数名的函数，因此没办法像普通函数那样被调用，需要保存到某个变量中或者作为立即执行函数

匿名函数多用于实现回调函数和闭包

```go
func main() {
	// 将匿名函数保存到变量中
	add := func(x, y int) {
		fmt.Println(x + y)
	}
	add(10, 20)

	// 自执行函数：匿名函数定义玩笔加()直接执行
	func(x, y int) {
		fmt.Println(x + y)
	}(10, 20)
}
```

闭包指一个函数和与其相关的引用环境组合而成的实体，即闭包=函数+引用环境，即闭包是一个引用了外层函数局部变量的函数

```go
func addrer() func(int) int {
	var x int
	return func(y int) int {
		// 将y的值累加到其外层函数的局部变量x
		x += y
		return x
	}
}

// 闭包的x变量也可以是外层函数的参数，在函数调用时由外部传入，更加灵活
func addrer2(x int) func(int) int {
	return func(y int) int {
		x += y
		return x
	}
}

func main() {
	var f = addrer()			// 在f的作用域中，闭包中的x也一直有效
	fmt.Println(f(1))			// 1	
	fmt.Println(f(2))			// 3
	fmt.Println(f(3))			// 6
}
```

一个闭包使用的实际性例子，理解上来说，可以实现对一个函数的复用，更常见的做法是将外部变量传入闭包中，从而使得该闭包在其生命周期内都都与该外部变量有关

```go
// makeSuffixFunc, 返回一个给文件名添加指定后缀的函数
func makeSuffixFunc(suffix string) func(string) string {
	return func(name string) string {
		// 判断文件名是否包含指定后缀
		if !strings.HasSuffix(name, suffix) {
			return name + suffix
		}
		return name
	}
}

func main() {
	// 创建一个添加.jpg后缀的函数
	jpgFunc := makeSuffixFunc(".jpg")
	// 创建一个添加.txt后缀的函数
	txtFunc := makeSuffixFunc(".txt")
	// 给文件名avatar添加.jpg后缀
	fmt.Println(jpgFunc("avatar")) // avatar.jpg
	// 给文件名readme添加.txt后缀
	fmt.Println(txtFunc("readme")) // readme.txt
}
```

- 内置函数

1. defer函数

defer语句能够延迟函数调用，这个被延迟调用的函数的所有返回值都需要被丢弃

被延迟处理的语句将按defer定义的逆序执行，即先被defer的语句最后被执行，最后被defer的语句最先被执行，类似于栈先进后出

defer函数常用于在读取文件内容时，延迟执行`f.Close()`，能够保证在函数运行玩笔后妥善地关闭文件，释放相关资源

```go
func readFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()		// 函数结束时关闭文件
	return nil
}
```

defer函数也经常用来做一些调试工作，如在函数的入口和出口执行特定的语句，记录执行时间和相关变量等

```go
func recordTime(funcName string) func() {
	start := time.Now()
	fmt.Printf("enter %s\n", funcName)
	return func() {
		fmt.Printf("exit %s, spend %v\n", funcName, time.Since(start))
	}
}

func readFromFile(filename string) error {
	defer recordTime("readFromFile")() // defer延迟调用，虽然recordTime函数return的没有变量进行接收，但是也会执行打印输出
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close() // 函数结束时关闭文件
	return nil
}
```

return语句不是原子操作，包含两步，首先是返回赋值，再执行RET指令进行返回，而defer的执行处于这两个步骤之间，因此可以更新函数的返回值变量

```go
// 首先将返回值赋给返回变量，接着修改x的值，最后返回，因此不会更新返回值
func f1() int {
	x := 5
	defer func () {
		x++				// 只修改x的值，不会更新返回值
	}()

	return x			// 返回值=5
}

// 首先x=5，然后执行defer会更新x的值，最后返回，因此此时会更新返回值，因为提前声明好了返回值变量
func f2() (x int) {
	defer func() {
		x++				// 会更新返回值
	}()
	return 5			// 返回值 x = 5
}
```

2. panic函数

当一个函数调用发生panic时，会立即进入退出阶段，开始执行所有延迟调用函数，然后程序会异常退出并输出相应的日志，日志中会包含函数调用信息和具体的panic值

panic函数可以接收任意值作为参数，但通常使用相关错误信息作为参数

panic会导致程序异常退出，因此一般只在程序出现不应该出现的逻辑错误时才会使用，应该使用Go语言提供的错误机制，而不是滥用panic机制

对于一个函数其非法的情况是不可能出现的，因此通常我们没必要去处理这种不可能出现的错误，此时就会用panic，并约定熟成的使用Muust作为函数名称的前缀

```go
// Compile函数是regexp包中用于将字符串参数编译成正则表达式的函数
// 正常情况下我们不会在代码中添加一个非法的正则表达式，因此对于错误情况直接使用panic进行处理
func MustCompile(str string) *Regexp {
	regexp, err := Compile(str)
	if err != nil {
		panic(`regexp: Compile(`+ quote(str) +`): ` + err.Error())
	}
	return regexp
}
```

3. recover函数

程序可以使用recover函数从panic状态恢复，通过在延迟调用函数中调用recover函数，可以消除当前goroutine中的一个panic，回到正常状态；recover函数旨在defer调用的函数中有效

在一个处于panic状态的goroutine退出之前，其中的panic不会蔓延到其他的goroutine；如果一个goroutine在panic状态下退出，则会使整个程序崩溃

```go
func funA() {
	fmt.Println("funA")
}

func funB() {
	// 延迟调用一个自执行的函数，必须得是panic了之后再恢复，因此recover必须配合defer
	defer func() {
		err := recover()
		// 如果程序出现了panic，那么可以通过recover恢复
		if err != nil {
			fmt.Printf("recover from panic: %v\n", err)
		}
	}()

	panic("funB panic")
}

func funC() {
	fmt.Println("funC")
}

func main() {
	funA()
	funB()
	funC()
	/*
		funA
		recover from panic: funB panic
		funC
	*/
}
```

使用recover函数从panic中恢复时需要特别注意以下几点：

①一定要在可能引发panic的语句之前调用recover函数

②recover函数必须搭配defer使用，并且不能直接使用defer延迟调用

③recover函数只能恢复当前goroutine的panic

除非明确知道程序应该使用panic和recover，否则一般不要使用它们

4. new和make函数

当我们声明一个值类型的变量时，Go语言会默认已经为其分配内存，而当我们声明指针或包含指针类型的变量时，不仅要声明类型，还要为它申请内存空间，否则没办法直接存储值，此时就要用到new和make函数

new函数不太常用，其分配内存后得到的是一个类型的指针，该指针对应的值就是该类型的零值

```go
func main() {
	a := new(int)		// 得到的是int类型的指针
	b := new(bool)		// 得到的是bool类型的指针
}
```

make函数更常用，但是他只用于切片、映射和通道三个类型，并且返回值也是这三个类型本身，而不是他们的指针类型

在使用切片、映射、通道时，都需要先使用make函数对他们进行初始化

```go
func main() {
	var s []int
	s = make([]int, 1) 			// make函数的返回值就是一个切片，分配一个内存
	s[0] = 100
	s = append(s, 200)			// 不会报错，切片会自动扩容
	fmt.Println(s[0], s[1])

	var m map[string]int
	m = make(map[string]int, 1) // make函数的返回值就是一个映射，分配两个内存
	m["zhangsan"] = 100

	m["zhangsi"] = 100 			// 不会报错，映射会自动扩容

	for k, v := range m {
		fmt.Println(k, v)
	}
}
```


## 指针

Go的指针不能偏移和运算，是安全指针

```go
func main() {
	// 使用取地址方式创建指针
	v := 10
	p := &v

	// 使用new创建指针
	p1 := new(int)

	// 指针取值
	*p1 = 10
}
```

## 结构体

1. 类型定义与类型声明

- 类型声明/定义

可以使用type关键字基于已存在的类型，定义新的类型，新类型与源类型的底层类型相同，支持显示转换

```go
func mian() {
	type MyInt int 			// 基于int类型声明一个MyInt类型
	var a int = 255
	var b MyInt
	// b = a 				// 不能将int类型的a直接赋值给MyInt类型的b
	b = MyInt(a)			// 支持显示类型转换
}
```

- 类型别名

类型别名比类型声明多一个等于号，rune类型和byte类型就属于类型别名，其中rune是int32的别名，完全等同于int32，byte是uint8的别名，完全等同于uint8；

类型别名和源类型本质上属于同一个类型，变量之间支持直接赋值，无需进行类型转换，类型别名只会存在于源代码中，在编译时会被自动替换成原来的类型

```go
func mian() {
	type MyInt = int
	var a = 255
	var b MyInt = 255
	b = a 		// int类型的a可以直接赋值给MyInt类型的b
}
```

- 二者区别

相当于类型声明是一个新的类型，而类型别名只是一个别名，本质上还是源类型

2. 结构体

当一个结构体类型中的所有字段都可以进行比较时，这个结构体类型就是支持比较的

若两个结构体中的字段相同但顺序不同，也属于不同的类型

以大写字母开头的字段名表示该字段是公开的，否则是私有的

字段的类型可以是结构体或结构体指针，一个结构体不能包含自己，但能包含自己的指针类型，一般用来定义一些递归的数据结构，如链表或树

```go
func main() {
	type Info struct {
		Email string
		Phone string
	}
	// tree 树结点
	type tree struct {
		value int
		left, right *tree	// 同一类型的可以放一行，只能包含自己的指针类型
		Contact Info		// 可以包含其他结构体
	}
}
```

声明结构体字段时，可以为其指定一个标签，该标签会在程序运行的时候通过反射机制被读取，标签内容由一个或多个键值对组成，用反引号括起来

字段标签在库中应用较多，比如数据库ORM、校验库validator、encoding/json库用来将JSON数据与GO对象相互转换，encoding/json在运行时会利用反射冲相关结构体变量的所有字段中根据json键去查找对应的值，然后使用这个值继续后续操作

```go
type Student1 struct {
	ID   int
	Name string
}

type Student2 struct {
	ID   int    `json:"stu_id"`
	Name string `json:"stu_name"`
}

func mian() {
	s1 := Student1{
		ID:   1,
		Name: "zhangsan",
	}

	s2 := Student2{
		ID:   2,
		Name: "zhangssi",
	}
	// 没有设置标签的，在序列化时默认使用结构体字段名称
	b1, err := json.Marshal(s1)
	fmt.Println(b1) // json:{"ID":1, "Name":"zhangsan"}
	// 设置标签的，在序列化时使用标签指定的字段名称
	b2, err := json.Marshal(s2)
	fmt.Println(b2) // json:{"ID":2, "stu_name":"zhangssi"}
}
```

3. 结构体变量和字面量

`var s Student`这样声明的结构体变量是零值，即每个字段均为对应类型的零值

结构体支持直接使用结构体字面量创建变量，此外还能按结构体字段顺序指定初始值

```go
s := {
	ID: 1,
	Name: "zhangsan",
	Gender: "男",
	Age: 18,
}

s := {
	ID: 1,
	Name: "zhangsan",
	Gender: "男",
	Age: 18}		// 这样则可以省略最后一个逗号

s := Student{1, "zhangsan", "男", 18}
```

4. 匿名结构体

在需要临时定义某些数据结构的场景中，可以使用匿名结构体，跟匿名函数一个道理，不需要调用，可以直接在运行过程中运行

```go
func main() {
	tmp := struct {
		ID int
		Info string
	}{
		ID: 1,
		Info: "zhangsan",
	}
	fmt.Println("tmp: ", tmp)
}
```

5. 结构体内存布局

一个结构体变量占用一块连续的内存，具体大小由结构体的字段决定，不是简单的将字段变量大小累加，因为CPU访问内存时，并不是逐个字节访问，而是以字长为单位访问（32位系统字长4字节，64位字长8字节），此外为了平台间的一致性和CPU内存访问的效率，编译器会对内存进行对齐操作，会进行操作系统对其和具体类型对齐；

> 对齐可以减少CPU访问内存的次数，加大CPU访问内存的吞吐量

对于操作系统来说，x86平台的对齐要求是4字节，x86_64平台的对齐要求是8字节

对齐可以参考[Go struct 内存对齐](https://geektutu.com/post/hpg-struct-alignment.html)讲的很详细，根据对齐倍数调整结构体的内存布局，会在结构体的相邻字段之间填充一些字节

> 结构体中字段的顺序会对对齐产生影响，因此要调整结构体内部字段变量的顺序，合理利用字段间的“填充”空间， 使得结构体的字段更加“紧凑”，从而缩小结构体的体积

当一个结构体的最后一个字段是空结构体（内存占用为0）时，编译器会自动在最后额外填充一些该结构体内存对齐要求的字节，可以放置对结构体的最后一个零内存占用字段进行取地址操作时发生越界，从而指向不相关的变量导致内存泄露，当空结构体出现在结构体的其他位置时，就不存在内存越界的风险，所以编译器不会对其进行额外的填充

6. 结构体指针

当结构体字段较多时，结构体体积变大，函数传参时会发生值拷贝，会导致开销较大，此外若需要在函数中修改结构体，此时就需要使用结构体指针，此外go支持直接对结构体指针使用`.`获取其字段，相当于go的语法糖

```go
func main() {
	// 直接使用&符号得到结构体指针
	s1 := &Student{
		ID:   1,
		Name: "zhangsan",
		Age:  18,
	}
	// 使用new函数直接创建结构体指针变量并使用
	var s1 = new(Student)
	s1.ID = 1				// 在底层相当于(*s1).ID
	s1.Name = "zhangsan"
	s1.Age = 18
}
```

7. 构造函数与结构体嵌套

构造函数相当于定义了一个函数，传入结构体字段的初始化值，由函数来实现对结构体的初始化，返回一个结构体指针（放置调用函数时的值拷贝性能开销过大）

```go
func NewStudent(id int, name string, age int) *Student {
	return &Student{
		ID:   id,
		Name: name,
		Age:  age,
	}
}
```

一个结构体中的字段可以是另外一个结构体类型，Go中也允许结构体中定义没有名只有类型的字段，这样的字段称为匿名字段，匿名字段只是把类型名当作字段名，但是在访问字段时却带来了方便，但是当结构体中的字段名与被嵌套的结构体中的字段名重复时，访问时被嵌套的结构体名称不可省略，不然就分不清访问的是哪个结构体的字段了

```go
func main() {
	t1 :=Teacher1 {
		Name: "zhangsan",
		Gender: "male",
		Info: Info {
			Email: "zhangsan@qq.com",
			Phone: "123456",
		},								// 逗号不能漏
	}

	var t2 Teacher2
	t2.Name = "lisi"
	t2.Gender = "female"
	t2.Email = "lisi@qq.com"			// 相当于t2.Info.Email，匿名字段带来访问字段的方便
	t2.Phone = "654321"
}
```

## 方法和接收者

Go语言允许为特定类型的变量设置专用的函数，即方法，不适用传统类的概念，而是使用方法接收者概念来绑定方法和对象，接收者的概念类似于其他语言中的this或self

1. 方法声明

方法声明与函数声明类似，但是需要在方法名前指定一个接收者参数

```go
func (接收者 接收者类型) 方法名(参数列表) (返回值列表) {
	方法体
}
```

接收者的变量名官方建议使用接收者类型的首字母小写，如Student类型的接收者变量应该命名为s

接收者类型必须是定义的类型（类型定义、结构体）或其他指针类型，最常用的是为结构体定义方法，对于数组、切片、映射这些类型，可以使用类型定义的方式声明一个新的类型，再为其声明方法；对于函数也可以基于函数定义新的类型，再为其添加方法；Go不支持为其他包的类型声明方法，也不支持为接口类型和基于指针定义的类型定义方法

```go
type Employee struct {
	name string
	salary int 
}

func (e Employee) SayHi() {			// 使用接收者类型首字母小写作为接收者变量名
	fmt.Println("Hi, my name is", e.name, "and my salary is", e.salary)
}

// 对于基本类型，数组、切片、映射、函数，通过类型定义的方式可以为其添加方法
type MyInt int 
func (m MyInt) SayHi() {
	fmt.Println("Hi, my name is", m)
}

type StringSet map[string]bool
func (ss StringSet) Has(key string) bool {			// Has判断集合中是否包含指定key的方法
	return ss[key]
}

type FilterFunc func(int)
func (ff FilterFunc) Do(n int) {					// Do执行函数的方法
	ff(n)
}
```

2. 方法调用

与访问结构体变量的字段类似，使用`.`进行调用

3. 值接收者与指针接收者

与函数传值传指针类似，涉及到在方法中修改结构体字段的值，或传递接收者值拷贝开销过大的时候，使用指针接收者会更好，此外还需要保证一致，如果一个类型的某个方法使用了指针类型接收者，那么该类型的其他方法也应该使用指针类型接收者

```go
type Employee struct {
	name string
	salary int 
}

func (e Employee) raise1(n int) {
	e.salary += n
}

func (e *Employee) raise2(n int) {
	e.salary += n
}

func main() {
	// Go语法糖
	// 对于声明时使用值接收者的方法，使用指针也可以直接调用
	e1 := &Employee{"zhangsan", 1000}
	e1.raise1(100)						// 编译器会自行根据指针e1获取实际值，即等价于(*e1).raise1(100)

	// 对于声明时使用指针接收者的方法，使用值也可以直接调用
	e2 := Employee{"lisi", 2000}
	e2.raise2(100)						// 编译器会自行根据值e2获取实际值，即等价于(&e2).raise2(100)
}
```

Go语言中类型的方法集定义了一组关联到具体类型的值或指针的方法

定义方法时使用的接收者的类型决定了这个方法是关联到值还是关联到指针

4. 组合

Go中通过结构体嵌套的方式进行组合，从而实现类似其他语言中面向对象的继承

在Teacher结构体中以匿名字段方式嵌入Info结构体，从此Teacher结构体就拥有了Info结构体的所有字段和方法，并支持直接访问

Go中进行方法调用时，会先查找当前结构体中是否声明了该方法，没有则依次从当前结构体的嵌入字段的方法中查找，若由重名，则要通过被嵌入字段调用该重名方法，不能省略

```go
type Info struct {
	Email string
	Phone string
}

func (i Info) Detail() string {
	return fmt.Sprintf("Email: %s, Phone: %s", i.Email, i.Phone)
}

type Teacher struct {
	Name string
	Gender string
	Info
}

func main() {
	t := Teacher{
		Name: "zhangsan",
		Gender: "male",
		Info: Info{
			Email: "zhangsan@qq.com",
			Phone: "123456",
		},
	}

	fmt.Println(t.Detail())			// 直接访问嵌套的结构体的方法
}
```

结构体还可以嵌入其他包内定义的结构体，下面展示通过在匿名结构体中嵌入sync.Mutex结构体进行代码优化

首先，借助映射和互斥锁sync.Mutex实现一个简易的缓存

```go
func main() {
	var (
		mux sync.Mutex						// 防止数据竞争的互斥锁
		mapping = make(map[string]string)	// 存储缓存数据的map
	)

	// Query 根据指定key查询缓存
	func Query(key string) string {
		mux.Lock()		// 加锁
		v := mapping[key]
		mux.Unlock()	// 解锁
		return v
	}
}
```

接下来通过在匿名结构体中嵌入sync.Mutex结构体，改写代码，使用更形象的cache变量，并且由于嵌入了sync.Mutex，可以直接调用Lock和Unlock两个方法，让逻辑更加直观和清晰，能清楚知道是对谁加锁解锁

```go
func main() {
	// cache 定义一个表示缓存的匿名结构体
	var cache = struct {
		sync.Mutex 		// 嵌入
		mapping map[string]string
	} {
		mapping: make(map[string]string),
	}

	// Query根据指定key查询缓存
	func Query(key string) string {
		cache.Lock()
		v := cache.mapping[key]
		cache.Unlock()
		return v
	}
}
```

5. 结构体和JSON序列化

JSON序列化：结构体-->JSON格式的字符串

JSON反序列化：JSON格式的字符串-->结构体

结构体的标签是结构体的元信息，可以在运行的时候通过反射的机制读取出来，标签由一个或多个键值对组成，键与值使用冒号分隔，值使用双引号括起来，同一个结构体字段可以设置多个键值对，不同键值对之间使用空格分隔

```go
type Student struct {
	ID   int 	`json:"stu_id"`		// 设置JSON标签，JSON序列化时将stu_id作为键名
	Gender string					// JSON序列化时默认将字段名作为key
	name string 					// 私有字段，不能被JSON包访问
}
```

## 接口

### 接口类型

从Go1.18开始，Go中接口分为一般接口和基本接口

如果一个接口类型的定义包含非接口类型，则称为一般接口，目前一般接口只能用于类型约束，不能实例化变量，参考泛型章节内容

如果一个接口类型的定义中只包含方法，则称为基本接口，是一组方法的集合，规定了需要实现的所有方法

```go
type 接口类型名 interface {
	方法名1(参数列表1) 返回值列表1
	方法名2(参数列表2) 返回值列表2
	...
}
```

一般会在单词后加er作为接口类型名，如Writer、Closer；当方法名的首字母大写且接口类型名的首字母页大写时，该方法可以被其他包访问；参数列表和返回值列表中的参数变量名可以省略

1. 实现接口

当满足以下条件时，称类型T实现了接口I：

- T不是接口：类型T是接口I代表的类型集中的一个成员
- T是接口：T接口代表的类型集是I代表的类型集的子集

一个接口类型只要实现了接口中规定的所有方法，就被称为实现了这个接口

2. 为什么要使用接口

假设猫和狗饿了会叫，此时又来了一只羊，饿了也会叫，随着动物的增多，他们都饿了会叫，可以将Sayer定义为一个接口，只要饿了就调用Say()方法，这样就把所有会叫的动物都当作Sayer类型来处理，只要实现了Say()方法就能被当成Sayer类型的变量来处理

通过接口，可以让我们专注于该类型提供的方法，而不是类型本身，从而写出更加通用和灵活的代码

```go
// 定义接口
type Sayer interface {
	Say()
}

type Cat struct {}
func (c Cat) Say() {		// 实现了Sayer接口
	fmt.Println("meow")
}

type Dog struct {}
func (d Dog) Say() {		// 实现了Sayer接口
	fmt.Println("woof")
}

type Sheep struct {}
func (s Sheep) Say() {		// 实现了Sayer接口
	fmt.Println("hello")
}

// 可以定义一个通用的饿肚子场景，将接口类型变量传入，这样就不用为猫、狗、羊分别定义一个饿肚子场景
func MakeHungry(s Sayer) {
	s.Say()
}

func main() {
	// 只要实现了Say()方法，就能被当成Sayer类型的变量来处理
	var c Cat
	MakeHungry(c)
	var d Dog
	MakeHungry(d)
	var s Sheep
	MakeHungry(s)
}
```

2. 面向接口编程

Go不像PHP、Java需要显式声明一个类实现了哪些接口，Go使用隐式声明的方式实现接口，即一个接口类型只要实现了接口中规定的所有方法，就实现了这个接口

对于接口的理解，相当于不同的结构体有一个相同的方法，通过将该方法定义成一个接口类型要实现的方法，从而不用为每个结构体定义一个调用该方法的方法，只需要通过一个方法 ，传入接口类型参数，就可以实现调用不同结构体的相同方法

```go
type ZhiFuBao struct {
	// 支付宝
}

// Pay 支付宝的支付方法
func (z *ZhiFuBao) Pay(amount int64) {
	fmt.Printf("使用支付宝付款：%.2f 元。\n", float64(amount)/100)
}

/* // Checkout 结账
func CheckoutZFB(obj *ZhiFuBao) {
	// 支付100元
	obj.Pay(100)
} */

// 随着业务拓展，支持用户使用微信支付
type WeChat struct {
	// 微信
}

// Pay 微信的支付方法
func (w *WeChat) Pay(amount int64) {
	fmt.Printf("使用微信付款：%.2f 元。\n", float64(amount)/100)
}

/* // Checkout 结账
func CheckoutWX(obj *WeChat) {
	// 支付100元
	obj.Pay(100)
} */

// 通过引入接口类型，能简化上面为不同的结构体定义不同的结账方法
type Payer interface {
	Pay(int64)
}

// Checkout 结账
func Checkout(obj Payer) {
	// 支付100元
	obj.Pay(100)
}

func main() {
	Checkout(&ZhiFuBao{})		// 使用支付宝支付
	Checkout(&WeChat{})			// 使用微信支付
}
```

### 值接收者和指针接收者

Go语言有对指针求值的语法糖，因此对于使用值接收者实现接口，不管是结构体类型还是对应的结构体指针类型（语法糖对指针求值），都可以将其变量赋值给该接口变量；

而由于我们并不是总能对一个值求址（可能没有初始化没有分配地址），因此使用指针接收者实现接口时，只能将结构体类型变量赋值给该接口变量，而不能将结构体指针类型变量赋值给该接口变量

```go
type Mover interface {
	Move()
}

type Dog struct {}
// 使用值接收者实现接口
func (d Dog) Move() {
	println("dog move")
}

type Cat struct {}
// 使用指针接收者实现接口
func (c *Cat) Move() {
	println("cat move")
}

func main() {
	var x Mover
	var d1 = Dog{}	// d1是Dog类型
	x = d1 			// 可以将d1赋值为变量x
	x.Move()
	var d2 = &Dog{}	// d2是Dog指针类型
	x = d2			// 也可以将d2赋值为变量x
	x.Move()

	var c1 = &Cat{}	// c1是Cat指针类型
	x = c1			// 也可以将c1赋值为变量x
	x.Move()
	var c2 = Cat{}	// c2是Cat类型
	x = c2			// 报错，不可以将c2赋值为变量x，即不可以将c2当成Mover类型
	x.Move()
}
```

### 类型与接口的关系

1. 一个类型实现多个接口

接口间彼此独立，互相不知道对方的实现；同一类型实现不同的接口互相不影响

```go
type Sayer interface {
	Say()
}

type Mover interface {
	Move()
}

type Dog struct{}
// 狗实现了Sayer和Mover接口
func (d Dog) Say() {
	fmt.Println("woof")
}
func (d Dog) Move() {
	fmt.Println("walk")
}

func main() {
	d := Dog{}
	d.Say()
	d.Move()
}
```

2. 多个类型实现同一接口

比如可以把狗和骑车都当成会动的类型来处理，不必关注他们具体是什么，只需要调用他们的Move方法即可

并且一个接口的所有方法可以通过在类型中嵌入其他类型或者结构体来实现

```go
type WashingMachine interface {
	wash()
	dry()
}

// 甩干机
type dryer struct{}
func (d dryer) dry() {
	fmt.Println("甩干")
}

// 海尔洗衣机
type haier struct{
	dryer		// 嵌入了甩干机
}
func (h haier) wash() {
	fmt.Println("洗衣")
}
// 因此海尔洗衣机只需要实现wash方法即可实现了WashingMachine接口
```

3. 接口组合

接口和接口可以组合形成新的接口，在新的接口中还可以额外定义新的方法，类似结构体嵌套那样

一个类型要实现组合接口，就要实现组合接口中所有接口的所有方法

使用接口组合的优势在于类型可以使用组合接口中的所有方法

```go
type Retriver interface {
	Get(url string) string
}

type Poster interface {
	Post(url string, form map[string]string) string
}

// 组合接口
type RetrieverPoster interface {
	Retriver
	Poster
	// 当然这里还可以添加别的方法
}

// 假设有个session函数，它需要一个参数，既是一个Retriver，又是一个Poster，此时就可以用到接口的组合，可以同时调用组合接口中的所有方法
func session(s RetrieverPoster) string {
	s.Post("http://www.baidu.com", map[string]string{"key": "value"})
	return s.Get("http://www.baidu.com")
}
```

结构体嵌入接口，该结构体将自动获得接口的所有方法，因此该结构体就实现了该接口

对于结构体的这个接口字段，任何实现了该接口的类型都可以赋值给该字段

```go
type Fooer interface {
	Foo() string
}

// Container结构体嵌入了Fooer接口，因此Container结构体能获得Fooer接口的所有方法，因此Container结构体实现了Fooer接口
type Container struct {
	Fooer
}

// 可以理解为Container结构体有这样一个转发的方法，其中cont.Fooer指的是任何实现了Fooer接口的对象
/* func (cont Container) Foo() string {
	return cont.Fooer.Foo()
} */

func sink(f Fooer) {
	fmt.Println("sink:", f.Foo())
}

// TheRealFoo结构体实现了Fooer接口
type TheRealFoo struct {}
func (trf TheRealFoo) Foo() string {
	return "TheRealFoo Foo"
}

// Container结构体的Fooer字段类型是任何实现了Fooer接口的对象
co := Container{Fooer: TheRealFoo{}}
sink(co)// Container结构体实现了Fooer接口，因此可以传入sink函数中
```

通过在结构体中嵌入接口，可以使该结构体直接实现该接口，但是一个需要注意的地方是，必须保证结构体中的接口属性不为nil，否则调用结构体中嵌入的接口的方法就会空指针panic，比如Go排序sort源码中reverse结构体就是不可导出的，只能使用sort.go中的Reverse函数让使用者创建reverse结构体示例，这样能保证结构体嵌入的接口字段不为空

```go
// src/sort/sort.go
// Interface定义通过索引对元素排序的接口类型
type Interface interface {
	Len() int 
	Less(i, j int) bool 
	Swap(i, j int) 
}

// reverse 结构体嵌入了Interface接口
type reverse struct {
	Interface
}

// sort.go通过定义一个可导出的Reverse函数让使用者创建reverse结构体实例，从而保证得到的reverse结构体中的Interface属性一定不为nil
func Reverse(data Interface) Interface {
	return &reverse{data}
}
```

此外，在结构体中嵌入接口，在访问接口方法的时候有区别，可以改写该接口类型的方法

```go
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
	w.Salary.func1()			// Salary func1
	w.Salary.func2()			// Salary func2
}
```

### 空接口

空接口指没有定义任何方法的接口类型，因此任何类型都可以认为实现了空接口，因此空接口类型变量可以存储任意类型的值

使用空接口可以实现可以接收任意类型的函数参数，使用空接口可以实现可以保存任意值的字典

```go
// 空接口作为函数参数，可以接收任意类型的函数参数
func shouw(a interface{}) {
	fmt.Println("type: %T value: %v\n", a, a)
}

func main() {
	// 使用空接口类型通常不需要使用type关键字声明
	var x interface{}		
	// 使用空接口实现可以保存任意值的字典
	var studentInfo = make(map[string]interface{})
	studentInfo["name"] = "zhangsan"
	studentInfo["age"] = 18
	studentInfo["married"] = false
}
```

### 接口值

接口值由两部分组成，一个是具体类型，一个是具体类型的值，称为接口的动态类型和动态值

接口值当且仅当动态类型和动态值两个豆相等时才相等

并且当接口值保存的动态类型相同时，如果动态类型不支持比较（例如切片），那么对他们进行比较会引发panic

```go
type Mover interface {
	Move()
}
type Dog struct {
	Name string
}
func (d *Dog) Move() {
	fmt.Println("dog move")
}
type Car struct {
	Brand string 
}
func (c *Car) Move() {
	fmt.Println("car move")
}

func main() {
	// 创建一个Mover接口类型的变量m，其类型和值都是nil
	var m Mover
	fmt.Println(m == nil) 		// true
	// 不能对一个空接口值调用任何方法，否则会产生panic

	// 将*Dog结构体指针赋给变量m，此时接口值m的动态类型会被设置为*Dog，动态值是结构体变量的拷贝
	m = &Dog{Name: "dog"}

	// 为接口变量m赋一个*Car类型的值，此时接口值m的动态类型会被设置为*Car，动态值为nil
	m = new(Car)
	fmt.Println(m == nil)		// false，因为只有动态值部分为nil，动态类型部分保存着对应值的类型
}
```

1. 类型断言

fmt包内部使用反射机制在程序运行时获取动态类型的名称

从接口值中获取对应的实际值需要使用类型断言

当一个接口值有多个实际类型要判断时，可以使用switch语句

```go
// 当一个接口值有多个实际类型需要判断时，可以使用switch
func justifyType(x interface{}) {
	switch v := x.(type) {
	case string:
		fmt.Println("x is string, value is ", v)
	case int:
		fmt.Println("x is int, value is ", v)
	case bool:
		fmt.Println("x is bool, value is ", v)
	default: 
		fmt.Println("x is unknown type")
	}
}

func main() {
	var n Mover = &Dog{Name: "dog"}
	// 类型断言语法：第一个参数是n转化为*Dog类型后的变量，第二个参数是一个布尔值，表示断言是否成功
	v, ok := n.(*Dog)
	if ok {
		fmt.Println("类型断言成功")
		v.Name = "富贵" 	// 变量v的类型是*Dog
	} else {
		fmt.Println("类型断言失败")
	}
}
```

## 反射

### 反射简介

反射是指计算机程序在运行时可以访问、检测和修改自身状态或行为的能力，Go通过反射机制，在程序运行期间获取变量的名称、类型等信息，用于动态语言等特性

编译阶段，变量会被转换为内存地址，但是变量名不会被编译器写入可执行部分，因此程序在运行阶段通常无法获取自身信息

Go通过标准库reflect提供了反射机制，从而能够在编译期将变量的反射信息，如字段名称、类型信息、结构体信息等整合到可执行文件中，并为程序提供访问反射信息的接口，从而程序就可以在运行阶段获取类型的反射信息并修改

反射就是通过接口值中保存的类型信息实现的，其本质是在运行时动态获取一个变量的类型信息和值信息

反射最常见的应用场景就是依赖用户输入的内容执行后续程序

- 结构体属性的tag支持用户动态定义属性名

- 根据用户输入的内容决定生成什么类型的对象或者调用什么方法

### reflect包

`reflect.TypeOf()`函数签名为`func TypeOf(i interface{}) Type`，用于获取任意值的类型对象`reflect.Type`接口类型

`reflect.Type`接口类型有很多方法，其中`reflect.Type.Name`方法用于获取声明的类型，`reflect.Type.Kind`方法用于获取底层类型

`reflect.ValueOf`函数签名为`func ValueOf(i interface{}) Value`，用于获取原始值的值信息`reflect.Value`结构体类型

`reflect.Type`结构体类型有很多方法，其中`func (v Value) IsNil() bool`用于检查v持有的值是否为nil；`func (v Value) IsValid() bool`用于检查v是否持有一个值

任何接口值都是由动态类型和动态值两部分组成，在反射中可以理解为由reflect.Type和reflect.Value两部分组成，reflect包提供reflect.TypeOf和reflect.ValueOf两个函数来获取任意值的Type和Value信息

1. TypeOf函数

可用于获取任意值的类型对象，从而可以访问任意值的类型信息

其函数签名为：`func TypeOf(i interface{}) Type`，入参为空接口类型，因此实参都会被转为空接口类型，从而函数得到的接口值就保存了原始值的类型信息和值信息，返回类型是reflect.Type，是一个接口类型

```go
func main() {
	var a float32 = 1.23
	ta := reflect.TypeOf(a)
	fmt.Printf("ta: %v\n", ta)		// type: float32

	var b float64 = 1.23
	tb := reflect.TypeOf(b)
	fmt.Printf("tb: %v\n", tb)		// type: float64
}
```

2. Type和King

Go语言中每个变量都有一个在编译阶段就确定的静态类型，通过反射得到的类型信息分为Type和Kind，前者指声明的类型，后者指语言底层的类型，可以通过调用reflect.Type的Name方法得到Type名称，调用Kind方法得到kind名称

```go
type myInt int64
func reflectType(x interface{}) {
	t := reflect.TypeOf(x)
	fmt.Printf("type: %v kind: %v \n", t.Name(), t.Kind())		//通过Name方法和Kind方法获取声明的类型和底层的类型
}


func main() {
	var a *float32   	// 指针
	var b myInt			// 自定义类型
	var c rune			// 类型别名
	reflectType(a)		// type: kind: ptr
	reflectType(b)		// type: myInt kind: int64
	reflectType(c)		// type: int32 kind: int32

	type person struct {
		name string
	}
	type book struct {
		title string
	}
	var d = person{
		name: "zhangsan",
	}
	var e = book{
		title: "book1",
	}

	reflectType(d)		// type: person kind: struct
	reflectType(e)		// type: book kind: struct
}
```

在Go的反射中，数组、切片、Map、指针等类型变量的Type名称都是空字符串，只有一个底层的Kind名称

Kind包含的底层类型（截止Go1.20版本）：`Invalid Kind = iota非法类型、Bool、Int、Int8、Int16、Int32、Int64、Uint、Uint8、Uint16、Uint32、Uint64、Uintptr指针、Float32、Float64、Complex64、Complex128、Array、Chan通道、Func、Interface、Map、Pointer指针、Slice、String、Struct、UnsafePointer底层指针`

3. reflect.ValueOf函数

返回一个reflect.Value结构体类型，包含了原始值的值信息，reflect.Value与原始值之间可以相互转换

```go
func reflectValue(x interface{}) {
	v := reflect.ValueOf(x)
	k := v.Kind()
	switch k {
	case reflect.Int64:
		// v.Int() 从反射中获取整型的原始值
		fmt.Printf("type is int64, value is %d\n", v.Int())
	case reflect.Float32:
		// v.Float() 从反射中获取浮点型的原始值，需要将类型转换为float32
		fmt.Printf("type is float32, value is %f\n", float32(v.Float()))
	case reflect.Float64:
		// v.Float() 从反射中获取浮点型的原始值，然后将类型转换为float64
		fmt.Printf("type is float64, value is %f\n", float64(v.Float()))
	}
}

func main() {
	var a float32 = 3.14
	var b int64 = 100
	reflectValue(a)		// type is float32, value is 3.140000
	reflectValue(b)		// type is int64, value is 100
	
	// 将int类型的原始值转换为reflect.Value类型
	c := reflect.ValueOf(10)
	fmt.Printf("c: %v type: %T\n", c, c)	// c: 10 type: reflect.Value
}
```

若向在函数中通过反射修改实参的值，则在函数传参的时候必须传入指针类型，否则得到的只是副本拷贝；指针类型可以通过反射传入，有专用的Elem方法来获取指针对应的值

```go
func reflectSetValue1(x interface{}) {
	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Int64 {
		fmt.Println("x is int64")
		v.SetInt(1000) // 修改的是副本，panic
	}
}

func reflectSetValue2(x interface{}) {
	v := reflect.ValueOf(x)
	// 反射中使用Elem()方法可以获取指针对应的值
	if v.Elem().Kind() == reflect.Int64 {
		fmt.Println("x is int64")
		v.Elem().SetInt(1000)
	}
}

func main() {
	var a int64 = 10
	reflectSetValue2(&a)
}
```

### reflect.Value结构体

定义如下

```go
type Value struct {
	// 类型信息
	typ *rtype
	// 真是数据的地址
	ptr unsafe.Pointer
	// 元信息标志位
	flag
}
```

该结构体有很多方法，通过这些方法可以直接操作其ptr字段所指向的实际数据，如`IsNil`用于判断Value是否为空指针；`IsValid`用于判断Value值是否有效

`func (v Value) IsNil() bool`返回v持有的值是否为nil，只能由通道、函数、接口、映射、指针或切片调用

`func (v Value) IsValid() bool`返回v是否持有一个值，若v是Value零值（此时调用除IsValid、String、Kind外的方法都会panic），则返回false，否则返回true

```go
func main() {
	// 空指针
	var a *int
	va := reflect.ValueOf(a)
	fmt.Println(va.IsNil())   // true
	fmt.Println(va.IsValid()) // true

	// nil值
	v := reflect.ValueOf(nil)
	fmt.Println(v.IsValid()) // false
	// fmt.Println(v.IsNil())   // false，不注释会panic

	// 匿名结构体变量
	b := struct{}{}
	// 尝试从结构体中查找"abc"字段
	vbf := reflect.ValueOf(b).FieldByName("abc")
	fmt.Println(vbf.IsValid()) // false
	// fmt.Println(vbf.IsNil())   // false，不注释会panic

	// 尝试从结构体中查找"abc"方法
	vbm := reflect.ValueOf(b).MethodByName("abc")
	fmt.Println(vbm.IsValid()) // false
	// fmt.Println(vbm.IsNil())   // false，不注释会panic
}
```

以上代码对于空的不注释会panic，但是书中说的是注释会panic，可能后续Go机制更改，策略变了，需要再调研一下

### 结构体反射

结构体反射是使用反射较多的场景，若反射对象的类型是结构体，则可以通过反射值对象（reflect.Type）的NumField（返回结构体成员字段数量）和Field（返回索引对应的结构体字段的信息）方法获取结构体成员的详细信息

reflect包中定义了一个StructField结构体类型来描述结构体中一个字段的相关信息

```go
type StructField struct{
	// Name是字段的名字，PkgPath是非导出字段的包路径，对导出字段该字段为""
	Name String
	PkgPath string
	Type Type		// 字段的类型
	Tag StructTag	// 字段的标签
	Offset uintptr  // 字段在结构体中的字节偏移量
	Index []int		// 用于Type.FieldByIndex时的索引切片
	Anonymous bool  // 是否为匿名字段
}
```

1. 结构体反射示例

需求分析：需要从一个文本文件info.txt中读取学生信息并赋值给结构体变量

```go
// info.txt内容
name=七米
age=18
```

示例一：结构体定义已经固定，但是文件内容的格式不确定，因此需要在程序运行的时候，读取文件内容，然后通过反射读取结构体字段的类型，再将内容转换称该字段对应的类型，最后进行赋值，从而实现读取文本文件最终对其进行解析，解析到对应结构体中

```go
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
// 思路是：一行一行读取文本文件的内容，并把每行内容都按照等号分割为键值对，然后按键名称去结构体变量中找到对应的结构体字段，根据字段的类型转换值的类型赋值给字段
func LoadInfo(s string, v interface{}) (err error) {
	// tInfo为v的类型对象，调用Kind方法获取底层类型
	tInfo := reflect.TypeOf(v)
	if tInfo.Kind() != reflect.Ptr {
		err = errors.New("Please pass into a struct ptr")
		return
	}
	// Elem()方法获取指针对应的值，再获取对应的底层类型
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
		// 去除字符串两端的空格
		key := strings.TrimSpace(kvList[0])
		value := strings.TrimSpace(kvList[1])
		// 遍历结构体字段的tag找到对应的key，NumField返回结构体成员字段数量
		for i := 0; i < tInfo.Elem().NumField(); i++ {
			// 返回索引对应的结构体字段的信息
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
		// 根据字段的类型将value转换为对应的类型赋给字段
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
```

示例二：根据传入的方法名在结构体变量的方法集中查找到对应的方法，然后调用该方法

```go
// Do调用变量v的name方法
// 实现思路：根据入参的方法名name在结构体变量的方法集中查找到对应的方法，然后再调用该方法
func Do(v interface{}, name string, arg interface{}) {
	tInfo := reflect.TypeOf(v)
	vInfo := reflect.ValueOf(v)

	// 返回该类型的方法集中方法的数量
	fmt.Println(tInfo.NumMethod())
	// 根据方法名返回该类型方法集中的方法
	m := vInfo.MethodByName(name)
	if !m.IsValid() || m.IsNil() {
		fmt.Printf("%s没有%s方法\n", tInfo.Name(), name)
		return
	}

	// 调用指定方法，通过反射调用方法传递的参数必须是[]reflect.Value类型的
	argVal := reflect.ValueOf(arg)
	// Call 方法是 reflect.Value 类型的一个方法，用于调用对应的方法
	// Call 方法需要一个 []reflect.Value 类型的切片作为参数，这个切片包含了要传递给方法的参数
	// arg就是Do函数传入的要传给name方法的参数
	m.Call([]reflect.Value{argVal})
}
```

关于Go中类型和值的关系，还是有点混乱，二轮学习的时候就要搞清楚这个

### 反射三大定律

- 反射是从接口值获取反射对象（reflect.Type和reflect.Value）的机制

- 反射也可以从反射对象得到接口值

- 如果要修改一个反射对象，那么它必须是可设置的

对于第三个定律，也即是，只有传入指针类型时得到的反射对象才是可设置的，在需要通过反射修改值的场景中，一定要为反射函数传入指针类型

```go
var x int64 = 10
v := reflect.ValueOf(x)
v.SetInx(100)			// panic，因为传递的是变量x的值拷贝，反射对象v把那个不能代表变量x

v := reflect.ValueOf(&x)
// 指针类型反射对象需要调用Elem方法得到原值
v.Elem().SetInt(100)
fmt.Println(v.Elem().Interface())	// 100
fmt.Println(x)						// 100
```

反射虽然灵活，但是不能滥用

- 基于反射的代码是机器脆弱的，反射中的类型错误只有在运行的时候才会引发panic

- 大量使用反射的代码通常难以理解

- 反射的性能较低，基于反射实现的代码通常比正常代码的运行速度慢一到两个数量级

## 并发编程

Go语言天生支持并发，能充分利用现代CPU的多核优势

### 并发编程简介

1. 串行、并行与并发

- 串行：多个任务顺序执行

- 并发：同一时间段内多个任务交替执行

- 并行：同一时间段内多个任务同时执行

2. 进程、线程和协程

- 进程：程序在操作系统中的一次执行过程，是系统进行资源分配和调度的独立单位

- 线程：操作系统基于进程开启的轻量级“进程”，是操作系统调度执行的最小单位

- 协程：用户态“线程”，比线程更轻量级

3. 并发模型

将并发编程的方法归纳为各种并发模型

- 线程与锁模型

- Actor模型

- CSP模型

- Fork与Join模型

而Go语言并发程序主要通过基于通信顺序过程的goroutine和通道channel实现，也支持传统的多线程共享内存的并发方式

### goroutine

goroutine是Go语言并发的核心，以一个很小的栈开始其生命周期，一般只需要2KB，并且不同于操作系统线程由系统内核调度，goroutine由Go运行时（runtime）调度

goroutine是Go程序中最基本的并发执行单元，每个Go程序都至少包含一个goroutine，也就是main goroutine

1. go关键字

在函数或方法调用前加上go关键字，该函数或方法就可以在新创建的goroutine中执行

```go
go f()		// 创建一个心的goroutine运行函数f
go func() {

}()			// 匿名函数也支持go关键字
```

2. 启动单个goroutine

main goroutine结束，所有由main goroutine创建的goroutine都会退出，因此需要让main goroutine等待，可以使用`time.Sleep(time.Second)`进行等待，但是不够优雅，因为创建goroutine执行函数需要一定的时间开销，因此使用`time.Sleep`进行等待也是不准确的

```go
func main() {
	// 串行执行
	hello()
	fmt.Println("你好")

	// 并发执行，只会给出你好，因为创建goroutine执行hello函数需要时间开销，然而main goroutine是继续执行的，main goroutine结束时所有由main goroutine创建的goroutine都会结束
	go hello()
	fmt.Println("你好")

	// 并发执行，使用time.Sleep等待一秒，你好在前hello在后，因为创建goroutine执行hello函数需要时间开销，然而main goroutine是继续执行的
	go hello()
	fmt.Println("你好")
	time.Sleep(time.Second)
}
```

3. 启动多个goroutine

若不关心并发操作的结果或者有其他方式收集并发操作的结果时，`WaitGroup`是实现等待一组并发操作完成的好方法，但是多个goroutine是并发执行的，而goroutine的调度是随机的，因此就会导致多个goroutine的执行顺序不确定

```go
// 声明全局等待组变量
var wg sync.WaitGroup

func hello(i int) {
	defer wg.Done()		// goroutine结束就-1
	fmt.Println("Hello World", i)
}

func main() {
	for i := 0; i < 10; i ++ {
		wg.Add(1)		// 启动一个goroutine就+1
		go hello(i)
	}
	wg.Wait()		// 等待所有goroutine结束
}
```

4. 动态栈

操作系统线程一般固定栈内存：2MB，goroutine的初始栈内存很小：2KB，同时goroutine的栈是不固定的，可以动态扩大或缩小，Go的运行时会自动为goroutine分配合适的栈空间

5. goroutine调度

操作系统的线程切换需要完整的上下文切换，开销较大，即操作系统调度时，会先挂起当前执行的线程，并将其寄存器内容保存到内存中，然后选出下一次要执行的线程，从内存中回复该线程的寄存器信息，最后恢复执行该线程的现场并开始执行线程

而goroutine的调度在Go运行时层面实现，完全由Go语言本身实现，其按照一定规则将所有goroutine调度到操作系统线程上执行，goroutine调度器使用GPM调度模型

![GMP调度模型](./images/GMP.jpg)

- G表示goroutine，包含要执行的函数和上下文信息

- P本地队列和全局队列都用来存放等待运行的G，但P本地队列数量不超过256个，新建的G优先加入P本地队列，满了才批量移动部分G到全局队列

- P表示goroutine执行所需要的资源，最多有GOMAXPROCS个

- M表示内核线程，会循环获取P和G，优先从P本地队列获取G，空了再向其他P本地队列或全局队列获取G

- goroutine调度器和操作系统调度器通过M结合起来，每个M都代表1个内核线程，操作系统调度器负责把内核线程分配到CPU的核上执行

Go语言的优势在于操作系统线程是由操作系统内核调度的，goroutine则是由Go运行时调度的，完全在用户态下完成，不涉及内核与用户态之间的频繁切换；内存的分配与释放也是在用户态维护一块大的内存池，不直接调用系统的malloc函数（除非内存池需要改变），成本比调度操作系统线程低很多；Go还充分利用多核硬件资源，近似的把若干goroutine均分在物理线程上，再加上自身goroutine的超轻量级，保证了goroutine调度的性能

GOMAXPROCS参数用于确定需要使用多少个操作系统线程来同时执行代码，默认是CPU核心数，可以通过`runtime.GOMAXPROCS`设置当前程序并发时占用的CPU逻辑核心数，Go1.5之前默认使用单核执行，从Go1.5之后默认使用全部CPU逻辑核心数

### 通道

通道用于并发函数之间进行数据交换，即让一个goroutine发送特定值到另一个goroutine的通信机制

共享内存进行数据交换时，在不同goroutine中容易发生竞态问题；而为了保证数据交换的正确性，很多并发模型都必须使用互斥量对内存加锁，造成性能问题

而Go语言采用CSP并发模型（通信顺序过程），提倡通过通信实现共享内存，而不是通过共享内存实现通信

通道遵循先进先出原则，保证收发数据的顺序，每个通道都是一个具有类型的导管，默认零值为nil，需要使用make函数初始化后才能使用，缓冲区大小可选

```go
var ch1 chan int		// 声明一个传递整型的通道
make(chan int, 10)	// 声明一个缓冲区大小为10的整型通道
```

1. 通道操作

发送、接收与关闭三种操作，通常由发送方执行关闭操作，并且只有在接收方明确等待通道关闭的信息时才需要执行关闭操作，即不同于关闭文件，关闭通道不是必须的

```go
ch := make(chan int)		// 定义一个int类型通道
ch <- 10					// 发送10到通道	

v := <- ch				// 从通道接收数据，并赋值给变量v
v, ok := <-ch				// 从通道接收数据，并赋值给变量v，ok为true表示通道未关闭，false表示通道已关闭

close(ch)				// 关闭通道
```

关闭后的通道有以下特点：

- 再发送值会导致panic

- 接收会一直获取值，直到通道为空

- 若通道已经没有值，此时接收会得到对应类型的零值

- 再次关闭会导致panic

2. 无缓冲的通道

又称为阻塞的通道，无缓冲通道只有在接收方能够接收值时才能发送成功，否则会一直等待发送，同理，对一个无缓冲通道执行接收操作时，若没有任何向通道中发送值的操作，也会导致接收操作阻塞，即没有缓冲，数据实时收发，会导致发送和接收的goroutine同步化，也成为同步通道

```go
func recv(c chan int) {
	ret := <- c
	fmt.Println("接收成功", ret)
}

func main() {
	ch := make(chan int)		// 无缓冲通道
	go recv(ch)					// 创建一个goroutine从通道接受值，会阻塞直到有值可以接收才会继续执行，即阻塞等待接收数据
	ch <- 10					// 发送10到通道，若没有正在等待接收的goroutine，也会阻塞等待有goroutine接收数据
	// 发送10之后，main goroutine和recv goroutine都会继续执行，从而同步
	fmt.Println("发送成功")
}
```

3. 有缓冲的通道

通道的容量表示该通道中能存放的最大元素数量，通道满时，就会像无缓冲通道那样，发送会阻塞，需要等待接收，可以使用len函数获取通道内元素的数量，使用cap函数获取通道的容量，但是一般不这么做

```go
func main() {
	ch := make(chan int, 1) 		// 创建一个容量为1的有缓冲通道
	ch <- 10
	fmt.Println("发送成功")
}
```

3. 多返回值模式

多返回值模式用于判断通道是否关闭

```go
v, ok := <- ch		// value表示从通道取出的值；ok为false表示通道关闭，true表示通道未关闭
```

通常使用for range循环从通道接受值，当通道被关闭后，通道内所有值被接收完毕后会自动退出循环

```go
func f3(ch chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}
```

4. 单向通道

用于限制某个函数只能发送或只能接收，可以保证数据安全，即Producer不会在其他地方被其他人调用，并向通道执行发送操作

对一个只接收通道执行close也是不被允许的，因为默认通道的关闭操作应该由发送方来完成

```go
// Producer返回一个只接收通道
func Producer() <-chan int {
	ch := make(chan int, 2)
	// 创建一个心的goroutine，执行发送数据的任务
	go func() {
		for i := 0; i < 10; i++ {
			if i % 2 == 0 {
				ch <- i
			}
		}
		close(ch)			// 任务完成关闭通道
	}()
	return ch
}

// 将通道进行传递，这样接收方才能知道要从哪个通道接收数据
func Consumer(ch <-chan int) int {
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

func main() {
	ch := Producer()
	res := Consumer(ch)
	fmt.Println(res)
}
```

在函数传参及赋值操作中，全向通道可以转换为单向通道，但是单向通道无法转换为全向通道

```go
var ch = make(chan int, 1)
ch <- 10
close(ch)
Consumer(ch)		// 函数传参时，会将ch转为单向通道

var ch1 = make(chan int, 1)
ch1 <- 10
var ch2 <-chan int	// 声明一个只接收通道
ch1 = ch2			// 变量赋值时将ch1转换成单向通道
```

不同状态的通道执行不同操作的结果如下表

| 操作状态 | nil | 无值 | 有值 | 通道已满 |
|  :---:  | :---:  | :---: | :---:  | :---:     | 
| 发送     | 阻塞 | 发送成功 | 发送成功 | 阻塞 |
| 接收	   | 阻塞 | 阻塞 | 接收成功 | 接收成功 |
| 关闭     | panic | 关闭成功 | 关闭成功 | 关闭成功 |

### select多路复用

在需要同时从多个通道接收数据的场景下，使用select同时响应多个通道的操作，每个case分支对应一个通道的通信（发送或接收）过程，select会一直等待，直到某个case的通信操作完成时，就会执行case分支对应的语句，每个case只要有能接收，就执行，或者只要通道为空，能发送，就执行

select多路复用有以下特点

- 可处理一个或多个channel的发送或接收操作

- 如果多个case同时满足，select会随机选择一个执行

- 没有case的select会一直阻塞，可用于阻塞main函数，防止退出

```go
func main() {
	ch := make(chan int, 1)
	for i := 1; i <= 10; i++ {
		select {
		// 创建变量x来接收通道中的值，这一条case只有在通道中有值时才会执行
		case x := <-ch:
			fmt.Println(x)
		// 这一条case只有在通道中没有值时才会执行
		case ch <- i:
		}
		// 因此程序会先向通道中发送至，然后x从通道中接收值并打印，因为通道缓冲为1，因此只会一次发送一次接收，打印奇数13579
	}
}
```

### 通道误用示例

示例一：匿名函数所在的goroutine的接收操作在通道被关闭后会一直接收零值，并不会退出，接收操作应该使用`task, ok := <-ch`，并判断ok为假时退出；或者使用select来处理通道

```go
// demo1 通道误用导致的bug
func demo1() {
	wg := sync.WaitGroup{}
	ch := make(chan int, 10)
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)

	wg.Add(3)
	for j := 0; j < 3; j++ {
		go func() {
			// 这里for会把通道中所有值都读取出来，当通道中无值，则会取类型对应的零值，导致死循环
			for {
				task := <-ch
				fmt.Println(task)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
```

示例二：select命中了超时逻辑，导致通道中没有消费者，而定义的通道为无缓冲通道，因此会一直阻塞，导致goroutine泄露

```go
// demo2通道误用导致的bug
func demo2() {
	ch := make(chan int)
	go func() {
		// 这里假设执行一些耗时的操作
		time.Sleep(time.Second)
		ch <- 1
	}()

	select {
	case result := <-ch:
		fmt.Println(result)
		case <-time.After(time.Second):		// 较短的超时时间，time.After返回一个通道，该通道会在指定的时间后发送当前时间
			return 
	}
}
```

### 并发安全和锁

当代码中存在多个goroutine同时操作一个资源（临界区）时，会引发竞态问题

以下代码每次执行都会输出不同的结果，因为某个goroutine对全局变量x的修改可能覆盖另外一个goroutine中的操作

```go
var (
	x int64
	wg sync.WaitGroup
)

// add对全局变量x执行5000次加1操作
func add() {
	for i := 0; i < 500; i++ {
		x++
	}
	wg.Done()
}

func main() {
	wg.Add(2)
	go add()
	go add()
	wg.Wait()
	fmt.Println(x)
}
```

1. 互斥锁

用来保证同一时间只有一个goroutine能够访问共享资源的机制，Go语言中使用sync包的Mutex类型来实现互斥锁，一个goroutine执行加锁操作之后，其他goroutine就无法操作共享资源了，只有等锁释放之后才可以操作共享资源

`func (m *Mutex) Lock()`：获取互斥锁

`func (m *Mutex) Unlock()`：释放互斥锁

当互斥锁释放后，等待的goroutine才能获取锁并进入临界区，当多个goroutine同时等待一个锁时，唤醒的策略是随机的

```go
var (
	x int64
	wg sync.WaitGroup		// 等待组
	m sync.Mutex			// 互斥锁	
)

// add对全局变量x执行5000次加1操作
func add() {
	for i := 0; i < 500; i++ {
		m.Lock()	// 修改x前加锁
		x++
		m.Unlock()	// 修改完解锁
	}
	wg.Done()
}

func main() {
	wg.Add(2)
	go add()
	go add()
	wg.Wait()
	fmt.Println(x)
}
```

2. 读写互斥锁

现实中很多场景是读多写少的，当并发的读取一个资源而不涉及资源修改时，没必要加互斥锁，可以使用读写互斥锁

对于读锁，当一个goroutine获取读锁之后，其他goroutine可以继续获取读锁，如果获取了写锁就会等待

对于写锁，当一个goroutine获取写锁之后，其他goroutine无论获取读锁还是写锁都会等待

`func (rw *RWMutex) Lock()`：获取写锁

`func (rw *RWMutex) Unlock()`：释放写锁

`func (rw *RWMutex) RLock()`：获取读锁

`func (rw *RWMutex) RUnlock()`：释放读锁

`func (rw *RWMutex) RLocker() Locker`：返回一个实现Locker接口的读写锁

```go
// writeWithLock使用读写互斥锁的写操作
func writeWithLock() {
	rwMutex.Lock()		// 加写锁
	x = x + 1
	time.Sleep(10 * time.Millisecond)		// 假设读操作耗时10ms
	rwMutex.Unlock()		// 解写锁
	wg.Done()
}
// readWithRWLock使用读写互斥锁的读操作
func readWithRWLock() {
	rwMutex.RLock()		// 加读锁
	time.Sleep(time.Millisecond)		// 假设读操作耗时1ms
	rwMutex.RUnlock()		// 解读锁
	wg.Done()
}
```

3. sync.WaitGroup

sync.WaitGroup是一个结构体，进行参数传递时要传递指针

`func (wg *WaitGroup) Add(delta int)`：等待组中加delta个等待任务

`func (wg *WaitGroup) Done()`：等待组中减1个等待任务

`func (wg *WaitGroup) Wait()`：等待所有任务完成

4. sync.Once

用于确保某些操作即使在高并发时也只会被执行一次，如只加载一次配置文件

如果要执行的函数f需要传递参数，就需要搭配闭包使用

`func (o *Once) Do(f func())`：执行并只执行一次f，Do是Once类型的一个方法

```go
var icons map[string]image.Image

// icons初始化函数
func loadIcons() {
	icons = map[string]image.Image{
		"left": loadIcon("left.png"),
		"up": loadIcon("up.png"),
		"right": loadIcon("right.png"),
		"down": loadIcon("down.png"),
	}
}

// Icon函数被多个goroutine调用时不是并发安全的
func Icon(name string) image.Image {
	// 若icons未初始化，则调用loadIcons进行初始化
	if icons == nil {
		loadIcons()
	}
	// 并根据name从字典中查找，返回对应的加载函数
	return icons[name]
}
```

以上代码不是并发安全的，因为现代编译器和CPU可以在保证所有goroutine串行一致的基础上自由地重排访问内存的顺序（不太懂），因此loadIcons可能被重排为以下结果，在这种情况下，即使判断icons不为nil，也不意味着变量初始化完成了，虽然可以添加互斥锁来保证初始化icons时不会被其他goroutine操作，但是这样会引发性能问题

```go
func loadIcons() {
	icons = make(map[string]image.Image)
	icons["left"] = loadIcon("left.png")
	icons["up"] = loadIcon("up.png")
	icons["right"] = loadIcon("right.png")
	icons["down"] = loadIcon("down.png")
}
```

因此可以使用sync.OnceA进行改造，sync.Once内部包含一个互斥锁和一个布尔值，互斥锁保证布尔值和数据的安全，布尔值用来记录初始化是否完成，这样就能保证在初始化操作时是并发安全的，并且初始化操作不会被执行多次（不太懂）

```go
var icons map[string]image.Image
var loadIconsOnce sync.Once

// icons初始化函数
func loadIcons() {
	icons = map[string]image.Image{
		"left": loadIcon("left.png"),
		"up": loadIcon("up.png"),
		"right": loadIcon("right.png"),
		"down": loadIcon("down.png"),
	}
}

// Icon函数被多个goroutine调用时是并发安全的
func Icon(name string) image.Image {
	loadIconsOnce.Do(loadIcons)
	// 并根据name从字典中查找，返回对应的加载函数
	return icons[name]
}
```

5. sync.Map

Go语言内置的map不是并发安全的，不能在多个goroutine中并发对内置的map进行读写操作，否则会出现竞态问题，因为可能会发生不同的goroutine同时对一个内存资源进行读写的情况

> 竞态问题：在并发程序中，多个线程或goroutine同时访问共享资源，并且至少有一个线程试图修改这个资源，最终的结果依赖于线程的执行顺序，导致程序行为不确定的问题

Go语言sync包提供了一个开箱即用的并发安全版map————sync.map，不用像内置的map一样使用make函数进行初始化，还内置了很多方法

- `func (m *Map) Store(key, value interface{})`：存储k-v数据
- `func (m *Map) Load(key interface{}) (value interface{}, ok bool)`：获取k-v数据
- `func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool)`：查询k-v数据，如果没有则存储k-v数据
- `func (m *Map) LoadAndDelete(key interface{}) (value interface{}, loaded bool)`：查询k-v数据，并删除k-v数据
- `func (m *Map) Delete(key interface{})`：删除k-v数据
- `func (m *Map) Range(f func(key, value interface{}) bool)`：对map中的每个key-value依次调用f

```go
// 并发安全的map
var m = sync.Map{}
func mian() {
	wg := sync.WaitGroup{}
	// 对m执行20次并发的读写操作
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			key := strconv.Itoa(n)
			m.Store(key, n)				// 存储key-value
			value, _ := m.Load(key)		// 根据key取值
			fmt.Printf("m[%s]=%v\n", key, value)
			wg.Done()	
		} (i)
	}
}
```

### 原子操作

使用原子操作来保证整数数据类型（int32、uint32、int64、uint64）的并发安全通常比使用锁操作效率更高，由内置的标准库sync/atomic提供

```go
type Counter interface {
	Inc()
	Load() int64
}
// 普通版，非并发安全
type CommonCounter struct {
	counter int64
}
func (c CommonCounter) Inc() {
	c.counter++
}
func (c CommonCounter) Load() int64 {
	return c.counter
}

// 互斥锁版，并发安全
type MutexCounter struct {
	counter int64
	lock sync.Mutex
}
func (m *MutexCounter) Inc() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.counter++
}
func (m *MutexCounter) Load() int64 {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.counter
}

// 原子操作版，并发安全且比互斥锁效率高
type AtomicCounter struct {
	counter int64
}
func (a *AtomicCounter) Inc() {
	atomic.AddInt64(&a.counter, 1)
}
func (a *AtomicCounter) Load() int64 {
	return atomic.LoadInt64(&a.counter)
}
```

## 泛型

Go1.18增加了对泛型的支持，这是Go语言自开源以来所做的最大改变

### 泛型简介

即在编写某些代码时或数据结构时不提供值的类型，而是在实例化的时候提供

例如实现一个反转切片的函数reverse，`func reverse(s []int) []int`，这个函数只能接收[]int类型的参数，如果要支持[]float64类型的参数，就要再定义一个函数，以此类推，对于不同的类型需要重复编写相同的功能，在Go1.18之前可以使用反射解决这个问题，但是反射在运行期间获取变量类型会降低代码的执行效率，并且跳过编译阶段的类型检查，还会使得程序变得晦涩难懂，可以使用泛型解决这个问题，如下所示

```go
func reverseWithGenerics[T any] (s []T) []T {
	l := len(s)
	r := make([]T, l)
	for i, e := range s {
		r[l-i-1] = e
	}
	return r
}
```

### 泛型语法

泛型为Go语言添加了3个重要特性

- 函数和类型的类型参数
- 将接口类型定义为类型集，包括没有方法的类型
- 类型推断，允许在调用函数时在许多情况下省略类型参数

#### 类型参数

指类型本身也可以作为一种参数

1. 类型形参和类型实参

Go语言的函数和类型还支持添加类型参数，看起来像普通的参数列表，但是用方括号而不是圆括号，下面例子中，min函数同时支持int和float64两种类型

```go
func min[T int | float64] (a, b T) T {
	if a <= b {
		return a
	}
	return b
}
```

2. 类型实例化

调用min函数时既可以传入int类型参数，也可以传入float64类型参数

```go
m1 := min[int](1, 2)
m2 := min[float64](1.1, 2.2)
```

向min函数提供类型参数称为实例化，有两个步骤

- 编译器在整个泛型函数或类型中将所有类型形参替换为它们各自的类型实参

- 编译器验证每个类型参数是否满足相应的约束

成功实例化后可以得到一个非泛型函数，可以像其他函数一样被调用

```go
fmin := min[float64]	// 类型实例化，编译器生成 T = float64 的min函数
m2 = fmin(1.2, 2.3)
```

3. 使用类型参数

类型参数列表也可以在类型中使用，要使用泛型，必须实例化

```go
// type用于定义新的类型
type Slice[T int | string] []T
type Map[K int | string, V float32 | float64] map[K]V
type Tree[T interface{}] struct {
	left, right *Tree[T]
	value T
}
// 为Tree实现一个查找元素的方法
func (t *Tree[T]) Lookup(x T) *Tree[T] {...}

// 以上新类型可以这样使用
var intSlice Slice[int] = make([]int, 0)
var stringSlice Slice[string] = make([]string, 0)
var intFloatMap Map[int, float64] = make(map[int]float64)
var intTree Tree[int]
intTree.Value = 42
intTree.left = &Tree[int]{}
intTree.right = &Tree[int]{value: 100}
```

4. 类型约束

类型参数列表每个类型参数都有一个类型约束，类型约束定义了类型集，只有在这个类型集中的类型才能用作类型实参，Go语言的类型约束是接口类型，以下是类型约束两种常见的方式

```go
// 类型约束字面量，通常可以省略外层interface{}
func min[T interface{ int | float64}] (a, b T) T {
	if a <= b {
		return a
	}
	return b
}

// 类型约束使用的接口类型可以事先定义并支持复用
type Value interface {
	int | float64
}
func min[T Value] (a, b T) T {
	if a <= b {
		return a
	}
	return b
}
```

暂时还不知道为什么以下三种表达会产生歧义

```go
type IntPtrSlice [T *int] []T
type IntPtrSlice[T *int] []T
type IntPtrSlice[T interface{ *int }] []T
```

#### 类型集

从Go1.18开始，接口类型的定义发生了改变，由过去的接口类型定义方法集，变成了接口类型定义类型集，即接口类型可以作为值的类型，也可以用作类型约束

把接口类型当作类型集相较于方法集有一个优势：可以显式地为集合添加类型，从而以新的方式控制类型集

Go语言扩展了接口类型的语法，让我们能够为接口类型添加类型

```go
type V interface {
	int | string | bool		// 包含int、string和bool类型的类型集
}
```

从Go1.18开始，一个接口不仅可以嵌入其他接口，还可以嵌入任何类型、类型的并集或共享相同底层类型的无限类型集合

- `|`符号：T1 | T2，表示类型约束T1和T2这两个类型的并集

- `~`符号：~T，表示所有底层类型是T的类型集合；并且`~`符号后面只能是基本类型

1. any接口

Go1.18引入一个新的预声明标识符来作为空接口类型的别名

以下代码定义了一个名为foo的函数，`S ~[]E`表示S被约束为一个切片类型，其元素类型是E，而E可以是任何类型

`func foo[S ~[]E, E any] () {}`

2. comparable接口

该接口是由所有可比较类型实现的接口，只能用作类型参数约束，不能用作变量的类型

> 所有可比较类型：布尔值、数字、字符串、指针、通道、可比较类型的数组、字段均为可比较类型的结构体

`type MyMap[KEY comparable, VALUE any] map[KEY]VALUE`

3. Ordered类型约束

除了可比较（支持==和!=操作）的类型，通常还会用到可比较大小（支持>、<、>=、<=操作）的类型，Go语言没有内置这种类型，可以通过自行定义或使用golang.org/x/exp/constraints包中的定义

在以下例子中，E表示所有可比较大小的类型，T必须是整数类型

```go
import "golang.org/x/exp/constraints"
type MySlice[E Constraints.Ordered] []E
func min[T constraints.Integer] (a, b T) T {
	if a <= b {
		return a
	}
	return b
}
```

#### 类型推断

类型推断是Go1.18随泛型增加的一个新的主要语言特征

1. 函数参数类型推断

编译器可以从普通参数推断T的类型实参，可以使代码更短，函数实参类型推断只适用于函数参数中使用的类型参数，不适用于仅在函数结果或函数体中使用的类型参数

```go
func min[T int | float64] (a, b T) T {
	if a <= b {
		return a
	}
	return b
}

var a, b, m float64
// m = min[float64] (a, b)		显式指定类型实参
m = min(a,b) 					// 无需指定类型实参
```

2. 约束类型推断

以下是一个泛型函数，适用于任何整数类型的切片

```go
// Scale 返回切片的每个元素都乘以c的副本切片
func Scale[E constraints.Integer](s []E, c E) []E {
	r := make([]E, len(s))
	for i, v := range s {
		r[i] = v * c
	}
	return r
}
```

下面进行调用会编译失败，因为Scale函数返回类型为[]E的值，当我们使用Point类型的值调用Scale时，返回的是[]int32类型的值，而不是Point类型，因此没有String方法，因此编译器会报错

```go
type Point []int32
func (p Point) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func ScaleAndPrint(p Point) {
	r := Scale(p, 2)
	fmt.Println(r.String())		// 会编译失败
}
```

要解决以上问题，就必须更改Scale函数

```go
func Scale[S ~[]E, E constraints.Integer](s S, c E) S {
	r := make(S, len(s))			// 需要更改make的方式
	for i, v := range s {
		r[i] = v * c
	}
	return r
}
```

此时就可以通过`Scale(p, 2)`进行调用，返回的也是Point类型，从而可以调用String方法，因为函数参数类型推断会让编译器推断S的类型是Point类型，也就是[]int32类型，而2是一个非类型化的常量，所以函数参数类型推断无法推断出E的正确类型

编译器从S的类型是[]int32进而推断出E的类型是int32的过程被称为约束类型推断，也即从类型参数约束中推导出类型实参，即当一个类型参数包含其他类型参数定义的约束时，会使用约束类型推断，当其中一个类型参数的类型实参已知时，便可以推断出另一个类型参数的类型实参，在例子中就是，S是~[]E，后面跟着一个用另一个类型参数写的类型[]E，如果我们知道S的类型实参，就可以推断出E的类型实参，S是一个切片类型，而E是该切片的元素类型

### 类型参数的适用场景

#### 应该使用类型参数的场景

1. 在使用语言定义的容器类型时

如slice、map和channel，若函数具有包含这些类型的参数，并且函数的代码不关心元素的类型

如下函数并不关心map中键的类型，也没有使用map值类型，适用于任何map类型，可以使用类型参数

```go
// MapKeys 返回m中所有key组成的切片
func MapKeys[Key comparable, Val any] (m map[Key]Val) []Key {
	s := make([]Key, 0, len(m)
	for k := range m ) {
		s = append(s, k)
	}
	return s
}
```

2. 通用数据结构

通用数据结构类似于slice或map，但不是Go语言内置的，如链表或二叉树，用类型参数替换特定的元素类型可以生成更通用的数据结构，用类型参数替换接口类型可以更有效地存储数据、节省内存资源，同时可以避免类型断言，并在构建时进行完整的类型检查

如以下例子中，树中每个结点都包含类型参数T的值，在使用特定类型参数实例化树时，该类型的值将直接存储在结点中，不会被存储为接口类型（因为要实现任意类型的话，要使用空接口）

```go
// Tree定义一个二叉树
type Tree[T any] struct {
	cmp func(T, T) int 
	root *node[T]
}
// 二叉树的一个节点
type node[T any] struct {
	left, right *node[T]
	val T
}
// find查找值
func (bt *Tree[T]) find(val T) **node[T] {
	pl := &bt.root
	for *pl != nil {
		switch cmp := bt.cmp(val, (*pl).val); {
		case cmp < 0:
			pl = &(*pl).left
		case cmp > 0:
			pl = &(*pl).right
		default:
			return pl
		}
	}
	return pl
}
```

3. 对于类型参数，优先选择函数而不是方法

上面Tree示例说明了一个一般原则：当需要比较函数等内容时，优先使用函数而不是方法；即将方法转换为函数要比将方法添加到类型中简单得多，因此对于通用数据类型，优先使用函数而不是编写需要方法的约束（不太懂）

4. 实现通用方法

当不同类型需要实现某些公共方法，而这些类型的实现看起来是相同的，如标准库sort.Interface它要求类型实现Len、Swap和Less三个方法

对以下这类代码使用类型参数是合适的，因为所有切片类型的方法看起来完全相同（get不到）

```go
// 泛型类型SliceFn为T类型切片实现sort.Interface
type SliceFn[T any] struct {
	s []T
	less func(T, T) bool
}
func (s SliceFn[T]) Len() int {
	return len(s.s)
}
func (s SliceFn[T]) Swap(i, j int) {
	s.s[i], s.s[j] = s.s[j], s.s[i]
}
func (s SliceFn[T]) Less(i, j int) bool {
	return s.less(s.s[i], s.s[j])
}
```

#### 不应该使用类型参数

1. 不要用类型参数替换接口类型

如果只需要调用某个类型的值的方法，则使用接口类型，而不是类型参数，比如io.Reader易读、高效，不需要使用类型参数，通过调用read方法就可以从值中读取数据，即不需要将以下第一个函数签名更改为第二个版本

Go中使用类型参数通常不会比使用接口类型快

```go
func ReadSome(r io.Reader) ([]byte, error)
func ReadSome[T io.Reader] (r T) ([]byte, error)
```

2. 如果方法实现不同，则不要使用类型参数

如果一个方法的实现对于所有类型都是相同的，则使用类型参数；如果每种类型的实现不同，则使用接口类型并编写不同的实现方法

例如从文件读取的实现和从随机数生成器读取的实现完全不同，我们应该编写两个不同的Read方法，并使用io.Reader这样的接口类型

3. 在适当的地方使用反射

反射能实现某种意义上的泛型编程，因为它可以编写适用于任何类型的代码

如果某些操作必须支持没有方法的类型（即不能使用接口类型），并且每个类型的操作都不同（不能使用类型参数），则应该使用反射，encoding/json包就是一个例子，不要求进行编码的每个类型都有一个MarshalJSON方法，所以不能使用接口类型，而对接口类型进行编码与对结构体类型进行编码完全不同，因此不能使用类型参数，因此encoding/json包使用反射实现

最终指导原则：发现自己多次编写几乎完全相同的代码，他们之间的唯一区别是使用的类型不同，就可以考虑是否可以使用类型参数

## 错误处理




## 字符串操作




## 字符串格式化




## JSON处理




## 时间处理




## 数字解析




## 进程信息



## 测试






**未完待续......**










