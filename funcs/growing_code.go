package funcs

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
	"gitlab.mvalley.com/common/adam/pkg/client"
	cfg "gitlab.mvalley.com/common/adam/pkg/config"
	"gitlab.mvalley.com/common/adam/pkg/go_utils"
)

var mutex sync.Mutex

const creationCodeKey = "creation_code_key"

func GenGrowingCode(size int) {
	mutex.Lock()
	fmt.Println("size ", size)
	for i := 0; i < size; i++ {
		t := time.Now().Unix()
		fmt.Println(t, " --- ", NumToBHex(t, 36))
		time.Sleep(time.Second)
	}
	mutex.Unlock()
}

var num2char = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NumToBHex(num int64, n int64) string {
	num_str := ""
	for num != 0 {
		yu := num % n
		num_str = string(num2char[yu]) + num_str
		num = num / n
	}
	return strings.ToUpper(num_str)
}

func AyncGenCode(size int) {
	redisCli, err := client.InitRedis(cfg.RedisConfiguration{
		Host: "192.168.88.122",
		Port: "6379",
	})
	if err != nil {
		panic(err)
	}

	// 如果以后存在多节点，得替换成分布式锁
	mutex.Lock()
	defer func() {
		mutex.Unlock()
	}()
	var codes []string
	t := time.Now().Unix()
	cacheT, err := getInt64(redisCli, creationCodeKey)
	if err != nil {
		panic(err)
	}
	if cacheT != 0 {
		t = cacheT
	}
	for i := 0; i < size; i++ {
		code := go_utils.NumToBHex(t+int64(i), 36)
		codes = append(codes, complementZero(code))
	}
	err = setKey(redisCli, creationCodeKey, t+int64(size)+1, time.Minute)
	if err != nil {
		panic(err)
	}
	fmt.Println(codes)
}

func complementZero(str string) string {
	numZero := 7 - len(str)
	zeroStr := ""
	for i := 0; i < numZero; i++ {
		zeroStr += "0"
	}

	return zeroStr + str
}
func getInt64(redisCli *redis.Client, key string) (int64, error) {
	value, err := redisCli.Get(key).Int64()
	if err != nil && err.Error() != redis.Nil.Error() {
		return 0, err
	}
	return value, nil

}

func setKey(redisCli *redis.Client, key string, value interface{}, expire time.Duration) error {
	err := redisCli.Set(key, value, expire).Err()
	if err != nil {
		return err
	}
	return nil
}
