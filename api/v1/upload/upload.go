package upload

import (
	"fmt"

	"github.com/lixvyang/betxin/internal/utils/upload"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		v1.SendResponse(c, errmsg.ERROR, nil)
	}
	fileSize := fileHeader.Size
	url, code := upload.UpLoadFile(file, fileSize)
	if code != errmsg.SUCCSE {
		fmt.Println("上传出错")
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, url)
}
