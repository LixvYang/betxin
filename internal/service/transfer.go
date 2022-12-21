package service

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/lixvyang/betxin/model"

	"github.com/lixvyang/betxin/internal/utils"
	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/avast/retry-go"
	fswap "github.com/fox-one/4swap-sdk-go"
	"github.com/fox-one/4swap-sdk-go/mtg"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

func TransferReturnWithRetry(ctx context.Context, client *mixin.Client, TraceId string, AssetID string, OpponentID string, Amount decimal.Decimal, Memo string) error {
	return retry.Do(
		func() error {
			err := TransferReturn(ctx, MixinClient, TraceId, AssetID, OpponentID, Amount, Memo)
			if err != nil {
				return err
			}
			return nil
		},
		retry.Delay(time.Second*2),
	)
}

func TransferReturn(ctx context.Context, client *mixin.Client, TraceId string, AssetID string, OpponentID string, Amount decimal.Decimal, Memo string) error {
	transferInput := &mixin.TransferInput{
		AssetID:    AssetID,
		OpponentID: OpponentID,
		Amount:     Amount,
		TraceID:    TraceId,
		Memo:       Memo,
	}

	tx, err := client.Transfer(ctx, transferInput, utils.Pin)
	if err != nil {
		return err
	}

	data := &model.SendBack{
		Type:       tx.Type,
		SnapshotId: tx.SnapshotID,
		OpponentID: tx.OpponentID,
		AssetID:    tx.AssetID,
		Amount:     tx.Amount,
		Memo:       tx.Memo,
	}

	if code := model.UpdateSendBack(TraceId, data); code != errmsg.SUCCSE {
		log.Println("error to update mixinnetwork snapshot")
	}
	return nil
}

// count: 尝试发送次数
// interval: 第一次失败重试间隔时间(后续间隔翻倍)
func TransferWithRetry(ctx context.Context, client *mixin.Client, TraceId string, AssetID string, OpponentID string, Amount decimal.Decimal, Memo string) error {
	return retry.Do(
		func() error {
			err := Transfer(ctx, MixinClient, TraceId, AssetID, OpponentID, Amount, Memo)
			if err != nil {
				return err
			}
			return nil
		},
		retry.Delay(time.Second*2),
	)
}

func Transfer(ctx context.Context, client *mixin.Client, TraceId string, AssetID string, OpponentID string, Amount decimal.Decimal, Memo string) error {
	transferInput := &mixin.TransferInput{
		AssetID:    AssetID,
		OpponentID: OpponentID,
		Amount:     Amount,
		TraceID:    TraceId,
		Memo:       Memo,
	}

	tx, err := client.Transfer(ctx, transferInput, utils.Pin)
	if err != nil {
		log.Println(err)
		log.Println("转账失败")
		return err
	}

	data := &model.MixinNetworkSnapshot{
		SnapshotId:     tx.SnapshotID,
		AssetId:        tx.AssetID,
		OpponentID:     tx.OpponentID,
		Amount:         tx.Amount,
		Memo:           tx.Memo,
		Type:           tx.Type,
		OpeningBalance: tx.OpeningBalance,
		ClosingBalance: tx.ClosingBalance,
	}

	if code := model.UpdateMixinNetworkSnapshot(TraceId, data); code != errmsg.SUCCSE {
		log.Println("error to update mixinnetwork snapshot")
	}
	return nil
}

// count: 尝试发送次数
// interval: 第一次失败重试间隔时间(后续间隔翻倍)
func TransactionWithRetry(ctx context.Context, client *mixin.Client, Amount decimal.Decimal, InputAssetID string) (*mixin.RawTransaction, error) {
	var rawTransaction *mixin.RawTransaction
	err := retry.Do(
		func() error {
			tx, err := Transaction(ctx, client, Amount, InputAssetID)
			if err != nil {
				return err
			}
			rawTransaction = tx
			return nil
		},
		// 4秒后重试
		retry.Delay(time.Second*4),
	)
	return rawTransaction, err
}

// 输入数量和输入资产id 输出交易单
func Transaction(ctx context.Context, client *mixin.Client, Amount decimal.Decimal, InputAssetID string) (*mixin.RawTransaction, error) {
	fswap.UseEndpoint(fswap.MtgEndpoint)
	// read the mtg group
	// the group information would change frequently
	// it's recommended to save it for later use
	group, err := fswap.ReadGroup(ctx)
	if err != nil {
		log.Println("读取组失败")
		return nil, err
	}
	pairs, _ := fswap.ListPairs(ctx)
	sort.Slice(pairs, func(i, j int) bool {
		aLiquidity := pairs[i].BaseValue.Add(pairs[i].QuoteValue)
		bLiquidity := pairs[j].BaseValue.Add(pairs[j].QuoteValue)
		return aLiquidity.GreaterThan(bLiquidity)
	})

	preOrder, err := fswap.Route(pairs, InputAssetID, utils.PUSD, Amount)
	if err != nil {
		log.Println("路由失败")
		return nil, err
	}

	followID, _ := uuid.NewV4()
	action := mtg.SwapAction(
		client.ClientID,
		followID.String(),
		utils.PUSD,
		preOrder.Routes,
		decimal.NewFromFloat(0.00000001),
	)

	// 生成 memo
	memo, err := action.Encode(group.PublicKey)
	if err != nil {
		log.Println("生成memo失败")
		return nil, err
	}

	tx, err := client.Transaction(ctx, &mixin.TransferInput{
		AssetID: InputAssetID,
		Amount:  Amount,
		TraceID: mixin.RandomTraceID(),
		Memo:    memo,
		OpponentMultisig: struct {
			Receivers []string `json:"receivers,omitempty"`
			Threshold uint8    `json:"threshold,omitempty"`
		}{
			Receivers: group.Members,
			Threshold: uint8(group.Threshold),
		},
	}, utils.Pin)
	if err != nil {
		log.Println("生成交易失败")
		return nil, err
	}
	return tx, nil
}

// 根据输入的资产id和资产数目计算出资产总价格
func CalculateTotalPriceByAssetId(ctx context.Context, AssedId string, amount decimal.Decimal) (decimal.Decimal, error) {
	decimal.DivisionPrecision = 8 // 保留两位小数，如有更多位，则进行四舍五入保留两位小数
	asset, code := model.GetCurrencyById(AssedId)
	if code != errmsg.SUCCSE {
		asset, err := mixin.ReadNetworkAsset(ctx, AssedId)
		if err != nil {
			return asset.PriceUSD.Mul(amount), nil
		}
		return decimal.NewFromFloat(0), err
	}
	return asset.PriceUsd.Mul(amount), nil
}

// 根据输入的资产symbol和资产数目计算出资产总价格
func CalculateTotalPriceBySymbol(ctx context.Context, Symbol string, amount decimal.Decimal) (decimal.Decimal, error) {
	decimal.DivisionPrecision = 8 // 保留两位小数，如有更多位，则进行四舍五入保留两位小数
	asset, code := model.GetCurrencyBySymbol(Symbol)
	if code != errmsg.SUCCSE {
		asset, err := mixin.ReadNetworkAsset(ctx, asset.AssetId)
		if err != nil {
			return asset.PriceUSD.Mul(amount), nil
		}
		return decimal.NewFromFloat(0), err
	}
	return asset.PriceUsd.Mul(amount), nil
}
