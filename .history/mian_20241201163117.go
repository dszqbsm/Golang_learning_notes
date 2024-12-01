package main

// 可以单个引入

type Student struct {
	ID     int    `json:"stu_id"` // 设置JSON标签，JSON序列化时将stu_id作为键名
	Gender string // JSON序列化时默认将字段名作为key
	name   string // 私有字段，不能被JSON包访问
}

func main() {

}
