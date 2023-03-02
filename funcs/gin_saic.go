package funcs

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"gitlab.mvalley.com/adam/common/pkg/es_utils"
)

type SaicKeywords struct {
	ID             string   `json:"kw.id"`
	Name           string   `json:"kw.name"`
	Keywords       []string `json:"ik.keywords"`
	IntegrityScore float64  `json:"fl.integrity_score"`
	Score          float64  `json:"_score"`
}

type SaicKeywordsV3 struct {
	ID             string   `json:"kw.id"`
	Name           string   `json:"ik.name"`
	Keywords       []string `json:"ik.keywords"`
	IntegrityScore float64  `json:"fl.integrity_score"`
	Score          float64  `json:"_score"`
}

var escli *elasticsearch.Client

func GinSaicKeywords() {
	escli = init23Es()
	r := gin.Default()
	r.GET("/query/:kw", func(c *gin.Context) {
		queryKWV4(c, "saic_keywords")
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})
	r.GET("/query_/:kw", func(c *gin.Context) {
		queryKWV4(c, "saic_keywords_v5")
	})


	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}

}

func init23Es() *elasticsearch.Client {
	var cfg = elasticsearch.Config{
		Addresses: []string{
			"http://192.168.88.201:23921",
		},
		Username: "",
		Password: "",
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Duration(10) * time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS11,
				InsecureSkipVerify: true,
			},
		},
	}

	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	return esClient
}

type ESCommonResponse struct {
	Hits struct {
		Total     TotalDocsInfo `json:"total"`
		Documents []DocumentCnt `json:"hits"`
	} `json:"hits"`
	Aggregations map[string]interface{} `json:"aggregations"`
}

//TotalDocsInfo ...
type TotalDocsInfo struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

//DocumentContent ...
type DocumentCnt struct {
	ID     string          `json:"_id"`
	Source json.RawMessage `json:"_source"`
	Score  float64         `json:"_score"`
}

func ESSearchCommon(esCli *elasticsearch.Client, esQuery interface{}, esIndex string) (*ESCommonResponse, error) {
	var res ESCommonResponse

	err := es_utils.PerformESQuery(esQuery, &res, esIndex, esCli)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func queryKW(c *gin.Context, index string) {
	kw := c.Param("kw")
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"function_score": map[string]interface{}{
				"query": map[string]interface{}{
					"match_phrase": map[string]interface{}{
						"ik.keywords": map[string]interface{}{
							"query": kw,
						},
					},
				},
				"script_score": map[string]interface{}{
					"script": "_score*0.8+doc['fl.integrity_score'].value*0.35",
				},
				"boost_mode": "replace",
			},
		},
		"size":    50,
		"_source": []string{"kw.name", "fl.integrity_score", "kw.id"},
	}
	esResp, err := ESSearchCommon(escli, esQuery, index)
	if err != nil {
		fmt.Println(err)
	}
	docs := esResp.Hits.Documents
	var respDocs []SaicKeywords
	for i := range docs {
		var doc SaicKeywords
		json.Unmarshal(docs[i].Source, &doc)
		respDocs = append(respDocs, doc)
	}
	c.JSON(200, respDocs)
}

func queryKWV3(c *gin.Context, index string) {
	kw := c.Param("kw")
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"function_score": map[string]interface{}{
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"should": []map[string]interface{}{
							{
								"match_phrase": map[string]interface{}{
									"ik.keywords": map[string]interface{}{
										"query": kw,
										"boost": 3,
									},
								},
							},
							{
								"match_phrase": map[string]interface{}{
									"ik.name": map[string]interface{}{
										"query": kw,
										"boost": 1,
									},
								},
							},
						},
					},
				},
				"script_score": map[string]interface{}{
					"script": "_score*0.45+doc['fl.integrity_score'].value*0.55",
				},
				"boost_mode": "replace",
			},
		},
		"size":    200,
		"_source": []string{"ik.name", "fl.integrity_score", "kw.id"},
	}
	esResp, err := ESSearchCommon(escli, esQuery, index)
	if err != nil {
		fmt.Println(err)
	}
	docs := esResp.Hits.Documents
	var respDocs []SaicKeywordsV3
	for i := range docs {
		var doc SaicKeywordsV3
		doc.Score = docs[i].Score
		json.Unmarshal(docs[i].Source, &doc)
		respDocs = append(respDocs, doc)
	}
	c.JSON(200, respDocs)
}

func queryKWV4(c *gin.Context, index string) {
	kw := c.Param("kw")
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"function_score": map[string]interface{}{
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"should": []map[string]interface{}{
							{
								"match_phrase": map[string]interface{}{
									"ik.keywords": map[string]interface{}{
										"query": kw,
										"boost": 2,
									},
								},
							},
							{
								"match_phrase": map[string]interface{}{
									"ik.verticals": map[string]interface{}{
										"query": kw,
										"boost": 1,
									},
								},
							},
							{
								"match_phrase": map[string]interface{}{
									"ik.name": map[string]interface{}{
										"query": kw,
										"boost": 1,
									},
								},
							},
						},
					},
				},
				"script_score": map[string]interface{}{
					"script": "_score*0.6+doc['fl.integrity_score'].value*0.45",
				},
				"boost_mode": "replace",
			},
		},
		"size":    200,
		"_source": []string{"ik.name", "fl.integrity_score", "kw.id"},
	}
	esResp, err := ESSearchCommon(escli, esQuery, index)
	if err != nil {
		fmt.Println(err)
	}
	docs := esResp.Hits.Documents
	var respDocs []SaicKeywordsV3
	for i := range docs {
		var doc SaicKeywordsV3
		doc.Score = docs[i].Score
		json.Unmarshal(docs[i].Source, &doc)
		respDocs = append(respDocs, doc)
	}
	c.JSON(200, respDocs)
}
