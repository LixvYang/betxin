package usertotopic

import (
	"sync"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

var userToTopicPool = sync.Pool{
	New: func() any {
		return new(model.UserToTopic)
	},
}

func CreateUserToTopic(c *gin.Context) {
	// var r model.UserToTopic
	r := userToTopicPool.Get().(*model.UserToTopic)
	if err := c.ShouldBindJSON(r); err != nil {
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	if code := model.CreateUserToTopic(r); code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	userToTopicPool.Put(r)
	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
