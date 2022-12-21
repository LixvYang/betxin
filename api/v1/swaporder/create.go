package swaporder

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

func CreateSwapOrder(c *gin.Context) {
	var r model.SwapOrder
	if err := c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	if code := model.CreateSwapOrder(&r); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, r)
}
