package mixpayorder

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

type ListResponse struct {
	TotalCount int                 `json:"totalCount"`
	List       []model.MixpayOrder `json:"list"`
}

type ListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func ListMixpayOrder(c *gin.Context) {
	var total int
	var data []model.MixpayOrder
	var code int
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

	data, total, code = model.ListMixpayOrder(r.Limit, r.Offset)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
		TotalCount: total,
		List:       data,
	})
}
