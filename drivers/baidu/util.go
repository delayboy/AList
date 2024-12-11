package baidu

import (
	"net/url"
	"strings"
	"time"
)

func getTime(t int64) *time.Time {
	tm := time.Unix(t, 0)
	return &tm
}

func encodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.ReplaceAll(r, "+", "%20") //万一遇上加号进行一个保底，用空格代替，理论上经过url编码后加号会被隐去，因此这个replace不会执行
	return r
}
