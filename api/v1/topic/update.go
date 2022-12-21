package topic

import (
	"log"
	"time"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
	"github.com/lixvyang/betxin/pkg/convert"

	"github.com/gin-gonic/gin"
)

func UpdateTopic(c *gin.Context) {
	tid := c.Param("id")
	var r CreateReqeust
	if err := c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
	}

	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", r.EndTime, time.Local)
	if err != nil {
		log.Println(err)
	}

	topic := &model.Topic{
		Tid:     r.Tid,
		Cid:     convert.StrToNum(r.Cid),
		Title:   r.Title,
		Intro:   r.Intro,
		ImgUrl:  r.ImgUrl,
		EndTime: endTime,
	}

	if code := model.UpdateTopic(tid, topic); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_UPDATE_TOPIC, nil)
		return
	}

	betxinredis.BatchDel("topic")

	v1.SendResponse(c, errmsg.SUCCSE, tid)
}
