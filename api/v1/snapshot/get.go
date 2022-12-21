package snapshot

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

func GetMixinNetworkSnapshot(c *gin.Context) {
	traceId := c.Param("traceId")

	if code := model.CheckMixinNetworkSnapshot(traceId); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	mixinNetworkSnapshot, code := model.GetMixinNetworkSnapshot(traceId)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	v1.SendResponse(c, errmsg.SUCCSE, mixinNetworkSnapshot)
}
