module github.com/fgq/go_temp

go 1.13

require (
	github.com/elastic/go-elasticsearch v0.0.0
	github.com/elastic/go-elasticsearch/v7 v7.11.0
	github.com/gin-gonic/gin v1.7.7
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-ego/gse v0.69.15
	github.com/go-redis/redis/v7 v7.4.1
	github.com/hashicorp/consul/api v1.12.0
	github.com/jinzhu/gorm v1.9.16
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v1.2.0
	github.com/spf13/viper v1.7.1
	gitlab.mvalley.com/adam/common v0.3.47-dev
	gitlab.mvalley.com/common/adam v0.0.176-release
	gitlab.mvalley.com/common/cain v0.0.74-dev
	gitlab.mvalley.com/common/disco v0.0.26-release
	gitlab.mvalley.com/datapack/cain v0.2.67-dev
	gitlab.mvalley.com/rime-index/common v0.1.63-release
	gitlab.mvalley.com/rime-index/quick-search-gateway v0.0.20-release
	go.mongodb.org/mongo-driver v1.4.6
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gorm.io/driver/mysql v1.0.4
	gorm.io/gorm v1.21.12
)

replace (
	github.com/hashicorp/consul/api v1.5.0 => /home/fgq/work/test/consul-1.5.0/api
	gitlab.mvalley.com/common/disco v0.0.20-release => /home/fgq/work/test/disco-v0.0.20-release
	gitlab.mvalley.com/rime-index/quick-search-gateway v0.0.20-release => /home/fgq/work/rx-workspace/quick-search-gateway
)
