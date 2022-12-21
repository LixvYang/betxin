package collect

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"

	"github.com/gin-gonic/gin"
)

type DeleteRequest struct {
	UserId string `json:"user_id"`
	Tid    string `json:"tid"`
}

func DeleteCollect(c *gin.Context) {
	var r DeleteRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}
	code := model.DeleteCollect(r.UserId, r.Tid)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_DELETE_CATENAME, nil)
		return
	}
	betxinredis.DelKeys(v1.COLLECT_GET_USER_LIST+r.UserId, v1.COLLECT_GET_USER_TOTAL+r.UserId, v1.COLLECT_LIST, v1.COLLECT_TOTAL)

	v1.SendResponse(c, errmsg.SUCCSE, r.Tid)
}
