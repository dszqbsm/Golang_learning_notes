# Go常见标准库使用

整理我在写代码过程中常用的标准库，以及一些常用的方法，并用例子展现使用说明，作为自己的标准库使用手册

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



## strconv

`strconv`包实现了基本数据类型和其字符串表示的相互转换

- `func ParseBool(str string) (value bool, err error)`：返回字符串表示的bool值

- `func ParseInt(s string, base int, bitSize int) (i int64, err error)`：返回字符串表示的整数值，接受正负号，base指定进制（2到36），如果base为0，则会从字符串前置判断，"0x"是16进制，"0"是8进制，否则是10进制；bitSize指定结果必须能无溢出赋值的整数类型，0、8、16、32、64分别代表int、int8、int16、int32、int64；返回的err是*NumError类型的，如果语法有误，err.Error = ErrSyntax；如果结果超出类型范围err.Error = ErrRange

- `func ParseUint(s string, base int, bitSize int) (uint64, error)`：返回字符串表示的无符号整数值

- `func ParseFloat(s string, bitSize int) (float64, error)`：返回字符串表示的浮点数值，bitSize指定了结果必须能无溢出赋值的整数类型，32或64；返回的err是*NumError类型的，如果语法有误，err.Error = ErrSyntax；如果结果超出类型范围err.Error = ErrRange

- `func FormatBool(b bool) string`：返回布尔类型b的字符串表示

- `func FormatInt(i int64, base int) string`：返回整数类型i的base进制的字符串表示

- `func FormatUint(i uint64, base int) string`：返回无符号整数类型i的base进制的字符串表示

- `func FormatFloat(f float64, fmt byte, prec, bitSize int) string`：返回浮点数类型f的字符串表示，fmt表示格式，prec表示精度，bitSize表示类型

- `func Atoi(s string) (int, error)`：Atoi是ParseInt的简写，返回字符串表示的整数值

- `func Itoa(i int) string`：Itoa是FormatInt的简写，返回整数类型i的字符串表示


## time





