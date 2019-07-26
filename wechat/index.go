package wechat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

const (
	host           = "https://wx.qq.com"
	listenScanPath = "/cgi-bin/mmwebwx-bin/login"
)

type Wechat struct {
	uuid        string
	signStatus  bool
	ticket      string
	scan        string
	redirectUri string
}

func New() *Wechat {
	return &Wechat{
		signStatus: false,
	}
}

func getTimeStamp() int64 {
	return time.Now().UnixNano() / 1e6
}

func (wechat *Wechat) Login() (err error) {
	if err = wechat.getQRC(); err != nil {
		return
	}
	fmt.Println("tips:")
	fmt.Println("请使用手机微信扫描二维码登录")
	tip := 0
	for {
		str, err := wechat.listeningScan(tip)
		if err != nil {
			return err
		}
		reg := regexp.MustCompile(`window.code=(\d+);`)
		result := reg.FindStringSubmatch(str)
		if len(result) >= 2 {
			code, err := strconv.Atoi(result[1])
			if err != nil {
				return err
			}
			if code == 201 {
				tip = 1
			}
			if code == 200 {
				reg = regexp.MustCompile(`window.redirect_uri="(\S+)";`)
				result2 := reg.FindStringSubmatch(str)
				urlResult, err := url.Parse(result2[1])
				if err != nil {
					return err
				}
				query := urlResult.Query()
				//ticket
				if ticket, ok := query["ticket"]; !ok {
					return errors.New("not found ticket,please try again later")
				} else {
					wechat.ticket = ticket[0]
				}
				//scan
				if scan, ok := query["scan"]; !ok {
					return errors.New("not found scan,please try again later")
				} else {
					wechat.scan = scan[0]
				}
				wechat.redirectUri = urlResult.Scheme + "://" + urlResult.Host + urlResult.Path
				wechat.signStatus = true
				return nil
			}
		}
		time.Sleep(3 * time.Second)
	}
}
func (wechat *Wechat) listeningScan(tip int) (str string, err error) {
	stamp := int(getTimeStamp())
	value := url.Values{"tip": {strconv.Itoa(tip)}, "uuid": {wechat.uuid}, "loginicon": {"false"}, "_": {strconv.Itoa(stamp)}}
	res, err := http.Get(host + listenScanPath + "?" + value.Encode())
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	return string(b), nil

}

func (wechat *Wechat) Init() (err error) {
	if !wechat.signStatus {
		return errors.New("微信登录状态异常")
	}
	value := url.Values{"ticket": {wechat.ticket}, "uuid": {wechat.uuid}, "lang": {"zh_CN"}, "scan": {wechat.scan}, "fun": {"new"}}
	res, err := http.Get(wechat.redirectUri + "?" + value.Encode())
	if err != nil {
		return
	}
	fmt.Println(res.Header)

	return
}
