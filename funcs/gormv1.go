package funcs

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitGormV1() *gorm.DB {
	connStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=True",
		"root", "root", "192.168.88.11", "3306", "pevc")
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		fmt.Println(err)
	}
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(20)
	db.LogMode(true)
	return db
}
