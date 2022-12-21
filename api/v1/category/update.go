package category

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
	"github.com/lixvyang/betxin/pkg/convert"

	"github.com/gin-gonic/gin"
)

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var cate *model.Category
	if err := c.ShouldBindJSON(&cate); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
	}

	code := model.CheckCategory(cate.CategoryName)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_CATENAME_USED, nil)
		return
	}

	code = model.UpdateCate(convert.StrToNum(id), cate.CategoryName)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_UPDATE_CATENAME, nil)
		return
	}

	// Delete redis store
	betxinredis.DelKeys(v1.CATEGORY_GET+id, v1.CATEGORY_LIST, v1.CATEGORY_TOTAL)

	v1.SendResponse(c, errmsg.SUCCSE, cate.CategoryName)
}
