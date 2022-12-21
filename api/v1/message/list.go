package message

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

type ListResponse struct {
	TotalCount int                  `json:"totalCount"`
	List       []model.MixinMessage `json:"list"`
}

type ListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	UserId string `json:"user_id"`
}

func ListMessages(c *gin.Context) {
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

	data, total, code := model.ListMixinMessageByUserId(r.UserId, r.Offset, r.Limit)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_LIST_CATEGORY, nil)
		return
	}
	v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
		TotalCount: total,
		List:       data,
	})
}
