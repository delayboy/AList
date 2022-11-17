package common

import (
	"fmt"
	"github.com/Xhofe/alist/drivers/base"
	"github.com/Xhofe/alist/model"
	"github.com/Xhofe/alist/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PathReq struct {
	Path     string `json:"path"`
	Password string `json:"password"`
	PageNum  int    `json:"page_num"`
	PageSize int    `json:"page_size"`
}

func ParsePath(rawPath string) (*model.Account, string, base.Driver, error) {
	rawPath = utils.ParsePath(rawPath)
	account, ok := model.GetBalancedAccount(rawPath)
	if !ok {
		return nil, "", nil, fmt.Errorf("path not found")
	}
	driver, ok := base.GetDriver(account.Type)
	if !ok {
		return nil, "", nil, fmt.Errorf("no [%s] driver", account.Type)
	}
	name := utils.ParsePath(account.Name)
	bIndex := strings.LastIndex(name, ".balance")
	if bIndex != -1 {
		name = name[:bIndex]
	}
	//if name == "/" {
	//	name = ""
	//}
	return &account, utils.ParsePath(strings.TrimPrefix(rawPath, name)), driver, nil
}

func ErrorResp(c *gin.Context, err error, code int) {
	log.Error(err.Error())
	c.JSON(200, Resp{
		Code:    code,
		Message: err.Error(),
		Data:    nil,
	})
	c.Abort()
}

func ErrorStrResp(c *gin.Context, str string, code int) {
	log.Error(str)
	c.JSON(200, Resp{
		Code:    code,
		Message: str,
		Data:    nil,
	})
	c.Abort()
}

func SuccessResp(c *gin.Context, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(200, Resp{
			Code:    200,
			Message: "success",
			Data:    nil,
		})
		return
	}
	c.JSON(200, Resp{
		Code:    200,
		Message: "success",
		Data:    data[0],
	})
}

func Hide(meta *model.Meta, files []model.File) []model.File {
	if meta == nil {
		return files
	}
	if meta.Hide != "" {
		tmpFiles := make([]model.File, 0)
		hideFiles := strings.Split(meta.Hide, ",")
		for _, item := range files {
			if !utils.IsContain(hideFiles, item.Name) {
				tmpFiles = append(tmpFiles, item)
			}
		}
		files = tmpFiles
	}
	if meta.OnlyShows != "" {
		tmpFiles := make([]model.File, 0)
		showFiles := strings.Split(meta.OnlyShows, ",")
		for _, item := range files {
			if utils.IsContain(showFiles, item.Name) {
				tmpFiles = append(tmpFiles, item)
			}
		}
		files = tmpFiles
	}
	return files
}
