// 记录topic结束, 从机器人转给用户的表
package model

import (
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type MixinNetworkSnapshot struct {
	SnapshotId     string          `gorm:"type:varchar(50)" json:"snapshot_id"`
	TraceId        string          `gorm:"type:varchar(50);not null;" json:"trace_id"`
	AssetId        string          `gorm:"type:varchar(50);index" json:"asset_id"`
	OpponentID     string          `gorm:"type:varchar(50)" json:"opponent_id"`
	Amount         decimal.Decimal `gorm:"type:decimal(16, 8)" json:"amount"`
	Memo           string          `gorm:"type:varchar(200)" json:"memo"`
	Type           string          `gorm:"type:varchar(200)" json:"type"`
	OpeningBalance decimal.Decimal `json:"opening_balance"`
	ClosingBalance decimal.Decimal `json:"closing_balance"`

	CreatedAt time.Time `gorm:"type:datetime(3)" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime(3)" json:"updated_at"`
}

func CheckMixinNetworkSnapshot(traceId string) int {
	var mixinNetworkSnapshot MixinNetworkSnapshot
	if err := db.First(&mixinNetworkSnapshot, "trace_id = ?", traceId).Error; err != nil || err == gorm.ErrRecordNotFound {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func CreateMixinNetworkSnapshot(data *MixinNetworkSnapshot) int {
	if err := db.Exec("insert into mixin_network_snapshot (snapshot_id, trace_id, asset_id, opponent_id, amount, memo, type, opening_balance, closing_balance, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		data.SnapshotId, data.TraceId, data.AssetId, data.OpponentID, data.Amount, data.Memo, data.Type, data.OpeningBalance, data.ClosingBalance, time.Now(), time.Now()); err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func DeleteMixinNetworkSnapshot(traceId string) int {
	if err := db.Delete(&MixinNetworkSnapshot{}, traceId).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetMixinNetworkSnapshot(traceId string) (MixinNetworkSnapshot, int) {
	var mixinNetworkSnapshot MixinNetworkSnapshot
	if err := db.Last(&mixinNetworkSnapshot, traceId).Error; err != nil {
		return mixinNetworkSnapshot, errmsg.ERROR
	}
	return mixinNetworkSnapshot, errmsg.SUCCSE
}

func UpdateMixinNetworkSnapshot(traceId string, data *MixinNetworkSnapshot) int {
	// 锁住指定 id 的 User 记录
	if err := db.Exec("update mixin_network_snapshot set snapshot_id = ?, asset_id = ?, opponent_id = ?, amount = ?, memo = ?, type = ?, opening_balance = ?, closing_balance = ? where trace_id = ?",
		data.SnapshotId, data.AssetId, data.OpponentID, data.Amount, data.Memo, data.Type, data.OpeningBalance, data.ClosingBalance, traceId).Error; err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSE
}

func ListMixinNetworkSnapshots(offset int, limit int) ([]MixinNetworkSnapshot, int, int) {
	var mixinNetworkSnapshot []MixinNetworkSnapshot
	var total int64
	var err error

	if err = db.Model(&MixinNetworkSnapshot{}).Count(&total).Error; err != nil {
		return mixinNetworkSnapshot, int(total), errmsg.ERROR
	}

	if err = db.Limit(limit).Offset(offset).Find(&mixinNetworkSnapshot).Error; err != nil {
		return mixinNetworkSnapshot, int(total), errmsg.ERROR
	}
	return mixinNetworkSnapshot, int(total), errmsg.SUCCSE
}

func ListMixinNetworkSnapshotsByUserId(userId string, offset int, limit int) ([]MixinNetworkSnapshot, int, int) {
	var mixinNetworkSnapshot []MixinNetworkSnapshot
	var total int64

	err := db.Find(&mixinNetworkSnapshot).Where("user_id = ?", userId).Error
	db.Model(mixinNetworkSnapshot).Count(&total)
	if err != nil {
		return nil, 0, errmsg.ERROR
	}
	return mixinNetworkSnapshot, int(total), errmsg.SUCCSE
}
