package ccip

//go:generate go run ../generation/wrap.go ccip Router router
//go:generate go run ../generation/wrap.go ccip CCIPHome ccip_home
//go:generate go run ../generation/wrap.go ccip OnRamp onramp
//go:generate go run ../generation/wrap.go ccip OffRamp offramp
//go:generate go run ../generation/wrap.go ccip FeeQuoter fee_quoter
//go:generate go run ../generation/wrap.go ccip NonceManager nonce_manager
//go:generate go run ../generation/wrap.go ccip MultiAggregateRateLimiter multi_aggregate_rate_limiter
//go:generate go run ../generation/wrap.go ccip TokenAdminRegistry token_admin_registry
//go:generate go run ../generation/wrap.go ccip RegistryModuleOwnerCustom registry_module_owner_custom
//go:generate go run ../generation/wrap.go ccip RMNProxy rmn_proxy_contract
//go:generate go run ../generation/wrap.go ccip RMNRemote rmn_remote
//go:generate go run ../generation/wrap.go ccip RMNHome rmn_home

// Pools
//go:generate go run ../generation/wrap.go ccip BurnMintTokenPool burn_mint_token_pool
//go:generate go run ../generation/wrap.go ccip BurnFromMintTokenPool burn_from_mint_token_pool
//go:generate go run ../generation/wrap.go ccip BurnWithFromMintTokenPool burn_with_from_mint_token_pool
//go:generate go run ../generation/wrap.go ccip LockReleaseTokenPool lock_release_token_pool
//go:generate go run ../generation/wrap.go ccip TokenPool token_pool
//go:generate go run ../generation/wrap.go ccip USDCTokenPool usdc_token_pool

// Helpers
//go:generate go run ../generation/wrap.go ccip MaybeRevertMessageReceiver maybe_revert_message_receiver
//go:generate go run ../generation/wrap.go ccip PingPongDemo ping_pong_demo
//go:generate go run ../generation/wrap.go ccip MessageHasher message_hasher
//go:generate go run ../generation/wrap.go ccip MultiOCR3Helper multi_ocr3_helper
//go:generate go run ../generation/wrap.go ccip USDCReaderTester usdc_reader_tester
//go:generate go run ../generation/wrap.go ccip ReportCodec report_codec
//go:generate go run ../generation/wrap.go ccip EtherSenderReceiver ether_sender_receiver
//go:generate go run ../generation/wrap.go ccip WETH9 weth9
//go:generate go run ../generation/wrap.go ccip MockE2EUSDCTokenMessenger mock_usdc_token_messenger
//go:generate go run ../generation/wrap.go ccip MockE2EUSDCTransmitter mock_usdc_token_transmitter
//go:generate go run ../generation/wrap.go ccip CCIPReaderTester ccip_reader_tester

// EncodingUtils
//go:generate go run ../generation/wrap.go ccip EncodingUtils ccip_encoding_utils
