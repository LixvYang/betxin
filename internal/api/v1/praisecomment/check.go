package praisecomment

import (
	"strconv"

	"github.com/lixvyang/betxin/internal/model"

	v1 "github.com/lixvyang/betxin/internal/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
)

type CheckPraiseRequest struct {
	Cid int    `json:"cid"` // commment_id
	Uid string `json:"uid"`
}

func CheckPraise(c *gin.Context) {
	var r CheckPraiseRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	if betxinredis.Exists(v1.PARISECOMMENT + strconv.Itoa(r.Cid)) {
		if betxinredis.SISMEMBER(v1.PARISECOMMENT+strconv.Itoa(r.Cid), v1.PARISECOMMENT_KEY+r.Uid) {
			v1.SendResponse(c, errmsg.SUCCSE, nil)
		}
		return
	}

	if code := model.CheckPraiseComment(r.Cid, r.Uid); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	betxinredis.SADD(v1.PARISECOMMENT+strconv.Itoa(r.Cid), v1.PARISECOMMENT_KEY+r.Uid)

	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
