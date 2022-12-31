package comment

import (
	"math"

	"github.com/go-redis/redis/v8"
	"github.com/lixvyang/betxin/internal/model"
	"github.com/lixvyang/betxin/pkg/convert"

	v1 "github.com/lixvyang/betxin/internal/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
)

func CreateComment(c *gin.Context) {
	var data model.Comment
	if err := c.ShouldBindJSON(&data); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	if code := model.CreateComment(&data); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	if betxinredis.Exists(v1.COMMENT_LIST + data.Tid) {
		betxinredis.ZADD(v1.COMMENT_LIST+data.Tid, &redis.Z{
			Score:  float64(int(math.Pow(10, 14))*(data.PraiseNum+1) + v1.COMMENT_MAXTIME - int(data.CreatedAt.UnixMilli())),
			Member: convert.Marshal(data),
		})
	}

	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
