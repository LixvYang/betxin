package administrator

import (
	"log"
	"strconv"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

func UpdateAdministrator(c *gin.Context) {
	var admin model.Administrator
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&admin); err != nil {
		log.Panicln(err)
	}

	code := model.UpdateAdministrator(id, &admin)
	if code != errmsg.SUCCSE {
		v1.SendResponse(c, errmsg.ERROR_UPDATE_USER, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, admin.Id)
}
