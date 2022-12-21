package praisecomment

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

type DeleteRequest struct {
	Cid int    `json:"cid"`
	Uid string `json:"uid"`
}

func DeletePraiseComment(c *gin.Context) {
	// uid cid
	var r DeleteRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	if code := model.DeletePraise(r.Cid, r.Uid); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
