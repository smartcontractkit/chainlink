package audit

type EventID string

// Static audit log event type constants
const (
	AuthLoginFailedEmail    EventID = "AUTH_LOGIN_FAILED_EMAIL"
	AuthLoginFailedPassword EventID = "AUTH_LOGIN_FAILED_PASSWORD"
	AuthLoginFailed2FA      EventID = "AUTH_LOGIN_FAILED_2FA"
	AuthLoginSuccessWith2FA EventID = "AUTH_LOGIN_SUCCESS_WITH_2FA"
	AuthLoginSuccessNo2FA   EventID = "AUTH_LOGIN_SUCCESS_NO_2FA"
	Auth2FAEnrolled         EventID = "AUTH_2FA_ENROLLED"
	AuthSessionDeleted      EventID = "SESSION_DELETED"

	PasswordResetAttemptFailedMismatch EventID = "PASSWORD_RESET_ATTEMPT_FAILED_MISMATCH"
	PasswordResetSuccess               EventID = "PASSWORD_RESET_SUCCESS"

	APITokenCreateAttemptPasswordMismatch EventID = "API_TOKEN_CREATE_ATTEMPT_PASSWORD_MISMATCH"
	APITokenCreated                       EventID = "API_TOKEN_CREATED"
	APITokenDeleteAttemptPasswordMismatch EventID = "API_TOKEN_DELETE_ATTEMPT_PASSWORD_MISMATCH"
	APITokenDeleted                       EventID = "API_TOKEN_DELETED"

	FeedsManCreated EventID = "FEEDS_MAN_CREATED"
	FeedsManUpdated EventID = "FEEDS_MAN_UPDATED"

	FeedsManChainConfigCreated EventID = "FEEDS_MAN_CHAIN_CONFIG_CREATED"
	FeedsManChainConfigUpdated EventID = "FEEDS_MAN_CHAIN_CONFIG_UPDATED"
	FeedsManChainConfigDeleted EventID = "FEEDS_MAN_CHAIN_CONFIG_DELETED"

	CSAKeyCreated  EventID = "CSA_KEY_CREATED"
	CSAKeyImported EventID = "CSA_KEY_IMPORTED"
	CSAKeyExported EventID = "CSA_KEY_EXPORTED"
	CSAKeyDeleted  EventID = "CSA_KEY_DELETED"

	OCRKeyBundleCreated  EventID = "OCR_KEY_BUNDLE_CREATED"
	OCRKeyBundleImported EventID = "OCR_KEY_BUNDLE_IMPORTED"
	OCRKeyBundleExported EventID = "OCR_KEY_BUNDLE_EXPORTED"
	OCRKeyBundleDeleted  EventID = "OCR_KEY_BUNDLE_DELETED"

	OCR2KeyBundleCreated  EventID = "OCR2_KEY_BUNDLE_CREATED"
	OCR2KeyBundleImported EventID = "OCR2_KEY_BUNDLE_IMPORTED"
	OCR2KeyBundleExported EventID = "OCR2_KEY_BUNDLE_EXPORTED"
	OCR2KeyBundleDeleted  EventID = "OCR2_KEY_BUNDLE_DELETED"

	KeyCreated  EventID = "KEY_CREATED"
	KeyUpdated  EventID = "KEY_UPDATED"
	KeyImported EventID = "KEY_IMPORTED"
	KeyExported EventID = "KEY_EXPORTED"
	KeyDeleted  EventID = "KEY_DELETED"

	EthTransactionCreated    EventID = "ETH_TRANSACTION_CREATED"
	CosmosTransactionCreated EventID = "COSMOS_TRANSACTION_CREATED"
	SolanaTransactionCreated EventID = "SOLANA_TRANSACTION_CREATED"

	JobCreated EventID = "JOB_CREATED"
	JobDeleted EventID = "JOB_DELETED"

	ChainAdded       EventID = "CHAIN_ADDED"
	ChainSpecUpdated EventID = "CHAIN_SPEC_UPDATED"
	ChainDeleted     EventID = "CHAIN_DELETED"

	ChainRpcNodeAdded   EventID = "CHAIN_RPC_NODE_ADDED"
	ChainRpcNodeDeleted EventID = "CHAIN_RPC_NODE_DELETED"

	BridgeCreated EventID = "BRIDGE_CREATED"
	BridgeUpdated EventID = "BRIDGE_UPDATED"
	BridgeDeleted EventID = "BRIDGE_DELETED"

	ForwarderCreated EventID = "FORWARDER_CREATED"
	ForwarderDeleted EventID = "FORWARDER_DELETED"

	ExternalInitiatorCreated EventID = "EXTERNAL_INITIATOR_CREATED"
	ExternalInitiatorDeleted EventID = "EXTERNAL_INITIATOR_DELETED"

	JobProposalSpecApproved EventID = "JOB_PROPOSAL_SPEC_APPROVED"
	JobProposalSpecUpdated  EventID = "JOB_PROPOSAL_SPEC_UPDATED"
	JobProposalSpecCanceled EventID = "JOB_PROPOSAL_SPEC_CANCELED"
	JobProposalSpecRejected EventID = "JOB_PROPOSAL_SPEC_REJECTED"

	ConfigUpdated            EventID = "CONFIG_UPDATED"
	ConfigSqlLoggingEnabled  EventID = "CONFIG_SQL_LOGGING_ENABLED"
	ConfigSqlLoggingDisabled EventID = "CONFIG_SQL_LOGGING_DISABLED"
	GlobalLogLevelSet        EventID = "GLOBAL_LOG_LEVEL_SET"

	JobErrorDismissed EventID = "JOB_ERROR_DISMISSED"
	JobRunSet         EventID = "JOB_RUN_SET"

	EnvNoncriticalEnvDumped EventID = "ENV_NONCRITICAL_ENV_DUMPED"

	UnauthedRunResumed EventID = "UNAUTHED_RUN_RESUMED"
)
