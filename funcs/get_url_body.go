package funcs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var errList []string

func WriteUrl() {
	connStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=True",
		"root", "root", "192.168.88.11", "3306", "bid_test")
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		fmt.Println(err)
	}
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(20)
	db.LogMode(true)
	var rowKeys []string
	var sources []int
	err = db.Table("tb_object_6501").Select("*").Where("F3_6501 in (829004009,829004010,829004011)").Limit(20000).Pluck("F8_6501",&rowKeys).Pluck("F1_6501", &sources).Error
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(rowKeys))
	t1 := time.Now()
	for i, rowKey := range rowKeys{
		pathStr := "files/"+strconv.Itoa(sources[i])
		pathMake(pathStr)
		time.Sleep(200)
		go func(rowKey string, i int) {
			writeToFile(rowKey, sources[i])
			wg.Done()
			fmt.Println(i)
		}(rowKey, i)

	}
	wg.Wait()
	fmt.Println("------------------------------")
	fmt.Println(len(rowKeys))
	fmt.Println(time.Since(t1).Seconds())
	fmt.Println(time.Since(t1).Seconds()/float64(len(rowKeys)))
	fmt.Println(len(errList))
	fmt.Println(errList)
}

func writeToFile(rowKey string, source int) {
	client := &http.Client{}

	req,err := http.NewRequest("GET", "http://news.windin.com/ns/imagebase/6501/"+ rowKey,nil)
	if err != nil{
		fmt.Println(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.2222")

	resp,err := client.Do(req)
	if err != nil{
		fmt.Println(err)
		errList = append(errList, err.Error())
		return
	}
	defer resp.Body.Close()




	buf := make([]byte, 1024)
	rowKey = strings.ReplaceAll(rowKey, "/", "")
	pathStr := "files/"+strconv.Itoa(source)+"/"
	f, err1 := os.OpenFile(pathStr+rowKey+".html", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)//可读写，追加的方式打开（或创建文件）
	if err1 != nil {
		panic(err1)
		return
	}
	defer f.Close()

	for {
		n, _ := resp.Body.Read(buf)
		if 0 == n {
			break
		}
		f.WriteString(string(buf[:n]))
	}
}


// 判断文件夹是否存在
func pathMake(path string){
	_, err := os.Stat(path)
	if err == nil {
		return
	}
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		panic(err)
	}

}