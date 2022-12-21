package topic

import (
	"fmt"
	"strconv"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
	"github.com/lixvyang/betxin/pkg/convert"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ListResponse struct {
	TotalCount int           `json:"totalCount"`
	List       []model.Topic `json:"list"`
}

type ListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Intro  string `json:"intro"`
	Cid    string `json:"cid"`
}

func ListTopics(c *gin.Context) {
	var r ListRequest
	var topics string
	var data []model.Topic
	var total int
	var code int
	var err error

	if err = c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	var totalRedis string = v1.TOPIC_TOTAL + "_" + strconv.Itoa(r.Limit) + "_" + strconv.Itoa(r.Offset)
	var topicsRedis string = v1.TOPIC_LIST + "_" + strconv.Itoa(r.Limit) + "_" + strconv.Itoa(r.Offset)

	total, _ = betxinredis.Get(totalRedis).Int()
	topics, err = betxinredis.Get(topicsRedis).Result()
	convert.Unmarshal(topics, &data)
	if err == redis.Nil {
		data, total, code = model.ListTopics(r.Offset, r.Limit)
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR_LIST_TOPIC, nil)
			return
		}
		topics = convert.Marshal(&data)
		betxinredis.Set(totalRedis, total, v1.REDISEXPIRE)
		betxinredis.Set(topicsRedis, topics, v1.REDISEXPIRE)

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

// GetTopicByCid 通过种类id获取信息
func GetTopicByCid(c *gin.Context) {
	var topics string
	var data []model.Topic
	var total int
	var code int
	var err error

	var r ListRequest
	if err = c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	cid := c.Param("cid")
	var totalRedis = v1.TOPIC_LIST_FROMCATE_TOTAL + cid
	var topicsRedis = v1.TOPIC_LIST_FROMCATE + cid + "_" + strconv.Itoa(r.Limit) + "_" + strconv.Itoa(r.Offset)
	total, _ = betxinredis.Get(totalRedis).Int()
	topics, err = betxinredis.Get(topicsRedis).Result()
	convert.Unmarshal(topics, &data)
	if err == redis.Nil {

		switch {
		case r.Offset >= 100:
			r.Offset = 100
		case r.Limit <= 0:
			r.Limit = 10
		}

		if r.Limit == 0 {
			r.Limit = 10
		}
		data, total, code = model.GetTopicByCid(convert.StrToNum(cid), r.Limit, r.Offset)
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR_GET_TOPIC, nil)
			return
		}

		topics = convert.Marshal(&data)
		betxinredis.Set(totalRedis, total, v1.REDISEXPIRE)
		betxinredis.Set(topicsRedis, topics, v1.REDISEXPIRE)

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

// GetTopicByTitle 通过标题获取信息
func GetTopicByTitle(c *gin.Context) {
	var data []model.Topic
	var total int
	var code int
	var err error

	var r ListRequest
	if err = c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	fmt.Println(r)

	switch {
	case r.Offset >= 100:
		r.Offset = 100
	case r.Limit <= 0:
		r.Limit = 10
	}

	if r.Limit == 0 {
		r.Limit = 10
	}
	data, total, code = model.SearchTopic(r.Offset, r.Limit, "intro LIKE  ?", "%"+r.Intro+"%")
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_GET_TOPIC, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, ListResponse{
		TotalCount: total,
		List:       data,
	})
}
