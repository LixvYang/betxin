package usertotopic

import (
	"context"

	"github.com/lixvyang/betxin/internal/service"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils"
	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type RefundRequest struct {
	UserId        string          `json:"user_id"`
	Tid           string          `json:"tid"`
	YesRatioPrice decimal.Decimal `json:"yes_ratio_price"`
	NoRatioPrice  decimal.Decimal `json:"no_ratio_price"`
}

// 每次扣除总资金的   5%
func RefundUserToTopic(c *gin.Context) {
	var r RefundRequest
	var yesFee decimal.Decimal
	var noFee decimal.Decimal
	if err := c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}
	usertotopic, code := model.GetUserToTopic(r.UserId, r.Tid)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	// 检查请求参数
	if r.NoRatioPrice.GreaterThan(usertotopic.NoRatioPrice.Mul(decimal.NewFromFloat(0.95))) || r.YesRatioPrice.GreaterThan(usertotopic.YesRatioPrice.Mul(decimal.NewFromFloat(0.95))) {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	data := model.UserToTopic{
		UserId:        r.UserId,
		Tid:           r.Tid,
		YesRatioPrice: r.YesRatioPrice,
		NoRatioPrice:  r.NoRatioPrice,
	}

	if yesFee, noFee, code = model.RefundUserToTopic(&data); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	if err := service.RefundUserToTopic(yesFee, noFee, data); err != nil {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	// 扣款
	if r.YesRatioPrice.GreaterThan(decimal.NewFromFloat(0)) {
		_ = service.Transfer(context.Background(), service.MixinClient, mixin.RandomTraceID(), utils.PUSD, "6a87e67f-02fb-47cf-b31f-32a13dd5b3d9", usertotopic.YesRatioPrice.Mul(decimal.NewFromFloat(0.05)), "退款手续费")
	}

	if r.NoRatioPrice.GreaterThan(decimal.NewFromFloat(0)) {
		_ = service.Transfer(context.Background(), service.MixinClient, mixin.RandomTraceID(), utils.PUSD, "6a87e67f-02fb-47cf-b31f-32a13dd5b3d9", usertotopic.NoRatioPrice.Mul(decimal.NewFromFloat(0.05)), "退款手续费")
	}

	service.CheckUserToTopicZero(r.UserId, r.Tid)

	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
