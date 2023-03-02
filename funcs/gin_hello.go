package funcs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.mvalley.com/common/cain/pkg/prime"
	md "gitlab.mvalley.com/rime-index/common/pkg/factor/models/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var eMDB *mongo.Database

func HelloWorld() {
	// mdb, err := ConnectToMongoDB("mongodb://root:pacman@192.168.88.201:33017/factor?authSource=admin", "factor")
	mdb, err := ConnectToMongoDB("mongodb://root:pacman@10.20.70.33:27017/factor?authSource=admin", "factor")
	if err != nil {
		panic(err)
	}
	eMDB = mdb
	r := gin.Default()
	//跨域
	r.Use(Cors)
	r.Static("/office", "./htmls")
	r.Static("/wps", "./wps-addon-build")
	r.GET("/match/:id", func(c *gin.Context) {
		ids := c.Param("id")
		fmt.Println(ids)
		var entityIDs []string
		err := json.Unmarshal([]byte(ids), &entityIDs)
		if err != nil {
			panic(err)
		}

		fmt.Println(entityIDs)
		d := getCompanyInfo(entityIDs)
		c.JSON(200, d)
	})
	// r.GET("/home.html", func(c *gin.Context) {
	// 	c.JSON(200, "is home")
	// })

	err = r.Run(":8888")
	if err != nil {
		panic(err)
	}
}

func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")

	// 允许的请求头， 如果有自定义的Header都要添加到这里，不然浏览器会报跨域错误
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Token, Timestamp. AppId, Referer")

	// 允许的请求方式都需要添加
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func getCompanyInfo(entityIDs []string) []map[string]interface{} {

	collection := eMDB.Collection("factor_test_0809")

	filters := bson.M{
		"entity_id": bson.M{
			"$in": entityIDs,
		},
	}
	var docs []md.FactorMongoDocument
	// 选择返回的列
	projection := make(bson.M, 0)
	projection["entity_id"] = 1
	projection["saic_legal_name"] = 1
	projection["founded_on"] = 1
	projection["registered_address"] = 1
	//t1 := time.Now()
	cur, err := collection.Find(context.TODO(), filters,
		&options.FindOptions{
			Projection: projection,
		})
	//fmt.Printf("[GetMongoResponseForPowerfulSheet] Find, collection: %s, time spent: %f \n", index, time.Now().Sub(t1).Seconds())
	defer cur.Close(context.TODO())
	if err != nil {
		panic(err)
	}

	//t2 := time.Now()

	err = cur.All(context.Background(), &docs)

	//fmt.Printf("[GetMongoResponseForPowerfulSheet] All, collection: %s, time spent: %f \n", index, time.Now().Sub(t2).Seconds())
	if err != nil {
		panic(err)
	}

	var idxMap = make(map[string]int)
	for i, id := range entityIDs {
		idxMap[id] = i
	}
	var mArr = make([]map[string]interface{}, len(entityIDs))
	for _, doc := range docs {
		m := make(map[string]interface{})
		m["saic_legal_name"] = ExcelTextTranscoder(doc.SaicLegalName)
		m["founded_on"] = ExcelTextTranscoder(doc.FoundedOn)
		m["registered_address"] = ExcelTextTranscoder(doc.RegisteredAddress)
		// mArr = append(mArr, m)
		mArr[idxMap[doc.EntityID]] = m
	}

	return mArr
}

func ExcelTextTranscoder(dataArray []prime.Data) string {
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
	if len(values) <= 0 {
		return ""
	}
	return strings.Join(values, "")
}
