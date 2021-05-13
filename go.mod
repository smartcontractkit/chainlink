module github.com/smartcontractkit/chainlink

go 1.15

require (
	github.com/DATA-DOG/go-txdb v0.1.4
	github.com/Depado/ginprom v1.2.1-0.20200115153638-53bbba851bd8
	github.com/araddon/dateparse v0.0.0-20190622164848-0fb0a474d195
	github.com/bitly/go-simplejson v0.5.0
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/coreos/go-semver v0.3.0
	github.com/danielkov/gin-helmet v0.0.0-20171108135313-1387e224435e
	github.com/ethereum/go-ethereum v1.10.2
	github.com/fatih/color v1.10.0
	github.com/fxamacker/cbor/v2 v2.2.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/expvar v0.0.0-20181230111036-f23b556cc79f
	github.com/gin-contrib/size v0.0.0-20190528085907-355431950c57
	github.com/gin-gonic/contrib v0.0.0-20190526021735-7fb7810ed2a0
	github.com/gin-gonic/gin v1.6.0
	github.com/gobuffalo/packr v1.30.1
	github.com/gofrs/uuid v3.4.0+incompatible
	github.com/google/uuid v1.2.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.4.2
	github.com/jackc/pgconn v1.8.1
	github.com/jackc/pgtype v1.7.0
	github.com/jackc/pgx/v4 v4.11.1-0.20210326152507-88ede6efb5b0
	github.com/jinzhu/gorm v1.9.16
	github.com/jpillora/backoff v1.0.0
	github.com/lib/pq v1.10.1
	github.com/libp2p/go-libp2p-core v0.8.5
	github.com/libp2p/go-libp2p-peerstore v0.2.7
	github.com/manyminds/api2go v0.0.0-20171030193247-e7b693844a6f
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/onsi/gomega v1.11.0
	github.com/pelletier/go-toml v1.9.0
	github.com/peterh/liner v1.2.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/smartcontractkit/libocr v0.0.0-20210507173830-a9684351c102
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tevino/abool v0.0.0-20170917061928-9b9efcf221b5
	github.com/tidwall/gjson v1.7.5
	github.com/tidwall/sjson v1.1.6
	github.com/ulule/limiter v0.0.0-20190417201358-7873d115fc4e
	github.com/unrolled/secure v0.0.0-20190624173513-716474489ad3
	github.com/urfave/cli v1.22.5
	go.dedis.ch/fixbuf v1.0.3
	go.dedis.ch/kyber/v3 v3.0.13
	go.uber.org/atomic v1.7.0
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887 // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1
	golang.org/x/text v0.3.6
	golang.org/x/tools v0.0.0-20201211185031-d93e913c1a58
	gonum.org/v1/gonum v0.9.1
	gopkg.in/guregu/null.v4 v4.0.0
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)

// To fix CVE: c16fb56d-9de6-4065-9fca-d2b4cfb13020
// See https://github.com/dgrijalva/jwt-go/issues/463
// If that happens to get released in a 3.X.X version, we can add a constraint to our go.mod
// for it. If its in 4.X.X, then we need all our transitive deps to upgrade to it.
replace github.com/dgrijalva/jwt-go => github.com/form3tech-oss/jwt-go v3.2.1+incompatible
