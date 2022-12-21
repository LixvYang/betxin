package mixinorder

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

func UpdateMixinOrder(c *gin.Context) {
	var mixinOrder *model.MixinOrder
	traceId := c.Param("traceId")
	if err := c.ShouldBindJSON(&mixinOrder); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	code := model.UpdateMixinOrder(traceId, mixinOrder)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_UPDATE_CATENAME, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, mixinOrder)
}
