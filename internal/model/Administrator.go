package model

import (
	"github.com/lixvyang/betxin/internal/utils/errmsg"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Administrator struct {
	Id        int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string         `gorm:"type:varchar(20);not null " json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password  string         `gorm:"type:varchar(500);not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	CreatedAt time.Time      `gorm:"type:datetime(3); not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime(3); not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}


// BeforeCreate 密码加密&权限控制
func (u *Administrator) BeforeCreate(_ *gorm.DB) (err error) {
	u.Password = ScryptPw(u.Password)
	return nil
}

func (u *Administrator) BeforeUpdate(_ *gorm.DB) (err error) {
	u.Password = ScryptPw(u.Password)
	return nil
}

// ScryptPw 生成密码
func ScryptPw(password string) string {
	const cost = 10

	HashPw, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Println(err)
	}
	return string(HashPw)
}


func CheckAdministrator(username string) int {
	var admin Administrator
	db.Select("id").Where("username = ?", username).First(&admin)
	if admin.Id > 0 {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// CreateAdministrator 新增管理员
func CreateAdministrator(data *Administrator) int {
	//data.Password = ScryptPw(data.Password)
	if err := db.Create(&data).Error; err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}

// Delete 管理员
func DeleteAdministrator(id int) int {
	if err := db.Where("id = ? ", id).Delete(&Administrator{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func GetAdministratorById(id int) (Administrator, int) {
	var admin Administrator
	if err := db.Where("id = ?", id).First(&admin).Error; err != nil {
		return admin, errmsg.ERROR
	}
	return admin, errmsg.SUCCSE
}

func ListAdministrators(offset int, limit int) ([]Administrator, int, int) {
	var admin []Administrator
	var total int64
	if err := db.Limit(limit).Offset(offset).Find(&admin).Limit(-1).Offset(-1).Count(&total).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, errmsg.ERROR
	}
	return admin, int(total), errmsg.SUCCSE
}

func UpdateAdministrator(id int, admin *Administrator) int {
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
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Last(&Administrator{}, id).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	maps["username"] = admin.Username
	maps["password"] = admin.Password

	if err := db.Model(&Administrator{}).Where("id = ? ", id).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// CheckLogin 后台登录验证
func CheckLogin(username string, password string) (Administrator, int) {
	var user Administrator
	var PasswordErr error

	db.Where("username = ?", username).First(&user)
	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if user.Id == 0 {
		return user, errmsg.ERROR_USER_NOT_EXIST
	}

	if PasswordErr != nil {
		return user, errmsg.ERROR_PASSWORD_WRONG
	}

	return user, errmsg.SUCCSE
}
