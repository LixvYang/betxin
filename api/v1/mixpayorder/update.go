package mixpayorder

import (
	"net/http"

	"github.com/lixvyang/betxin/internal/service/mixpay"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

type UpdateRequest struct {
	OrderId string `json:"orderId"`
	PayeeId string `json:"payeeId"`
	TraceId string `json:"traceId"`
}

func UpdateMixpayOrder(c *gin.Context) {
	var mixpayorder model.MixpayOrder
	var u UpdateRequest
	if err := c.ShouldBindJSON(&u); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	mixpayorder = model.MixpayOrder{
		OrderId: u.OrderId,
		PayeeId: u.PayeeId,
		TraceId: u.TraceId,
	}

	if code := model.UpdateMixpayOrder(&mixpayorder); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	mixpayRes, err := mixpay.GetMixpayResult(u.OrderId, u.PayeeId)
	if err != nil {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	// 查询Mixpay支付信息　比如
	mixpayOrder, code := model.GetMixpayOrder(mixpayorder.TraceId)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	if err := mixpay.Worker(mixpayOrder, mixpayRes); err != nil {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "SUCCESS",
	})
}
