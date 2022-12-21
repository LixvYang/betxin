package model

import (
	"sync"
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type UserToTopic struct {
	Id            int             `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	UserId        string          `gorm:"type:varchar(50);index:useid_topicid_index;index:userid_yes_no_index" json:"user_id"`
	Tid           string          `gorm:"type:varchar(36);not null;index:useid_topicid_index" json:"tid"`
	Topic         Topic           `gorm:"foreignKey:Tid;references:Tid;" json:"topic"`
	YesRatioPrice decimal.Decimal `gorm:"type:decimal(16,8);index:userid_yes_no_index" json:"yes_ratio_price"`
	NoRatioPrice  decimal.Decimal `gorm:"type:decimal(16,8);index:userid_yes_no_index" json:"no_ratio_price"`

	CreatedAt time.Time `gorm:"type:datetime(3)" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime(3)" json:"updated_at"`
}

func CheckUserToTopic(userId, tid string) int {
	var userToTopic UserToTopic
	db.Select("id").Where("user_id = ? AND tid = ?", userId, tid).Last(&userToTopic)
	if userToTopic.Id == 0 {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetUserToTopic(userId, tid string) (UserToTopic, int) {
	var userToTopic UserToTopic
	db.Model(&UserToTopic{}).Where("user_id = ? AND tid = ?", userId, tid).Last(&userToTopic)
	if userToTopic.Id == 0 {
		return userToTopic, errmsg.ERROR
	}
	return userToTopic, errmsg.SUCCSE
}

func CreateUserToTopic(data *UserToTopic) int {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	if err := db.Exec("insert into user_to_topic (user_id, tid, yes_ratio_price, no_ratio_price, created_at, updated_at) values (?, ?, ?, ?, ?, ?)", data.UserId, data.Tid, data.YesRatioPrice, data.NoRatioPrice, time.Now(), time.Now()).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func RefundUserToTopic(data *UserToTopic) (decimal.Decimal, decimal.Decimal, int) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, errmsg.ERROR
	}

	// 锁住指定 id 的 User 记录
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ? AND tid = ?", data.UserId, data.Tid).Error; err != nil {
		tx.Rollback()
		return decimal.Decimal{}, decimal.Decimal{}, errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	var userToTopic UserToTopic
	var yesFee decimal.Decimal
	var noFee decimal.Decimal
	err := db.Model(&UserToTopic{}).Where("user_id = ? AND tid = ?", data.UserId, data.Tid).Last(&userToTopic).Error
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, errmsg.ERROR
	}

	// 先扣除手续费
	if data.YesRatioPrice.GreaterThan(decimal.NewFromFloat(0)) {
		yesFee = userToTopic.YesRatioPrice.Mul(decimal.NewFromFloat(0.05))
		db.Model(&UserToTopic{}).Where("user_id = ? AND tid = ?", data.UserId, data.Tid).Update("yes_ratio_price", gorm.Expr("yes_ratio_price * 0.95"))
	}

	if data.NoRatioPrice.GreaterThan(decimal.NewFromFloat(0)) {
		noFee = userToTopic.NoRatioPrice.Mul(decimal.NewFromFloat(0.05))
		db.Model(&UserToTopic{}).Where("user_id = ? AND tid = ?", data.UserId, data.Tid).Update("no_ratio_price", gorm.Expr("no_ratio_price * 0.95"))
	}

	maps["YesRatioPrice"] = gorm.Expr("yes_ratio_price - ?", data.YesRatioPrice)
	maps["NoRatioPrice"] = gorm.Expr("no_ratio_price - ?", data.NoRatioPrice)

	if err := db.Model(&UserToTopic{}).Where("user_id = ? AND tid = ?", data.UserId, data.Tid).Updates(maps).Error; err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, errmsg.ERROR
	}
	return yesFee, noFee, errmsg.SUCCSE
}

func UpdateUserToTopic(data *UserToTopic) int {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return errmsg.ERROR
	}

	// 锁住指定 id 的 User 记录
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ? AND tid = ?", data.UserId, data.Tid).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	maps["YesRatioPrice"] = gorm.Expr("yes_ratio_price + ?", data.YesRatioPrice)
	maps["NoRatioPrice"] = gorm.Expr("no_ratio_price + ?", data.NoRatioPrice)
	if err := db.Model(&UserToTopic{}).Where("user_id = ? AND tid = ?", data.UserId, data.Tid).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func DeleteUserToTopic(userId, tid string) int {
	if err := db.Where("user_id = ? AND tid = ?", userId, tid).Delete(&UserToTopic{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func ListUserToTopicsByUserId(userId string, offset, limit int) ([]UserToTopic, int, int) {
	var userToTopics []UserToTopic
	var count int64

	if err := db.Preload("Topic").Model(&userToTopics).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return userToTopics, 0, errmsg.ERROR
	}

	if err := db.Preload("Topic").Model(&userToTopics).Where("user_id = ?").Limit(limit).Offset(offset).Order("created_at DESC").Find(userToTopics).Error; err != nil {
		return userToTopics, 0, errmsg.ERROR
	}

	return userToTopics, int(count), errmsg.SUCCSE
}

func ListUserToTopicsByTopicId(tid string, offset, limit int) ([]UserToTopic, int, int) {
	var userToTopics []UserToTopic
	var count int64

	if err := db.Preload("Topic").Model(&userToTopics).Where("tid = ?", tid).Count(&count).Error; err != nil {
		return userToTopics, 0, errmsg.ERROR
	}

	if err := db.Preload("Topic").Model(&userToTopics).Where("tid = ?").Limit(limit).Offset(offset).Order("created_at DESC").Find(userToTopics).Error; err != nil {
		return userToTopics, 0, errmsg.ERROR
	}

	return userToTopics, int(count), errmsg.SUCCSE
}

func ListUserToTopics(offset, limit int) ([]UserToTopic, int, int) {
	var userToTopics []UserToTopic
	var count int64

	if err := db.Preload("Topic").Model(&UserToTopic{}).Count(&count).Error; err != nil {
		return userToTopics, 0, errmsg.ERROR
	}

	if err := db.Select("user_to_topic.tid, id,user_to_topic.yes_ratio_price, user_to_topic.no_ratio_price, user_id, user_to_topic.updated_at,user_to_topic.created_at, Topic.cid").Limit(limit).Offset(offset).Joins("Topic").Find(&userToTopics).Error; err != nil {
		return nil, 0, errmsg.ERROR
	}

	return userToTopics, int(count), errmsg.SUCCSE
}

// 列出话题下的哪些用户赢了
func ListUserToTopicsWin(tid string, win string) ([]UserToTopic, int, int) {
	var userToTopics []UserToTopic
	var count int64

	if err := db.Model(&userToTopics).Count(&count).Error; err != nil {
		return userToTopics, 0, errmsg.ERROR
	}

	var mutex sync.Mutex
	mutex.Lock()
	if win == "yes_win" {
		db = db.Where("yes_ratio_price > 0")
	} else {
		db = db.Where("no_ratio_price > 0")
	}
	if err := db.Where("tid = ?", tid).Find(&userToTopics).Error; err != nil {
		return userToTopics, 0, errmsg.ERROR
	}
	mutex.Unlock()

	return userToTopics, int(count), errmsg.SUCCSE
}

func ListUserToTopicsByUserIdNoLimit(userId string) ([]UserToTopic, int, int) {
	var userToTopics []UserToTopic
	var count int64

	if err := db.Preload("Topic").Model(&userToTopics).Where("user_id = ?", userId).Order("created_at DESC").Count(&count).Find(&userToTopics).Error; err != nil {
		return userToTopics, 0, errmsg.ERROR
	}

	return userToTopics, int(count), errmsg.SUCCSE
}
