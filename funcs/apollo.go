package funcs

import (
	"flag"
	"fmt"
	"gitlab.mvalley.com/common/adam/pkg/apollo"
	"gitlab.mvalley.com/common/adam/pkg/config"
)

var configFileName = flag.String("cfn", "config", "name of configs file")
var configFilePath = flag.String("cfp", "./configs", "path of configs file")

var CONFIG ConfigType

type ConfigType struct {
	MySQLConfig MySQLConfiguration
}

type MySQLConfiguration struct {
	FactorMySQLConfig cfg.MySQLConfiguration
}

func Apollo() {
	apolloCli := apollo.Init(&CONFIG)
	_, err := apolloCli.Start()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", CONFIG)
}
