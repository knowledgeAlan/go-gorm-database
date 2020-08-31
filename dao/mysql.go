package dao

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Like struct {
	ID        int    `gorm:"primary_key"`
	Ip        string `gorm:"type:varchar(20) ; not null"`
	Ua        string `gorm:"type:varchar(256);"`
	Title     string `gorm:"type:varchar(128) ; not null"`
	Hash      uint64 `gorm:"unique_index:hash_idx;"`
	CreatedAt time.Time
}

func InitMySql() (err error) {
	fmt.Println("init mysql....")

	if db == nil {
		db, err = gorm.Open("mysql", "root:123456@/gorm_demo?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			fmt.Println("init mysql.... err:", err)

			return err
		}
		db.DB().SetMaxOpenConns(10)

	}

	//CreateLike()
	return
}

func CreateLike() error {

	db.Table("like").CreateTable(&Like{})
	return nil
}

func SaveLike(like *Like) error {

	if err := db.Table("like").Create(like).Error; err != nil {
		fmt.Println("save like  error", err)
		return err
	}
	return nil
}

func DeleteLike(like *Like) error {

	if err := db.Table("like").Delete(Like{}).Error; err != nil {
		fmt.Println("delete like error", err)
		return err
	}
	return nil
}

func SelectLike(like *Like) (likeList []Like, err error) {

	if err := db.Table("like").Find(&likeList).Error; err != nil {
		fmt.Println("select like error", err)
		return nil, err
	}
	return likeList, nil
}
