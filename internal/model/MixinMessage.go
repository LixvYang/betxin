package model

import (
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
)

type MixinMessage struct {
	UserId         string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	ConversationId string    `gorm:"type:varchar(50);not null;" json:"conversation_id"`
	Category       string    `gorm:"type:varchar(50); not null" json:"category"`
	MessageId      string    `gorm:"type:varchar(50);not null;comment:UUID;uniqueIndex" json:"message_id"`
	Content        string    `gorm:"type:varchar(50);comment:decrepted data;" json:"content"`
	CreatedAt      time.Time `gorm:"type:datetime(3); not null;" json:"created_at"`
	UpdatedAt      time.Time `gorm:"type:datetime(3);not null;" json:"updated_at"`
}

func (m *MixinMessage) New() {

}

func CreateMixinMessage(data *MixinMessage) int {
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 根据message_id获取message
func GetMixinMessageById(message_id int) (MixinMessage, int) {
	var message MixinMessage
	if err := db.Where("message_id = ?", message_id).First(&message).Error; err != nil {
		return message, errmsg.ERROR
	}
	return message, errmsg.SUCCSE
}

// 根据user_id获取messages
func ListMixinMessageByUserId(user_id string, offset, limit int) ([]MixinMessage, int, int) {
	var message []MixinMessage
	var err error
	var total int64

	err = db.Where("user_id = ?", user_id).Offset(offset).Limit(limit).Error
	db.Model(&message).Count(&total)
	if err != nil {
		return nil, 0, errmsg.ERROR
	}
	return message, int(total), errmsg.SUCCSE
}

func DeleteMixinMessageByMessageId(message_id string) int {
	if err := db.Where("message_id = ?", message_id).Delete(&MixinMessage{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func UpdateMixinMessageByMsgId(message_id string, msg *MixinMessage) int {
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
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Last(&Category{}, message_id).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	maps["Content"] = msg.Content
	maps["Category"] = msg.Category

	if err := db.Model(&Category{}).Where("message_id = ? ", message_id).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
