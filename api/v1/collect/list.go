package collect

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

type ListResponse struct {
	TotalCount int             `json:"totalCount"`
	List       []model.Collect `json:"list"`
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
func ListCollects(c *gin.Context) {
	var total int
	var data []model.Collect
	// var err error
	// var collect string
	var code int

	// total, _ = betxinredis.Get(v1.COLLECT_TOTAL).Int()
	// collect, err = betxinredis.Get(v1.COLLECT_LIST).Result()
	// convert.Unmarshal(collect, &data)
	// if err == redis.Nil {
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
	//
	// collect = convert.Marshal(&data)
	// betxinredis.Set(v1.COLLECT_TOTAL, total, v1.REDISEXPIRE)
	// betxinredis.Set(v1.COLLECT_LIST, collect, v1.REDISEXPIRE)

	data, total, code = model.ListCollects(r.Offset, r.Limit)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_LIST_CATEGORY, nil)
		return
	}

	// v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
	// 	TotalCount: total,
	// 	List:       data,
	// })
	// } else if err != nil {
	// 	v1.SendResponse(c, errmsg.ERROR, nil)
	// 	return
	// } else {
	// 	fmt.Println("从redis拿数据")
	v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
		TotalCount: total,
		List:       data,
	})
	// }

}
