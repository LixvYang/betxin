package usertotopic

import (
	"fmt"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUserToTopic(c *gin.Context) {
	tid := c.Param("id")
	session := sessions.Default(c)
	user := session.Get("userId")
	userId := fmt.Sprintf("%v", user)
	data, code := model.GetUserToTopic(userId, tid)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, data)
}
