package category

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
	TotalCount int              `json:"totalCount"`
	List       []model.Category `json:"list"`
}

type ListRequest struct {
	CategoryName string `json:"category_name"`
	Offset       int    `json:"offset"`
	Limit        int    `json:"limit"`
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
func ListCategories(c *gin.Context) {
	var categoryies string
	var total int
	var code int
	var err error

	var data []model.Category

	total, _ = betxinredis.Get(v1.CATEGORY_TOTAL).Int()
	categoryies, err = betxinredis.Get(v1.CATEGORY_LIST).Result()
	convert.Unmarshal(categoryies, &data)
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

		if r.CategoryName != "" {
			data, total, code = model.SearchCategory(r.CategoryName, r.Offset, r.Limit)
			if code != errmsg.SUCCSE {
				v1.SendResponse(c, errmsg.ERROR_LIST_CATEGORY, nil)
				return
			}
		} else {
			data, total, code = model.ListCategories(r.Offset, r.Limit)
			if code != errmsg.SUCCSE {
				v1.SendResponse(c, errmsg.ERROR_LIST_CATEGORY, nil)
				return
			}
		}

		categoryies = convert.Marshal(&data)
		betxinredis.Set(v1.CATEGORY_TOTAL, total, v1.REDISEXPIRE)
		betxinredis.Set(v1.CATEGORY_LIST, categoryies, v1.REDISEXPIRE)
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
