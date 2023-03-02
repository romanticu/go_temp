package funcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
)

type CommonESResponse struct {
	Hits struct {
		Documents []DocumentContent `json:"hits"`
	} `json:"hits"`
	ScrollID string `json:"_scroll_id"`
}

type DocumentContent struct {
	ID string `json:"_id"`
	//Content es.FactorDocument `json:"_source"`
}

func FirstScrollQuery(index string, esClient *elasticsearch.Client) (string, []string, error) {
	query := map[string]interface{}{
		"size":    10000,
		"_source": "",
	}
	var reqBody bytes.Buffer
	err := json.NewEncoder(&reqBody).Encode(query)
	if err != nil {
		err = fmt.Errorf("encode query failed, %v", err)
		return "", nil, errors.WithStack(err)
	}
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(index),
		esClient.Search.WithBody(&reqBody),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithPretty(),
		esClient.Search.WithTimeout(10*time.Second),
		esClient.Search.WithScroll(time.Minute),
	)
	if err != nil {
		err = fmt.Errorf("error getting response: %s", err)
		return "", nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			err = fmt.Errorf("error parsing the response body: %s", err)
		} else {
			err = fmt.Errorf("[%s] %s: %s", res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"])
		}
		return "", nil, errors.WithStack(err)
	}

	var result CommonESResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		err = errors.WithStack(fmt.Errorf("Error parsing the response body: %s", err))
		return "", nil, err
	}
	var ids []string
	for _, item := range result.Hits.Documents {
		ids = append(ids, item.ID)
	}
	return result.ScrollID, ids, nil
}

func AfterScrollQuery(scrollID string, esClient *elasticsearch.Client) (string, []string, error) {
	res, err := esClient.Scroll(
		esClient.Scroll.WithContext(context.Background()),
		esClient.Scroll.WithPretty(),
		esClient.Scroll.WithScrollID(scrollID),
		esClient.Scroll.WithScroll(time.Minute),
	)
	if err != nil {
		err = fmt.Errorf("error getting response: %s", err)
		return "", nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			err = fmt.Errorf("error parsing the response body: %s", err)
		} else {
			err = fmt.Errorf("[%s] %s: %s", res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"])
		}
		return "", nil, errors.WithStack(err)
	}

	var result CommonESResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		err = fmt.Errorf("error parsing the response body: %s", err)
		return "", nil, err
	}
	var ids []string
	for _, item := range result.Hits.Documents {
		ids = append(ids, item.ID)
	}
	return result.ScrollID, ids, nil
}
