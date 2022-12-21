package currency

import (
	"time"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
	"github.com/lixvyang/betxin/pkg/convert"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ListResponse struct {
	TotalCount int              `json:"totalCount"`
	List       []model.Currency `json:"list"`
}

func ListCurrencies(c *gin.Context) {
	var data []model.Currency
	var total int
	var code int
	var err error
	var currencies string
	total, _ = betxinredis.Get(v1.CURRENCY_TOTAL).Int()
	currencies, err = betxinredis.Get(v1.CURRENCY_LIST).Result()
	convert.Unmarshal(currencies, &data)
	if err == redis.Nil {
		data, total, code = model.ListCurrencies()
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR, nil)
			return
		}

		currencies = convert.Marshal(&data)
		betxinredis.Set(v1.CURRENCY_TOTAL, total, time.Minute/6)
		betxinredis.Set(v1.CURRENCY_LIST, currencies, time.Minute/6)

		v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
			TotalCount: total,
			List:       data,
		})
	} else if err != nil {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	} else {
		v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
			TotalCount: total,
			List:       data,
		})
	}
}
