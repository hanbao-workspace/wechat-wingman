package main

import (
	"fmt"
	_ "github.com/skip2/go-qrcode"
	"wechat-wingman/wechat"
	//"wechat-wingman/wechat"
)

func main() {
	//wechat.GetQRC("http://www.baidu.com")
	a := wechat.GetQRC()
	fmt.Println(a)
}
