package topic

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"

	"github.com/gin-gonic/gin"
)

func DeleteTopic(c *gin.Context) {
	tid := c.Param("id")

	if code := model.DeleteTopic(tid); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_DELETE_TOPIC, nil)
		return
	}
	betxinredis.BatchDel("topic")

	betxinredis.DelKeys(v1.TOPIC_LIST, v1.TOPIC_TOTAL, v1.TOPIC_GET+tid)
	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
