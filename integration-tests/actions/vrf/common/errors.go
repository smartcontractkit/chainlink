package common

const (
	ErrGenericFormat             = "%s, err %w"
	ErrNodePrimaryKey            = "error getting node's primary ETH key"
	ErrNodeNewTxKey              = "error creating node's EVM transaction key"
	ErrCreatingProvingKeyHash    = "error creating a keyHash from the proving key"
	ErrRegisteringProvingKey     = "error registering a proving key on Coordinator contract"
	ErrRegisterProvingKey        = "error registering proving keys"
	ErrEncodingProvingKey        = "error encoding proving key"
	ErrDeployBlockHashStore      = "error deploying blockhash store"
	ErrDeployBatchBlockHashStore = "error deploying batch blockhash store"

	ErrABIEncodingFunding      = "error Abi encoding subscriptionID"
	ErrSendingLinkToken        = "error sending Link token"
	ErrCreatingBHSJob          = "error creating BHS job"
	ErrParseJob                = "error parsing job definition"
	ErrSetVRFCoordinatorConfig = "error setting config for VRF Coordinator contract"
	ErrCreateVRFSubscription   = "error creating VRF Subscription"
	ErrAddConsumerToSub        = "error adding consumer to VRF Subscription"
	ErrFundSubWithLinkToken    = "error funding subscription with Link tokens"
	ErrRestartCLNode           = "error restarting CL node"
	ErrWaitTXsComplete         = "error waiting for TXs to complete"
	ErrRequestRandomness       = "error requesting randomness"
	ErrLoadingCoordinator      = "error loading coordinator contract"
	ErrCreatingVRFKey          = "error creating VRF key"

	ErrWaitRandomWordsRequestedEvent = "error waiting for RandomWordsRequested event"
	ErrWaitRandomWordsFulfilledEvent = "error waiting for RandomWordsFulfilled event"
)
