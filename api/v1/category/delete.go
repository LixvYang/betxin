package category

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
	"github.com/lixvyang/betxin/pkg/convert"

	"github.com/gin-gonic/gin"
)

// @Summary 根据id删除分类
// @Description 根据id删除分类
// @Tags	category
// @Accept  json
// @Produce  json
// @Param id path int true "分类的数据库id"
// @Success 200 {object} v1.Response "{"code":200,"message":"OK","data":null}"
// @Router /v1/category/{id} [delete]
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	code := model.DeleteCategory(convert.StrToNum(id))
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_DELETE_CATENAME, nil)
		return
	}

	betxinredis.DelKeys(v1.CATEGORY_GET+id, v1.CATEGORY_LIST, v1.CATEGORY_TOTAL)
	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
