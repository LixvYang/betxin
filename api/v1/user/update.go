package user

import (
	"log"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"

	"github.com/gin-gonic/gin"
)

func UpdateUser(c *gin.Context) {
	var user *model.User
	userId := c.Param("userId")
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Panicln(err)
	}
	code := model.CheckUser(userId)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	code = model.UpdateUser(userId, user)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	betxinredis.DelKeys(v1.USER_INFO+userId, v1.USER_LIST, v1.USER_TOTAL)

	v1.SendResponse(c, errmsg.SUCCSE, userId)
}
