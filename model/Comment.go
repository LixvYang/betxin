package model

import (
	"fmt"
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"gorm.io/gorm"
)

type Comment struct {
	Id             int    `gorm:"primarykey" json:"id"`
	Tid            string `gorm:"type:varchar(36);index:comment_tid" json:"tid"`
	Content        string `gorm:"longtext" json:"content"`
	FromUid        string `gorm:"type:varchar(36);index:comment_from_uid" json:"from_uid"`
	FromUserName   string `gorm:"type:varchar(36)" json:"from_user_name"`
	FromUserAvatar string `gorm:"type:varchar(255)" json:"from_user_avatar"`
	PraiseNum      int    `gorm:"type:int(8); default 0" json:"praise_num"`
	ToUid          string `gorm:"type:varchar(36);default null" json:"to_uid"`
	ToUserName     string `gorm:"type:varchar(36);default null" json:"to_user_name"`

	CreatedAt time.Time `gorm:"type:datetime(3); index:comment_created_at" json:"created_at"`
}

func CreateComment(data *Comment) int {
	if err := db.Create(&data).Error; err != nil {
		fmt.Println(err)
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func ListCommentByTid(tid string, limit int, offset int) ([]Comment, int, int) {
	var comments []Comment
	var total int64

	if err := db.Model(&Comment{}).Where("tid = ?", tid).Count(&total).Error; err != nil {
		return comments, int(total), errmsg.ERROR
	}

	if err := db.Model(&Comment{}).Where("tid = ?", tid).Order("praise_num").Offset(offset).Limit(limit).Find(&comments).Error; err != nil {
		return nil, 0, errmsg.ERROR
	}
	return comments, int(total), errmsg.SUCCSE
}

func ListComment(limit, offset int) ([]Comment, int, int) {
	var comments []Comment
	var total int64
	if err := db.Model(&comments).Count(&total).Error; err != nil {
		return comments, 0, errmsg.ERROR
	}

	if err := db.Limit(limit).Offset(offset).Find(&comments).Error; err != nil && err != gorm.ErrRecordNotFound {
		return comments, 0, errmsg.ERROR
	}
	return comments, int(total), errmsg.SUCCSE
}

func GetCommentById(id int) (Comment, int) {
	var comment Comment
	err := db.Model(&comment).Where("id = ?", id).Last(&comment).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return comment, errmsg.ERROR
	}
	return comment, errmsg.SUCCSE
}
