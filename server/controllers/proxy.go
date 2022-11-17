package controllers

import (
	"fmt"
	"github.com/Xhofe/alist/conf"
	"github.com/Xhofe/alist/drivers/base"
	"github.com/Xhofe/alist/drivers/operate"
	"github.com/Xhofe/alist/server/common"
	"github.com/Xhofe/alist/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
)

func Proxy(c *gin.Context) {
	rawPath := c.Param("path")
	rawPath = utils.ParsePath(rawPath)
	log.Debugf("proxy: %s", rawPath)
	account, path, driver, err := common.ParsePath(rawPath)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	// 只有以下几种情况允许中转：
	// 1. 账号开启中转
	// 2. driver只能中转
	// 3. 是文本类型文件
	// 4. 开启webdav中转（需要验证sign）
	if !account.Proxy && !driver.Config().OnlyProxy &&
		utils.GetFileType(filepath.Ext(rawPath)) != conf.TEXT &&
		!utils.IsContain(conf.DProxyTypes, utils.Ext(rawPath)) {
		// 只开启了webdav中转，验证sign
		ok := false
		if account.WebdavProxy {
			_, ok = c.Get("sign")
		}
		if !ok {
			common.ErrorStrResp(c, fmt.Sprintf("[%s] not allowed proxy", account.Name), 403)
			return
		}
	}
	// 中转时有中转机器使用中转机器，若携带标志位则表明不能再走中转机器了
	if account.DownProxyUrl != "" && c.Query("d") != "1" {
		name := utils.Base(rawPath)
		link := fmt.Sprintf("%s%s?sign=%s", strings.Split(account.DownProxyUrl, "\n")[0], rawPath, utils.SignWithToken(name, conf.Token))
		c.Redirect(302, link)
		return
	}
	// 检查文件
	file, err := operate.File(driver, account, path)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	// 对于中转，不需要重设IP
	link, err := driver.Link(base.Args{Path: path, Header: c.Request.Header}, account)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	err = common.Proxy(c.Writer, c.Request, link, file)
	log.Debugln("web proxy error:", err)
	if err != nil {
		common.ErrorResp(c, err, 500)
	}
}

var client *resty.Client

func init() {
	client = resty.New()
	client.SetRetryCount(3).SetTimeout(base.DefaultTimeout)
}

func Text(c *gin.Context, link *base.Link) {
	req := client.R()
	if link.Headers != nil {
		for _, header := range link.Headers {
			req.SetHeader(header.Name, header.Value)
		}
	}
	res, err := req.Get(link.Url)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	text := res.String()
	t := utils.GetStrCoding(res.Body())
	log.Debugf("text type: %s", t)
	if t != utils.UTF8 {
		body, err := utils.GbkToUtf8(res.Body())
		if err != nil {
			common.ErrorResp(c, err, 500)
			return
		}
		text = string(body)
	}
	c.String(200, text)
}
