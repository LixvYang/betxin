package model

import (
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
)

type FeedBack struct {
	Id        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    string    `gorm:"type:varchar(36);index" json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	CreatedAt time.Time `gorm:"type:datetime(3); not null;" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime(3);not null;" json:"updated_at"`
}

func CreateFeedBack(data *FeedBack) int {
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func ListFeedBack(offset, limit int) ([]FeedBack, int, int) {
	var feedback []FeedBack
	var total int64
	var err error
	err = db.Model(&FeedBack{}).Count(&total).Error
	err = db.Model(&FeedBack{}).Limit(limit).Offset(offset).Order("created_at DESC").Find(&feedback).Error
	if err != nil {
		return nil, 0, errmsg.ERROR
	}
	return feedback, int(total), errmsg.SUCCSE
}

func ListFeedBackByUserId(user_id string) ([]FeedBack, int, int) {
	var message []FeedBack
	var err error
	var total int64

	err = db.Where("user_id = ?", user_id).Error
	db.Model(&message).Count(&total)
	if err != nil {
		return nil, 0, errmsg.ERROR
	}
	return message, int(total), errmsg.SUCCSE
}

func DeleteFeedBackById(id string) int {
	if err := db.Where("id = ?", id).Delete(&FeedBack{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
