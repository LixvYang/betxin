package upload

import (
	"context"
	"mime/multipart"

	"github.com/lixvyang/betxin/internal/utils"
	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func UpLoadFile(file multipart.File, fileSize int64) (string, int) {
	putPolicy := storage.PutPolicy{
		Scope: utils.Bucket,
	}
	mac := qbox.NewMac(utils.AccessKey, utils.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}

	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		return "", errmsg.ERROR
	}
	url := utils.QiniuSever + ret.Key
	return url, errmsg.SUCCSE

}
