package model

import (
	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"gorm.io/gorm"
)

type Category struct {
	Id           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryName string `gorm:"type:varchar(20);not null;" json:"category_name"`
}

// 查询分类是否存在
func CheckCategory(category_name string) int {
	var cate Category
	db.Select("id").Where("category_name = ?", category_name).First(&cate)
	if cate.Id > 0 {
		return errmsg.ERROR_CATENAME_USED
	}
	return errmsg.SUCCSE
}

// Create category
func CreateCatrgory(data *Category) int {
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 根据种类id获取话题数据.
func GetCategoryById(id int) (Category, int) {
	var cate Category
	if err := db.Where("id = ?", id).First(&cate).Error; err != nil {
		return cate, errmsg.ERROR
	}
	return cate, errmsg.SUCCSE
}

// GetCate 查询分类列表
func ListCategories(offset int, limit int) ([]Category, int, int) {
	var cate []Category
	var total int64
	if err := db.Model(&cate).Count(&total).Error; err != nil {
		return nil, 0, errmsg.ERROR
	}

	if err := db.Limit(limit).Offset(offset).Find(&cate).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, errmsg.ERROR
	}
	return cate, int(total), errmsg.SUCCSE
}

// EditCate 编辑分类信息
func UpdateCate(id int, categoryName string) int {
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
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Last(&Category{}, id).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	maps["category_name"] = categoryName

	if err := db.Model(&Category{}).Where("id = ? ", id).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// Delete category by id
func DeleteCategory(id int) int {
	if err := db.Where("id = ?", id).Delete(&Category{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func SearchCategory(categoryName string, offset int, limit int) ([]Category, int, int) {
	var cate []Category
	var total int64

	if err := db.Model(&Category{}).Where("category_name LIKE ?", categoryName+"%").Offset(offset).Limit(limit).Count(&total).Error; err != nil {
		return nil, 0, errmsg.ERROR
	}

	if err := db.Model(&Category{}).Offset(offset).Limit(limit).Where("category_name LIKE ?", categoryName+"%").Find(&cate).Error; err != nil {
		return nil, 0, errmsg.ERROR
	}
	return cate, int(total), errmsg.SUCCSE
}
