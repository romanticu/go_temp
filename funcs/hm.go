package funcs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type LogicType string

var AND LogicType = "and"
var OR LogicType = "or"

type OperatorType string

func (o *OperatorType) String() string {
	if o == nil {
		return ""
	}
	return string(*o)
}

var EqualNumOperator OperatorType = "eq-num"
var EqualStrOperator OperatorType = "eq-str"
var NotEqualStrOperator OperatorType = "neq-str"
var ContainOperator OperatorType = "contain"
var GreaterThanOperator OperatorType = "gt"
var IsHistoryOperator OperatorType = "is-history"
var NotHistoryOperator OperatorType = "not-history"

type DateOperatorType string

func (d *DateOperatorType) String() string {
	if d == nil {
		return ""
	}
	return string(*d)
}

var TodayOperator DateOperatorType = "today"
var YesterdayOperator DateOperatorType = "yesterday"
var TomorrowOperator DateOperatorType = "tomorrow"
var CurrentDayOperator DateOperatorType = "current-today"
var BeforeDayOperator DateOperatorType = "before-today"
var AfterDayOperator DateOperatorType = "after-today"
var ThisWeekOperator DateOperatorType = "this-week"
var LastWeekOperator DateOperatorType = "last-week"
var NextWeekOperator DateOperatorType = "next-week"
var CurrentWeekOperator DateOperatorType = "current-week"
var BeforeWeekOperator DateOperatorType = "before-week"
var AfterWeekOperator DateOperatorType = "after-week"
var ThisYearOperator DateOperatorType = "this-year"
var LastYearOperator DateOperatorType = "last-year"
var NextYearOperator DateOperatorType = "next-year"
var CurrentYearOperator DateOperatorType = "current-year"
var BeforeYearOperator DateOperatorType = "before-year"
var AfterYearOperator DateOperatorType = "after-year"
var SelectDateOperator DateOperatorType = "select-date"
var SelectMonthDayRangeOperator DateOperatorType = "select-month-day-range"
var SelectDateRangeOperator DateOperatorType = "select-date-range"
var FromTodayAddOperator DateOperatorType = "from-today-add"
var FromTodaySubOperator DateOperatorType = "from-today-sub"

type DateRangeType int

var Current DateRangeType = 0
var History DateRangeType = 1

var Nested = "nested"

type FilterReq struct {
	Filters Filter
}

type DateMode struct {
	IsOrNot     bool
	DateOptions []DateOperatorType
	DateValues  []string
}

type Filter struct {
	Field       *string
	Values      []string
	Operator    *OperatorType // Operator or nil
	Type        *string       // nested or nil
	Conjunction *LogicType    // Conjunction or nil
	DateRange   DateRangeType
	DateSelect  DateMode
	FilterSet   []Filter
}

type F = func(db *gorm.DB) *gorm.DB

var DBB *gorm.DB
var recIDs []string

func ConstructSql() {
	url := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=True",
		"root",
		"3682561289d64b63b24dd85f41776ba2",
		"192.168.88.201",
		"21306",
		"morpheus",
	)
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	DBB = db
	var companyType = "company_type"
	var uploadList = "upload_list"
	var regionName = "region_name"
	var name = "entity_primary_name"
	var foundedOn = "founded_on"
	var vertical = "vertical"
	fq := FilterReq{Filters: Filter{
		Conjunction: &AND,
		FilterSet: []Filter{
			{
				Field:     &companyType,
				Values:    []string{"高新技术", "科技中小企业"},
				Operator:  &ContainOperator,
				DateRange: Current,
			},
			{
				Type:        &Nested,
				Conjunction: &AND,
				FilterSet: []Filter{
					{
						Field:     &uploadList,
						Values:    []string{"1"},
						Operator:  &ContainOperator,
						DateRange: Current,
					},
					{
						Field:     &regionName,
						Values:    []string{"大连市甘井子区"},
						Operator:  &IsHistoryOperator,
						DateRange: History,
						DateSelect: DateMode{
							IsOrNot:       true,
							DateOptions: []DateOperatorType{TodayOperator, CurrentDayOperator},
							DateValues:    nil,
						},
					},
					{
						Field:     &name,
						Values:    []string{"京东方集团"},
						Operator:  &NotHistoryOperator,
						DateRange: History,
						DateSelect: DateMode{
							IsOrNot:       true,
							DateOptions: []DateOperatorType{ThisYearOperator, SelectMonthDayRangeOperator},
							DateValues:    []string{"10-01", "12-31"},
						},
					},
					{
						Field:     &foundedOn,
						Values:    nil,
						Operator:  nil,
						DateRange: Current,
						DateSelect: DateMode{
							IsOrNot:       true,
							DateOptions: []DateOperatorType{SelectDateOperator},
							DateValues:    []string{"2021-01-01"},
						},
					},
					{
						Type:        &Nested,
						Conjunction: &OR,
						FilterSet: []Filter{
							{
								Field:     &vertical,
								Values:    []string{"自动驾驶", "人工智能"},
								Operator:  &ContainOperator,
								DateRange: Current,
							},
							{
								Field:     &regionName,
								Values:    []string{"北京"},
								Operator:  &EqualStrOperator,
								DateRange: Current,
							},
						},
					},
				},
			},
		},
	}}
	fmt.Println(sql(fq.Filters.FilterSet, fq.Filters.Conjunction, "futian_table_id_fund_of_fund_investments"))
	// company_type in ("高新技术","科技中小企业") and (upload_list in ("1") and rec_id in ("046c9ad8-a461-42d0-9e47-57e474e68340") and founded_on = "2021-01-01")
}

func sql(filterSet []Filter, logic *LogicType, tableID string) string {

	var where []string
	for i := range filterSet {
		if filterSet[i].Type != nil && *filterSet[i].Type == Nested {
			whereStr := sql(filterSet[i].FilterSet, filterSet[i].Conjunction, tableID)
			where = append(where, "("+whereStr+")")
		} else {
			// 普通字段处理
			if filterSet[i].DateRange == Current && filterSet[i].Operator != nil {
				gcf := generalConstructionFunc[(*filterSet[i].Operator).String()]
				if gcf == nil {
					continue
				}
				w := gcf(*filterSet[i].Field, filterSet[i].Values)
				where = append(where, w)
			} else if filterSet[i].DateRange == Current && filterSet[i].Operator == nil {
				// 普通日期字段
				dcf := dateConstructionFunc[dateOperatorsToKey(filterSet[i].DateSelect.DateOptions)]
				if dcf == nil {
					continue
				}
				w := dcf(*filterSet[i].Field, filterSet[i].DateSelect)
				where = append(where, w)
			} else if filterSet[i].DateRange == History && filterSet[i].Operator != nil {
				// 历史数据
				dcf := dateConstructionFunc[dateOperatorsToKey(filterSet[i].DateSelect.DateOptions)]
				if dcf == nil {
					continue
				}
				mWhere := dcf("m.updated_at", filterSet[i].DateSelect)
				fmt.Println(mWhere)
				rIDs := HistoryOperatorConstruction(tableID, *filterSet[i].Field, mWhere, filterSet[i].Values, *filterSet[i].Operator)
				fmt.Println(rIDs)
				if len(rIDs) != 0 {
					where = append(where, "rec_id in ("+arrStrToSqlStr(rIDs)+")")
				}

			}

		}
	}

	return strings.Join(where, string(" "+*logic+" "))
}

var generalConstructionFunc map[string]func(field string, values []string) string
var dateConstructionFunc map[string]func(field string, dateSelect DateMode) string

func init() {
	generalConstructionFunc = make(map[string]func(field string, values []string) string)
	dateConstructionFunc = make(map[string]func(field string, dateSelect DateMode) string)

	generalConstructionFunc[ContainOperator.String()] = ContainOperatorConstruction
	generalConstructionFunc[EqualStrOperator.String()] = EqualStrOperatorConstruction

	dateConstructionFunc[SelectDateOperator.String()] = SelectDateOperatorConstruction
	dateConstructionFunc[TodayOperator.String()+CurrentDayOperator.String()] = TodayCurrentDayOperatorConstruction
	dateConstructionFunc[ThisYearOperator.String()+SelectMonthDayRangeOperator.String()] = ThisYearSelectMonthDayOperatorConstruction
}

func ContainOperatorConstruction(field string, values []string) string {
	return fmt.Sprintf("%s in (%s)", field, arrStrToSqlStr(values))
}

func EqualStrOperatorConstruction(field string, values []string) string {
	return fmt.Sprintf("%s = \"%s\"", field, values[0])
}

func SelectDateOperatorConstruction(field string, dateSelect DateMode) string {
	if dateSelect.IsOrNot {
		return fmt.Sprintf("%s = %s", field, arrStrToSqlStr(dateSelect.DateValues))
	}

	return fmt.Sprintf("%s != %s", field, arrStrToSqlStr(dateSelect.DateValues))
}

func TodayCurrentDayOperatorConstruction(field string, dateSelect DateMode) string {
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	t2 := t1.AddDate(0, 0, 1)
	if dateSelect.IsOrNot {
		return fmt.Sprintf(`(%s >= "%s" and %s < "%s")`, field, t1.String(), field, t2.String())
	}
	return fmt.Sprintf(`(%s < "%s" and %s >= "%s")`, field, t1.String(), field, t2.String())
}

func ThisYearSelectMonthDayOperatorConstruction(field string, dateSelect DateMode) string {
	dateArr := strings.Split(dateSelect.DateValues[0], "-")
	sMonth, err := strconv.Atoi(dateArr[0])
	if err != nil {
		fmt.Println(err)
	}
	sDay, err := strconv.Atoi(dateArr[1])
	if err != nil {
		fmt.Println(err)
	}
	t := time.Now()
	t1 := time.Date(t.Year(), time.Month(sMonth), sDay, 0, 0, 0, 0, t.Location())
	dateArr = strings.Split(dateSelect.DateValues[1], "-")
	eMonth, err := strconv.Atoi(dateArr[0])
	if err != nil {
		fmt.Println(err)
	}
	eDay, err := strconv.Atoi(dateArr[1])
	if err != nil {
		fmt.Println(err)
	}
	t2 := time.Date(t.Year(), time.Month(eMonth), eDay, 0, 0, 0, 0, t.Location())
	if dateSelect.IsOrNot {
		return fmt.Sprintf(`(%s >= "%s" and %s < "%s")`, field, t1.String(), field, t2.String())
	}
	return fmt.Sprintf(`(%s < "%s" and %s >= "%s")`, field, t1.String(), field, t2.String())
}

type mf struct {
	RecordID string
}

func HistoryOperatorConstruction(tableID, field, dateWhere string, values []string, history OperatorType) []string {
	var ids []string
	var rows []mf
	var op = "="
	if history != IsHistoryOperator {
		op = "!="
	}
	err := DBB.Debug().Raw(`
			SELECT MAX(updated_at), record_id FROM modifications m where table_id = ? and column_id = ? and m.after `+op+` ? and 
	`+dateWhere+" group BY record_id ", tableID, field, values[0]).Scan(&rows).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, item := range rows {
		ids = append(ids, item.RecordID)
	}

	return ids
}

func arrStrToSqlStr(values []string) string {
	if len(values) == 0 {
		return ""
	}
	b, _ := json.Marshal(values)
	s := string(b)
	fmt.Println(s)
	return s[1 : len(s)-1]
}

func dateOperatorsToKey(dateOperators []DateOperatorType) string {
	var str string
	for _, item := range dateOperators {
		str += string(item)
	}

	return str
}

func strOperatorConstruction(op OperatorType) string {
	switch op {
	case EqualStrOperator:
		return ""
	}

	return ""
}
