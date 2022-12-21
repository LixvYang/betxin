package usertotopic

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

type CheckUserToTid struct {
	UserId string `json:"user_id"`
	Tid    string `json:"tid"`
}

func CheckUserToTopic(c *gin.Context) {
	var r CheckUserToTid
	if err := c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	code := model.CheckUserToTopic(r.UserId, r.Tid)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	v1.SendResponse(c, errmsg.SUCCSE, r.Tid)
}
