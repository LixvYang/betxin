// 结束topic返还给用户的钱
package model

import (
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Bonuse struct {
	Id        uint            `gorm:"primarykey" json:"id"`
	UserId    string          `gorm:"type:varchar(50);index;" json:"user_id"`
	Tid       string          `gorm:"varchar(50);" json:"tid"`
	AssetId   string          `gorm:"type:varchar(50);" json:"asset_id"`
	Amount    decimal.Decimal `gorm:"type:decimal(16, 8)" json:"amount"`
	Memo      string          `gorm:"type:varchar(200);" json:"memo"`
	TraceId   string          `gorm:"type:varchar(50);not null;uniqueIndex;" json:"trace_id"`
	CreatedAt time.Time       `gorm:"datatime(3)" json:"created_at"`
	UpdatedAt time.Time       `gorm:"datatime(3)" json:"updated_at"`
}

func CreateBonuse(data *Bonuse) int {
	if err := db.Exec("insert into bonuse (user_id, tid, asset_id, amount, memo, trace_id, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?, ?)", data.UserId, data.Tid, data.AssetId, data.Amount, data.Memo, data.TraceId, time.Now(), time.Now()).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetBonuseByTraceId(trace_id string) (Bonuse, int) {
	var bonuse Bonuse
	if err := db.Where("trace_id = ?", trace_id).Last(&bonuse).Error; err != nil {
		return bonuse, errmsg.ERROR
	}
	return bonuse, errmsg.SUCCSE
}

func ListBonuses(offset int, limit int) ([]Bonuse, int, int) {
	var bonuse []Bonuse
	var total int64
	err := db.Find(&bonuse).Limit(limit).Offset(offset).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, errmsg.ERROR
	}
	return bonuse, int(total), errmsg.SUCCSE
}

func UpdateBonuse(id int, data *Bonuse) int {
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
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Last(&Bonuse{}, id).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	maps["asset_id"] = data.AssetId
	maps["amount"] = data.Amount
	maps["memo"] = data.Memo
	maps["trace_id"] = data.TraceId
	maps["user_id"] = data.UserId

	if err := db.Model(&Bonuse{}).Where("id = ? ", id).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func DeleteBonuse(id string) int {
	if err := db.Where("id = ?", id).Delete(&Bonuse{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetBonusesByUserId(user_id int) ([]Bonuse, int) {
	var bonuse []Bonuse
	var total int64
	db.Model(&bonuse).Count(&total)
	if err := db.Find(&bonuse).Where("user_id = ?", user_id).Error; err != nil {
		return nil, 0
	}
	return bonuse, int(total)
}
