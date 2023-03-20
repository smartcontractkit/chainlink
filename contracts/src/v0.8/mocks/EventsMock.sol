// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

contract EventsMock {
	event VRFCoordinatorV2_RandomWordsRequested_1( bytes32  indexed  keyHash, uint256   requestId, uint256   preSeed, uint64  indexed  subId, uint16   minimumRequestConfirmations, uint32   callbackGasLimit, uint32   numWords, address  indexed  sender);
	event VRFV2WrapperConsumerExample_WrappedRequestFulfilled_1( uint256   requestId, uint256[]   randomWords, uint256   payment);
	event FunctionsClient_RequestSent_1( bytes32  indexed  id);
	event FunctionsClientExample_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event FunctionsOracleWithInit_Transmitted_1( bytes32   configDigest, uint32   epoch);
	event KeeperRegistryLogic1_3_KeepersUpdated_1( address[]   keepers, address[]   payees);
	event KeeperRegistryLogic1_3_UpkeepReceived_1( uint256  indexed  id, uint256   startingBalance, address   importedFrom);
	event VRFV2TransparentUpgradeableProxy_BeaconUpgraded_1( address  indexed  beacon);
	event VRFV2TransparentUpgradeableProxy_Upgraded_1( address  indexed  implementation);
	event BatchVRFCoordinatorV2_ErrorReturned_1( uint256  indexed  requestId, string   reason);
	event KeeperRegistry1_3_UpkeepMigrated_1( uint256  indexed  id, uint256   remainingBalance, address   destination);
	event KeeperRegistryLogic1_3_UpkeepRegistered_1( uint256  indexed  id, uint32   executeGas, address   admin);
	event OVM_GasPriceOracle_GasPriceUpdated_1( uint256   pizhiq);
	event VRFCoordinatorV2TestHelper_RandomWordsFulfilled_1( uint256  indexed  requestId, uint256   outputSeed, uint96   payment, bool   success);
	event ConfirmedOwnerUpgradeable_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event KeeperRegistryBase2_0_FundsAdded_1( uint256  indexed  id, address  indexed  from, uint96   amount);
	event ProxyAdmin_OwnershipTransferred_1( address  indexed  previousOwner, address  indexed  newOwner);
	event KeeperRegistryLogic2_0_ReorgedUpkeepReport_1( uint256  indexed  id);
	event VRFV2Wrapper_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event KeeperRegistry1_2_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event KeeperRegistryBase2_0_UpkeepGasLimitSet_1( uint256  indexed  id, uint96   gasLimit);
	event TransparentUpgradeableProxy_AdminChanged_1( address   previousAdmin, address   newAdmin);
	event Pausable_Unpaused_1( address   account);
	event CronUpkeepFactory_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event FunctionsBillingRegistry_SubscriptionCanceled_1( uint64  indexed  subscriptionId, address   to, uint256   amount);
	event KeeperRegistry1_2_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event KeeperRegistry2_0_UpkeepAdminTransferred_1( uint256  indexed  id, address  indexed  from, address  indexed  to);
	event KeeperRegistryBase2_0_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event ConfirmedOwnerUpgradeable_Initialized_1( uint8   version);
	event FunctionsBillingRegistry_Unpaused_1( address   account);
	event FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferred_1( uint64  indexed  subscriptionId, address   from, address   to);
	event VRFCoordinatorV2_SubscriptionCanceled_1( uint64  indexed  subId, address   to, uint256   amount);
	event FunctionsOracleWithInit_AuthorizedSendersChanged_1( address[]   senders, address   changedBy);
	event KeeperRegistry1_3_Unpaused_1( address   account);
	event KeeperRegistryBase2_0_UpkeepCanceled_1( uint256  indexed  id, uint64  indexed  atBlockHeight);
	event OVM_GasPriceOracle_L1BaseFeeUpdated_1( uint256   qakwnd);
	event VRFConsumerBaseV2Upgradeable_Initialized_1( uint8   version);
	event KeeperRegistry1_2_KeepersUpdated_1( address[]   keepers, address[]   payees);
	event KeeperRegistryLogic1_3_PaymentWithdrawn_1( address  indexed  keeper, uint256  indexed  amount, address  indexed  to, address   payee);
	event VRFCoordinatorV2TestHelper_SubscriptionCanceled_1( uint64  indexed  subId, address   to, uint256   amount);
	event VRFV2WrapperConsumerExample_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event ERC1967Upgrade_Upgraded_1( address  indexed  implementation);
	event FunctionsBillingRegistryWithInit_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event FunctionsBillingRegistryWithInit_Unpaused_1( address   account);
	event KeeperRegistry1_2_FundsAdded_1( uint256  indexed  id, address  indexed  from, uint96   amount);
	event KeeperRegistryBase1_3_UpkeepAdminTransferred_1( uint256  indexed  id, address  indexed  from, address  indexed  to);
	event KeeperRegistry2_0_Transmitted_1( bytes32   configDigest, uint32   epoch);
	event KeeperRegistryBase2_0_UpkeepPaused_1( uint256  indexed  id);
	event VRFCoordinatorV2_ProvingKeyRegistered_1( bytes32   keyHash, address  indexed  oracle);
	event KeeperRegistryLogic1_3_Unpaused_1( address   account);
	event KeeperRegistryLogic2_0_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event KeeperRegistryLogic2_0_UpkeepPerformed_1( uint256  indexed  id, bool  indexed  success, uint32   checkBlockNumber, uint256   gasUsed, uint256   gasOverhead, uint96   totalPayment);
	event ChainlinkClient_ChainlinkFulfilled_1( bytes32  indexed  id);
	event ERC1967Proxy_BeaconUpgraded_1( address  indexed  beacon);
	event FunctionsOracle_OracleResponse_1( bytes32  indexed  requestId);
	event FunctionsOracleWithInit_OracleRequest_1( bytes32  indexed  requestId, address   requestingContract, address   requestInitiator, uint64   subscriptionId, address   subscriptionOwner, bytes   data);
	event KeeperRegistry1_3_UpkeepPerformed_1( uint256  indexed  id, bool  indexed  success, address  indexed  from, uint96   payment, bytes   performData);
	event VRFCoordinatorV2TestHelper_ConfigSet_1( uint16   minimumRequestConfirmations, uint32   maxGasLimit, uint32   stalenessSeconds, uint32   gasAfterPaymentCalculation, int256   fallbackWeiPerUnitLink, uint32   feeConfig_fulfillmentFlatFeeLinkPPMTier1, uint32   feeConfig_fulfillmentFlatFeeLinkPPMTier2, uint32   feeConfig_fulfillmentFlatFeeLinkPPMTier3, uint32   feeConfig_fulfillmentFlatFeeLinkPPMTier4, uint32   feeConfig_fulfillmentFlatFeeLinkPPMTier5, uint24   feeConfig_reqsForTier2, uint24   feeConfig_reqsForTier3, uint24   feeConfig_reqsForTier4, uint24   feeConfig_reqsForTier5);
	event VRFCoordinatorV2TestHelper_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event Verifier_ReportVerified_1( bytes32  indexed  feedId, bytes32   reportHash, address   requester);
	event VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferred_1( uint64  indexed  subId, address   from, address   to);
	event Verifier_ConfigActivated_1( bytes32  indexed  feedId, bytes32   configDigest);
	event Verifier_ConfigDeactivated_1( bytes32  indexed  feedId, bytes32   configDigest);
	event AuthorizedReceiver_AuthorizedSendersChanged_1( address[]   senders, address   changedBy);
	event KeeperRegistry1_3_Paused_1( address   account);
	event KeeperRegistryBase1_3_FundsWithdrawn_1( uint256  indexed  id, uint256   amount, address   to);
	event KeeperRegistryLogic1_3_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event VRFCoordinatorV2TestHelper_FundsRecovered_1( address   to, uint256   amount);
	event KeeperRegistryLogic2_0_UpkeepPaused_1( uint256  indexed  id);
	event OVM_GasPriceOracle_OwnershipTransferred_1( address  indexed  previousOwner, address  indexed  newOwner);
	event VRFCoordinatorV2TestHelper_ProvingKeyDeregistered_1( bytes32   keyHash, address  indexed  oracle);
	event AggregatorV2V3Interface_NewRound_1( uint256  indexed  roundId, address  indexed  startedBy, uint256   startedAt);
	event FunctionsClient_RequestFulfilled_1( bytes32  indexed  id);
	event FunctionsOracleWithInit_ConfigSet_1( uint32   previousConfigBlockNumber, bytes32   configDigest, uint64   configCount, address[]   signers, address[]   transmitters, uint8   f, bytes   onchainConfig, uint64   offchainConfigVersion, bytes   offchainConfig);
	event KeeperRegistry2_0_FundsAdded_1( uint256  indexed  id, address  indexed  from, uint96   amount);
	event KeeperRegistryBase2_0_UpkeepAdminTransferRequested_1( uint256  indexed  id, address  indexed  from, address  indexed  to);
	event VRFCoordinatorV2TestHelper_RandomWordsRequested_1( bytes32  indexed  keyHash, uint256   requestId, uint256   preSeed, uint64  indexed  subId, uint16   minimumRequestConfirmations, uint32   callbackGasLimit, uint32   numWords, address  indexed  sender);
	event CronUpkeepFactory_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event KeeperRegistryBase1_3_FundsAdded_1( uint256  indexed  id, address  indexed  from, uint96   amount);
	event KeeperRegistryBase2_0_UpkeepAdminTransferred_1( uint256  indexed  id, address  indexed  from, address  indexed  to);
	event KeeperRegistryLogic2_0_InsufficientFundsUpkeepReport_1( uint256  indexed  id);
	event KeeperRegistryLogic2_0_PayeeshipTransferRequested_1( address  indexed  transmitter, address  indexed  from, address  indexed  to);
	event ENSInterface_NewTTL_1( bytes32  indexed  node, uint64   ttl);
	event FunctionsOracleWithInit_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event KeeperRegistryBase2_0_PayeesUpdated_1( address[]   transmitters, address[]   payees);
	event KeeperRegistryLogic2_0_PayeesUpdated_1( address[]   transmitters, address[]   payees);
	event KeeperRegistryLogic2_0_UpkeepAdminTransferred_1( uint256  indexed  id, address  indexed  from, address  indexed  to);
	event KeeperRegistryLogic2_0_UpkeepReceived_1( uint256  indexed  id, uint256   startingBalance, address   importedFrom);
	event KeeperRegistryLogic2_0_UpkeepUnpaused_1( uint256  indexed  id);
	event AggregatorInterface_AnswerUpdated_1( int256  indexed  current, uint256  indexed  roundId, uint256   updatedAt);
	event CronUpkeepFactory_NewCronUpkeepCreated_1( address   upkeep, address   owner);
	event FunctionsOracle_UserCallbackError_1( bytes32  indexed  requestId, string   reason);
	event KeeperRegistry2_0_StaleUpkeepReport_1( uint256  indexed  id);
	event KeeperRegistryLogic2_0_UpkeepGasLimitSet_1( uint256  indexed  id, uint96   gasLimit);
	event FunctionsOracleWithInit_InvalidRequestID_1( bytes32  indexed  requestId);
	event KeeperRegistryBase1_3_UpkeepPerformed_1( uint256  indexed  id, bool  indexed  success, address  indexed  from, uint96   payment, bytes   performData);
	event KeeperRegistryLogic2_0_PaymentWithdrawn_1( address  indexed  transmitter, uint256  indexed  amount, address  indexed  to, address   payee);
	event KeeperRegistryLogic2_0_UpkeepOffchainConfigSet_1( uint256  indexed  id, bytes   offchainConfig);
	event VRFCoordinatorMock_RandomnessRequest_1( address  indexed  sender, bytes32  indexed  keyHash, uint256  indexed  seed);
	event VRFCoordinatorV2_SubscriptionCreated_1( uint64  indexed  subId, address   owner);
	event VerifierProxy_VerifierSet_1( bytes32   oldConfigDigest, bytes32   newConfigDigest, address   verifierAddress);
	event AuthorizedOriginReceiverUpgradeable_AuthorizedSendersActive_1( address   account);
	event FunctionsBillingRegistryWithInit_AuthorizedSendersChanged_1( address[]   senders, address   changedBy);
	event KeeperRegistryBase1_3_UpkeepCanceled_1( uint256  indexed  id, uint64  indexed  atBlockHeight);
	event VRFCoordinatorV2_FundsRecovered_1( address   to, uint256   amount);
	event VRFCoordinatorV2_ProvingKeyDeregistered_1( bytes32   keyHash, address  indexed  oracle);
	event KeeperRegistry1_2_PayeeshipTransferRequested_1( address  indexed  keeper, address  indexed  from, address  indexed  to);
	event KeeperRegistry1_2_UpkeepRegistered_1( uint256  indexed  id, uint32   executeGas, address   admin);
	event KeeperRegistry1_3_UpkeepReceived_1( uint256  indexed  id, uint256   startingBalance, address   importedFrom);
	event AuthorizedOriginReceiverUpgradeable_AuthorizedSendersDeactive_1( address   account);
	event ConfirmedOwnerWithProposal_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event ConfirmedOwnerWithProposal_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event FunctionsBillingRegistryWithInit_BillingEnd_1( bytes32  indexed  requestId, uint64   subscriptionId, uint96   signerPayment, uint96   transmitterPayment, uint96   totalCost, bool   success);
	event FunctionsBillingRegistryWithInit_RequestTimedOut_1( bytes32  indexed  requestId);
	event KeeperRegistry2_0_PayeesUpdated_1( address[]   transmitters, address[]   payees);
	event KeeperRegistry2_0_PayeeshipTransferred_1( address  indexed  transmitter, address  indexed  from, address  indexed  to);
	event KeeperRegistryLogic1_3_UpkeepUnpaused_1( uint256  indexed  id);
	event KeeperRegistryLogic2_0_FundsWithdrawn_1( uint256  indexed  id, uint256   amount, address   to);
	event ENSInterface_NewResolver_1( bytes32  indexed  node, address   resolver);
	event FunctionsOracleWithInit_AuthorizedSendersDeactive_1( address   account);
	event KeeperRegistry1_3_FundsAdded_1( uint256  indexed  id, address  indexed  from, uint96   amount);
	event Verifier_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event OCR2Abstract_ConfigSet_1( uint32   previousConfigBlockNumber, bytes32   configDigest, uint64   configCount, address[]   signers, address[]   transmitters, uint8   f, bytes   onchainConfig, uint64   offchainConfigVersion, bytes   offchainConfig);
	event ENSInterface_NewOwner_1( bytes32  indexed  node, bytes32  indexed  label, address   owner);
	event FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferRequested_1( uint64  indexed  subscriptionId, address   from, address   to);
	event FunctionsClientExample_RequestFulfilled_1( bytes32  indexed  id);
	event KeeperRegistry2_0_Paused_1( address   account);
	event KeeperRegistryLogic1_3_UpkeepGasLimitSet_1( uint256  indexed  id, uint96   gasLimit);
	event VRFCoordinatorV2TestHelper_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event VRFV2TransparentUpgradeableProxy_AdminChanged_1( address   previousAdmin, address   newAdmin);
	event FunctionsBillingRegistry_FundsRecovered_1( address   to, uint256   amount);
	event KeeperRegistry1_3_UpkeepCheckDataUpdated_1( uint256  indexed  id, bytes   newCheckData);
	event KeeperRegistry2_0_UpkeepPaused_1( uint256  indexed  id);
	event KeeperRegistryBase2_0_UpkeepOffchainConfigSet_1( uint256  indexed  id, bytes   offchainConfig);
	event OVM_GasPriceOracle_DecimalsUpdated_1( uint256   nlwwiv);
	event VerifierProxy_VerifierUnset_1( bytes32   configDigest, address   verifierAddress);
	event FunctionsBillingRegistry_Paused_1( address   account);
	event Initializable_Initialized_1( uint8   version);
	event KeeperRegistryLogic1_3_FundsWithdrawn_1( uint256  indexed  id, uint256   amount, address   to);
	event KeeperRegistryLogic1_3_Paused_1( address   account);
	event KeeperRegistryLogic2_0_PayeeshipTransferred_1( address  indexed  transmitter, address  indexed  from, address  indexed  to);
	event KeeperRegistry1_3_UpkeepGasLimitSet_1( uint256  indexed  id, uint96   gasLimit);
	event FunctionsOracle_AuthorizedSendersDeactive_1( address   account);
	event VRFCoordinatorV2TestHelper_ProvingKeyRegistered_1( bytes32   keyHash, address  indexed  oracle);
	event Verifier_ConfigSet_1( bytes32  indexed  feedId, uint32   previousConfigBlockNumber, bytes32   configDigest, uint64   configCount, address[]   signers, bytes32[]   offchainTransmitters, uint8   f, bytes   onchainConfig, uint64   offchainConfigVersion, bytes   offchainConfig);
	event KeeperRegistryBase2_0_PaymentWithdrawn_1( address  indexed  transmitter, uint256  indexed  amount, address  indexed  to, address   payee);
	event KeeperRegistryLogic2_0_UpkeepCanceled_1( uint256  indexed  id, uint64  indexed  atBlockHeight);
	event OCR2BaseUpgradeable_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event AggregatorInterface_NewRound_1( uint256  indexed  roundId, address  indexed  startedBy, uint256   startedAt);
	event AggregatorV2V3Interface_AnswerUpdated_1( int256  indexed  current, uint256  indexed  roundId, uint256   updatedAt);
	event CronUpkeep_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event FunctionsBillingRegistry_RequestTimedOut_1( bytes32  indexed  requestId);
	event FunctionsBillingRegistry_SubscriptionFunded_1( uint64  indexed  subscriptionId, uint256   oldBalance, uint256   newBalance);
	event Ownable_OwnershipTransferred_1( address  indexed  previousOwner, address  indexed  newOwner);
	event VRFCoordinatorV2_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event VRFV2WrapperConsumerExample_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event KeeperRegistryLogic1_3_OwnerFundsWithdrawn_1( uint96   amount);
	event KeeperRegistryLogic2_0_StaleUpkeepReport_1( uint256  indexed  id);
	event VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferRequested_1( uint64  indexed  subId, address   from, address   to);
	event KeeperRegistry1_3_PaymentWithdrawn_1( address  indexed  keeper, uint256  indexed  amount, address  indexed  to, address   payee);
	event KeeperRegistry2_0_UpkeepAdminTransferRequested_1( uint256  indexed  id, address  indexed  from, address  indexed  to);
	event KeeperRegistryBase1_3_UpkeepCheckDataUpdated_1( uint256  indexed  id, bytes   newCheckData);
	event KeeperRegistryLogic2_0_UpkeepMigrated_1( uint256  indexed  id, uint256   remainingBalance, address   destination);
	event FunctionsBillingRegistry_SubscriptionConsumerRemoved_1( uint64  indexed  subscriptionId, address   consumer);
	event KeeperRegistryBase2_0_UpkeepCheckDataUpdated_1( uint256  indexed  id, bytes   newCheckData);
	event KeeperRegistryLogic1_3_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event KeeperRegistryLogic2_0_UpkeepRegistered_1( uint256  indexed  id, uint32   executeGas, address   admin);
	event CronUpkeep_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event FunctionsBillingRegistry_SubscriptionOwnerTransferRequested_1( uint64  indexed  subscriptionId, address   from, address   to);
	event KeeperRegistry1_3_KeepersUpdated_1( address[]   keepers, address[]   payees);
	event KeeperRegistryBase1_3_Unpaused_1( address   account);
	event KeeperRegistryBase1_3_UpkeepReceived_1( uint256  indexed  id, uint256   startingBalance, address   importedFrom);
	event CronUpkeep_CronJobDeleted_1( uint256  indexed  id);
	event ERC1967Upgrade_AdminChanged_1( address   previousAdmin, address   newAdmin);
	event KeeperRegistryBase2_0_PayeeshipTransferRequested_1( address  indexed  transmitter, address  indexed  from, address  indexed  to);
	event KeeperRegistryLogic2_0_UpkeepCheckDataUpdated_1( uint256  indexed  id, bytes   newCheckData);
	event KeeperRegistryBase2_0_UpkeepReceived_1( uint256  indexed  id, uint256   startingBalance, address   importedFrom);
	event KeeperRegistryLogic2_0_Paused_1( address   account);
	event VRFCoordinatorV2TestHelper_SubscriptionFunded_1( uint64  indexed  subId, uint256   oldBalance, uint256   newBalance);
	event FunctionsBillingRegistry_BillingEnd_1( bytes32  indexed  requestId, uint64   subscriptionId, uint96   signerPayment, uint96   transmitterPayment, uint96   totalCost, bool   success);
	event FunctionsBillingRegistryWithInit_Paused_1( address   account);
	event FunctionsOracleWithInit_OracleResponse_1( bytes32  indexed  requestId);
	event FunctionsOracleWithInit_UserCallbackRawError_1( bytes32  indexed  requestId, bytes   lowLevelData);
	event KeeperRegistryBase1_3_Paused_1( address   account);
	event TransparentUpgradeableProxy_Upgraded_1( address  indexed  implementation);
	event FunctionsBillingRegistry_BillingStart_1( bytes32  indexed  requestId, uint64   commitment_subscriptionId, address   commitment_client, uint32   commitment_gasLimit, uint256   commitment_gasPrice, address   commitment_don, uint96   commitment_donFee, uint96   commitment_registryFee, uint96   commitment_estimatedCost, uint256   commitment_timestamp);
	event KeeperRegistryBase1_3_PayeeshipTransferred_1( address  indexed  keeper, address  indexed  from, address  indexed  to);
	event KeeperRegistryLogic2_0_FundsAdded_1( uint256  indexed  id, address  indexed  from, uint96   amount);
	event Pausable_Paused_1( address   account);
	event PausableUpgradeable_Paused_1( address   account);
	event KeeperRegistry2_0_PaymentWithdrawn_1( address  indexed  transmitter, uint256  indexed  amount, address  indexed  to, address   payee);
	event KeeperRegistryBase2_0_PayeeshipTransferred_1( address  indexed  transmitter, address  indexed  from, address  indexed  to);
	event VRFLoadTestExternalSubOwner_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event OCR2BaseUpgradeable_Transmitted_1( bytes32   configDigest, uint32   epoch);
	event VRFLoadTestExternalSubOwner_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event ContextUpgradeable_Initialized_1( uint8   version);
	event FunctionsBillingRegistryWithInit_FundsRecovered_1( address   to, uint256   amount);
	event FunctionsBillingRegistryWithInit_SubscriptionCanceled_1( uint64  indexed  subscriptionId, address   to, uint256   amount);
	event FunctionsClientExample_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event KeeperRegistryLogic1_3_UpkeepAdminTransferred_1( uint256  indexed  id, address  indexed  from, address  indexed  to);
	event OCR2BaseUpgradeable_Initialized_1( uint8   version);
	event PausableUpgradeable_Initialized_1( uint8   version);
	event VRFCoordinatorV2_ConfigSet_1( uint16   minimumRequestConfirmations, uint32   maxGasLimit, uint32   stalenessSeconds, uint32   gasAfterPaymentCalculation, int256   fallbackWeiPerUnitLink, uint32   feeConfig_fulfillmentFlatFeeLinkPPMTier1, uint32   feeConfig_fulfillmentFlatFeeLinkPPMTier2, uint32   feeConfig_fulfillmentFlatFeeLinkPPMTier3, uint32   feeConfig_fulfillmentFlatFeeLinkPPMTier4, uint32   feeConfig_fulfillmentFlatFeeLinkPPMTier5, uint24   feeConfig_reqsForTier2, uint24   feeConfig_reqsForTier3, uint24   feeConfig_reqsForTier4, uint24   feeConfig_reqsForTier5);
	event ERC1967Proxy_Upgraded_1( address  indexed  implementation);
	event FunctionsBillingRegistry_Initialized_1( uint8   version);
	event KeeperRegistry1_2_PaymentWithdrawn_1( address  indexed  keeper, uint256  indexed  amount, address  indexed  to, address   payee);
	event KeeperRegistry2_0_PayeeshipTransferRequested_1( address  indexed  transmitter, address  indexed  from, address  indexed  to);
	event KeeperRegistryLogic1_3_UpkeepMigrated_1( uint256  indexed  id, uint256   remainingBalance, address   destination);
	event KeeperRegistryLogic1_3_PayeeshipTransferRequested_1( address  indexed  keeper, address  indexed  from, address  indexed  to);
	event KeeperRegistryLogic1_3_PayeeshipTransferred_1( address  indexed  keeper, address  indexed  from, address  indexed  to);
	event KeeperRegistryLogic2_0_Unpaused_1( address   account);
	event KeeperRegistryBase2_0_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event KeeperRegistryBase2_0_UpkeepMigrated_1( uint256  indexed  id, uint256   remainingBalance, address   destination);
	event KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event FunctionsBillingRegistry_SubscriptionConsumerAdded_1( uint64  indexed  subscriptionId, address   consumer);
	event FunctionsOracle_InvalidRequestID_1( bytes32  indexed  requestId);
	event FunctionsOracleWithInit_UserCallbackError_1( bytes32  indexed  requestId, string   reason);
	event KeeperRegistry2_0_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event KeeperRegistry2_0_UpkeepUnpaused_1( uint256  indexed  id);
	event KeeperRegistryLogic1_3_UpkeepAdminTransferRequested_1( uint256  indexed  id, address  indexed  from, address  indexed  to);
	event LogEmitter_Log3_1( string   nelhzo);
	event VRFCoordinatorV2_SubscriptionFunded_1( uint64  indexed  subId, uint256   oldBalance, uint256   newBalance);
	event VRFCoordinatorV2TestHelper_SubscriptionConsumerAdded_1( uint64  indexed  subId, address   consumer);
	event ERC1967Upgrade_BeaconUpgraded_1( address  indexed  beacon);
	event KeeperRegistry2_0_FundsWithdrawn_1( uint256  indexed  id, uint256   amount, address   to);
	event KeeperRegistryBase2_0_OwnerFundsWithdrawn_1( uint96   amount);
	event ChainlinkClient_ChainlinkRequested_1( bytes32  indexed  id);
	event FunctionsOracle_Initialized_1( uint8   version);
	event KeeperRegistry1_2_PayeeshipTransferred_1( address  indexed  keeper, address  indexed  from, address  indexed  to);
	event KeeperRegistryBase2_0_UpkeepUnpaused_1( uint256  indexed  id);
	event KeeperRegistryLogic1_3_UpkeepPaused_1( uint256  indexed  id);
	event FunctionsClientExample_RequestSent_1( bytes32  indexed  id);
	event KeeperRegistry2_0_UpkeepPerformed_1( uint256  indexed  id, bool  indexed  success, uint32   checkBlockNumber, uint256   gasUsed, uint256   gasOverhead, uint96   totalPayment);
	event KeeperRegistryBase2_0_FundsWithdrawn_1( uint256  indexed  id, uint256   amount, address   to);
	event LogEmitter_Log1_1( uint256   bevkse);
	event VRFCoordinatorV2TestHelper_SubscriptionCreated_1( uint64  indexed  subId, address   owner);
	event ChainlinkClient_ChainlinkCancelled_1( bytes32  indexed  id);
	event KeeperRegistry1_2_UpkeepCanceled_1( uint256  indexed  id, uint64  indexed  atBlockHeight);
	event KeeperRegistry2_0_OwnerFundsWithdrawn_1( uint96   amount);
	event KeeperRegistry1_3_UpkeepCanceled_1( uint256  indexed  id, uint64  indexed  atBlockHeight);
	event KeeperRegistry2_0_UpkeepOffchainConfigSet_1( uint256  indexed  id, bytes   offchainConfig);
	event KeeperRegistryLogic2_0_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event ConfirmedOwnerUpgradeable_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event ERC1967Proxy_AdminChanged_1( address   previousAdmin, address   newAdmin);
	event FunctionsBillingRegistryWithInit_SubscriptionConsumerRemoved_1( uint64  indexed  subscriptionId, address   consumer);
	event KeeperRegistry1_2_UpkeepMigrated_1( uint256  indexed  id, uint256   remainingBalance, address   destination);
	event KeeperRegistry1_2_UpkeepReceived_1( uint256  indexed  id, uint256   startingBalance, address   importedFrom);
	event KeeperRegistryBase1_3_ConfigSet_1( uint32   config_paymentPremiumPPB, uint32   config_flatFeeMicroLink, uint24   config_blockCountPerTurn, uint32   config_checkGasLimit, uint24   config_stalenessSeconds, uint16   config_gasCeilingMultiplier, uint96   config_minUpkeepSpend, uint32   config_maxPerformGas, uint256   config_fallbackGasPrice, uint256   config_fallbackLinkPrice, address   config_transcoder, address   config_registrar);
	event KeeperRegistryLogic2_0_UpkeepAdminTransferRequested_1( uint256  indexed  id, address  indexed  from, address  indexed  to);
	event PausableUpgradeable_Unpaused_1( address   account);
	event VRFV2ProxyAdmin_OwnershipTransferred_1( address  indexed  previousOwner, address  indexed  newOwner);
	event VRFV2WrapperConsumerExample_WrapperRequestMade_1( uint256  indexed  requestId, uint256   paid);
	event FunctionsBillingRegistry_ConfigSet_1( uint32   maxGasLimit, uint32   stalenessSeconds, uint256   gasAfterPaymentCalculation, int256   fallbackWeiPerUnitLink, uint32   gasOverhead);
	event FunctionsOracle_ConfigSet_1( uint32   previousConfigBlockNumber, bytes32   configDigest, uint64   configCount, address[]   signers, address[]   transmitters, uint8   f, bytes   onchainConfig, uint64   offchainConfigVersion, bytes   offchainConfig);
	event KeeperRegistryBase2_0_UpkeepRegistered_1( uint256  indexed  id, uint32   executeGas, address   admin);
	event KeeperRegistryLogic1_3_ConfigSet_1( uint32   config_paymentPremiumPPB, uint32   config_flatFeeMicroLink, uint24   config_blockCountPerTurn, uint32   config_checkGasLimit, uint24   config_stalenessSeconds, uint16   config_gasCeilingMultiplier, uint96   config_minUpkeepSpend, uint32   config_maxPerformGas, uint256   config_fallbackGasPrice, uint256   config_fallbackLinkPrice, address   config_transcoder, address   config_registrar);
	event KeeperRegistryLogic1_3_UpkeepCanceled_1( uint256  indexed  id, uint64  indexed  atBlockHeight);
	event VerifierProxy_AccessControllerSet_1( address   oldAccessController, address   newAccessController);
	event CronUpkeep_CronJobCreated_1( uint256  indexed  id, address   target, bytes   handler);
	event CronUpkeep_Paused_1( address   account);
	event KeeperRegistry1_2_Unpaused_1( address   account);
	event KeeperRegistry1_3_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event KeeperRegistryBase1_3_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event BatchVRFCoordinatorV2_RawErrorReturned_1( uint256  indexed  requestId, bytes   lowLevelData);
	event KeeperRegistry2_0_UpkeepReceived_1( uint256  indexed  id, uint256   startingBalance, address   importedFrom);
	event KeeperRegistryBase2_0_Unpaused_1( address   account);
	event KeeperRegistryLogic2_0_OwnerFundsWithdrawn_1( uint96   amount);
	event KeeperRegistryBase1_3_PayeeshipTransferRequested_1( address  indexed  keeper, address  indexed  from, address  indexed  to);
	event KeeperRegistryBase2_0_Paused_1( address   account);
	event VerifierProxy_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event KeeperRegistry1_3_OwnerFundsWithdrawn_1( uint96   amount);
	event KeeperRegistry1_3_UpkeepAdminTransferRequested_1( uint256  indexed  id, address  indexed  from, address  indexed  to);
	event KeeperRegistry1_3_UpkeepUnpaused_1( uint256  indexed  id);
	event KeeperRegistry2_0_Unpaused_1( address   account);
	event KeeperRegistryLogic2_0_CancelledUpkeepReport_1( uint256  indexed  id);
	event KeeperRegistry1_3_ConfigSet_1( uint32   config_paymentPremiumPPB, uint32   config_flatFeeMicroLink, uint24   config_blockCountPerTurn, uint32   config_checkGasLimit, uint24   config_stalenessSeconds, uint16   config_gasCeilingMultiplier, uint96   config_minUpkeepSpend, uint32   config_maxPerformGas, uint256   config_fallbackGasPrice, uint256   config_fallbackLinkPrice, address   config_transcoder, address   config_registrar);
	event KeeperRegistry1_3_FundsWithdrawn_1( uint256  indexed  id, uint256   amount, address   to);
	event KeeperRegistry2_0_UpkeepCanceled_1( uint256  indexed  id, uint64  indexed  atBlockHeight);
	event AuthorizedOriginReceiverUpgradeable_AuthorizedSendersChanged_1( address[]   senders, address   changedBy);
	event ConfirmedOwner_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event FunctionsBillingRegistry_SubscriptionCreated_1( uint64  indexed  subscriptionId, address   owner);
	event FunctionsOracle_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event KeeperRegistry1_2_ConfigSet_1( uint32   config_paymentPremiumPPB, uint32   config_flatFeeMicroLink, uint24   config_blockCountPerTurn, uint32   config_checkGasLimit, uint24   config_stalenessSeconds, uint16   config_gasCeilingMultiplier, uint96   config_minUpkeepSpend, uint32   config_maxPerformGas, uint256   config_fallbackGasPrice, uint256   config_fallbackLinkPrice, address   config_transcoder, address   config_registrar);
	event KeeperRegistry2_0_UpkeepCheckDataUpdated_1( uint256  indexed  id, bytes   newCheckData);
	event KeeperRegistryBase1_3_UpkeepRegistered_1( uint256  indexed  id, uint32   executeGas, address   admin);
	event KeeperRegistryLogic1_3_UpkeepPerformed_1( uint256  indexed  id, bool  indexed  success, address  indexed  from, uint96   payment, bytes   performData);
	event FunctionsOracle_OracleRequest_1( bytes32  indexed  requestId, address   requestingContract, address   requestInitiator, uint64   subscriptionId, address   subscriptionOwner, bytes   data);
	event FunctionsOracleWithInit_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event KeeperRegistry1_2_UpkeepGasLimitSet_1( uint256  indexed  id, uint96   gasLimit);
	event KeeperRegistry1_3_UpkeepAdminTransferred_1( uint256  indexed  id, address  indexed  from, address  indexed  to);
	event VRFCoordinatorV2_SubscriptionConsumerAdded_1( uint64  indexed  subId, address   consumer);
	event FunctionsBillingRegistry_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event FunctionsBillingRegistry_SubscriptionOwnerTransferred_1( uint64  indexed  subscriptionId, address   from, address   to);
	event FunctionsBillingRegistryWithInit_BillingStart_1( bytes32  indexed  requestId, uint64   commitment_subscriptionId, address   commitment_client, uint32   commitment_gasLimit, uint256   commitment_gasPrice, address   commitment_don, uint96   commitment_donFee, uint96   commitment_registryFee, uint96   commitment_estimatedCost, uint256   commitment_timestamp);
	event VRFV2Wrapper_WrapperFulfillmentFailed_1( uint256  indexed  requestId, address  indexed  consumer);
	event TransparentUpgradeableProxy_BeaconUpgraded_1( address  indexed  beacon);
	event FunctionsOracle_Transmitted_1( bytes32   configDigest, uint32   epoch);
	event KeeperRegistryBase1_3_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event KeeperRegistryBase1_3_UpkeepPaused_1( uint256  indexed  id);
	event KeeperRegistryBase2_0_InsufficientFundsUpkeepReport_1( uint256  indexed  id);
	event KeeperRegistryBase2_0_StaleUpkeepReport_1( uint256  indexed  id);
	event FunctionsOracle_AuthorizedSendersChanged_1( address[]   senders, address   changedBy);
	event KeeperRegistry1_3_PayeeshipTransferRequested_1( address  indexed  keeper, address  indexed  from, address  indexed  to);
	event VerifierProxy_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event FunctionsBillingRegistryWithInit_SubscriptionConsumerAdded_1( uint64  indexed  subscriptionId, address   consumer);
	event KeeperRegistry1_2_OwnerFundsWithdrawn_1( uint96   amount);
	event KeeperRegistryBase1_3_OwnerFundsWithdrawn_1( uint96   amount);
	event OCR2BaseUpgradeable_ConfigSet_1( uint32   previousConfigBlockNumber, bytes32   configDigest, uint64   configCount, address[]   signers, address[]   transmitters, uint8   f, bytes   onchainConfig, uint64   offchainConfigVersion, bytes   offchainConfig);
	event KeeperRegistryBase1_3_UpkeepMigrated_1( uint256  indexed  id, uint256   remainingBalance, address   destination);
	event LogEmitter_Log2_1( uint256  indexed  mfhzxy);
	event VRFCoordinatorV2_SubscriptionOwnerTransferred_1( uint64  indexed  subId, address   from, address   to);
	event FunctionsBillingRegistryWithInit_SubscriptionFunded_1( uint64  indexed  subscriptionId, uint256   oldBalance, uint256   newBalance);
	event FunctionsOracle_AuthorizedSendersActive_1( address   account);
	event FunctionsOracle_UserCallbackRawError_1( bytes32  indexed  requestId, bytes   lowLevelData);
	event FunctionsOracleWithInit_AuthorizedSendersActive_1( address   account);
	event KeeperRegistry1_2_Paused_1( address   account);
	event VRFV2Wrapper_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event Verifier_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event KeeperRegistry2_0_InsufficientFundsUpkeepReport_1( uint256  indexed  id);
	event KeeperRegistry2_0_ReorgedUpkeepReport_1( uint256  indexed  id);
	event KeeperRegistryBase1_3_UpkeepGasLimitSet_1( uint256  indexed  id, uint96   gasLimit);
	event FunctionsBillingRegistryWithInit_ConfigSet_1( uint32   maxGasLimit, uint32   stalenessSeconds, uint256   gasAfterPaymentCalculation, int256   fallbackWeiPerUnitLink, uint32   gasOverhead);
	event KeeperRegistry1_2_UpkeepPerformed_1( uint256  indexed  id, bool  indexed  success, address  indexed  from, uint96   payment, bytes   performData);
	event KeeperRegistry1_3_UpkeepRegistered_1( uint256  indexed  id, uint32   executeGas, address   admin);
	event KeeperRegistry2_0_CancelledUpkeepReport_1( uint256  indexed  id);
	event KeeperRegistry2_0_ConfigSet_1( uint32   previousConfigBlockNumber, bytes32   configDigest, uint64   configCount, address[]   signers, address[]   transmitters, uint8   f, bytes   onchainConfig, uint64   offchainConfigVersion, bytes   offchainConfig);
	event KeeperRegistryBase1_3_UpkeepUnpaused_1( uint256  indexed  id);
	event KeeperRegistryLogic1_3_UpkeepCheckDataUpdated_1( uint256  indexed  id, bytes   newCheckData);
	event OCR2BaseUpgradeable_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event VRFCoordinatorV2_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event FunctionsOracle_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event KeeperRegistryBase2_0_UpkeepPerformed_1( uint256  indexed  id, bool  indexed  success, uint32   checkBlockNumber, uint256   gasUsed, uint256   gasOverhead, uint96   totalPayment);
	event VRFCoordinatorV2_RandomWordsFulfilled_1( uint256  indexed  requestId, uint256   outputSeed, uint96   payment, bool   success);
	event KeeperRegistryLogic1_3_FundsAdded_1( uint256  indexed  id, address  indexed  from, uint96   amount);
	event FunctionsBillingRegistryWithInit_SubscriptionCreated_1( uint64  indexed  subscriptionId, address   owner);
	event KeeperRegistry2_0_UpkeepGasLimitSet_1( uint256  indexed  id, uint96   gasLimit);
	event KeeperRegistry2_0_UpkeepMigrated_1( uint256  indexed  id, uint256   remainingBalance, address   destination);
	event KeeperRegistryBase1_3_KeepersUpdated_1( address[]   keepers, address[]   payees);
	event KeeperRegistryBase1_3_PaymentWithdrawn_1( address  indexed  keeper, uint256  indexed  amount, address  indexed  to, address   payee);
	event KeeperRegistryBase2_0_CancelledUpkeepReport_1( uint256  indexed  id);
	event VRFConsumerV2UpgradeableExample_Initialized_1( uint8   version);
	event VRFCoordinatorV2_SubscriptionConsumerRemoved_1( uint64  indexed  subId, address   consumer);
	event CronUpkeep_CronJobExecuted_1( uint256  indexed  id, uint256   timestamp);
	event FunctionsBillingRegistryWithInit_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event KeeperRegistry1_2_FundsWithdrawn_1( uint256  indexed  id, uint256   amount, address   to);
	event KeeperRegistry1_3_UpkeepPaused_1( uint256  indexed  id);
	event KeeperRegistry2_0_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event VRFCoordinatorV2TestHelper_SubscriptionConsumerRemoved_1( uint64  indexed  subId, address   consumer);
	event ConfirmedOwner_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event CronUpkeep_CronJobUpdated_1( uint256  indexed  id, address   target, bytes   handler);
	event ENSInterface_Transfer_1( bytes32  indexed  node, address   owner);
	event FunctionsOracleWithInit_Initialized_1( uint8   version);
	event KeeperRegistry1_3_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	event OVM_GasPriceOracle_OverheadUpdated_1( uint256   vhectp);
	event FunctionsBillingRegistry_AuthorizedSendersChanged_1( address[]   senders, address   changedBy);
	event FunctionsBillingRegistry_OwnershipTransferRequested_1( address  indexed  from, address  indexed  to);
	event FunctionsBillingRegistryWithInit_Initialized_1( uint8   version);
	event KeeperRegistryBase1_3_UpkeepAdminTransferRequested_1( uint256  indexed  id, address  indexed  from, address  indexed  to);
	event KeeperRegistryBase2_0_ReorgedUpkeepReport_1( uint256  indexed  id);
	event OCR2Abstract_Transmitted_1( bytes32   configDigest, uint32   epoch);
	event OVM_GasPriceOracle_ScalarUpdated_1( uint256   upncvl);
	event VRFCoordinatorV2_SubscriptionOwnerTransferRequested_1( uint64  indexed  subId, address   from, address   to);
	event AuthorizedOriginReceiverUpgradeable_Initialized_1( uint8   version);
	event CronUpkeep_Unpaused_1( address   account);
	event KeeperRegistry1_3_PayeeshipTransferred_1( address  indexed  keeper, address  indexed  from, address  indexed  to);
	event KeeperRegistry2_0_UpkeepRegistered_1( uint256  indexed  id, uint32   executeGas, address   admin);
	event KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferred_1( address  indexed  from, address  indexed  to);
	function emitVRFCoordinatorV2_RandomWordsRequested_1(bytes32 keyHash,uint256 requestId,uint256 preSeed,uint64 subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address sender) public {
		emit VRFCoordinatorV2_RandomWordsRequested_1(keyHash,requestId,preSeed,subId,minimumRequestConfirmations,callbackGasLimit,numWords,sender);
	}
	function emitVRFV2WrapperConsumerExample_WrappedRequestFulfilled_1(uint256 requestId,uint256[] memory randomWords,uint256 payment) public {
		emit VRFV2WrapperConsumerExample_WrappedRequestFulfilled_1(requestId,randomWords,payment);
	}
	function emitFunctionsClient_RequestSent_1(bytes32 id) public {
		emit FunctionsClient_RequestSent_1(id);
	}
	function emitFunctionsClientExample_OwnershipTransferred_1(address from,address to) public {
		emit FunctionsClientExample_OwnershipTransferred_1(from,to);
	}
	function emitFunctionsOracleWithInit_Transmitted_1(bytes32 configDigest,uint32 epoch) public {
		emit FunctionsOracleWithInit_Transmitted_1(configDigest,epoch);
	}
	function emitKeeperRegistryLogic1_3_KeepersUpdated_1(address[] memory keepers,address[] memory payees) public {
		emit KeeperRegistryLogic1_3_KeepersUpdated_1(keepers,payees);
	}
	function emitKeeperRegistryLogic1_3_UpkeepReceived_1(uint256 id,uint256 startingBalance,address importedFrom) public {
		emit KeeperRegistryLogic1_3_UpkeepReceived_1(id,startingBalance,importedFrom);
	}
	function emitVRFV2TransparentUpgradeableProxy_BeaconUpgraded_1(address beacon) public {
		emit VRFV2TransparentUpgradeableProxy_BeaconUpgraded_1(beacon);
	}
	function emitVRFV2TransparentUpgradeableProxy_Upgraded_1(address implementation) public {
		emit VRFV2TransparentUpgradeableProxy_Upgraded_1(implementation);
	}
	function emitBatchVRFCoordinatorV2_ErrorReturned_1(uint256 requestId,string memory reason) public {
		emit BatchVRFCoordinatorV2_ErrorReturned_1(requestId,reason);
	}
	function emitKeeperRegistry1_3_UpkeepMigrated_1(uint256 id,uint256 remainingBalance,address destination) public {
		emit KeeperRegistry1_3_UpkeepMigrated_1(id,remainingBalance,destination);
	}
	function emitKeeperRegistryLogic1_3_UpkeepRegistered_1(uint256 id,uint32 executeGas,address admin) public {
		emit KeeperRegistryLogic1_3_UpkeepRegistered_1(id,executeGas,admin);
	}
	function emitOVM_GasPriceOracle_GasPriceUpdated_1(uint256 pizhiq) public {
		emit OVM_GasPriceOracle_GasPriceUpdated_1(pizhiq);
	}
	function emitVRFCoordinatorV2TestHelper_RandomWordsFulfilled_1(uint256 requestId,uint256 outputSeed,uint96 payment,bool success) public {
		emit VRFCoordinatorV2TestHelper_RandomWordsFulfilled_1(requestId,outputSeed,payment,success);
	}
	function emitConfirmedOwnerUpgradeable_OwnershipTransferred_1(address from,address to) public {
		emit ConfirmedOwnerUpgradeable_OwnershipTransferred_1(from,to);
	}
	function emitKeeperRegistryBase2_0_FundsAdded_1(uint256 id,address from,uint96 amount) public {
		emit KeeperRegistryBase2_0_FundsAdded_1(id,from,amount);
	}
	function emitProxyAdmin_OwnershipTransferred_1(address previousOwner,address newOwner) public {
		emit ProxyAdmin_OwnershipTransferred_1(previousOwner,newOwner);
	}
	function emitKeeperRegistryLogic2_0_ReorgedUpkeepReport_1(uint256 id) public {
		emit KeeperRegistryLogic2_0_ReorgedUpkeepReport_1(id);
	}
	function emitVRFV2Wrapper_OwnershipTransferRequested_1(address from,address to) public {
		emit VRFV2Wrapper_OwnershipTransferRequested_1(from,to);
	}
	function emitKeeperRegistry1_2_OwnershipTransferRequested_1(address from,address to) public {
		emit KeeperRegistry1_2_OwnershipTransferRequested_1(from,to);
	}
	function emitKeeperRegistryBase2_0_UpkeepGasLimitSet_1(uint256 id,uint96 gasLimit) public {
		emit KeeperRegistryBase2_0_UpkeepGasLimitSet_1(id,gasLimit);
	}
	function emitTransparentUpgradeableProxy_AdminChanged_1(address previousAdmin,address newAdmin) public {
		emit TransparentUpgradeableProxy_AdminChanged_1(previousAdmin,newAdmin);
	}
	function emitPausable_Unpaused_1(address account) public {
		emit Pausable_Unpaused_1(account);
	}
	function emitCronUpkeepFactory_OwnershipTransferRequested_1(address from,address to) public {
		emit CronUpkeepFactory_OwnershipTransferRequested_1(from,to);
	}
	function emitFunctionsBillingRegistry_SubscriptionCanceled_1(uint64 subscriptionId,address to,uint256 amount) public {
		emit FunctionsBillingRegistry_SubscriptionCanceled_1(subscriptionId,to,amount);
	}
	function emitKeeperRegistry1_2_OwnershipTransferred_1(address from,address to) public {
		emit KeeperRegistry1_2_OwnershipTransferred_1(from,to);
	}
	function emitKeeperRegistry2_0_UpkeepAdminTransferred_1(uint256 id,address from,address to) public {
		emit KeeperRegistry2_0_UpkeepAdminTransferred_1(id,from,to);
	}
	function emitKeeperRegistryBase2_0_OwnershipTransferred_1(address from,address to) public {
		emit KeeperRegistryBase2_0_OwnershipTransferred_1(from,to);
	}
	function emitConfirmedOwnerUpgradeable_Initialized_1(uint8 version) public {
		emit ConfirmedOwnerUpgradeable_Initialized_1(version);
	}
	function emitFunctionsBillingRegistry_Unpaused_1(address account) public {
		emit FunctionsBillingRegistry_Unpaused_1(account);
	}
	function emitFunctionsBillingRegistryWithInit_SubscriptionOwnerTransferred_1(uint64 subscriptionId,address from,address to) public {
		emit FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferred_1(subscriptionId,from,to);
	}
	function emitVRFCoordinatorV2_SubscriptionCanceled_1(uint64 subId,address to,uint256 amount) public {
		emit VRFCoordinatorV2_SubscriptionCanceled_1(subId,to,amount);
	}
	function emitFunctionsOracleWithInit_AuthorizedSendersChanged_1(address[] memory senders,address changedBy) public {
		emit FunctionsOracleWithInit_AuthorizedSendersChanged_1(senders,changedBy);
	}
	function emitKeeperRegistry1_3_Unpaused_1(address account) public {
		emit KeeperRegistry1_3_Unpaused_1(account);
	}
	function emitKeeperRegistryBase2_0_UpkeepCanceled_1(uint256 id,uint64 atBlockHeight) public {
		emit KeeperRegistryBase2_0_UpkeepCanceled_1(id,atBlockHeight);
	}
	function emitOVM_GasPriceOracle_L1BaseFeeUpdated_1(uint256 qakwnd) public {
		emit OVM_GasPriceOracle_L1BaseFeeUpdated_1(qakwnd);
	}
	function emitVRFConsumerBaseV2Upgradeable_Initialized_1(uint8 version) public {
		emit VRFConsumerBaseV2Upgradeable_Initialized_1(version);
	}
	function emitKeeperRegistry1_2_KeepersUpdated_1(address[] memory keepers,address[] memory payees) public {
		emit KeeperRegistry1_2_KeepersUpdated_1(keepers,payees);
	}
	function emitKeeperRegistryLogic1_3_PaymentWithdrawn_1(address keeper,uint256 amount,address to,address payee) public {
		emit KeeperRegistryLogic1_3_PaymentWithdrawn_1(keeper,amount,to,payee);
	}
	function emitVRFCoordinatorV2TestHelper_SubscriptionCanceled_1(uint64 subId,address to,uint256 amount) public {
		emit VRFCoordinatorV2TestHelper_SubscriptionCanceled_1(subId,to,amount);
	}
	function emitVRFV2WrapperConsumerExample_OwnershipTransferRequested_1(address from,address to) public {
		emit VRFV2WrapperConsumerExample_OwnershipTransferRequested_1(from,to);
	}
	function emitERC1967Upgrade_Upgraded_1(address implementation) public {
		emit ERC1967Upgrade_Upgraded_1(implementation);
	}
	function emitFunctionsBillingRegistryWithInit_OwnershipTransferRequested_1(address from,address to) public {
		emit FunctionsBillingRegistryWithInit_OwnershipTransferRequested_1(from,to);
	}
	function emitFunctionsBillingRegistryWithInit_Unpaused_1(address account) public {
		emit FunctionsBillingRegistryWithInit_Unpaused_1(account);
	}
	function emitKeeperRegistry1_2_FundsAdded_1(uint256 id,address from,uint96 amount) public {
		emit KeeperRegistry1_2_FundsAdded_1(id,from,amount);
	}
	function emitKeeperRegistryBase1_3_UpkeepAdminTransferred_1(uint256 id,address from,address to) public {
		emit KeeperRegistryBase1_3_UpkeepAdminTransferred_1(id,from,to);
	}
	function emitKeeperRegistry2_0_Transmitted_1(bytes32 configDigest,uint32 epoch) public {
		emit KeeperRegistry2_0_Transmitted_1(configDigest,epoch);
	}
	function emitKeeperRegistryBase2_0_UpkeepPaused_1(uint256 id) public {
		emit KeeperRegistryBase2_0_UpkeepPaused_1(id);
	}
	function emitVRFCoordinatorV2_ProvingKeyRegistered_1(bytes32 keyHash,address oracle) public {
		emit VRFCoordinatorV2_ProvingKeyRegistered_1(keyHash,oracle);
	}
	function emitKeeperRegistryLogic1_3_Unpaused_1(address account) public {
		emit KeeperRegistryLogic1_3_Unpaused_1(account);
	}
	function emitKeeperRegistryLogic2_0_OwnershipTransferred_1(address from,address to) public {
		emit KeeperRegistryLogic2_0_OwnershipTransferred_1(from,to);
	}
	function emitKeeperRegistryLogic2_0_UpkeepPerformed_1(uint256 id,bool success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment) public {
		emit KeeperRegistryLogic2_0_UpkeepPerformed_1(id,success,checkBlockNumber,gasUsed,gasOverhead,totalPayment);
	}
	function emitChainlinkClient_ChainlinkFulfilled_1(bytes32 id) public {
		emit ChainlinkClient_ChainlinkFulfilled_1(id);
	}
	function emitERC1967Proxy_BeaconUpgraded_1(address beacon) public {
		emit ERC1967Proxy_BeaconUpgraded_1(beacon);
	}
	function emitFunctionsOracle_OracleResponse_1(bytes32 requestId) public {
		emit FunctionsOracle_OracleResponse_1(requestId);
	}
	function emitFunctionsOracleWithInit_OracleRequest_1(bytes32 requestId,address requestingContract,address requestInitiator,uint64 subscriptionId,address subscriptionOwner,bytes memory data) public {
		emit FunctionsOracleWithInit_OracleRequest_1(requestId,requestingContract,requestInitiator,subscriptionId,subscriptionOwner,data);
	}
	function emitKeeperRegistry1_3_UpkeepPerformed_1(uint256 id,bool success,address from,uint96 payment,bytes memory performData) public {
		emit KeeperRegistry1_3_UpkeepPerformed_1(id,success,from,payment,performData);
	}
	function emitVRFCoordinatorV2TestHelper_ConfigSet_1(uint16 minimumRequestConfirmations,uint32 maxGasLimit,uint32 stalenessSeconds,uint32 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,uint32 feeConfig_fulfillmentFlatFeeLinkPPMTier1,uint32 feeConfig_fulfillmentFlatFeeLinkPPMTier2,uint32 feeConfig_fulfillmentFlatFeeLinkPPMTier3,uint32 feeConfig_fulfillmentFlatFeeLinkPPMTier4,uint32 feeConfig_fulfillmentFlatFeeLinkPPMTier5,uint24 feeConfig_reqsForTier2,uint24 feeConfig_reqsForTier3,uint24 feeConfig_reqsForTier4,uint24 feeConfig_reqsForTier5) public {
		emit VRFCoordinatorV2TestHelper_ConfigSet_1(minimumRequestConfirmations,maxGasLimit,stalenessSeconds,gasAfterPaymentCalculation,fallbackWeiPerUnitLink,feeConfig_fulfillmentFlatFeeLinkPPMTier1,feeConfig_fulfillmentFlatFeeLinkPPMTier2,feeConfig_fulfillmentFlatFeeLinkPPMTier3,feeConfig_fulfillmentFlatFeeLinkPPMTier4,feeConfig_fulfillmentFlatFeeLinkPPMTier5,feeConfig_reqsForTier2,feeConfig_reqsForTier3,feeConfig_reqsForTier4,feeConfig_reqsForTier5);
	}
	function emitVRFCoordinatorV2TestHelper_OwnershipTransferRequested_1(address from,address to) public {
		emit VRFCoordinatorV2TestHelper_OwnershipTransferRequested_1(from,to);
	}
	function emitVerifier_ReportVerified_1(bytes32 feedId,bytes32 reportHash,address requester) public {
		emit Verifier_ReportVerified_1(feedId,reportHash,requester);
	}
	function emitVRFCoordinatorV2TestHelper_SubscriptionOwnerTransferred_1(uint64 subId,address from,address to) public {
		emit VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferred_1(subId,from,to);
	}
	function emitVerifier_ConfigActivated_1(bytes32 feedId,bytes32 configDigest) public {
		emit Verifier_ConfigActivated_1(feedId,configDigest);
	}
	function emitVerifier_ConfigDeactivated_1(bytes32 feedId,bytes32 configDigest) public {
		emit Verifier_ConfigDeactivated_1(feedId,configDigest);
	}
	function emitAuthorizedReceiver_AuthorizedSendersChanged_1(address[] memory senders,address changedBy) public {
		emit AuthorizedReceiver_AuthorizedSendersChanged_1(senders,changedBy);
	}
	function emitKeeperRegistry1_3_Paused_1(address account) public {
		emit KeeperRegistry1_3_Paused_1(account);
	}
	function emitKeeperRegistryBase1_3_FundsWithdrawn_1(uint256 id,uint256 amount,address to) public {
		emit KeeperRegistryBase1_3_FundsWithdrawn_1(id,amount,to);
	}
	function emitKeeperRegistryLogic1_3_OwnershipTransferred_1(address from,address to) public {
		emit KeeperRegistryLogic1_3_OwnershipTransferred_1(from,to);
	}
	function emitVRFCoordinatorV2TestHelper_FundsRecovered_1(address to,uint256 amount) public {
		emit VRFCoordinatorV2TestHelper_FundsRecovered_1(to,amount);
	}
	function emitKeeperRegistryLogic2_0_UpkeepPaused_1(uint256 id) public {
		emit KeeperRegistryLogic2_0_UpkeepPaused_1(id);
	}
	function emitOVM_GasPriceOracle_OwnershipTransferred_1(address previousOwner,address newOwner) public {
		emit OVM_GasPriceOracle_OwnershipTransferred_1(previousOwner,newOwner);
	}
	function emitVRFCoordinatorV2TestHelper_ProvingKeyDeregistered_1(bytes32 keyHash,address oracle) public {
		emit VRFCoordinatorV2TestHelper_ProvingKeyDeregistered_1(keyHash,oracle);
	}
	function emitAggregatorV2V3Interface_NewRound_1(uint256 roundId,address startedBy,uint256 startedAt) public {
		emit AggregatorV2V3Interface_NewRound_1(roundId,startedBy,startedAt);
	}
	function emitFunctionsClient_RequestFulfilled_1(bytes32 id) public {
		emit FunctionsClient_RequestFulfilled_1(id);
	}
	function emitFunctionsOracleWithInit_ConfigSet_1(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
		emit FunctionsOracleWithInit_ConfigSet_1(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
	}
	function emitKeeperRegistry2_0_FundsAdded_1(uint256 id,address from,uint96 amount) public {
		emit KeeperRegistry2_0_FundsAdded_1(id,from,amount);
	}
	function emitKeeperRegistryBase2_0_UpkeepAdminTransferRequested_1(uint256 id,address from,address to) public {
		emit KeeperRegistryBase2_0_UpkeepAdminTransferRequested_1(id,from,to);
	}
	function emitVRFCoordinatorV2TestHelper_RandomWordsRequested_1(bytes32 keyHash,uint256 requestId,uint256 preSeed,uint64 subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address sender) public {
		emit VRFCoordinatorV2TestHelper_RandomWordsRequested_1(keyHash,requestId,preSeed,subId,minimumRequestConfirmations,callbackGasLimit,numWords,sender);
	}
	function emitCronUpkeepFactory_OwnershipTransferred_1(address from,address to) public {
		emit CronUpkeepFactory_OwnershipTransferred_1(from,to);
	}
	function emitKeeperRegistryBase1_3_FundsAdded_1(uint256 id,address from,uint96 amount) public {
		emit KeeperRegistryBase1_3_FundsAdded_1(id,from,amount);
	}
	function emitKeeperRegistryBase2_0_UpkeepAdminTransferred_1(uint256 id,address from,address to) public {
		emit KeeperRegistryBase2_0_UpkeepAdminTransferred_1(id,from,to);
	}
	function emitKeeperRegistryLogic2_0_InsufficientFundsUpkeepReport_1(uint256 id) public {
		emit KeeperRegistryLogic2_0_InsufficientFundsUpkeepReport_1(id);
	}
	function emitKeeperRegistryLogic2_0_PayeeshipTransferRequested_1(address transmitter,address from,address to) public {
		emit KeeperRegistryLogic2_0_PayeeshipTransferRequested_1(transmitter,from,to);
	}
	function emitENSInterface_NewTTL_1(bytes32 node,uint64 ttl) public {
		emit ENSInterface_NewTTL_1(node,ttl);
	}
	function emitFunctionsOracleWithInit_OwnershipTransferRequested_1(address from,address to) public {
		emit FunctionsOracleWithInit_OwnershipTransferRequested_1(from,to);
	}
	function emitKeeperRegistryBase2_0_PayeesUpdated_1(address[] memory transmitters,address[] memory payees) public {
		emit KeeperRegistryBase2_0_PayeesUpdated_1(transmitters,payees);
	}
	function emitKeeperRegistryLogic2_0_PayeesUpdated_1(address[] memory transmitters,address[] memory payees) public {
		emit KeeperRegistryLogic2_0_PayeesUpdated_1(transmitters,payees);
	}
	function emitKeeperRegistryLogic2_0_UpkeepAdminTransferred_1(uint256 id,address from,address to) public {
		emit KeeperRegistryLogic2_0_UpkeepAdminTransferred_1(id,from,to);
	}
	function emitKeeperRegistryLogic2_0_UpkeepReceived_1(uint256 id,uint256 startingBalance,address importedFrom) public {
		emit KeeperRegistryLogic2_0_UpkeepReceived_1(id,startingBalance,importedFrom);
	}
	function emitKeeperRegistryLogic2_0_UpkeepUnpaused_1(uint256 id) public {
		emit KeeperRegistryLogic2_0_UpkeepUnpaused_1(id);
	}
	function emitAggregatorInterface_AnswerUpdated_1(int256 current,uint256 roundId,uint256 updatedAt) public {
		emit AggregatorInterface_AnswerUpdated_1(current,roundId,updatedAt);
	}
	function emitCronUpkeepFactory_NewCronUpkeepCreated_1(address upkeep,address owner) public {
		emit CronUpkeepFactory_NewCronUpkeepCreated_1(upkeep,owner);
	}
	function emitFunctionsOracle_UserCallbackError_1(bytes32 requestId,string memory reason) public {
		emit FunctionsOracle_UserCallbackError_1(requestId,reason);
	}
	function emitKeeperRegistry2_0_StaleUpkeepReport_1(uint256 id) public {
		emit KeeperRegistry2_0_StaleUpkeepReport_1(id);
	}
	function emitKeeperRegistryLogic2_0_UpkeepGasLimitSet_1(uint256 id,uint96 gasLimit) public {
		emit KeeperRegistryLogic2_0_UpkeepGasLimitSet_1(id,gasLimit);
	}
	function emitFunctionsOracleWithInit_InvalidRequestID_1(bytes32 requestId) public {
		emit FunctionsOracleWithInit_InvalidRequestID_1(requestId);
	}
	function emitKeeperRegistryBase1_3_UpkeepPerformed_1(uint256 id,bool success,address from,uint96 payment,bytes memory performData) public {
		emit KeeperRegistryBase1_3_UpkeepPerformed_1(id,success,from,payment,performData);
	}
	function emitKeeperRegistryLogic2_0_PaymentWithdrawn_1(address transmitter,uint256 amount,address to,address payee) public {
		emit KeeperRegistryLogic2_0_PaymentWithdrawn_1(transmitter,amount,to,payee);
	}
	function emitKeeperRegistryLogic2_0_UpkeepOffchainConfigSet_1(uint256 id,bytes memory offchainConfig) public {
		emit KeeperRegistryLogic2_0_UpkeepOffchainConfigSet_1(id,offchainConfig);
	}
	function emitVRFCoordinatorMock_RandomnessRequest_1(address sender,bytes32 keyHash,uint256 seed) public {
		emit VRFCoordinatorMock_RandomnessRequest_1(sender,keyHash,seed);
	}
	function emitVRFCoordinatorV2_SubscriptionCreated_1(uint64 subId,address owner) public {
		emit VRFCoordinatorV2_SubscriptionCreated_1(subId,owner);
	}
	function emitVerifierProxy_VerifierSet_1(bytes32 oldConfigDigest,bytes32 newConfigDigest,address verifierAddress) public {
		emit VerifierProxy_VerifierSet_1(oldConfigDigest,newConfigDigest,verifierAddress);
	}
	function emitAuthorizedOriginReceiverUpgradeable_AuthorizedSendersActive_1(address account) public {
		emit AuthorizedOriginReceiverUpgradeable_AuthorizedSendersActive_1(account);
	}
	function emitFunctionsBillingRegistryWithInit_AuthorizedSendersChanged_1(address[] memory senders,address changedBy) public {
		emit FunctionsBillingRegistryWithInit_AuthorizedSendersChanged_1(senders,changedBy);
	}
	function emitKeeperRegistryBase1_3_UpkeepCanceled_1(uint256 id,uint64 atBlockHeight) public {
		emit KeeperRegistryBase1_3_UpkeepCanceled_1(id,atBlockHeight);
	}
	function emitVRFCoordinatorV2_FundsRecovered_1(address to,uint256 amount) public {
		emit VRFCoordinatorV2_FundsRecovered_1(to,amount);
	}
	function emitVRFCoordinatorV2_ProvingKeyDeregistered_1(bytes32 keyHash,address oracle) public {
		emit VRFCoordinatorV2_ProvingKeyDeregistered_1(keyHash,oracle);
	}
	function emitKeeperRegistry1_2_PayeeshipTransferRequested_1(address keeper,address from,address to) public {
		emit KeeperRegistry1_2_PayeeshipTransferRequested_1(keeper,from,to);
	}
	function emitKeeperRegistry1_2_UpkeepRegistered_1(uint256 id,uint32 executeGas,address admin) public {
		emit KeeperRegistry1_2_UpkeepRegistered_1(id,executeGas,admin);
	}
	function emitKeeperRegistry1_3_UpkeepReceived_1(uint256 id,uint256 startingBalance,address importedFrom) public {
		emit KeeperRegistry1_3_UpkeepReceived_1(id,startingBalance,importedFrom);
	}
	function emitAuthorizedOriginReceiverUpgradeable_AuthorizedSendersDeactive_1(address account) public {
		emit AuthorizedOriginReceiverUpgradeable_AuthorizedSendersDeactive_1(account);
	}
	function emitConfirmedOwnerWithProposal_OwnershipTransferRequested_1(address from,address to) public {
		emit ConfirmedOwnerWithProposal_OwnershipTransferRequested_1(from,to);
	}
	function emitConfirmedOwnerWithProposal_OwnershipTransferred_1(address from,address to) public {
		emit ConfirmedOwnerWithProposal_OwnershipTransferred_1(from,to);
	}
	function emitFunctionsBillingRegistryWithInit_BillingEnd_1(bytes32 requestId,uint64 subscriptionId,uint96 signerPayment,uint96 transmitterPayment,uint96 totalCost,bool success) public {
		emit FunctionsBillingRegistryWithInit_BillingEnd_1(requestId,subscriptionId,signerPayment,transmitterPayment,totalCost,success);
	}
	function emitFunctionsBillingRegistryWithInit_RequestTimedOut_1(bytes32 requestId) public {
		emit FunctionsBillingRegistryWithInit_RequestTimedOut_1(requestId);
	}
	function emitKeeperRegistry2_0_PayeesUpdated_1(address[] memory transmitters,address[] memory payees) public {
		emit KeeperRegistry2_0_PayeesUpdated_1(transmitters,payees);
	}
	function emitKeeperRegistry2_0_PayeeshipTransferred_1(address transmitter,address from,address to) public {
		emit KeeperRegistry2_0_PayeeshipTransferred_1(transmitter,from,to);
	}
	function emitKeeperRegistryLogic1_3_UpkeepUnpaused_1(uint256 id) public {
		emit KeeperRegistryLogic1_3_UpkeepUnpaused_1(id);
	}
	function emitKeeperRegistryLogic2_0_FundsWithdrawn_1(uint256 id,uint256 amount,address to) public {
		emit KeeperRegistryLogic2_0_FundsWithdrawn_1(id,amount,to);
	}
	function emitENSInterface_NewResolver_1(bytes32 node,address resolver) public {
		emit ENSInterface_NewResolver_1(node,resolver);
	}
	function emitFunctionsOracleWithInit_AuthorizedSendersDeactive_1(address account) public {
		emit FunctionsOracleWithInit_AuthorizedSendersDeactive_1(account);
	}
	function emitKeeperRegistry1_3_FundsAdded_1(uint256 id,address from,uint96 amount) public {
		emit KeeperRegistry1_3_FundsAdded_1(id,from,amount);
	}
	function emitVerifier_OwnershipTransferRequested_1(address from,address to) public {
		emit Verifier_OwnershipTransferRequested_1(from,to);
	}
	function emitOCR2Abstract_ConfigSet_1(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
		emit OCR2Abstract_ConfigSet_1(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
	}
	function emitENSInterface_NewOwner_1(bytes32 node,bytes32 label,address owner) public {
		emit ENSInterface_NewOwner_1(node,label,owner);
	}
	function emitFunctionsBillingRegistryWithInit_SubscriptionOwnerTransferRequested_1(uint64 subscriptionId,address from,address to) public {
		emit FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferRequested_1(subscriptionId,from,to);
	}
	function emitFunctionsClientExample_RequestFulfilled_1(bytes32 id) public {
		emit FunctionsClientExample_RequestFulfilled_1(id);
	}
	function emitKeeperRegistry2_0_Paused_1(address account) public {
		emit KeeperRegistry2_0_Paused_1(account);
	}
	function emitKeeperRegistryLogic1_3_UpkeepGasLimitSet_1(uint256 id,uint96 gasLimit) public {
		emit KeeperRegistryLogic1_3_UpkeepGasLimitSet_1(id,gasLimit);
	}
	function emitVRFCoordinatorV2TestHelper_OwnershipTransferred_1(address from,address to) public {
		emit VRFCoordinatorV2TestHelper_OwnershipTransferred_1(from,to);
	}
	function emitVRFV2TransparentUpgradeableProxy_AdminChanged_1(address previousAdmin,address newAdmin) public {
		emit VRFV2TransparentUpgradeableProxy_AdminChanged_1(previousAdmin,newAdmin);
	}
	function emitFunctionsBillingRegistry_FundsRecovered_1(address to,uint256 amount) public {
		emit FunctionsBillingRegistry_FundsRecovered_1(to,amount);
	}
	function emitKeeperRegistry1_3_UpkeepCheckDataUpdated_1(uint256 id,bytes memory newCheckData) public {
		emit KeeperRegistry1_3_UpkeepCheckDataUpdated_1(id,newCheckData);
	}
	function emitKeeperRegistry2_0_UpkeepPaused_1(uint256 id) public {
		emit KeeperRegistry2_0_UpkeepPaused_1(id);
	}
	function emitKeeperRegistryBase2_0_UpkeepOffchainConfigSet_1(uint256 id,bytes memory offchainConfig) public {
		emit KeeperRegistryBase2_0_UpkeepOffchainConfigSet_1(id,offchainConfig);
	}
	function emitOVM_GasPriceOracle_DecimalsUpdated_1(uint256 nlwwiv) public {
		emit OVM_GasPriceOracle_DecimalsUpdated_1(nlwwiv);
	}
	function emitVerifierProxy_VerifierUnset_1(bytes32 configDigest,address verifierAddress) public {
		emit VerifierProxy_VerifierUnset_1(configDigest,verifierAddress);
	}
	function emitFunctionsBillingRegistry_Paused_1(address account) public {
		emit FunctionsBillingRegistry_Paused_1(account);
	}
	function emitInitializable_Initialized_1(uint8 version) public {
		emit Initializable_Initialized_1(version);
	}
	function emitKeeperRegistryLogic1_3_FundsWithdrawn_1(uint256 id,uint256 amount,address to) public {
		emit KeeperRegistryLogic1_3_FundsWithdrawn_1(id,amount,to);
	}
	function emitKeeperRegistryLogic1_3_Paused_1(address account) public {
		emit KeeperRegistryLogic1_3_Paused_1(account);
	}
	function emitKeeperRegistryLogic2_0_PayeeshipTransferred_1(address transmitter,address from,address to) public {
		emit KeeperRegistryLogic2_0_PayeeshipTransferred_1(transmitter,from,to);
	}
	function emitKeeperRegistry1_3_UpkeepGasLimitSet_1(uint256 id,uint96 gasLimit) public {
		emit KeeperRegistry1_3_UpkeepGasLimitSet_1(id,gasLimit);
	}
	function emitFunctionsOracle_AuthorizedSendersDeactive_1(address account) public {
		emit FunctionsOracle_AuthorizedSendersDeactive_1(account);
	}
	function emitVRFCoordinatorV2TestHelper_ProvingKeyRegistered_1(bytes32 keyHash,address oracle) public {
		emit VRFCoordinatorV2TestHelper_ProvingKeyRegistered_1(keyHash,oracle);
	}
	function emitVerifier_ConfigSet_1(bytes32 feedId,uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,bytes32[] memory offchainTransmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
		emit Verifier_ConfigSet_1(feedId,previousConfigBlockNumber,configDigest,configCount,signers,offchainTransmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
	}
	function emitKeeperRegistryBase2_0_PaymentWithdrawn_1(address transmitter,uint256 amount,address to,address payee) public {
		emit KeeperRegistryBase2_0_PaymentWithdrawn_1(transmitter,amount,to,payee);
	}
	function emitKeeperRegistryLogic2_0_UpkeepCanceled_1(uint256 id,uint64 atBlockHeight) public {
		emit KeeperRegistryLogic2_0_UpkeepCanceled_1(id,atBlockHeight);
	}
	function emitOCR2BaseUpgradeable_OwnershipTransferRequested_1(address from,address to) public {
		emit OCR2BaseUpgradeable_OwnershipTransferRequested_1(from,to);
	}
	function emitAggregatorInterface_NewRound_1(uint256 roundId,address startedBy,uint256 startedAt) public {
		emit AggregatorInterface_NewRound_1(roundId,startedBy,startedAt);
	}
	function emitAggregatorV2V3Interface_AnswerUpdated_1(int256 current,uint256 roundId,uint256 updatedAt) public {
		emit AggregatorV2V3Interface_AnswerUpdated_1(current,roundId,updatedAt);
	}
	function emitCronUpkeep_OwnershipTransferred_1(address from,address to) public {
		emit CronUpkeep_OwnershipTransferred_1(from,to);
	}
	function emitFunctionsBillingRegistry_RequestTimedOut_1(bytes32 requestId) public {
		emit FunctionsBillingRegistry_RequestTimedOut_1(requestId);
	}
	function emitFunctionsBillingRegistry_SubscriptionFunded_1(uint64 subscriptionId,uint256 oldBalance,uint256 newBalance) public {
		emit FunctionsBillingRegistry_SubscriptionFunded_1(subscriptionId,oldBalance,newBalance);
	}
	function emitOwnable_OwnershipTransferred_1(address previousOwner,address newOwner) public {
		emit Ownable_OwnershipTransferred_1(previousOwner,newOwner);
	}
	function emitVRFCoordinatorV2_OwnershipTransferred_1(address from,address to) public {
		emit VRFCoordinatorV2_OwnershipTransferred_1(from,to);
	}
	function emitVRFV2WrapperConsumerExample_OwnershipTransferred_1(address from,address to) public {
		emit VRFV2WrapperConsumerExample_OwnershipTransferred_1(from,to);
	}
	function emitKeeperRegistryLogic1_3_OwnerFundsWithdrawn_1(uint96 amount) public {
		emit KeeperRegistryLogic1_3_OwnerFundsWithdrawn_1(amount);
	}
	function emitKeeperRegistryLogic2_0_StaleUpkeepReport_1(uint256 id) public {
		emit KeeperRegistryLogic2_0_StaleUpkeepReport_1(id);
	}
	function emitVRFCoordinatorV2TestHelper_SubscriptionOwnerTransferRequested_1(uint64 subId,address from,address to) public {
		emit VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferRequested_1(subId,from,to);
	}
	function emitKeeperRegistry1_3_PaymentWithdrawn_1(address keeper,uint256 amount,address to,address payee) public {
		emit KeeperRegistry1_3_PaymentWithdrawn_1(keeper,amount,to,payee);
	}
	function emitKeeperRegistry2_0_UpkeepAdminTransferRequested_1(uint256 id,address from,address to) public {
		emit KeeperRegistry2_0_UpkeepAdminTransferRequested_1(id,from,to);
	}
	function emitKeeperRegistryBase1_3_UpkeepCheckDataUpdated_1(uint256 id,bytes memory newCheckData) public {
		emit KeeperRegistryBase1_3_UpkeepCheckDataUpdated_1(id,newCheckData);
	}
	function emitKeeperRegistryLogic2_0_UpkeepMigrated_1(uint256 id,uint256 remainingBalance,address destination) public {
		emit KeeperRegistryLogic2_0_UpkeepMigrated_1(id,remainingBalance,destination);
	}
	function emitFunctionsBillingRegistry_SubscriptionConsumerRemoved_1(uint64 subscriptionId,address consumer) public {
		emit FunctionsBillingRegistry_SubscriptionConsumerRemoved_1(subscriptionId,consumer);
	}
	function emitKeeperRegistryBase2_0_UpkeepCheckDataUpdated_1(uint256 id,bytes memory newCheckData) public {
		emit KeeperRegistryBase2_0_UpkeepCheckDataUpdated_1(id,newCheckData);
	}
	function emitKeeperRegistryLogic1_3_OwnershipTransferRequested_1(address from,address to) public {
		emit KeeperRegistryLogic1_3_OwnershipTransferRequested_1(from,to);
	}
	function emitKeeperRegistryLogic2_0_UpkeepRegistered_1(uint256 id,uint32 executeGas,address admin) public {
		emit KeeperRegistryLogic2_0_UpkeepRegistered_1(id,executeGas,admin);
	}
	function emitCronUpkeep_OwnershipTransferRequested_1(address from,address to) public {
		emit CronUpkeep_OwnershipTransferRequested_1(from,to);
	}
	function emitFunctionsBillingRegistry_SubscriptionOwnerTransferRequested_1(uint64 subscriptionId,address from,address to) public {
		emit FunctionsBillingRegistry_SubscriptionOwnerTransferRequested_1(subscriptionId,from,to);
	}
	function emitKeeperRegistry1_3_KeepersUpdated_1(address[] memory keepers,address[] memory payees) public {
		emit KeeperRegistry1_3_KeepersUpdated_1(keepers,payees);
	}
	function emitKeeperRegistryBase1_3_Unpaused_1(address account) public {
		emit KeeperRegistryBase1_3_Unpaused_1(account);
	}
	function emitKeeperRegistryBase1_3_UpkeepReceived_1(uint256 id,uint256 startingBalance,address importedFrom) public {
		emit KeeperRegistryBase1_3_UpkeepReceived_1(id,startingBalance,importedFrom);
	}
	function emitCronUpkeep_CronJobDeleted_1(uint256 id) public {
		emit CronUpkeep_CronJobDeleted_1(id);
	}
	function emitERC1967Upgrade_AdminChanged_1(address previousAdmin,address newAdmin) public {
		emit ERC1967Upgrade_AdminChanged_1(previousAdmin,newAdmin);
	}
	function emitKeeperRegistryBase2_0_PayeeshipTransferRequested_1(address transmitter,address from,address to) public {
		emit KeeperRegistryBase2_0_PayeeshipTransferRequested_1(transmitter,from,to);
	}
	function emitKeeperRegistryLogic2_0_UpkeepCheckDataUpdated_1(uint256 id,bytes memory newCheckData) public {
		emit KeeperRegistryLogic2_0_UpkeepCheckDataUpdated_1(id,newCheckData);
	}
	function emitKeeperRegistryBase2_0_UpkeepReceived_1(uint256 id,uint256 startingBalance,address importedFrom) public {
		emit KeeperRegistryBase2_0_UpkeepReceived_1(id,startingBalance,importedFrom);
	}
	function emitKeeperRegistryLogic2_0_Paused_1(address account) public {
		emit KeeperRegistryLogic2_0_Paused_1(account);
	}
	function emitVRFCoordinatorV2TestHelper_SubscriptionFunded_1(uint64 subId,uint256 oldBalance,uint256 newBalance) public {
		emit VRFCoordinatorV2TestHelper_SubscriptionFunded_1(subId,oldBalance,newBalance);
	}
	function emitFunctionsBillingRegistry_BillingEnd_1(bytes32 requestId,uint64 subscriptionId,uint96 signerPayment,uint96 transmitterPayment,uint96 totalCost,bool success) public {
		emit FunctionsBillingRegistry_BillingEnd_1(requestId,subscriptionId,signerPayment,transmitterPayment,totalCost,success);
	}
	function emitFunctionsBillingRegistryWithInit_Paused_1(address account) public {
		emit FunctionsBillingRegistryWithInit_Paused_1(account);
	}
	function emitFunctionsOracleWithInit_OracleResponse_1(bytes32 requestId) public {
		emit FunctionsOracleWithInit_OracleResponse_1(requestId);
	}
	function emitFunctionsOracleWithInit_UserCallbackRawError_1(bytes32 requestId,bytes memory lowLevelData) public {
		emit FunctionsOracleWithInit_UserCallbackRawError_1(requestId,lowLevelData);
	}
	function emitKeeperRegistryBase1_3_Paused_1(address account) public {
		emit KeeperRegistryBase1_3_Paused_1(account);
	}
	function emitTransparentUpgradeableProxy_Upgraded_1(address implementation) public {
		emit TransparentUpgradeableProxy_Upgraded_1(implementation);
	}
	function emitFunctionsBillingRegistry_BillingStart_1(bytes32 requestId,uint64 commitment_subscriptionId,address commitment_client,uint32 commitment_gasLimit,uint256 commitment_gasPrice,address commitment_don,uint96 commitment_donFee,uint96 commitment_registryFee,uint96 commitment_estimatedCost,uint256 commitment_timestamp) public {
		emit FunctionsBillingRegistry_BillingStart_1(requestId,commitment_subscriptionId,commitment_client,commitment_gasLimit,commitment_gasPrice,commitment_don,commitment_donFee,commitment_registryFee,commitment_estimatedCost,commitment_timestamp);
	}
	function emitKeeperRegistryBase1_3_PayeeshipTransferred_1(address keeper,address from,address to) public {
		emit KeeperRegistryBase1_3_PayeeshipTransferred_1(keeper,from,to);
	}
	function emitKeeperRegistryLogic2_0_FundsAdded_1(uint256 id,address from,uint96 amount) public {
		emit KeeperRegistryLogic2_0_FundsAdded_1(id,from,amount);
	}
	function emitPausable_Paused_1(address account) public {
		emit Pausable_Paused_1(account);
	}
	function emitPausableUpgradeable_Paused_1(address account) public {
		emit PausableUpgradeable_Paused_1(account);
	}
	function emitKeeperRegistry2_0_PaymentWithdrawn_1(address transmitter,uint256 amount,address to,address payee) public {
		emit KeeperRegistry2_0_PaymentWithdrawn_1(transmitter,amount,to,payee);
	}
	function emitKeeperRegistryBase2_0_PayeeshipTransferred_1(address transmitter,address from,address to) public {
		emit KeeperRegistryBase2_0_PayeeshipTransferred_1(transmitter,from,to);
	}
	function emitVRFLoadTestExternalSubOwner_OwnershipTransferRequested_1(address from,address to) public {
		emit VRFLoadTestExternalSubOwner_OwnershipTransferRequested_1(from,to);
	}
	function emitOCR2BaseUpgradeable_Transmitted_1(bytes32 configDigest,uint32 epoch) public {
		emit OCR2BaseUpgradeable_Transmitted_1(configDigest,epoch);
	}
	function emitVRFLoadTestExternalSubOwner_OwnershipTransferred_1(address from,address to) public {
		emit VRFLoadTestExternalSubOwner_OwnershipTransferred_1(from,to);
	}
	function emitContextUpgradeable_Initialized_1(uint8 version) public {
		emit ContextUpgradeable_Initialized_1(version);
	}
	function emitFunctionsBillingRegistryWithInit_FundsRecovered_1(address to,uint256 amount) public {
		emit FunctionsBillingRegistryWithInit_FundsRecovered_1(to,amount);
	}
	function emitFunctionsBillingRegistryWithInit_SubscriptionCanceled_1(uint64 subscriptionId,address to,uint256 amount) public {
		emit FunctionsBillingRegistryWithInit_SubscriptionCanceled_1(subscriptionId,to,amount);
	}
	function emitFunctionsClientExample_OwnershipTransferRequested_1(address from,address to) public {
		emit FunctionsClientExample_OwnershipTransferRequested_1(from,to);
	}
	function emitKeeperRegistryLogic1_3_UpkeepAdminTransferred_1(uint256 id,address from,address to) public {
		emit KeeperRegistryLogic1_3_UpkeepAdminTransferred_1(id,from,to);
	}
	function emitOCR2BaseUpgradeable_Initialized_1(uint8 version) public {
		emit OCR2BaseUpgradeable_Initialized_1(version);
	}
	function emitPausableUpgradeable_Initialized_1(uint8 version) public {
		emit PausableUpgradeable_Initialized_1(version);
	}
	function emitVRFCoordinatorV2_ConfigSet_1(uint16 minimumRequestConfirmations,uint32 maxGasLimit,uint32 stalenessSeconds,uint32 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,uint32 feeConfig_fulfillmentFlatFeeLinkPPMTier1,uint32 feeConfig_fulfillmentFlatFeeLinkPPMTier2,uint32 feeConfig_fulfillmentFlatFeeLinkPPMTier3,uint32 feeConfig_fulfillmentFlatFeeLinkPPMTier4,uint32 feeConfig_fulfillmentFlatFeeLinkPPMTier5,uint24 feeConfig_reqsForTier2,uint24 feeConfig_reqsForTier3,uint24 feeConfig_reqsForTier4,uint24 feeConfig_reqsForTier5) public {
		emit VRFCoordinatorV2_ConfigSet_1(minimumRequestConfirmations,maxGasLimit,stalenessSeconds,gasAfterPaymentCalculation,fallbackWeiPerUnitLink,feeConfig_fulfillmentFlatFeeLinkPPMTier1,feeConfig_fulfillmentFlatFeeLinkPPMTier2,feeConfig_fulfillmentFlatFeeLinkPPMTier3,feeConfig_fulfillmentFlatFeeLinkPPMTier4,feeConfig_fulfillmentFlatFeeLinkPPMTier5,feeConfig_reqsForTier2,feeConfig_reqsForTier3,feeConfig_reqsForTier4,feeConfig_reqsForTier5);
	}
	function emitERC1967Proxy_Upgraded_1(address implementation) public {
		emit ERC1967Proxy_Upgraded_1(implementation);
	}
	function emitFunctionsBillingRegistry_Initialized_1(uint8 version) public {
		emit FunctionsBillingRegistry_Initialized_1(version);
	}
	function emitKeeperRegistry1_2_PaymentWithdrawn_1(address keeper,uint256 amount,address to,address payee) public {
		emit KeeperRegistry1_2_PaymentWithdrawn_1(keeper,amount,to,payee);
	}
	function emitKeeperRegistry2_0_PayeeshipTransferRequested_1(address transmitter,address from,address to) public {
		emit KeeperRegistry2_0_PayeeshipTransferRequested_1(transmitter,from,to);
	}
	function emitKeeperRegistryLogic1_3_UpkeepMigrated_1(uint256 id,uint256 remainingBalance,address destination) public {
		emit KeeperRegistryLogic1_3_UpkeepMigrated_1(id,remainingBalance,destination);
	}
	function emitKeeperRegistryLogic1_3_PayeeshipTransferRequested_1(address keeper,address from,address to) public {
		emit KeeperRegistryLogic1_3_PayeeshipTransferRequested_1(keeper,from,to);
	}
	function emitKeeperRegistryLogic1_3_PayeeshipTransferred_1(address keeper,address from,address to) public {
		emit KeeperRegistryLogic1_3_PayeeshipTransferred_1(keeper,from,to);
	}
	function emitKeeperRegistryLogic2_0_Unpaused_1(address account) public {
		emit KeeperRegistryLogic2_0_Unpaused_1(account);
	}
	function emitKeeperRegistryBase2_0_OwnershipTransferRequested_1(address from,address to) public {
		emit KeeperRegistryBase2_0_OwnershipTransferRequested_1(from,to);
	}
	function emitKeeperRegistryBase2_0_UpkeepMigrated_1(uint256 id,uint256 remainingBalance,address destination) public {
		emit KeeperRegistryBase2_0_UpkeepMigrated_1(id,remainingBalance,destination);
	}
	function emitKeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferRequested_1(address from,address to) public {
		emit KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferRequested_1(from,to);
	}
	function emitFunctionsBillingRegistry_SubscriptionConsumerAdded_1(uint64 subscriptionId,address consumer) public {
		emit FunctionsBillingRegistry_SubscriptionConsumerAdded_1(subscriptionId,consumer);
	}
	function emitFunctionsOracle_InvalidRequestID_1(bytes32 requestId) public {
		emit FunctionsOracle_InvalidRequestID_1(requestId);
	}
	function emitFunctionsOracleWithInit_UserCallbackError_1(bytes32 requestId,string memory reason) public {
		emit FunctionsOracleWithInit_UserCallbackError_1(requestId,reason);
	}
	function emitKeeperRegistry2_0_OwnershipTransferred_1(address from,address to) public {
		emit KeeperRegistry2_0_OwnershipTransferred_1(from,to);
	}
	function emitKeeperRegistry2_0_UpkeepUnpaused_1(uint256 id) public {
		emit KeeperRegistry2_0_UpkeepUnpaused_1(id);
	}
	function emitKeeperRegistryLogic1_3_UpkeepAdminTransferRequested_1(uint256 id,address from,address to) public {
		emit KeeperRegistryLogic1_3_UpkeepAdminTransferRequested_1(id,from,to);
	}
	function emitLogEmitter_Log3_1(string memory nelhzo) public {
		emit LogEmitter_Log3_1(nelhzo);
	}
	function emitVRFCoordinatorV2_SubscriptionFunded_1(uint64 subId,uint256 oldBalance,uint256 newBalance) public {
		emit VRFCoordinatorV2_SubscriptionFunded_1(subId,oldBalance,newBalance);
	}
	function emitVRFCoordinatorV2TestHelper_SubscriptionConsumerAdded_1(uint64 subId,address consumer) public {
		emit VRFCoordinatorV2TestHelper_SubscriptionConsumerAdded_1(subId,consumer);
	}
	function emitERC1967Upgrade_BeaconUpgraded_1(address beacon) public {
		emit ERC1967Upgrade_BeaconUpgraded_1(beacon);
	}
	function emitKeeperRegistry2_0_FundsWithdrawn_1(uint256 id,uint256 amount,address to) public {
		emit KeeperRegistry2_0_FundsWithdrawn_1(id,amount,to);
	}
	function emitKeeperRegistryBase2_0_OwnerFundsWithdrawn_1(uint96 amount) public {
		emit KeeperRegistryBase2_0_OwnerFundsWithdrawn_1(amount);
	}
	function emitChainlinkClient_ChainlinkRequested_1(bytes32 id) public {
		emit ChainlinkClient_ChainlinkRequested_1(id);
	}
	function emitFunctionsOracle_Initialized_1(uint8 version) public {
		emit FunctionsOracle_Initialized_1(version);
	}
	function emitKeeperRegistry1_2_PayeeshipTransferred_1(address keeper,address from,address to) public {
		emit KeeperRegistry1_2_PayeeshipTransferred_1(keeper,from,to);
	}
	function emitKeeperRegistryBase2_0_UpkeepUnpaused_1(uint256 id) public {
		emit KeeperRegistryBase2_0_UpkeepUnpaused_1(id);
	}
	function emitKeeperRegistryLogic1_3_UpkeepPaused_1(uint256 id) public {
		emit KeeperRegistryLogic1_3_UpkeepPaused_1(id);
	}
	function emitFunctionsClientExample_RequestSent_1(bytes32 id) public {
		emit FunctionsClientExample_RequestSent_1(id);
	}
	function emitKeeperRegistry2_0_UpkeepPerformed_1(uint256 id,bool success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment) public {
		emit KeeperRegistry2_0_UpkeepPerformed_1(id,success,checkBlockNumber,gasUsed,gasOverhead,totalPayment);
	}
	function emitKeeperRegistryBase2_0_FundsWithdrawn_1(uint256 id,uint256 amount,address to) public {
		emit KeeperRegistryBase2_0_FundsWithdrawn_1(id,amount,to);
	}
	function emitLogEmitter_Log1_1(uint256 bevkse) public {
		emit LogEmitter_Log1_1(bevkse);
	}
	function emitVRFCoordinatorV2TestHelper_SubscriptionCreated_1(uint64 subId,address owner) public {
		emit VRFCoordinatorV2TestHelper_SubscriptionCreated_1(subId,owner);
	}
	function emitChainlinkClient_ChainlinkCancelled_1(bytes32 id) public {
		emit ChainlinkClient_ChainlinkCancelled_1(id);
	}
	function emitKeeperRegistry1_2_UpkeepCanceled_1(uint256 id,uint64 atBlockHeight) public {
		emit KeeperRegistry1_2_UpkeepCanceled_1(id,atBlockHeight);
	}
	function emitKeeperRegistry2_0_OwnerFundsWithdrawn_1(uint96 amount) public {
		emit KeeperRegistry2_0_OwnerFundsWithdrawn_1(amount);
	}
	function emitKeeperRegistry1_3_UpkeepCanceled_1(uint256 id,uint64 atBlockHeight) public {
		emit KeeperRegistry1_3_UpkeepCanceled_1(id,atBlockHeight);
	}
	function emitKeeperRegistry2_0_UpkeepOffchainConfigSet_1(uint256 id,bytes memory offchainConfig) public {
		emit KeeperRegistry2_0_UpkeepOffchainConfigSet_1(id,offchainConfig);
	}
	function emitKeeperRegistryLogic2_0_OwnershipTransferRequested_1(address from,address to) public {
		emit KeeperRegistryLogic2_0_OwnershipTransferRequested_1(from,to);
	}
	function emitConfirmedOwnerUpgradeable_OwnershipTransferRequested_1(address from,address to) public {
		emit ConfirmedOwnerUpgradeable_OwnershipTransferRequested_1(from,to);
	}
	function emitERC1967Proxy_AdminChanged_1(address previousAdmin,address newAdmin) public {
		emit ERC1967Proxy_AdminChanged_1(previousAdmin,newAdmin);
	}
	function emitFunctionsBillingRegistryWithInit_SubscriptionConsumerRemoved_1(uint64 subscriptionId,address consumer) public {
		emit FunctionsBillingRegistryWithInit_SubscriptionConsumerRemoved_1(subscriptionId,consumer);
	}
	function emitKeeperRegistry1_2_UpkeepMigrated_1(uint256 id,uint256 remainingBalance,address destination) public {
		emit KeeperRegistry1_2_UpkeepMigrated_1(id,remainingBalance,destination);
	}
	function emitKeeperRegistry1_2_UpkeepReceived_1(uint256 id,uint256 startingBalance,address importedFrom) public {
		emit KeeperRegistry1_2_UpkeepReceived_1(id,startingBalance,importedFrom);
	}
	function emitKeeperRegistryBase1_3_ConfigSet_1(uint32 config_paymentPremiumPPB,uint32 config_flatFeeMicroLink,uint24 config_blockCountPerTurn,uint32 config_checkGasLimit,uint24 config_stalenessSeconds,uint16 config_gasCeilingMultiplier,uint96 config_minUpkeepSpend,uint32 config_maxPerformGas,uint256 config_fallbackGasPrice,uint256 config_fallbackLinkPrice,address config_transcoder,address config_registrar) public {
		emit KeeperRegistryBase1_3_ConfigSet_1(config_paymentPremiumPPB,config_flatFeeMicroLink,config_blockCountPerTurn,config_checkGasLimit,config_stalenessSeconds,config_gasCeilingMultiplier,config_minUpkeepSpend,config_maxPerformGas,config_fallbackGasPrice,config_fallbackLinkPrice,config_transcoder,config_registrar);
	}
	function emitKeeperRegistryLogic2_0_UpkeepAdminTransferRequested_1(uint256 id,address from,address to) public {
		emit KeeperRegistryLogic2_0_UpkeepAdminTransferRequested_1(id,from,to);
	}
	function emitPausableUpgradeable_Unpaused_1(address account) public {
		emit PausableUpgradeable_Unpaused_1(account);
	}
	function emitVRFV2ProxyAdmin_OwnershipTransferred_1(address previousOwner,address newOwner) public {
		emit VRFV2ProxyAdmin_OwnershipTransferred_1(previousOwner,newOwner);
	}
	function emitVRFV2WrapperConsumerExample_WrapperRequestMade_1(uint256 requestId,uint256 paid) public {
		emit VRFV2WrapperConsumerExample_WrapperRequestMade_1(requestId,paid);
	}
	function emitFunctionsBillingRegistry_ConfigSet_1(uint32 maxGasLimit,uint32 stalenessSeconds,uint256 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,uint32 gasOverhead) public {
		emit FunctionsBillingRegistry_ConfigSet_1(maxGasLimit,stalenessSeconds,gasAfterPaymentCalculation,fallbackWeiPerUnitLink,gasOverhead);
	}
	function emitFunctionsOracle_ConfigSet_1(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
		emit FunctionsOracle_ConfigSet_1(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
	}
	function emitKeeperRegistryBase2_0_UpkeepRegistered_1(uint256 id,uint32 executeGas,address admin) public {
		emit KeeperRegistryBase2_0_UpkeepRegistered_1(id,executeGas,admin);
	}
	function emitKeeperRegistryLogic1_3_ConfigSet_1(uint32 config_paymentPremiumPPB,uint32 config_flatFeeMicroLink,uint24 config_blockCountPerTurn,uint32 config_checkGasLimit,uint24 config_stalenessSeconds,uint16 config_gasCeilingMultiplier,uint96 config_minUpkeepSpend,uint32 config_maxPerformGas,uint256 config_fallbackGasPrice,uint256 config_fallbackLinkPrice,address config_transcoder,address config_registrar) public {
		emit KeeperRegistryLogic1_3_ConfigSet_1(config_paymentPremiumPPB,config_flatFeeMicroLink,config_blockCountPerTurn,config_checkGasLimit,config_stalenessSeconds,config_gasCeilingMultiplier,config_minUpkeepSpend,config_maxPerformGas,config_fallbackGasPrice,config_fallbackLinkPrice,config_transcoder,config_registrar);
	}
	function emitKeeperRegistryLogic1_3_UpkeepCanceled_1(uint256 id,uint64 atBlockHeight) public {
		emit KeeperRegistryLogic1_3_UpkeepCanceled_1(id,atBlockHeight);
	}
	function emitVerifierProxy_AccessControllerSet_1(address oldAccessController,address newAccessController) public {
		emit VerifierProxy_AccessControllerSet_1(oldAccessController,newAccessController);
	}
	function emitCronUpkeep_CronJobCreated_1(uint256 id,address target,bytes memory handler) public {
		emit CronUpkeep_CronJobCreated_1(id,target,handler);
	}
	function emitCronUpkeep_Paused_1(address account) public {
		emit CronUpkeep_Paused_1(account);
	}
	function emitKeeperRegistry1_2_Unpaused_1(address account) public {
		emit KeeperRegistry1_2_Unpaused_1(account);
	}
	function emitKeeperRegistry1_3_OwnershipTransferRequested_1(address from,address to) public {
		emit KeeperRegistry1_3_OwnershipTransferRequested_1(from,to);
	}
	function emitKeeperRegistryBase1_3_OwnershipTransferred_1(address from,address to) public {
		emit KeeperRegistryBase1_3_OwnershipTransferred_1(from,to);
	}
	function emitBatchVRFCoordinatorV2_RawErrorReturned_1(uint256 requestId,bytes memory lowLevelData) public {
		emit BatchVRFCoordinatorV2_RawErrorReturned_1(requestId,lowLevelData);
	}
	function emitKeeperRegistry2_0_UpkeepReceived_1(uint256 id,uint256 startingBalance,address importedFrom) public {
		emit KeeperRegistry2_0_UpkeepReceived_1(id,startingBalance,importedFrom);
	}
	function emitKeeperRegistryBase2_0_Unpaused_1(address account) public {
		emit KeeperRegistryBase2_0_Unpaused_1(account);
	}
	function emitKeeperRegistryLogic2_0_OwnerFundsWithdrawn_1(uint96 amount) public {
		emit KeeperRegistryLogic2_0_OwnerFundsWithdrawn_1(amount);
	}
	function emitKeeperRegistryBase1_3_PayeeshipTransferRequested_1(address keeper,address from,address to) public {
		emit KeeperRegistryBase1_3_PayeeshipTransferRequested_1(keeper,from,to);
	}
	function emitKeeperRegistryBase2_0_Paused_1(address account) public {
		emit KeeperRegistryBase2_0_Paused_1(account);
	}
	function emitVerifierProxy_OwnershipTransferRequested_1(address from,address to) public {
		emit VerifierProxy_OwnershipTransferRequested_1(from,to);
	}
	function emitKeeperRegistry1_3_OwnerFundsWithdrawn_1(uint96 amount) public {
		emit KeeperRegistry1_3_OwnerFundsWithdrawn_1(amount);
	}
	function emitKeeperRegistry1_3_UpkeepAdminTransferRequested_1(uint256 id,address from,address to) public {
		emit KeeperRegistry1_3_UpkeepAdminTransferRequested_1(id,from,to);
	}
	function emitKeeperRegistry1_3_UpkeepUnpaused_1(uint256 id) public {
		emit KeeperRegistry1_3_UpkeepUnpaused_1(id);
	}
	function emitKeeperRegistry2_0_Unpaused_1(address account) public {
		emit KeeperRegistry2_0_Unpaused_1(account);
	}
	function emitKeeperRegistryLogic2_0_CancelledUpkeepReport_1(uint256 id) public {
		emit KeeperRegistryLogic2_0_CancelledUpkeepReport_1(id);
	}
	function emitKeeperRegistry1_3_ConfigSet_1(uint32 config_paymentPremiumPPB,uint32 config_flatFeeMicroLink,uint24 config_blockCountPerTurn,uint32 config_checkGasLimit,uint24 config_stalenessSeconds,uint16 config_gasCeilingMultiplier,uint96 config_minUpkeepSpend,uint32 config_maxPerformGas,uint256 config_fallbackGasPrice,uint256 config_fallbackLinkPrice,address config_transcoder,address config_registrar) public {
		emit KeeperRegistry1_3_ConfigSet_1(config_paymentPremiumPPB,config_flatFeeMicroLink,config_blockCountPerTurn,config_checkGasLimit,config_stalenessSeconds,config_gasCeilingMultiplier,config_minUpkeepSpend,config_maxPerformGas,config_fallbackGasPrice,config_fallbackLinkPrice,config_transcoder,config_registrar);
	}
	function emitKeeperRegistry1_3_FundsWithdrawn_1(uint256 id,uint256 amount,address to) public {
		emit KeeperRegistry1_3_FundsWithdrawn_1(id,amount,to);
	}
	function emitKeeperRegistry2_0_UpkeepCanceled_1(uint256 id,uint64 atBlockHeight) public {
		emit KeeperRegistry2_0_UpkeepCanceled_1(id,atBlockHeight);
	}
	function emitAuthorizedOriginReceiverUpgradeable_AuthorizedSendersChanged_1(address[] memory senders,address changedBy) public {
		emit AuthorizedOriginReceiverUpgradeable_AuthorizedSendersChanged_1(senders,changedBy);
	}
	function emitConfirmedOwner_OwnershipTransferRequested_1(address from,address to) public {
		emit ConfirmedOwner_OwnershipTransferRequested_1(from,to);
	}
	function emitFunctionsBillingRegistry_SubscriptionCreated_1(uint64 subscriptionId,address owner) public {
		emit FunctionsBillingRegistry_SubscriptionCreated_1(subscriptionId,owner);
	}
	function emitFunctionsOracle_OwnershipTransferred_1(address from,address to) public {
		emit FunctionsOracle_OwnershipTransferred_1(from,to);
	}
	function emitKeeperRegistry1_2_ConfigSet_1(uint32 config_paymentPremiumPPB,uint32 config_flatFeeMicroLink,uint24 config_blockCountPerTurn,uint32 config_checkGasLimit,uint24 config_stalenessSeconds,uint16 config_gasCeilingMultiplier,uint96 config_minUpkeepSpend,uint32 config_maxPerformGas,uint256 config_fallbackGasPrice,uint256 config_fallbackLinkPrice,address config_transcoder,address config_registrar) public {
		emit KeeperRegistry1_2_ConfigSet_1(config_paymentPremiumPPB,config_flatFeeMicroLink,config_blockCountPerTurn,config_checkGasLimit,config_stalenessSeconds,config_gasCeilingMultiplier,config_minUpkeepSpend,config_maxPerformGas,config_fallbackGasPrice,config_fallbackLinkPrice,config_transcoder,config_registrar);
	}
	function emitKeeperRegistry2_0_UpkeepCheckDataUpdated_1(uint256 id,bytes memory newCheckData) public {
		emit KeeperRegistry2_0_UpkeepCheckDataUpdated_1(id,newCheckData);
	}
	function emitKeeperRegistryBase1_3_UpkeepRegistered_1(uint256 id,uint32 executeGas,address admin) public {
		emit KeeperRegistryBase1_3_UpkeepRegistered_1(id,executeGas,admin);
	}
	function emitKeeperRegistryLogic1_3_UpkeepPerformed_1(uint256 id,bool success,address from,uint96 payment,bytes memory performData) public {
		emit KeeperRegistryLogic1_3_UpkeepPerformed_1(id,success,from,payment,performData);
	}
	function emitFunctionsOracle_OracleRequest_1(bytes32 requestId,address requestingContract,address requestInitiator,uint64 subscriptionId,address subscriptionOwner,bytes memory data) public {
		emit FunctionsOracle_OracleRequest_1(requestId,requestingContract,requestInitiator,subscriptionId,subscriptionOwner,data);
	}
	function emitFunctionsOracleWithInit_OwnershipTransferred_1(address from,address to) public {
		emit FunctionsOracleWithInit_OwnershipTransferred_1(from,to);
	}
	function emitKeeperRegistry1_2_UpkeepGasLimitSet_1(uint256 id,uint96 gasLimit) public {
		emit KeeperRegistry1_2_UpkeepGasLimitSet_1(id,gasLimit);
	}
	function emitKeeperRegistry1_3_UpkeepAdminTransferred_1(uint256 id,address from,address to) public {
		emit KeeperRegistry1_3_UpkeepAdminTransferred_1(id,from,to);
	}
	function emitVRFCoordinatorV2_SubscriptionConsumerAdded_1(uint64 subId,address consumer) public {
		emit VRFCoordinatorV2_SubscriptionConsumerAdded_1(subId,consumer);
	}
	function emitFunctionsBillingRegistry_OwnershipTransferred_1(address from,address to) public {
		emit FunctionsBillingRegistry_OwnershipTransferred_1(from,to);
	}
	function emitFunctionsBillingRegistry_SubscriptionOwnerTransferred_1(uint64 subscriptionId,address from,address to) public {
		emit FunctionsBillingRegistry_SubscriptionOwnerTransferred_1(subscriptionId,from,to);
	}
	function emitFunctionsBillingRegistryWithInit_BillingStart_1(bytes32 requestId,uint64 commitment_subscriptionId,address commitment_client,uint32 commitment_gasLimit,uint256 commitment_gasPrice,address commitment_don,uint96 commitment_donFee,uint96 commitment_registryFee,uint96 commitment_estimatedCost,uint256 commitment_timestamp) public {
		emit FunctionsBillingRegistryWithInit_BillingStart_1(requestId,commitment_subscriptionId,commitment_client,commitment_gasLimit,commitment_gasPrice,commitment_don,commitment_donFee,commitment_registryFee,commitment_estimatedCost,commitment_timestamp);
	}
	function emitVRFV2Wrapper_WrapperFulfillmentFailed_1(uint256 requestId,address consumer) public {
		emit VRFV2Wrapper_WrapperFulfillmentFailed_1(requestId,consumer);
	}
	function emitTransparentUpgradeableProxy_BeaconUpgraded_1(address beacon) public {
		emit TransparentUpgradeableProxy_BeaconUpgraded_1(beacon);
	}
	function emitFunctionsOracle_Transmitted_1(bytes32 configDigest,uint32 epoch) public {
		emit FunctionsOracle_Transmitted_1(configDigest,epoch);
	}
	function emitKeeperRegistryBase1_3_OwnershipTransferRequested_1(address from,address to) public {
		emit KeeperRegistryBase1_3_OwnershipTransferRequested_1(from,to);
	}
	function emitKeeperRegistryBase1_3_UpkeepPaused_1(uint256 id) public {
		emit KeeperRegistryBase1_3_UpkeepPaused_1(id);
	}
	function emitKeeperRegistryBase2_0_InsufficientFundsUpkeepReport_1(uint256 id) public {
		emit KeeperRegistryBase2_0_InsufficientFundsUpkeepReport_1(id);
	}
	function emitKeeperRegistryBase2_0_StaleUpkeepReport_1(uint256 id) public {
		emit KeeperRegistryBase2_0_StaleUpkeepReport_1(id);
	}
	function emitFunctionsOracle_AuthorizedSendersChanged_1(address[] memory senders,address changedBy) public {
		emit FunctionsOracle_AuthorizedSendersChanged_1(senders,changedBy);
	}
	function emitKeeperRegistry1_3_PayeeshipTransferRequested_1(address keeper,address from,address to) public {
		emit KeeperRegistry1_3_PayeeshipTransferRequested_1(keeper,from,to);
	}
	function emitVerifierProxy_OwnershipTransferred_1(address from,address to) public {
		emit VerifierProxy_OwnershipTransferred_1(from,to);
	}
	function emitFunctionsBillingRegistryWithInit_SubscriptionConsumerAdded_1(uint64 subscriptionId,address consumer) public {
		emit FunctionsBillingRegistryWithInit_SubscriptionConsumerAdded_1(subscriptionId,consumer);
	}
	function emitKeeperRegistry1_2_OwnerFundsWithdrawn_1(uint96 amount) public {
		emit KeeperRegistry1_2_OwnerFundsWithdrawn_1(amount);
	}
	function emitKeeperRegistryBase1_3_OwnerFundsWithdrawn_1(uint96 amount) public {
		emit KeeperRegistryBase1_3_OwnerFundsWithdrawn_1(amount);
	}
	function emitOCR2BaseUpgradeable_ConfigSet_1(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
		emit OCR2BaseUpgradeable_ConfigSet_1(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
	}
	function emitKeeperRegistryBase1_3_UpkeepMigrated_1(uint256 id,uint256 remainingBalance,address destination) public {
		emit KeeperRegistryBase1_3_UpkeepMigrated_1(id,remainingBalance,destination);
	}
	function emitLogEmitter_Log2_1(uint256 mfhzxy) public {
		emit LogEmitter_Log2_1(mfhzxy);
	}
	function emitVRFCoordinatorV2_SubscriptionOwnerTransferred_1(uint64 subId,address from,address to) public {
		emit VRFCoordinatorV2_SubscriptionOwnerTransferred_1(subId,from,to);
	}
	function emitFunctionsBillingRegistryWithInit_SubscriptionFunded_1(uint64 subscriptionId,uint256 oldBalance,uint256 newBalance) public {
		emit FunctionsBillingRegistryWithInit_SubscriptionFunded_1(subscriptionId,oldBalance,newBalance);
	}
	function emitFunctionsOracle_AuthorizedSendersActive_1(address account) public {
		emit FunctionsOracle_AuthorizedSendersActive_1(account);
	}
	function emitFunctionsOracle_UserCallbackRawError_1(bytes32 requestId,bytes memory lowLevelData) public {
		emit FunctionsOracle_UserCallbackRawError_1(requestId,lowLevelData);
	}
	function emitFunctionsOracleWithInit_AuthorizedSendersActive_1(address account) public {
		emit FunctionsOracleWithInit_AuthorizedSendersActive_1(account);
	}
	function emitKeeperRegistry1_2_Paused_1(address account) public {
		emit KeeperRegistry1_2_Paused_1(account);
	}
	function emitVRFV2Wrapper_OwnershipTransferred_1(address from,address to) public {
		emit VRFV2Wrapper_OwnershipTransferred_1(from,to);
	}
	function emitVerifier_OwnershipTransferred_1(address from,address to) public {
		emit Verifier_OwnershipTransferred_1(from,to);
	}
	function emitKeeperRegistry2_0_InsufficientFundsUpkeepReport_1(uint256 id) public {
		emit KeeperRegistry2_0_InsufficientFundsUpkeepReport_1(id);
	}
	function emitKeeperRegistry2_0_ReorgedUpkeepReport_1(uint256 id) public {
		emit KeeperRegistry2_0_ReorgedUpkeepReport_1(id);
	}
	function emitKeeperRegistryBase1_3_UpkeepGasLimitSet_1(uint256 id,uint96 gasLimit) public {
		emit KeeperRegistryBase1_3_UpkeepGasLimitSet_1(id,gasLimit);
	}
	function emitFunctionsBillingRegistryWithInit_ConfigSet_1(uint32 maxGasLimit,uint32 stalenessSeconds,uint256 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,uint32 gasOverhead) public {
		emit FunctionsBillingRegistryWithInit_ConfigSet_1(maxGasLimit,stalenessSeconds,gasAfterPaymentCalculation,fallbackWeiPerUnitLink,gasOverhead);
	}
	function emitKeeperRegistry1_2_UpkeepPerformed_1(uint256 id,bool success,address from,uint96 payment,bytes memory performData) public {
		emit KeeperRegistry1_2_UpkeepPerformed_1(id,success,from,payment,performData);
	}
	function emitKeeperRegistry1_3_UpkeepRegistered_1(uint256 id,uint32 executeGas,address admin) public {
		emit KeeperRegistry1_3_UpkeepRegistered_1(id,executeGas,admin);
	}
	function emitKeeperRegistry2_0_CancelledUpkeepReport_1(uint256 id) public {
		emit KeeperRegistry2_0_CancelledUpkeepReport_1(id);
	}
	function emitKeeperRegistry2_0_ConfigSet_1(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
		emit KeeperRegistry2_0_ConfigSet_1(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
	}
	function emitKeeperRegistryBase1_3_UpkeepUnpaused_1(uint256 id) public {
		emit KeeperRegistryBase1_3_UpkeepUnpaused_1(id);
	}
	function emitKeeperRegistryLogic1_3_UpkeepCheckDataUpdated_1(uint256 id,bytes memory newCheckData) public {
		emit KeeperRegistryLogic1_3_UpkeepCheckDataUpdated_1(id,newCheckData);
	}
	function emitOCR2BaseUpgradeable_OwnershipTransferred_1(address from,address to) public {
		emit OCR2BaseUpgradeable_OwnershipTransferred_1(from,to);
	}
	function emitVRFCoordinatorV2_OwnershipTransferRequested_1(address from,address to) public {
		emit VRFCoordinatorV2_OwnershipTransferRequested_1(from,to);
	}
	function emitFunctionsOracle_OwnershipTransferRequested_1(address from,address to) public {
		emit FunctionsOracle_OwnershipTransferRequested_1(from,to);
	}
	function emitKeeperRegistryBase2_0_UpkeepPerformed_1(uint256 id,bool success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment) public {
		emit KeeperRegistryBase2_0_UpkeepPerformed_1(id,success,checkBlockNumber,gasUsed,gasOverhead,totalPayment);
	}
	function emitVRFCoordinatorV2_RandomWordsFulfilled_1(uint256 requestId,uint256 outputSeed,uint96 payment,bool success) public {
		emit VRFCoordinatorV2_RandomWordsFulfilled_1(requestId,outputSeed,payment,success);
	}
	function emitKeeperRegistryLogic1_3_FundsAdded_1(uint256 id,address from,uint96 amount) public {
		emit KeeperRegistryLogic1_3_FundsAdded_1(id,from,amount);
	}
	function emitFunctionsBillingRegistryWithInit_SubscriptionCreated_1(uint64 subscriptionId,address owner) public {
		emit FunctionsBillingRegistryWithInit_SubscriptionCreated_1(subscriptionId,owner);
	}
	function emitKeeperRegistry2_0_UpkeepGasLimitSet_1(uint256 id,uint96 gasLimit) public {
		emit KeeperRegistry2_0_UpkeepGasLimitSet_1(id,gasLimit);
	}
	function emitKeeperRegistry2_0_UpkeepMigrated_1(uint256 id,uint256 remainingBalance,address destination) public {
		emit KeeperRegistry2_0_UpkeepMigrated_1(id,remainingBalance,destination);
	}
	function emitKeeperRegistryBase1_3_KeepersUpdated_1(address[] memory keepers,address[] memory payees) public {
		emit KeeperRegistryBase1_3_KeepersUpdated_1(keepers,payees);
	}
	function emitKeeperRegistryBase1_3_PaymentWithdrawn_1(address keeper,uint256 amount,address to,address payee) public {
		emit KeeperRegistryBase1_3_PaymentWithdrawn_1(keeper,amount,to,payee);
	}
	function emitKeeperRegistryBase2_0_CancelledUpkeepReport_1(uint256 id) public {
		emit KeeperRegistryBase2_0_CancelledUpkeepReport_1(id);
	}
	function emitVRFConsumerV2UpgradeableExample_Initialized_1(uint8 version) public {
		emit VRFConsumerV2UpgradeableExample_Initialized_1(version);
	}
	function emitVRFCoordinatorV2_SubscriptionConsumerRemoved_1(uint64 subId,address consumer) public {
		emit VRFCoordinatorV2_SubscriptionConsumerRemoved_1(subId,consumer);
	}
	function emitCronUpkeep_CronJobExecuted_1(uint256 id,uint256 timestamp) public {
		emit CronUpkeep_CronJobExecuted_1(id,timestamp);
	}
	function emitFunctionsBillingRegistryWithInit_OwnershipTransferred_1(address from,address to) public {
		emit FunctionsBillingRegistryWithInit_OwnershipTransferred_1(from,to);
	}
	function emitKeeperRegistry1_2_FundsWithdrawn_1(uint256 id,uint256 amount,address to) public {
		emit KeeperRegistry1_2_FundsWithdrawn_1(id,amount,to);
	}
	function emitKeeperRegistry1_3_UpkeepPaused_1(uint256 id) public {
		emit KeeperRegistry1_3_UpkeepPaused_1(id);
	}
	function emitKeeperRegistry2_0_OwnershipTransferRequested_1(address from,address to) public {
		emit KeeperRegistry2_0_OwnershipTransferRequested_1(from,to);
	}
	function emitVRFCoordinatorV2TestHelper_SubscriptionConsumerRemoved_1(uint64 subId,address consumer) public {
		emit VRFCoordinatorV2TestHelper_SubscriptionConsumerRemoved_1(subId,consumer);
	}
	function emitConfirmedOwner_OwnershipTransferred_1(address from,address to) public {
		emit ConfirmedOwner_OwnershipTransferred_1(from,to);
	}
	function emitCronUpkeep_CronJobUpdated_1(uint256 id,address target,bytes memory handler) public {
		emit CronUpkeep_CronJobUpdated_1(id,target,handler);
	}
	function emitENSInterface_Transfer_1(bytes32 node,address owner) public {
		emit ENSInterface_Transfer_1(node,owner);
	}
	function emitFunctionsOracleWithInit_Initialized_1(uint8 version) public {
		emit FunctionsOracleWithInit_Initialized_1(version);
	}
	function emitKeeperRegistry1_3_OwnershipTransferred_1(address from,address to) public {
		emit KeeperRegistry1_3_OwnershipTransferred_1(from,to);
	}
	function emitOVM_GasPriceOracle_OverheadUpdated_1(uint256 vhectp) public {
		emit OVM_GasPriceOracle_OverheadUpdated_1(vhectp);
	}
	function emitFunctionsBillingRegistry_AuthorizedSendersChanged_1(address[] memory senders,address changedBy) public {
		emit FunctionsBillingRegistry_AuthorizedSendersChanged_1(senders,changedBy);
	}
	function emitFunctionsBillingRegistry_OwnershipTransferRequested_1(address from,address to) public {
		emit FunctionsBillingRegistry_OwnershipTransferRequested_1(from,to);
	}
	function emitFunctionsBillingRegistryWithInit_Initialized_1(uint8 version) public {
		emit FunctionsBillingRegistryWithInit_Initialized_1(version);
	}
	function emitKeeperRegistryBase1_3_UpkeepAdminTransferRequested_1(uint256 id,address from,address to) public {
		emit KeeperRegistryBase1_3_UpkeepAdminTransferRequested_1(id,from,to);
	}
	function emitKeeperRegistryBase2_0_ReorgedUpkeepReport_1(uint256 id) public {
		emit KeeperRegistryBase2_0_ReorgedUpkeepReport_1(id);
	}
	function emitOCR2Abstract_Transmitted_1(bytes32 configDigest,uint32 epoch) public {
		emit OCR2Abstract_Transmitted_1(configDigest,epoch);
	}
	function emitOVM_GasPriceOracle_ScalarUpdated_1(uint256 upncvl) public {
		emit OVM_GasPriceOracle_ScalarUpdated_1(upncvl);
	}
	function emitVRFCoordinatorV2_SubscriptionOwnerTransferRequested_1(uint64 subId,address from,address to) public {
		emit VRFCoordinatorV2_SubscriptionOwnerTransferRequested_1(subId,from,to);
	}
	function emitAuthorizedOriginReceiverUpgradeable_Initialized_1(uint8 version) public {
		emit AuthorizedOriginReceiverUpgradeable_Initialized_1(version);
	}
	function emitCronUpkeep_Unpaused_1(address account) public {
		emit CronUpkeep_Unpaused_1(account);
	}
	function emitKeeperRegistry1_3_PayeeshipTransferred_1(address keeper,address from,address to) public {
		emit KeeperRegistry1_3_PayeeshipTransferred_1(keeper,from,to);
	}
	function emitKeeperRegistry2_0_UpkeepRegistered_1(uint256 id,uint32 executeGas,address admin) public {
		emit KeeperRegistry2_0_UpkeepRegistered_1(id,executeGas,admin);
	}
	function emitKeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferred_1(address from,address to) public {
		emit KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferred_1(from,to);
	}
}
