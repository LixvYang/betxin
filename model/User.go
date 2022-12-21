package model

import (
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
)

type User struct {
	Id             int    `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	IdentityNumber string `gorm:"type:varchar(50);not null;index" json:"identity_number"`
	MixinUuid      string `gorm:"type:varchar(36);index;" json:"mixin_uuid"`
	FullName       string `gorm:"type:varchar(50);not null" json:"full_name"`
	AvatarUrl      string `gorm:"type:varchar(255);not null" json:"avatar_url"`
	MixinId        string `gorm:"type:varchar(50);not null;index;" json:"mixin_id"`
	SessionId      string `gorm:"type:varchar(50);" json:"session_id"`
	Phone          string `gorm:"type: varchar(30);" json:"phone"`

	CreatedAt time.Time `gorm:"type:datetime(3)" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime(3)" json:"updated_at"`
}

// CheckUser 查询用户是否存在
func CheckUser(user_id string) int {
	var user User
	db.Model(&User{}).Where("mixin_uuid = ?", user_id).Last(&user)
	if user.Id == 0 {
		return errmsg.ERROR //1001
	}
	return errmsg.SUCCSE
}

// Create user
func CreateUser(data *User) int {
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetUserById(user_id string) (User, int) {
	var user User
	if err := db.Model(&user).Where("mixin_uuid = ?", user_id).First(&user).Error; err != nil {
		return User{}, errmsg.ERROR
	}
	return user, errmsg.SUCCSE
}

// Delete user
func DeleteUser(user_id string) int {
	if err := db.Where("user_id = ?", user_id).Delete(User{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// EditUser 编辑用户信息
func UpdateUser(user_id string, data *User) int {
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
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ?", user_id).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	maps["full_name"] = data.FullName
	if err := db.Model(&User{}).Where("identity_number = ? ", user_id).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// GetUserByName gets an user by the username.
func GetUserByName(full_name string) (*User, int) {
	var user *User
	if err := db.Where("full_name LIKE ?", full_name+"%").First(&user).Error; err != nil {
		return nil, errmsg.ERROR
	}
	return user, errmsg.SUCCSE
}

// List users
func ListUser(offset, limit int) ([]User, int, int) {
	var users []User
	var count int64
	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		return users, 0, errmsg.ERROR
	}

	if err := db.Where("").Offset(offset).Limit(limit).Order("id desc").Find(&users).Error; err != nil {
		return users, 0, errmsg.ERROR
	}
	return users, int(count), errmsg.SUCCSE
}
