module github.com/smartcontractkit/chainlink/v2

go 1.24.5

require (
	github.com/Depado/ginprom v1.8.0
	github.com/Masterminds/semver/v3 v3.3.1
	github.com/NethermindEth/juno v0.12.5
	github.com/NethermindEth/starknet.go v0.8.0
	github.com/andybalholm/brotli v1.1.1
	github.com/aptos-labs/aptos-go-sdk v1.9.1-0.20250613185448-581cb03acb8f
	github.com/avast/retry-go/v4 v4.6.1
	github.com/btcsuite/btcd/btcec/v2 v2.3.4
	github.com/coreos/go-oidc/v3 v3.11.0
	github.com/cosmos/cosmos-sdk v0.50.13
	github.com/danielkov/gin-helmet v0.0.0-20171108135313-1387e224435e
	github.com/deckarep/golang-set/v2 v2.6.0
	github.com/dominikbraun/graph v0.23.0
	github.com/doyensec/safeurl v0.2.1
	github.com/esote/minmaxheap v1.0.0
	github.com/ethereum/go-ethereum v1.15.3
	github.com/fatih/color v1.18.0
	github.com/fbsobreira/gotron-sdk v0.0.0-20250403083053-2943ce8c759b
	github.com/fxamacker/cbor/v2 v2.7.0
	github.com/gagliardetto/binary v0.8.0
	github.com/gagliardetto/solana-go v1.13.0
	github.com/getsentry/sentry-go v0.27.0
	github.com/gibson042/canonicaljson-go v1.0.3
	github.com/gin-contrib/cors v1.7.2
	github.com/gin-contrib/expvar v0.0.1
	github.com/gin-contrib/sessions v0.0.5
	github.com/gin-contrib/size v0.0.0-20230212012657-e14a14094dc4
	github.com/gin-gonic/gin v1.10.0
	github.com/go-ldap/ldap/v3 v3.4.6
	github.com/go-viper/mapstructure/v2 v2.4.0
	github.com/go-webauthn/webauthn v0.9.4
	github.com/golang-jwt/jwt/v5 v5.2.3
	github.com/google/pprof v0.0.0-20241210010833-40e02aabc2ad
	github.com/google/uuid v1.6.0
	github.com/gorilla/securecookie v1.1.2
	github.com/gorilla/sessions v1.2.2
	github.com/gorilla/websocket v1.5.3
	github.com/grafana/pyroscope-go v1.1.2
	github.com/graph-gophers/dataloader v5.0.0+incompatible
	github.com/graph-gophers/graphql-go v1.5.0
	github.com/hashicorp/go-envparse v0.1.0
	github.com/hashicorp/go-plugin v1.6.3
	github.com/hashicorp/go-retryablehttp v0.7.7
	github.com/hdevalence/ed25519consensus v0.2.0
	github.com/holiman/uint256 v1.3.2
	github.com/imdario/mergo v0.3.16
	github.com/jackc/pgconn v1.14.3
	github.com/jackc/pgx/v4 v4.18.3
	github.com/jmoiron/sqlx v1.4.0
	github.com/jonboulle/clockwork v0.5.0
	github.com/jpillora/backoff v1.0.0
	github.com/kylelemons/godebug v1.1.0
	github.com/leanovate/gopter v0.2.11
	github.com/lib/pq v1.10.9
	github.com/manyminds/api2go v0.0.0-20171030193247-e7b693844a6f
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mr-tron/base58 v1.2.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/onsi/gomega v1.36.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pelletier/go-toml v1.9.5
	github.com/pelletier/go-toml/v2 v2.2.4
	github.com/pkg/errors v0.9.1
	github.com/pressly/goose/v3 v3.21.1
	github.com/prometheus/client_golang v1.22.0
	github.com/prometheus/client_model v0.6.2
	github.com/prometheus/common v0.65.0
	github.com/prometheus/prometheus v0.304.2
	github.com/robfig/cron/v3 v3.0.1
	github.com/rogpeppe/go-internal v1.13.1
	github.com/rs/zerolog v1.33.0
	github.com/scylladb/go-reflectx v1.0.1
	github.com/shirou/gopsutil/v3 v3.24.3
	github.com/shopspring/decimal v1.4.0
	github.com/sigurn/crc16 v0.0.0-20211026045750-20ab5afb07e3
	github.com/smartcontractkit/chain-selectors v1.0.67
	github.com/smartcontractkit/chainlink-aptos v0.0.0-20250818164129-fa2e60d95157
	github.com/smartcontractkit/chainlink-automation v0.8.1
	github.com/smartcontractkit/chainlink-ccip v0.1.1-solana.0.20250825135846-84f0d5167f8f
	github.com/smartcontractkit/chainlink-ccip/chains/solana v0.0.0-20250805210128-7f8a0f403c3a
	github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings v0.0.0-20250804184440-c0506474fc44
	github.com/smartcontractkit/chainlink-common v0.9.4
	github.com/smartcontractkit/chainlink-data-streams v0.1.2
	github.com/smartcontractkit/chainlink-evm v0.3.2
	github.com/smartcontractkit/chainlink-evm/gethwrappers v0.0.0-20250827130336-5922343458be
	github.com/smartcontractkit/chainlink-feeds v0.1.2-0.20250227211209-7cd000095135
	github.com/smartcontractkit/chainlink-framework/capabilities v0.0.0-20250818175541-3389ac08a563
	github.com/smartcontractkit/chainlink-framework/chains v0.0.0-20250717121125-2350c82883e2
	github.com/smartcontractkit/chainlink-framework/multinode v0.0.0-20250729142306-508e798f6a5d
	github.com/smartcontractkit/chainlink-protos/billing/go v0.0.0-20250722175102-6dcdf5122683
	github.com/smartcontractkit/chainlink-protos/cre/go v0.0.0-20250819150450-95ef563f6e6d
	github.com/smartcontractkit/chainlink-protos/orchestrator v0.8.1
	github.com/smartcontractkit/chainlink-protos/storage-service v0.3.0
	github.com/smartcontractkit/chainlink-protos/workflows/go v0.0.0-20250822025801-598d3d86f873
	github.com/smartcontractkit/chainlink-solana v1.1.2-0.20250820135304-632bebc0e802
	github.com/smartcontractkit/chainlink-tron/relayer v0.0.11-0.20250815105909-75499abc4335
	github.com/smartcontractkit/cre-sdk-go v0.5.1-0.20250818135829-9ea58491207f
	github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http v0.5.1-0.20250818141131-0b979c98bab0
	github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron v0.5.1-0.20250818141131-0b979c98bab0
	github.com/smartcontractkit/freeport v0.1.3-0.20250716200817-cb5dfd0e369e
	github.com/smartcontractkit/libocr v0.0.0-20250707144819-babe0ec4e358
	github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin v0.0.0-20241009055228-33d0c0bf38de
	github.com/smartcontractkit/tdh2/go/tdh2 v0.0.0-20250624150019-e49f7e125e6b
	github.com/smartcontractkit/wsrpc v0.8.5-0.20250502134807-c57d3d995945
	github.com/spf13/cast v1.7.1
	github.com/spf13/cobra v1.9.1
	github.com/stretchr/testify v1.10.0
	github.com/theodesp/go-heaps v0.0.0-20190520121037-88e35354fe0a
	github.com/tidwall/gjson v1.18.0
	github.com/ulule/limiter/v3 v3.11.2
	github.com/umbracle/ethgo v0.1.3
	github.com/unrolled/secure v1.13.0
	github.com/urfave/cli v1.22.14
	github.com/wk8/go-ordered-map/v2 v2.1.8
	github.com/xssnick/tonutils-go v1.13.0
	go.dedis.ch/fixbuf v1.0.3
	go.dedis.ch/kyber/v3 v3.1.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.49.0
	go.opentelemetry.io/otel v1.37.0
	go.opentelemetry.io/otel/metric v1.37.0
	go.opentelemetry.io/otel/sdk/metric v1.37.0
	go.opentelemetry.io/otel/trace v1.37.0
	go.uber.org/atomic v1.11.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.40.0
	golang.org/x/exp v0.0.0-20250711185948-6ae5c78190dc
	golang.org/x/mod v0.26.0
	golang.org/x/oauth2 v0.30.0
	golang.org/x/sync v0.16.0
	golang.org/x/term v0.33.0
	golang.org/x/time v0.12.0
	gonum.org/v1/gonum v0.16.0
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.7
	gopkg.in/guregu/null.v4 v4.0.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/utils v0.0.0-20241104163129-6fe5fd82f078
)

require (
	cosmossdk.io/api v0.7.6 // indirect
	cosmossdk.io/collections v0.4.0 // indirect
	cosmossdk.io/core v0.11.0 // indirect
	cosmossdk.io/depinject v1.1.0 // indirect
	cosmossdk.io/errors v1.0.1 // indirect
	cosmossdk.io/log v1.4.1 // indirect
	cosmossdk.io/math v1.4.0 // indirect
	cosmossdk.io/store v1.1.1 // indirect
	cosmossdk.io/x/tx v0.13.7 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.1 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/DataDog/zstd v1.5.6-0.20230824185856-869dae002e5e // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/VictoriaMetrics/fastcache v1.12.2 // indirect
	github.com/XSAM/otelsql v0.37.0 // indirect
	github.com/apache/arrow-go/v18 v18.3.0 // indirect
	github.com/atombender/go-jsonschema v0.16.1-0.20240916205339-a74cd4e2851c // indirect
	github.com/aybabtme/rgbterm v0.0.0-20170906152045-cc83f3b3ce59 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.1.1-0.20220910012023-760eaf8b6816 // indirect
	github.com/bits-and-blooms/bitset v1.20.0 // indirect
	github.com/blendle/zapdriver v1.3.1 // indirect
	github.com/btcsuite/btcd v0.24.2 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.6 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/buger/goterm v1.0.4 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/bytecodealliance/wasmtime-go/v28 v28.0.0 // indirect
	github.com/bytedance/sonic v1.12.3 // indirect
	github.com/bytedance/sonic/loader v0.2.0 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudevents/sdk-go/binding/format/protobuf/v2 v2.16.1 // indirect
	github.com/cloudevents/sdk-go/v2 v2.16.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/fifo v0.0.0-20240606204812-0bbfbd93a7ce // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v1.1.5 // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/coder/websocket v1.8.12 // indirect
	github.com/cometbft/cometbft v0.38.17 // indirect
	github.com/cometbft/cometbft-db v1.0.1 // indirect
	github.com/consensys/gnark-crypto v0.18.0 // indirect
	github.com/containerd/continuity v0.4.3 // indirect
	github.com/cosmos/btcutil v1.0.5 // indirect
	github.com/cosmos/cosmos-db v1.1.1 // indirect
	github.com/cosmos/cosmos-proto v1.0.0-beta.5 // indirect
	github.com/cosmos/go-bip39 v1.0.0 // indirect
	github.com/cosmos/gogoproto v1.7.0 // indirect
	github.com/cosmos/ics23/go v0.11.0 // indirect
	github.com/cosmos/ledger-cosmos-go v0.14.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20240724233137-53bbb0ceb27a // indirect
	github.com/crate-crypto/go-kzg-4844 v1.1.0 // indirect
	github.com/danieljoos/wincred v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/dgraph-io/badger/v4 v4.7.0 // indirect
	github.com/dgraph-io/ristretto/v2 v2.2.0 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.7.0 // indirect
	github.com/ethereum/c-kzg-4844 v1.0.3 // indirect
	github.com/ethereum/go-verkle v0.2.2 // indirect
	github.com/expr-lang/expr v1.17.5 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gagliardetto/treeout v0.1.4 // indirect
	github.com/gagliardetto/utilz v0.1.3 // indirect
	github.com/gedex/inflector v0.0.0-20170307190818-16278e9db813 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.5 // indirect
	github.com/go-jose/go-jose/v4 v4.0.5 // indirect
	github.com/go-json-experiment/json v0.0.0-20250223041408-d3c622f1b874 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.25.0 // indirect
	github.com/go-webauthn/x v0.1.5 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/goccy/go-yaml v1.12.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gofrs/flock v0.12.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/flatbuffers v25.2.10+incompatible // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-tpm v0.9.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/handlers v1.5.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/grafana/pyroscope-go/godeltaprof v0.1.8 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus v1.0.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hako/durafmt v0.0.0-20200710122514-c0fb7b4da026 // indirect
	github.com/hashicorp/go-bexpr v0.1.10 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-metrics v0.5.4 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/hasura/go-graphql-client v0.13.1 // indirect
	github.com/holiman/billy v0.0.0-20240216141850-2abb0c79d3c4 // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/huin/goupnp v1.3.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/invopop/jsonschema v0.13.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.4 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/linxGnu/grocksdb v1.9.3 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/magiconair/properties v1.8.10 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/marcboeker/go-duckdb v1.8.5 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/miekg/dns v1.1.65 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20220423185008-bf980b35cac4 // indirect
	github.com/mitchellh/pointerstructure v1.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/zstdpool-freelist v0.0.0-20201229113212-927304c0c3b1 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oasisprotocol/curve25519-voi v0.0.0-20230904125328-1f23a7beb09a // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/petermattis/goid v0.0.0-20240813172612-4fcff4a6cae7 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pion/dtls/v2 v2.2.12 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/stun/v2 v2.0.0 // indirect
	github.com/pion/transport/v2 v2.2.10 // indirect
	github.com/pion/transport/v3 v3.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/procfs v0.16.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sanity-io/litter v1.5.5 // indirect
	github.com/santhosh-tekuri/jsonschema/v5 v5.3.1 // indirect
	github.com/sasha-s/go-deadlock v0.3.5 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/smartcontractkit/chainlink-common/pkg/chipingress v0.0.1 // indirect
	github.com/smartcontractkit/chainlink-framework/metrics v0.0.0-20250717121125-2350c82883e2 // indirect
	github.com/smartcontractkit/chainlink-protos/rmn/v1.6/go v0.0.0-20250131130834-15e0d4cde2a6 // indirect
	github.com/smartcontractkit/chainlink-protos/svr v1.1.0 // indirect
	github.com/smartcontractkit/grpc-proxy v0.0.0-20240830132753-a7e17fec5ab7 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/stephenlacy/go-ethereum-hdwallet v0.0.0-20230913225845-a4fa94429863 // indirect
	github.com/streamingfast/logging v0.0.0-20230608130331-f22c91403091 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/supranational/blst v0.3.14 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tendermint/go-amino v0.16.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/umbracle/fastrlp v0.0.0-20220527094140-59d5dd30e722 // indirect
	github.com/urfave/cli/v2 v2.27.6 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	github.com/zksync-sdk/zksync2-go v1.1.1-0.20250620124214-2c742ee399c6 // indirect
	github.com/zondax/hid v0.9.2 // indirect
	github.com/zondax/ledger-go v0.14.3 // indirect
	go.dedis.ch/protobuf v1.0.11 // indirect
	go.etcd.io/bbolt v1.4.0 // indirect
	go.mongodb.org/mongo-driver v1.17.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.61.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.11.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.11.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.35.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.35.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.35.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.35.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.35.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.13.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.35.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.35.0 // indirect
	go.opentelemetry.io/otel/log v0.13.0 // indirect
	go.opentelemetry.io/otel/sdk v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.13.0 // indirect
	go.opentelemetry.io/proto/otlp v1.5.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/ratelimit v0.3.1 // indirect
	golang.org/x/arch v0.11.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	golang.org/x/tools v0.35.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/genproto v0.0.0-20250122153221-138b5a5a4fd4 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250519155744-55703ea1f237 // indirect
	gopkg.in/guregu/null.v2 v2.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gotest.tools/v3 v3.5.2 // indirect
	pgregory.net/rapid v1.1.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace github.com/fbsobreira/gotron-sdk => github.com/smartcontractkit/chainlink-tron/relayer/gotron-sdk v0.0.5-0.20250528121202-292529af39df
