package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lixvyang/betxin/model"

	"github.com/lixvyang/betxin/internal/utils"
	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinmq "github.com/lixvyang/betxin/pkg/mq"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
	"github.com/lixvyang/betxin/pkg/timewheel"

	"github.com/fox-one/mixin-sdk-go"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

var (
	BETXIN_WORKER = "BETXIN_WORKER"
	mq            = betxinmq.NewMQClient()
)

func init() {
	mq.SetConditions(100)
	ch, err := mq.Subscribe(BETXIN_WORKER)
	if err != nil {
		log.Println("subscribe failed")
		return
	}
	go WorkerSub(ch, mq)
}

func WorkerSub(m <-chan interface{}, c *betxinmq.Client) {
	for {
		val, ok := c.GetPayLoad(m).(mixin.Snapshot)
		if !ok {
			return
		}
		_ = HandlerNewMixinSnapshot(context.Background(), val)
	}
}

type Memo struct {
	Tid      string `json:"tid"`
	YesRatio bool   `json:"yes_ratio"`
	NoRatio  bool   `json:"no_ratio"`
}

type Stats struct {
	preCreatedAt time.Time
}

func (s *Stats) getPrevSnapshotCreatedAt() time.Time {
	return s.preCreatedAt
}

func (s *Stats) updatePrevSnapshotCreatedAt(time time.Time) {
	s.preCreatedAt = time
}

func getTopSnapshotCreatedAt(client *mixin.Client, c context.Context) (time.Time, error) {
	snapshots, err := client.ReadSnapshots(c, "", time.Now(), "", 1)
	if err != nil {
		return time.Now(), err
	}
	if len(snapshots) == 0 {
		return time.Now(), nil
	}
	return snapshots[0].CreatedAt, nil
}

func getTopHundredCreated(client *mixin.Client, c context.Context) ([]mixin.Snapshot, error) {
	snapshots, err := client.ReadSnapshots(c, "", time.Now(), "", 50)
	if err != nil {
		return nil, err
	}
	var snapshot []mixin.Snapshot
	for i := 0; i < len(snapshots); i++ {
		snapshot = append(snapshot, *snapshots[i])
	}
	return snapshot, nil
}

func sendTopCreatedAtToChannel(ctx context.Context, stats *Stats) {
	preCreatedAt := stats.getPrevSnapshotCreatedAt()
	snapshots, err := getTopHundredCreated(MixinClient, ctx)
	if err != nil {
		log.Printf("getTopHundredCreated error")
		return
	}
	var wg sync.WaitGroup
	for _, snapshot := range snapshots {
		wg.Add(1)
		if snapshot.CreatedAt.After(preCreatedAt) {
			stats.updatePrevSnapshotCreatedAt(snapshot.CreatedAt)
			if snapshot.Amount.Cmp(decimal.NewFromInt(0)) == 1 && snapshot.Type == "transfer" {
				go func(snapshot mixin.Snapshot) {
					defer wg.Done()
					// _ = HandlerNewMixinSnapshot(ctx, client, snapshot)
					fmt.Println("来新账单啦")
					mq.Publish(BETXIN_WORKER, snapshot)
				}(snapshot)
			}
		}
	}
	wg.Wait()
}

func Worker(ctx context.Context) error {
	createdAt, err := getTopSnapshotCreatedAt(MixinClient, ctx)
	if err != nil {
		return err
	}
	stats := &Stats{createdAt}
	timewheel.Every(time.Second*2, func() {
		go sendTopCreatedAtToChannel(ctx, stats)
	})
	return nil
}

// func Worker(ctx context.Context, client *mixin.Client) error {
// 	createdAt, err := getTopSnapshotCreatedAt(client, ctx)
// 	if err != nil {
// 		return err
// 	}
// 	stats := &Stats{createdAt}
// 	gocron.Every(2).Second().Do(sendTopCreatedAtToChannel, ctx, stats, client)
// 	<-gocron.Start()
// 	return nil
// }

func HandlerNewMixinSnapshot(ctx context.Context, snapshot mixin.Snapshot) error {
	fmt.Println("开始处理HandlerNewMixinSnapshot")
	if snapshot.Memo == "" {
		log.Println("memo 为空退出")
		return nil
	}

	r := model.MixinOrder{
		Type:       snapshot.Type,
		AssetId:    snapshot.AssetID,
		Amount:     snapshot.Amount,
		TraceId:    snapshot.TraceID,
		Memo:       snapshot.Memo,
		SnapshotId: snapshot.SnapshotID,
	}

	if code := model.CreateMixinOrder(&r); code != errmsg.SUCCSE {
		log.Println("创建CreateMixinOrder错误")
		return errors.New("")
	}
	// 用户传过来的memo是经过base64加密的  yes或no  再加上trace_id 的json
	///  memo  traceId:不应该是随机id 应该是把userid和买的topic id yesorno放在一起
	tx := &mixin.RawTransaction{}
	if snapshot.AssetID != utils.PUSD {
		tx = SwapOrderToPusd(ctx, snapshot.Amount, snapshot.AssetID, snapshot)
	} else {
		tx.Amount = snapshot.Amount.String()
		tx.AssetID = snapshot.AssetID
	}
	amount, err := decimal.NewFromString(tx.Amount)
	if err != nil {
		log.Println(err)
		log.Println("计算失败")
	}

	// 用户投入的总价格
	userTotalPrice, err := CalculateTotalPriceByAssetId(ctx, tx.AssetID, amount.Abs())
	if err != nil {
		log.Println("计算失败")
	}

	memoMsg, err := base64.StdEncoding.DecodeString(snapshot.Memo)
	if err != nil {
		return errors.New("解码memo失败")
	}

	memo := &Memo{}
	if err := json.Unmarshal(memoMsg, &memo); err != nil {
		return errors.New("解构memo失败")
	}

	var data model.UserToTopic
	var selectWin string
	data.UserId = snapshot.OpponentID
	if memo.YesRatio {
		selectWin = "yes_win"
		data.YesRatioPrice = userTotalPrice
	} else {
		selectWin = "no_win"
		data.NoRatioPrice = userTotalPrice
	}
	data.Tid = memo.Tid

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

	betxinredis.BatchDel("topic")

	return nil
}

func SwapOrderToPusd(ctx context.Context, Amount decimal.Decimal, InputAssetId string, snapshot mixin.Snapshot) *mixin.RawTransaction {
	tx, err := TransactionWithRetry(ctx, MixinClient, Amount, InputAssetId)
	if err != nil {
		uuid := uuid.NewV4()
		model.CreateSendBack(&model.SendBack{TraceId: uuid.String()})
		err := TransferReturnWithRetry(ctx, MixinClient, uuid.String(), InputAssetId, snapshot.OpponentID, snapshot.Amount, "Swap 失败")
		switch {
		case mixin.IsErrorCodes(err, mixin.InsufficientBalance):
			log.Println("insufficient balance")
		default:
			log.Printf("transfer: %v", err)
		}
	}
	amount, _ := decimal.NewFromString(tx.Amount)
	data := &model.SwapOrder{
		Type:       tx.Type,
		SnapshotId: tx.SnapshotID,
		AssetID:    tx.AssetID,
		Amount:     amount,
		TraceId:    tx.TraceID,
		Memo:       tx.Memo,
		State:      tx.State,
	}

	if code := model.CreateSwapOrder(data); code != errmsg.SUCCSE {
		return nil
	}
	return tx
}
