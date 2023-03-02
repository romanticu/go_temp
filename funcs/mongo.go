package funcs

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func QueryRecurments() {
	// db, err := ConnectToMongoDB("mongodb://root:pacman@192.168.88.121:33017/saic?authSource=admin", "recruitment_datapack")
	db, err := ConnectToMongoDB("mongodb://root:pacman@10.20.70.33:27017/saic?authSource=admin", "recruitment_datapack")
	if err != nil {
		panic(err)
	}
	collection := db.Collection("recruitment")
	var batch int64 = 3000
	opts := options.Find().SetSort(bson.M{"_id": 1}).SetLimit(batch)

	var nextId string
	var count int64
	for {

		var filters bson.M
		if nextId != "" {
			oid, _ := primitive.ObjectIDFromHex(nextId)
			filters = bson.M{
				"_id": bson.M{
					"$gt": oid,
				},
			}
		}
		cur, err := collection.Find(context.TODO(), filters, opts)
		if err != nil {
			panic(err)
		}
		var recs []RecruitmentDoc
		err = cur.All(context.Background(), &recs)
		if err != nil {
			panic(err)
		}
		cur.Close(context.TODO())

		count += int64(len(recs))
		fmt.Println("count ", count, " --- batch ", len(recs))

		if len(recs) < int(batch) {
			break
		}
		nextId = recs[len(recs)-1].Id

	}
}

// ConnectToMongoDB ...
func ConnectToMongoDB(host string, dbName string) (*mongo.Database, error) {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(host)
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	return db, nil
}

func MongoDBRollbackTest() {
	cli, err := ConnectToMongoDB("mongodb://root:pacman@192.168.88.201:33017/test?authSource=admin", "test")
	if err != nil {
		panic(err)
	}

	col := cli.Collection("test")
	var docs []interface{}
	doc1 := map[string]string{
		"doc_id": "d0",
		"num": "2",
	}
	doc2 := map[string]string{
		"doc_id": "d2",
		"num": "2",
	}
	docs = append(docs, doc1, doc2)
	err = cli.Client().UseSession(context.TODO(), func (sc mongo.SessionContext) error {
		//开启事务
		if err := sc.StartTransaction(); err != nil {
			return err
		}
		_, err := col.InsertMany(context.TODO(), docs)
		if err != nil {
			if err := sc.AbortTransaction(sc); err != nil {
				return err
			}
			return err
		}
		//最后的时候，关闭会话
		defer sc.EndSession(context.TODO())
		return sc.CommitTransaction(context.TODO())
	})
	if err != nil {
		fmt.Println(err)
	}
	
}

type RecruitmentDoc struct {
	Id    string `bson:"_id,omitempty"`
	RecID string `bson:"rec_id,omitempty"`

	CompanyID *string `bson:"company_id,omitempty"`
	// 公司名称
	CompanyName string `bson:"company_name,omitempty"`
	// 招聘渠道名称
	SourceName string `bson:"source_name,omitempty"`
	// 首次招聘发布时间
	FirstPublishOn string `bson:"first_publish_on,omitempty"`
	// 招聘发布时间
	PublishOn string `bson:"publish_on,omitempty"`
	// 职位名称
	JobName string `bson:"job_name,omitempty"`
	// 职位关键词
	JobKeyword *string `bson:"job_keyword,omitempty"`
	// 职位 ID
	JobID int `bson:"job_id"`
	// 岗位职责
	JobDescription string `bson:"jd,omitempty"`
	// 职能类型
	JobType *int `bson:"job_type,omitempty"`
	// 职能类别
	JobFunction *string `bson:"job_function,omitempty"`
	// 职位
	JobTitleType int `bson:"job_title_type,omitempty"`
	// 所属部门
	Department *string `bson:"department,omitempty"`
	// 招聘人数
	HeaderCount int `bson:"header_count,omitempty"`
	// 城市
	CityName *string `bson:"city_name,omitempty"`
	// 区
	AreaName *string `bson:"area_name,omitempty"`
	// 工作地点
	LocationName *string `bson:"location_name,omitempty"`
	// 地区代码
	LocationCode *string `bson:"location_code,omitempty"`
	// 学历要求
	EducationRequirement *int `bson:"education_requirement,omitempty"`
	// 工作经验要求
	WorkingRequirement float64 `bson:"working_requirement"`
	// 语言要求
	LanguageRequirement *string `bson:"language_requirement,omitempty"`
	// 最低年龄
	MinAge int `bson:"min_age"`
	// 最高年龄
	MaxAge int `bson:"max_age"`
	// 最低薪水
	MinCompensation *float64 `bson:"min_compensation,omitempty"`
	// 最高薪水
	MaxCompensation *float64 `bson:"max_compensation,omitempty"`
	// 平均薪水
	AvgCompensation *float64 `bson:"avg_compensation,omitempty"`
	// 福利
	Allowance *string `bson:"allowance,omitempty"`
}
