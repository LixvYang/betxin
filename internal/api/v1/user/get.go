package user

import (
	"fmt"

	"github.com/lixvyang/betxin/internal/model"
	"github.com/lixvyang/betxin/pkg/convert"

	v1 "github.com/lixvyang/betxin/internal/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUserInfoByUserId(c *gin.Context) {
	var code int
	var data model.User
	session := sessions.Default(c)
	user := session.Get("userId")
	userId := fmt.Sprintf("%v", user)

	userInfo := v1.USERINFO + userId
	if betxinredis.Exists(userInfo) {
		info := betxinredis.Get(userInfo).Val()
		convert.Unmarshal(info, &data)
		v1.SendResponse(c, errmsg.SUCCSE, data)
	} else {
		data, code = model.GetUserById(userId)
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR, nil)
			return
		}
		info := convert.Marshal(data)
		betxinredis.Set(userInfo, info, v1.REDISEXPIRE)
		v1.SendResponse(c, errmsg.SUCCSE, data)
	}
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
