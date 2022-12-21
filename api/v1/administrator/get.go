package administrator

import (
	"log"
	"strconv"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// @Summary 用 administrator id获取分类
// @Description 用administrator id获取分类信息
// @Tags category
// @Accept  json
// @Produce  json
// @Param id path string true "Id"
// @Success 200 {object} model.Category "{"code":200,"message":"OK","data":{}}"
// @Router /v1/category/{id} [get]
func GetAdministratorInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("获取参数id错误")
	}
	category, code := model.GetAdministratorById(id)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	v1.SendResponse(c, errmsg.SUCCSE, category)
}
