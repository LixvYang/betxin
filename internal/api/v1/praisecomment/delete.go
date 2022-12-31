package praisecomment

import (
	"strconv"

	"github.com/lixvyang/betxin/internal/model"

	v1 "github.com/lixvyang/betxin/internal/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"

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

	if betxinredis.Exists(v1.PARISECOMMENT + strconv.Itoa(r.Cid)) {
		betxinredis.SREM(v1.PARISECOMMENT+strconv.Itoa(r.Cid), v1.PARISECOMMENT_KEY+r.Uid)
	}

	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
