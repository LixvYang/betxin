package user

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"

	"github.com/gin-gonic/gin"
)

func DeleteUser(c *gin.Context) {
	userId := c.Param("userId")
	code := model.DeleteUser(userId)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_DELETE_CATENAME, nil)
		return
	}
	betxinredis.DelKeys(v1.USER_INFO+userId, v1.USER_LIST, v1.USER_TOTAL)
	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
