package model

import (
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
)

type Collect struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserId    string    `gorm:"type:varchar(50);not null;index:user_collect_topic" json:"user_id"`
	Tid       string    `gorm:"index:user_collect_topic;type:varchar(36);not null;uniqueKey" json:"tid"`
	Topic     Topic     `gorm:"foreignKey:Tid;references:Tid;" json:"topic"`
	CreatedAt time.Time `gorm:"datetime(3)" json:"created_at"`
	UpdatedAt time.Time `gorm:"datatime(3)" json:"updated_at"`
}

// check Collect
func CheckCollect(userId string, Tid string) int {
	var collect Collect
	if err := db.Model(&Collect{}).Where("user_id = ? AND tid = ?", userId, Tid).Last(&collect).Error; err != nil || collect.ID == 0 {
		// 没有收藏或者查询失败
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// Create Collect
func CreateCollect(data *Collect) int {
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// list Collect
func ListCollects(offset int, limit int) ([]Collect, int, int) {
	var collects []Collect
	var count int64

	if err := db.Preload("Topic").Model(&Collect{}).Count(&count).Error; err != nil {
		return collects, int(count), errmsg.ERROR
	}
	if err := db.Preload("Topic").Where("").Offset(offset).Limit(limit).Order("id desc").Find(&collects).Error; err != nil {
		return collects, int(count), errmsg.ERROR
	}

	return collects, int(count), errmsg.SUCCSE
}

// 根据标签user_id获取收藏数据.
func GetCollectByUserId(userId string) ([]Collect, int, int) {
	var collects []Collect
	var count int64

	if err := db.Preload("Topic").Model(&collects).Where("user_id = ?", userId).Count(&count).Order("id desc").Find(&collects).Error; err != nil {
		return collects, int(count), errmsg.ERROR
	}

	return collects, int(count), errmsg.SUCCSE
}

// Delete collect by id
func DeleteCollect(user_id string, tid string) int {
	if err := db.Where("user_id = ? and tid = ?", user_id, tid).Delete(&Collect{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func DeleteCollectByTid(tid string) int {
	if err := db.Where("tid = ?", tid).Delete(&Collect{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
