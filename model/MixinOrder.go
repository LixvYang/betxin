// 记录mixin转账到机器人的交易  接收用户转账
package model

import (
	"sync"
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type MixinOrder struct {
	Type       string          `gorm:"type:varchar(20)" json:"type"`
	SnapshotId string          `gorm:"type:varchar(50)" json:"snapshot_id"`
	AssetId    string          `gorm:"type:varchar(50);not null;index" json:"asset_id"`
	Amount     decimal.Decimal `gorm:"type:decimal(16, 8)" json:"amount"`
	TraceId    string          `gorm:"type:varchar(50);not null;index" json:"trace_id"`
	Memo       string          `gorm:"type:varchar(255);" json:"memo"`
	CreatedAt  time.Time       `gorm:"type:datetime(3)" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"type:datetime(3)" json:"updated_at"`
}

func CreateMixinOrder(data *MixinOrder) int {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func DeleteMixinOrder(traceId string) int {
	if err := db.Delete(&User{}, traceId).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetMixinOrderByTraceId(traceId string) (MixinOrder, int) {
	var mixinOrder MixinOrder
	if err := db.First(&mixinOrder, traceId).Error; err != nil {
		return MixinOrder{}, errmsg.ERROR
	}
	return mixinOrder, errmsg.SUCCSE
}

func UpdateMixinOrder(traceId string, data *MixinOrder) int {
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
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Last(&MixinOrder{}, traceId).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	maps["asset_id"] = data.AssetId
	maps["amount"] = data.Amount
	maps["memo"] = data.Memo

	if err := db.Model(&Category{}).Where("trace_id = ? ", traceId).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func ListMixinOrder(limit, offset int, query interface{}, args ...interface{}) ([]MixinOrder, int, int) {
	var mixinOrder []MixinOrder
	var total int64
	if query != "" {
		db.Where(query, args...)
	}
	err := db.Limit(limit).Offset(offset).Order("created_at DESC").Count(&total).Find(&mixinOrder).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, errmsg.ERROR
	}
	return mixinOrder, int(total), errmsg.SUCCSE
}

func ListMixinOrderNoLimit(limit, offset int) ([]MixinOrder, int, int) {
	var mixinOrder []MixinOrder
	var total int64
	err := db.Model(&MixinOrder{}).Count(&total).Error
	err = db.Model(&MixinOrder{}).Limit(limit).Offset(offset).Order("created_at DESC").Find(&mixinOrder).Error
	if err != nil {
		return nil, 0, errmsg.ERROR
	}
	return mixinOrder, int(total), errmsg.SUCCSE
}
