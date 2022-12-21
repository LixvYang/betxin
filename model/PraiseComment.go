package model

import (
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"gorm.io/gorm"
)

// 点赞评论的人的列表
type PraiseComment struct {
	Id  int    `gorm:"primarykey" json:"id"`
	Cid int    `gorm:"type:int;index:cid_uid_praise" json:"cid"` // 评论id
	Uid string `gorm:"type:varchar(36);index:cid_uid_praise" json:"uid"`

	CreatedAt time.Time `gorm:"type:datetime(3); index:comment_created_at" json:"created_at"`
}

// check Collect
func CheckPraiseComment(cid int, uid string) int {
	var pc PraiseComment
	if err := db.Model(&PraiseComment{}).Where("cid = ? AND uid = ?", cid, uid).Last(&pc).Error; err != nil || pc.Id == 0 {
		// 没有收藏或者查询失败
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func CreatePraiseComment(data *PraiseComment) int {
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	db.Model(&Comment{}).Where("id = ?", data.Cid).UpdateColumn("praise_num", gorm.Expr("praise_num + ?", 1))
	return errmsg.SUCCSE
}

func DeletePraise(cid int, uid string) int {
	if err := db.Where("cid = ? AND uid = ?", cid, uid).Delete(&PraiseComment{}).Error; err != nil {
		return errmsg.ERROR
	}
	db.Model(&Comment{}).Where("id = ?", cid).UpdateColumn("praise_num", gorm.Expr("praise_num - ?", 1))

	return errmsg.SUCCSE
}

func ListPraiseCommentByUid(uid string) ([]PraiseComment, int, int) {
	var pc []PraiseComment
	var total int64

	if err := db.Where("uid = ?", uid).Count(&total).Find(&pc).Error; err != nil {
		return pc, int(total), errmsg.ERROR
	}
	return pc, int(total), errmsg.SUCCSE
}

func ListPraiseCommentByCid(cid string) ([]PraiseComment, int, int) {
	var pc []PraiseComment
	var total int64
	if err := db.Where("cid = ?", cid).Count(&total).Find(&pc).Error; err != nil {
		return pc, int(total), errmsg.ERROR
	}
	return pc, int(total), errmsg.SUCCSE
}

func ListPraiseComment(limit, offset int) ([]PraiseComment, int, int) {
	var praiseComment []PraiseComment
	var total int64
	var err error
	err = db.Model(&PraiseComment{}).Count(&total).Error
	err = db.Model(&PraiseComment{}).Limit(limit).Offset(offset).Order("created_at DESC").Find(&praiseComment).Error
	if err != nil {
		return nil, 0, errmsg.ERROR
	}
	return praiseComment, int(total), errmsg.SUCCSE
}
