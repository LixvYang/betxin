package user

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"

	"github.com/gin-gonic/gin"
)

// @Summary 创建分类
// @Description 创建分类
// @Tags category
// @Accept  json
// @Produce  json
// @Param category body category.model true "创建分类"
// @Success 200 {object} v1.Response "{"code":200,"message":"OK","data":null}"
// @Router /v1/category/add [post]
func CreateUser(c *gin.Context) {
	var r model.User
	if err := c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	code := model.CheckUser(r.IdentityNumber)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_CATENAME_USED, nil)
		return
	}

	if code = model.CreateUser(&r); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	betxinredis.DelKeys(v1.USER_LIST, v1.USER_TOTAL)
	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
