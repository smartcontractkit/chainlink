package src

import (
	"encoding/json"
	"os"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type generateLocalOCR3Config struct{}

func NewGenerateLocalOCR3ConfigCommand() *generateLocalOCR3Config {
	return &generateLocalOCR3Config{}
}

func (g *generateLocalOCR3Config) Name() string {
	return "generate-local-ocr3-config"
}

func (g *generateLocalOCR3Config) Run(args []string) {
	publicKeys := []byte(`[
	{
		"EthAddress": "0xF4e7e516146c8567F8E8be0ED1f1A92798628d35",
		"P2PPeerID": "12D3KooWNmhKZL1XW4Vv3rNjLXzJ6mqcVerihdijjGYuexPrFUFZ",
		"OCR2BundleID": "2f92c96da20fbe39c89e59516e3a7473254523316887394e406527c72071d3db",
		"OCR2OnchainPublicKey": "a2402db8e549f094ea31e1c0edd77623f4ca5b12",
		"OCR2OffchainPublicKey": "3ca9918cd2787de8f9aff91f220f30a5cc54c394f73e173b12c93368bd7072ad",
		"OCR2ConfigPublicKey": "19904debd03994fe9ea411cda7a6b2f01f20a3fe803df0fed67aaf00cc99113f",
		"CSAPublicKey": "csa_dbae6965bad0b0fa95ecc34a602eee1c0c570ddc29b56502e400d18574b8c3df"
	},
	{
		"EthAddress": "0x8B60FDcc9CAC8ea476b31d17011CB204471431d9",
		"P2PPeerID": "12D3KooWFUjV73ZYkAMhS2cVwte3kXDWD8Ybyx3u9CEDHNoeEhBH",
		"OCR2BundleID": "b3df4d8748b67731a1112e8b45a764941974f5590c93672eebbc4f3504dd10ed",
		"OCR2OnchainPublicKey": "4af19c802b244d1d085492c3946391c965e10519",
		"OCR2OffchainPublicKey": "365b9e1c3c945fc3f51afb25772f0a5a1f1547935a4b5dc89c012f590709fefe",
		"OCR2ConfigPublicKey": "15ff12569d11b8ff9f17f8999ea928d03a439f3fb116661cbc4669a0a3192775",
		"CSAPublicKey": "csa_c5cc655a9c19b69626519c4a72c44a94a3675daeba9c16cc23e010a7a6dac1be"
	},
	{
		"EthAddress": "0x6620F516F29979B214e2451498a057FDd3a0A85d",
		"P2PPeerID": "12D3KooWRTtH2WWrztD87Do1kXePSmGjyU4r7mZVWThmqTGgdbUC",
		"OCR2BundleID": "38459ae37f29f2c1fde0f25972a973322be8cada82acf43f464756836725be97",
		"OCR2OnchainPublicKey": "61925685d2b80b121537341d063c4e57b2f9323c",
		"OCR2OffchainPublicKey": "7fe2dbd9f9fb96f7dbbe0410e32d435ad67dae6c91410189fe5664cf3057ef10",
		"OCR2ConfigPublicKey": "2f02fd80b362e1c7acf91680fd48c062718233acd595a6ae7cbe434e118e6a4f",
		"CSAPublicKey": "csa_7407fc90c70895c0fb2bdf385e2e4918364bec1f7a74bad7fdf696bffafbcab8"
	},
	{
		"EthAddress": "0xFeB61E22FCf4F9740c9D96b05199F195bd61A7c2",
		"P2PPeerID": "12D3KooWMTZnZtcVK4EJsjkKsV9qXNoNRSjT62CZi3tKkXGaCsGh",
		"OCR2BundleID": "b5dbc4c9da983cddde2e3226b85807eb7beaf818694a22576af4d80f352702ed",
		"OCR2OnchainPublicKey": "fd97efd53fc20acc098fcd746c04d8d7540d97e0",
		"OCR2OffchainPublicKey": "91b393bb5e6bd6fd9de23845bcd0e0d9b0dd28a1d65d3cfb1fce9f91bd3d8c19",
		"OCR2ConfigPublicKey": "09eb53924ff8b33a08b4eae2f3819015314ce6e8864ac4f86e97caafd4181506",
		"CSAPublicKey": "csa_ef55caf17eefc2a9d547b5a3978d396bd237c73af99cd849a4758701122e3cba"
		},
		{
		"EthAddress": "0x882Fd04D78A7e7D386Dd5b550f19479E5494B0B2",
		"P2PPeerID": "12D3KooWRsM9yordRQDhLgbErH8WMMGz1bC1J4hR5gAGvMWu8goN",
		"OCR2BundleID": "260d5c1a618cdf5324509d7db95f5a117511864ebb9e1f709e8969339eb225af",
		"OCR2OnchainPublicKey": "a0b67dc5345a71d02b396147ae2cb75dda63cbe9",
		"OCR2OffchainPublicKey": "4f42ef42e5cc351dbbd79c29ef33af25c0250cac84837c1ff997bc111199d07e",
		"OCR2ConfigPublicKey": "3b90249731beb9e4f598371f0b96c3babf47bcc62121ebc9c195e3c33e4fd708",
		"CSAPublicKey": "csa_1b874ac2d54b966cec5a8358678ca6f030261aabf3372ce9dbea2d4eb9cdab3d"
	}]`)

	var pubKeys []changeset.NodeKeys
	err := json.Unmarshal(publicKeys, &pubKeys)
	if err != nil {
		panic(err)
	}

	config := []byte(`{
		"MaxQueryLengthBytes": 1000000,
		"MaxObservationLengthBytes": 1000000,
		"MaxReportLengthBytes": 1000000,
		"MaxBatchSize": 20,
		"UniqueReports": true,
		"DeltaProgressMillis": 5000,
		"DeltaResendMillis": 5000,
		"DeltaInitialMillis": 5000,
		"DeltaRoundMillis": 2000,
		"DeltaGraceMillis": 500,
		"DeltaCertifiedCommitRequestMillis": 1000,
		"DeltaStageMillis": 30000,
		"MaxRoundsPerEpoch": 10,
		"TransmissionSchedule": [4],
		"MaxDurationQueryMillis": 1000,
		"MaxDurationObservationMillis": 1000,
		"MaxDurationReportMillis": 1000,
		"MaxDurationShouldAcceptMillis": 1000,
		"MaxDurationShouldTransmitMillis": 1000,
		"MaxFaultyOracles": 1}`)
	var cfg changeset.OracleConfig
	err = json.Unmarshal(config, &cfg)
	if err != nil {
		panic(err)
	}

	ocrConfig, err := changeset.GenerateOCR3Config(cfg, pubKeys, deployment.XXXGenerateTestOCRSecrets())
	if err != nil {
		panic(err)
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(ocrConfig, "", "  ")
	if err != nil {
		panic(err)
	}

	// Generate a json file in `./OCR3Config.json`
	err = os.WriteFile("./OCR3LocalConfig.json", jsonData, 0600)
	if err != nil {
		panic(err)
	}
}
