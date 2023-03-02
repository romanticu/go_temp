package funcs

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
)

// PerformESQuery return response body
func PerformESQuery() ([]byte, error) {
	query := strToMap()

	index := "company_search"
	esClient := InitElasticsearch()
	var err error

	var reqBody bytes.Buffer
	err = json.NewEncoder(&reqBody).Encode(query)
	if err != nil {
		err = fmt.Errorf("encode query failed, %v", err)
		return nil, err
	}
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(index),
		esClient.Search.WithBody(&reqBody),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithPretty(),
		esClient.Search.WithTimeout(10*time.Second),
	)

	if err != nil {
		err = fmt.Errorf("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			err = fmt.Errorf("Error parsing the response body: %s", err)
		} else {
			err = fmt.Errorf("[%s] %s: %s", res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"])
		}
		return nil, err
	}

	var builder []byte
	buf := make([]byte, 256)
	for {
		n, err := res.Body.Read(buf)
		builder = append(builder, buf[:n]...)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("read error:%+v", err)
			}
			break
		}
	}

	return builder, nil
}

func InitElasticsearch() *elasticsearch.Client {
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

func strToMap() map[string]interface{} {
	filter := `{
		"size": 0,
		"query": {
		  "bool": {
			"must": [
			  {
				"match_phrase": {
				  "keywords": "咖啡"
				}
			  }
			]
		  }
		},
		"aggs": {
		  "deal_nest": {
			"nested": {
			  "path": "deals"
			},
			"aggs": {
			  "count_investor": {
				"cardinality": {
				  "field": "deals.investor_ids"
				}
			  },
			  "count_exits": {
				"filter": {
				  "terms": {
					"deals.deal_type": [
					  "105002201",
					  "105002202",
					  "105002502",
					  "105002905",
					  "105002508",
					  "105002510",
					  "105002511",
					  "105002503"
					]
				  }
				}
			  },
			  "capital_invested": {
				"sum": {
				  "field": "deals.size"
				}
			  },
			  "med_post_val": {
				"percentiles": {
				  "field": "deals.post_money_valuation",
				  "percents": [
					50
				  ]
				}
			  }
			}
		  },
		  "count_company": {
			"cardinality": {
			  "field": "ni.id"
			}
		  }
		}
	  }`

	q := make(map[string]interface{})
	err := json.Unmarshal([]byte(filter), &q)
	if err != nil {
		panic(err)
	}

	return q
}

func GetESField(m map[string]interface{}, path string) interface{} {
	if path == "" {
		return nil
	}

	fields := strings.Split(path, ".")
	lastField := fields[len(fields)-1]
	if len(fields) == 1 {
		return m[lastField]
	}
	fields = fields[:len(fields)-1]

	for _, f := range fields {
		m = m[f].(map[string]interface{})
	}

	return m[lastField]
}

func GetIndexCount(index string, query map[string]interface{}, esClient *elasticsearch.Client) (int, error) {
	q, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}
	fmt.Println("++++")
	req := esapi.CountRequest{
		Index: []string{index},
		Body:    strings.NewReader(string(q)),
	}
	// Perform the request with the client.
	res, err := req.Do(context.Background(), esClient)
	fmt.Println(err, "----")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	fmt.Println(err, "----")
	defer res.Body.Close()
	// if res.IsError() {
	// 	return 0, err
	// }

	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return 0, err
	}
	fmt.Println("---- ", r)
	fmt.Println("++++ ", int(r["count"].(float64)))

	return 0, nil
}
