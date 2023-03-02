package funcs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gitlab.mvalley.com/datapack/cain/pkg/essentials"
)

// Recruitment 招聘职位信息表，对应 NJ 的 6497 表
type Recruitment struct {
	essentials.BasicModel

	CompanyID *string `gorm:"type:VARCHAR(100);index"`
	// 公司名称
	CompanyName string `gorm:"type:VARCHAR(500);not null"`
	// 招聘渠道名称
	SourceName string `gorm:"type:VARCHAR(255)"`
	// 首次招聘发布时间
	FirstPublishOn string `gorm:"type:VARCHAR(8)"`
	// 招聘发布时间
	PublishOn string `gorm:"type:VARCHAR(8)"`
	// 职位名称
	JobName string `gorm:"type:VARCHAR(500)"`
	// 职位关键词
	JobKeyword string `gorm:"type:VARCHAR(2000)"`
	// 职位 ID
	JobID int `gorm:"type:INT"`
	// 岗位职责
	JobDescription string `gorm:"type:VARCHAR(4000)"`
	// 职能类型
	JobType string `gorm:"type:VARCHAR(100)"`
	// 职能类别
	JobFunction string `gorm:"type:VARCHAR(100)"`
	// 所属部门
	Department string `gorm:"type:VARCHAR(500)"`
	// 招聘人数
	HeaderCount int `gorm:"type:INT"`
	// 工作地点
	LocationName string `gorm:"type:VARCHAR(1000)"`
	// 地区代码
	LocationCode string `gorm:"type:VARCHAR(10)"`
	// 学历要求
	EducationRequirement string `gorm:"type:VARCHAR(100)"`
	// 工作经验要求
	WorkingRequirement string `gorm:"type:VARCHAR(100)"`
	// 工作经验要求
	WorkingRequirementFormat float64 `gorm:"type:FLOAT(10,4)"`
	// 语言要求
	LanguageRequirement string `gorm:"type:VARCHAR(100)"`
	// 最低年龄
	MinAge int `gorm:"type:INT"`
	// 最高年龄
	MaxAge int `gorm:"type:INT"`
	// 最低薪水
	MinCompensation string `gorm:"type:VARCHAR(100)"`
	// 最高薪水
	MaxCompensation string `gorm:"type:VARCHAR(100)"`
	// 福利
	Allowance string `gorm:"type:VARCHAR(500)"`
}

func MatchYear() {
	db, err := InitGormV2()
	if err != nil {
		fmt.Println(err)
		return
	}
	var recs = make([]Recruitment, 0)
	err = db.Model(&Recruitment{}).Find(&recs).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("len: ", len(recs))

	noExpRe := regexp.MustCompile("(不限|无|应届|在校)")
	oneYearExpRe := regexp.MustCompile("(一年以下|1年以内|1年以下|个月)")
	otherYearExpRe := regexp.MustCompile(`^(\d+)(年经验|年以上|年及以上|年以上经验)?$`)
	rangeYearExpRe := regexp.MustCompile(`(\d+)-(\d+)年*`)
	//subMatch := rangeYearExpRe.FindStringSubmatch("5-10年")
	//start, _ := strconv.ParseFloat(subMatch[1], 64)
	//end, _ := strconv.ParseFloat(subMatch[2], 64)
	//fmt.Println(subMatch)
	//fmt.Println((start + end) / 2)
	noMatch := make([]string, 0)
	s := time.Now()
	fmt.Println(s)
	for i, item := range recs {
		if noExpRe.MatchString(item.WorkingRequirement) || item.WorkingRequirement == "" {
			recs[i].WorkingRequirementFormat = 0
		} else if oneYearExpRe.MatchString(item.WorkingRequirement) {
			recs[i].WorkingRequirementFormat = 0.5
		} else if subMatch := otherYearExpRe.FindStringSubmatch(item.WorkingRequirement); len(subMatch) > 1 {
			year, err := strconv.ParseFloat(subMatch[1], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			recs[i].WorkingRequirementFormat = year
		} else if subMatch := rangeYearExpRe.FindStringSubmatch(item.WorkingRequirement); len(subMatch) > 1 {
			start, err := strconv.ParseFloat(subMatch[1], 64)
			end, err := strconv.ParseFloat(subMatch[2], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}

			recs[i].WorkingRequirementFormat = (start + end) / 2
			//if item.WorkingRequirement == "5-10年" {
			//	fmt.Println((start + end) / 2)
			//	fmt.Println(recs[i].WorkingRequirementFormat)
			//}
		} else {
			noMatch = append(noMatch, recs[i].WorkingRequirement)
		}
	}

	elapsed := time.Since(s)
	fmt.Println("used: ", elapsed)

	fmt.Println(noMatch)
	//err = db.Model(&Recruitment{}).Where("1=1").Delete(&Recruitment{}).Error
	//if err != nil {
	//	fmt.Println(err)
	//}
	//err = db.CreateInBatches(&recs, 1000).Error
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func MatchCity() {
	// db, err := InitGormV2()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// var recs = make([]Recruitment, 0)
	// err = db.Model(&Recruitment{}).Find(&recs).Error
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("len: ", len(recs))
	s := " "
	locationExpRe := regexp.MustCompile("^([\u4e00-\u9fa5]+)-([\u4e00-\u9fa5]+)$")
	fmt.Println(locationExpRe.FindStringSubmatch(s))
	//fmt.Println(locationExpRe.FindStringSubmatch("上海"))

	//noMatch := make([]string, 0)
	// for i, item := range recs {
	// 	if subMatch := locationExpRe.FindStringSubmatch(item.LocationName); len(subMatch) > 1 {
	// 		recs[i].LocationName = subMatch[1]
	// 	}
	// }
	// err = db.Model(&Recruitment{}).Where("1=1").Delete(&Recruitment{}).Error
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// err = db.CreateInBatches(&recs, 1000).Error
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

var MatchEnglishCompanyName = regexp.MustCompile("^[a-zA-Z0-9]+$")

func IsAbroadCompany(name string) bool {
	afterName := strings.ReplaceAll(name, " ", "")
	return MatchEnglishCompanyName.MatchString(afterName)
}
