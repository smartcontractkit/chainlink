package ccip

//go:generate go run ../generation/wrap.go ccip Router router latest
//go:generate go run ../generation/wrap.go ccip CCIPHome ccip_home latest
//go:generate go run ../generation/wrap.go ccip OnRamp onramp latest
//go:generate go run ../generation/wrap.go ccip OffRamp offramp latest
//go:generate go run ../generation/wrap.go ccip OnRampWithMessageTransformer onramp_with_message_transformer latest
//go:generate go run ../generation/wrap.go ccip OffRampWithMessageTransformer offramp_with_message_transformer latest
//go:generate go run ../generation/wrap.go ccip FeeQuoter fee_quoter latest
//go:generate go run ../generation/wrap.go ccip NonceManager nonce_manager latest
//go:generate go run ../generation/wrap.go ccip MultiAggregateRateLimiter multi_aggregate_rate_limiter latest
//go:generate go run ../generation/wrap.go ccip TokenAdminRegistry token_admin_registry latest
//go:generate go run ../generation/wrap.go ccip RegistryModuleOwnerCustom registry_module_owner_custom latest
//go:generate go run ../generation/wrap.go ccip RMNProxy rmn_proxy_contract latest
//go:generate go run ../generation/wrap.go ccip RMNRemote rmn_remote latest
//go:generate go run ../generation/wrap.go ccip RMNHome rmn_home latest

// Pools
//go:generate go run ../generation/wrap.go ccip BurnMintTokenPool burn_mint_token_pool latest
//go:generate go run ../generation/wrap.go ccip BurnFromMintTokenPool burn_from_mint_token_pool latest
//go:generate go run ../generation/wrap.go ccip BurnWithFromMintTokenPool burn_with_from_mint_token_pool latest
//go:generate go run ../generation/wrap.go ccip LockReleaseTokenPool lock_release_token_pool latest
//go:generate go run ../generation/wrap.go ccip TokenPool token_pool latest
//go:generate go run ../generation/wrap.go ccip USDCTokenPool usdc_token_pool latest
//go:generate go run ../generation/wrap.go ccip SiloedLockReleaseTokenPool siloed_lock_release_token_pool latest
//go:generate go run ../generation/wrap.go ccip BurnToAddressMintTokenPool burn_to_address_mint_token_pool latest

// Helpers
//go:generate go run ../generation/wrap.go ccip MaybeRevertMessageReceiver maybe_revert_message_receiver latest
//go:generate go run ../generation/wrap.go ccip LogMessageDataReceiver log_message_data_receiver latest
//go:generate go run ../generation/wrap.go ccip PingPongDemo ping_pong_demo latest
//go:generate go run ../generation/wrap.go ccip MessageHasher message_hasher latest
//go:generate go run ../generation/wrap.go ccip MultiOCR3Helper multi_ocr3_helper latest
//go:generate go run ../generation/wrap.go ccip USDCReaderTester usdc_reader_tester latest
//go:generate go run ../generation/wrap.go ccip ReportCodec report_codec latest
//go:generate go run ../generation/wrap.go ccip EtherSenderReceiver ether_sender_receiver latest
//go:generate go run ../generation/wrap.go ccip MockE2EUSDCTokenMessenger mock_usdc_token_messenger latest
//go:generate go run ../generation/wrap.go ccip MockE2EUSDCTransmitter mock_usdc_token_transmitter latest
//go:generate go run ../generation/wrap.go ccip CCIPReaderTester ccip_reader_tester latest

// EncodingUtils
//go:generate go run ../generation/wrap.go ccip EncodingUtils ccip_encoding_utils latest
