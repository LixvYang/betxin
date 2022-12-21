package swaporder

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

type ListResponse struct {
	TotalCount int               `json:"totalCount"`
	List       []model.SwapOrder `json:"list"`
}

type ListRequest struct {
	AssetId string `json:"asset_id"`
	TraceId string `json:"trace_id"`
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
}

func ListSwapOrder(c *gin.Context) {
	var data []model.SwapOrder
	var total int
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

	if r.AssetId != "" && r.TraceId != "" {
		data, total, code = model.ListSwapOrders(r.Offset, r.Limit, "assed_id = ? AND trace_id = ?", r.AssetId, r.TraceId)
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR_LIST_CATEGORY, nil)
			return
		}
	} else if r.AssetId != "" && r.TraceId == "" {
		data, total, code = model.ListSwapOrders(r.Offset, r.Limit, "assed_id = ?", r.AssetId)
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR_LIST_CATEGORY, nil)
			return
		}
	} else if r.AssetId == "" && r.TraceId == "" {
		data, total, code = model.ListSwapOrdersNoLimit(r.Offset, r.Limit)
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR_LIST_CATEGORY, nil)
			return
		}
	}

	v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
		TotalCount: total,
		List:       data,
	})
}

func ListSwapOrderNoLimit(c *gin.Context) {
	var data []model.SwapOrder
	var total int
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

	data, total, code = model.ListSwapOrdersNoLimit(r.Offset, r.Limit)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_LIST_CATEGORY, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
		TotalCount: total,
		List:       data,
	})
}
