package administrator

import (
	"strconv"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// @Summary 根据id删除管理员
// @Description 根据id删除管理员
// @Tags	administrator
// @Accept  json
// @Produce  json
// @Param id path int true "管理员的数据库id"
// @Success 200 {object} v1.Response "{"code":200,"message":"OK","data":null}"
// @Router /v1/administrator/{id} [delete]
func DeleteAdministrator(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DeleteAdministrator(id)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_DELETE_CATENAME, nil)
		return
	}
	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
