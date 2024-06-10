module github.com/smartcontractkit/chainlink/v2

go 1.21.7

require (
	github.com/Depado/ginprom v1.8.0
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/Masterminds/sprig/v3 v3.2.3
	github.com/NethermindEth/juno v0.3.1
	github.com/NethermindEth/starknet.go v0.7.1-0.20240401080518-34a506f3cfdb
	github.com/XSAM/otelsql v0.27.0
	github.com/avast/retry-go/v4 v4.5.1
	github.com/btcsuite/btcd/btcec/v2 v2.3.2
	github.com/cometbft/cometbft v0.37.2
	github.com/cosmos/cosmos-sdk v0.47.4
	github.com/danielkov/gin-helmet v0.0.0-20171108135313-1387e224435e
	github.com/dominikbraun/graph v0.23.0
	github.com/esote/minmaxheap v1.0.0
	github.com/ethereum/go-ethereum v1.13.8
	github.com/fatih/color v1.16.0
	github.com/fxamacker/cbor/v2 v2.5.0
	github.com/gagliardetto/solana-go v1.8.4
	github.com/getsentry/sentry-go v0.19.0
	github.com/gin-contrib/cors v1.5.0
	github.com/gin-contrib/expvar v0.0.1
	github.com/gin-contrib/sessions v0.0.5
	github.com/gin-contrib/size v0.0.0-20230212012657-e14a14094dc4
	github.com/gin-gonic/gin v1.9.1
	github.com/go-ldap/ldap/v3 v3.4.6
	github.com/go-webauthn/webauthn v0.9.4
	github.com/google/pprof v0.0.0-20231023181126-ff6d637d2a7b
	github.com/google/uuid v1.6.0
	github.com/gorilla/securecookie v1.1.2
	github.com/gorilla/sessions v1.2.2
	github.com/gorilla/websocket v1.5.1
	github.com/grafana/pyroscope-go v1.1.1
	github.com/graph-gophers/dataloader v5.0.0+incompatible
	github.com/graph-gophers/graphql-go v1.3.0
	github.com/hashicorp/consul/sdk v0.16.0
	github.com/hashicorp/go-envparse v0.1.0
	github.com/hashicorp/go-plugin v1.6.0
	github.com/hashicorp/go-retryablehttp v0.7.5
	github.com/hdevalence/ed25519consensus v0.1.0
	github.com/jackc/pgconn v1.14.3
	github.com/jackc/pgtype v1.14.0
	github.com/jackc/pgx/v4 v4.18.2
	github.com/jmoiron/sqlx v1.3.5
	github.com/jonboulle/clockwork v0.4.0
	github.com/jpillora/backoff v1.0.0
	github.com/kylelemons/godebug v1.1.0
	github.com/leanovate/gopter v0.2.10-0.20210127095200-9abe2343507a
	github.com/lib/pq v1.10.9
	github.com/manyminds/api2go v0.0.0-20171030193247-e7b693844a6f
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mr-tron/base58 v1.2.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/onsi/gomega v1.30.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pelletier/go-toml v1.9.5
	github.com/pelletier/go-toml/v2 v2.1.1
	github.com/pkg/errors v0.9.1
	github.com/pressly/goose/v3 v3.16.0
	github.com/prometheus/client_golang v1.17.0
	github.com/prometheus/client_model v0.5.0
	github.com/prometheus/common v0.45.0
	github.com/prometheus/prometheus v0.48.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/rogpeppe/go-internal v1.11.0
	github.com/scylladb/go-reflectx v1.0.1
	github.com/shirou/gopsutil/v3 v3.24.3
	github.com/shopspring/decimal v1.3.1
	github.com/smartcontractkit/chain-selectors v1.0.10
	github.com/smartcontractkit/chainlink-automation v1.0.3
	github.com/smartcontractkit/chainlink-common v0.1.7-0.20240607135320-c9bc0a2ac0ce
	github.com/smartcontractkit/chainlink-cosmos v0.4.1-0.20240524214833-c362c2ebbd2d
	github.com/smartcontractkit/chainlink-data-streams v0.0.0-20240220203239-09be0ea34540
	github.com/smartcontractkit/chainlink-feeds v0.0.0-20240522213638-159fb2d99917
	github.com/smartcontractkit/chainlink-solana v1.0.3-0.20240605170242-555ff582f36a
	github.com/smartcontractkit/chainlink-starknet/relayer v0.0.1-beta-test.0.20240531021326-99118e47f696
	github.com/smartcontractkit/chainlink-vrf v0.0.0-20231120191722-fef03814f868
	github.com/smartcontractkit/libocr v0.0.0-20240419185742-fd3cab206b2c
	github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin v0.0.0-20230906073235-9e478e5e19f1
	github.com/smartcontractkit/tdh2/go/tdh2 v0.0.0-20230906073235-9e478e5e19f1
	github.com/smartcontractkit/wsrpc v0.8.1
	github.com/spf13/cast v1.6.0
	github.com/stretchr/testify v1.9.0
	github.com/theodesp/go-heaps v0.0.0-20190520121037-88e35354fe0a
	github.com/tidwall/gjson v1.17.0
	github.com/ugorji/go/codec v1.2.12
	github.com/ulule/limiter/v3 v3.11.2
	github.com/umbracle/ethgo v0.1.3
	github.com/unrolled/secure v1.13.0
	github.com/urfave/cli v1.22.14
	go.dedis.ch/fixbuf v1.0.3
	go.dedis.ch/kyber/v3 v3.1.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.49.0
	go.opentelemetry.io/otel v1.24.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.26.0
	golang.org/x/crypto v0.22.0
	golang.org/x/exp v0.0.0-20240213143201-ec583247a57a
	golang.org/x/mod v0.15.0
	golang.org/x/net v0.24.0
	golang.org/x/sync v0.6.0
	golang.org/x/term v0.19.0
	golang.org/x/text v0.14.0
	golang.org/x/time v0.5.0
	golang.org/x/tools v0.18.0
	gonum.org/v1/gonum v0.14.0
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.33.0
	gopkg.in/guregu/null.v4 v4.0.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5 // indirect
	cosmossdk.io/api v0.3.1 // indirect
	cosmossdk.io/core v0.5.1 // indirect
	cosmossdk.io/depinject v1.0.0-alpha.3 // indirect
	cosmossdk.io/errors v1.0.0 // indirect
	cosmossdk.io/math v1.0.1 // indirect
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.1 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/ChainSafe/go-schnorrkel v0.0.0-20200405005733-88cbf1b4c40d // indirect
	github.com/CosmWasm/wasmd v0.40.1 // indirect
	github.com/CosmWasm/wasmvm v1.2.4 // indirect
	github.com/DataDog/zstd v1.5.2 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/VictoriaMetrics/fastcache v1.12.1 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.1.1-0.20220910012023-760eaf8b6816 // indirect
	github.com/bits-and-blooms/bitset v1.10.0 // indirect
	github.com/blendle/zapdriver v1.3.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/bytedance/sonic v1.10.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.0 // indirect
	github.com/cockroachdb/errors v1.9.1 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v0.0.0-20230928194634-aa077af62593 // indirect
	github.com/cockroachdb/redact v1.1.3 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/cometbft/cometbft-db v0.7.0 // indirect
	github.com/confio/ics23/go v0.9.0 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.12.1 // indirect
	github.com/cosmos/btcutil v1.0.5 // indirect
	github.com/cosmos/cosmos-proto v1.0.0-beta.2 // indirect
	github.com/cosmos/go-bip39 v1.0.0 // indirect
	github.com/cosmos/gogoproto v1.4.11 // indirect
	github.com/cosmos/iavl v0.20.0 // indirect
	github.com/cosmos/ibc-go/v7 v7.0.1 // indirect
	github.com/cosmos/ics23/go v0.9.1-0.20221207100636-b1abd8678aab // indirect
	github.com/cosmos/ledger-cosmos-go v0.12.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20231025140028-3c0104f4b233 // indirect
	github.com/crate-crypto/go-kzg-4844 v0.7.0 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/deckarep/golang-set/v2 v2.3.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/dfuse-io/logging v0.0.0-20210109005628-b97a57253f70 // indirect
	github.com/dgraph-io/badger/v2 v2.2007.4 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.7.0 // indirect
	github.com/ethereum/c-kzg-4844 v0.4.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gagliardetto/binary v0.7.7 // indirect
	github.com/gagliardetto/treeout v0.1.4 // indirect
	github.com/gballet/go-libpcsclite v0.0.0-20191108122812-4678299bea08 // indirect
	github.com/gballet/go-verkle v0.1.1-0.20231031103413-a67434b50f46 // indirect
	github.com/gedex/inflector v0.0.0-20170307190818-16278e9db813 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.5 // indirect
	github.com/go-json-experiment/json v0.0.0-20231102232822-2e55bd4e08b0 // indirect
	github.com/go-kit/kit v0.12.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.15.5 // indirect
	github.com/go-webauthn/x v0.1.5 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.3 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/golang/glog v1.1.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-tpm v0.9.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/grafana/pyroscope-go/godeltaprof v0.1.6 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.0-rc.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.0 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/gtank/merlin v0.1.1 // indirect
	github.com/gtank/ristretto255 v0.1.2 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/golang-lru v0.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/holiman/billy v0.0.0-20230718173358-1c7e68d277a7 // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/holiman/uint256 v1.2.4 // indirect
	github.com/huandu/skiplist v1.2.0 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/huin/goupnp v1.3.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/invopop/jsonschema v0.12.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.3 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/libp2p/go-buffer-pool v0.1.0 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/mimoo/StrobeGo v0.0.0-20210601165009-122bf33a46e0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/zstdpool-freelist v0.0.0-20201229113212-927304c0c3b1 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/mwitkow/grpc-proxy v0.0.0-20230212185441-f345521cb9c9 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/petermattis/goid v0.0.0-20230317030725-371a4b8eda08 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v5 v5.3.1 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cobra v1.8.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.15.0 // indirect
	github.com/status-im/keycard-go v0.2.0 // indirect
	github.com/streamingfast/logging v0.0.0-20220405224725-2755dab2ce75 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/supranational/blst v0.3.11 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tecbot/gorocksdb v0.0.0-20191217155057-f0fad39f321c // indirect
	github.com/tendermint/go-amino v0.16.0 // indirect
	github.com/teris-io/shortid v0.0.0-20201117134242-e59966efd125 // indirect
	github.com/test-go/testify v1.1.4 // indirect
	github.com/tidwall/btree v1.6.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/umbracle/fastrlp v0.0.0-20220527094140-59d5dd30e722 // indirect
	github.com/valyala/fastjson v1.4.1 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zondax/hid v0.9.1 // indirect
	github.com/zondax/ledger-go v0.14.1 // indirect
	go.dedis.ch/protobuf v1.0.11 // indirect
	go.etcd.io/bbolt v1.3.7 // indirect
	go.mongodb.org/mongo-driver v1.15.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.46.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/sdk v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	go.uber.org/ratelimit v0.3.0 // indirect
	golang.org/x/arch v0.7.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	google.golang.org/api v0.149.0 // indirect
	google.golang.org/genproto v0.0.0-20231030173426-d783a09b4405 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231016165738-49dd2c1f3d0b // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
	gopkg.in/guregu/null.v2 v2.1.2 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	pgregory.net/rapid v0.5.5 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace (
	// replicating the replace directive on cosmos SDK
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

	// until merged upstream: https://github.com/hashicorp/go-plugin/pull/257
	github.com/hashicorp/go-plugin => github.com/smartcontractkit/go-plugin v0.0.0-20240208201424-b3b91517de16

	// until merged upstream: https://github.com/mwitkow/grpc-proxy/pull/69
	github.com/mwitkow/grpc-proxy => github.com/smartcontractkit/grpc-proxy v0.0.0-20230731113816-f1be6620749f

)
