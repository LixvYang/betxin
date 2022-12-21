package model

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"

	"gorm.io/gorm"

	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type Topic struct {
	Tid           string          `gorm:"type:varchar(36);index;" json:"tid"`
	Cid           int             `gorm:"type:int;not null" json:"cid"`
	Category      Category        `gorm:"foreignKey:Cid" json:"category"`
	Title         string          `gorm:"type:varchar(50);not null;index:title_intro_topic_index" json:"title"`
	Intro         string          `gorm:"type:varchar(255);not null;index:title_intro_topic_index" json:"intro"`
	CollectCount  int             `gorm:"type:int;default 0" json:"collect_count"`
	YesRatio      decimal.Decimal `gorm:"type:decimal(5,2);" json:"yes_ratio"`
	NoRatio       decimal.Decimal `gorm:"type:decimal(5,2);" json:"no_ratio"`
	YesRatioPrice decimal.Decimal `gorm:"type:decimal(16,8);default 0" json:"yes_ratio_price"`
	NoRatioPrice  decimal.Decimal `gorm:"type:decimal(16,8);default 0" json:"no_ratio_price"`
	TotalPrice    decimal.Decimal `gorm:"type:decimal(32,8);default 0;" json:"total_price"`
	ReadCount     int             `gorm:"type:int;default:0" json:"read_count"`
	ImgUrl        string          `gorm:"varchar(255);" json:"img_url"`
	IsStop        int             `gorm:"type:int;default 0;" json:"is_stop"`
	EndTime       time.Time       `json:"end_time"`

	CreatedAt time.Time `gorm:"type:datetime(3)" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime(3)" json:"updated_at"`
}

func (t *Topic) BeforeCreate(tx *gorm.DB) error {
	t.Tid = uuid.NewV4().String()
	t.YesRatio = decimal.NewFromFloat(0.5)
	t.NoRatio = decimal.NewFromFloat(0.5)
	return nil
}

func (t *Topic) BeforeUpdate(tx *gorm.DB) error {
	if t.IsStop == 1 || time.Now().After(t.EndTime) {
		return errors.New("话题已经停止")
	}
	return nil
}

func CheckTopic(title string) int {
	var topic Topic
	db.Select("tid").Where("title = ?", title).First(&topic)
	if topic.Intro != "" {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 将某个话题停止     update topic set is_stop = 1 where tid = '9b47f4d6-9f34-4485-9d4d-c4538e5d25a8'
func StopTopic(tid string) int {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	var maps = make(map[string]interface{})
	maps["is_stop"] = 1
	fmt.Println(tid)

	if err := db.Model(&Topic{}).Exec("update topic set is_stop = 1 where tid = ?", tid).Error; err != nil {
		fmt.Println("更新话题id  stop出错")
		return errmsg.ERROR
	}
	betxinredis.BatchDel("topic")

	return errmsg.SUCCSE
}

func GetTopicTotalPrice(tid string) (decimal.Decimal, int) {
	decimal.DivisionPrecision = 2
	var topic Topic
	if err := db.Model(&Topic{}).Where("tid = ?", tid).First(&topic).Error; err != nil {
		return topic.TotalPrice, errmsg.ERROR
	}
	return topic.TotalPrice, errmsg.SUCCSE
}

// GetCateArt 查询分类下的所有话题
func GetTopicByCid(cid int, limit int, offset int) ([]Topic, int, int) {
	var topicList []Topic
	var total int64
	err := db.Preload("Category").Limit(limit).Offset(offset).Where("cid =?", cid).Order("Created_At DESC").Find(&topicList).Error
	db.Model(&topicList).Where("cid =?", cid).Count(&total)
	if err != nil {
		return nil, errmsg.ERROR, 0
	}

	return topicList, int(total), errmsg.SUCCSE
}

// 根据uuid获取话题数据.  查询单个话题
func GetTopicById(tid string) (Topic, int) {
	var topic Topic
	err := db.Where("tid = ?", tid).Preload("Category").Joins("Category").Select("tid, cid, title, intro, collect_count, yes_ratio, no_ratio, yes_ratio_price, no_ratio_price, total_price, read_count,img_url, end_time, Category.category_name, Category.id, created_at, updated_at, is_stop").
		First(&topic).Error
	db.Model(&topic).Where("tid = ?", tid).UpdateColumn("read_count", gorm.Expr("read_count + ?", 1))
	if err != nil {
		return topic, errmsg.ERROR
	}
	return topic, errmsg.SUCCSE
}

// 创建新话题
func CreateTopic(data *Topic) int {
	if err := db.Create(data).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 根据tid删除标签
func DeleteTopic(tid string) int {
	if err := db.Where("tid = ?", tid).Delete(&Topic{}).Error; err != nil {
		return errmsg.ERROR
	}

	if err := db.Where("tid = ?", tid).Delete(&Collect{}).Error; err != nil {
		return errmsg.ERROR
	}

	if err := db.Where("tid = ?", tid).Delete(&UserToTopic{}).Error; err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSE
}

func UpdateTopic(tid string, data *Topic) int {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return errmsg.ERROR
	}

	// 锁住指定 id 记录
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&Topic{}).Where("tid = ?", tid).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	maps["cid"] = data.Cid
	maps["intro"] = data.Intro
	maps["title"] = data.Title
	maps["ImgUrl"] = data.ImgUrl
	maps["EndTime"] = data.EndTime
	maps["IsStop"] = data.IsStop

	if err := db.Model(&Topic{}).Where("tid = ?", tid).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 更新话题的总价钱
func UpdateTopicTotalPrice(tid string, selectWin string, plusPrice decimal.Decimal) int {
	// selectWin yes_ratio, no_ratio
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return errmsg.ERROR
	}

	// 锁住指定 id 记录
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&Topic{}).Where("tid = ?", tid).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}
	fmt.Println("更新话题总价钱")
	// db.Model(&Topic{}).Where("tid = ?", tid).Update("total_price", gorm.Expr("total_price + ?", plusPrice))
	var mutex sync.Mutex
	mutex.Lock()
	err := db.Exec("update topic set total_price = total_price + ? where tid = ?", plusPrice, tid).Error
	if err != nil {
		log.Fatal("error: ", err)
	}
	mutex.Unlock()
	if selectWin == "yes_win" {
		err = db.Exec("update topic set yes_ratio = (? + yes_ratio_price)/total_price where tid = ?", plusPrice, tid).Error
		if err != nil {
			log.Fatal("error: ", err)
		}
		err = db.Exec("update topic set yes_ratio_price = yes_ratio_price + ? where tid = ?", plusPrice, tid).Error
		if err != nil {
			log.Fatal("error: ", err)
		}
		err = db.Exec("update topic set no_ratio = no_ratio_price/total_price where tid = ?", tid).Error
		if err != nil {
			log.Fatal("error: ", err)
		}
	} else {
		err = db.Exec("update topic set no_ratio = (? + no_ratio_price)/total_price where tid = ?", plusPrice, tid).Error
		if err != nil {
			log.Fatal("error: ", err)
		}
		err = db.Exec("update topic set no_ratio_price = no_ratio_price + ? where tid = ?", plusPrice, tid).Error
		if err != nil {
			log.Fatal("error: ", err)
		}
		err = db.Exec("update topic set yes_ratio = yes_ratio_price/total_price where tid = ?", tid).Error
		if err != nil {
			log.Fatal("error: ", err)
		}
	}

	betxinredis.BatchDel("topic")

	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 减少话题的总价钱
func RefundTopicTotalPrice(data *Refund, selected string, fee decimal.Decimal) int {
	// selectWin yes_ratio, no_ratio
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return errmsg.ERROR
	}

	// 锁住指定 id 记录
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&Topic{}).Where("tid = ?", data.Tid).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}
	fmt.Println("更新话题总价钱")
	var mutex sync.Mutex
	mutex.Lock()
	err := db.Exec("update topic set total_price = total_price - ? where tid = ?", data.RefundPrice.Add(fee), data.Tid).Error
	if err != nil {
		log.Fatal("error: ", err)
	}
	mutex.Unlock()

	// 先扣除手续费
	if selected == "yes" {
		err = db.Exec("update topic set yes_ratio = (yes_ratio_price - ?)/total_price where tid = ?", data.RefundPrice.Add(fee), data.Tid).Error
		if err != nil {
			log.Fatal("error: ", err)
		}
		err = db.Exec("update topic set yes_ratio_price = yes_ratio_price - ? where tid = ?", data.RefundPrice.Add(fee), data.Tid).Error
		if err != nil {
			log.Fatal("error: ", err)
		}
		err = db.Exec("update topic set no_ratio = no_ratio_price/total_price where tid = ?", data.Tid).Error
		if err != nil {
			log.Fatal("error: ", err)
		}
	} else {
		err = db.Exec("update topic set no_ratio = (no_ratio_price - ?)/total_price where tid = ?", data.RefundPrice.Add(fee), data.Tid).Error
		if err != nil {
			log.Fatal("error: ", err)
		}
		err = db.Exec("update topic set no_ratio_price = no_ratio_price - ? where tid = ?", data.RefundPrice.Add(fee), data.Tid).Error
		if err != nil {
			log.Fatal("error: ", err)
		}
		err = db.Exec("update topic set yes_ratio = yes_ratio_price/total_price where tid = ?", data.Tid).Error
		if err != nil {
			log.Fatal("error: ", err)
		}
	}

	betxinredis.BatchDel("topic")

	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// GetArt 查询话题列表
func ListTopics(offset int, limit int) ([]Topic, int, int) {
	var topicList []Topic
	var err error
	var total int64
	err = db.Select("tid, cid, title, intro, collect_count, yes_ratio, no_ratio, yes_ratio_price, no_ratio_price, total_price, read_count,img_url, end_time, Category.category_name, Category.id, created_at, updated_at, is_stop").
		Limit(limit).Offset(offset).Order("Created_At DESC").Joins("Category").Find(&topicList).Error
	// 单独计数
	db.Model(&topicList).Count(&total)
	if err != nil {
		return nil, 0, errmsg.ERROR
	}
	return topicList, int(total), errmsg.SUCCSE
}

// 搜索标题
func SearchTopic(offset int, limit int, query interface{}, args ...interface{}) ([]Topic, int, int) {
	var topicList []Topic
	var err error
	var total int64
	err = db.Select("tid, cid, title, intro, collect_count, yes_ratio, no_ratio, yes_ratio_price, no_ratio_price, total_price, read_count,img_url, end_time, Category.category_name, Category.id, created_at, updated_at, is_stop").
		Order("Created_At DESC").Joins("Category").Where(query, args...).Limit(limit).Offset(offset).Find(&topicList).Count(&total).Error
	if err != nil {
		return nil, int(total), errmsg.ERROR
	}
	return topicList, int(total), errmsg.SUCCSE
}

// 查询某个话题的赢了的总钱数
func SearchTopicWinTopic(tid string, win string) (decimal.Decimal, int) {
	var winPrice decimal.Decimal
	var err error
	if win == "yes_win" {
		err = db.Raw("select SUM(yes_ratio_price) from topic where tid = ?", tid).Scan(&winPrice).Error
		if err != nil {
			return winPrice, errmsg.ERROR
		}
	} else {
		err = db.Raw("select SUM(no_ratio_price) from topic where tid = ?", tid).Scan(&winPrice).Error
		if err != nil {
			return winPrice, errmsg.ERROR
		}
	}
	return winPrice, errmsg.SUCCSE
}
