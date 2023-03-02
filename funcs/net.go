package funcs

import (
	"fmt"
	"net"

	"gitlab.mvalley.com/common/disco/pkg/consulcommon"
)

func NetTest() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d",
		consulcommon.GenRandomNetworkPort()))
	if err != nil {
		panic(err)
	}
	fmt.Println(addr.IP)
	fmt.Println(addr.Port)
}
