// 从4swap交换出来tx的结构保存到数据库
package model

import (
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type SendBack struct {
	TraceId    string          `gorm:"type:varchar(50);not null" json:"trace_id"`
	Type       string          `gorm:"type:varchar(20);" json:"type"`
	OpponentID string          `gorm:"type:varchar(50)" json:"opponent_id"`
	SnapshotId string          `gorm:"type:varchar(50)" json:"snapshot_id,omitempty"`
	AssetID    string          `gorm:"type:varchar(50)" json:"asset_id"`
	Amount     decimal.Decimal `gorm:"type:decimal(16, 8)" json:"amount"`
	Memo       string          `gorm:"type:varchar(255)" json:"memo"`

	CreatedAt time.Time `gorm:"type:datetime(3)" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime(3)" json:"updated_at"`
}

func CreateSendBack(data *SendBack) int {
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func DeleteSendBackByTraceId(traceId string) int {
	if err := db.Where("trace_id = ?", traceId).Delete(&SendBack{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func UpdateSendBack(traceId string, data *SendBack) int {
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
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("trace_id = ?", traceId).Last(&SendBack{}).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	maps["snapshot_id"] = data.SnapshotId
	maps["opponent_id"] = data.OpponentID
	maps["asset_id"] = data.AssetID
	maps["amount"] = data.Amount
	maps["memo"] = data.Memo
	maps["type"] = data.Type

	if err := db.Model(&MixinNetworkSnapshot{}).Where("trace_id = ? ", traceId).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}

	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetSendBack(traceId string) (SendBack, int) {
	var swapOrder SendBack
	if err := db.First(&swapOrder, traceId).Error; err != nil {
		return SendBack{}, errmsg.ERROR
	}
	return swapOrder, errmsg.SUCCSE
}

func ListSendBacks(offset int, limit int, query interface{}, args ...interface{}) ([]SendBack, int, int) {
	var swapOrders []SendBack
	var total int64
	if query != "" {
		db.Where(query, args...)
	}
	db.Model(&swapOrders).Count(&total)
	err := db.Find(&swapOrders).Limit(limit).Offset(offset).Order("created_at DESC").Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, errmsg.ERROR
	}
	return swapOrders, int(total), errmsg.SUCCSE
}

func ListSendBacksNoLimit(offset int, limit int) ([]SendBack, int, int) {
	var sendBack []SendBack
	var total int64
	db.Model(&sendBack).Count(&total)
	err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&sendBack).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, errmsg.ERROR
	}
	return sendBack, int(total), errmsg.SUCCSE
}
