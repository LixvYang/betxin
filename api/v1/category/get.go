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

// @Summary 用分类id获取分类
// @Description 用分类id获取分类信息
// @Tags category
// @Accept  json
// @Produce  json
// @Param id path string true "Id"
// @Success 200 {object} model.Category "{"code":200,"message":"OK","data":{}}"
// @Router /v1/category/{id} [get]
func GetCategoryInfo(c *gin.Context) {
	id := c.Param("id")
	var data model.Category
	var code int
	var category string
	var err error

	category, err = betxinredis.Get(v1.CATEGORY_GET + id).Result()
	convert.Unmarshal(category, &data)
	if err == redis.Nil {

		data, code = model.GetCategoryById(convert.StrToNum(id))
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR, nil)
			return
		}
		// 将数据存入redis
		category = convert.Marshal(&data)
		betxinredis.Set(v1.CATEGORY_GET+id, category, v1.REDISEXPIRE)

		v1.SendResponse(c, errmsg.SUCCSE, model.Category{
			Id:           data.Id,
			CategoryName: data.CategoryName,
		})
	} else if err != nil {
		v1.SendResponse(c, errmsg.ERROR, nil)
	} else {
		v1.SendResponse(c, errmsg.SUCCSE, model.Category{
			Id:           data.Id,
			CategoryName: data.CategoryName,
		})
	}
}
