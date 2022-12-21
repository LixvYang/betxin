package comment

import (
	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	"github.com/lixvyang/betxin/pkg/convert"

	"github.com/gin-gonic/gin"
)

func GetCommentById(c *gin.Context) {
	id := c.Param("id")

	comment, code := model.GetCommentById(convert.StrToNum(id))
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, comment)
}
