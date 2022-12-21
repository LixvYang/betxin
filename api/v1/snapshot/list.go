package snapshot

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

type ListResponse struct {
	TotalCount int                          `json:"totalCount"`
	List       []model.MixinNetworkSnapshot `json:"list"`
}

type ListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	UserId string `json:"user_id"`
}

func ListMixinNetworkSnapshots(c *gin.Context) {
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

	var data []model.MixinNetworkSnapshot
	var total int
	var code int
	if r.UserId == "" {
		data, total, code = model.ListMixinNetworkSnapshots(r.Offset, r.Limit)
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR, nil)
			return
		}
	} else {
		data, total, code = model.ListMixinNetworkSnapshotsByUserId(r.UserId, r.Offset, r.Limit)
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR, nil)
			return
		}
	}

	v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
		TotalCount: total,
		List:       data,
	})
}
