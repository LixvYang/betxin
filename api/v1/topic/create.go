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

type CreateReqeust struct {
	Tid     string `json:"tid"`
	Cid     string `json:"cid"`
	Title   string `json:"title"`
	Intro   string `json:"intro"`
	ImgUrl  string `json:"img_url"`
	EndTime string `json:"end_time"`
}

func CreateTopic(c *gin.Context) {
	var r CreateReqeust
	if err := c.ShouldBindJSON(&r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", r.EndTime, time.Local)
	if err != nil {
		log.Println(err)
	}
	t := &model.Topic{
		Tid:     r.Tid,
		Cid:     convert.StrToNum(r.Cid),
		Title:   r.Title,
		Intro:   r.Intro,
		ImgUrl:  r.ImgUrl,
		EndTime: endTime,
	}

	code := model.CreateTopic(t)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	betxinredis.BatchDel("topic")

	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
