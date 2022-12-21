package mixpayorder

import (
	"fmt"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Tid      string `json:"tid"`
	YesRatio bool   `json:"yes_ratio"`
	NoRatio  bool   `json:"no_ratio"`
	OrderId  string `json:"orderId"`
	PayeeId  string `json:"payeeId"`
}

// 在用户点击时创建
func CreateMixinpayOrder(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("userId")
	userId := fmt.Sprintf("%v", user)
	var r CreateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	data := &model.MixpayOrder{
		Uid:      userId,
		Tid:      r.Tid,
		YesRatio: r.YesRatio,
		NoRatio:  r.NoRatio,
		OrderId:  r.OrderId,
		PayeeId:  r.PayeeId,
	}

	if code := model.CreateMixpayOrder(data); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, r)
}
