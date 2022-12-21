package user

import (
	"fmt"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUserInfoByUserId(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("userId")
	userId := fmt.Sprintf("%v", user)

	data, code := model.GetUserById(userId)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	// user := convert.Marshal(&data)
	// betxinredis.Set(v1.USER_INFO+userId, user, v1.REDISEXPIRE)
	v1.SendResponse(c, errmsg.SUCCSE, data)
}

func GetUserInfoByUserFullName(c *gin.Context) {
	fullName := c.Param("fullName")

	user, code := model.GetUserByName(fullName)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	v1.SendResponse(c, errmsg.SUCCSE, user)
}
