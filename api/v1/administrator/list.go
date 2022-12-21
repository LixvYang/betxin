package administrator

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

type ListResponse struct {
	TotalCount int                   `json:"totalCount"`
	List       []model.Administrator `json:"list"`
}

type ListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// @Summary 获得 administrator 列表
// @Description 获取administrator列表
// @Tags administrator
// @Accept  json
// @Produce  json
// @Param   offset      query    int     true     "Offset"
// @Param   limit      query    int     true      "Limit"
// @Success 200 {object} administrator.ListResponse "{"code":200,"message":"OK","data":{"totalCount":1,"list":[]}"
// @Router /v1/administrator [get]
func ListAdministrators(c *gin.Context) {
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

	data, total, code := model.ListAdministrators(r.Offset, r.Limit)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_LIST_USER, nil)
		return
	}
	v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
		TotalCount: total,
		List:       data,
	})
}
