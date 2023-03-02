package funcs

import (
	"fmt"

	"gitlab.mvalley.com/common/adam/pkg/client"
	cfg "gitlab.mvalley.com/common/adam/pkg/config"
)

func SaveStringJSON() {
	// t := time.Now()
	// et := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1)
	cli, err := client.InitRedis(cfg.RedisConfiguration{
		Host: "127.0.0.1",
		Port: "6379",
	})
	if err != nil {
		panic(err)
	}
	// fmt.Println()
	ret := cli.MGet("ccc", "aaa", "bbb")
	rets, _ := ret.Result()
	for _, item := range rets {
		fmt.Println(item)
	}
}
