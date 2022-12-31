package comment

import (
	"fmt"
	"math"

	"github.com/go-redis/redis/v8"
	"github.com/lixvyang/betxin/internal/model"
	"github.com/lixvyang/betxin/pkg/convert"

	v1 "github.com/lixvyang/betxin/internal/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
)

type ListResponse struct {
	TotalCount int             `json:"totalCount"`
	List       []model.Comment `json:"list"`
}

type ListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func ListCommentByTid(c *gin.Context) {
	var err error
	var total int
	var data []model.Comment
	var code int
	var comments []string

	tid := c.Param("tid")
	var r ListRequest
	if err = c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	var totalCommentsByTid = v1.COMMENT_LIST + tid
	total = betxinredis.ZCARD(totalCommentsByTid)
	comments, _ = betxinredis.ZREVRANGE(totalCommentsByTid, r.Offset, r.Limit)
	data = []model.Comment{}
	for _, comment := range comments {
		var pc model.Comment
		convert.Unmarshal(comment, &pc)
		data = append(data, pc)
	}
	// redis
	if !betxinredis.Exists(totalCommentsByTid) {
		data, _, code = model.ListComments(tid)
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR, nil)
			return
		}

		// 将数据存入redis
		var members []*redis.Z
		for _, pc := range data {
			Z := &redis.Z{
				Score:  float64(int(math.Pow(10, 14))*(pc.PraiseNum+1) + v1.COMMENT_MAXTIME - int(pc.CreatedAt.UnixMilli())),
				Member: convert.Marshal(pc),
			}
			members = append(members, Z)
			fmt.Println(int(pc.CreatedAt.UnixMilli()))
			fmt.Println(int(math.Pow(10, 13))*(pc.PraiseNum+1) + v1.COMMENT_MAXTIME - int(pc.CreatedAt.UnixMilli()))
		}
		betxinredis.ZADD(totalCommentsByTid, members...)

		// 查找 redis中的数据
		// 点赞数最多的在最后面
		total = betxinredis.ZCARD(totalCommentsByTid)
		comments, _ = betxinredis.ZREVRANGE(totalCommentsByTid, r.Offset, r.Limit)
		data = []model.Comment{}
		for _, comment := range comments {
			var pc model.Comment
			convert.Unmarshal(comment, &pc)
			data = append(data, pc)
		}
		v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
			TotalCount: total,
			List:       data,
		})
	} else {
		v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
			TotalCount: total,
			List:       data,
		})
	}
}

func ListComment(c *gin.Context) {
	var total int
	var data []model.Comment
	var code int
	var r ListRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}
	switch {
	case r.Offset >= 100:
		r.Offset = 100
	case r.Limit <= 0:
		r.Limit = 10
	}

	if r.Limit == 0 {
		r.Limit = 10
	}

	data, total, code = model.ListComment(r.Limit, r.Offset)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
		TotalCount: total,
		List:       data,
	})
}
