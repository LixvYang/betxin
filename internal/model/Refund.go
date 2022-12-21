// 结束topic返还给用户的钱
package model

import (
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Refund struct {
	Id          uint            `gorm:"primarykey" json:"id"`
	UserId      string          `gorm:"type:varchar(50);index;" json:"user_id"`
	Tid         string          `gorm:"varchar(50);" json:"tid"`
	AssetId     string          `gorm:"type:varchar(50);" json:"asset_id"`
	RefundPrice decimal.Decimal `gorm:"type:decimal(16, 8)" json:"refund_price"`
	Memo        string          `gorm:"type:varchar(200);" json:"memo"`
	TraceId     string          `gorm:"type:varchar(50);not null;uniqueIndex;" json:"trace_id"`
	CreatedAt   time.Time       `gorm:"datatime(3)" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"datatime(3)" json:"updated_at"`
}

func CreateRefund(data *Refund) int {
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetRefundByTraceId(trace_id string) (Refund, int) {
	var refund Refund
	if err := db.Where("trace_id = ?", trace_id).Last(&refund).Error; err != nil {
		return refund, errmsg.ERROR
	}
	return refund, errmsg.SUCCSE
}

func ListRefunds(offset int, limit int) ([]Refund, int, int) {
	var refund []Refund
	var total int64
	err := db.Find(&refund).Limit(limit).Offset(offset).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, errmsg.ERROR
	}
	return refund, int(total), errmsg.SUCCSE
}

func UpdateRefund(tracdId string, data *Refund) int {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return errmsg.ERROR
	}
	refund := Refund{}
	// 锁住指定 id 的 User 记录
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("trace_id = ?", tracdId).Last(&refund).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	maps["user_id"] = data.UserId
	maps["tid"] = data.Tid
	maps["asset_id"] = data.AssetId
	maps["refund_price"] = data.RefundPrice
	maps["memo"] = data.Memo

	if err := db.Model(&Refund{}).Where("trace_id = ? ", tracdId).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func DeleteRefund(id string) int {
	if err := db.Where("id = ?", id).Delete(&Refund{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetRefundsByUserId(user_id int) ([]Refund, int) {
	var refund []Refund
	var total int64
	db.Model(&refund).Count(&total)
	if err := db.Find(&refund).Where("user_id = ?", user_id).Error; err != nil {
		return nil, 0
	}
	return refund, int(total)
}
