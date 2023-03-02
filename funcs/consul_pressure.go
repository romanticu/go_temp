package funcs

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	cu "gitlab.mvalley.com/common/adam/pkg/context_utils"
	"gitlab.mvalley.com/common/disco/pkg/sd"
	"gitlab.mvalley.com/rime-index/common/rpc/micro_service"
	"gitlab.mvalley.com/rime-index/quick-search-gateway/rpc/qksh_common"
	qkshgw "gitlab.mvalley.com/rime-index/quick-search-gateway/rpc/qksh_gateway"
	"log"
	"time"
)

func getConsulClient() *api.Client {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", "192.168.88.203", "48500")
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	return client
}

func ConsulPressure() {
	tc := sd.NewServiceConsul(sd.WithConsulAddr("192.168.88.203", "48500"))
	count := 0
	for {
		count++
		log.Println(count)
		twc := sd.NewTwirpHTTPClient(sd.GetDefaultHttpClient(30), tc)
		qk := qkshgw.NewQKSHGatewayProtobufClient(micro_service.QuickSearchGateway, cu.NewTransferContextClient(twc))
		ctx := context.Background()
		ctx = cu.SetUserID(ctx, "923472934793")
		_, err := qk.SearchEntity(ctx, &qksh_common.SearchEntityRequest{})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond * 50)
	}

}
