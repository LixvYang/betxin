package praisecomment

import (
	"strconv"

	"github.com/lixvyang/betxin/internal/model"

	v1 "github.com/lixvyang/betxin/internal/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"

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
	betxinredis.SADD(v1.PARISECOMMENT+strconv.Itoa(data.Cid), v1.PARISECOMMENT_KEY+data.Uid)

	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
