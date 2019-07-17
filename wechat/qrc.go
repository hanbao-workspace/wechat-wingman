package wechat

import (
	"errors"
	"fmt"
	"github.com/skip2/go-qrcode"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	BrightBlack = "\033[48;5;0m  \033[0m"

	BrightWhite = "\033[48;5;7m  \033[0m"
)
const host = "https://wx.qq.com"

func CreateQRC(content string) {
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		fmt.Println(err)
		return
	}
	for ir, row := range qr.Bitmap() {
		lr := len(row)

		if ir == 0 || ir == 1 || ir == 2 ||
			ir == lr-1 || ir == lr-2 || ir == lr-3 {
			continue
		}

		for ic, col := range row {
			lc := len(qr.Bitmap())
			if ic == 0 || ic == 1 || ic == 2 ||
				ic == lc-1 || ic == lc-2 || ic == lc-3 {
				continue
			}
			if col {
				fmt.Print(BrightBlack)
			} else {
				fmt.Print(BrightWhite)
			}
		}
		fmt.Println()
	}
}
func getTimeStamp() int64 {
	return time.Now().UnixNano() / 1e6
}
func getUUID() (string, error) {
	stamp := int(getTimeStamp())
	value := url.Values{"fun": {"new"}, "lang": {"zh_CN"}, "appid": {"wx782c26e4c19acffb"}, "_": {strconv.Itoa(stamp)}}
	res, err := http.Post(host+"/jslogin", "application/x-www-form-urlencoded", strings.NewReader(value.Encode()))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	reg := regexp.MustCompile(`window.QRLogin.code = (\d+); window.QRLogin.uuid = "(.+)";`)
	result := reg.FindStringSubmatch(string(b))
	if result[1] != "200" || result[2] == "" {
		return "", errors.New("获取wechatuuid失败")
	}
	return result[2], nil
}
func GetQRC() error {
	uuid, err := getUUID()
	if err != nil {
		return err
	}
	res, err := http.Get(host + "/qrcode/" + uuid)
	if err != nil {
		return err
	}
	img, err := jpeg.Decode(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(img.Bounds())
	img.
	return nil
}
