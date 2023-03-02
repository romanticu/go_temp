package funcs

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type BasicModel struct {
	ID        int        `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `gorm:"index"`
}
type InstitutionalInvestorADSHUpdate struct {
	BasicModel
	EntityID string `gorm:"type:VARCHAR(100);unique_index"`
}

var DB *gorm.DB

func initDB() {
	connStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=True",
		"root", "root", "192.168.88.11", "3306", "pevc")
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		fmt.Println(err)
	}
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(20)
	db.LogMode(true)
	DB = db
}

func initUpdateDB() {
	connStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=True",
		"root", "root", "192.168.88.11", "3306", "pevc_update")
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		fmt.Println(err)
	}
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(20)
	db.LogMode(true)
	DB = db
}

func TestInsert() {
	initDB()
	DB.Save(&InstitutionalInvestorADSHUpdate{
		BasicModel: BasicModel{
			CreatedAt: time.Now(),
		},
		EntityID: "6666",
	})
	DB.Save(&InstitutionalInvestorADSHUpdate{
		BasicModel: BasicModel{
			CreatedAt: time.Now(),
		},
		EntityID: "8888",
	})
}

func TestDelete() {
	initDB()
	ids := []string{"6666", "8888"}
	err := DB.Unscoped().Where("entity_id IN (?)", ids).Delete(&InstitutionalInvestorADSHUpdate{}).Error
	if err != nil {
		fmt.Println(err)
	}
}

type EntityRecord struct {
	ID              int       `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt       time.Time `gorm:"not null;index"`
	HasBeenConsumed int       `gorm:"type:TINYINT(1);index;not null"`
	EntityID        string    `gorm:"type:VARCHAR(100)"`
}

func TestBatchUpdates() {
	initUpdateDB()
	recordIDs := []int{1, 2, 3, 4}
	DB.Table("institutional_investor_delta").Where("id IN (?)", recordIDs).
		Updates(EntityRecord{
			HasBeenConsumed: 1,
		})
}

func TestBatchDel() {
	initDB()
	err := DB.Exec(`TRUNCATE TABLE vertical_deals`).Error

	if err != nil {
		fmt.Println(err)
	}
}

func TestPluck() {
	initDB()
	var ids []string
	err := DB.Table("organizations").Pluck("id", &ids).Error
	if err != nil {
		fmt.Println(err)
	}
}

func Test_() {
	fmt.Println("Test")
}

func TestAMAC() {
	initDB()
	var c int
	err := DB.Table("amac_disclosed_members adm").
		Select(`
			adm.*,
			pemf.*`,
		).Joins("left join private_equity_management_firms pemf on pemf.entity_id = adm.org_id and pemf.deleted_at is null").
		Where("pemf.entity_type = ?", 109006022).
		Where("pemf.parent_institutional_investor_entity_type=?", 109006022).
		Where("adm.person_id = ? and adm.deleted_at is null", 666).Count(&c).Error
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
}


