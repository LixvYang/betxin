package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lixvyang/betxin/model"

	"github.com/lixvyang/betxin/internal/utils"
	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/fox-one/mixin-sdk-go"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type UserBounse struct {
	percentage decimal.Decimal
	TraceId    string
	UserId     string
	Memo       string
}

func EndOfTopic(c context.Context, tid string, win string) {
	var code int
	var userTotopics []model.UserToTopic
	var totalPrice decimal.Decimal
	var winTotalPrice decimal.Decimal
	var userBounses []UserBounse
	data := &model.Bonuse{}

	totalPrice, code = model.GetTopicTotalPrice(tid)
	if code != errmsg.SUCCSE {
		return
	}

	// 收取5%的金钱
	_ = TransferWithRetry(context.Background(), MixinClient, mixin.RandomTraceID(), utils.PUSD, "6a87e67f-02fb-47cf-b31f-32a13dd5b3d9", totalPrice.Mul(decimal.NewFromFloat(0.05)), "话题收取手续费")
	_ = Transfer(context.Background(), MixinClient, mixin.RandomTraceID(), utils.PUSD, "6a87e67f-02fb-47cf-b31f-32a13dd5b3d9", totalPrice.Mul(decimal.NewFromFloat(0.05)), "话题手续费")
	totalPrice = totalPrice.Mul(decimal.NewFromFloat(0.95))

	userTotopics, _, code = model.ListUserToTopicsWin(tid, win)
	if code != errmsg.SUCCSE {
		log.Println("列出赢了的用户失败")
		return
	}

	for _, userToTopic := range userTotopics {
		data.Tid = tid
		data.AssetId = utils.PUSD
		data.Memo = fmt.Sprintln("bonuse from betxin" + userToTopic.Topic.Intro)
		data.UserId = userToTopic.UserId
		data.TraceId = uuid.NewV4().String()
		winTotalPrice, code = model.SearchTopicWinTopic(tid, win)
		if code != errmsg.SUCCSE {
			log.Println("计算赢了总价格失败")
		}

		if win == "yes_win" {
			// 占赢了的百分比
			percentage := userToTopic.YesRatioPrice.Div(winTotalPrice)
			data.Amount = percentage.Mul(totalPrice)
			userBounses = append(userBounses, UserBounse{percentage: percentage, UserId: data.UserId, TraceId: data.TraceId, Memo: data.Memo})
		} else {
			percentage := userToTopic.NoRatioPrice.Div(winTotalPrice)
			data.Amount = percentage.Mul(totalPrice)
			userBounses = append(userBounses, UserBounse{percentage: percentage, UserId: data.UserId, TraceId: data.TraceId, Memo: data.Memo})
		}

		if code = model.CreateBonuse(data); code != errmsg.SUCCSE {
			log.Println("创建奖金出错")
			return
		}
		snapShot := &model.MixinNetworkSnapshot{
			TraceId: data.TraceId,
		}
		model.CreateMixinNetworkSnapshot(snapShot)
	}

	// send for users
	for _, userBounse := range userBounses {
		_ = TransferWithRetry(c, MixinClient, userBounse.TraceId, utils.PUSD, userBounse.UserId, userBounse.percentage.Mul(totalPrice), userBounse.Memo)
		time.Sleep(1 * time.Second)
	}
}
