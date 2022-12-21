package feedback

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

func CreateFeedback(c *gin.Context) {
	var r model.FeedBack
	var err error
	if err = c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	if code := model.CreateFeedBack(&r); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
