package usertotopic

import (
	"log"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

func UpdateUserToTopic(c *gin.Context) {
	var userToTopic *model.UserToTopic
	if err := c.ShouldBindJSON(&userToTopic); err != nil {
		log.Panicln(err)
	}

	code := model.UpdateUserToTopic(userToTopic)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_UPDATE_CATENAME, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, userToTopic.Id)
}
