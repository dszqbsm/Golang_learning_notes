package main

import "fmt"

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
	Checkout(&ZhiFuBao{}) // 使用支付宝支付
	Checkout(&WeChat{})   // 使用微信支付
}
