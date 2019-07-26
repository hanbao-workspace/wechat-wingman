package wechat

import (
	"errors"
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/skip2/go-qrcode"
	scanCode "github.com/tuotoo/qrcode"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	BrightBlack = "\033[48;5;0m  \033[0m"
	BrightWhite = "\033[48;5;7m  \033[0m"
)
const loginHost = "http://login.weixin.qq.com"

func PrintQRC(bitMap [][]bool) error {
	var s string
	for ir, row := range bitMap {
		lr := len(row)
		if ir == 0 || ir == 1 || ir == 2 ||
			ir == lr-1 || ir == lr-2 || ir == lr-3 {
			continue
		}

		for ic, col := range row {
			lc := len(bitMap)
			if ic == 0 || ic == 1 || ic == 2 ||
				ic == lc-1 || ic == lc-2 || ic == lc-3 {
				continue
			}
			if col {
				s += BrightBlack
			} else {
				s += BrightWhite
			}
		}
		s += fmt.Sprintln()
	}
	outer := colorable.NewColorableStdout()
	if _, err := fmt.Fprint(outer, s); err != nil {
		return err
	}
	return nil
}

func (wechat *Wechat) getUUID() (err error) {
	stamp := int(getTimeStamp())
	value := url.Values{"fun": {"new"}, "lang": {"zh_CN"}, "appid": {"wx782c26e4c19acffb"}, "_": {strconv.Itoa(stamp)}}
	res, err := http.Post(host+"/jslogin", "application/x-www-form-urlencoded", strings.NewReader(value.Encode()))
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	reg := regexp.MustCompile(`window.QRLogin.code = (\d+); window.QRLogin.uuid = "(.+)";`)
	result := reg.FindStringSubmatch(string(b))
	if result[1] != "200" || result[2] == "" {
		return errors.New("获取wechatuuid失败")
	}
	wechat.uuid = result[2]
	return nil
}

func (wechat *Wechat) getQRC() (err error) {
	if err = wechat.getUUID(); err != nil {
		return
	}

	res, err := http.Get(host + "/qrcode/" + wechat.uuid)
	if err != nil && res.StatusCode != 200 {
		return
	}
	defer res.Body.Close()
	scanResult, err := scanCode.Decode(res.Body)
	if err != nil {
		return
	}
	qr, err := qrcode.New(scanResult.Content, qrcode.Medium)
	if err != nil {
		return
	}
	if err = PrintQRC(qr.Bitmap()); err != nil {
		return
	}
	return
}
