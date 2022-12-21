package dailycurrency

import (
	"context"
	"sync"
	"time"

	"github.com/lixvyang/betxin/model"

	"github.com/lixvyang/betxin/internal/utils"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"

	"github.com/fox-one/mixin-sdk-go"
)

var AllCurrency = [...]string{
	utils.PUSD,
	utils.BTC,
	utils.BOX,
	utils.XIN,
	utils.ETH,
	utils.MOB,
	utils.USDC,
	utils.USDT,
	utils.EOS,
	utils.SOL,
	utils.UNI,
	utils.DOGE,
	utils.RUM,
	utils.DOT,
	utils.WOO,
	utils.ZEC,
	utils.LTC,
	utils.SHIB,
	utils.BCH,
	utils.MANA,
	utils.FIL,
	utils.BNB,
	utils.XRP,
	utils.SC,
	utils.MATIC,
	utils.ETC,
	utils.XMR,
	utils.DCR,
	utils.TRX,
	utils.ATOM,
	utils.CKB,
	utils.LINK,
	utils.GTC,
	utils.HNS,
	utils.DASH,
	utils.XLM,
}

func updateRedisCurrency(ctx context.Context) {
	var wg sync.WaitGroup

	for _, currency := range AllCurrency {
		wg.Add(1)
		go func(currency string) {
			defer wg.Done()
			asset, err := mixin.ReadNetworkAsset(ctx, currency)
			if err != nil {
				return
			}
			currencies := &model.Currency{
				AssetId:  asset.AssetID,
				PriceUsd: asset.PriceUSD,
				PriceBtc: asset.PriceBTC,
				ChainId:  asset.ChainID,
				IconUrl:  asset.IconURL,
				Symbol:   asset.Symbol,
			}
			// 有值
			// if model.CheckCurrency(asset.AssetID) != errmsg.SUCCSE {
			model.UpdateCurrency(currencies)
			// } else {
			// 	model.CreateCurrency(currencies)
			// }
			betxinredis.Del(asset.Name + "_" + currency + "_" + "price")
			betxinredis.Set(asset.Name+"_"+currency+"_"+"price", asset.PriceUSD, time.Minute/6)
		}(currency)
	}
	wg.Wait()
}

func DailyCurrency(ctx context.Context) {
	for {
		updateRedisCurrency(ctx)
		time.Sleep(time.Minute / 6)
	}
}
