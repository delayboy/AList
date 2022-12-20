package controllers

import (
	"github.com/Xhofe/alist/conf"
	"github.com/Xhofe/alist/server/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ClearCache() error {
	err := conf.Cache.Clear(conf.Ctx)
	log.Info("cache has been cleared, we will redo file sync")
	return err
}

func ClearCacheEntry(c *gin.Context) {
	err := ClearCache()
	if err == nil {
		common.SuccessResp(c, 0)
	} else {
		common.ErrorResp(c, err, -1)
	}

}
