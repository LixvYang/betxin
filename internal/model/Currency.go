package model

import (
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Currency struct {
	AssetId  string          `gorm:"uniqueIndex;type:varchar(255);not null" json:"asset_id"`
	PriceUsd decimal.Decimal `gorm:"type:decimal(20,10); not null" json:"price_usd"`
	PriceBtc decimal.Decimal `gorm:"type:decimal(20,10); not null" json:"price_btc"`
	ChainId  string          `gorm:"type:varchar(255); not null" json:"chain_id"`
	IconUrl  string          `gorm:"type:varchar(255)" json:"icon_url"`
	Symbol   string          `gorm:"type:varchar(255)" json:"symbol"`

	CreatedAt time.Time `gorm:"type:datetime(3)" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime(3)" json:"updated_at"`
}

func CheckCurrency(assetId string) int {
	var currency Currency
	if err := db.Model(&Currency{}).Where("asset_id = ?", assetId).Find(&currency).Error; err != nil || currency.AssetId != "" {
		// 有值
		return errmsg.ERROR
	}
	// 无值
	return errmsg.SUCCSE
}

func CreateCurrency(data *Currency) int {
	if err := db.Model(&Currency{}).Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func UpdateCurrency(data *Currency) int {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return errmsg.ERROR
	}

	// 锁住指定 id 的 记录
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&Currency{}).Where("asset_id = ?", data.AssetId).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}
	var maps = make(map[string]interface{})
	maps["PriceUsd"] = data.PriceUsd
	maps["PriceBtc"] = data.PriceBtc
	maps["ChainId"] = data.ChainId
	maps["Symbol"] = data.Symbol
	if err := db.Model(&Currency{}).Where("asset_id = ?", data.AssetId).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 查询列表
func ListCurrencies() ([]Currency, int, int) {
	var currency []Currency
	var total int64
	db.Model(&currency).Count(&total)
	err := db.Find(&currency).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, errmsg.ERROR
	}
	return currency, int(total), errmsg.SUCCSE
}

func GetCurrencyById(asset_id string) (Currency, int) {
	var currency Currency
	if err := db.Where("asset_id = ?", asset_id).First(&currency).Error; err != nil {
		return Currency{}, errmsg.ERROR
	}
	return currency, errmsg.SUCCSE
}

func GetCurrencyBySymbol(symbol string) (Currency, int) {
	var currency Currency
	if err := db.Where("symbol = ?", symbol).First(&currency).Error; err != nil {
		return Currency{}, errmsg.ERROR
	}
	return currency, errmsg.SUCCSE
}
