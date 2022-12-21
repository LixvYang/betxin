package message

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

func UpdateMessage(c *gin.Context) {
	var msg *model.MixinMessage
	msgId := c.Param("id")
	if err := c.ShouldBindJSON(&msg); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	code := model.UpdateMixinMessageByMsgId(msgId, msg)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_UPDATE_CATENAME, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, msgId)
}
