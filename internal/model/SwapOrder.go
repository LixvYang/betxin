// 从4swap交换出来tx的结构保存到数据库
package model

import (
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type SwapOrder struct {
	Type       string          `gorm:"type:varchar(20);" json:"type"`
	SnapshotId string          `gorm:"type:varchar(50)" json:"snapshot_id,omitempty"`
	AssetID    string          `gorm:"type:varchar(50)" json:"asset_id"`
	Amount     decimal.Decimal `gorm:"type:decimal(16, 8)" json:"amount"`
	TraceId    string          `gorm:"type:varchar(50);not null" json:"trace_id"`
	Memo       string          `gorm:"type:varchar(255)" json:"memo"`
	State      string          `gorm:"type:varchar(20)" json:"state"`
	CreatedAt  time.Time       `gorm:"type:datetime(3)" json:"created_at"`
}

func CreateSwapOrder(data *SwapOrder) int {
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func DeleteSwapOrder(traceId string) int {
	if err := db.Delete(&SwapOrder{}, traceId).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func UpdateSwapOrder(traceId string, data *SwapOrder) int {
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
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Last(&SwapOrder{}, traceId).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	maps["amount"] = data.Amount
	maps["asset_id"] = data.AssetID
	maps["memo"] = data.Memo
	maps["state"] = data.State
	maps["type"] = data.Type
	if err := db.Model(&UserToTopic{}).Where("trace_id = ?", traceId).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetSwapOrder(traceId string) (SwapOrder, int) {
	var swapOrder SwapOrder
	if err := db.First(&swapOrder, traceId).Error; err != nil {
		return SwapOrder{}, errmsg.ERROR
	}
	return swapOrder, errmsg.SUCCSE
}

func ListSwapOrders(offset int, limit int, query interface{}, args ...interface{}) ([]SwapOrder, int, int) {
	var swapOrders []SwapOrder
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

func ListSwapOrdersNoLimit(offset int, limit int) ([]SwapOrder, int, int) {
	var swapOrders []SwapOrder
	var total int64
	db.Model(&swapOrders).Count(&total)
	err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&swapOrders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, errmsg.ERROR
	}
	return swapOrders, int(total), errmsg.SUCCSE
}
