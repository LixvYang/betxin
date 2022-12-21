package user

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
	"github.com/lixvyang/betxin/pkg/convert"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ListResponse struct {
	TotalCount int          `json:"totalCount"`
	List       []model.User `json:"list"`
}

type ListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// @Summary 获得category列表
// @Description 获取category列表
// @Tags category
// @Accept  json
// @Produce  json
// @Param   offset      query    int     true     "Offset"
// @Param   limit      query    int     true      "Limit"
// @Success 200 {object} category.ListResponse "{"code":200,"message":"OK","data":{"totalCount":1,"list":[]}"
// @Router /v1/category [get]
func ListUser(c *gin.Context) {
	var err error
	var data []model.User
	var code int
	var total int
	var users string

	total, _ = betxinredis.Get(v1.USER_TOTAL).Int()
	users, err = betxinredis.Get(v1.USER_LIST).Result()
	convert.Unmarshal(users, &data)
	if err == redis.Nil {
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

		data, total, code = model.ListUser(r.Offset, r.Limit)
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR, nil)
			return
		}
		//
		users = convert.Marshal(&data)
		betxinredis.Set(v1.USER_TOTAL, total, v1.REDISEXPIRE)
		betxinredis.Set(v1.USER_LIST, users, v1.REDISEXPIRE)
		v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
			TotalCount: total,
			List:       data,
		})
	} else if err != nil {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	} else {
		v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
			TotalCount: total,
			List:       data,
		})
	}
}
