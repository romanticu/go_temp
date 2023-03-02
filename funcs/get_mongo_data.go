package funcs

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"gitlab.mvalley.com/common/cain/pkg/bedrock"
	mt "gitlab.mvalley.com/common/cain/pkg/mongo"
	"gitlab.mvalley.com/common/cain/pkg/prime"
	"gitlab.mvalley.com/rime-index/common/pkg/exhibitdata"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var fMongoDB *mongo.Database

func initMG() {
	fMongoDB = mt.ConnectToMongoDB("mongodb://root:pacman@192.168.88.201:33017/factor?authSource=admin", "factor")
}
func GetMongoDB(orgIDs []string, columnIds []string) {
	initMG()
	handlers := R.GetHandlers(columnIds)
	for _, h := range handlers {
		err := GetMongoResponseForSheet(fMongoDB, h.GetMongoCollection(), orgIDs, h.GetDisplayColumnIDs(), h.GetMongoDBModels())
		if err != nil {
			panic(err)
		}
		h.InitMongoDocsMap()
	}
	var rowsData []*exhibitdata.RowData
	for _, id := range orgIDs {
		rowData := new(exhibitdata.RowData)
		for _, h := range handlers {
			docsMap := h.GetMongoDocsMap()
			doc, ok := docsMap[id]
			if !ok {
				continue
			}
			// jd, _ := json.Marshal(doc)
			// fmt.Println("----- doc")
			// fmt.Println(string(jd))
			err := getRowDataFromMongoDocument(doc, rowData, h.GetDisplayColumnIDs(), h.GetTranscoderMap())
			if err != nil {
				panic(err)
			}
		}
		//s, _ := rowData.Marshal()
		//fmt.Println("rowData")
		//fmt.Println(s)
		rowsData = append(rowsData, rowData)
	}

}

type Router struct {
}

func (r *Router) GetHandlers(columnIDs []string) []SheetHandler {
	var initHandlers = initHandlers()
	var handlers []SheetHandler

	for i := range initHandlers {
		if initHandlers[i].ContainColumn(columnIDs) {
			handlers = append(handlers, initHandlers[i])
		}
	}

	return handlers
}

var R Router

func initHandlers() []SheetHandler {
	var handlers []SheetHandler
	handlers = append(handlers, &factorHandler{})
	return handlers
}

type SheetHandler interface {
	GetTranscoderMap() map[string]Transcoder
	ContainColumn(columnIDs []string) bool
	GetDisplayColumnIDs() []string
	GetMongoCollection() string
	InitMongoDocsMap()
	GetMongoDocsMap() map[string]interface{}
	GetMongoDBModels() []interface{}
}

type factorHandler struct {
	displayColumnIDs  []string
	mongoDocuments    []*FactorMongoDocument
	mongoDocumentsArr []*[]*FactorMongoDocument
	mongoDocumentsMap map[string]interface{}
}

// GetTranscoderMap ...
func (f *factorHandler) GetTranscoderMap() map[string]Transcoder {
	return map[string]Transcoder{
		"saic_legal_name": TextTranscoder{},
		"founded_on":      DateTranscoder{},
		"former_names":    TextTranscoder{},
	}
}

func (f *factorHandler) ContainColumn(columnIDs []string) bool {
	tMap := f.GetTranscoderMap()
	var contain bool
	var displayColumnIDs []string
	for _, item := range columnIDs {
		_, ok := tMap[item]
		if ok {
			contain = true
		}
		displayColumnIDs = append(displayColumnIDs, item)
	}
	f.displayColumnIDs = displayColumnIDs
	return contain
}

func (f *factorHandler) GetDisplayColumnIDs() []string {
	return f.displayColumnIDs
}

func (f *factorHandler) GetMongoCollection() string {
	return "factor_0322"
}
func (f *factorHandler) InitMongoDocsMap() {
	fmt.Println(f.mongoDocumentsArr)
	docsMap := make(map[string]interface{})
	for i := range f.mongoDocumentsArr {
		docs := *f.mongoDocumentsArr[i]
		for j := range docs {
			docsMap[docs[j].EntityID] = docs[j]
		}

	}
	f.mongoDocumentsMap = docsMap
}
func (f *factorHandler) GetMongoDocsMap() map[string]interface{} {
	return f.mongoDocumentsMap
}

func (f *factorHandler) GetMongoDBModels() []interface{} {
	pointerArr := make([]interface{}, 0)
	f.mongoDocumentsArr = nil
	for i := 0; i < goroutineNum; i++ {
		d := make([]*FactorMongoDocument, 0)
		pointerArr = append(pointerArr, &d)
		f.mongoDocumentsArr = append(f.mongoDocumentsArr, &d)
	}

	return pointerArr
}

var goroutineNum = 5

func GetMongoResponseForSheet(mdb *mongo.Database, index string, entityIDs []string, columnIDs []string, docs []interface{}) error {
	collection := mdb.Collection(index)
	if len(entityIDs) == 0 {
		return errors.New("query entity id must not null")
	}

	// 选择返回的列
	projection := make(bson.M, 0)
	projection["entity_id"] = 1
	projection["entity_type"] = 1

	for _, col := range columnIDs {
		projection[col] = 1
	}
	batchSize := len(entityIDs)/goroutineNum + 1
	j := 0
	var wg sync.WaitGroup

	for i := 0; i < len(entityIDs); i += batchSize {
		e := i + batchSize
		if e > len(entityIDs) {
			e = len(entityIDs)
		}
		filters := bson.M{
			"entity_id": bson.M{
				"$in": entityIDs[i:e],
			},
		}
		wg.Add(1)
		index := j
		j++
		go func(index int) {
			cur, err := collection.Find(context.TODO(), filters, &options.FindOptions{
				Projection: projection,
			})
			defer cur.Close(context.TODO())
			if err != nil {
				fmt.Println(err)
			}
			err = cur.All(context.Background(), docs[index])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(index, " ", docs[index])
			wg.Done()
		}(index)

	}

	wg.Wait()
	return nil
}

func getRowDataFromMongoDocument(doc interface{}, rowData *exhibitdata.RowData,
	columnIDs []string, transcoderMap map[string]Transcoder) error {
	set := make(map[string]struct{})
	for _, v := range columnIDs {
		set[v] = struct{}{}
	}
	if doc == nil {
		return errors.New("variable [doc] nil error")
	}

	typ := reflect.TypeOf(doc)
	val := reflect.ValueOf(doc)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	// 跳过前两个字段
	for i := 2; i < typ.NumField(); i++ {
		jsonTag := typ.Field(i).Tag.Get("json")
		if jsonTag == "" {
			return errors.New(fmt.Sprintf("json tag is empty field index:[%d]", i))
		}
		columnID := strings.Split(jsonTag, ",")[0]
		if _, ex := set[columnID]; !ex {
			continue
		}
		//fmt.Println("1")
		colValue := val.Field(i)
		dataArray, ok := colValue.Interface().([]prime.Data)
		if !ok {
			return errors.New(fmt.Sprintf("variable [field] type is not prime.Data, cur type: %T", colValue))
		}

		t, ok := transcoderMap[columnID]
		if !ok {
			return errors.New(fmt.Sprintf("parameter [columnID] invalid error, can not find transcoder, cur columnID: %s", columnID))
		}

		cellData, err := t.Transcode(dataArray)
		if err != nil {
			return err
		}
		rowData.WithCell(columnID, cellData)
	}

	return nil
}

type FactorMongoDocument struct {
	EntityID   string             `json:"entity_id" bson:"entity_id"`
	EntityType bedrock.EntityType `json:"entity_type" bson:"entity_type"`
	// 	企业名称 - 全称
	SaicLegalName []prime.Data `json:"saic_legal_name" bson:"saic_legal_name,omitempty"`
	// 工商曾用名
	FormerNames []prime.Data `json:"former_names" bson:"former_names,omitempty"`
	// 成立日期
	FoundedOn []prime.Data `json:"founded_on" bson:"founded_on,omitempty"`
}

// Transcoder ...
type Transcoder interface {
	Transcode(dataArray []prime.Data) (exhibitdata.CellData, error)
}

type TextTranscoder struct{}

// Transcode ...
func (t TextTranscoder) Transcode(dataArray []prime.Data) (exhibitdata.CellData, error) {
	var values []string
	for i := range dataArray {
		if dataArray[i].TextValue == nil {
			continue
		}
		value := strings.TrimSpace(dataArray[i].TextValue.Value)
		if len(value) == 0 {
			continue
		}
		values = append(values, value)

	}
	return exhibitdata.CellData{
		DisplayValues: values,
	}, nil
}

// DateTranscoder ... 日期，可以搜索
type DateTranscoder struct{}

// Transcode ...
func (t DateTranscoder) Transcode(dataArray []prime.Data) (exhibitdata.CellData, error) {
	var ids []int64
	var items []string

	for i := range dataArray {
		if dataArray[i].TextValue == nil {
			continue
		}
		value := strings.TrimSpace(dataArray[i].TextValue.Value)
		if len(value) == 0 {
			continue
		}
		t, err := time.Parse("2006-01-02", value)
		if err != nil {
			return exhibitdata.CellData{}, err
		}
		ids = append(ids, t.Unix())
		items = append(items, value)
	}
	return exhibitdata.CellData{
		SearchIDs:     ids,
		DisplayValues: items,
	}, nil
}
