package model

import (
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"gorm.io/gorm"
)

type MixpayOrder struct {
	Tid      string `gorm:"type:varchar(36)" json:"tid"`
	YesRatio bool   `json:"yes_ratio"`
	NoRatio  bool   `json:"no_ratio"`
	Uid      string `gorm:"varchar(36)" json:"uid"`
	OrderId  string `gorm:"type:varchar(50)" json:"order_id"`
	TraceId  string `gorm:"type:varchar(50)" json:"trace_id"`
	PayeeId  string `gorm:"type:varchar(50)" json:"payee_id"`

	CreatedAt time.Time `gorm:"type:datetime(3)" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime(3)" json:"updated_at"`
}

func CreateMixpayOrder(data *MixpayOrder) int {
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func UpdateMixpayOrder(data *MixpayOrder) int {
	if err := db.Model(&MixpayOrder{}).Where("order_id = ? AND payee_id = ?", data.OrderId, data.PayeeId).Update("trace_id", data.TraceId).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func DeleteMixpayOrder(orderId string) int {
	if err := db.Where("order_id = ?", orderId).Delete(&MixpayOrder{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetMixpayOrder(traceId string) (MixpayOrder, int) {
	var mixpayOrder MixpayOrder
	if err := db.Model(&MixpayOrder{}).Where("trace_id = ?", traceId).Last(&mixpayOrder).Error; err != nil {
		return mixpayOrder, errmsg.ERROR
	}
	return mixpayOrder, errmsg.SUCCSE
}

func ListMixpayOrder(limit, offset int) ([]MixpayOrder, int, int) {
	var mixinpay []MixpayOrder
	var total int64
	if err := db.Model(&mixinpay).Count(&total).Error; err != nil {
		return mixinpay, 0, errmsg.ERROR
	}

	if err := db.Limit(limit).Offset(offset).Find(&mixinpay).Error; err != nil && err != gorm.ErrRecordNotFound {
		return mixinpay, 0, errmsg.ERROR
	}
	return mixinpay, int(total), errmsg.SUCCSE
}
