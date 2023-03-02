package funcs

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	UserID string `gorm:"type:VARCHAR(100);primaryKey"`
	Name   string `gorm:"type:VARCHAR(100)"`
	Age    int    `gorm:"type:INT"`
}

func InitGormV2() (*gorm.DB, error) {
	url := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=True",
		"root",
		"root",
		"192.168.88.11",
		"3306",
		"gfp",
	)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)

	return db, nil
}
