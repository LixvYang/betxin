package mixpayorder

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

func GetMixpayOrder(c *gin.Context) {
	tracdId := c.Param("traceId")
	mixpayorder, code := model.GetMixpayOrder(tracdId)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	v1.SendResponse(c, errmsg.SUCCSE, mixpayorder)
}
