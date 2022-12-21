package service

import (
	"context"
	"errors"
	"log"

	"github.com/lixvyang/betxin/model"

	"github.com/lixvyang/betxin/internal/utils"
	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/shopspring/decimal"
)

// 退还用户资金  5%的手续费
func RefundUserToTopic(yesFee decimal.Decimal, noFee decimal.Decimal, userToTopic model.UserToTopic) error {
	// 退yes
	if userToTopic.YesRatioPrice.GreaterThan(decimal.NewFromFloat(0)) {
		traceId := mixin.RandomTraceID()
		refund := model.Refund{TraceId: traceId}
		if code := model.CreateRefund(&refund); code != errmsg.SUCCSE {
			return errors.New("create refund error")
		}

		err := TransferReturnWithRetry(context.Background(), MixinClient, traceId, utils.PUSD, userToTopic.UserId, userToTopic.YesRatioPrice, "Refund YES Money")
		if err != nil {
			return err
		}
		data := model.Refund{
			UserId:      userToTopic.UserId,
			Tid:         userToTopic.Tid,
			AssetId:     utils.PUSD,
			RefundPrice: userToTopic.YesRatioPrice,
			Memo:        "Refund YES Money",
		}
		if code := model.UpdateRefund(traceId, &data); code != errmsg.SUCCSE {
			log.Println("UpdateRefund error!")
		}
		if code := model.RefundTopicTotalPrice(&data, "yes", yesFee); code != errmsg.SUCCSE {
			log.Println("RefundTopicTotalPrice error!")
		}
	}

	if userToTopic.NoRatioPrice.GreaterThan(decimal.NewFromFloat(0)) {
		traceId := mixin.RandomTraceID()
		refund := model.Refund{TraceId: traceId}
		if code := model.CreateRefund(&refund); code != errmsg.SUCCSE {
			return errors.New("create refund error")
		}
		err := TransferReturnWithRetry(context.Background(), MixinClient, traceId, utils.PUSD, userToTopic.UserId, userToTopic.NoRatioPrice, "Refund No Money")
		if err != nil {
			return err
		}
		data := model.Refund{
			UserId:      userToTopic.UserId,
			Tid:         userToTopic.Tid,
			AssetId:     utils.PUSD,
			RefundPrice: userToTopic.NoRatioPrice,
			Memo:        "Refund YES Money",
		}
		if code := model.UpdateRefund(traceId, &data); code != errmsg.SUCCSE {
			log.Println("UpdateRefund error!")
		}
		if code := model.RefundTopicTotalPrice(&data, "no", noFee); code != errmsg.SUCCSE {
			log.Println("RefundTopicTotalPrice error!")
		}
	}

	return nil
}

// true -> delete usertopic
func CheckUserToTopicZero(uid, tid string) {
	userToTopic, code := model.GetUserToTopic(uid, tid)
	if code != errmsg.SUCCSE {
		return
	}
	zeroPrice := decimal.NewFromFloat(0)
	if userToTopic.YesRatioPrice == zeroPrice && userToTopic.NoRatioPrice == zeroPrice {
		_ = model.DeleteUserToTopic(uid, tid)
		return
	}
}
