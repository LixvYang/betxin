package mixpay

import (
	"context"
	"log"
	"time"

	"github.com/lixvyang/betxin/internal/service"

	"github.com/lixvyang/betxin/model"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/shopspring/decimal"
)

func Worker(mixpayorder model.MixpayOrder, mixpayRes MixpayResult) error {
	// Mixpay支付成功
	// 将用户信息加入到user to topic 里面
	var data model.UserToTopic
	var selectWin string
	var err error
	var userTotalPrice decimal.Decimal
	data.Tid = mixpayorder.Tid
	data.UserId = mixpayorder.Uid

	if mixpayorder.YesRatio {
		selectWin = "yes_win"
		payAmount, err := decimal.NewFromString(mixpayRes.Data.PaymentAmount)
		if err != nil {
			log.Println("支付价转换失败")
		}
		userTotalPrice, err = service.CalculateTotalPriceBySymbol(context.Background(), mixpayRes.Data.PaymentSymbol, payAmount)
		if err != nil {
			log.Println("计算价格转换失败")
		}
		data.YesRatioPrice = userTotalPrice
	} else {
		selectWin = "no_win"
		payAmount, err := decimal.NewFromString(mixpayRes.Data.PaymentAmount)
		if err != nil {
			log.Println("支付价转换失败")
		}
		userTotalPrice, err = service.CalculateTotalPriceBySymbol(context.Background(), mixpayRes.Data.PaymentSymbol, payAmount)
		if err != nil {
			log.Println("计算价格转换失败")
		}
		data.NoRatioPrice = userTotalPrice
	}

	// 已经买过了
	if code := model.CheckUserToTopic(data.UserId, data.Tid); code != errmsg.ERROR {
		code = model.UpdateUserToTopic(&data)
		if code != errmsg.SUCCSE {
			log.Println("CreateUserToTopic错误")
			return err
		}
	} else {
		code = model.CreateUserToTopic(&data)
		if code != errmsg.SUCCSE {
			log.Println("CreateUserToTopic错误")
			return err
		}
	}

	if code := model.UpdateTopicTotalPrice(data.Tid, selectWin, userTotalPrice); code != errmsg.SUCCSE {
		log.Println("UpdateTopicTotalPrice错误")
		return err
	}
	return nil
}

// delete day out
func Delete() {
	for {
		// TODO
		time.Sleep(time.Duration(time.Now().Day()))
	}
}
