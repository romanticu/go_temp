

var LogicAnd = "and"
var LogicOr = "or"

/*  
    number operator
*/
const EqualNumOperator = "eq-num"
const NotEqualNumOperator = "neq-num"
const GreaterThanNumOperator = "gt-num"
const GreaterThanEqualNumOperator = "gte-num"
const LessThanNumOperator = "lt-num"
const LessThanEqualNumOperator = "lte-num"
const IsEmptyNumOperator = "is-empty-num"
const IsNotEmptyNumOperator = "is-not-empty-num"

/*  
    text operator
*/
const ContainStrOperator = "contain-str"
const DoesNotContainStrOperator = "does-not-contain-str"
const EqualStrOperator = "eq-str"
const NotEqualStrOperator = "neq-str"
const IsEmptyStrOperator = "is-empty-str"
const IsNotEmptyStrOperator = "is-not-empty-str"


/* 
    muti tag operator
*/
// 包含任意
const HasAnyOfMTagOperator = "has-any-of-mtag"
// 包含全部
const HasAllOfMTagOperator = "has-all-of-mtag"
// 完全等于
const IsExactlyMTagOperator = "is-exactly-mtag"
// 不包含任意
const HasNoneOfMTagOperator = "has-none-of-mtag"
const IsEmptyMTagOperator = "is-empty-mtag"
const IsNotEmptyMTagOperator = "is-not-empty-mtag"

/*
    single tag operator
*/
const IsAnyOfSTagOperator = "is-any-of-stag"
const IsNoneOfSTagOperator = "is-none-of-stag"
const IsEmptySTagOperator = "is-empty-stag"
const IsNotEmptySTagOperator = "is-not-empty-stag"

/*
    date operator
*/

const EqualsTimeOperator = "eq-time"
const NotEqualsTimeOperator = "neq-time"

const TodayOption = "today"
const YesterdayOption = "yesterday"
const TomorrowOption = "tomorrow"
const CurrentTodayOption = "current-today"
const BeforeTodayOption = "before-today"
const AfterTodayOption = "after-today"


const ThisWeekOption = "this-week"
const LastWeekOption = "last-week"
const NextWeekOption = "next-week"
const CurrentWeekOption = "current-week"
const BeforeWeekOption = "before-week"
const AfterWeekOption = "after-week"

const ThisMonthOption = "this-month"
const LastMonthOption = "last-month"
const NextMonthOption = "next-month"
const CurrentMonthOption = "current-month"
const BeforeMonthOption = "before-month"
const AfterMonthOption = "after-month"

const ThisQuarterOption = "this-quarter"
const LastQuarterOption = "last-quarter"
const NextQuarterOption = "next-quarter"
const CurrentQuarterOption = "current-quarter"
const BeforeQuarterOption = "before-quarter"
const AfterQuarterOption = "after-quarter"

const ThisYearOption = "this-year"
const LastYearOption = "last-year"
const NextYearOption = "next-year"
const CurrentYearOption = "current-year"
const BeforeYearOption = "before-year"
const AfterYearOption = "after-year"

const FromTodayAddOption = "from-today-add"

const FromTodaySubOption = "from-today-sub"

const SelectMonthDayRangeOption = "select-month-day-range"

const SelectDateOption = "select-date"

/**
 *  const
 */
const FilterNested = "nested"
const CurrentDateRange = "current"
const HistoryDateRange = "history"
const CollisionFieldType = "collision"
const EntityFieldType = "entity"
const LinkFieldType = "link"
const GeneralFieldType = "general"
const DayTimestamp = 24*60*60*1000
const WeekTimestamp = 7*24*60*60*1000
/**
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

type DateMode struct {
    Operator *OperatorType
	DateOptions []DateOperatorType
	DateValues  []string
}
 */

function filterDemo() {
    let companyType = "company_type"
    let uploadList = "upload_list"
    let regionName = "region_name"
    let investmentTarget = "investment_target_rime_entity"
	let name = "entity"
	let foundedOn = "founded_on"
	let vertical = "vertical"
    let f = {
        Conjunction: LogicOr,
        FilterSet: [
            {
                Field: companyType,
                FieldType: GeneralFieldType,
                Values:    ["高新技术", "科技中小企业"],
                Operator: HasAnyOfMTagOperator,
                DateRange: CurrentDateRange,
            },
            {
                Type:        FilterNested,
				Conjunction: LogicOr,
                FilterSet: [
                    {
						Field:     uploadList,
						Values:    ["1"],
						Operator:  ContainStrOperator,
						DateRange: CurrentDateRange,
					},
                    {
						Field:     regionName,
                        FieldType: CollisionFieldType,
						Values:    ["深圳市福田区"],
						Operator:  EqualStrOperator,
						DateRange: CurrentDateRange,
						DateSelect: {
                            Operator: EqualsTimeOperator,
							DateOptions: [TodayOption, CurrentDayOption],
						},
					},
                    {
						Field:     name,
                        FieldType: EntityFieldType,
						Values:    ["京东方集团"],
						Operator:  NotEqualStrOperator,
						DateRange: HistoryDateRange,
						DateSelect: {
                            Operator: EqualsTimeOperator,
							DateOptions: [ThisYearOption, SelectMonthDayRangeOption],
							DateValues:    ["10-01", "12-31"],
						},
					},
                    {
						Field:     foundedOn,
                        FieldType: CollisionFieldType,
						DateRange: CurrentDateRange,
						DateSelect: {
                            Operator: EqualsTimeOperator,
							DateOptions: [SelectDateOption],
							DateValues:    ["2021-01-01"],
						},
					},
                    {
						Type:        FilterNested,
						Conjunction: LogicOr,
						FilterSet: [
							{
								Field:     vertical,
                                FieldType: CollisionFieldType,
								Values:    ["自动驾驶", "人工智能"],
								Operator:  HasAnyOfMTagOperator,
								DateRange: CurrentDateRange,
							},
							{
								Field:     investmentTarget,
                                FieldType: LinkFieldType,
								Values:    ["大米科技"],
								Operator:  EqualStrOperator,
								DateRange: CurrentDateRange,
							},
                        ],
					},
                ]
            }
        ]
    }

    console.log(searchRecords(f.FilterSet, f.Conjunction))
}

function searchRecords(filterSet, logic) {
    let recIDs = []
    let whereIDs = []
    for (let i = 0; i < filterSet.length; i++) {
        if (filterSet[i].Type != undefined && filterSet[i].Type == FilterNested) {
            let ids = searchRecords(filterSet[i].FilterSet, filterSet[i].Conjunction)
            whereIDs.push(ids)
        } else {
            // 普通字段处理
			if (filterSet[i].DateRange == CurrentDateRange && filterSet[i].Operator != undefined) {
                let gfh = generalFieldHandlerMap[filterSet[i].Operator]
                if (gfh == undefined) {
                    continue
                }
                let ids = dfh(filterSet[i].Field, filterSet[i].FieldType, filterSet[i].Values)
                whereIDs.push(ids)
            
            // 普通日期字段 func(field, fieldType, dateSelect)
            } else if (filterSet[i].DateRange == CurrentDateRange && filterSet[i].Operator == undefined &&
                filterSet[i].DateSelect != undefined && filterSet[i].DateSelect.Operator != undefined) {
                let dh = dateHandlerMap[filterSet[i].DateSelect.DateOptions.join("")]
                if (dh != undefined) {
                    continue
                }
                let ids = dh(filterSet[i].Field, filterSet[i].FieldType, filterSet[i].DateSelect, recordsRaw, false)
                whereIDs.push(ids)

            // 历史数据范围搜索，先进行条件搜索，再进行时间搜索
            } else if (filterSet[i].DateRange == HistoryDateRange && filterSet[i].Operator != undefined && 
                filterSet[i].DateSelect != undefined && filterSet[i].DateSelect.Operator != undefined) {
                let gfh = generalFieldHandlerMap[filterSet[i].Operator]
                if (gfh != undefined) {
                    continue
                }
                // todo historyRecords
                let historyRecords = dfh(filterSet[i].Field, filterSet[i].FieldType, filterSet[i].Values)

                let dh = dateHandlerMap[filterSet[i].DateSelect.DateOptions.join("")]
                if (dh != undefined) {
                    continue
                }
                let ids = dh(filterSet[i].Field, filterSet[i].FieldType, filterSet[i].DateSelect, historyRecords, true)
                whereIDs.push(ids)
                
            }
        }
    }

    if (logic == LogicAnd) {
        recIDs = andWhereHandler(whereIDs)
    } else {
        recIDs = orWhereHandler(whereIDs)
    }
    return recIDs
}

/**
 * 
 * 对 whereIDs 二维数组取交集，二维数组中的一个数组即代表一个条件
 * @param {[][]int} whereIDs
 * @returns 
 */
function andWhereHandler(whereIDs) {
    let idsMap = {}
    let len = whereIDs.length
    let recIDs = []
    whereIDs.forEach((ids, i) => {
        ids.forEach((id, j) => {
            if (idsMap[id] == undefined) {
                idsMap[id] = 1
            } else {
                idsMap[id]++
            }
        })
    })

    for (let k in idsMap) {
        if (idsMap[k] == len) {
            recIDs.push(k)
        }
    }

    return recIDs
}

/**
 * 
 * 对 whereIDs 二维数组取并集，二维数组中的一个数组即代表一个条件
 * @param {[][]int} whereIDs
 * @returns 
 */
function orWhereHandler(whereIDs) {
    let idsMap = {}
    let len = whereIDs.length
    let recIDs = []
    whereIDs.forEach((ids, i) => {
        ids.forEach((id, j) => {
            idsMap[id] = 1
        })
    })

    for (let k in idsMap) {
        recIDs.push(k)
    }

    return recIDs
}

var generalFieldHandlerMap = {}
var dateHandlerMap = {}

function initFuncMap() {
    generalFieldHandlerMap[EqualStrOperator] = equalStrOperatorHandler



    dateHandlerMap[TodayOption+CurrentTodayOption] = dayOptionHandler
    dateHandlerMap[TodayOption+BeforeTodayOption] = dayOptionHandler
    dateHandlerMap[TodayOption+AfterTodayOption] = dayOptionHandler
    dateHandlerMap[YesterdayOption+CurrentTodayOption] = dayOptionHandler
    dateHandlerMap[YesterdayOption+BeforeTodayOption] = dayOptionHandler
    dateHandlerMap[YesterdayOption+AfterTodayOption] = dayOptionHandler
    dateHandlerMap[TomorrowOption+CurrentTodayOption] = dayOptionHandler
    dateHandlerMap[TomorrowOption+BeforeTodayOption] = dayOptionHandler
    dateHandlerMap[TomorrowOption+AfterTodayOption] = dayOptionHandler

    dateHandlerMap[ThisWeekOption+CurrentWeekOption] = weekOptionHandler
    dateHandlerMap[ThisWeekOption+BeforeWeekOption] = weekOptionHandler
    dateHandlerMap[ThisWeekOption+AfterWeekOption] = weekOptionHandler
    dateHandlerMap[LastWeekOption+CurrentWeekOption] = weekOptionHandler
    dateHandlerMap[LastWeekOption+BeforeWeekOption] = weekOptionHandler
    dateHandlerMap[LastWeekOption+AfterWeekOption] = weekOptionHandler
    dateHandlerMap[NextWeekOption+CurrentWeekOption] = weekOptionHandler
    dateHandlerMap[NextWeekOption+BeforeWeekOption] = weekOptionHandler
    dateHandlerMap[NextWeekOption+AfterWeekOption] = weekOptionHandler

    dateHandlerMap[ThisMonthOption+CurrentMonthOption] = monthOptionHandler
    dateHandlerMap[ThisMonthOption+BeforeMonthOption] = monthOptionHandler
    dateHandlerMap[ThisMonthOption+AfterMonthOption] = monthOptionHandler
    dateHandlerMap[LastMonthOption+CurrentMonthOption] = monthOptionHandler
    dateHandlerMap[LastMonthOption+BeforeMonthOption] = monthOptionHandler
    dateHandlerMap[LastMonthOption+AfterMonthOption] = monthOptionHandler
    dateHandlerMap[NextMonthOption+CurrentMonthOption] = monthOptionHandler
    dateHandlerMap[NextMonthOption+BeforeMonthOption] = monthOptionHandler
    dateHandlerMap[NextMonthOption+AfterMonthOption] = monthOptionHandler

    dateHandlerMap[ThisQuarterOption+CurrentQuarterOption] = quarterOptionHandler
    dateHandlerMap[ThisQuarterOption+BeforeQuarterOption] = quarterOptionHandler
    dateHandlerMap[ThisQuarterOption+AfterQuarterOption] = quarterOptionHandler
    dateHandlerMap[LastQuarterOption+CurrentQuarterOption] = quarterOptionHandler
    dateHandlerMap[LastQuarterOption+BeforeQuarterOption] = quarterOptionHandler
    dateHandlerMap[LastQuarterOption+AfterQuarterOption] = quarterOptionHandler
    dateHandlerMap[NextQuarterOption+CurrentQuarterOption] = quarterOptionHandler
    dateHandlerMap[NextQuarterOption+BeforeQuarterOption] = quarterOptionHandler
    dateHandlerMap[NextQuarterOption+AfterQuarterOption] = quarterOptionHandler

    dateHandlerMap[ThisYearOption+CurrentYearOption] = yearCBAOptionHandler
    dateHandlerMap[ThisYearOption+BeforeYearOption] = yearCBAOptionHandler
    dateHandlerMap[ThisYearOption+AfterYearOption] = yearCBAOptionHandler
    dateHandlerMap[LastYearOption+CurrentYearOption] = yearCBAOptionHandler
    dateHandlerMap[LastYearOption+BeforeYearOption] = yearCBAOptionHandler
    dateHandlerMap[LastYearOption+AfterYearOption] = yearCBAOptionHandler
    dateHandlerMap[NextYearOption+CurrentYearOption] = yearCBAOptionHandler
    dateHandlerMap[NextYearOption+BeforeYearOption] = yearCBAOptionHandler
    dateHandlerMap[NextYearOption+AfterYearOption] = yearCBAOptionHandler

    dateHandlerMap[ThisYearOption+SelectDateOption+CurrentTodayOption] = yearMonthDayCBAOptionHandler
    dateHandlerMap[ThisYearOption+SelectDateOption+BeforeTodayOption] = yearMonthDayCBAOptionHandler
    dateHandlerMap[ThisYearOption+SelectDateOption+AfterTodayOption] = yearMonthDayCBAOptionHandler
    dateHandlerMap[LastYearOption+SelectDateOption+CurrentTodayOption] = yearMonthDayCBAOptionHandler
    dateHandlerMap[LastYearOption+SelectDateOption+BeforeTodayOption] = yearMonthDayCBAOptionHandler
    dateHandlerMap[LastYearOption+SelectDateOption+AfterTodayOption] = yearMonthDayCBAOptionHandler
    dateHandlerMap[NextYearOption+SelectDateOption+CurrentTodayOption] = yearMonthDayCBAOptionHandler
    dateHandlerMap[NextYearOption+SelectDateOption+BeforeTodayOption] = yearMonthDayCBAOptionHandler
    dateHandlerMap[NextYearOption+SelectDateOption+AfterTodayOption] = yearMonthDayCBAOptionHandler

    dateHandlerMap[ThisYearOption+SelectMonthDayRangeOption] = yearMonthDayRangeOptionHandler
    dateHandlerMap[ThisYearOption+SelectMonthDayRangeOption] = yearMonthDayRangeOptionHandler
    dateHandlerMap[ThisYearOption+SelectMonthDayRangeOption] = yearMonthDayRangeOptionHandler
    dateHandlerMap[LastYearOption+SelectMonthDayRangeOption] = yearMonthDayRangeOptionHandler
    dateHandlerMap[LastYearOption+SelectMonthDayRangeOption] = yearMonthDayRangeOptionHandler
    dateHandlerMap[LastYearOption+SelectMonthDayRangeOption] = yearMonthDayRangeOptionHandler
    dateHandlerMap[NextYearOption+SelectMonthDayRangeOption] = yearMonthDayRangeOptionHandler
    dateHandlerMap[NextYearOption+SelectMonthDayRangeOption] = yearMonthDayRangeOptionHandler
    dateHandlerMap[NextYearOption+SelectMonthDayRangeOption] = yearMonthDayRangeOptionHandler

    dateHandlerMap[FromTodayAddOption+CurrentTodayOption] = fromTodayOptionHandler
    dateHandlerMap[FromTodayAddOption+BeforeTodayOption] = fromTodayOptionHandler
    dateHandlerMap[FromTodayAddOption+AfterTodayOption] = fromTodayOptionHandler
    dateHandlerMap[FromTodaySubOption+CurrentTodayOption] = fromTodayOptionHandler
    dateHandlerMap[FromTodaySubOption+BeforeTodayOption] = fromTodayOptionHandler
    dateHandlerMap[FromTodaySubOption+AfterTodayOption] = fromTodayOptionHandler

    dateHandlerMap[SelectDateOption+CurrentTodayOption] = selectDateOptionHandler
    dateHandlerMap[SelectDateOption+BeforeTodayOption] = selectDateOptionHandler
    dateHandlerMap[SelectDateOption+AfterTodayOption] = selectDateOptionHandler

    dateHandlerMap[SelectMonthDayRangeOption] = selectDateRangeOptionHandler
}

function equalStrOperatorHandler(field, fieldType, values) {
    let recIds = []
    let val = values[0]
    records.forEach((item, i) => {
        let recVal = getFieldValue(item, field, fieldType)
        if (recVal == val){
            recIds.push(item.recordId)
        }
    })
    return recIds
}

function dayOptionHandler(field, fieldType, dateSelect, records, findHistory) {
    let startTimestamp, endTimestamp
    let option0, option1
    let recIds = []
    option0 = dateSelect.DateOptions[0]
    option1 = dateSelect.DateOptions[1]
    if (option0 == TodayOption) {
        startTimestamp = new Date(new Date().setHours(0, 0, 0, 0)).getTime()
    } else if (option0 == YesterdayOption) {
        startTimestamp = new Date(new Date(new Date().setHours(0, 0, 0, 0)).getTime() - DayTimestamp).getTime()
    } else if (option0 == TomorrowOption) {
        startTimestamp = new Date(new Date(new Date().setHours(0, 0, 0, 0)).getTime() + DayTimestamp).getTime()
    }

    endTimestamp = new Date(startTimestamp + DayTimestamp)

    records.forEach((item) => {
        // todo item.values 
        if (findHistory && matchCBAHistory(dateSelect.Operator, option1, item.values, startTimestamp, endTimestamp, CurrentTodayOption, BeforeTodayOption, AfterTodayOption)) {
            recIds.push(item.recordId)
            return
        }
        let recVal = getFieldValue(item, field, fieldType)
        if (matchCBA(dateSelect.Operator, option1, recVal, startTimestamp, endTimestamp, CurrentTodayOption, BeforeTodayOption, AfterTodayOption)) {
            recIds.push(item.recordId)
        }        
    })
    return recIds
}

function weekOptionHandler(field, fieldType, dateSelect, records, findHistory) {
    let startTimestamp, endTimestamp
    let option0, option1
    let recIds = []
    option0 = dateSelect.DateOptions[0]
    option1 = dateSelect.DateOptions[1]

    if (option0 == ThisWeekOption) {
        startTimestamp = getThisWeekMonday()
    } else if (option0 == LastWeekOption) {
        startTimestamp = new Date(getThisWeekMonday() - WeekTimestamp).getTime()
    } else if (option0 == NextWeekOption) {
        startTimestamp = new Date(getThisWeekMonday() + WeekTimestamp).getTime()
    }

    endTimestamp = new Date(startTimestamp + WeekTimestamp)

    records.forEach((item) => {
        // todo item.values 
        if (findHistory && matchCBAHistory(dateSelect.Operator, option1, item.values, startTimestamp, endTimestamp, CurrentWeekOption, BeforeWeekOption, AfterWeekOption)) {
            recIds.push(item.recordId)
            return
        }
        let recVal = getFieldValue(item, field, fieldType)
        if (matchCBA(dateSelect.Operator, option1, recVal, startTimestamp, endTimestamp, CurrentWeekOption, BeforeWeekOption, AfterWeekOption)) {
            recIds.push(item.recordId)
        }
    })

    return recIds
}

function monthOptionHandler(field, fieldType, dateSelect, records, findHistory) {
    let startTimestamp, endTimestamp
    let option0, option1
    let recIds = []
    option0 = dateSelect.DateOptions[0]
    option1 = dateSelect.DateOptions[1]

    if (option0 == ThisMonthOption) {
        startTimestamp = getThisMonth1st()
        endTimestamp = getNextMonth1st()
    } else if (option0 == LastMonthOption) {
        startTimestamp = getLastMonth1st()
        endTimestamp = getThisMonth1st()
    } else if (option0 == NextMonthOption) {
        startTimestamp = getNextMonth1st()
        endTimestamp = getAfterNextMonth1st()
    }

    records.forEach((item) => {
        // todo item.values 
        if (findHistory && matchCBAHistory(dateSelect.Operator, option1, item.values, startTimestamp, endTimestamp, CurrentMonthOption, BeforeMonthOption, AfterMonthOption)) {
            recIds.push(item.recordId)
            return
        }
        let recVal = getFieldValue(item, field, fieldType)
        if (matchCBA(dateSelect.Operator, option1, recVal, startTimestamp, endTimestamp, CurrentMonthOption, BeforeMonthOption, AfterMonthOption)) {
            recIds.push(item.recordId)
        }
    })
}

function quarterOptionHandler(field, fieldType, dateSelect, records, findHistory) {
    let startTimestamp, endTimestamp
    let option0, option1
    let recIds = []
    option0 = dateSelect.DateOptions[0]
    option1 = dateSelect.DateOptions[1]

    if (option0 == ThisQuarterOption) {
        startTimestamp = getThisQuarter1st()
        endTimestamp = getNextQuarter1st()
    } else if (option0 == LastQuarterOption) {
        startTimestamp = getLastQuarter1st()
        endTimestamp = getThisQuarter1st()
    } else if (option0 == NextQuarterOption) {
        startTimestamp = getNextQuarter1st()
        endTimestamp = getAfterNextQuarter1st()
    }
    records.forEach((item) => {
        // todo item.values 
        if (findHistory && matchCBAHistory(dateSelect.Operator, option1, item.values, startTimestamp, endTimestamp, CurrentQuarterOption, BeforeQuarterOption, AfterQuarterOption)) {
            recIds.push(item.recordId)
            return
        }
        let recVal = getFieldValue(item, field, fieldType)
        if (matchCBA(dateSelect.Operator, option1, recVal, startTimestamp, endTimestamp, CurrentQuarterOption, BeforeQuarterOption, AfterQuarterOption)) {
            recIds.push(item.recordId)
        }
    })

    return recIDs
}

// 今年/去年/明年 
function yearCBAOptionHandler(field, fieldType, dateSelect, records, findHistory) {
    let startTimestamp, endTimestamp
    let option0, option1
    let recIds = []
    option0 = dateSelect.DateOptions[0]
    option1 = dateSelect.DateOptions[1]
    
    if (option0 == ThisYearOption) {
        startTimestamp = getThisYear1st()
        endTimestamp = getNextYear1st()
    } else if (option0 == LastYearOption) {
        startTimestamp = getLastYear1st()
        endTimestamp = getThisYear1st()
    } else if (option0 == NextYearOption) {
        startTimestamp = getNextYear1st()
        endTimestamp = getAfterNextYear1st()
    }

    records.forEach((item) => {
        // todo item.values 
        if (findHistory && matchCBAHistory(dateSelect.Operator, option1, item.values, startTimestamp, endTimestamp, CurrentYearOption, BeforeYearOption, AfterYearOption)) {
            recIds.push(item.recordId)
            return
        }
        let recVal = getFieldValue(item, field, fieldType)
        if (matchCBA(dateSelect.Operator, option1, recVal, startTimestamp, endTimestamp, CurrentYearOption, BeforeYearOption, AfterYearOption)) {
            recIds.push(item.recordId)
        }
    })

    return recIDs
}

// 今年/去年/明年  具体日期 月日选择 当天/之前/之后
function yearMonthDayCBAOptionHandler(field, fieldType, dateSelect, records, findHistory) {
    let startTimestamp, endTimestamp
    let option0, option1
    let recIds = []
    let timeArr = dateSelect.DateValues[0].split("-")
    let month = timeArr[0] - 1
    let date = timeArr[1]
    option0 = dateSelect.DateOptions[0]
    option1 = dateSelect.DateOptions[2]
    if (option0 == ThisYearOption) {
        startTimestamp = getThisYear1st(month, date)
    } else if (option0 == LastYearOption) {
        startTimestamp = getLastYear1st(month, date)
    } else if (option0 == NextYearOption) {
        startTimestamp = getNextYear1st(month, date)
    }
    endTimestamp = new Date(startTimestamp + DayTimestamp)
    records.forEach((item, i) => {
        // todo item.values 
        if (findHistory && matchCBAHistory(dateSelect.Operator, option1, item.values, startTimestamp, endTimestamp, CurrentTodayOption, BeforeTodayOption, AfterTodayOption)) {
            recIds.push(item.recordId)
            return
        }
        let recVal = getFieldValue(item, field, fieldType)
        if (matchCBA(dateSelect.Operator, option1, recVal, startTimestamp, endTimestamp, CurrentTodayOption, BeforeTodayOption, AfterTodayOption)) {
            recIds.push(item.recordId)
        }
    })

    return recIDs
}

// 今年/去年/明年 具体日期区间 两个月日选择
function yearMonthDayRangeOptionHandler(field, fieldType, dateSelect, records, findHistory) {
    let startTimestamp, endTimestamp
    let option = dateSelect.DateOptions[0]
    let recIds = []
    let month0, date0, month1, date1
    let firstMonthDateArr = dateSelect.DateValues[0].split("-")
    let secondMonthDateArr = dateSelect.DateValues[1].split("-")
    month0 = firstMonthDateArr[0] - 1
    date0 = firstMonthDateArr[1]
    month1 = secondMonthDateArr[0] - 1
    date1 = secondMonthDateArr[1]
    if (option == ThisYearOption) {
        startTimestamp = getThisYear1st(month0, date0)
        endTimestamp = getThisYear1st(month1, date1)
    } else if (option == LastYearOption) {
        startTimestamp = getLastYear1st(month0, date0)
        endTimestamp = getLastYear1st(month1, date1)
    } else if (option == NextYearOption) {
        startTimestamp = getNextYear1st(month0, date0)
        endTimestamp = getNextYear1st(month1, date1)
    }

    records.forEach((item) => {
        // todo item.values 
        if (findHistory && matchDateRangeHistory(dateSelect.Operator, item.values, startTimestamp, endTimestamp)) {
            recIds.push(item.recordId)
            return
        }
        let recVal = getFieldValue(item, field, fieldType)
        if (matchDateRange(dateSelect.Operator, recVal, startTimestamp, endTimestamp)) {
            recIds.push(item.recordId)
        }
    })

    return recIds
}

function fromTodayOptionHandler(field, fieldType, dateSelect, records, findHistory) {
    let startTimestamp, endTimestamp
    let option0, option1
    let recIds = []
    let diffDay = dateSelect.DateValues[0]
    option0 = dateSelect.DateOptions[0]
    option1 = dateSelect.DateOptions[1]
    if (option0 == FromTodayAddOption) {
        startTimestamp = new Date(new Date(new Date().setHours(0, 0, 0, 0)).getTime() + (diffDay * DayTimestamp)).getTime()
    } else if (option0 == FromTodaySubOption) {
        startTimestamp = new Date(new Date(new Date().setHours(0, 0, 0, 0)).getTime() - (diffDay * DayTimestamp)).getTime()
    }

    endTimestamp = new Date(startTimestamp + DayTimestamp)
  
    records.forEach((item) => {
        // todo item.values 
        if (findHistory && matchCBAHistory(dateSelect.Operator, option1, item.values, startTimestamp, endTimestamp, CurrentTodayOption, BeforeTodayOption, AfterTodayOption)) {
            recIds.push(item.recordId)
            return
        }
        let recVal = getFieldValue(item, field, fieldType)
        if (matchCBA(dateSelect.Operator, option1, recVal, startTimestamp, endTimestamp, CurrentTodayOption, BeforeTodayOption, AfterTodayOption)) {
            recIds.push(item.recordId)
        }
    })
    

    return recIDs
}

// 具体日期 年月日选择 当天/之前/之后
function selectDateOptionHandler(field, fieldType, dateSelect, records, findHistory) {
    let startTimestamp, endTimestamp
    let option, year, month, date
    let recIds = []
    let fullDateArr = dateSelect.DateValues[0].split("-") // 2022-01-06
    year = fullDateArr[0]
    month = fullDateArr[1] - 1
    date = fullDateArr[2]
    option = dateSelect.DateOptions[1]

    startTimestamp = new Date(year, month, date)
    endTimestamp = new Date(startTimestamp + DayTimestamp)
    records.forEach((item) => {
        // todo item.values 
        if (findHistory && matchCBAHistory(dateSelect.Operator, option, item.values, startTimestamp, endTimestamp, CurrentTodayOption, BeforeTodayOption, AfterTodayOption)) {
            recIds.push(item.recordId)
            return
        }
        let recVal = getFieldValue(item, field, fieldType)
        if (matchCBA(dateSelect.Operator, option, recVal, startTimestamp, endTimestamp, CurrentTodayOption, BeforeTodayOption, AfterTodayOption)) {
            recIds.push(item.recordId)
        }
    })

    return recIDs
}

// 具体日期区间 两个年月日选择
function selectDateRangeOptionHandler(field, fieldType, dateSelect, records, findHistory) {
    let startTimestamp, endTimestamp
    let recIds = []
    let year0, month0, date0, year1, month1, date1
    let firstFullDateArr = dateSelect.DateValues[0].split("-")
    let secondFullDateArr = dateSelect.DateValues[1].split("-")
    year0 = firstFullDateArr[0]
    month0 = firstFullDateArr[1] - 1
    date0 = firstFullDateArr[2]
    year1 = secondFullDateArr[0]
    month1 = secondFullDateArr[1] - 1
    date1 = secondFullDateArr[2]

    startTimestamp = new Date(year0, month0, date0)
    endTimestamp = new Date(year1, month1, date1)
    records.forEach((item) => {
        // todo item.values 
        if (findHistory && matchDateRangeHistory(dateSelect.Operator, item.values, startTimestamp, endTimestamp)) {
            recIds.push(item.recordId)
            return
        }
        let recVal = getFieldValue(item, field, fieldType)
        if (matchDateRange(dateSelect.Operator, recVal, startTimestamp, endTimestamp)) {
            recIds.push(item.recordId)
        }
    })
}


/**
 * CBA =》 current before after
 */
function matchCBA(operator, option, timestamp, startTimestamp, endTimestamp, currentOption, beforeOption, afterOption) {
    if (operator == EqualsTimeOperator) {
        if (option == currentOption && startTimestamp <= timestamp && timestamp < endTimestamp) {
            // 在当天/周/月/季
            return true
        } else if (option == beforeOption && timestamp < startTimestamp) {
             // 在之前
             return true
        } else if (option == afterOption && timestamp >= endTimestamp) {
            return true
        }
    } else if (operator == NotEqualsTimeOperator){
        if (option == currentOption && (timestamp < startTimestamp || timestamp >= endTimestamp)) {
            // 不在当天/周/月/季
            return true
        } else if (option == beforeOption && timestamp >= startTimestamp) {
            // 不在之前
            return true
       } else if (option == afterOption && timestamp < endTimestamp) {
           // 不在之后
           return true
       }
    }

    return false
}

function matchCBAHistory(operator, option, historyValues, startTimestamp, endTimestamp, currentOption, beforeOption, afterOption) {
    for (let i = 0; i < historyValues.length; i++) {
        let timestamp 
        // totdo timestamp = historyValues[i]
        if (matchCBA(operator, option, timestamp, startTimestamp, endTimestamp, currentOption, beforeOption, afterOption)) {
            return true
        }
    }
    
    return false
}

function matchDateRange(operator, timestamp, startTimestamp, endTimestamp) {
    if (operator == EqualsTimeOperator && startTimestamp <= timestamp && timestamp < endTimestamp) {
        return true
    } else if (operator == NotEqualsTimeOperator && (startTimestamp > timestamp || timestamp >= endTimestamp)) {
        return true
    }

    return false
}

function matchDateRangeHistory(operator, historyValues, startTimestamp, endTimestamp) {
    for (let i = 0; i < historyValues.length; i++) {
        let timestamp 
        // totdo timestamp = historyValues[i]
        if (matchDateRange(operator, timestamp, startTimestamp, endTimestamp)) {
            return true
        }
    }

    return false
}

function getThisWeekMonday() {
    let today = new Date(new Date().setHours(0, 0, 0, 0))
    let day = today.getDay() || 7
    let monday = new Date(today.setDate(today.getDate()- (day - 1))).getTime()
    return monday
}

function getThisMonth1st() {
    let today = new Date(new Date().setHours(0, 0, 0, 0))
    let first = new Date(today.setDate(1)).getTime()
    return first
}

function getLastMonth1st() {
    let today = new Date(new Date().setHours(0, 0, 0, 0))
    let month = today.getMonth()
    let year = today.getFullYear()
    if (month == 0) {
        month = 11
        year -= 1
    } else {
        month -= 1
    }
    let first = new Date(year, month, 1).getTime()
    return first
}

function getNextMonth1st() {
    let today = new Date(new Date().setHours(0, 0, 0, 0))
    let month = today.getMonth()
    let year = today.getFullYear()
    if (month == 11) {
        month = 0
        year += 1
    } else {
        month += 1
    }
    let first = new Date(year, month, 1).getTime()
    return first
}

function getAfterNextMonth1st() {
    let nextMonth = new Date(getNextMonth1st())
    let month = nextMonth.getMonth()
    let year = nextMonth.getFullYear()
    if (month == 11) {
        month = 0
        year += 1
    } else {
        month += 1
    }
    let first = new Date(year, month, 1).getTime()
    return first

}

function getThisQuarter1st() {
    let today = new Date(new Date().setHours(0, 0, 0, 0))
    let month = today.getMonth()
    let year = today.getFullYear()
    let first = new Date(year, getQuarterStartMonth(month), 1)
    return first
}

function getLastQuarter1st() {
    let thisQuarter = new Date(getThisQuarter1st())
    let month = thisQuarter.getMonth()
    let year = thisQuarter.getFullYear()
    if (month == 0) {
        month = 11
        year -= 1
    } else {
        month -= 3
    } 
    let first = new Date(year, getQuarterStartMonth(month), 1)
    return first
}

function getNextQuarter1st() {
    let thisQuarter = new Date(getThisQuarter1st())
    let month = thisQuarter.getMonth()
    let year = thisQuarter.getFullYear()
    if (month == 11) {
        month = 0
        year += 1
    } else {
        month += 3
    }
    let first = new Date(year, getQuarterStartMonth(month), 1)
    return first
}

function getAfterNextQuarter1st() {
    let nextQuarter = new Date(getNextQuarter1st())
    let month = nextQuarter.getMonth()
    let year = nextQuarter.getFullYear()
    if (month == 11) {
        month = 0
        year += 1
    } else {
        month += 3
    }
    let first = new Date(year, getQuarterStartMonth(month), 1)
    return first
}

function getQuarterStartMonth(nowMonth) {
    let quarterStartMonth = 0;
    if (nowMonth < 3) {
        quarterStartMonth = 0;
    }
    if (2 < nowMonth && nowMonth < 6) {
        quarterStartMonth = 3;
    }
    if (5 < nowMonth && nowMonth < 9) {
        quarterStartMonth = 6;
    }
    if (nowMonth > 8) {
        quarterStartMonth = 9;
    }
    return quarterStartMonth;
}

function getThisYear1st(month, date) {
    let today = new Date()
    let year = today.getFullYear()
    if (month == undefined || date == undefined) {
        month = 0
        date = 1
    }
    let first = new Date(year, month, date)
    return first
}

function getLastYear1st(month, date) {
    let today = new Date()
    let year = today.getFullYear() - 1
    if (month == undefined || date == undefined) {
        month = 0
        date = 1
    }
    let first = new Date(year, month, date)
    return first
}

function getNextYear1st(month, date) {
    let today = new Date()
    let year = today.getFullYear() + 1
    if (month == undefined || date == undefined) {
        month = 0
        date = 1
    }
    let first = new Date(year, month, date)
    return first
}

function getAfterNextYear1st(month, date) {
    let today = new Date()
    let year = today.getFullYear() + 2
    if (month == undefined || date == undefined) {
        month = 0
        date = 1
    }
    let first = new Date(year, month, date)
    return first
}

var recordsRaw = [
    {
        "recordId": "046c9ad8-a461-42d0-9e47-57e474e68340",
        "cellValuesByColumnId": {
            "actual_investment_amount": 100,
            "entity": {
                "value": {
                    "displayName": "小米2",
                    "entityId": "046c9ad8-a461-42d0-9e47-57e474e68340",
                    "entityType": "EVENT"
                },
                "recordId": "046c9ad8-a461-42d0-9e47-57e474e68340"
            },
            "entity_primary_name": "小米2",
            "entity_type": "EVENT",
            "exit_amount": 0,
            "founded_on": {
                "raw": 1418227200,
                "shadow": 1418227200,
                "indicator": {
                    "value": 1418227200,
                    "recordId": "1030666186-ORGANIZATION"
                }
            },
            "fund_of_fund_rec_id": [
                {
                    "value": {
                        "displayName": "福田引导",
                        "entityId": "6f4f4917-e323-4c94-b6e9-bcaddcecdf7d",
                        "entityType": "ORGANIZATION"
                    },

                    
                    "recordId": "6f4f4917-e323-4c94-b6e9-bcaddcecdf7d"
                }
            ],
            "indicator_founded_on": {
                "value": 1418227200,
                "recordId": "1030666186-ORGANIZATION"
            },
            "indicator_region_name": {
                "value": "大连市甘井子区",
                "recordId": "1030666186-ORGANIZATION"
            },
            "indicator_registered_address": {
                "value": "辽宁省大连市甘井子区凌水镇栾金村",
                "recordId": "1030666186-ORGANIZATION"
            },
            "indicator_social_credit_id": {
                "value": "91210231311465219U",
                "recordId": "1030666186-ORGANIZATION"
            },
            "indicator_vertical": {
                "value": "电子商务",
                "recordId": "1030666186-ORGANIZATION"
            },
            "investment_target_rec_id": [
                {
                    "value": {
                        "displayName": "大米科技",
                        "entityId": "c2eda5b9-e4bb-47c8-85f8-dcf29ecaeb3a",
                        "entityType": "ORGANIZATION"
                    },
                    "recordId": "c2eda5b9-e4bb-47c8-85f8-dcf29ecaeb3a"
                }
            ],
            "investment_target_rime_entity": [
                {
                    "value": {
                        "displayName": "大连大米科技有限公司",
                        "entityId": "1030666186",
                        "entityType": "ORGANIZATION"
                    },
                    "recordId": "c2eda5b9-e4bb-47c8-85f8-dcf29ecaeb3a"
                }
            ],
            "investment_target_rime_entity_id": [
                "1030666186"
            ],
            "investment_target_rime_entity_type": [
                "ORGANIZATION"
            ],
            "raw_founded_on": 1418227200,
            "raw_region_name": "大连市甘井子区",
            "raw_registered_address": "辽宁省大连市甘井子区凌水镇栾金村",
            "raw_social_credit_id": "91210231311465219U",
            "raw_vertical": "电子商务",
            "rec_id": "046c9ad8-a461-42d0-9e47-57e474e68340",
            "region_name": {
                "raw": "大连市甘井子区",
                "shadow": "大连市甘井子区",
                "indicator": {
                    "value": "大连市甘井子区",
                    "recordId": "1030666186-ORGANIZATION"
                }
            },
            "registered_address": {
                "raw": "辽宁省大连市甘井子区凌水镇栾金村",
                "shadow": "辽宁省大连市甘井子区凌水镇栾金村",
                "indicator": {
                    "value": "辽宁省大连市甘井子区凌水镇栾金村",
                    "recordId": "1030666186-ORGANIZATION"
                }
            },
            "shadow_founded_on": 1418227200,
            "shadow_region_name": "大连市甘井子区",
            "shadow_registered_address": "辽宁省大连市甘井子区凌水镇栾金村",
            "shadow_social_credit_id": "91210231311465219U",
            "shadow_vertical": "电子商务",
            "social_credit_id": {
                "raw": "91210231311465219U",
                "shadow": "91210231311465219U",
                "indicator": {
                    "value": "91210231311465219U",
                    "recordId": "1030666186-ORGANIZATION"
                }
            },
            "vertical": {
                "raw": "电子商务",
                "shadow": "电子商务",
                "indicator": {
                    "value": "电子商务",
                    "recordId": "1030666186-ORGANIZATION"
                }
            }
        }
    },
    {
        "recordId": "29dcf113-951f-4340-b706-0c2911a15859",
        "cellValuesByColumnId": {
            "development_stage": "210104",
            "entity": {
                "value": {
                    "displayName": "丝路B轮",
                    "entityId": "29dcf113-951f-4340-b706-0c2911a15859",
                    "entityType": "EVENT"
                },
                "recordId": "29dcf113-951f-4340-b706-0c2911a15859"
            },
            "entity_primary_name": "丝路B轮",
            "entity_type": "EVENT",
            "exit_amount": 0,
            "founded_on": {
                "raw": 953740800,
                "shadow": 953740800,
                "indicator": {
                    "value": 953740800,
                    "recordId": "1045824165-ORGANIZATION"
                }
            },
            "fund_of_fund_rec_id": [
                {
                    "value": {
                        "displayName": "福田引导",
                        "entityId": "6f4f4917-e323-4c94-b6e9-bcaddcecdf7d",
                        "entityType": "ORGANIZATION"
                    },
                    "recordId": "6f4f4917-e323-4c94-b6e9-bcaddcecdf7d"
                }
            ],
            "in_early_stage": "Y",
            "indicator_founded_on": {
                "value": 953740800,
                "recordId": "1045824165-ORGANIZATION"
            },
            "indicator_region_name": {
                "value": "深圳市福田区",
                "recordId": "1045824165-ORGANIZATION"
            },
            "indicator_registered_address": {
                "value": "深圳市福田区福强路3030号文化体育产业总部大厦17楼",
                "recordId": "1045824165-ORGANIZATION"
            },
            "indicator_social_credit_id": {
                "value": "914403007152851426",
                "recordId": "1045824165-ORGANIZATION"
            },
            "indicator_vertical": {
                "value": "智能营销",
                "recordId": "1045824165-ORGANIZATION"
            },
            "investment_deal_type": "10105",
            "investment_target_rec_id": [
                {
                    "value": {
                        "displayName": "丝路视觉",
                        "entityId": "108a46ec-7b8e-41d7-ad41-32fe3734eccf",
                        "entityType": "ORGANIZATION"
                    },
                    "recordId": "108a46ec-7b8e-41d7-ad41-32fe3734eccf"
                }
            ],
            "investment_target_rime_entity": [
                {
                    "value": {
                        "displayName": "丝路视觉科技股份有限公司",
                        "entityId": "1045824165",
                        "entityType": "ORGANIZATION"
                    },
                    "recordId": "108a46ec-7b8e-41d7-ad41-32fe3734eccf"
                }
            ],
            "investment_target_rime_entity_id": [
                "1045824165"
            ],
            "investment_target_rime_entity_type": [
                "ORGANIZATION"
            ],
            "is_unicorn": "Y",
            "major_businesses": "卖报",
            "raw_founded_on": 953740800,
            "raw_region_name": "深圳市福田区",
            "raw_registered_address": "深圳市福田区福强路3030号文化体育产业总部大厦17楼",
            "raw_social_credit_id": "914403007152851426",
            "raw_vertical": "智能营销，金融，人工智能",
            "rec_id": "29dcf113-951f-4340-b706-0c2911a15859",
            "region_name": {
                "raw": "深圳市福田区",
                "shadow": "深圳市福田区",
                "indicator": {
                    "value": "深圳市福田区",
                    "recordId": "1045824165-ORGANIZATION"
                }
            },
            "registered_address": {
                "raw": "深圳市福田区福强路3030号文化体育产业总部大厦17楼",
                "shadow": "深圳市福田区福强路3030号文化体育产业总部大厦17楼",
                "indicator": {
                    "value": "深圳市福田区福强路3030号文化体育产业总部大厦17楼",
                    "recordId": "1045824165-ORGANIZATION"
                }
            },
            "shadow_founded_on": 953740800,
            "shadow_region_name": "深圳市福田区",
            "shadow_registered_address": "深圳市福田区福强路3030号文化体育产业总部大厦17楼",
            "shadow_social_credit_id": "914403007152851426",
            "shadow_vertical": "智能营销",
            "social_credit_id": {
                "raw": "914403007152851426",
                "shadow": "914403007152851426",
                "indicator": {
                    "value": "914403007152851426",
                    "recordId": "1045824165-ORGANIZATION"
                }
            },
            "sub_vertical": "硬科技",
            "vertical": {
                "raw": "智能营销，金融，人工智能",
                "shadow": "智能营销",
                "indicator": {
                    "value": "智能营销",
                    "recordId": "1045824165-ORGANIZATION"
                }
            }
        }
    },
    {
        "recordId": "3a87c276-69ab-4147-bf6c-0324dd31737c",
        "cellValuesByColumnId": {
            "actual_investment_amount": 122345,
            "current_valuation": 666666,
            "entity": {
                "value": {
                    "displayName": "小米科技",
                    "entityId": "3a87c276-69ab-4147-bf6c-0324dd31737c",
                    "entityType": "EVENT"
                },
                "recordId": "3a87c276-69ab-4147-bf6c-0324dd31737c"
            },
            "entity_primary_name": "小米科技",
            "entity_type": "EVENT",
            "exit_amount": 0,
            "former_valuation": 3345,
            "founded_on": {
                "raw": 1418227200,
                "shadow": 1418227200,
                "indicator": {
                    "value": 1418227200,
                    "recordId": "1030666186-ORGANIZATION"
                }
            },
            "fund_of_fund_rec_id": [
                {
                    "value": {
                        "displayName": "福田引导",
                        "entityId": "6f4f4917-e323-4c94-b6e9-bcaddcecdf7d",
                        "entityType": "ORGANIZATION"
                    },
                    "recordId": "6f4f4917-e323-4c94-b6e9-bcaddcecdf7d"
                }
            ],
            "indicator_founded_on": {
                "value": 1418227200,
                "recordId": "1030666186-ORGANIZATION"
            },
            "indicator_region_name": {
                "value": "大连市甘井子区",
                "recordId": "1030666186-ORGANIZATION"
            },
            "indicator_registered_address": {
                "value": "辽宁省大连市甘井子区凌水镇栾金村",
                "recordId": "1030666186-ORGANIZATION"
            },
            "indicator_social_credit_id": {
                "value": "91210231311465219U",
                "recordId": "1030666186-ORGANIZATION"
            },
            "indicator_vertical": {
                "value": "电子商务",
                "recordId": "1030666186-ORGANIZATION"
            },
            "intended_investment_amount": 2222,
            "intended_shareholding_percentage": 0,
            "invested_on": 1356105600,
            "investment_target_rec_id": [
                {
                    "value": {
                        "displayName": "大米科技",
                        "entityId": "c2eda5b9-e4bb-47c8-85f8-dcf29ecaeb3a",
                        "entityType": "ORGANIZATION"
                    },
                    "recordId": "c2eda5b9-e4bb-47c8-85f8-dcf29ecaeb3a"
                }
            ],
            "investment_target_rime_entity": [
                {
                    "value": {
                        "displayName": "大连大米科技有限公司",
                        "entityId": "1030666186",
                        "entityType": "ORGANIZATION"
                    },
                    "recordId": "c2eda5b9-e4bb-47c8-85f8-dcf29ecaeb3a"
                }
            ],
            "investment_target_rime_entity_id": [
                "1030666186"
            ],
            "investment_target_rime_entity_type": [
                "ORGANIZATION"
            ],
            "latest_shareholding_value": 356784,
            "raw_founded_on": 1418227200,
            "raw_region_name": "大连市甘井子区",
            "raw_registered_address": "辽宁省大连市甘井子区凌水镇栾金村",
            "raw_social_credit_id": "91210231311465219U",
            "raw_vertical": "电子商务",
            "rec_id": "3a87c276-69ab-4147-bf6c-0324dd31737c",
            "region_name": {
                "raw": "大连市甘井子区",
                "shadow": "大连市甘井子区",
                "indicator": {
                    "value": "大连市甘井子区",
                    "recordId": "1030666186-ORGANIZATION"
                }
            },
            "registered_address": {
                "raw": "辽宁省大连市甘井子区凌水镇栾金村",
                "shadow": "辽宁省大连市甘井子区凌水镇栾金村",
                "indicator": {
                    "value": "辽宁省大连市甘井子区凌水镇栾金村",
                    "recordId": "1030666186-ORGANIZATION"
                }
            },
            "shadow_founded_on": 1418227200,
            "shadow_region_name": "大连市甘井子区",
            "shadow_registered_address": "辽宁省大连市甘井子区凌水镇栾金村",
            "shadow_social_credit_id": "91210231311465219U",
            "shadow_vertical": "电子商务",
            "social_credit_id": {
                "raw": "91210231311465219U",
                "shadow": "91210231311465219U",
                "indicator": {
                    "value": "91210231311465219U",
                    "recordId": "1030666186-ORGANIZATION"
                }
            },
            "vertical": {
                "raw": "电子商务",
                "shadow": "电子商务",
                "indicator": {
                    "value": "电子商务",
                    "recordId": "1030666186-ORGANIZATION"
                }
            }
        }
    }
]


function getFieldValue(item, field, fieldType) {
    let val
    if (fieldType == CollisionFieldType) {
        val = getCollisionValue(item, field)
    } else if (fieldType == EntityFieldType) {
        val = getEntityValue(item, field)
    } else if (fieldType == LinkFieldType) {
        val = getLinkCurrentValue(item, field)
    } else {
        val = item.cellValuesByColumnId[field]
    }
    return val
}

function getCollisionValue(item, field) {
    return item.cellValuesByColumnId[field].raw
}

function getEntityValue(item, field) {
    return item.cellValuesByColumnId[field].value.displayName
}

function getLinkCurrentValue(item, field) {
    let arr = item.cellValuesByColumnId[field]
    if (arr.length != 0) {
        return arr[0].value.displayName
    }

    return null
}

initFuncMap()
filterDemo()