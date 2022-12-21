package topic

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
	"github.com/lixvyang/betxin/pkg/convert"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// GetArtInfo 查询单个话题信息
func GetTopicInfoById(c *gin.Context) {
	tid := c.Param("tid")
	var data model.Topic
	var code int
	var topic string
	var err error

	topic, err = betxinredis.Get(v1.TOPIC_GET + tid).Result()
	convert.Unmarshal(topic, &data)
	if err == redis.Nil {
		data, code = model.GetTopicById(tid)
		if code != errmsg.SUCCSE {
			v1.SendResponse(c, errmsg.ERROR_GET_TOPIC, nil)
			return
		}
		topic = convert.Marshal(&data)
		betxinredis.Set(v1.TOPIC_GET+tid, topic, v1.REDISEXPIRE)
		v1.SendResponse(c, errmsg.SUCCSE, data)
	} else if err != nil {
		v1.SendResponse(c, errmsg.ERROR, nil)
	} else {
		v1.SendResponse(c, errmsg.SUCCSE, data)
	}
}
