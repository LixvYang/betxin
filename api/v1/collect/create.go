package collect

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"

	"github.com/gin-gonic/gin"
)

func CreateCollect(c *gin.Context) {
	var r model.Collect
	if err := c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	if code := model.CreateCollect(&r); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	betxinredis.DelKeys(v1.COLLECT_LIST, v1.COLLECT_TOTAL, v1.COLLECT_GET_USER_LIST+r.UserId, v1.COLLECT_GET_USER_TOTAL+r.UserId)
	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
