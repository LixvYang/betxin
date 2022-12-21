package praisecomment

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

func CreatePraiseComment(c *gin.Context) {
	var data model.PraiseComment
	if err := c.ShouldBindJSON(&data); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	if code := model.CreatePraiseComment(&data); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
