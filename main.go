package main

import (
	"fmt"
	"log"
	"wechat-wingman/wechat"
)

func main() {
	//defer func() {
	//	for {
	//	}
	//}()
	wechat := wechat.New()
	if err := wechat.Login(); err != nil {
		log.Panicln(err)
	}
	fmt.Println("登录成功")
	wechat.Init()

}
