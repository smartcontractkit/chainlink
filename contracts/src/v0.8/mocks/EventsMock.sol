// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

contract EventsMock {
    struct EventsMock_FunctionsBillingRegistryWithInit_Commitment {uint64 subscriptionId;address client;uint32 gasLimit;uint256 gasPrice;address don;uint96 donFee;uint96 registryFee;uint96 estimatedCost;uint256 timestamp; }
    struct EventsMock_FunctionsBillingRegistry_Commitment {uint64 subscriptionId;address client;uint32 gasLimit;uint256 gasPrice;address don;uint96 donFee;uint96 registryFee;uint96 estimatedCost;uint256 timestamp; }
    struct EventsMock_KeeperRegistry1_2_Config {uint32 paymentPremiumPPB;uint32 flatFeeMicroLink;uint24 blockCountPerTurn;uint32 checkGasLimit;uint24 stalenessSeconds;uint16 gasCeilingMultiplier;uint96 minUpkeepSpend;uint32 maxPerformGas;uint256 fallbackGasPrice;uint256 fallbackLinkPrice;address transcoder;address registrar; }
    struct EventsMock_KeeperRegistry1_3_Config {uint32 paymentPremiumPPB;uint32 flatFeeMicroLink;uint24 blockCountPerTurn;uint32 checkGasLimit;uint24 stalenessSeconds;uint16 gasCeilingMultiplier;uint96 minUpkeepSpend;uint32 maxPerformGas;uint256 fallbackGasPrice;uint256 fallbackLinkPrice;address transcoder;address registrar; }
    struct EventsMock_KeeperRegistryBase1_3_Config {uint32 paymentPremiumPPB;uint32 flatFeeMicroLink;uint24 blockCountPerTurn;uint32 checkGasLimit;uint24 stalenessSeconds;uint16 gasCeilingMultiplier;uint96 minUpkeepSpend;uint32 maxPerformGas;uint256 fallbackGasPrice;uint256 fallbackLinkPrice;address transcoder;address registrar; }
    struct EventsMock_KeeperRegistryLogic1_3_Config {uint32 paymentPremiumPPB;uint32 flatFeeMicroLink;uint24 blockCountPerTurn;uint32 checkGasLimit;uint24 stalenessSeconds;uint16 gasCeilingMultiplier;uint96 minUpkeepSpend;uint32 maxPerformGas;uint256 fallbackGasPrice;uint256 fallbackLinkPrice;address transcoder;address registrar; }
    struct EventsMock_VRFCoordinatorV2TestHelper_FeeConfig {uint32 fulfillmentFlatFeeLinkPPMTier1;uint32 fulfillmentFlatFeeLinkPPMTier2;uint32 fulfillmentFlatFeeLinkPPMTier3;uint32 fulfillmentFlatFeeLinkPPMTier4;uint32 fulfillmentFlatFeeLinkPPMTier5;uint24 reqsForTier2;uint24 reqsForTier3;uint24 reqsForTier4;uint24 reqsForTier5; }
    struct EventsMock_VRFCoordinatorV2_FeeConfig {uint32 fulfillmentFlatFeeLinkPPMTier1;uint32 fulfillmentFlatFeeLinkPPMTier2;uint32 fulfillmentFlatFeeLinkPPMTier3;uint32 fulfillmentFlatFeeLinkPPMTier4;uint32 fulfillmentFlatFeeLinkPPMTier5;uint24 reqsForTier2;uint24 reqsForTier3;uint24 reqsForTier4;uint24 reqsForTier5; }
    struct FunctionsBillingRegistry_Commitment {uint64 subscriptionId;address client;uint32 gasLimit;uint256 gasPrice;address don;uint96 donFee;uint96 registryFee;uint96 estimatedCost;uint256 timestamp; }
    struct FunctionsBillingRegistryWithInit_Commitment {uint64 subscriptionId;address client;uint32 gasLimit;uint256 gasPrice;address don;uint96 donFee;uint96 registryFee;uint96 estimatedCost;uint256 timestamp; }
    struct KeeperRegistry1_2_Config {uint32 paymentPremiumPPB;uint32 flatFeeMicroLink;uint24 blockCountPerTurn;uint32 checkGasLimit;uint24 stalenessSeconds;uint16 gasCeilingMultiplier;uint96 minUpkeepSpend;uint32 maxPerformGas;uint256 fallbackGasPrice;uint256 fallbackLinkPrice;address transcoder;address registrar; }
    struct KeeperRegistry1_3_Config {uint32 paymentPremiumPPB;uint32 flatFeeMicroLink;uint24 blockCountPerTurn;uint32 checkGasLimit;uint24 stalenessSeconds;uint16 gasCeilingMultiplier;uint96 minUpkeepSpend;uint32 maxPerformGas;uint256 fallbackGasPrice;uint256 fallbackLinkPrice;address transcoder;address registrar; }
    struct KeeperRegistryBase1_3_Config {uint32 paymentPremiumPPB;uint32 flatFeeMicroLink;uint24 blockCountPerTurn;uint32 checkGasLimit;uint24 stalenessSeconds;uint16 gasCeilingMultiplier;uint96 minUpkeepSpend;uint32 maxPerformGas;uint256 fallbackGasPrice;uint256 fallbackLinkPrice;address transcoder;address registrar; }
    struct KeeperRegistryLogic1_3_Config {uint32 paymentPremiumPPB;uint32 flatFeeMicroLink;uint24 blockCountPerTurn;uint32 checkGasLimit;uint24 stalenessSeconds;uint16 gasCeilingMultiplier;uint96 minUpkeepSpend;uint32 maxPerformGas;uint256 fallbackGasPrice;uint256 fallbackLinkPrice;address transcoder;address registrar; }
    struct VRFCoordinatorV2_FeeConfig {uint32 fulfillmentFlatFeeLinkPPMTier1;uint32 fulfillmentFlatFeeLinkPPMTier2;uint32 fulfillmentFlatFeeLinkPPMTier3;uint32 fulfillmentFlatFeeLinkPPMTier4;uint32 fulfillmentFlatFeeLinkPPMTier5;uint24 reqsForTier2;uint24 reqsForTier3;uint24 reqsForTier4;uint24 reqsForTier5; }
    struct VRFCoordinatorV2TestHelper_FeeConfig {uint32 fulfillmentFlatFeeLinkPPMTier1;uint32 fulfillmentFlatFeeLinkPPMTier2;uint32 fulfillmentFlatFeeLinkPPMTier3;uint32 fulfillmentFlatFeeLinkPPMTier4;uint32 fulfillmentFlatFeeLinkPPMTier5;uint24 reqsForTier2;uint24 reqsForTier3;uint24 reqsForTier4;uint24 reqsForTier5; }
    event AggregatorInterface_AnswerUpdated(int256 indexed current,uint256 indexed roundId,uint256 updatedAt);
    event AggregatorInterface_NewRound(uint256 indexed roundId,address indexed startedBy,uint256 startedAt);
    event AggregatorV2V3Interface_AnswerUpdated(int256 indexed current,uint256 indexed roundId,uint256 updatedAt);
    event AggregatorV2V3Interface_NewRound(uint256 indexed roundId,address indexed startedBy,uint256 startedAt);
    event AuthorizedOriginReceiverUpgradeable_AuthorizedSendersActive(address account);
    event AuthorizedOriginReceiverUpgradeable_AuthorizedSendersChanged(address[] senders,address changedBy);
    event AuthorizedOriginReceiverUpgradeable_AuthorizedSendersDeactive(address account);
    event AuthorizedOriginReceiverUpgradeable_Initialized(uint8 version);
    event AuthorizedReceiver_AuthorizedSendersChanged(address[] senders,address changedBy);
    event BatchVRFCoordinatorV2_ErrorReturned(uint256 indexed requestId,string reason);
    event BatchVRFCoordinatorV2_RawErrorReturned(uint256 indexed requestId,bytes lowLevelData);
    event ChainlinkClient_ChainlinkCancelled(bytes32 indexed id);
    event ChainlinkClient_ChainlinkFulfilled(bytes32 indexed id);
    event ChainlinkClient_ChainlinkRequested(bytes32 indexed id);
    event ConfirmedOwner_OwnershipTransferRequested(address indexed from,address indexed to);
    event ConfirmedOwner_OwnershipTransferred(address indexed from,address indexed to);
    event ConfirmedOwnerUpgradeable_Initialized(uint8 version);
    event ConfirmedOwnerUpgradeable_OwnershipTransferRequested(address indexed from,address indexed to);
    event ConfirmedOwnerUpgradeable_OwnershipTransferred(address indexed from,address indexed to);
    event ConfirmedOwnerWithProposal_OwnershipTransferRequested(address indexed from,address indexed to);
    event ConfirmedOwnerWithProposal_OwnershipTransferred(address indexed from,address indexed to);
    event ContextUpgradeable_Initialized(uint8 version);
    event CronUpkeep_CronJobCreated(uint256 indexed id,address target,bytes handler);
    event CronUpkeep_CronJobDeleted(uint256 indexed id);
    event CronUpkeep_CronJobExecuted(uint256 indexed id,uint256 timestamp);
    event CronUpkeep_CronJobUpdated(uint256 indexed id,address target,bytes handler);
    event CronUpkeep_OwnershipTransferRequested(address indexed from,address indexed to);
    event CronUpkeep_OwnershipTransferred(address indexed from,address indexed to);
    event CronUpkeep_Paused(address account);
    event CronUpkeep_Unpaused(address account);
    event CronUpkeepFactory_NewCronUpkeepCreated(address upkeep,address owner);
    event CronUpkeepFactory_OwnershipTransferRequested(address indexed from,address indexed to);
    event CronUpkeepFactory_OwnershipTransferred(address indexed from,address indexed to);
    event ENSInterface_NewOwner(bytes32 indexed node,bytes32 indexed label,address owner);
    event ENSInterface_NewResolver(bytes32 indexed node,address resolver);
    event ENSInterface_NewTTL(bytes32 indexed node,uint64 ttl);
    event ENSInterface_Transfer(bytes32 indexed node,address owner);
    event ERC1967Proxy_AdminChanged(address previousAdmin,address newAdmin);
    event ERC1967Proxy_BeaconUpgraded(address indexed beacon);
    event ERC1967Proxy_Upgraded(address indexed implementation);
    event ERC1967Upgrade_AdminChanged(address previousAdmin,address newAdmin);
    event ERC1967Upgrade_BeaconUpgraded(address indexed beacon);
    event ERC1967Upgrade_Upgraded(address indexed implementation);
    event EventsMock_AggregatorInterface_AnswerUpdated(int256 indexed current,uint256 indexed roundId,uint256 updatedAt);
    event EventsMock_AggregatorInterface_NewRound(uint256 indexed roundId,address indexed startedBy,uint256 startedAt);
    event EventsMock_AggregatorV2V3Interface_AnswerUpdated(int256 indexed current,uint256 indexed roundId,uint256 updatedAt);
    event EventsMock_AggregatorV2V3Interface_NewRound(uint256 indexed roundId,address indexed startedBy,uint256 startedAt);
    event EventsMock_AuthorizedOriginReceiverUpgradeable_AuthorizedSendersActive(address account);
    event EventsMock_AuthorizedOriginReceiverUpgradeable_AuthorizedSendersChanged(address[] senders,address changedBy);
    event EventsMock_AuthorizedOriginReceiverUpgradeable_AuthorizedSendersDeactive(address account);
    event EventsMock_AuthorizedOriginReceiverUpgradeable_Initialized(uint8 version);
    event EventsMock_AuthorizedReceiver_AuthorizedSendersChanged(address[] senders,address changedBy);
    event EventsMock_BatchVRFCoordinatorV2_ErrorReturned(uint256 indexed requestId,string reason);
    event EventsMock_BatchVRFCoordinatorV2_RawErrorReturned(uint256 indexed requestId,bytes lowLevelData);
    event EventsMock_ChainlinkClient_ChainlinkCancelled(bytes32 indexed id);
    event EventsMock_ChainlinkClient_ChainlinkFulfilled(bytes32 indexed id);
    event EventsMock_ChainlinkClient_ChainlinkRequested(bytes32 indexed id);
    event EventsMock_ConfirmedOwnerUpgradeable_Initialized(uint8 version);
    event EventsMock_ConfirmedOwnerUpgradeable_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_ConfirmedOwnerUpgradeable_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_ConfirmedOwnerWithProposal_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_ConfirmedOwnerWithProposal_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_ConfirmedOwner_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_ConfirmedOwner_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_ContextUpgradeable_Initialized(uint8 version);
    event EventsMock_CronUpkeepFactory_NewCronUpkeepCreated(address upkeep,address owner);
    event EventsMock_CronUpkeepFactory_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_CronUpkeepFactory_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_CronUpkeep_CronJobCreated(uint256 indexed id,address target,bytes handler);
    event EventsMock_CronUpkeep_CronJobDeleted(uint256 indexed id);
    event EventsMock_CronUpkeep_CronJobExecuted(uint256 indexed id,uint256 timestamp);
    event EventsMock_CronUpkeep_CronJobUpdated(uint256 indexed id,address target,bytes handler);
    event EventsMock_CronUpkeep_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_CronUpkeep_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_CronUpkeep_Paused(address account);
    event EventsMock_CronUpkeep_Unpaused(address account);
    event EventsMock_ENSInterface_NewOwner(bytes32 indexed node,bytes32 indexed label,address owner);
    event EventsMock_ENSInterface_NewResolver(bytes32 indexed node,address resolver);
    event EventsMock_ENSInterface_NewTTL(bytes32 indexed node,uint64 ttl);
    event EventsMock_ENSInterface_Transfer(bytes32 indexed node,address owner);
    event EventsMock_ERC1967Proxy_AdminChanged(address previousAdmin,address newAdmin);
    event EventsMock_ERC1967Proxy_BeaconUpgraded(address indexed beacon);
    event EventsMock_ERC1967Proxy_Upgraded(address indexed implementation);
    event EventsMock_ERC1967Upgrade_AdminChanged(address previousAdmin,address newAdmin);
    event EventsMock_ERC1967Upgrade_BeaconUpgraded(address indexed beacon);
    event EventsMock_ERC1967Upgrade_Upgraded(address indexed implementation);
    event EventsMock_FunctionsBillingRegistryWithInit_AuthorizedSendersChanged(address[] senders,address changedBy);
    event EventsMock_FunctionsBillingRegistryWithInit_BillingEnd(bytes32 indexed requestId,uint64 subscriptionId,uint96 signerPayment,uint96 transmitterPayment,uint96 totalCost,bool success);
    event EventsMock_FunctionsBillingRegistryWithInit_BillingStart(bytes32 indexed requestId,EventsMock_FunctionsBillingRegistryWithInit_Commitment commitment);
    event EventsMock_FunctionsBillingRegistryWithInit_ConfigSet(uint32 maxGasLimit,uint32 stalenessSeconds,uint256 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,uint32 gasOverhead);
    event EventsMock_FunctionsBillingRegistryWithInit_FundsRecovered(address to,uint256 amount);
    event EventsMock_FunctionsBillingRegistryWithInit_Initialized(uint8 version);
    event EventsMock_FunctionsBillingRegistryWithInit_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_FunctionsBillingRegistryWithInit_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_FunctionsBillingRegistryWithInit_Paused(address account);
    event EventsMock_FunctionsBillingRegistryWithInit_RequestTimedOut(bytes32 indexed requestId);
    event EventsMock_FunctionsBillingRegistryWithInit_SubscriptionCanceled(uint64 indexed subscriptionId,address to,uint256 amount);
    event EventsMock_FunctionsBillingRegistryWithInit_SubscriptionConsumerAdded(uint64 indexed subscriptionId,address consumer);
    event EventsMock_FunctionsBillingRegistryWithInit_SubscriptionConsumerRemoved(uint64 indexed subscriptionId,address consumer);
    event EventsMock_FunctionsBillingRegistryWithInit_SubscriptionCreated(uint64 indexed subscriptionId,address owner);
    event EventsMock_FunctionsBillingRegistryWithInit_SubscriptionFunded(uint64 indexed subscriptionId,uint256 oldBalance,uint256 newBalance);
    event EventsMock_FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferRequested(uint64 indexed subscriptionId,address from,address to);
    event EventsMock_FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferred(uint64 indexed subscriptionId,address from,address to);
    event EventsMock_FunctionsBillingRegistryWithInit_Unpaused(address account);
    event EventsMock_FunctionsBillingRegistry_AuthorizedSendersChanged(address[] senders,address changedBy);
    event EventsMock_FunctionsBillingRegistry_BillingEnd(bytes32 indexed requestId,uint64 subscriptionId,uint96 signerPayment,uint96 transmitterPayment,uint96 totalCost,bool success);
    event EventsMock_FunctionsBillingRegistry_BillingStart(bytes32 indexed requestId,EventsMock_FunctionsBillingRegistry_Commitment commitment);
    event EventsMock_FunctionsBillingRegistry_ConfigSet(uint32 maxGasLimit,uint32 stalenessSeconds,uint256 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,uint32 gasOverhead);
    event EventsMock_FunctionsBillingRegistry_FundsRecovered(address to,uint256 amount);
    event EventsMock_FunctionsBillingRegistry_Initialized(uint8 version);
    event EventsMock_FunctionsBillingRegistry_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_FunctionsBillingRegistry_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_FunctionsBillingRegistry_Paused(address account);
    event EventsMock_FunctionsBillingRegistry_RequestTimedOut(bytes32 indexed requestId);
    event EventsMock_FunctionsBillingRegistry_SubscriptionCanceled(uint64 indexed subscriptionId,address to,uint256 amount);
    event EventsMock_FunctionsBillingRegistry_SubscriptionConsumerAdded(uint64 indexed subscriptionId,address consumer);
    event EventsMock_FunctionsBillingRegistry_SubscriptionConsumerRemoved(uint64 indexed subscriptionId,address consumer);
    event EventsMock_FunctionsBillingRegistry_SubscriptionCreated(uint64 indexed subscriptionId,address owner);
    event EventsMock_FunctionsBillingRegistry_SubscriptionFunded(uint64 indexed subscriptionId,uint256 oldBalance,uint256 newBalance);
    event EventsMock_FunctionsBillingRegistry_SubscriptionOwnerTransferRequested(uint64 indexed subscriptionId,address from,address to);
    event EventsMock_FunctionsBillingRegistry_SubscriptionOwnerTransferred(uint64 indexed subscriptionId,address from,address to);
    event EventsMock_FunctionsBillingRegistry_Unpaused(address account);
    event EventsMock_FunctionsClientExample_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_FunctionsClientExample_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_FunctionsClientExample_RequestFulfilled(bytes32 indexed id);
    event EventsMock_FunctionsClientExample_RequestSent(bytes32 indexed id);
    event EventsMock_FunctionsClient_RequestFulfilled(bytes32 indexed id);
    event EventsMock_FunctionsClient_RequestSent(bytes32 indexed id);
    event EventsMock_FunctionsOracleWithInit_AuthorizedSendersActive(address account);
    event EventsMock_FunctionsOracleWithInit_AuthorizedSendersChanged(address[] senders,address changedBy);
    event EventsMock_FunctionsOracleWithInit_AuthorizedSendersDeactive(address account);
    event EventsMock_FunctionsOracleWithInit_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] signers,address[] transmitters,uint8 f,bytes onchainConfig,uint64 offchainConfigVersion,bytes offchainConfig);
    event EventsMock_FunctionsOracleWithInit_Initialized(uint8 version);
    event EventsMock_FunctionsOracleWithInit_InvalidRequestID(bytes32 indexed requestId);
    event EventsMock_FunctionsOracleWithInit_OracleRequest(bytes32 indexed requestId,address requestingContract,address requestInitiator,uint64 subscriptionId,address subscriptionOwner,bytes data);
    event EventsMock_FunctionsOracleWithInit_OracleResponse(bytes32 indexed requestId);
    event EventsMock_FunctionsOracleWithInit_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_FunctionsOracleWithInit_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_FunctionsOracleWithInit_Transmitted(bytes32 configDigest,uint32 epoch);
    event EventsMock_FunctionsOracleWithInit_UserCallbackError(bytes32 indexed requestId,string reason);
    event EventsMock_FunctionsOracleWithInit_UserCallbackRawError(bytes32 indexed requestId,bytes lowLevelData);
    event EventsMock_FunctionsOracle_AuthorizedSendersActive(address account);
    event EventsMock_FunctionsOracle_AuthorizedSendersChanged(address[] senders,address changedBy);
    event EventsMock_FunctionsOracle_AuthorizedSendersDeactive(address account);
    event EventsMock_FunctionsOracle_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] signers,address[] transmitters,uint8 f,bytes onchainConfig,uint64 offchainConfigVersion,bytes offchainConfig);
    event EventsMock_FunctionsOracle_Initialized(uint8 version);
    event EventsMock_FunctionsOracle_InvalidRequestID(bytes32 indexed requestId);
    event EventsMock_FunctionsOracle_OracleRequest(bytes32 indexed requestId,address requestingContract,address requestInitiator,uint64 subscriptionId,address subscriptionOwner,bytes data);
    event EventsMock_FunctionsOracle_OracleResponse(bytes32 indexed requestId);
    event EventsMock_FunctionsOracle_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_FunctionsOracle_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_FunctionsOracle_Transmitted(bytes32 configDigest,uint32 epoch);
    event EventsMock_FunctionsOracle_UserCallbackError(bytes32 indexed requestId,string reason);
    event EventsMock_FunctionsOracle_UserCallbackRawError(bytes32 indexed requestId,bytes lowLevelData);
    event EventsMock_Initializable_Initialized(uint8 version);
    event EventsMock_KeeperRegistry1_2_ConfigSet(EventsMock_KeeperRegistry1_2_Config config);
    event EventsMock_KeeperRegistry1_2_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event EventsMock_KeeperRegistry1_2_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event EventsMock_KeeperRegistry1_2_KeepersUpdated(address[] keepers,address[] payees);
    event EventsMock_KeeperRegistry1_2_OwnerFundsWithdrawn(uint96 amount);
    event EventsMock_KeeperRegistry1_2_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_KeeperRegistry1_2_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_KeeperRegistry1_2_Paused(address account);
    event EventsMock_KeeperRegistry1_2_PayeeshipTransferRequested(address indexed keeper,address indexed from,address indexed to);
    event EventsMock_KeeperRegistry1_2_PayeeshipTransferred(address indexed keeper,address indexed from,address indexed to);
    event EventsMock_KeeperRegistry1_2_PaymentWithdrawn(address indexed keeper,uint256 indexed amount,address indexed to,address payee);
    event EventsMock_KeeperRegistry1_2_Unpaused(address account);
    event EventsMock_KeeperRegistry1_2_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event EventsMock_KeeperRegistry1_2_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event EventsMock_KeeperRegistry1_2_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event EventsMock_KeeperRegistry1_2_UpkeepPerformed(uint256 indexed id,bool indexed success,address indexed from,uint96 payment,bytes performData);
    event EventsMock_KeeperRegistry1_2_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event EventsMock_KeeperRegistry1_2_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event EventsMock_KeeperRegistry1_3_ConfigSet(EventsMock_KeeperRegistry1_3_Config config);
    event EventsMock_KeeperRegistry1_3_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event EventsMock_KeeperRegistry1_3_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event EventsMock_KeeperRegistry1_3_KeepersUpdated(address[] keepers,address[] payees);
    event EventsMock_KeeperRegistry1_3_OwnerFundsWithdrawn(uint96 amount);
    event EventsMock_KeeperRegistry1_3_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_KeeperRegistry1_3_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_KeeperRegistry1_3_Paused(address account);
    event EventsMock_KeeperRegistry1_3_PayeeshipTransferRequested(address indexed keeper,address indexed from,address indexed to);
    event EventsMock_KeeperRegistry1_3_PayeeshipTransferred(address indexed keeper,address indexed from,address indexed to);
    event EventsMock_KeeperRegistry1_3_PaymentWithdrawn(address indexed keeper,uint256 indexed amount,address indexed to,address payee);
    event EventsMock_KeeperRegistry1_3_Unpaused(address account);
    event EventsMock_KeeperRegistry1_3_UpkeepAdminTransferRequested(uint256 indexed id,address indexed from,address indexed to);
    event EventsMock_KeeperRegistry1_3_UpkeepAdminTransferred(uint256 indexed id,address indexed from,address indexed to);
    event EventsMock_KeeperRegistry1_3_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event EventsMock_KeeperRegistry1_3_UpkeepCheckDataUpdated(uint256 indexed id,bytes newCheckData);
    event EventsMock_KeeperRegistry1_3_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event EventsMock_KeeperRegistry1_3_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event EventsMock_KeeperRegistry1_3_UpkeepPaused(uint256 indexed id);
    event EventsMock_KeeperRegistry1_3_UpkeepPerformed(uint256 indexed id,bool indexed success,address indexed from,uint96 payment,bytes performData);
    event EventsMock_KeeperRegistry1_3_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event EventsMock_KeeperRegistry1_3_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event EventsMock_KeeperRegistry1_3_UpkeepUnpaused(uint256 indexed id);
    event EventsMock_KeeperRegistry2_0_CancelledUpkeepReport(uint256 indexed id);
    event EventsMock_KeeperRegistry2_0_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] signers,address[] transmitters,uint8 f,bytes onchainConfig,uint64 offchainConfigVersion,bytes offchainConfig);
    event EventsMock_KeeperRegistry2_0_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event EventsMock_KeeperRegistry2_0_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event EventsMock_KeeperRegistry2_0_InsufficientFundsUpkeepReport(uint256 indexed id);
    event EventsMock_KeeperRegistry2_0_OwnerFundsWithdrawn(uint96 amount);
    event EventsMock_KeeperRegistry2_0_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_KeeperRegistry2_0_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_KeeperRegistry2_0_Paused(address account);
    event EventsMock_KeeperRegistry2_0_PayeesUpdated(address[] transmitters,address[] payees);
    event EventsMock_KeeperRegistry2_0_PayeeshipTransferRequested(address indexed transmitter,address indexed from,address indexed to);
    event EventsMock_KeeperRegistry2_0_PayeeshipTransferred(address indexed transmitter,address indexed from,address indexed to);
    event EventsMock_KeeperRegistry2_0_PaymentWithdrawn(address indexed transmitter,uint256 indexed amount,address indexed to,address payee);
    event EventsMock_KeeperRegistry2_0_ReorgedUpkeepReport(uint256 indexed id);
    event EventsMock_KeeperRegistry2_0_StaleUpkeepReport(uint256 indexed id);
    event EventsMock_KeeperRegistry2_0_Transmitted(bytes32 configDigest,uint32 epoch);
    event EventsMock_KeeperRegistry2_0_Unpaused(address account);
    event EventsMock_KeeperRegistry2_0_UpkeepAdminTransferRequested(uint256 indexed id,address indexed from,address indexed to);
    event EventsMock_KeeperRegistry2_0_UpkeepAdminTransferred(uint256 indexed id,address indexed from,address indexed to);
    event EventsMock_KeeperRegistry2_0_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event EventsMock_KeeperRegistry2_0_UpkeepCheckDataUpdated(uint256 indexed id,bytes newCheckData);
    event EventsMock_KeeperRegistry2_0_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event EventsMock_KeeperRegistry2_0_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event EventsMock_KeeperRegistry2_0_UpkeepOffchainConfigSet(uint256 indexed id,bytes offchainConfig);
    event EventsMock_KeeperRegistry2_0_UpkeepPaused(uint256 indexed id);
    event EventsMock_KeeperRegistry2_0_UpkeepPerformed(uint256 indexed id,bool indexed success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment);
    event EventsMock_KeeperRegistry2_0_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event EventsMock_KeeperRegistry2_0_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event EventsMock_KeeperRegistry2_0_UpkeepUnpaused(uint256 indexed id);
    event EventsMock_KeeperRegistryBase1_3_ConfigSet(EventsMock_KeeperRegistryBase1_3_Config config);
    event EventsMock_KeeperRegistryBase1_3_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event EventsMock_KeeperRegistryBase1_3_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event EventsMock_KeeperRegistryBase1_3_KeepersUpdated(address[] keepers,address[] payees);
    event EventsMock_KeeperRegistryBase1_3_OwnerFundsWithdrawn(uint96 amount);
    event EventsMock_KeeperRegistryBase1_3_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_KeeperRegistryBase1_3_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_KeeperRegistryBase1_3_Paused(address account);
    event EventsMock_KeeperRegistryBase1_3_PayeeshipTransferRequested(address indexed keeper,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryBase1_3_PayeeshipTransferred(address indexed keeper,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryBase1_3_PaymentWithdrawn(address indexed keeper,uint256 indexed amount,address indexed to,address payee);
    event EventsMock_KeeperRegistryBase1_3_Unpaused(address account);
    event EventsMock_KeeperRegistryBase1_3_UpkeepAdminTransferRequested(uint256 indexed id,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryBase1_3_UpkeepAdminTransferred(uint256 indexed id,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryBase1_3_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event EventsMock_KeeperRegistryBase1_3_UpkeepCheckDataUpdated(uint256 indexed id,bytes newCheckData);
    event EventsMock_KeeperRegistryBase1_3_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event EventsMock_KeeperRegistryBase1_3_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event EventsMock_KeeperRegistryBase1_3_UpkeepPaused(uint256 indexed id);
    event EventsMock_KeeperRegistryBase1_3_UpkeepPerformed(uint256 indexed id,bool indexed success,address indexed from,uint96 payment,bytes performData);
    event EventsMock_KeeperRegistryBase1_3_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event EventsMock_KeeperRegistryBase1_3_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event EventsMock_KeeperRegistryBase1_3_UpkeepUnpaused(uint256 indexed id);
    event EventsMock_KeeperRegistryBase2_0_CancelledUpkeepReport(uint256 indexed id);
    event EventsMock_KeeperRegistryBase2_0_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event EventsMock_KeeperRegistryBase2_0_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event EventsMock_KeeperRegistryBase2_0_InsufficientFundsUpkeepReport(uint256 indexed id);
    event EventsMock_KeeperRegistryBase2_0_OwnerFundsWithdrawn(uint96 amount);
    event EventsMock_KeeperRegistryBase2_0_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_KeeperRegistryBase2_0_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_KeeperRegistryBase2_0_Paused(address account);
    event EventsMock_KeeperRegistryBase2_0_PayeesUpdated(address[] transmitters,address[] payees);
    event EventsMock_KeeperRegistryBase2_0_PayeeshipTransferRequested(address indexed transmitter,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryBase2_0_PayeeshipTransferred(address indexed transmitter,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryBase2_0_PaymentWithdrawn(address indexed transmitter,uint256 indexed amount,address indexed to,address payee);
    event EventsMock_KeeperRegistryBase2_0_ReorgedUpkeepReport(uint256 indexed id);
    event EventsMock_KeeperRegistryBase2_0_StaleUpkeepReport(uint256 indexed id);
    event EventsMock_KeeperRegistryBase2_0_Unpaused(address account);
    event EventsMock_KeeperRegistryBase2_0_UpkeepAdminTransferRequested(uint256 indexed id,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryBase2_0_UpkeepAdminTransferred(uint256 indexed id,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryBase2_0_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event EventsMock_KeeperRegistryBase2_0_UpkeepCheckDataUpdated(uint256 indexed id,bytes newCheckData);
    event EventsMock_KeeperRegistryBase2_0_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event EventsMock_KeeperRegistryBase2_0_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event EventsMock_KeeperRegistryBase2_0_UpkeepOffchainConfigSet(uint256 indexed id,bytes offchainConfig);
    event EventsMock_KeeperRegistryBase2_0_UpkeepPaused(uint256 indexed id);
    event EventsMock_KeeperRegistryBase2_0_UpkeepPerformed(uint256 indexed id,bool indexed success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment);
    event EventsMock_KeeperRegistryBase2_0_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event EventsMock_KeeperRegistryBase2_0_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event EventsMock_KeeperRegistryBase2_0_UpkeepUnpaused(uint256 indexed id);
    event EventsMock_KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic1_3_ConfigSet(EventsMock_KeeperRegistryLogic1_3_Config config);
    event EventsMock_KeeperRegistryLogic1_3_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event EventsMock_KeeperRegistryLogic1_3_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event EventsMock_KeeperRegistryLogic1_3_KeepersUpdated(address[] keepers,address[] payees);
    event EventsMock_KeeperRegistryLogic1_3_OwnerFundsWithdrawn(uint96 amount);
    event EventsMock_KeeperRegistryLogic1_3_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic1_3_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic1_3_Paused(address account);
    event EventsMock_KeeperRegistryLogic1_3_PayeeshipTransferRequested(address indexed keeper,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic1_3_PayeeshipTransferred(address indexed keeper,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic1_3_PaymentWithdrawn(address indexed keeper,uint256 indexed amount,address indexed to,address payee);
    event EventsMock_KeeperRegistryLogic1_3_Unpaused(address account);
    event EventsMock_KeeperRegistryLogic1_3_UpkeepAdminTransferRequested(uint256 indexed id,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic1_3_UpkeepAdminTransferred(uint256 indexed id,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic1_3_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event EventsMock_KeeperRegistryLogic1_3_UpkeepCheckDataUpdated(uint256 indexed id,bytes newCheckData);
    event EventsMock_KeeperRegistryLogic1_3_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event EventsMock_KeeperRegistryLogic1_3_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event EventsMock_KeeperRegistryLogic1_3_UpkeepPaused(uint256 indexed id);
    event EventsMock_KeeperRegistryLogic1_3_UpkeepPerformed(uint256 indexed id,bool indexed success,address indexed from,uint96 payment,bytes performData);
    event EventsMock_KeeperRegistryLogic1_3_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event EventsMock_KeeperRegistryLogic1_3_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event EventsMock_KeeperRegistryLogic1_3_UpkeepUnpaused(uint256 indexed id);
    event EventsMock_KeeperRegistryLogic2_0_CancelledUpkeepReport(uint256 indexed id);
    event EventsMock_KeeperRegistryLogic2_0_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event EventsMock_KeeperRegistryLogic2_0_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event EventsMock_KeeperRegistryLogic2_0_InsufficientFundsUpkeepReport(uint256 indexed id);
    event EventsMock_KeeperRegistryLogic2_0_OwnerFundsWithdrawn(uint96 amount);
    event EventsMock_KeeperRegistryLogic2_0_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic2_0_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic2_0_Paused(address account);
    event EventsMock_KeeperRegistryLogic2_0_PayeesUpdated(address[] transmitters,address[] payees);
    event EventsMock_KeeperRegistryLogic2_0_PayeeshipTransferRequested(address indexed transmitter,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic2_0_PayeeshipTransferred(address indexed transmitter,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic2_0_PaymentWithdrawn(address indexed transmitter,uint256 indexed amount,address indexed to,address payee);
    event EventsMock_KeeperRegistryLogic2_0_ReorgedUpkeepReport(uint256 indexed id);
    event EventsMock_KeeperRegistryLogic2_0_StaleUpkeepReport(uint256 indexed id);
    event EventsMock_KeeperRegistryLogic2_0_Unpaused(address account);
    event EventsMock_KeeperRegistryLogic2_0_UpkeepAdminTransferRequested(uint256 indexed id,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic2_0_UpkeepAdminTransferred(uint256 indexed id,address indexed from,address indexed to);
    event EventsMock_KeeperRegistryLogic2_0_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event EventsMock_KeeperRegistryLogic2_0_UpkeepCheckDataUpdated(uint256 indexed id,bytes newCheckData);
    event EventsMock_KeeperRegistryLogic2_0_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event EventsMock_KeeperRegistryLogic2_0_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event EventsMock_KeeperRegistryLogic2_0_UpkeepOffchainConfigSet(uint256 indexed id,bytes offchainConfig);
    event EventsMock_KeeperRegistryLogic2_0_UpkeepPaused(uint256 indexed id);
    event EventsMock_KeeperRegistryLogic2_0_UpkeepPerformed(uint256 indexed id,bool indexed success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment);
    event EventsMock_KeeperRegistryLogic2_0_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event EventsMock_KeeperRegistryLogic2_0_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event EventsMock_KeeperRegistryLogic2_0_UpkeepUnpaused(uint256 indexed id);
    event EventsMock_LogEmitter_Log1(uint256 ibdrbwbi);
    event EventsMock_LogEmitter_Log2(uint256 indexed sopulnfk);
    event EventsMock_LogEmitter_Log3(string sqdzywth);
    event EventsMock_OCR2Abstract_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] signers,address[] transmitters,uint8 f,bytes onchainConfig,uint64 offchainConfigVersion,bytes offchainConfig);
    event EventsMock_OCR2Abstract_Transmitted(bytes32 configDigest,uint32 epoch);
    event EventsMock_OCR2BaseUpgradeable_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] signers,address[] transmitters,uint8 f,bytes onchainConfig,uint64 offchainConfigVersion,bytes offchainConfig);
    event EventsMock_OCR2BaseUpgradeable_Initialized(uint8 version);
    event EventsMock_OCR2BaseUpgradeable_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_OCR2BaseUpgradeable_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_OCR2BaseUpgradeable_Transmitted(bytes32 configDigest,uint32 epoch);
    event EventsMock_OVM_GasPriceOracle_DecimalsUpdated(uint256 wwbbrgcs);
    event EventsMock_OVM_GasPriceOracle_GasPriceUpdated(uint256 bzwkuldh);
    event EventsMock_OVM_GasPriceOracle_L1BaseFeeUpdated(uint256 bwjknner);
    event EventsMock_OVM_GasPriceOracle_OverheadUpdated(uint256 riulbpbh);
    event EventsMock_OVM_GasPriceOracle_OwnershipTransferred(address indexed previousOwner,address indexed newOwner);
    event EventsMock_OVM_GasPriceOracle_ScalarUpdated(uint256 whsnlsbn);
    event EventsMock_Ownable_OwnershipTransferred(address indexed previousOwner,address indexed newOwner);
    event EventsMock_PausableUpgradeable_Initialized(uint8 version);
    event EventsMock_PausableUpgradeable_Paused(address account);
    event EventsMock_PausableUpgradeable_Unpaused(address account);
    event EventsMock_Pausable_Paused(address account);
    event EventsMock_Pausable_Unpaused(address account);
    event EventsMock_ProxyAdmin_OwnershipTransferred(address indexed previousOwner,address indexed newOwner);
    event EventsMock_TransparentUpgradeableProxy_AdminChanged(address previousAdmin,address newAdmin);
    event EventsMock_TransparentUpgradeableProxy_BeaconUpgraded(address indexed beacon);
    event EventsMock_TransparentUpgradeableProxy_Upgraded(address indexed implementation);
    event EventsMock_VRFConsumerBaseV2Upgradeable_Initialized(uint8 version);
    event EventsMock_VRFConsumerV2UpgradeableExample_Initialized(uint8 version);
    event EventsMock_VRFCoordinatorMock_RandomnessRequest(address indexed sender,bytes32 indexed keyHash,uint256 indexed seed);
    event EventsMock_VRFCoordinatorV2TestHelper_ConfigSet(uint16 minimumRequestConfirmations,uint32 maxGasLimit,uint32 stalenessSeconds,uint32 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,EventsMock_VRFCoordinatorV2TestHelper_FeeConfig feeConfig);
    event EventsMock_VRFCoordinatorV2TestHelper_FundsRecovered(address to,uint256 amount);
    event EventsMock_VRFCoordinatorV2TestHelper_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_VRFCoordinatorV2TestHelper_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_VRFCoordinatorV2TestHelper_ProvingKeyDeregistered(bytes32 keyHash,address indexed oracle);
    event EventsMock_VRFCoordinatorV2TestHelper_ProvingKeyRegistered(bytes32 keyHash,address indexed oracle);
    event EventsMock_VRFCoordinatorV2TestHelper_RandomWordsFulfilled(uint256 indexed requestId,uint256 outputSeed,uint96 payment,bool success);
    event EventsMock_VRFCoordinatorV2TestHelper_RandomWordsRequested(bytes32 indexed keyHash,uint256 requestId,uint256 preSeed,uint64 indexed subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address indexed sender);
    event EventsMock_VRFCoordinatorV2TestHelper_SubscriptionCanceled(uint64 indexed subId,address to,uint256 amount);
    event EventsMock_VRFCoordinatorV2TestHelper_SubscriptionConsumerAdded(uint64 indexed subId,address consumer);
    event EventsMock_VRFCoordinatorV2TestHelper_SubscriptionConsumerRemoved(uint64 indexed subId,address consumer);
    event EventsMock_VRFCoordinatorV2TestHelper_SubscriptionCreated(uint64 indexed subId,address owner);
    event EventsMock_VRFCoordinatorV2TestHelper_SubscriptionFunded(uint64 indexed subId,uint256 oldBalance,uint256 newBalance);
    event EventsMock_VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferRequested(uint64 indexed subId,address from,address to);
    event EventsMock_VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferred(uint64 indexed subId,address from,address to);
    event EventsMock_VRFCoordinatorV2_ConfigSet(uint16 minimumRequestConfirmations,uint32 maxGasLimit,uint32 stalenessSeconds,uint32 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,EventsMock_VRFCoordinatorV2_FeeConfig feeConfig);
    event EventsMock_VRFCoordinatorV2_FundsRecovered(address to,uint256 amount);
    event EventsMock_VRFCoordinatorV2_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_VRFCoordinatorV2_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_VRFCoordinatorV2_ProvingKeyDeregistered(bytes32 keyHash,address indexed oracle);
    event EventsMock_VRFCoordinatorV2_ProvingKeyRegistered(bytes32 keyHash,address indexed oracle);
    event EventsMock_VRFCoordinatorV2_RandomWordsFulfilled(uint256 indexed requestId,uint256 outputSeed,uint96 payment,bool success);
    event EventsMock_VRFCoordinatorV2_RandomWordsRequested(bytes32 indexed keyHash,uint256 requestId,uint256 preSeed,uint64 indexed subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address indexed sender);
    event EventsMock_VRFCoordinatorV2_SubscriptionCanceled(uint64 indexed subId,address to,uint256 amount);
    event EventsMock_VRFCoordinatorV2_SubscriptionConsumerAdded(uint64 indexed subId,address consumer);
    event EventsMock_VRFCoordinatorV2_SubscriptionConsumerRemoved(uint64 indexed subId,address consumer);
    event EventsMock_VRFCoordinatorV2_SubscriptionCreated(uint64 indexed subId,address owner);
    event EventsMock_VRFCoordinatorV2_SubscriptionFunded(uint64 indexed subId,uint256 oldBalance,uint256 newBalance);
    event EventsMock_VRFCoordinatorV2_SubscriptionOwnerTransferRequested(uint64 indexed subId,address from,address to);
    event EventsMock_VRFCoordinatorV2_SubscriptionOwnerTransferred(uint64 indexed subId,address from,address to);
    event EventsMock_VRFLoadTestExternalSubOwner_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_VRFLoadTestExternalSubOwner_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_VRFV2ProxyAdmin_OwnershipTransferred(address indexed previousOwner,address indexed newOwner);
    event EventsMock_VRFV2TransparentUpgradeableProxy_AdminChanged(address previousAdmin,address newAdmin);
    event EventsMock_VRFV2TransparentUpgradeableProxy_BeaconUpgraded(address indexed beacon);
    event EventsMock_VRFV2TransparentUpgradeableProxy_Upgraded(address indexed implementation);
    event EventsMock_VRFV2WrapperConsumerExample_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_VRFV2WrapperConsumerExample_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_VRFV2WrapperConsumerExample_WrappedRequestFulfilled(uint256 requestId,uint256[] randomWords,uint256 payment);
    event EventsMock_VRFV2WrapperConsumerExample_WrapperRequestMade(uint256 indexed requestId,uint256 paid);
    event EventsMock_VRFV2Wrapper_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_VRFV2Wrapper_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_VRFV2Wrapper_WrapperFulfillmentFailed(uint256 indexed requestId,address indexed consumer);
    event EventsMock_VerifierProxy_AccessControllerSet(address oldAccessController,address newAccessController);
    event EventsMock_VerifierProxy_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_VerifierProxy_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_VerifierProxy_VerifierSet(bytes32 oldConfigDigest,bytes32 newConfigDigest,address verifierAddress);
    event EventsMock_VerifierProxy_VerifierUnset(bytes32 configDigest,address verifierAddress);
    event EventsMock_Verifier_ConfigActivated(bytes32 indexed feedId,bytes32 configDigest);
    event EventsMock_Verifier_ConfigDeactivated(bytes32 indexed feedId,bytes32 configDigest);
    event EventsMock_Verifier_ConfigSet(bytes32 indexed feedId,uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] signers,bytes32[] offchainTransmitters,uint8 f,bytes onchainConfig,uint64 offchainConfigVersion,bytes offchainConfig);
    event EventsMock_Verifier_OwnershipTransferRequested(address indexed from,address indexed to);
    event EventsMock_Verifier_OwnershipTransferred(address indexed from,address indexed to);
    event EventsMock_Verifier_ReportVerified(bytes32 indexed feedId,bytes32 reportHash,address requester);
    event FunctionsBillingRegistry_AuthorizedSendersChanged(address[] senders,address changedBy);
    event FunctionsBillingRegistry_BillingEnd(bytes32 indexed requestId,uint64 subscriptionId,uint96 signerPayment,uint96 transmitterPayment,uint96 totalCost,bool success);
    event FunctionsBillingRegistry_BillingStart(bytes32 indexed requestId,FunctionsBillingRegistry_Commitment commitment);
    event FunctionsBillingRegistry_ConfigSet(uint32 maxGasLimit,uint32 stalenessSeconds,uint256 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,uint32 gasOverhead);
    event FunctionsBillingRegistry_FundsRecovered(address to,uint256 amount);
    event FunctionsBillingRegistry_Initialized(uint8 version);
    event FunctionsBillingRegistry_OwnershipTransferRequested(address indexed from,address indexed to);
    event FunctionsBillingRegistry_OwnershipTransferred(address indexed from,address indexed to);
    event FunctionsBillingRegistry_Paused(address account);
    event FunctionsBillingRegistry_RequestTimedOut(bytes32 indexed requestId);
    event FunctionsBillingRegistry_SubscriptionCanceled(uint64 indexed subscriptionId,address to,uint256 amount);
    event FunctionsBillingRegistry_SubscriptionConsumerAdded(uint64 indexed subscriptionId,address consumer);
    event FunctionsBillingRegistry_SubscriptionConsumerRemoved(uint64 indexed subscriptionId,address consumer);
    event FunctionsBillingRegistry_SubscriptionCreated(uint64 indexed subscriptionId,address owner);
    event FunctionsBillingRegistry_SubscriptionFunded(uint64 indexed subscriptionId,uint256 oldBalance,uint256 newBalance);
    event FunctionsBillingRegistry_SubscriptionOwnerTransferRequested(uint64 indexed subscriptionId,address from,address to);
    event FunctionsBillingRegistry_SubscriptionOwnerTransferred(uint64 indexed subscriptionId,address from,address to);
    event FunctionsBillingRegistry_Unpaused(address account);
    event FunctionsBillingRegistryWithInit_AuthorizedSendersChanged(address[] senders,address changedBy);
    event FunctionsBillingRegistryWithInit_BillingEnd(bytes32 indexed requestId,uint64 subscriptionId,uint96 signerPayment,uint96 transmitterPayment,uint96 totalCost,bool success);
    event FunctionsBillingRegistryWithInit_BillingStart(bytes32 indexed requestId,FunctionsBillingRegistryWithInit_Commitment commitment);
    event FunctionsBillingRegistryWithInit_ConfigSet(uint32 maxGasLimit,uint32 stalenessSeconds,uint256 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,uint32 gasOverhead);
    event FunctionsBillingRegistryWithInit_FundsRecovered(address to,uint256 amount);
    event FunctionsBillingRegistryWithInit_Initialized(uint8 version);
    event FunctionsBillingRegistryWithInit_OwnershipTransferRequested(address indexed from,address indexed to);
    event FunctionsBillingRegistryWithInit_OwnershipTransferred(address indexed from,address indexed to);
    event FunctionsBillingRegistryWithInit_Paused(address account);
    event FunctionsBillingRegistryWithInit_RequestTimedOut(bytes32 indexed requestId);
    event FunctionsBillingRegistryWithInit_SubscriptionCanceled(uint64 indexed subscriptionId,address to,uint256 amount);
    event FunctionsBillingRegistryWithInit_SubscriptionConsumerAdded(uint64 indexed subscriptionId,address consumer);
    event FunctionsBillingRegistryWithInit_SubscriptionConsumerRemoved(uint64 indexed subscriptionId,address consumer);
    event FunctionsBillingRegistryWithInit_SubscriptionCreated(uint64 indexed subscriptionId,address owner);
    event FunctionsBillingRegistryWithInit_SubscriptionFunded(uint64 indexed subscriptionId,uint256 oldBalance,uint256 newBalance);
    event FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferRequested(uint64 indexed subscriptionId,address from,address to);
    event FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferred(uint64 indexed subscriptionId,address from,address to);
    event FunctionsBillingRegistryWithInit_Unpaused(address account);
    event FunctionsClient_RequestFulfilled(bytes32 indexed id);
    event FunctionsClient_RequestSent(bytes32 indexed id);
    event FunctionsClientExample_OwnershipTransferRequested(address indexed from,address indexed to);
    event FunctionsClientExample_OwnershipTransferred(address indexed from,address indexed to);
    event FunctionsClientExample_RequestFulfilled(bytes32 indexed id);
    event FunctionsClientExample_RequestSent(bytes32 indexed id);
    event FunctionsOracle_AuthorizedSendersActive(address account);
    event FunctionsOracle_AuthorizedSendersChanged(address[] senders,address changedBy);
    event FunctionsOracle_AuthorizedSendersDeactive(address account);
    event FunctionsOracle_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] signers,address[] transmitters,uint8 f,bytes onchainConfig,uint64 offchainConfigVersion,bytes offchainConfig);
    event FunctionsOracle_Initialized(uint8 version);
    event FunctionsOracle_InvalidRequestID(bytes32 indexed requestId);
    event FunctionsOracle_OracleRequest(bytes32 indexed requestId,address requestingContract,address requestInitiator,uint64 subscriptionId,address subscriptionOwner,bytes data);
    event FunctionsOracle_OracleResponse(bytes32 indexed requestId);
    event FunctionsOracle_OwnershipTransferRequested(address indexed from,address indexed to);
    event FunctionsOracle_OwnershipTransferred(address indexed from,address indexed to);
    event FunctionsOracle_Transmitted(bytes32 configDigest,uint32 epoch);
    event FunctionsOracle_UserCallbackError(bytes32 indexed requestId,string reason);
    event FunctionsOracle_UserCallbackRawError(bytes32 indexed requestId,bytes lowLevelData);
    event FunctionsOracleWithInit_AuthorizedSendersActive(address account);
    event FunctionsOracleWithInit_AuthorizedSendersChanged(address[] senders,address changedBy);
    event FunctionsOracleWithInit_AuthorizedSendersDeactive(address account);
    event FunctionsOracleWithInit_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] signers,address[] transmitters,uint8 f,bytes onchainConfig,uint64 offchainConfigVersion,bytes offchainConfig);
    event FunctionsOracleWithInit_Initialized(uint8 version);
    event FunctionsOracleWithInit_InvalidRequestID(bytes32 indexed requestId);
    event FunctionsOracleWithInit_OracleRequest(bytes32 indexed requestId,address requestingContract,address requestInitiator,uint64 subscriptionId,address subscriptionOwner,bytes data);
    event FunctionsOracleWithInit_OracleResponse(bytes32 indexed requestId);
    event FunctionsOracleWithInit_OwnershipTransferRequested(address indexed from,address indexed to);
    event FunctionsOracleWithInit_OwnershipTransferred(address indexed from,address indexed to);
    event FunctionsOracleWithInit_Transmitted(bytes32 configDigest,uint32 epoch);
    event FunctionsOracleWithInit_UserCallbackError(bytes32 indexed requestId,string reason);
    event FunctionsOracleWithInit_UserCallbackRawError(bytes32 indexed requestId,bytes lowLevelData);
    event Initializable_Initialized(uint8 version);
    event KeeperRegistry1_2_ConfigSet(KeeperRegistry1_2_Config config);
    event KeeperRegistry1_2_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event KeeperRegistry1_2_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event KeeperRegistry1_2_KeepersUpdated(address[] keepers,address[] payees);
    event KeeperRegistry1_2_OwnerFundsWithdrawn(uint96 amount);
    event KeeperRegistry1_2_OwnershipTransferRequested(address indexed from,address indexed to);
    event KeeperRegistry1_2_OwnershipTransferred(address indexed from,address indexed to);
    event KeeperRegistry1_2_Paused(address account);
    event KeeperRegistry1_2_PayeeshipTransferRequested(address indexed keeper,address indexed from,address indexed to);
    event KeeperRegistry1_2_PayeeshipTransferred(address indexed keeper,address indexed from,address indexed to);
    event KeeperRegistry1_2_PaymentWithdrawn(address indexed keeper,uint256 indexed amount,address indexed to,address payee);
    event KeeperRegistry1_2_Unpaused(address account);
    event KeeperRegistry1_2_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event KeeperRegistry1_2_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event KeeperRegistry1_2_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event KeeperRegistry1_2_UpkeepPerformed(uint256 indexed id,bool indexed success,address indexed from,uint96 payment,bytes performData);
    event KeeperRegistry1_2_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event KeeperRegistry1_2_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event KeeperRegistry1_3_ConfigSet(KeeperRegistry1_3_Config config);
    event KeeperRegistry1_3_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event KeeperRegistry1_3_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event KeeperRegistry1_3_KeepersUpdated(address[] keepers,address[] payees);
    event KeeperRegistry1_3_OwnerFundsWithdrawn(uint96 amount);
    event KeeperRegistry1_3_OwnershipTransferRequested(address indexed from,address indexed to);
    event KeeperRegistry1_3_OwnershipTransferred(address indexed from,address indexed to);
    event KeeperRegistry1_3_Paused(address account);
    event KeeperRegistry1_3_PayeeshipTransferRequested(address indexed keeper,address indexed from,address indexed to);
    event KeeperRegistry1_3_PayeeshipTransferred(address indexed keeper,address indexed from,address indexed to);
    event KeeperRegistry1_3_PaymentWithdrawn(address indexed keeper,uint256 indexed amount,address indexed to,address payee);
    event KeeperRegistry1_3_Unpaused(address account);
    event KeeperRegistry1_3_UpkeepAdminTransferRequested(uint256 indexed id,address indexed from,address indexed to);
    event KeeperRegistry1_3_UpkeepAdminTransferred(uint256 indexed id,address indexed from,address indexed to);
    event KeeperRegistry1_3_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event KeeperRegistry1_3_UpkeepCheckDataUpdated(uint256 indexed id,bytes newCheckData);
    event KeeperRegistry1_3_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event KeeperRegistry1_3_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event KeeperRegistry1_3_UpkeepPaused(uint256 indexed id);
    event KeeperRegistry1_3_UpkeepPerformed(uint256 indexed id,bool indexed success,address indexed from,uint96 payment,bytes performData);
    event KeeperRegistry1_3_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event KeeperRegistry1_3_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event KeeperRegistry1_3_UpkeepUnpaused(uint256 indexed id);
    event KeeperRegistry2_0_CancelledUpkeepReport(uint256 indexed id);
    event KeeperRegistry2_0_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] signers,address[] transmitters,uint8 f,bytes onchainConfig,uint64 offchainConfigVersion,bytes offchainConfig);
    event KeeperRegistry2_0_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event KeeperRegistry2_0_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event KeeperRegistry2_0_InsufficientFundsUpkeepReport(uint256 indexed id);
    event KeeperRegistry2_0_OwnerFundsWithdrawn(uint96 amount);
    event KeeperRegistry2_0_OwnershipTransferRequested(address indexed from,address indexed to);
    event KeeperRegistry2_0_OwnershipTransferred(address indexed from,address indexed to);
    event KeeperRegistry2_0_Paused(address account);
    event KeeperRegistry2_0_PayeesUpdated(address[] transmitters,address[] payees);
    event KeeperRegistry2_0_PayeeshipTransferRequested(address indexed transmitter,address indexed from,address indexed to);
    event KeeperRegistry2_0_PayeeshipTransferred(address indexed transmitter,address indexed from,address indexed to);
    event KeeperRegistry2_0_PaymentWithdrawn(address indexed transmitter,uint256 indexed amount,address indexed to,address payee);
    event KeeperRegistry2_0_ReorgedUpkeepReport(uint256 indexed id);
    event KeeperRegistry2_0_StaleUpkeepReport(uint256 indexed id);
    event KeeperRegistry2_0_Transmitted(bytes32 configDigest,uint32 epoch);
    event KeeperRegistry2_0_Unpaused(address account);
    event KeeperRegistry2_0_UpkeepAdminTransferRequested(uint256 indexed id,address indexed from,address indexed to);
    event KeeperRegistry2_0_UpkeepAdminTransferred(uint256 indexed id,address indexed from,address indexed to);
    event KeeperRegistry2_0_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event KeeperRegistry2_0_UpkeepCheckDataUpdated(uint256 indexed id,bytes newCheckData);
    event KeeperRegistry2_0_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event KeeperRegistry2_0_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event KeeperRegistry2_0_UpkeepOffchainConfigSet(uint256 indexed id,bytes offchainConfig);
    event KeeperRegistry2_0_UpkeepPaused(uint256 indexed id);
    event KeeperRegistry2_0_UpkeepPerformed(uint256 indexed id,bool indexed success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment);
    event KeeperRegistry2_0_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event KeeperRegistry2_0_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event KeeperRegistry2_0_UpkeepUnpaused(uint256 indexed id);
    event KeeperRegistryBase1_3_ConfigSet(KeeperRegistryBase1_3_Config config);
    event KeeperRegistryBase1_3_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event KeeperRegistryBase1_3_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event KeeperRegistryBase1_3_KeepersUpdated(address[] keepers,address[] payees);
    event KeeperRegistryBase1_3_OwnerFundsWithdrawn(uint96 amount);
    event KeeperRegistryBase1_3_OwnershipTransferRequested(address indexed from,address indexed to);
    event KeeperRegistryBase1_3_OwnershipTransferred(address indexed from,address indexed to);
    event KeeperRegistryBase1_3_Paused(address account);
    event KeeperRegistryBase1_3_PayeeshipTransferRequested(address indexed keeper,address indexed from,address indexed to);
    event KeeperRegistryBase1_3_PayeeshipTransferred(address indexed keeper,address indexed from,address indexed to);
    event KeeperRegistryBase1_3_PaymentWithdrawn(address indexed keeper,uint256 indexed amount,address indexed to,address payee);
    event KeeperRegistryBase1_3_Unpaused(address account);
    event KeeperRegistryBase1_3_UpkeepAdminTransferRequested(uint256 indexed id,address indexed from,address indexed to);
    event KeeperRegistryBase1_3_UpkeepAdminTransferred(uint256 indexed id,address indexed from,address indexed to);
    event KeeperRegistryBase1_3_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event KeeperRegistryBase1_3_UpkeepCheckDataUpdated(uint256 indexed id,bytes newCheckData);
    event KeeperRegistryBase1_3_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event KeeperRegistryBase1_3_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event KeeperRegistryBase1_3_UpkeepPaused(uint256 indexed id);
    event KeeperRegistryBase1_3_UpkeepPerformed(uint256 indexed id,bool indexed success,address indexed from,uint96 payment,bytes performData);
    event KeeperRegistryBase1_3_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event KeeperRegistryBase1_3_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event KeeperRegistryBase1_3_UpkeepUnpaused(uint256 indexed id);
    event KeeperRegistryBase2_0_CancelledUpkeepReport(uint256 indexed id);
    event KeeperRegistryBase2_0_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event KeeperRegistryBase2_0_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event KeeperRegistryBase2_0_InsufficientFundsUpkeepReport(uint256 indexed id);
    event KeeperRegistryBase2_0_OwnerFundsWithdrawn(uint96 amount);
    event KeeperRegistryBase2_0_OwnershipTransferRequested(address indexed from,address indexed to);
    event KeeperRegistryBase2_0_OwnershipTransferred(address indexed from,address indexed to);
    event KeeperRegistryBase2_0_Paused(address account);
    event KeeperRegistryBase2_0_PayeesUpdated(address[] transmitters,address[] payees);
    event KeeperRegistryBase2_0_PayeeshipTransferRequested(address indexed transmitter,address indexed from,address indexed to);
    event KeeperRegistryBase2_0_PayeeshipTransferred(address indexed transmitter,address indexed from,address indexed to);
    event KeeperRegistryBase2_0_PaymentWithdrawn(address indexed transmitter,uint256 indexed amount,address indexed to,address payee);
    event KeeperRegistryBase2_0_ReorgedUpkeepReport(uint256 indexed id);
    event KeeperRegistryBase2_0_StaleUpkeepReport(uint256 indexed id);
    event KeeperRegistryBase2_0_Unpaused(address account);
    event KeeperRegistryBase2_0_UpkeepAdminTransferRequested(uint256 indexed id,address indexed from,address indexed to);
    event KeeperRegistryBase2_0_UpkeepAdminTransferred(uint256 indexed id,address indexed from,address indexed to);
    event KeeperRegistryBase2_0_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event KeeperRegistryBase2_0_UpkeepCheckDataUpdated(uint256 indexed id,bytes newCheckData);
    event KeeperRegistryBase2_0_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event KeeperRegistryBase2_0_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event KeeperRegistryBase2_0_UpkeepOffchainConfigSet(uint256 indexed id,bytes offchainConfig);
    event KeeperRegistryBase2_0_UpkeepPaused(uint256 indexed id);
    event KeeperRegistryBase2_0_UpkeepPerformed(uint256 indexed id,bool indexed success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment);
    event KeeperRegistryBase2_0_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event KeeperRegistryBase2_0_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event KeeperRegistryBase2_0_UpkeepUnpaused(uint256 indexed id);
    event KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferRequested(address indexed from,address indexed to);
    event KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferred(address indexed from,address indexed to);
    event KeeperRegistryLogic1_3_ConfigSet(KeeperRegistryLogic1_3_Config config);
    event KeeperRegistryLogic1_3_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event KeeperRegistryLogic1_3_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event KeeperRegistryLogic1_3_KeepersUpdated(address[] keepers,address[] payees);
    event KeeperRegistryLogic1_3_OwnerFundsWithdrawn(uint96 amount);
    event KeeperRegistryLogic1_3_OwnershipTransferRequested(address indexed from,address indexed to);
    event KeeperRegistryLogic1_3_OwnershipTransferred(address indexed from,address indexed to);
    event KeeperRegistryLogic1_3_Paused(address account);
    event KeeperRegistryLogic1_3_PayeeshipTransferRequested(address indexed keeper,address indexed from,address indexed to);
    event KeeperRegistryLogic1_3_PayeeshipTransferred(address indexed keeper,address indexed from,address indexed to);
    event KeeperRegistryLogic1_3_PaymentWithdrawn(address indexed keeper,uint256 indexed amount,address indexed to,address payee);
    event KeeperRegistryLogic1_3_Unpaused(address account);
    event KeeperRegistryLogic1_3_UpkeepAdminTransferRequested(uint256 indexed id,address indexed from,address indexed to);
    event KeeperRegistryLogic1_3_UpkeepAdminTransferred(uint256 indexed id,address indexed from,address indexed to);
    event KeeperRegistryLogic1_3_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event KeeperRegistryLogic1_3_UpkeepCheckDataUpdated(uint256 indexed id,bytes newCheckData);
    event KeeperRegistryLogic1_3_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event KeeperRegistryLogic1_3_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event KeeperRegistryLogic1_3_UpkeepPaused(uint256 indexed id);
    event KeeperRegistryLogic1_3_UpkeepPerformed(uint256 indexed id,bool indexed success,address indexed from,uint96 payment,bytes performData);
    event KeeperRegistryLogic1_3_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event KeeperRegistryLogic1_3_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event KeeperRegistryLogic1_3_UpkeepUnpaused(uint256 indexed id);
    event KeeperRegistryLogic2_0_CancelledUpkeepReport(uint256 indexed id);
    event KeeperRegistryLogic2_0_FundsAdded(uint256 indexed id,address indexed from,uint96 amount);
    event KeeperRegistryLogic2_0_FundsWithdrawn(uint256 indexed id,uint256 amount,address to);
    event KeeperRegistryLogic2_0_InsufficientFundsUpkeepReport(uint256 indexed id);
    event KeeperRegistryLogic2_0_OwnerFundsWithdrawn(uint96 amount);
    event KeeperRegistryLogic2_0_OwnershipTransferRequested(address indexed from,address indexed to);
    event KeeperRegistryLogic2_0_OwnershipTransferred(address indexed from,address indexed to);
    event KeeperRegistryLogic2_0_Paused(address account);
    event KeeperRegistryLogic2_0_PayeesUpdated(address[] transmitters,address[] payees);
    event KeeperRegistryLogic2_0_PayeeshipTransferRequested(address indexed transmitter,address indexed from,address indexed to);
    event KeeperRegistryLogic2_0_PayeeshipTransferred(address indexed transmitter,address indexed from,address indexed to);
    event KeeperRegistryLogic2_0_PaymentWithdrawn(address indexed transmitter,uint256 indexed amount,address indexed to,address payee);
    event KeeperRegistryLogic2_0_ReorgedUpkeepReport(uint256 indexed id);
    event KeeperRegistryLogic2_0_StaleUpkeepReport(uint256 indexed id);
    event KeeperRegistryLogic2_0_Unpaused(address account);
    event KeeperRegistryLogic2_0_UpkeepAdminTransferRequested(uint256 indexed id,address indexed from,address indexed to);
    event KeeperRegistryLogic2_0_UpkeepAdminTransferred(uint256 indexed id,address indexed from,address indexed to);
    event KeeperRegistryLogic2_0_UpkeepCanceled(uint256 indexed id,uint64 indexed atBlockHeight);
    event KeeperRegistryLogic2_0_UpkeepCheckDataUpdated(uint256 indexed id,bytes newCheckData);
    event KeeperRegistryLogic2_0_UpkeepGasLimitSet(uint256 indexed id,uint96 gasLimit);
    event KeeperRegistryLogic2_0_UpkeepMigrated(uint256 indexed id,uint256 remainingBalance,address destination);
    event KeeperRegistryLogic2_0_UpkeepOffchainConfigSet(uint256 indexed id,bytes offchainConfig);
    event KeeperRegistryLogic2_0_UpkeepPaused(uint256 indexed id);
    event KeeperRegistryLogic2_0_UpkeepPerformed(uint256 indexed id,bool indexed success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment);
    event KeeperRegistryLogic2_0_UpkeepReceived(uint256 indexed id,uint256 startingBalance,address importedFrom);
    event KeeperRegistryLogic2_0_UpkeepRegistered(uint256 indexed id,uint32 executeGas,address admin);
    event KeeperRegistryLogic2_0_UpkeepUnpaused(uint256 indexed id);
    event LogEmitter_Log1(uint256 jxztvtdu);
    event LogEmitter_Log2(uint256 indexed jbysnosu);
    event LogEmitter_Log3(string njjusihh);
    event OCR2Abstract_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] signers,address[] transmitters,uint8 f,bytes onchainConfig,uint64 offchainConfigVersion,bytes offchainConfig);
    event OCR2Abstract_Transmitted(bytes32 configDigest,uint32 epoch);
    event OCR2BaseUpgradeable_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] signers,address[] transmitters,uint8 f,bytes onchainConfig,uint64 offchainConfigVersion,bytes offchainConfig);
    event OCR2BaseUpgradeable_Initialized(uint8 version);
    event OCR2BaseUpgradeable_OwnershipTransferRequested(address indexed from,address indexed to);
    event OCR2BaseUpgradeable_OwnershipTransferred(address indexed from,address indexed to);
    event OCR2BaseUpgradeable_Transmitted(bytes32 configDigest,uint32 epoch);
    event OVM_GasPriceOracle_DecimalsUpdated(uint256 ofmpqceu);
    event OVM_GasPriceOracle_GasPriceUpdated(uint256 lflivotg);
    event OVM_GasPriceOracle_L1BaseFeeUpdated(uint256 ohwqxozp);
    event OVM_GasPriceOracle_OverheadUpdated(uint256 jvbubple);
    event OVM_GasPriceOracle_OwnershipTransferred(address indexed previousOwner,address indexed newOwner);
    event OVM_GasPriceOracle_ScalarUpdated(uint256 hdoeyfyg);
    event Ownable_OwnershipTransferred(address indexed previousOwner,address indexed newOwner);
    event Pausable_Paused(address account);
    event Pausable_Unpaused(address account);
    event PausableUpgradeable_Initialized(uint8 version);
    event PausableUpgradeable_Paused(address account);
    event PausableUpgradeable_Unpaused(address account);
    event ProxyAdmin_OwnershipTransferred(address indexed previousOwner,address indexed newOwner);
    event TransparentUpgradeableProxy_AdminChanged(address previousAdmin,address newAdmin);
    event TransparentUpgradeableProxy_BeaconUpgraded(address indexed beacon);
    event TransparentUpgradeableProxy_Upgraded(address indexed implementation);
    event VRFConsumerBaseV2Upgradeable_Initialized(uint8 version);
    event VRFConsumerV2UpgradeableExample_Initialized(uint8 version);
    event VRFCoordinatorMock_RandomnessRequest(address indexed sender,bytes32 indexed keyHash,uint256 indexed seed);
    event VRFCoordinatorV2_ConfigSet(uint16 minimumRequestConfirmations,uint32 maxGasLimit,uint32 stalenessSeconds,uint32 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,VRFCoordinatorV2_FeeConfig feeConfig);
    event VRFCoordinatorV2_FundsRecovered(address to,uint256 amount);
    event VRFCoordinatorV2_OwnershipTransferRequested(address indexed from,address indexed to);
    event VRFCoordinatorV2_OwnershipTransferred(address indexed from,address indexed to);
    event VRFCoordinatorV2_ProvingKeyDeregistered(bytes32 keyHash,address indexed oracle);
    event VRFCoordinatorV2_ProvingKeyRegistered(bytes32 keyHash,address indexed oracle);
    event VRFCoordinatorV2_RandomWordsFulfilled(uint256 indexed requestId,uint256 outputSeed,uint96 payment,bool success);
    event VRFCoordinatorV2_RandomWordsRequested(bytes32 indexed keyHash,uint256 requestId,uint256 preSeed,uint64 indexed subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address indexed sender);
    event VRFCoordinatorV2_SubscriptionCanceled(uint64 indexed subId,address to,uint256 amount);
    event VRFCoordinatorV2_SubscriptionConsumerAdded(uint64 indexed subId,address consumer);
    event VRFCoordinatorV2_SubscriptionConsumerRemoved(uint64 indexed subId,address consumer);
    event VRFCoordinatorV2_SubscriptionCreated(uint64 indexed subId,address owner);
    event VRFCoordinatorV2_SubscriptionFunded(uint64 indexed subId,uint256 oldBalance,uint256 newBalance);
    event VRFCoordinatorV2_SubscriptionOwnerTransferRequested(uint64 indexed subId,address from,address to);
    event VRFCoordinatorV2_SubscriptionOwnerTransferred(uint64 indexed subId,address from,address to);
    event VRFCoordinatorV2TestHelper_ConfigSet(uint16 minimumRequestConfirmations,uint32 maxGasLimit,uint32 stalenessSeconds,uint32 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,VRFCoordinatorV2TestHelper_FeeConfig feeConfig);
    event VRFCoordinatorV2TestHelper_FundsRecovered(address to,uint256 amount);
    event VRFCoordinatorV2TestHelper_OwnershipTransferRequested(address indexed from,address indexed to);
    event VRFCoordinatorV2TestHelper_OwnershipTransferred(address indexed from,address indexed to);
    event VRFCoordinatorV2TestHelper_ProvingKeyDeregistered(bytes32 keyHash,address indexed oracle);
    event VRFCoordinatorV2TestHelper_ProvingKeyRegistered(bytes32 keyHash,address indexed oracle);
    event VRFCoordinatorV2TestHelper_RandomWordsFulfilled(uint256 indexed requestId,uint256 outputSeed,uint96 payment,bool success);
    event VRFCoordinatorV2TestHelper_RandomWordsRequested(bytes32 indexed keyHash,uint256 requestId,uint256 preSeed,uint64 indexed subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address indexed sender);
    event VRFCoordinatorV2TestHelper_SubscriptionCanceled(uint64 indexed subId,address to,uint256 amount);
    event VRFCoordinatorV2TestHelper_SubscriptionConsumerAdded(uint64 indexed subId,address consumer);
    event VRFCoordinatorV2TestHelper_SubscriptionConsumerRemoved(uint64 indexed subId,address consumer);
    event VRFCoordinatorV2TestHelper_SubscriptionCreated(uint64 indexed subId,address owner);
    event VRFCoordinatorV2TestHelper_SubscriptionFunded(uint64 indexed subId,uint256 oldBalance,uint256 newBalance);
    event VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferRequested(uint64 indexed subId,address from,address to);
    event VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferred(uint64 indexed subId,address from,address to);
    event VRFLoadTestExternalSubOwner_OwnershipTransferRequested(address indexed from,address indexed to);
    event VRFLoadTestExternalSubOwner_OwnershipTransferred(address indexed from,address indexed to);
    event VRFV2ProxyAdmin_OwnershipTransferred(address indexed previousOwner,address indexed newOwner);
    event VRFV2TransparentUpgradeableProxy_AdminChanged(address previousAdmin,address newAdmin);
    event VRFV2TransparentUpgradeableProxy_BeaconUpgraded(address indexed beacon);
    event VRFV2TransparentUpgradeableProxy_Upgraded(address indexed implementation);
    event VRFV2Wrapper_OwnershipTransferRequested(address indexed from,address indexed to);
    event VRFV2Wrapper_OwnershipTransferred(address indexed from,address indexed to);
    event VRFV2Wrapper_WrapperFulfillmentFailed(uint256 indexed requestId,address indexed consumer);
    event VRFV2WrapperConsumerExample_OwnershipTransferRequested(address indexed from,address indexed to);
    event VRFV2WrapperConsumerExample_OwnershipTransferred(address indexed from,address indexed to);
    event VRFV2WrapperConsumerExample_WrappedRequestFulfilled(uint256 requestId,uint256[] randomWords,uint256 payment);
    event VRFV2WrapperConsumerExample_WrapperRequestMade(uint256 indexed requestId,uint256 paid);
    event Verifier_ConfigActivated(bytes32 indexed feedId,bytes32 configDigest);
    event Verifier_ConfigDeactivated(bytes32 indexed feedId,bytes32 configDigest);
    event Verifier_ConfigSet(bytes32 indexed feedId,uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] signers,bytes32[] offchainTransmitters,uint8 f,bytes onchainConfig,uint64 offchainConfigVersion,bytes offchainConfig);
    event Verifier_OwnershipTransferRequested(address indexed from,address indexed to);
    event Verifier_OwnershipTransferred(address indexed from,address indexed to);
    event Verifier_ReportVerified(bytes32 indexed feedId,bytes32 reportHash,address requester);
    event VerifierProxy_AccessControllerSet(address oldAccessController,address newAccessController);
    event VerifierProxy_OwnershipTransferRequested(address indexed from,address indexed to);
    event VerifierProxy_OwnershipTransferred(address indexed from,address indexed to);
    event VerifierProxy_VerifierSet(bytes32 oldConfigDigest,bytes32 newConfigDigest,address verifierAddress);
    event VerifierProxy_VerifierUnset(bytes32 configDigest,address verifierAddress);
    function emitAggregatorInterface_AnswerUpdated(int256 current,uint256 roundId,uint256 updatedAt) public {
        emit AggregatorInterface_AnswerUpdated(current,roundId,updatedAt);
    }
    function emitAggregatorInterface_NewRound(uint256 roundId,address startedBy,uint256 startedAt) public {
        emit AggregatorInterface_NewRound(roundId,startedBy,startedAt);
    }
    function emitAggregatorV2V3Interface_AnswerUpdated(int256 current,uint256 roundId,uint256 updatedAt) public {
        emit AggregatorV2V3Interface_AnswerUpdated(current,roundId,updatedAt);
    }
    function emitAggregatorV2V3Interface_NewRound(uint256 roundId,address startedBy,uint256 startedAt) public {
        emit AggregatorV2V3Interface_NewRound(roundId,startedBy,startedAt);
    }
    function emitAuthorizedOriginReceiverUpgradeable_AuthorizedSendersActive(address account) public {
        emit AuthorizedOriginReceiverUpgradeable_AuthorizedSendersActive(account);
    }
    function emitAuthorizedOriginReceiverUpgradeable_AuthorizedSendersChanged(address[] memory senders,address changedBy) public {
        emit AuthorizedOriginReceiverUpgradeable_AuthorizedSendersChanged(senders,changedBy);
    }
    function emitAuthorizedOriginReceiverUpgradeable_AuthorizedSendersDeactive(address account) public {
        emit AuthorizedOriginReceiverUpgradeable_AuthorizedSendersDeactive(account);
    }
    function emitAuthorizedOriginReceiverUpgradeable_Initialized(uint8 version) public {
        emit AuthorizedOriginReceiverUpgradeable_Initialized(version);
    }
    function emitAuthorizedReceiver_AuthorizedSendersChanged(address[] memory senders,address changedBy) public {
        emit AuthorizedReceiver_AuthorizedSendersChanged(senders,changedBy);
    }
    function emitBatchVRFCoordinatorV2_ErrorReturned(uint256 requestId,string memory reason) public {
        emit BatchVRFCoordinatorV2_ErrorReturned(requestId,reason);
    }
    function emitBatchVRFCoordinatorV2_RawErrorReturned(uint256 requestId,bytes memory lowLevelData) public {
        emit BatchVRFCoordinatorV2_RawErrorReturned(requestId,lowLevelData);
    }
    function emitChainlinkClient_ChainlinkCancelled(bytes32 id) public {
        emit ChainlinkClient_ChainlinkCancelled(id);
    }
    function emitChainlinkClient_ChainlinkFulfilled(bytes32 id) public {
        emit ChainlinkClient_ChainlinkFulfilled(id);
    }
    function emitChainlinkClient_ChainlinkRequested(bytes32 id) public {
        emit ChainlinkClient_ChainlinkRequested(id);
    }
    function emitConfirmedOwner_OwnershipTransferRequested(address from,address to) public {
        emit ConfirmedOwner_OwnershipTransferRequested(from,to);
    }
    function emitConfirmedOwner_OwnershipTransferred(address from,address to) public {
        emit ConfirmedOwner_OwnershipTransferred(from,to);
    }
    function emitConfirmedOwnerUpgradeable_Initialized(uint8 version) public {
        emit ConfirmedOwnerUpgradeable_Initialized(version);
    }
    function emitConfirmedOwnerUpgradeable_OwnershipTransferRequested(address from,address to) public {
        emit ConfirmedOwnerUpgradeable_OwnershipTransferRequested(from,to);
    }
    function emitConfirmedOwnerUpgradeable_OwnershipTransferred(address from,address to) public {
        emit ConfirmedOwnerUpgradeable_OwnershipTransferred(from,to);
    }
    function emitConfirmedOwnerWithProposal_OwnershipTransferRequested(address from,address to) public {
        emit ConfirmedOwnerWithProposal_OwnershipTransferRequested(from,to);
    }
    function emitConfirmedOwnerWithProposal_OwnershipTransferred(address from,address to) public {
        emit ConfirmedOwnerWithProposal_OwnershipTransferred(from,to);
    }
    function emitContextUpgradeable_Initialized(uint8 version) public {
        emit ContextUpgradeable_Initialized(version);
    }
    function emitCronUpkeep_CronJobCreated(uint256 id,address target,bytes memory handler) public {
        emit CronUpkeep_CronJobCreated(id,target,handler);
    }
    function emitCronUpkeep_CronJobDeleted(uint256 id) public {
        emit CronUpkeep_CronJobDeleted(id);
    }
    function emitCronUpkeep_CronJobExecuted(uint256 id,uint256 timestamp) public {
        emit CronUpkeep_CronJobExecuted(id,timestamp);
    }
    function emitCronUpkeep_CronJobUpdated(uint256 id,address target,bytes memory handler) public {
        emit CronUpkeep_CronJobUpdated(id,target,handler);
    }
    function emitCronUpkeep_OwnershipTransferRequested(address from,address to) public {
        emit CronUpkeep_OwnershipTransferRequested(from,to);
    }
    function emitCronUpkeep_OwnershipTransferred(address from,address to) public {
        emit CronUpkeep_OwnershipTransferred(from,to);
    }
    function emitCronUpkeep_Paused(address account) public {
        emit CronUpkeep_Paused(account);
    }
    function emitCronUpkeep_Unpaused(address account) public {
        emit CronUpkeep_Unpaused(account);
    }
    function emitCronUpkeepFactory_NewCronUpkeepCreated(address upkeep,address owner) public {
        emit CronUpkeepFactory_NewCronUpkeepCreated(upkeep,owner);
    }
    function emitCronUpkeepFactory_OwnershipTransferRequested(address from,address to) public {
        emit CronUpkeepFactory_OwnershipTransferRequested(from,to);
    }
    function emitCronUpkeepFactory_OwnershipTransferred(address from,address to) public {
        emit CronUpkeepFactory_OwnershipTransferred(from,to);
    }
    function emitENSInterface_NewOwner(bytes32 node,bytes32 label,address owner) public {
        emit ENSInterface_NewOwner(node,label,owner);
    }
    function emitENSInterface_NewResolver(bytes32 node,address resolver) public {
        emit ENSInterface_NewResolver(node,resolver);
    }
    function emitENSInterface_NewTTL(bytes32 node,uint64 ttl) public {
        emit ENSInterface_NewTTL(node,ttl);
    }
    function emitENSInterface_Transfer(bytes32 node,address owner) public {
        emit ENSInterface_Transfer(node,owner);
    }
    function emitERC1967Proxy_AdminChanged(address previousAdmin,address newAdmin) public {
        emit ERC1967Proxy_AdminChanged(previousAdmin,newAdmin);
    }
    function emitERC1967Proxy_BeaconUpgraded(address beacon) public {
        emit ERC1967Proxy_BeaconUpgraded(beacon);
    }
    function emitERC1967Proxy_Upgraded(address implementation) public {
        emit ERC1967Proxy_Upgraded(implementation);
    }
    function emitERC1967Upgrade_AdminChanged(address previousAdmin,address newAdmin) public {
        emit ERC1967Upgrade_AdminChanged(previousAdmin,newAdmin);
    }
    function emitERC1967Upgrade_BeaconUpgraded(address beacon) public {
        emit ERC1967Upgrade_BeaconUpgraded(beacon);
    }
    function emitERC1967Upgrade_Upgraded(address implementation) public {
        emit ERC1967Upgrade_Upgraded(implementation);
    }
    function emitEventsMock_AggregatorInterface_AnswerUpdated(int256 current,uint256 roundId,uint256 updatedAt) public {
        emit EventsMock_AggregatorInterface_AnswerUpdated(current,roundId,updatedAt);
    }
    function emitEventsMock_AggregatorInterface_NewRound(uint256 roundId,address startedBy,uint256 startedAt) public {
        emit EventsMock_AggregatorInterface_NewRound(roundId,startedBy,startedAt);
    }
    function emitEventsMock_AggregatorV2V3Interface_AnswerUpdated(int256 current,uint256 roundId,uint256 updatedAt) public {
        emit EventsMock_AggregatorV2V3Interface_AnswerUpdated(current,roundId,updatedAt);
    }
    function emitEventsMock_AggregatorV2V3Interface_NewRound(uint256 roundId,address startedBy,uint256 startedAt) public {
        emit EventsMock_AggregatorV2V3Interface_NewRound(roundId,startedBy,startedAt);
    }
    function emitEventsMock_AuthorizedOriginReceiverUpgradeable_AuthorizedSendersActive(address account) public {
        emit EventsMock_AuthorizedOriginReceiverUpgradeable_AuthorizedSendersActive(account);
    }
    function emitEventsMock_AuthorizedOriginReceiverUpgradeable_AuthorizedSendersChanged(address[] memory senders,address changedBy) public {
        emit EventsMock_AuthorizedOriginReceiverUpgradeable_AuthorizedSendersChanged(senders,changedBy);
    }
    function emitEventsMock_AuthorizedOriginReceiverUpgradeable_AuthorizedSendersDeactive(address account) public {
        emit EventsMock_AuthorizedOriginReceiverUpgradeable_AuthorizedSendersDeactive(account);
    }
    function emitEventsMock_AuthorizedOriginReceiverUpgradeable_Initialized(uint8 version) public {
        emit EventsMock_AuthorizedOriginReceiverUpgradeable_Initialized(version);
    }
    function emitEventsMock_AuthorizedReceiver_AuthorizedSendersChanged(address[] memory senders,address changedBy) public {
        emit EventsMock_AuthorizedReceiver_AuthorizedSendersChanged(senders,changedBy);
    }
    function emitEventsMock_BatchVRFCoordinatorV2_ErrorReturned(uint256 requestId,string memory reason) public {
        emit EventsMock_BatchVRFCoordinatorV2_ErrorReturned(requestId,reason);
    }
    function emitEventsMock_BatchVRFCoordinatorV2_RawErrorReturned(uint256 requestId,bytes memory lowLevelData) public {
        emit EventsMock_BatchVRFCoordinatorV2_RawErrorReturned(requestId,lowLevelData);
    }
    function emitEventsMock_ChainlinkClient_ChainlinkCancelled(bytes32 id) public {
        emit EventsMock_ChainlinkClient_ChainlinkCancelled(id);
    }
    function emitEventsMock_ChainlinkClient_ChainlinkFulfilled(bytes32 id) public {
        emit EventsMock_ChainlinkClient_ChainlinkFulfilled(id);
    }
    function emitEventsMock_ChainlinkClient_ChainlinkRequested(bytes32 id) public {
        emit EventsMock_ChainlinkClient_ChainlinkRequested(id);
    }
    function emitEventsMock_ConfirmedOwnerUpgradeable_Initialized(uint8 version) public {
        emit EventsMock_ConfirmedOwnerUpgradeable_Initialized(version);
    }
    function emitEventsMock_ConfirmedOwnerUpgradeable_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_ConfirmedOwnerUpgradeable_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_ConfirmedOwnerUpgradeable_OwnershipTransferred(address from,address to) public {
        emit EventsMock_ConfirmedOwnerUpgradeable_OwnershipTransferred(from,to);
    }
    function emitEventsMock_ConfirmedOwnerWithProposal_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_ConfirmedOwnerWithProposal_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_ConfirmedOwnerWithProposal_OwnershipTransferred(address from,address to) public {
        emit EventsMock_ConfirmedOwnerWithProposal_OwnershipTransferred(from,to);
    }
    function emitEventsMock_ConfirmedOwner_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_ConfirmedOwner_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_ConfirmedOwner_OwnershipTransferred(address from,address to) public {
        emit EventsMock_ConfirmedOwner_OwnershipTransferred(from,to);
    }
    function emitEventsMock_ContextUpgradeable_Initialized(uint8 version) public {
        emit EventsMock_ContextUpgradeable_Initialized(version);
    }
    function emitEventsMock_CronUpkeepFactory_NewCronUpkeepCreated(address upkeep,address owner) public {
        emit EventsMock_CronUpkeepFactory_NewCronUpkeepCreated(upkeep,owner);
    }
    function emitEventsMock_CronUpkeepFactory_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_CronUpkeepFactory_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_CronUpkeepFactory_OwnershipTransferred(address from,address to) public {
        emit EventsMock_CronUpkeepFactory_OwnershipTransferred(from,to);
    }
    function emitEventsMock_CronUpkeep_CronJobCreated(uint256 id,address target,bytes memory handler) public {
        emit EventsMock_CronUpkeep_CronJobCreated(id,target,handler);
    }
    function emitEventsMock_CronUpkeep_CronJobDeleted(uint256 id) public {
        emit EventsMock_CronUpkeep_CronJobDeleted(id);
    }
    function emitEventsMock_CronUpkeep_CronJobExecuted(uint256 id,uint256 timestamp) public {
        emit EventsMock_CronUpkeep_CronJobExecuted(id,timestamp);
    }
    function emitEventsMock_CronUpkeep_CronJobUpdated(uint256 id,address target,bytes memory handler) public {
        emit EventsMock_CronUpkeep_CronJobUpdated(id,target,handler);
    }
    function emitEventsMock_CronUpkeep_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_CronUpkeep_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_CronUpkeep_OwnershipTransferred(address from,address to) public {
        emit EventsMock_CronUpkeep_OwnershipTransferred(from,to);
    }
    function emitEventsMock_CronUpkeep_Paused(address account) public {
        emit EventsMock_CronUpkeep_Paused(account);
    }
    function emitEventsMock_CronUpkeep_Unpaused(address account) public {
        emit EventsMock_CronUpkeep_Unpaused(account);
    }
    function emitEventsMock_ENSInterface_NewOwner(bytes32 node,bytes32 label,address owner) public {
        emit EventsMock_ENSInterface_NewOwner(node,label,owner);
    }
    function emitEventsMock_ENSInterface_NewResolver(bytes32 node,address resolver) public {
        emit EventsMock_ENSInterface_NewResolver(node,resolver);
    }
    function emitEventsMock_ENSInterface_NewTTL(bytes32 node,uint64 ttl) public {
        emit EventsMock_ENSInterface_NewTTL(node,ttl);
    }
    function emitEventsMock_ENSInterface_Transfer(bytes32 node,address owner) public {
        emit EventsMock_ENSInterface_Transfer(node,owner);
    }
    function emitEventsMock_ERC1967Proxy_AdminChanged(address previousAdmin,address newAdmin) public {
        emit EventsMock_ERC1967Proxy_AdminChanged(previousAdmin,newAdmin);
    }
    function emitEventsMock_ERC1967Proxy_BeaconUpgraded(address beacon) public {
        emit EventsMock_ERC1967Proxy_BeaconUpgraded(beacon);
    }
    function emitEventsMock_ERC1967Proxy_Upgraded(address implementation) public {
        emit EventsMock_ERC1967Proxy_Upgraded(implementation);
    }
    function emitEventsMock_ERC1967Upgrade_AdminChanged(address previousAdmin,address newAdmin) public {
        emit EventsMock_ERC1967Upgrade_AdminChanged(previousAdmin,newAdmin);
    }
    function emitEventsMock_ERC1967Upgrade_BeaconUpgraded(address beacon) public {
        emit EventsMock_ERC1967Upgrade_BeaconUpgraded(beacon);
    }
    function emitEventsMock_ERC1967Upgrade_Upgraded(address implementation) public {
        emit EventsMock_ERC1967Upgrade_Upgraded(implementation);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_AuthorizedSendersChanged(address[] memory senders,address changedBy) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_AuthorizedSendersChanged(senders,changedBy);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_BillingEnd(bytes32 requestId,uint64 subscriptionId,uint96 signerPayment,uint96 transmitterPayment,uint96 totalCost,bool success) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_BillingEnd(requestId,subscriptionId,signerPayment,transmitterPayment,totalCost,success);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_BillingStart(bytes32 requestId,EventsMock_FunctionsBillingRegistryWithInit_Commitment memory commitment) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_BillingStart(requestId,commitment);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_ConfigSet(uint32 maxGasLimit,uint32 stalenessSeconds,uint256 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,uint32 gasOverhead) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_ConfigSet(maxGasLimit,stalenessSeconds,gasAfterPaymentCalculation,fallbackWeiPerUnitLink,gasOverhead);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_FundsRecovered(address to,uint256 amount) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_FundsRecovered(to,amount);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_Initialized(uint8 version) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_Initialized(version);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_OwnershipTransferred(address from,address to) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_OwnershipTransferred(from,to);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_Paused(address account) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_Paused(account);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_RequestTimedOut(bytes32 requestId) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_RequestTimedOut(requestId);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_SubscriptionCanceled(uint64 subscriptionId,address to,uint256 amount) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_SubscriptionCanceled(subscriptionId,to,amount);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_SubscriptionConsumerAdded(uint64 subscriptionId,address consumer) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_SubscriptionConsumerAdded(subscriptionId,consumer);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_SubscriptionConsumerRemoved(uint64 subscriptionId,address consumer) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_SubscriptionConsumerRemoved(subscriptionId,consumer);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_SubscriptionCreated(uint64 subscriptionId,address owner) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_SubscriptionCreated(subscriptionId,owner);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_SubscriptionFunded(uint64 subscriptionId,uint256 oldBalance,uint256 newBalance) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_SubscriptionFunded(subscriptionId,oldBalance,newBalance);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferRequested(uint64 subscriptionId,address from,address to) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferRequested(subscriptionId,from,to);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferred(uint64 subscriptionId,address from,address to) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferred(subscriptionId,from,to);
    }
    function emitEventsMock_FunctionsBillingRegistryWithInit_Unpaused(address account) public {
        emit EventsMock_FunctionsBillingRegistryWithInit_Unpaused(account);
    }
    function emitEventsMock_FunctionsBillingRegistry_AuthorizedSendersChanged(address[] memory senders,address changedBy) public {
        emit EventsMock_FunctionsBillingRegistry_AuthorizedSendersChanged(senders,changedBy);
    }
    function emitEventsMock_FunctionsBillingRegistry_BillingEnd(bytes32 requestId,uint64 subscriptionId,uint96 signerPayment,uint96 transmitterPayment,uint96 totalCost,bool success) public {
        emit EventsMock_FunctionsBillingRegistry_BillingEnd(requestId,subscriptionId,signerPayment,transmitterPayment,totalCost,success);
    }
    function emitEventsMock_FunctionsBillingRegistry_BillingStart(bytes32 requestId,EventsMock_FunctionsBillingRegistry_Commitment memory commitment) public {
        emit EventsMock_FunctionsBillingRegistry_BillingStart(requestId,commitment);
    }
    function emitEventsMock_FunctionsBillingRegistry_ConfigSet(uint32 maxGasLimit,uint32 stalenessSeconds,uint256 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,uint32 gasOverhead) public {
        emit EventsMock_FunctionsBillingRegistry_ConfigSet(maxGasLimit,stalenessSeconds,gasAfterPaymentCalculation,fallbackWeiPerUnitLink,gasOverhead);
    }
    function emitEventsMock_FunctionsBillingRegistry_FundsRecovered(address to,uint256 amount) public {
        emit EventsMock_FunctionsBillingRegistry_FundsRecovered(to,amount);
    }
    function emitEventsMock_FunctionsBillingRegistry_Initialized(uint8 version) public {
        emit EventsMock_FunctionsBillingRegistry_Initialized(version);
    }
    function emitEventsMock_FunctionsBillingRegistry_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_FunctionsBillingRegistry_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_FunctionsBillingRegistry_OwnershipTransferred(address from,address to) public {
        emit EventsMock_FunctionsBillingRegistry_OwnershipTransferred(from,to);
    }
    function emitEventsMock_FunctionsBillingRegistry_Paused(address account) public {
        emit EventsMock_FunctionsBillingRegistry_Paused(account);
    }
    function emitEventsMock_FunctionsBillingRegistry_RequestTimedOut(bytes32 requestId) public {
        emit EventsMock_FunctionsBillingRegistry_RequestTimedOut(requestId);
    }
    function emitEventsMock_FunctionsBillingRegistry_SubscriptionCanceled(uint64 subscriptionId,address to,uint256 amount) public {
        emit EventsMock_FunctionsBillingRegistry_SubscriptionCanceled(subscriptionId,to,amount);
    }
    function emitEventsMock_FunctionsBillingRegistry_SubscriptionConsumerAdded(uint64 subscriptionId,address consumer) public {
        emit EventsMock_FunctionsBillingRegistry_SubscriptionConsumerAdded(subscriptionId,consumer);
    }
    function emitEventsMock_FunctionsBillingRegistry_SubscriptionConsumerRemoved(uint64 subscriptionId,address consumer) public {
        emit EventsMock_FunctionsBillingRegistry_SubscriptionConsumerRemoved(subscriptionId,consumer);
    }
    function emitEventsMock_FunctionsBillingRegistry_SubscriptionCreated(uint64 subscriptionId,address owner) public {
        emit EventsMock_FunctionsBillingRegistry_SubscriptionCreated(subscriptionId,owner);
    }
    function emitEventsMock_FunctionsBillingRegistry_SubscriptionFunded(uint64 subscriptionId,uint256 oldBalance,uint256 newBalance) public {
        emit EventsMock_FunctionsBillingRegistry_SubscriptionFunded(subscriptionId,oldBalance,newBalance);
    }
    function emitEventsMock_FunctionsBillingRegistry_SubscriptionOwnerTransferRequested(uint64 subscriptionId,address from,address to) public {
        emit EventsMock_FunctionsBillingRegistry_SubscriptionOwnerTransferRequested(subscriptionId,from,to);
    }
    function emitEventsMock_FunctionsBillingRegistry_SubscriptionOwnerTransferred(uint64 subscriptionId,address from,address to) public {
        emit EventsMock_FunctionsBillingRegistry_SubscriptionOwnerTransferred(subscriptionId,from,to);
    }
    function emitEventsMock_FunctionsBillingRegistry_Unpaused(address account) public {
        emit EventsMock_FunctionsBillingRegistry_Unpaused(account);
    }
    function emitEventsMock_FunctionsClientExample_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_FunctionsClientExample_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_FunctionsClientExample_OwnershipTransferred(address from,address to) public {
        emit EventsMock_FunctionsClientExample_OwnershipTransferred(from,to);
    }
    function emitEventsMock_FunctionsClientExample_RequestFulfilled(bytes32 id) public {
        emit EventsMock_FunctionsClientExample_RequestFulfilled(id);
    }
    function emitEventsMock_FunctionsClientExample_RequestSent(bytes32 id) public {
        emit EventsMock_FunctionsClientExample_RequestSent(id);
    }
    function emitEventsMock_FunctionsClient_RequestFulfilled(bytes32 id) public {
        emit EventsMock_FunctionsClient_RequestFulfilled(id);
    }
    function emitEventsMock_FunctionsClient_RequestSent(bytes32 id) public {
        emit EventsMock_FunctionsClient_RequestSent(id);
    }
    function emitEventsMock_FunctionsOracleWithInit_AuthorizedSendersActive(address account) public {
        emit EventsMock_FunctionsOracleWithInit_AuthorizedSendersActive(account);
    }
    function emitEventsMock_FunctionsOracleWithInit_AuthorizedSendersChanged(address[] memory senders,address changedBy) public {
        emit EventsMock_FunctionsOracleWithInit_AuthorizedSendersChanged(senders,changedBy);
    }
    function emitEventsMock_FunctionsOracleWithInit_AuthorizedSendersDeactive(address account) public {
        emit EventsMock_FunctionsOracleWithInit_AuthorizedSendersDeactive(account);
    }
    function emitEventsMock_FunctionsOracleWithInit_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
        emit EventsMock_FunctionsOracleWithInit_ConfigSet(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
    }
    function emitEventsMock_FunctionsOracleWithInit_Initialized(uint8 version) public {
        emit EventsMock_FunctionsOracleWithInit_Initialized(version);
    }
    function emitEventsMock_FunctionsOracleWithInit_InvalidRequestID(bytes32 requestId) public {
        emit EventsMock_FunctionsOracleWithInit_InvalidRequestID(requestId);
    }
    function emitEventsMock_FunctionsOracleWithInit_OracleRequest(bytes32 requestId,address requestingContract,address requestInitiator,uint64 subscriptionId,address subscriptionOwner,bytes memory data) public {
        emit EventsMock_FunctionsOracleWithInit_OracleRequest(requestId,requestingContract,requestInitiator,subscriptionId,subscriptionOwner,data);
    }
    function emitEventsMock_FunctionsOracleWithInit_OracleResponse(bytes32 requestId) public {
        emit EventsMock_FunctionsOracleWithInit_OracleResponse(requestId);
    }
    function emitEventsMock_FunctionsOracleWithInit_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_FunctionsOracleWithInit_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_FunctionsOracleWithInit_OwnershipTransferred(address from,address to) public {
        emit EventsMock_FunctionsOracleWithInit_OwnershipTransferred(from,to);
    }
    function emitEventsMock_FunctionsOracleWithInit_Transmitted(bytes32 configDigest,uint32 epoch) public {
        emit EventsMock_FunctionsOracleWithInit_Transmitted(configDigest,epoch);
    }
    function emitEventsMock_FunctionsOracleWithInit_UserCallbackError(bytes32 requestId,string memory reason) public {
        emit EventsMock_FunctionsOracleWithInit_UserCallbackError(requestId,reason);
    }
    function emitEventsMock_FunctionsOracleWithInit_UserCallbackRawError(bytes32 requestId,bytes memory lowLevelData) public {
        emit EventsMock_FunctionsOracleWithInit_UserCallbackRawError(requestId,lowLevelData);
    }
    function emitEventsMock_FunctionsOracle_AuthorizedSendersActive(address account) public {
        emit EventsMock_FunctionsOracle_AuthorizedSendersActive(account);
    }
    function emitEventsMock_FunctionsOracle_AuthorizedSendersChanged(address[] memory senders,address changedBy) public {
        emit EventsMock_FunctionsOracle_AuthorizedSendersChanged(senders,changedBy);
    }
    function emitEventsMock_FunctionsOracle_AuthorizedSendersDeactive(address account) public {
        emit EventsMock_FunctionsOracle_AuthorizedSendersDeactive(account);
    }
    function emitEventsMock_FunctionsOracle_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
        emit EventsMock_FunctionsOracle_ConfigSet(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
    }
    function emitEventsMock_FunctionsOracle_Initialized(uint8 version) public {
        emit EventsMock_FunctionsOracle_Initialized(version);
    }
    function emitEventsMock_FunctionsOracle_InvalidRequestID(bytes32 requestId) public {
        emit EventsMock_FunctionsOracle_InvalidRequestID(requestId);
    }
    function emitEventsMock_FunctionsOracle_OracleRequest(bytes32 requestId,address requestingContract,address requestInitiator,uint64 subscriptionId,address subscriptionOwner,bytes memory data) public {
        emit EventsMock_FunctionsOracle_OracleRequest(requestId,requestingContract,requestInitiator,subscriptionId,subscriptionOwner,data);
    }
    function emitEventsMock_FunctionsOracle_OracleResponse(bytes32 requestId) public {
        emit EventsMock_FunctionsOracle_OracleResponse(requestId);
    }
    function emitEventsMock_FunctionsOracle_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_FunctionsOracle_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_FunctionsOracle_OwnershipTransferred(address from,address to) public {
        emit EventsMock_FunctionsOracle_OwnershipTransferred(from,to);
    }
    function emitEventsMock_FunctionsOracle_Transmitted(bytes32 configDigest,uint32 epoch) public {
        emit EventsMock_FunctionsOracle_Transmitted(configDigest,epoch);
    }
    function emitEventsMock_FunctionsOracle_UserCallbackError(bytes32 requestId,string memory reason) public {
        emit EventsMock_FunctionsOracle_UserCallbackError(requestId,reason);
    }
    function emitEventsMock_FunctionsOracle_UserCallbackRawError(bytes32 requestId,bytes memory lowLevelData) public {
        emit EventsMock_FunctionsOracle_UserCallbackRawError(requestId,lowLevelData);
    }
    function emitEventsMock_Initializable_Initialized(uint8 version) public {
        emit EventsMock_Initializable_Initialized(version);
    }
    function emitEventsMock_KeeperRegistry1_2_ConfigSet(EventsMock_KeeperRegistry1_2_Config memory config) public {
        emit EventsMock_KeeperRegistry1_2_ConfigSet(config);
    }
    function emitEventsMock_KeeperRegistry1_2_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit EventsMock_KeeperRegistry1_2_FundsAdded(id,from,amount);
    }
    function emitEventsMock_KeeperRegistry1_2_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit EventsMock_KeeperRegistry1_2_FundsWithdrawn(id,amount,to);
    }
    function emitEventsMock_KeeperRegistry1_2_KeepersUpdated(address[] memory keepers,address[] memory payees) public {
        emit EventsMock_KeeperRegistry1_2_KeepersUpdated(keepers,payees);
    }
    function emitEventsMock_KeeperRegistry1_2_OwnerFundsWithdrawn(uint96 amount) public {
        emit EventsMock_KeeperRegistry1_2_OwnerFundsWithdrawn(amount);
    }
    function emitEventsMock_KeeperRegistry1_2_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_KeeperRegistry1_2_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_KeeperRegistry1_2_OwnershipTransferred(address from,address to) public {
        emit EventsMock_KeeperRegistry1_2_OwnershipTransferred(from,to);
    }
    function emitEventsMock_KeeperRegistry1_2_Paused(address account) public {
        emit EventsMock_KeeperRegistry1_2_Paused(account);
    }
    function emitEventsMock_KeeperRegistry1_2_PayeeshipTransferRequested(address keeper,address from,address to) public {
        emit EventsMock_KeeperRegistry1_2_PayeeshipTransferRequested(keeper,from,to);
    }
    function emitEventsMock_KeeperRegistry1_2_PayeeshipTransferred(address keeper,address from,address to) public {
        emit EventsMock_KeeperRegistry1_2_PayeeshipTransferred(keeper,from,to);
    }
    function emitEventsMock_KeeperRegistry1_2_PaymentWithdrawn(address keeper,uint256 amount,address to,address payee) public {
        emit EventsMock_KeeperRegistry1_2_PaymentWithdrawn(keeper,amount,to,payee);
    }
    function emitEventsMock_KeeperRegistry1_2_Unpaused(address account) public {
        emit EventsMock_KeeperRegistry1_2_Unpaused(account);
    }
    function emitEventsMock_KeeperRegistry1_2_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit EventsMock_KeeperRegistry1_2_UpkeepCanceled(id,atBlockHeight);
    }
    function emitEventsMock_KeeperRegistry1_2_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit EventsMock_KeeperRegistry1_2_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitEventsMock_KeeperRegistry1_2_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit EventsMock_KeeperRegistry1_2_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitEventsMock_KeeperRegistry1_2_UpkeepPerformed(uint256 id,bool success,address from,uint96 payment,bytes memory performData) public {
        emit EventsMock_KeeperRegistry1_2_UpkeepPerformed(id,success,from,payment,performData);
    }
    function emitEventsMock_KeeperRegistry1_2_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit EventsMock_KeeperRegistry1_2_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitEventsMock_KeeperRegistry1_2_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit EventsMock_KeeperRegistry1_2_UpkeepRegistered(id,executeGas,admin);
    }
    function emitEventsMock_KeeperRegistry1_3_ConfigSet(EventsMock_KeeperRegistry1_3_Config memory config) public {
        emit EventsMock_KeeperRegistry1_3_ConfigSet(config);
    }
    function emitEventsMock_KeeperRegistry1_3_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit EventsMock_KeeperRegistry1_3_FundsAdded(id,from,amount);
    }
    function emitEventsMock_KeeperRegistry1_3_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit EventsMock_KeeperRegistry1_3_FundsWithdrawn(id,amount,to);
    }
    function emitEventsMock_KeeperRegistry1_3_KeepersUpdated(address[] memory keepers,address[] memory payees) public {
        emit EventsMock_KeeperRegistry1_3_KeepersUpdated(keepers,payees);
    }
    function emitEventsMock_KeeperRegistry1_3_OwnerFundsWithdrawn(uint96 amount) public {
        emit EventsMock_KeeperRegistry1_3_OwnerFundsWithdrawn(amount);
    }
    function emitEventsMock_KeeperRegistry1_3_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_KeeperRegistry1_3_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_KeeperRegistry1_3_OwnershipTransferred(address from,address to) public {
        emit EventsMock_KeeperRegistry1_3_OwnershipTransferred(from,to);
    }
    function emitEventsMock_KeeperRegistry1_3_Paused(address account) public {
        emit EventsMock_KeeperRegistry1_3_Paused(account);
    }
    function emitEventsMock_KeeperRegistry1_3_PayeeshipTransferRequested(address keeper,address from,address to) public {
        emit EventsMock_KeeperRegistry1_3_PayeeshipTransferRequested(keeper,from,to);
    }
    function emitEventsMock_KeeperRegistry1_3_PayeeshipTransferred(address keeper,address from,address to) public {
        emit EventsMock_KeeperRegistry1_3_PayeeshipTransferred(keeper,from,to);
    }
    function emitEventsMock_KeeperRegistry1_3_PaymentWithdrawn(address keeper,uint256 amount,address to,address payee) public {
        emit EventsMock_KeeperRegistry1_3_PaymentWithdrawn(keeper,amount,to,payee);
    }
    function emitEventsMock_KeeperRegistry1_3_Unpaused(address account) public {
        emit EventsMock_KeeperRegistry1_3_Unpaused(account);
    }
    function emitEventsMock_KeeperRegistry1_3_UpkeepAdminTransferRequested(uint256 id,address from,address to) public {
        emit EventsMock_KeeperRegistry1_3_UpkeepAdminTransferRequested(id,from,to);
    }
    function emitEventsMock_KeeperRegistry1_3_UpkeepAdminTransferred(uint256 id,address from,address to) public {
        emit EventsMock_KeeperRegistry1_3_UpkeepAdminTransferred(id,from,to);
    }
    function emitEventsMock_KeeperRegistry1_3_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit EventsMock_KeeperRegistry1_3_UpkeepCanceled(id,atBlockHeight);
    }
    function emitEventsMock_KeeperRegistry1_3_UpkeepCheckDataUpdated(uint256 id,bytes memory newCheckData) public {
        emit EventsMock_KeeperRegistry1_3_UpkeepCheckDataUpdated(id,newCheckData);
    }
    function emitEventsMock_KeeperRegistry1_3_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit EventsMock_KeeperRegistry1_3_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitEventsMock_KeeperRegistry1_3_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit EventsMock_KeeperRegistry1_3_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitEventsMock_KeeperRegistry1_3_UpkeepPaused(uint256 id) public {
        emit EventsMock_KeeperRegistry1_3_UpkeepPaused(id);
    }
    function emitEventsMock_KeeperRegistry1_3_UpkeepPerformed(uint256 id,bool success,address from,uint96 payment,bytes memory performData) public {
        emit EventsMock_KeeperRegistry1_3_UpkeepPerformed(id,success,from,payment,performData);
    }
    function emitEventsMock_KeeperRegistry1_3_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit EventsMock_KeeperRegistry1_3_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitEventsMock_KeeperRegistry1_3_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit EventsMock_KeeperRegistry1_3_UpkeepRegistered(id,executeGas,admin);
    }
    function emitEventsMock_KeeperRegistry1_3_UpkeepUnpaused(uint256 id) public {
        emit EventsMock_KeeperRegistry1_3_UpkeepUnpaused(id);
    }
    function emitEventsMock_KeeperRegistry2_0_CancelledUpkeepReport(uint256 id) public {
        emit EventsMock_KeeperRegistry2_0_CancelledUpkeepReport(id);
    }
    function emitEventsMock_KeeperRegistry2_0_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
        emit EventsMock_KeeperRegistry2_0_ConfigSet(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
    }
    function emitEventsMock_KeeperRegistry2_0_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit EventsMock_KeeperRegistry2_0_FundsAdded(id,from,amount);
    }
    function emitEventsMock_KeeperRegistry2_0_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit EventsMock_KeeperRegistry2_0_FundsWithdrawn(id,amount,to);
    }
    function emitEventsMock_KeeperRegistry2_0_InsufficientFundsUpkeepReport(uint256 id) public {
        emit EventsMock_KeeperRegistry2_0_InsufficientFundsUpkeepReport(id);
    }
    function emitEventsMock_KeeperRegistry2_0_OwnerFundsWithdrawn(uint96 amount) public {
        emit EventsMock_KeeperRegistry2_0_OwnerFundsWithdrawn(amount);
    }
    function emitEventsMock_KeeperRegistry2_0_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_KeeperRegistry2_0_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_KeeperRegistry2_0_OwnershipTransferred(address from,address to) public {
        emit EventsMock_KeeperRegistry2_0_OwnershipTransferred(from,to);
    }
    function emitEventsMock_KeeperRegistry2_0_Paused(address account) public {
        emit EventsMock_KeeperRegistry2_0_Paused(account);
    }
    function emitEventsMock_KeeperRegistry2_0_PayeesUpdated(address[] memory transmitters,address[] memory payees) public {
        emit EventsMock_KeeperRegistry2_0_PayeesUpdated(transmitters,payees);
    }
    function emitEventsMock_KeeperRegistry2_0_PayeeshipTransferRequested(address transmitter,address from,address to) public {
        emit EventsMock_KeeperRegistry2_0_PayeeshipTransferRequested(transmitter,from,to);
    }
    function emitEventsMock_KeeperRegistry2_0_PayeeshipTransferred(address transmitter,address from,address to) public {
        emit EventsMock_KeeperRegistry2_0_PayeeshipTransferred(transmitter,from,to);
    }
    function emitEventsMock_KeeperRegistry2_0_PaymentWithdrawn(address transmitter,uint256 amount,address to,address payee) public {
        emit EventsMock_KeeperRegistry2_0_PaymentWithdrawn(transmitter,amount,to,payee);
    }
    function emitEventsMock_KeeperRegistry2_0_ReorgedUpkeepReport(uint256 id) public {
        emit EventsMock_KeeperRegistry2_0_ReorgedUpkeepReport(id);
    }
    function emitEventsMock_KeeperRegistry2_0_StaleUpkeepReport(uint256 id) public {
        emit EventsMock_KeeperRegistry2_0_StaleUpkeepReport(id);
    }
    function emitEventsMock_KeeperRegistry2_0_Transmitted(bytes32 configDigest,uint32 epoch) public {
        emit EventsMock_KeeperRegistry2_0_Transmitted(configDigest,epoch);
    }
    function emitEventsMock_KeeperRegistry2_0_Unpaused(address account) public {
        emit EventsMock_KeeperRegistry2_0_Unpaused(account);
    }
    function emitEventsMock_KeeperRegistry2_0_UpkeepAdminTransferRequested(uint256 id,address from,address to) public {
        emit EventsMock_KeeperRegistry2_0_UpkeepAdminTransferRequested(id,from,to);
    }
    function emitEventsMock_KeeperRegistry2_0_UpkeepAdminTransferred(uint256 id,address from,address to) public {
        emit EventsMock_KeeperRegistry2_0_UpkeepAdminTransferred(id,from,to);
    }
    function emitEventsMock_KeeperRegistry2_0_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit EventsMock_KeeperRegistry2_0_UpkeepCanceled(id,atBlockHeight);
    }
    function emitEventsMock_KeeperRegistry2_0_UpkeepCheckDataUpdated(uint256 id,bytes memory newCheckData) public {
        emit EventsMock_KeeperRegistry2_0_UpkeepCheckDataUpdated(id,newCheckData);
    }
    function emitEventsMock_KeeperRegistry2_0_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit EventsMock_KeeperRegistry2_0_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitEventsMock_KeeperRegistry2_0_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit EventsMock_KeeperRegistry2_0_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitEventsMock_KeeperRegistry2_0_UpkeepOffchainConfigSet(uint256 id,bytes memory offchainConfig) public {
        emit EventsMock_KeeperRegistry2_0_UpkeepOffchainConfigSet(id,offchainConfig);
    }
    function emitEventsMock_KeeperRegistry2_0_UpkeepPaused(uint256 id) public {
        emit EventsMock_KeeperRegistry2_0_UpkeepPaused(id);
    }
    function emitEventsMock_KeeperRegistry2_0_UpkeepPerformed(uint256 id,bool success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment) public {
        emit EventsMock_KeeperRegistry2_0_UpkeepPerformed(id,success,checkBlockNumber,gasUsed,gasOverhead,totalPayment);
    }
    function emitEventsMock_KeeperRegistry2_0_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit EventsMock_KeeperRegistry2_0_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitEventsMock_KeeperRegistry2_0_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit EventsMock_KeeperRegistry2_0_UpkeepRegistered(id,executeGas,admin);
    }
    function emitEventsMock_KeeperRegistry2_0_UpkeepUnpaused(uint256 id) public {
        emit EventsMock_KeeperRegistry2_0_UpkeepUnpaused(id);
    }
    function emitEventsMock_KeeperRegistryBase1_3_ConfigSet(EventsMock_KeeperRegistryBase1_3_Config memory config) public {
        emit EventsMock_KeeperRegistryBase1_3_ConfigSet(config);
    }
    function emitEventsMock_KeeperRegistryBase1_3_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit EventsMock_KeeperRegistryBase1_3_FundsAdded(id,from,amount);
    }
    function emitEventsMock_KeeperRegistryBase1_3_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit EventsMock_KeeperRegistryBase1_3_FundsWithdrawn(id,amount,to);
    }
    function emitEventsMock_KeeperRegistryBase1_3_KeepersUpdated(address[] memory keepers,address[] memory payees) public {
        emit EventsMock_KeeperRegistryBase1_3_KeepersUpdated(keepers,payees);
    }
    function emitEventsMock_KeeperRegistryBase1_3_OwnerFundsWithdrawn(uint96 amount) public {
        emit EventsMock_KeeperRegistryBase1_3_OwnerFundsWithdrawn(amount);
    }
    function emitEventsMock_KeeperRegistryBase1_3_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_KeeperRegistryBase1_3_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_KeeperRegistryBase1_3_OwnershipTransferred(address from,address to) public {
        emit EventsMock_KeeperRegistryBase1_3_OwnershipTransferred(from,to);
    }
    function emitEventsMock_KeeperRegistryBase1_3_Paused(address account) public {
        emit EventsMock_KeeperRegistryBase1_3_Paused(account);
    }
    function emitEventsMock_KeeperRegistryBase1_3_PayeeshipTransferRequested(address keeper,address from,address to) public {
        emit EventsMock_KeeperRegistryBase1_3_PayeeshipTransferRequested(keeper,from,to);
    }
    function emitEventsMock_KeeperRegistryBase1_3_PayeeshipTransferred(address keeper,address from,address to) public {
        emit EventsMock_KeeperRegistryBase1_3_PayeeshipTransferred(keeper,from,to);
    }
    function emitEventsMock_KeeperRegistryBase1_3_PaymentWithdrawn(address keeper,uint256 amount,address to,address payee) public {
        emit EventsMock_KeeperRegistryBase1_3_PaymentWithdrawn(keeper,amount,to,payee);
    }
    function emitEventsMock_KeeperRegistryBase1_3_Unpaused(address account) public {
        emit EventsMock_KeeperRegistryBase1_3_Unpaused(account);
    }
    function emitEventsMock_KeeperRegistryBase1_3_UpkeepAdminTransferRequested(uint256 id,address from,address to) public {
        emit EventsMock_KeeperRegistryBase1_3_UpkeepAdminTransferRequested(id,from,to);
    }
    function emitEventsMock_KeeperRegistryBase1_3_UpkeepAdminTransferred(uint256 id,address from,address to) public {
        emit EventsMock_KeeperRegistryBase1_3_UpkeepAdminTransferred(id,from,to);
    }
    function emitEventsMock_KeeperRegistryBase1_3_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit EventsMock_KeeperRegistryBase1_3_UpkeepCanceled(id,atBlockHeight);
    }
    function emitEventsMock_KeeperRegistryBase1_3_UpkeepCheckDataUpdated(uint256 id,bytes memory newCheckData) public {
        emit EventsMock_KeeperRegistryBase1_3_UpkeepCheckDataUpdated(id,newCheckData);
    }
    function emitEventsMock_KeeperRegistryBase1_3_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit EventsMock_KeeperRegistryBase1_3_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitEventsMock_KeeperRegistryBase1_3_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit EventsMock_KeeperRegistryBase1_3_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitEventsMock_KeeperRegistryBase1_3_UpkeepPaused(uint256 id) public {
        emit EventsMock_KeeperRegistryBase1_3_UpkeepPaused(id);
    }
    function emitEventsMock_KeeperRegistryBase1_3_UpkeepPerformed(uint256 id,bool success,address from,uint96 payment,bytes memory performData) public {
        emit EventsMock_KeeperRegistryBase1_3_UpkeepPerformed(id,success,from,payment,performData);
    }
    function emitEventsMock_KeeperRegistryBase1_3_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit EventsMock_KeeperRegistryBase1_3_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitEventsMock_KeeperRegistryBase1_3_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit EventsMock_KeeperRegistryBase1_3_UpkeepRegistered(id,executeGas,admin);
    }
    function emitEventsMock_KeeperRegistryBase1_3_UpkeepUnpaused(uint256 id) public {
        emit EventsMock_KeeperRegistryBase1_3_UpkeepUnpaused(id);
    }
    function emitEventsMock_KeeperRegistryBase2_0_CancelledUpkeepReport(uint256 id) public {
        emit EventsMock_KeeperRegistryBase2_0_CancelledUpkeepReport(id);
    }
    function emitEventsMock_KeeperRegistryBase2_0_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit EventsMock_KeeperRegistryBase2_0_FundsAdded(id,from,amount);
    }
    function emitEventsMock_KeeperRegistryBase2_0_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit EventsMock_KeeperRegistryBase2_0_FundsWithdrawn(id,amount,to);
    }
    function emitEventsMock_KeeperRegistryBase2_0_InsufficientFundsUpkeepReport(uint256 id) public {
        emit EventsMock_KeeperRegistryBase2_0_InsufficientFundsUpkeepReport(id);
    }
    function emitEventsMock_KeeperRegistryBase2_0_OwnerFundsWithdrawn(uint96 amount) public {
        emit EventsMock_KeeperRegistryBase2_0_OwnerFundsWithdrawn(amount);
    }
    function emitEventsMock_KeeperRegistryBase2_0_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_KeeperRegistryBase2_0_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_KeeperRegistryBase2_0_OwnershipTransferred(address from,address to) public {
        emit EventsMock_KeeperRegistryBase2_0_OwnershipTransferred(from,to);
    }
    function emitEventsMock_KeeperRegistryBase2_0_Paused(address account) public {
        emit EventsMock_KeeperRegistryBase2_0_Paused(account);
    }
    function emitEventsMock_KeeperRegistryBase2_0_PayeesUpdated(address[] memory transmitters,address[] memory payees) public {
        emit EventsMock_KeeperRegistryBase2_0_PayeesUpdated(transmitters,payees);
    }
    function emitEventsMock_KeeperRegistryBase2_0_PayeeshipTransferRequested(address transmitter,address from,address to) public {
        emit EventsMock_KeeperRegistryBase2_0_PayeeshipTransferRequested(transmitter,from,to);
    }
    function emitEventsMock_KeeperRegistryBase2_0_PayeeshipTransferred(address transmitter,address from,address to) public {
        emit EventsMock_KeeperRegistryBase2_0_PayeeshipTransferred(transmitter,from,to);
    }
    function emitEventsMock_KeeperRegistryBase2_0_PaymentWithdrawn(address transmitter,uint256 amount,address to,address payee) public {
        emit EventsMock_KeeperRegistryBase2_0_PaymentWithdrawn(transmitter,amount,to,payee);
    }
    function emitEventsMock_KeeperRegistryBase2_0_ReorgedUpkeepReport(uint256 id) public {
        emit EventsMock_KeeperRegistryBase2_0_ReorgedUpkeepReport(id);
    }
    function emitEventsMock_KeeperRegistryBase2_0_StaleUpkeepReport(uint256 id) public {
        emit EventsMock_KeeperRegistryBase2_0_StaleUpkeepReport(id);
    }
    function emitEventsMock_KeeperRegistryBase2_0_Unpaused(address account) public {
        emit EventsMock_KeeperRegistryBase2_0_Unpaused(account);
    }
    function emitEventsMock_KeeperRegistryBase2_0_UpkeepAdminTransferRequested(uint256 id,address from,address to) public {
        emit EventsMock_KeeperRegistryBase2_0_UpkeepAdminTransferRequested(id,from,to);
    }
    function emitEventsMock_KeeperRegistryBase2_0_UpkeepAdminTransferred(uint256 id,address from,address to) public {
        emit EventsMock_KeeperRegistryBase2_0_UpkeepAdminTransferred(id,from,to);
    }
    function emitEventsMock_KeeperRegistryBase2_0_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit EventsMock_KeeperRegistryBase2_0_UpkeepCanceled(id,atBlockHeight);
    }
    function emitEventsMock_KeeperRegistryBase2_0_UpkeepCheckDataUpdated(uint256 id,bytes memory newCheckData) public {
        emit EventsMock_KeeperRegistryBase2_0_UpkeepCheckDataUpdated(id,newCheckData);
    }
    function emitEventsMock_KeeperRegistryBase2_0_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit EventsMock_KeeperRegistryBase2_0_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitEventsMock_KeeperRegistryBase2_0_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit EventsMock_KeeperRegistryBase2_0_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitEventsMock_KeeperRegistryBase2_0_UpkeepOffchainConfigSet(uint256 id,bytes memory offchainConfig) public {
        emit EventsMock_KeeperRegistryBase2_0_UpkeepOffchainConfigSet(id,offchainConfig);
    }
    function emitEventsMock_KeeperRegistryBase2_0_UpkeepPaused(uint256 id) public {
        emit EventsMock_KeeperRegistryBase2_0_UpkeepPaused(id);
    }
    function emitEventsMock_KeeperRegistryBase2_0_UpkeepPerformed(uint256 id,bool success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment) public {
        emit EventsMock_KeeperRegistryBase2_0_UpkeepPerformed(id,success,checkBlockNumber,gasUsed,gasOverhead,totalPayment);
    }
    function emitEventsMock_KeeperRegistryBase2_0_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit EventsMock_KeeperRegistryBase2_0_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitEventsMock_KeeperRegistryBase2_0_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit EventsMock_KeeperRegistryBase2_0_UpkeepRegistered(id,executeGas,admin);
    }
    function emitEventsMock_KeeperRegistryBase2_0_UpkeepUnpaused(uint256 id) public {
        emit EventsMock_KeeperRegistryBase2_0_UpkeepUnpaused(id);
    }
    function emitEventsMock_KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferred(address from,address to) public {
        emit EventsMock_KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferred(from,to);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_ConfigSet(EventsMock_KeeperRegistryLogic1_3_Config memory config) public {
        emit EventsMock_KeeperRegistryLogic1_3_ConfigSet(config);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit EventsMock_KeeperRegistryLogic1_3_FundsAdded(id,from,amount);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit EventsMock_KeeperRegistryLogic1_3_FundsWithdrawn(id,amount,to);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_KeepersUpdated(address[] memory keepers,address[] memory payees) public {
        emit EventsMock_KeeperRegistryLogic1_3_KeepersUpdated(keepers,payees);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_OwnerFundsWithdrawn(uint96 amount) public {
        emit EventsMock_KeeperRegistryLogic1_3_OwnerFundsWithdrawn(amount);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_KeeperRegistryLogic1_3_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_OwnershipTransferred(address from,address to) public {
        emit EventsMock_KeeperRegistryLogic1_3_OwnershipTransferred(from,to);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_Paused(address account) public {
        emit EventsMock_KeeperRegistryLogic1_3_Paused(account);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_PayeeshipTransferRequested(address keeper,address from,address to) public {
        emit EventsMock_KeeperRegistryLogic1_3_PayeeshipTransferRequested(keeper,from,to);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_PayeeshipTransferred(address keeper,address from,address to) public {
        emit EventsMock_KeeperRegistryLogic1_3_PayeeshipTransferred(keeper,from,to);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_PaymentWithdrawn(address keeper,uint256 amount,address to,address payee) public {
        emit EventsMock_KeeperRegistryLogic1_3_PaymentWithdrawn(keeper,amount,to,payee);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_Unpaused(address account) public {
        emit EventsMock_KeeperRegistryLogic1_3_Unpaused(account);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_UpkeepAdminTransferRequested(uint256 id,address from,address to) public {
        emit EventsMock_KeeperRegistryLogic1_3_UpkeepAdminTransferRequested(id,from,to);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_UpkeepAdminTransferred(uint256 id,address from,address to) public {
        emit EventsMock_KeeperRegistryLogic1_3_UpkeepAdminTransferred(id,from,to);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit EventsMock_KeeperRegistryLogic1_3_UpkeepCanceled(id,atBlockHeight);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_UpkeepCheckDataUpdated(uint256 id,bytes memory newCheckData) public {
        emit EventsMock_KeeperRegistryLogic1_3_UpkeepCheckDataUpdated(id,newCheckData);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit EventsMock_KeeperRegistryLogic1_3_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit EventsMock_KeeperRegistryLogic1_3_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_UpkeepPaused(uint256 id) public {
        emit EventsMock_KeeperRegistryLogic1_3_UpkeepPaused(id);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_UpkeepPerformed(uint256 id,bool success,address from,uint96 payment,bytes memory performData) public {
        emit EventsMock_KeeperRegistryLogic1_3_UpkeepPerformed(id,success,from,payment,performData);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit EventsMock_KeeperRegistryLogic1_3_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit EventsMock_KeeperRegistryLogic1_3_UpkeepRegistered(id,executeGas,admin);
    }
    function emitEventsMock_KeeperRegistryLogic1_3_UpkeepUnpaused(uint256 id) public {
        emit EventsMock_KeeperRegistryLogic1_3_UpkeepUnpaused(id);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_CancelledUpkeepReport(uint256 id) public {
        emit EventsMock_KeeperRegistryLogic2_0_CancelledUpkeepReport(id);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit EventsMock_KeeperRegistryLogic2_0_FundsAdded(id,from,amount);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit EventsMock_KeeperRegistryLogic2_0_FundsWithdrawn(id,amount,to);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_InsufficientFundsUpkeepReport(uint256 id) public {
        emit EventsMock_KeeperRegistryLogic2_0_InsufficientFundsUpkeepReport(id);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_OwnerFundsWithdrawn(uint96 amount) public {
        emit EventsMock_KeeperRegistryLogic2_0_OwnerFundsWithdrawn(amount);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_KeeperRegistryLogic2_0_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_OwnershipTransferred(address from,address to) public {
        emit EventsMock_KeeperRegistryLogic2_0_OwnershipTransferred(from,to);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_Paused(address account) public {
        emit EventsMock_KeeperRegistryLogic2_0_Paused(account);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_PayeesUpdated(address[] memory transmitters,address[] memory payees) public {
        emit EventsMock_KeeperRegistryLogic2_0_PayeesUpdated(transmitters,payees);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_PayeeshipTransferRequested(address transmitter,address from,address to) public {
        emit EventsMock_KeeperRegistryLogic2_0_PayeeshipTransferRequested(transmitter,from,to);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_PayeeshipTransferred(address transmitter,address from,address to) public {
        emit EventsMock_KeeperRegistryLogic2_0_PayeeshipTransferred(transmitter,from,to);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_PaymentWithdrawn(address transmitter,uint256 amount,address to,address payee) public {
        emit EventsMock_KeeperRegistryLogic2_0_PaymentWithdrawn(transmitter,amount,to,payee);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_ReorgedUpkeepReport(uint256 id) public {
        emit EventsMock_KeeperRegistryLogic2_0_ReorgedUpkeepReport(id);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_StaleUpkeepReport(uint256 id) public {
        emit EventsMock_KeeperRegistryLogic2_0_StaleUpkeepReport(id);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_Unpaused(address account) public {
        emit EventsMock_KeeperRegistryLogic2_0_Unpaused(account);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_UpkeepAdminTransferRequested(uint256 id,address from,address to) public {
        emit EventsMock_KeeperRegistryLogic2_0_UpkeepAdminTransferRequested(id,from,to);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_UpkeepAdminTransferred(uint256 id,address from,address to) public {
        emit EventsMock_KeeperRegistryLogic2_0_UpkeepAdminTransferred(id,from,to);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit EventsMock_KeeperRegistryLogic2_0_UpkeepCanceled(id,atBlockHeight);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_UpkeepCheckDataUpdated(uint256 id,bytes memory newCheckData) public {
        emit EventsMock_KeeperRegistryLogic2_0_UpkeepCheckDataUpdated(id,newCheckData);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit EventsMock_KeeperRegistryLogic2_0_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit EventsMock_KeeperRegistryLogic2_0_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_UpkeepOffchainConfigSet(uint256 id,bytes memory offchainConfig) public {
        emit EventsMock_KeeperRegistryLogic2_0_UpkeepOffchainConfigSet(id,offchainConfig);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_UpkeepPaused(uint256 id) public {
        emit EventsMock_KeeperRegistryLogic2_0_UpkeepPaused(id);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_UpkeepPerformed(uint256 id,bool success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment) public {
        emit EventsMock_KeeperRegistryLogic2_0_UpkeepPerformed(id,success,checkBlockNumber,gasUsed,gasOverhead,totalPayment);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit EventsMock_KeeperRegistryLogic2_0_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit EventsMock_KeeperRegistryLogic2_0_UpkeepRegistered(id,executeGas,admin);
    }
    function emitEventsMock_KeeperRegistryLogic2_0_UpkeepUnpaused(uint256 id) public {
        emit EventsMock_KeeperRegistryLogic2_0_UpkeepUnpaused(id);
    }
    function emitEventsMock_LogEmitter_Log1(uint256 ibdrbwbi) public {
        emit EventsMock_LogEmitter_Log1(ibdrbwbi);
    }
    function emitEventsMock_LogEmitter_Log2(uint256 sopulnfk) public {
        emit EventsMock_LogEmitter_Log2(sopulnfk);
    }
    function emitEventsMock_LogEmitter_Log3(string memory sqdzywth) public {
        emit EventsMock_LogEmitter_Log3(sqdzywth);
    }
    function emitEventsMock_OCR2Abstract_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
        emit EventsMock_OCR2Abstract_ConfigSet(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
    }
    function emitEventsMock_OCR2Abstract_Transmitted(bytes32 configDigest,uint32 epoch) public {
        emit EventsMock_OCR2Abstract_Transmitted(configDigest,epoch);
    }
    function emitEventsMock_OCR2BaseUpgradeable_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
        emit EventsMock_OCR2BaseUpgradeable_ConfigSet(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
    }
    function emitEventsMock_OCR2BaseUpgradeable_Initialized(uint8 version) public {
        emit EventsMock_OCR2BaseUpgradeable_Initialized(version);
    }
    function emitEventsMock_OCR2BaseUpgradeable_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_OCR2BaseUpgradeable_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_OCR2BaseUpgradeable_OwnershipTransferred(address from,address to) public {
        emit EventsMock_OCR2BaseUpgradeable_OwnershipTransferred(from,to);
    }
    function emitEventsMock_OCR2BaseUpgradeable_Transmitted(bytes32 configDigest,uint32 epoch) public {
        emit EventsMock_OCR2BaseUpgradeable_Transmitted(configDigest,epoch);
    }
    function emitEventsMock_OVM_GasPriceOracle_DecimalsUpdated(uint256 wwbbrgcs) public {
        emit EventsMock_OVM_GasPriceOracle_DecimalsUpdated(wwbbrgcs);
    }
    function emitEventsMock_OVM_GasPriceOracle_GasPriceUpdated(uint256 bzwkuldh) public {
        emit EventsMock_OVM_GasPriceOracle_GasPriceUpdated(bzwkuldh);
    }
    function emitEventsMock_OVM_GasPriceOracle_L1BaseFeeUpdated(uint256 bwjknner) public {
        emit EventsMock_OVM_GasPriceOracle_L1BaseFeeUpdated(bwjknner);
    }
    function emitEventsMock_OVM_GasPriceOracle_OverheadUpdated(uint256 riulbpbh) public {
        emit EventsMock_OVM_GasPriceOracle_OverheadUpdated(riulbpbh);
    }
    function emitEventsMock_OVM_GasPriceOracle_OwnershipTransferred(address previousOwner,address newOwner) public {
        emit EventsMock_OVM_GasPriceOracle_OwnershipTransferred(previousOwner,newOwner);
    }
    function emitEventsMock_OVM_GasPriceOracle_ScalarUpdated(uint256 whsnlsbn) public {
        emit EventsMock_OVM_GasPriceOracle_ScalarUpdated(whsnlsbn);
    }
    function emitEventsMock_Ownable_OwnershipTransferred(address previousOwner,address newOwner) public {
        emit EventsMock_Ownable_OwnershipTransferred(previousOwner,newOwner);
    }
    function emitEventsMock_PausableUpgradeable_Initialized(uint8 version) public {
        emit EventsMock_PausableUpgradeable_Initialized(version);
    }
    function emitEventsMock_PausableUpgradeable_Paused(address account) public {
        emit EventsMock_PausableUpgradeable_Paused(account);
    }
    function emitEventsMock_PausableUpgradeable_Unpaused(address account) public {
        emit EventsMock_PausableUpgradeable_Unpaused(account);
    }
    function emitEventsMock_Pausable_Paused(address account) public {
        emit EventsMock_Pausable_Paused(account);
    }
    function emitEventsMock_Pausable_Unpaused(address account) public {
        emit EventsMock_Pausable_Unpaused(account);
    }
    function emitEventsMock_ProxyAdmin_OwnershipTransferred(address previousOwner,address newOwner) public {
        emit EventsMock_ProxyAdmin_OwnershipTransferred(previousOwner,newOwner);
    }
    function emitEventsMock_TransparentUpgradeableProxy_AdminChanged(address previousAdmin,address newAdmin) public {
        emit EventsMock_TransparentUpgradeableProxy_AdminChanged(previousAdmin,newAdmin);
    }
    function emitEventsMock_TransparentUpgradeableProxy_BeaconUpgraded(address beacon) public {
        emit EventsMock_TransparentUpgradeableProxy_BeaconUpgraded(beacon);
    }
    function emitEventsMock_TransparentUpgradeableProxy_Upgraded(address implementation) public {
        emit EventsMock_TransparentUpgradeableProxy_Upgraded(implementation);
    }
    function emitEventsMock_VRFConsumerBaseV2Upgradeable_Initialized(uint8 version) public {
        emit EventsMock_VRFConsumerBaseV2Upgradeable_Initialized(version);
    }
    function emitEventsMock_VRFConsumerV2UpgradeableExample_Initialized(uint8 version) public {
        emit EventsMock_VRFConsumerV2UpgradeableExample_Initialized(version);
    }
    function emitEventsMock_VRFCoordinatorMock_RandomnessRequest(address sender,bytes32 keyHash,uint256 seed) public {
        emit EventsMock_VRFCoordinatorMock_RandomnessRequest(sender,keyHash,seed);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_ConfigSet(uint16 minimumRequestConfirmations,uint32 maxGasLimit,uint32 stalenessSeconds,uint32 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,EventsMock_VRFCoordinatorV2TestHelper_FeeConfig memory feeConfig) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_ConfigSet(minimumRequestConfirmations,maxGasLimit,stalenessSeconds,gasAfterPaymentCalculation,fallbackWeiPerUnitLink,feeConfig);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_FundsRecovered(address to,uint256 amount) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_FundsRecovered(to,amount);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_OwnershipTransferred(address from,address to) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_OwnershipTransferred(from,to);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_ProvingKeyDeregistered(bytes32 keyHash,address oracle) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_ProvingKeyDeregistered(keyHash,oracle);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_ProvingKeyRegistered(bytes32 keyHash,address oracle) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_ProvingKeyRegistered(keyHash,oracle);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_RandomWordsFulfilled(uint256 requestId,uint256 outputSeed,uint96 payment,bool success) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_RandomWordsFulfilled(requestId,outputSeed,payment,success);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_RandomWordsRequested(bytes32 keyHash,uint256 requestId,uint256 preSeed,uint64 subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address sender) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_RandomWordsRequested(keyHash,requestId,preSeed,subId,minimumRequestConfirmations,callbackGasLimit,numWords,sender);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_SubscriptionCanceled(uint64 subId,address to,uint256 amount) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_SubscriptionCanceled(subId,to,amount);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_SubscriptionConsumerAdded(uint64 subId,address consumer) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_SubscriptionConsumerAdded(subId,consumer);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_SubscriptionConsumerRemoved(uint64 subId,address consumer) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_SubscriptionConsumerRemoved(subId,consumer);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_SubscriptionCreated(uint64 subId,address owner) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_SubscriptionCreated(subId,owner);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_SubscriptionFunded(uint64 subId,uint256 oldBalance,uint256 newBalance) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_SubscriptionFunded(subId,oldBalance,newBalance);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferRequested(uint64 subId,address from,address to) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferRequested(subId,from,to);
    }
    function emitEventsMock_VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferred(uint64 subId,address from,address to) public {
        emit EventsMock_VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferred(subId,from,to);
    }
    function emitEventsMock_VRFCoordinatorV2_ConfigSet(uint16 minimumRequestConfirmations,uint32 maxGasLimit,uint32 stalenessSeconds,uint32 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,EventsMock_VRFCoordinatorV2_FeeConfig memory feeConfig) public {
        emit EventsMock_VRFCoordinatorV2_ConfigSet(minimumRequestConfirmations,maxGasLimit,stalenessSeconds,gasAfterPaymentCalculation,fallbackWeiPerUnitLink,feeConfig);
    }
    function emitEventsMock_VRFCoordinatorV2_FundsRecovered(address to,uint256 amount) public {
        emit EventsMock_VRFCoordinatorV2_FundsRecovered(to,amount);
    }
    function emitEventsMock_VRFCoordinatorV2_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_VRFCoordinatorV2_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_VRFCoordinatorV2_OwnershipTransferred(address from,address to) public {
        emit EventsMock_VRFCoordinatorV2_OwnershipTransferred(from,to);
    }
    function emitEventsMock_VRFCoordinatorV2_ProvingKeyDeregistered(bytes32 keyHash,address oracle) public {
        emit EventsMock_VRFCoordinatorV2_ProvingKeyDeregistered(keyHash,oracle);
    }
    function emitEventsMock_VRFCoordinatorV2_ProvingKeyRegistered(bytes32 keyHash,address oracle) public {
        emit EventsMock_VRFCoordinatorV2_ProvingKeyRegistered(keyHash,oracle);
    }
    function emitEventsMock_VRFCoordinatorV2_RandomWordsFulfilled(uint256 requestId,uint256 outputSeed,uint96 payment,bool success) public {
        emit EventsMock_VRFCoordinatorV2_RandomWordsFulfilled(requestId,outputSeed,payment,success);
    }
    function emitEventsMock_VRFCoordinatorV2_RandomWordsRequested(bytes32 keyHash,uint256 requestId,uint256 preSeed,uint64 subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address sender) public {
        emit EventsMock_VRFCoordinatorV2_RandomWordsRequested(keyHash,requestId,preSeed,subId,minimumRequestConfirmations,callbackGasLimit,numWords,sender);
    }
    function emitEventsMock_VRFCoordinatorV2_SubscriptionCanceled(uint64 subId,address to,uint256 amount) public {
        emit EventsMock_VRFCoordinatorV2_SubscriptionCanceled(subId,to,amount);
    }
    function emitEventsMock_VRFCoordinatorV2_SubscriptionConsumerAdded(uint64 subId,address consumer) public {
        emit EventsMock_VRFCoordinatorV2_SubscriptionConsumerAdded(subId,consumer);
    }
    function emitEventsMock_VRFCoordinatorV2_SubscriptionConsumerRemoved(uint64 subId,address consumer) public {
        emit EventsMock_VRFCoordinatorV2_SubscriptionConsumerRemoved(subId,consumer);
    }
    function emitEventsMock_VRFCoordinatorV2_SubscriptionCreated(uint64 subId,address owner) public {
        emit EventsMock_VRFCoordinatorV2_SubscriptionCreated(subId,owner);
    }
    function emitEventsMock_VRFCoordinatorV2_SubscriptionFunded(uint64 subId,uint256 oldBalance,uint256 newBalance) public {
        emit EventsMock_VRFCoordinatorV2_SubscriptionFunded(subId,oldBalance,newBalance);
    }
    function emitEventsMock_VRFCoordinatorV2_SubscriptionOwnerTransferRequested(uint64 subId,address from,address to) public {
        emit EventsMock_VRFCoordinatorV2_SubscriptionOwnerTransferRequested(subId,from,to);
    }
    function emitEventsMock_VRFCoordinatorV2_SubscriptionOwnerTransferred(uint64 subId,address from,address to) public {
        emit EventsMock_VRFCoordinatorV2_SubscriptionOwnerTransferred(subId,from,to);
    }
    function emitEventsMock_VRFLoadTestExternalSubOwner_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_VRFLoadTestExternalSubOwner_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_VRFLoadTestExternalSubOwner_OwnershipTransferred(address from,address to) public {
        emit EventsMock_VRFLoadTestExternalSubOwner_OwnershipTransferred(from,to);
    }
    function emitEventsMock_VRFV2ProxyAdmin_OwnershipTransferred(address previousOwner,address newOwner) public {
        emit EventsMock_VRFV2ProxyAdmin_OwnershipTransferred(previousOwner,newOwner);
    }
    function emitEventsMock_VRFV2TransparentUpgradeableProxy_AdminChanged(address previousAdmin,address newAdmin) public {
        emit EventsMock_VRFV2TransparentUpgradeableProxy_AdminChanged(previousAdmin,newAdmin);
    }
    function emitEventsMock_VRFV2TransparentUpgradeableProxy_BeaconUpgraded(address beacon) public {
        emit EventsMock_VRFV2TransparentUpgradeableProxy_BeaconUpgraded(beacon);
    }
    function emitEventsMock_VRFV2TransparentUpgradeableProxy_Upgraded(address implementation) public {
        emit EventsMock_VRFV2TransparentUpgradeableProxy_Upgraded(implementation);
    }
    function emitEventsMock_VRFV2WrapperConsumerExample_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_VRFV2WrapperConsumerExample_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_VRFV2WrapperConsumerExample_OwnershipTransferred(address from,address to) public {
        emit EventsMock_VRFV2WrapperConsumerExample_OwnershipTransferred(from,to);
    }
    function emitEventsMock_VRFV2WrapperConsumerExample_WrappedRequestFulfilled(uint256 requestId,uint256[] memory randomWords,uint256 payment) public {
        emit EventsMock_VRFV2WrapperConsumerExample_WrappedRequestFulfilled(requestId,randomWords,payment);
    }
    function emitEventsMock_VRFV2WrapperConsumerExample_WrapperRequestMade(uint256 requestId,uint256 paid) public {
        emit EventsMock_VRFV2WrapperConsumerExample_WrapperRequestMade(requestId,paid);
    }
    function emitEventsMock_VRFV2Wrapper_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_VRFV2Wrapper_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_VRFV2Wrapper_OwnershipTransferred(address from,address to) public {
        emit EventsMock_VRFV2Wrapper_OwnershipTransferred(from,to);
    }
    function emitEventsMock_VRFV2Wrapper_WrapperFulfillmentFailed(uint256 requestId,address consumer) public {
        emit EventsMock_VRFV2Wrapper_WrapperFulfillmentFailed(requestId,consumer);
    }
    function emitEventsMock_VerifierProxy_AccessControllerSet(address oldAccessController,address newAccessController) public {
        emit EventsMock_VerifierProxy_AccessControllerSet(oldAccessController,newAccessController);
    }
    function emitEventsMock_VerifierProxy_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_VerifierProxy_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_VerifierProxy_OwnershipTransferred(address from,address to) public {
        emit EventsMock_VerifierProxy_OwnershipTransferred(from,to);
    }
    function emitEventsMock_VerifierProxy_VerifierSet(bytes32 oldConfigDigest,bytes32 newConfigDigest,address verifierAddress) public {
        emit EventsMock_VerifierProxy_VerifierSet(oldConfigDigest,newConfigDigest,verifierAddress);
    }
    function emitEventsMock_VerifierProxy_VerifierUnset(bytes32 configDigest,address verifierAddress) public {
        emit EventsMock_VerifierProxy_VerifierUnset(configDigest,verifierAddress);
    }
    function emitEventsMock_Verifier_ConfigActivated(bytes32 feedId,bytes32 configDigest) public {
        emit EventsMock_Verifier_ConfigActivated(feedId,configDigest);
    }
    function emitEventsMock_Verifier_ConfigDeactivated(bytes32 feedId,bytes32 configDigest) public {
        emit EventsMock_Verifier_ConfigDeactivated(feedId,configDigest);
    }
    function emitEventsMock_Verifier_ConfigSet(bytes32 feedId,uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,bytes32[] memory offchainTransmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
        emit EventsMock_Verifier_ConfigSet(feedId,previousConfigBlockNumber,configDigest,configCount,signers,offchainTransmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
    }
    function emitEventsMock_Verifier_OwnershipTransferRequested(address from,address to) public {
        emit EventsMock_Verifier_OwnershipTransferRequested(from,to);
    }
    function emitEventsMock_Verifier_OwnershipTransferred(address from,address to) public {
        emit EventsMock_Verifier_OwnershipTransferred(from,to);
    }
    function emitEventsMock_Verifier_ReportVerified(bytes32 feedId,bytes32 reportHash,address requester) public {
        emit EventsMock_Verifier_ReportVerified(feedId,reportHash,requester);
    }
    function emitFunctionsBillingRegistry_AuthorizedSendersChanged(address[] memory senders,address changedBy) public {
        emit FunctionsBillingRegistry_AuthorizedSendersChanged(senders,changedBy);
    }
    function emitFunctionsBillingRegistry_BillingEnd(bytes32 requestId,uint64 subscriptionId,uint96 signerPayment,uint96 transmitterPayment,uint96 totalCost,bool success) public {
        emit FunctionsBillingRegistry_BillingEnd(requestId,subscriptionId,signerPayment,transmitterPayment,totalCost,success);
    }
    function emitFunctionsBillingRegistry_BillingStart(bytes32 requestId,FunctionsBillingRegistry_Commitment memory commitment) public {
        emit FunctionsBillingRegistry_BillingStart(requestId,commitment);
    }
    function emitFunctionsBillingRegistry_ConfigSet(uint32 maxGasLimit,uint32 stalenessSeconds,uint256 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,uint32 gasOverhead) public {
        emit FunctionsBillingRegistry_ConfigSet(maxGasLimit,stalenessSeconds,gasAfterPaymentCalculation,fallbackWeiPerUnitLink,gasOverhead);
    }
    function emitFunctionsBillingRegistry_FundsRecovered(address to,uint256 amount) public {
        emit FunctionsBillingRegistry_FundsRecovered(to,amount);
    }
    function emitFunctionsBillingRegistry_Initialized(uint8 version) public {
        emit FunctionsBillingRegistry_Initialized(version);
    }
    function emitFunctionsBillingRegistry_OwnershipTransferRequested(address from,address to) public {
        emit FunctionsBillingRegistry_OwnershipTransferRequested(from,to);
    }
    function emitFunctionsBillingRegistry_OwnershipTransferred(address from,address to) public {
        emit FunctionsBillingRegistry_OwnershipTransferred(from,to);
    }
    function emitFunctionsBillingRegistry_Paused(address account) public {
        emit FunctionsBillingRegistry_Paused(account);
    }
    function emitFunctionsBillingRegistry_RequestTimedOut(bytes32 requestId) public {
        emit FunctionsBillingRegistry_RequestTimedOut(requestId);
    }
    function emitFunctionsBillingRegistry_SubscriptionCanceled(uint64 subscriptionId,address to,uint256 amount) public {
        emit FunctionsBillingRegistry_SubscriptionCanceled(subscriptionId,to,amount);
    }
    function emitFunctionsBillingRegistry_SubscriptionConsumerAdded(uint64 subscriptionId,address consumer) public {
        emit FunctionsBillingRegistry_SubscriptionConsumerAdded(subscriptionId,consumer);
    }
    function emitFunctionsBillingRegistry_SubscriptionConsumerRemoved(uint64 subscriptionId,address consumer) public {
        emit FunctionsBillingRegistry_SubscriptionConsumerRemoved(subscriptionId,consumer);
    }
    function emitFunctionsBillingRegistry_SubscriptionCreated(uint64 subscriptionId,address owner) public {
        emit FunctionsBillingRegistry_SubscriptionCreated(subscriptionId,owner);
    }
    function emitFunctionsBillingRegistry_SubscriptionFunded(uint64 subscriptionId,uint256 oldBalance,uint256 newBalance) public {
        emit FunctionsBillingRegistry_SubscriptionFunded(subscriptionId,oldBalance,newBalance);
    }
    function emitFunctionsBillingRegistry_SubscriptionOwnerTransferRequested(uint64 subscriptionId,address from,address to) public {
        emit FunctionsBillingRegistry_SubscriptionOwnerTransferRequested(subscriptionId,from,to);
    }
    function emitFunctionsBillingRegistry_SubscriptionOwnerTransferred(uint64 subscriptionId,address from,address to) public {
        emit FunctionsBillingRegistry_SubscriptionOwnerTransferred(subscriptionId,from,to);
    }
    function emitFunctionsBillingRegistry_Unpaused(address account) public {
        emit FunctionsBillingRegistry_Unpaused(account);
    }
    function emitFunctionsBillingRegistryWithInit_AuthorizedSendersChanged(address[] memory senders,address changedBy) public {
        emit FunctionsBillingRegistryWithInit_AuthorizedSendersChanged(senders,changedBy);
    }
    function emitFunctionsBillingRegistryWithInit_BillingEnd(bytes32 requestId,uint64 subscriptionId,uint96 signerPayment,uint96 transmitterPayment,uint96 totalCost,bool success) public {
        emit FunctionsBillingRegistryWithInit_BillingEnd(requestId,subscriptionId,signerPayment,transmitterPayment,totalCost,success);
    }
    function emitFunctionsBillingRegistryWithInit_BillingStart(bytes32 requestId,FunctionsBillingRegistryWithInit_Commitment memory commitment) public {
        emit FunctionsBillingRegistryWithInit_BillingStart(requestId,commitment);
    }
    function emitFunctionsBillingRegistryWithInit_ConfigSet(uint32 maxGasLimit,uint32 stalenessSeconds,uint256 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,uint32 gasOverhead) public {
        emit FunctionsBillingRegistryWithInit_ConfigSet(maxGasLimit,stalenessSeconds,gasAfterPaymentCalculation,fallbackWeiPerUnitLink,gasOverhead);
    }
    function emitFunctionsBillingRegistryWithInit_FundsRecovered(address to,uint256 amount) public {
        emit FunctionsBillingRegistryWithInit_FundsRecovered(to,amount);
    }
    function emitFunctionsBillingRegistryWithInit_Initialized(uint8 version) public {
        emit FunctionsBillingRegistryWithInit_Initialized(version);
    }
    function emitFunctionsBillingRegistryWithInit_OwnershipTransferRequested(address from,address to) public {
        emit FunctionsBillingRegistryWithInit_OwnershipTransferRequested(from,to);
    }
    function emitFunctionsBillingRegistryWithInit_OwnershipTransferred(address from,address to) public {
        emit FunctionsBillingRegistryWithInit_OwnershipTransferred(from,to);
    }
    function emitFunctionsBillingRegistryWithInit_Paused(address account) public {
        emit FunctionsBillingRegistryWithInit_Paused(account);
    }
    function emitFunctionsBillingRegistryWithInit_RequestTimedOut(bytes32 requestId) public {
        emit FunctionsBillingRegistryWithInit_RequestTimedOut(requestId);
    }
    function emitFunctionsBillingRegistryWithInit_SubscriptionCanceled(uint64 subscriptionId,address to,uint256 amount) public {
        emit FunctionsBillingRegistryWithInit_SubscriptionCanceled(subscriptionId,to,amount);
    }
    function emitFunctionsBillingRegistryWithInit_SubscriptionConsumerAdded(uint64 subscriptionId,address consumer) public {
        emit FunctionsBillingRegistryWithInit_SubscriptionConsumerAdded(subscriptionId,consumer);
    }
    function emitFunctionsBillingRegistryWithInit_SubscriptionConsumerRemoved(uint64 subscriptionId,address consumer) public {
        emit FunctionsBillingRegistryWithInit_SubscriptionConsumerRemoved(subscriptionId,consumer);
    }
    function emitFunctionsBillingRegistryWithInit_SubscriptionCreated(uint64 subscriptionId,address owner) public {
        emit FunctionsBillingRegistryWithInit_SubscriptionCreated(subscriptionId,owner);
    }
    function emitFunctionsBillingRegistryWithInit_SubscriptionFunded(uint64 subscriptionId,uint256 oldBalance,uint256 newBalance) public {
        emit FunctionsBillingRegistryWithInit_SubscriptionFunded(subscriptionId,oldBalance,newBalance);
    }
    function emitFunctionsBillingRegistryWithInit_SubscriptionOwnerTransferRequested(uint64 subscriptionId,address from,address to) public {
        emit FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferRequested(subscriptionId,from,to);
    }
    function emitFunctionsBillingRegistryWithInit_SubscriptionOwnerTransferred(uint64 subscriptionId,address from,address to) public {
        emit FunctionsBillingRegistryWithInit_SubscriptionOwnerTransferred(subscriptionId,from,to);
    }
    function emitFunctionsBillingRegistryWithInit_Unpaused(address account) public {
        emit FunctionsBillingRegistryWithInit_Unpaused(account);
    }
    function emitFunctionsClient_RequestFulfilled(bytes32 id) public {
        emit FunctionsClient_RequestFulfilled(id);
    }
    function emitFunctionsClient_RequestSent(bytes32 id) public {
        emit FunctionsClient_RequestSent(id);
    }
    function emitFunctionsClientExample_OwnershipTransferRequested(address from,address to) public {
        emit FunctionsClientExample_OwnershipTransferRequested(from,to);
    }
    function emitFunctionsClientExample_OwnershipTransferred(address from,address to) public {
        emit FunctionsClientExample_OwnershipTransferred(from,to);
    }
    function emitFunctionsClientExample_RequestFulfilled(bytes32 id) public {
        emit FunctionsClientExample_RequestFulfilled(id);
    }
    function emitFunctionsClientExample_RequestSent(bytes32 id) public {
        emit FunctionsClientExample_RequestSent(id);
    }
    function emitFunctionsOracle_AuthorizedSendersActive(address account) public {
        emit FunctionsOracle_AuthorizedSendersActive(account);
    }
    function emitFunctionsOracle_AuthorizedSendersChanged(address[] memory senders,address changedBy) public {
        emit FunctionsOracle_AuthorizedSendersChanged(senders,changedBy);
    }
    function emitFunctionsOracle_AuthorizedSendersDeactive(address account) public {
        emit FunctionsOracle_AuthorizedSendersDeactive(account);
    }
    function emitFunctionsOracle_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
        emit FunctionsOracle_ConfigSet(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
    }
    function emitFunctionsOracle_Initialized(uint8 version) public {
        emit FunctionsOracle_Initialized(version);
    }
    function emitFunctionsOracle_InvalidRequestID(bytes32 requestId) public {
        emit FunctionsOracle_InvalidRequestID(requestId);
    }
    function emitFunctionsOracle_OracleRequest(bytes32 requestId,address requestingContract,address requestInitiator,uint64 subscriptionId,address subscriptionOwner,bytes memory data) public {
        emit FunctionsOracle_OracleRequest(requestId,requestingContract,requestInitiator,subscriptionId,subscriptionOwner,data);
    }
    function emitFunctionsOracle_OracleResponse(bytes32 requestId) public {
        emit FunctionsOracle_OracleResponse(requestId);
    }
    function emitFunctionsOracle_OwnershipTransferRequested(address from,address to) public {
        emit FunctionsOracle_OwnershipTransferRequested(from,to);
    }
    function emitFunctionsOracle_OwnershipTransferred(address from,address to) public {
        emit FunctionsOracle_OwnershipTransferred(from,to);
    }
    function emitFunctionsOracle_Transmitted(bytes32 configDigest,uint32 epoch) public {
        emit FunctionsOracle_Transmitted(configDigest,epoch);
    }
    function emitFunctionsOracle_UserCallbackError(bytes32 requestId,string memory reason) public {
        emit FunctionsOracle_UserCallbackError(requestId,reason);
    }
    function emitFunctionsOracle_UserCallbackRawError(bytes32 requestId,bytes memory lowLevelData) public {
        emit FunctionsOracle_UserCallbackRawError(requestId,lowLevelData);
    }
    function emitFunctionsOracleWithInit_AuthorizedSendersActive(address account) public {
        emit FunctionsOracleWithInit_AuthorizedSendersActive(account);
    }
    function emitFunctionsOracleWithInit_AuthorizedSendersChanged(address[] memory senders,address changedBy) public {
        emit FunctionsOracleWithInit_AuthorizedSendersChanged(senders,changedBy);
    }
    function emitFunctionsOracleWithInit_AuthorizedSendersDeactive(address account) public {
        emit FunctionsOracleWithInit_AuthorizedSendersDeactive(account);
    }
    function emitFunctionsOracleWithInit_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
        emit FunctionsOracleWithInit_ConfigSet(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
    }
    function emitFunctionsOracleWithInit_Initialized(uint8 version) public {
        emit FunctionsOracleWithInit_Initialized(version);
    }
    function emitFunctionsOracleWithInit_InvalidRequestID(bytes32 requestId) public {
        emit FunctionsOracleWithInit_InvalidRequestID(requestId);
    }
    function emitFunctionsOracleWithInit_OracleRequest(bytes32 requestId,address requestingContract,address requestInitiator,uint64 subscriptionId,address subscriptionOwner,bytes memory data) public {
        emit FunctionsOracleWithInit_OracleRequest(requestId,requestingContract,requestInitiator,subscriptionId,subscriptionOwner,data);
    }
    function emitFunctionsOracleWithInit_OracleResponse(bytes32 requestId) public {
        emit FunctionsOracleWithInit_OracleResponse(requestId);
    }
    function emitFunctionsOracleWithInit_OwnershipTransferRequested(address from,address to) public {
        emit FunctionsOracleWithInit_OwnershipTransferRequested(from,to);
    }
    function emitFunctionsOracleWithInit_OwnershipTransferred(address from,address to) public {
        emit FunctionsOracleWithInit_OwnershipTransferred(from,to);
    }
    function emitFunctionsOracleWithInit_Transmitted(bytes32 configDigest,uint32 epoch) public {
        emit FunctionsOracleWithInit_Transmitted(configDigest,epoch);
    }
    function emitFunctionsOracleWithInit_UserCallbackError(bytes32 requestId,string memory reason) public {
        emit FunctionsOracleWithInit_UserCallbackError(requestId,reason);
    }
    function emitFunctionsOracleWithInit_UserCallbackRawError(bytes32 requestId,bytes memory lowLevelData) public {
        emit FunctionsOracleWithInit_UserCallbackRawError(requestId,lowLevelData);
    }
    function emitInitializable_Initialized(uint8 version) public {
        emit Initializable_Initialized(version);
    }
    function emitKeeperRegistry1_2_ConfigSet(KeeperRegistry1_2_Config memory config) public {
        emit KeeperRegistry1_2_ConfigSet(config);
    }
    function emitKeeperRegistry1_2_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit KeeperRegistry1_2_FundsAdded(id,from,amount);
    }
    function emitKeeperRegistry1_2_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit KeeperRegistry1_2_FundsWithdrawn(id,amount,to);
    }
    function emitKeeperRegistry1_2_KeepersUpdated(address[] memory keepers,address[] memory payees) public {
        emit KeeperRegistry1_2_KeepersUpdated(keepers,payees);
    }
    function emitKeeperRegistry1_2_OwnerFundsWithdrawn(uint96 amount) public {
        emit KeeperRegistry1_2_OwnerFundsWithdrawn(amount);
    }
    function emitKeeperRegistry1_2_OwnershipTransferRequested(address from,address to) public {
        emit KeeperRegistry1_2_OwnershipTransferRequested(from,to);
    }
    function emitKeeperRegistry1_2_OwnershipTransferred(address from,address to) public {
        emit KeeperRegistry1_2_OwnershipTransferred(from,to);
    }
    function emitKeeperRegistry1_2_Paused(address account) public {
        emit KeeperRegistry1_2_Paused(account);
    }
    function emitKeeperRegistry1_2_PayeeshipTransferRequested(address keeper,address from,address to) public {
        emit KeeperRegistry1_2_PayeeshipTransferRequested(keeper,from,to);
    }
    function emitKeeperRegistry1_2_PayeeshipTransferred(address keeper,address from,address to) public {
        emit KeeperRegistry1_2_PayeeshipTransferred(keeper,from,to);
    }
    function emitKeeperRegistry1_2_PaymentWithdrawn(address keeper,uint256 amount,address to,address payee) public {
        emit KeeperRegistry1_2_PaymentWithdrawn(keeper,amount,to,payee);
    }
    function emitKeeperRegistry1_2_Unpaused(address account) public {
        emit KeeperRegistry1_2_Unpaused(account);
    }
    function emitKeeperRegistry1_2_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit KeeperRegistry1_2_UpkeepCanceled(id,atBlockHeight);
    }
    function emitKeeperRegistry1_2_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit KeeperRegistry1_2_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitKeeperRegistry1_2_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit KeeperRegistry1_2_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitKeeperRegistry1_2_UpkeepPerformed(uint256 id,bool success,address from,uint96 payment,bytes memory performData) public {
        emit KeeperRegistry1_2_UpkeepPerformed(id,success,from,payment,performData);
    }
    function emitKeeperRegistry1_2_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit KeeperRegistry1_2_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitKeeperRegistry1_2_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit KeeperRegistry1_2_UpkeepRegistered(id,executeGas,admin);
    }
    function emitKeeperRegistry1_3_ConfigSet(KeeperRegistry1_3_Config memory config) public {
        emit KeeperRegistry1_3_ConfigSet(config);
    }
    function emitKeeperRegistry1_3_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit KeeperRegistry1_3_FundsAdded(id,from,amount);
    }
    function emitKeeperRegistry1_3_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit KeeperRegistry1_3_FundsWithdrawn(id,amount,to);
    }
    function emitKeeperRegistry1_3_KeepersUpdated(address[] memory keepers,address[] memory payees) public {
        emit KeeperRegistry1_3_KeepersUpdated(keepers,payees);
    }
    function emitKeeperRegistry1_3_OwnerFundsWithdrawn(uint96 amount) public {
        emit KeeperRegistry1_3_OwnerFundsWithdrawn(amount);
    }
    function emitKeeperRegistry1_3_OwnershipTransferRequested(address from,address to) public {
        emit KeeperRegistry1_3_OwnershipTransferRequested(from,to);
    }
    function emitKeeperRegistry1_3_OwnershipTransferred(address from,address to) public {
        emit KeeperRegistry1_3_OwnershipTransferred(from,to);
    }
    function emitKeeperRegistry1_3_Paused(address account) public {
        emit KeeperRegistry1_3_Paused(account);
    }
    function emitKeeperRegistry1_3_PayeeshipTransferRequested(address keeper,address from,address to) public {
        emit KeeperRegistry1_3_PayeeshipTransferRequested(keeper,from,to);
    }
    function emitKeeperRegistry1_3_PayeeshipTransferred(address keeper,address from,address to) public {
        emit KeeperRegistry1_3_PayeeshipTransferred(keeper,from,to);
    }
    function emitKeeperRegistry1_3_PaymentWithdrawn(address keeper,uint256 amount,address to,address payee) public {
        emit KeeperRegistry1_3_PaymentWithdrawn(keeper,amount,to,payee);
    }
    function emitKeeperRegistry1_3_Unpaused(address account) public {
        emit KeeperRegistry1_3_Unpaused(account);
    }
    function emitKeeperRegistry1_3_UpkeepAdminTransferRequested(uint256 id,address from,address to) public {
        emit KeeperRegistry1_3_UpkeepAdminTransferRequested(id,from,to);
    }
    function emitKeeperRegistry1_3_UpkeepAdminTransferred(uint256 id,address from,address to) public {
        emit KeeperRegistry1_3_UpkeepAdminTransferred(id,from,to);
    }
    function emitKeeperRegistry1_3_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit KeeperRegistry1_3_UpkeepCanceled(id,atBlockHeight);
    }
    function emitKeeperRegistry1_3_UpkeepCheckDataUpdated(uint256 id,bytes memory newCheckData) public {
        emit KeeperRegistry1_3_UpkeepCheckDataUpdated(id,newCheckData);
    }
    function emitKeeperRegistry1_3_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit KeeperRegistry1_3_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitKeeperRegistry1_3_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit KeeperRegistry1_3_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitKeeperRegistry1_3_UpkeepPaused(uint256 id) public {
        emit KeeperRegistry1_3_UpkeepPaused(id);
    }
    function emitKeeperRegistry1_3_UpkeepPerformed(uint256 id,bool success,address from,uint96 payment,bytes memory performData) public {
        emit KeeperRegistry1_3_UpkeepPerformed(id,success,from,payment,performData);
    }
    function emitKeeperRegistry1_3_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit KeeperRegistry1_3_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitKeeperRegistry1_3_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit KeeperRegistry1_3_UpkeepRegistered(id,executeGas,admin);
    }
    function emitKeeperRegistry1_3_UpkeepUnpaused(uint256 id) public {
        emit KeeperRegistry1_3_UpkeepUnpaused(id);
    }
    function emitKeeperRegistry2_0_CancelledUpkeepReport(uint256 id) public {
        emit KeeperRegistry2_0_CancelledUpkeepReport(id);
    }
    function emitKeeperRegistry2_0_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
        emit KeeperRegistry2_0_ConfigSet(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
    }
    function emitKeeperRegistry2_0_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit KeeperRegistry2_0_FundsAdded(id,from,amount);
    }
    function emitKeeperRegistry2_0_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit KeeperRegistry2_0_FundsWithdrawn(id,amount,to);
    }
    function emitKeeperRegistry2_0_InsufficientFundsUpkeepReport(uint256 id) public {
        emit KeeperRegistry2_0_InsufficientFundsUpkeepReport(id);
    }
    function emitKeeperRegistry2_0_OwnerFundsWithdrawn(uint96 amount) public {
        emit KeeperRegistry2_0_OwnerFundsWithdrawn(amount);
    }
    function emitKeeperRegistry2_0_OwnershipTransferRequested(address from,address to) public {
        emit KeeperRegistry2_0_OwnershipTransferRequested(from,to);
    }
    function emitKeeperRegistry2_0_OwnershipTransferred(address from,address to) public {
        emit KeeperRegistry2_0_OwnershipTransferred(from,to);
    }
    function emitKeeperRegistry2_0_Paused(address account) public {
        emit KeeperRegistry2_0_Paused(account);
    }
    function emitKeeperRegistry2_0_PayeesUpdated(address[] memory transmitters,address[] memory payees) public {
        emit KeeperRegistry2_0_PayeesUpdated(transmitters,payees);
    }
    function emitKeeperRegistry2_0_PayeeshipTransferRequested(address transmitter,address from,address to) public {
        emit KeeperRegistry2_0_PayeeshipTransferRequested(transmitter,from,to);
    }
    function emitKeeperRegistry2_0_PayeeshipTransferred(address transmitter,address from,address to) public {
        emit KeeperRegistry2_0_PayeeshipTransferred(transmitter,from,to);
    }
    function emitKeeperRegistry2_0_PaymentWithdrawn(address transmitter,uint256 amount,address to,address payee) public {
        emit KeeperRegistry2_0_PaymentWithdrawn(transmitter,amount,to,payee);
    }
    function emitKeeperRegistry2_0_ReorgedUpkeepReport(uint256 id) public {
        emit KeeperRegistry2_0_ReorgedUpkeepReport(id);
    }
    function emitKeeperRegistry2_0_StaleUpkeepReport(uint256 id) public {
        emit KeeperRegistry2_0_StaleUpkeepReport(id);
    }
    function emitKeeperRegistry2_0_Transmitted(bytes32 configDigest,uint32 epoch) public {
        emit KeeperRegistry2_0_Transmitted(configDigest,epoch);
    }
    function emitKeeperRegistry2_0_Unpaused(address account) public {
        emit KeeperRegistry2_0_Unpaused(account);
    }
    function emitKeeperRegistry2_0_UpkeepAdminTransferRequested(uint256 id,address from,address to) public {
        emit KeeperRegistry2_0_UpkeepAdminTransferRequested(id,from,to);
    }
    function emitKeeperRegistry2_0_UpkeepAdminTransferred(uint256 id,address from,address to) public {
        emit KeeperRegistry2_0_UpkeepAdminTransferred(id,from,to);
    }
    function emitKeeperRegistry2_0_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit KeeperRegistry2_0_UpkeepCanceled(id,atBlockHeight);
    }
    function emitKeeperRegistry2_0_UpkeepCheckDataUpdated(uint256 id,bytes memory newCheckData) public {
        emit KeeperRegistry2_0_UpkeepCheckDataUpdated(id,newCheckData);
    }
    function emitKeeperRegistry2_0_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit KeeperRegistry2_0_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitKeeperRegistry2_0_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit KeeperRegistry2_0_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitKeeperRegistry2_0_UpkeepOffchainConfigSet(uint256 id,bytes memory offchainConfig) public {
        emit KeeperRegistry2_0_UpkeepOffchainConfigSet(id,offchainConfig);
    }
    function emitKeeperRegistry2_0_UpkeepPaused(uint256 id) public {
        emit KeeperRegistry2_0_UpkeepPaused(id);
    }
    function emitKeeperRegistry2_0_UpkeepPerformed(uint256 id,bool success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment) public {
        emit KeeperRegistry2_0_UpkeepPerformed(id,success,checkBlockNumber,gasUsed,gasOverhead,totalPayment);
    }
    function emitKeeperRegistry2_0_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit KeeperRegistry2_0_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitKeeperRegistry2_0_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit KeeperRegistry2_0_UpkeepRegistered(id,executeGas,admin);
    }
    function emitKeeperRegistry2_0_UpkeepUnpaused(uint256 id) public {
        emit KeeperRegistry2_0_UpkeepUnpaused(id);
    }
    function emitKeeperRegistryBase1_3_ConfigSet(KeeperRegistryBase1_3_Config memory config) public {
        emit KeeperRegistryBase1_3_ConfigSet(config);
    }
    function emitKeeperRegistryBase1_3_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit KeeperRegistryBase1_3_FundsAdded(id,from,amount);
    }
    function emitKeeperRegistryBase1_3_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit KeeperRegistryBase1_3_FundsWithdrawn(id,amount,to);
    }
    function emitKeeperRegistryBase1_3_KeepersUpdated(address[] memory keepers,address[] memory payees) public {
        emit KeeperRegistryBase1_3_KeepersUpdated(keepers,payees);
    }
    function emitKeeperRegistryBase1_3_OwnerFundsWithdrawn(uint96 amount) public {
        emit KeeperRegistryBase1_3_OwnerFundsWithdrawn(amount);
    }
    function emitKeeperRegistryBase1_3_OwnershipTransferRequested(address from,address to) public {
        emit KeeperRegistryBase1_3_OwnershipTransferRequested(from,to);
    }
    function emitKeeperRegistryBase1_3_OwnershipTransferred(address from,address to) public {
        emit KeeperRegistryBase1_3_OwnershipTransferred(from,to);
    }
    function emitKeeperRegistryBase1_3_Paused(address account) public {
        emit KeeperRegistryBase1_3_Paused(account);
    }
    function emitKeeperRegistryBase1_3_PayeeshipTransferRequested(address keeper,address from,address to) public {
        emit KeeperRegistryBase1_3_PayeeshipTransferRequested(keeper,from,to);
    }
    function emitKeeperRegistryBase1_3_PayeeshipTransferred(address keeper,address from,address to) public {
        emit KeeperRegistryBase1_3_PayeeshipTransferred(keeper,from,to);
    }
    function emitKeeperRegistryBase1_3_PaymentWithdrawn(address keeper,uint256 amount,address to,address payee) public {
        emit KeeperRegistryBase1_3_PaymentWithdrawn(keeper,amount,to,payee);
    }
    function emitKeeperRegistryBase1_3_Unpaused(address account) public {
        emit KeeperRegistryBase1_3_Unpaused(account);
    }
    function emitKeeperRegistryBase1_3_UpkeepAdminTransferRequested(uint256 id,address from,address to) public {
        emit KeeperRegistryBase1_3_UpkeepAdminTransferRequested(id,from,to);
    }
    function emitKeeperRegistryBase1_3_UpkeepAdminTransferred(uint256 id,address from,address to) public {
        emit KeeperRegistryBase1_3_UpkeepAdminTransferred(id,from,to);
    }
    function emitKeeperRegistryBase1_3_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit KeeperRegistryBase1_3_UpkeepCanceled(id,atBlockHeight);
    }
    function emitKeeperRegistryBase1_3_UpkeepCheckDataUpdated(uint256 id,bytes memory newCheckData) public {
        emit KeeperRegistryBase1_3_UpkeepCheckDataUpdated(id,newCheckData);
    }
    function emitKeeperRegistryBase1_3_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit KeeperRegistryBase1_3_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitKeeperRegistryBase1_3_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit KeeperRegistryBase1_3_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitKeeperRegistryBase1_3_UpkeepPaused(uint256 id) public {
        emit KeeperRegistryBase1_3_UpkeepPaused(id);
    }
    function emitKeeperRegistryBase1_3_UpkeepPerformed(uint256 id,bool success,address from,uint96 payment,bytes memory performData) public {
        emit KeeperRegistryBase1_3_UpkeepPerformed(id,success,from,payment,performData);
    }
    function emitKeeperRegistryBase1_3_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit KeeperRegistryBase1_3_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitKeeperRegistryBase1_3_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit KeeperRegistryBase1_3_UpkeepRegistered(id,executeGas,admin);
    }
    function emitKeeperRegistryBase1_3_UpkeepUnpaused(uint256 id) public {
        emit KeeperRegistryBase1_3_UpkeepUnpaused(id);
    }
    function emitKeeperRegistryBase2_0_CancelledUpkeepReport(uint256 id) public {
        emit KeeperRegistryBase2_0_CancelledUpkeepReport(id);
    }
    function emitKeeperRegistryBase2_0_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit KeeperRegistryBase2_0_FundsAdded(id,from,amount);
    }
    function emitKeeperRegistryBase2_0_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit KeeperRegistryBase2_0_FundsWithdrawn(id,amount,to);
    }
    function emitKeeperRegistryBase2_0_InsufficientFundsUpkeepReport(uint256 id) public {
        emit KeeperRegistryBase2_0_InsufficientFundsUpkeepReport(id);
    }
    function emitKeeperRegistryBase2_0_OwnerFundsWithdrawn(uint96 amount) public {
        emit KeeperRegistryBase2_0_OwnerFundsWithdrawn(amount);
    }
    function emitKeeperRegistryBase2_0_OwnershipTransferRequested(address from,address to) public {
        emit KeeperRegistryBase2_0_OwnershipTransferRequested(from,to);
    }
    function emitKeeperRegistryBase2_0_OwnershipTransferred(address from,address to) public {
        emit KeeperRegistryBase2_0_OwnershipTransferred(from,to);
    }
    function emitKeeperRegistryBase2_0_Paused(address account) public {
        emit KeeperRegistryBase2_0_Paused(account);
    }
    function emitKeeperRegistryBase2_0_PayeesUpdated(address[] memory transmitters,address[] memory payees) public {
        emit KeeperRegistryBase2_0_PayeesUpdated(transmitters,payees);
    }
    function emitKeeperRegistryBase2_0_PayeeshipTransferRequested(address transmitter,address from,address to) public {
        emit KeeperRegistryBase2_0_PayeeshipTransferRequested(transmitter,from,to);
    }
    function emitKeeperRegistryBase2_0_PayeeshipTransferred(address transmitter,address from,address to) public {
        emit KeeperRegistryBase2_0_PayeeshipTransferred(transmitter,from,to);
    }
    function emitKeeperRegistryBase2_0_PaymentWithdrawn(address transmitter,uint256 amount,address to,address payee) public {
        emit KeeperRegistryBase2_0_PaymentWithdrawn(transmitter,amount,to,payee);
    }
    function emitKeeperRegistryBase2_0_ReorgedUpkeepReport(uint256 id) public {
        emit KeeperRegistryBase2_0_ReorgedUpkeepReport(id);
    }
    function emitKeeperRegistryBase2_0_StaleUpkeepReport(uint256 id) public {
        emit KeeperRegistryBase2_0_StaleUpkeepReport(id);
    }
    function emitKeeperRegistryBase2_0_Unpaused(address account) public {
        emit KeeperRegistryBase2_0_Unpaused(account);
    }
    function emitKeeperRegistryBase2_0_UpkeepAdminTransferRequested(uint256 id,address from,address to) public {
        emit KeeperRegistryBase2_0_UpkeepAdminTransferRequested(id,from,to);
    }
    function emitKeeperRegistryBase2_0_UpkeepAdminTransferred(uint256 id,address from,address to) public {
        emit KeeperRegistryBase2_0_UpkeepAdminTransferred(id,from,to);
    }
    function emitKeeperRegistryBase2_0_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit KeeperRegistryBase2_0_UpkeepCanceled(id,atBlockHeight);
    }
    function emitKeeperRegistryBase2_0_UpkeepCheckDataUpdated(uint256 id,bytes memory newCheckData) public {
        emit KeeperRegistryBase2_0_UpkeepCheckDataUpdated(id,newCheckData);
    }
    function emitKeeperRegistryBase2_0_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit KeeperRegistryBase2_0_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitKeeperRegistryBase2_0_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit KeeperRegistryBase2_0_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitKeeperRegistryBase2_0_UpkeepOffchainConfigSet(uint256 id,bytes memory offchainConfig) public {
        emit KeeperRegistryBase2_0_UpkeepOffchainConfigSet(id,offchainConfig);
    }
    function emitKeeperRegistryBase2_0_UpkeepPaused(uint256 id) public {
        emit KeeperRegistryBase2_0_UpkeepPaused(id);
    }
    function emitKeeperRegistryBase2_0_UpkeepPerformed(uint256 id,bool success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment) public {
        emit KeeperRegistryBase2_0_UpkeepPerformed(id,success,checkBlockNumber,gasUsed,gasOverhead,totalPayment);
    }
    function emitKeeperRegistryBase2_0_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit KeeperRegistryBase2_0_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitKeeperRegistryBase2_0_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit KeeperRegistryBase2_0_UpkeepRegistered(id,executeGas,admin);
    }
    function emitKeeperRegistryBase2_0_UpkeepUnpaused(uint256 id) public {
        emit KeeperRegistryBase2_0_UpkeepUnpaused(id);
    }
    function emitKeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferRequested(address from,address to) public {
        emit KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferRequested(from,to);
    }
    function emitKeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferred(address from,address to) public {
        emit KeeperRegistryCheckUpkeepGasUsageWrapper1_2_OwnershipTransferred(from,to);
    }
    function emitKeeperRegistryLogic1_3_ConfigSet(KeeperRegistryLogic1_3_Config memory config) public {
        emit KeeperRegistryLogic1_3_ConfigSet(config);
    }
    function emitKeeperRegistryLogic1_3_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit KeeperRegistryLogic1_3_FundsAdded(id,from,amount);
    }
    function emitKeeperRegistryLogic1_3_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit KeeperRegistryLogic1_3_FundsWithdrawn(id,amount,to);
    }
    function emitKeeperRegistryLogic1_3_KeepersUpdated(address[] memory keepers,address[] memory payees) public {
        emit KeeperRegistryLogic1_3_KeepersUpdated(keepers,payees);
    }
    function emitKeeperRegistryLogic1_3_OwnerFundsWithdrawn(uint96 amount) public {
        emit KeeperRegistryLogic1_3_OwnerFundsWithdrawn(amount);
    }
    function emitKeeperRegistryLogic1_3_OwnershipTransferRequested(address from,address to) public {
        emit KeeperRegistryLogic1_3_OwnershipTransferRequested(from,to);
    }
    function emitKeeperRegistryLogic1_3_OwnershipTransferred(address from,address to) public {
        emit KeeperRegistryLogic1_3_OwnershipTransferred(from,to);
    }
    function emitKeeperRegistryLogic1_3_Paused(address account) public {
        emit KeeperRegistryLogic1_3_Paused(account);
    }
    function emitKeeperRegistryLogic1_3_PayeeshipTransferRequested(address keeper,address from,address to) public {
        emit KeeperRegistryLogic1_3_PayeeshipTransferRequested(keeper,from,to);
    }
    function emitKeeperRegistryLogic1_3_PayeeshipTransferred(address keeper,address from,address to) public {
        emit KeeperRegistryLogic1_3_PayeeshipTransferred(keeper,from,to);
    }
    function emitKeeperRegistryLogic1_3_PaymentWithdrawn(address keeper,uint256 amount,address to,address payee) public {
        emit KeeperRegistryLogic1_3_PaymentWithdrawn(keeper,amount,to,payee);
    }
    function emitKeeperRegistryLogic1_3_Unpaused(address account) public {
        emit KeeperRegistryLogic1_3_Unpaused(account);
    }
    function emitKeeperRegistryLogic1_3_UpkeepAdminTransferRequested(uint256 id,address from,address to) public {
        emit KeeperRegistryLogic1_3_UpkeepAdminTransferRequested(id,from,to);
    }
    function emitKeeperRegistryLogic1_3_UpkeepAdminTransferred(uint256 id,address from,address to) public {
        emit KeeperRegistryLogic1_3_UpkeepAdminTransferred(id,from,to);
    }
    function emitKeeperRegistryLogic1_3_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit KeeperRegistryLogic1_3_UpkeepCanceled(id,atBlockHeight);
    }
    function emitKeeperRegistryLogic1_3_UpkeepCheckDataUpdated(uint256 id,bytes memory newCheckData) public {
        emit KeeperRegistryLogic1_3_UpkeepCheckDataUpdated(id,newCheckData);
    }
    function emitKeeperRegistryLogic1_3_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit KeeperRegistryLogic1_3_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitKeeperRegistryLogic1_3_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit KeeperRegistryLogic1_3_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitKeeperRegistryLogic1_3_UpkeepPaused(uint256 id) public {
        emit KeeperRegistryLogic1_3_UpkeepPaused(id);
    }
    function emitKeeperRegistryLogic1_3_UpkeepPerformed(uint256 id,bool success,address from,uint96 payment,bytes memory performData) public {
        emit KeeperRegistryLogic1_3_UpkeepPerformed(id,success,from,payment,performData);
    }
    function emitKeeperRegistryLogic1_3_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit KeeperRegistryLogic1_3_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitKeeperRegistryLogic1_3_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit KeeperRegistryLogic1_3_UpkeepRegistered(id,executeGas,admin);
    }
    function emitKeeperRegistryLogic1_3_UpkeepUnpaused(uint256 id) public {
        emit KeeperRegistryLogic1_3_UpkeepUnpaused(id);
    }
    function emitKeeperRegistryLogic2_0_CancelledUpkeepReport(uint256 id) public {
        emit KeeperRegistryLogic2_0_CancelledUpkeepReport(id);
    }
    function emitKeeperRegistryLogic2_0_FundsAdded(uint256 id,address from,uint96 amount) public {
        emit KeeperRegistryLogic2_0_FundsAdded(id,from,amount);
    }
    function emitKeeperRegistryLogic2_0_FundsWithdrawn(uint256 id,uint256 amount,address to) public {
        emit KeeperRegistryLogic2_0_FundsWithdrawn(id,amount,to);
    }
    function emitKeeperRegistryLogic2_0_InsufficientFundsUpkeepReport(uint256 id) public {
        emit KeeperRegistryLogic2_0_InsufficientFundsUpkeepReport(id);
    }
    function emitKeeperRegistryLogic2_0_OwnerFundsWithdrawn(uint96 amount) public {
        emit KeeperRegistryLogic2_0_OwnerFundsWithdrawn(amount);
    }
    function emitKeeperRegistryLogic2_0_OwnershipTransferRequested(address from,address to) public {
        emit KeeperRegistryLogic2_0_OwnershipTransferRequested(from,to);
    }
    function emitKeeperRegistryLogic2_0_OwnershipTransferred(address from,address to) public {
        emit KeeperRegistryLogic2_0_OwnershipTransferred(from,to);
    }
    function emitKeeperRegistryLogic2_0_Paused(address account) public {
        emit KeeperRegistryLogic2_0_Paused(account);
    }
    function emitKeeperRegistryLogic2_0_PayeesUpdated(address[] memory transmitters,address[] memory payees) public {
        emit KeeperRegistryLogic2_0_PayeesUpdated(transmitters,payees);
    }
    function emitKeeperRegistryLogic2_0_PayeeshipTransferRequested(address transmitter,address from,address to) public {
        emit KeeperRegistryLogic2_0_PayeeshipTransferRequested(transmitter,from,to);
    }
    function emitKeeperRegistryLogic2_0_PayeeshipTransferred(address transmitter,address from,address to) public {
        emit KeeperRegistryLogic2_0_PayeeshipTransferred(transmitter,from,to);
    }
    function emitKeeperRegistryLogic2_0_PaymentWithdrawn(address transmitter,uint256 amount,address to,address payee) public {
        emit KeeperRegistryLogic2_0_PaymentWithdrawn(transmitter,amount,to,payee);
    }
    function emitKeeperRegistryLogic2_0_ReorgedUpkeepReport(uint256 id) public {
        emit KeeperRegistryLogic2_0_ReorgedUpkeepReport(id);
    }
    function emitKeeperRegistryLogic2_0_StaleUpkeepReport(uint256 id) public {
        emit KeeperRegistryLogic2_0_StaleUpkeepReport(id);
    }
    function emitKeeperRegistryLogic2_0_Unpaused(address account) public {
        emit KeeperRegistryLogic2_0_Unpaused(account);
    }
    function emitKeeperRegistryLogic2_0_UpkeepAdminTransferRequested(uint256 id,address from,address to) public {
        emit KeeperRegistryLogic2_0_UpkeepAdminTransferRequested(id,from,to);
    }
    function emitKeeperRegistryLogic2_0_UpkeepAdminTransferred(uint256 id,address from,address to) public {
        emit KeeperRegistryLogic2_0_UpkeepAdminTransferred(id,from,to);
    }
    function emitKeeperRegistryLogic2_0_UpkeepCanceled(uint256 id,uint64 atBlockHeight) public {
        emit KeeperRegistryLogic2_0_UpkeepCanceled(id,atBlockHeight);
    }
    function emitKeeperRegistryLogic2_0_UpkeepCheckDataUpdated(uint256 id,bytes memory newCheckData) public {
        emit KeeperRegistryLogic2_0_UpkeepCheckDataUpdated(id,newCheckData);
    }
    function emitKeeperRegistryLogic2_0_UpkeepGasLimitSet(uint256 id,uint96 gasLimit) public {
        emit KeeperRegistryLogic2_0_UpkeepGasLimitSet(id,gasLimit);
    }
    function emitKeeperRegistryLogic2_0_UpkeepMigrated(uint256 id,uint256 remainingBalance,address destination) public {
        emit KeeperRegistryLogic2_0_UpkeepMigrated(id,remainingBalance,destination);
    }
    function emitKeeperRegistryLogic2_0_UpkeepOffchainConfigSet(uint256 id,bytes memory offchainConfig) public {
        emit KeeperRegistryLogic2_0_UpkeepOffchainConfigSet(id,offchainConfig);
    }
    function emitKeeperRegistryLogic2_0_UpkeepPaused(uint256 id) public {
        emit KeeperRegistryLogic2_0_UpkeepPaused(id);
    }
    function emitKeeperRegistryLogic2_0_UpkeepPerformed(uint256 id,bool success,uint32 checkBlockNumber,uint256 gasUsed,uint256 gasOverhead,uint96 totalPayment) public {
        emit KeeperRegistryLogic2_0_UpkeepPerformed(id,success,checkBlockNumber,gasUsed,gasOverhead,totalPayment);
    }
    function emitKeeperRegistryLogic2_0_UpkeepReceived(uint256 id,uint256 startingBalance,address importedFrom) public {
        emit KeeperRegistryLogic2_0_UpkeepReceived(id,startingBalance,importedFrom);
    }
    function emitKeeperRegistryLogic2_0_UpkeepRegistered(uint256 id,uint32 executeGas,address admin) public {
        emit KeeperRegistryLogic2_0_UpkeepRegistered(id,executeGas,admin);
    }
    function emitKeeperRegistryLogic2_0_UpkeepUnpaused(uint256 id) public {
        emit KeeperRegistryLogic2_0_UpkeepUnpaused(id);
    }
    function emitLogEmitter_Log1(uint256 jxztvtdu) public {
        emit LogEmitter_Log1(jxztvtdu);
    }
    function emitLogEmitter_Log2(uint256 jbysnosu) public {
        emit LogEmitter_Log2(jbysnosu);
    }
    function emitLogEmitter_Log3(string memory njjusihh) public {
        emit LogEmitter_Log3(njjusihh);
    }
    function emitOCR2Abstract_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
        emit OCR2Abstract_ConfigSet(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
    }
    function emitOCR2Abstract_Transmitted(bytes32 configDigest,uint32 epoch) public {
        emit OCR2Abstract_Transmitted(configDigest,epoch);
    }
    function emitOCR2BaseUpgradeable_ConfigSet(uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,address[] memory transmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
        emit OCR2BaseUpgradeable_ConfigSet(previousConfigBlockNumber,configDigest,configCount,signers,transmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
    }
    function emitOCR2BaseUpgradeable_Initialized(uint8 version) public {
        emit OCR2BaseUpgradeable_Initialized(version);
    }
    function emitOCR2BaseUpgradeable_OwnershipTransferRequested(address from,address to) public {
        emit OCR2BaseUpgradeable_OwnershipTransferRequested(from,to);
    }
    function emitOCR2BaseUpgradeable_OwnershipTransferred(address from,address to) public {
        emit OCR2BaseUpgradeable_OwnershipTransferred(from,to);
    }
    function emitOCR2BaseUpgradeable_Transmitted(bytes32 configDigest,uint32 epoch) public {
        emit OCR2BaseUpgradeable_Transmitted(configDigest,epoch);
    }
    function emitOVM_GasPriceOracle_DecimalsUpdated(uint256 ofmpqceu) public {
        emit OVM_GasPriceOracle_DecimalsUpdated(ofmpqceu);
    }
    function emitOVM_GasPriceOracle_GasPriceUpdated(uint256 lflivotg) public {
        emit OVM_GasPriceOracle_GasPriceUpdated(lflivotg);
    }
    function emitOVM_GasPriceOracle_L1BaseFeeUpdated(uint256 ohwqxozp) public {
        emit OVM_GasPriceOracle_L1BaseFeeUpdated(ohwqxozp);
    }
    function emitOVM_GasPriceOracle_OverheadUpdated(uint256 jvbubple) public {
        emit OVM_GasPriceOracle_OverheadUpdated(jvbubple);
    }
    function emitOVM_GasPriceOracle_OwnershipTransferred(address previousOwner,address newOwner) public {
        emit OVM_GasPriceOracle_OwnershipTransferred(previousOwner,newOwner);
    }
    function emitOVM_GasPriceOracle_ScalarUpdated(uint256 hdoeyfyg) public {
        emit OVM_GasPriceOracle_ScalarUpdated(hdoeyfyg);
    }
    function emitOwnable_OwnershipTransferred(address previousOwner,address newOwner) public {
        emit Ownable_OwnershipTransferred(previousOwner,newOwner);
    }
    function emitPausable_Paused(address account) public {
        emit Pausable_Paused(account);
    }
    function emitPausable_Unpaused(address account) public {
        emit Pausable_Unpaused(account);
    }
    function emitPausableUpgradeable_Initialized(uint8 version) public {
        emit PausableUpgradeable_Initialized(version);
    }
    function emitPausableUpgradeable_Paused(address account) public {
        emit PausableUpgradeable_Paused(account);
    }
    function emitPausableUpgradeable_Unpaused(address account) public {
        emit PausableUpgradeable_Unpaused(account);
    }
    function emitProxyAdmin_OwnershipTransferred(address previousOwner,address newOwner) public {
        emit ProxyAdmin_OwnershipTransferred(previousOwner,newOwner);
    }
    function emitTransparentUpgradeableProxy_AdminChanged(address previousAdmin,address newAdmin) public {
        emit TransparentUpgradeableProxy_AdminChanged(previousAdmin,newAdmin);
    }
    function emitTransparentUpgradeableProxy_BeaconUpgraded(address beacon) public {
        emit TransparentUpgradeableProxy_BeaconUpgraded(beacon);
    }
    function emitTransparentUpgradeableProxy_Upgraded(address implementation) public {
        emit TransparentUpgradeableProxy_Upgraded(implementation);
    }
    function emitVRFConsumerBaseV2Upgradeable_Initialized(uint8 version) public {
        emit VRFConsumerBaseV2Upgradeable_Initialized(version);
    }
    function emitVRFConsumerV2UpgradeableExample_Initialized(uint8 version) public {
        emit VRFConsumerV2UpgradeableExample_Initialized(version);
    }
    function emitVRFCoordinatorMock_RandomnessRequest(address sender,bytes32 keyHash,uint256 seed) public {
        emit VRFCoordinatorMock_RandomnessRequest(sender,keyHash,seed);
    }
    function emitVRFCoordinatorV2_ConfigSet(uint16 minimumRequestConfirmations,uint32 maxGasLimit,uint32 stalenessSeconds,uint32 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,VRFCoordinatorV2_FeeConfig memory feeConfig) public {
        emit VRFCoordinatorV2_ConfigSet(minimumRequestConfirmations,maxGasLimit,stalenessSeconds,gasAfterPaymentCalculation,fallbackWeiPerUnitLink,feeConfig);
    }
    function emitVRFCoordinatorV2_FundsRecovered(address to,uint256 amount) public {
        emit VRFCoordinatorV2_FundsRecovered(to,amount);
    }
    function emitVRFCoordinatorV2_OwnershipTransferRequested(address from,address to) public {
        emit VRFCoordinatorV2_OwnershipTransferRequested(from,to);
    }
    function emitVRFCoordinatorV2_OwnershipTransferred(address from,address to) public {
        emit VRFCoordinatorV2_OwnershipTransferred(from,to);
    }
    function emitVRFCoordinatorV2_ProvingKeyDeregistered(bytes32 keyHash,address oracle) public {
        emit VRFCoordinatorV2_ProvingKeyDeregistered(keyHash,oracle);
    }
    function emitVRFCoordinatorV2_ProvingKeyRegistered(bytes32 keyHash,address oracle) public {
        emit VRFCoordinatorV2_ProvingKeyRegistered(keyHash,oracle);
    }
    function emitVRFCoordinatorV2_RandomWordsFulfilled(uint256 requestId,uint256 outputSeed,uint96 payment,bool success) public {
        emit VRFCoordinatorV2_RandomWordsFulfilled(requestId,outputSeed,payment,success);
    }
    function emitVRFCoordinatorV2_RandomWordsRequested(bytes32 keyHash,uint256 requestId,uint256 preSeed,uint64 subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address sender) public {
        emit VRFCoordinatorV2_RandomWordsRequested(keyHash,requestId,preSeed,subId,minimumRequestConfirmations,callbackGasLimit,numWords,sender);
    }
    function emitVRFCoordinatorV2_SubscriptionCanceled(uint64 subId,address to,uint256 amount) public {
        emit VRFCoordinatorV2_SubscriptionCanceled(subId,to,amount);
    }
    function emitVRFCoordinatorV2_SubscriptionConsumerAdded(uint64 subId,address consumer) public {
        emit VRFCoordinatorV2_SubscriptionConsumerAdded(subId,consumer);
    }
    function emitVRFCoordinatorV2_SubscriptionConsumerRemoved(uint64 subId,address consumer) public {
        emit VRFCoordinatorV2_SubscriptionConsumerRemoved(subId,consumer);
    }
    function emitVRFCoordinatorV2_SubscriptionCreated(uint64 subId,address owner) public {
        emit VRFCoordinatorV2_SubscriptionCreated(subId,owner);
    }
    function emitVRFCoordinatorV2_SubscriptionFunded(uint64 subId,uint256 oldBalance,uint256 newBalance) public {
        emit VRFCoordinatorV2_SubscriptionFunded(subId,oldBalance,newBalance);
    }
    function emitVRFCoordinatorV2_SubscriptionOwnerTransferRequested(uint64 subId,address from,address to) public {
        emit VRFCoordinatorV2_SubscriptionOwnerTransferRequested(subId,from,to);
    }
    function emitVRFCoordinatorV2_SubscriptionOwnerTransferred(uint64 subId,address from,address to) public {
        emit VRFCoordinatorV2_SubscriptionOwnerTransferred(subId,from,to);
    }
    function emitVRFCoordinatorV2TestHelper_ConfigSet(uint16 minimumRequestConfirmations,uint32 maxGasLimit,uint32 stalenessSeconds,uint32 gasAfterPaymentCalculation,int256 fallbackWeiPerUnitLink,VRFCoordinatorV2TestHelper_FeeConfig memory feeConfig) public {
        emit VRFCoordinatorV2TestHelper_ConfigSet(minimumRequestConfirmations,maxGasLimit,stalenessSeconds,gasAfterPaymentCalculation,fallbackWeiPerUnitLink,feeConfig);
    }
    function emitVRFCoordinatorV2TestHelper_FundsRecovered(address to,uint256 amount) public {
        emit VRFCoordinatorV2TestHelper_FundsRecovered(to,amount);
    }
    function emitVRFCoordinatorV2TestHelper_OwnershipTransferRequested(address from,address to) public {
        emit VRFCoordinatorV2TestHelper_OwnershipTransferRequested(from,to);
    }
    function emitVRFCoordinatorV2TestHelper_OwnershipTransferred(address from,address to) public {
        emit VRFCoordinatorV2TestHelper_OwnershipTransferred(from,to);
    }
    function emitVRFCoordinatorV2TestHelper_ProvingKeyDeregistered(bytes32 keyHash,address oracle) public {
        emit VRFCoordinatorV2TestHelper_ProvingKeyDeregistered(keyHash,oracle);
    }
    function emitVRFCoordinatorV2TestHelper_ProvingKeyRegistered(bytes32 keyHash,address oracle) public {
        emit VRFCoordinatorV2TestHelper_ProvingKeyRegistered(keyHash,oracle);
    }
    function emitVRFCoordinatorV2TestHelper_RandomWordsFulfilled(uint256 requestId,uint256 outputSeed,uint96 payment,bool success) public {
        emit VRFCoordinatorV2TestHelper_RandomWordsFulfilled(requestId,outputSeed,payment,success);
    }
    function emitVRFCoordinatorV2TestHelper_RandomWordsRequested(bytes32 keyHash,uint256 requestId,uint256 preSeed,uint64 subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address sender) public {
        emit VRFCoordinatorV2TestHelper_RandomWordsRequested(keyHash,requestId,preSeed,subId,minimumRequestConfirmations,callbackGasLimit,numWords,sender);
    }
    function emitVRFCoordinatorV2TestHelper_SubscriptionCanceled(uint64 subId,address to,uint256 amount) public {
        emit VRFCoordinatorV2TestHelper_SubscriptionCanceled(subId,to,amount);
    }
    function emitVRFCoordinatorV2TestHelper_SubscriptionConsumerAdded(uint64 subId,address consumer) public {
        emit VRFCoordinatorV2TestHelper_SubscriptionConsumerAdded(subId,consumer);
    }
    function emitVRFCoordinatorV2TestHelper_SubscriptionConsumerRemoved(uint64 subId,address consumer) public {
        emit VRFCoordinatorV2TestHelper_SubscriptionConsumerRemoved(subId,consumer);
    }
    function emitVRFCoordinatorV2TestHelper_SubscriptionCreated(uint64 subId,address owner) public {
        emit VRFCoordinatorV2TestHelper_SubscriptionCreated(subId,owner);
    }
    function emitVRFCoordinatorV2TestHelper_SubscriptionFunded(uint64 subId,uint256 oldBalance,uint256 newBalance) public {
        emit VRFCoordinatorV2TestHelper_SubscriptionFunded(subId,oldBalance,newBalance);
    }
    function emitVRFCoordinatorV2TestHelper_SubscriptionOwnerTransferRequested(uint64 subId,address from,address to) public {
        emit VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferRequested(subId,from,to);
    }
    function emitVRFCoordinatorV2TestHelper_SubscriptionOwnerTransferred(uint64 subId,address from,address to) public {
        emit VRFCoordinatorV2TestHelper_SubscriptionOwnerTransferred(subId,from,to);
    }
    function emitVRFLoadTestExternalSubOwner_OwnershipTransferRequested(address from,address to) public {
        emit VRFLoadTestExternalSubOwner_OwnershipTransferRequested(from,to);
    }
    function emitVRFLoadTestExternalSubOwner_OwnershipTransferred(address from,address to) public {
        emit VRFLoadTestExternalSubOwner_OwnershipTransferred(from,to);
    }
    function emitVRFV2ProxyAdmin_OwnershipTransferred(address previousOwner,address newOwner) public {
        emit VRFV2ProxyAdmin_OwnershipTransferred(previousOwner,newOwner);
    }
    function emitVRFV2TransparentUpgradeableProxy_AdminChanged(address previousAdmin,address newAdmin) public {
        emit VRFV2TransparentUpgradeableProxy_AdminChanged(previousAdmin,newAdmin);
    }
    function emitVRFV2TransparentUpgradeableProxy_BeaconUpgraded(address beacon) public {
        emit VRFV2TransparentUpgradeableProxy_BeaconUpgraded(beacon);
    }
    function emitVRFV2TransparentUpgradeableProxy_Upgraded(address implementation) public {
        emit VRFV2TransparentUpgradeableProxy_Upgraded(implementation);
    }
    function emitVRFV2Wrapper_OwnershipTransferRequested(address from,address to) public {
        emit VRFV2Wrapper_OwnershipTransferRequested(from,to);
    }
    function emitVRFV2Wrapper_OwnershipTransferred(address from,address to) public {
        emit VRFV2Wrapper_OwnershipTransferred(from,to);
    }
    function emitVRFV2Wrapper_WrapperFulfillmentFailed(uint256 requestId,address consumer) public {
        emit VRFV2Wrapper_WrapperFulfillmentFailed(requestId,consumer);
    }
    function emitVRFV2WrapperConsumerExample_OwnershipTransferRequested(address from,address to) public {
        emit VRFV2WrapperConsumerExample_OwnershipTransferRequested(from,to);
    }
    function emitVRFV2WrapperConsumerExample_OwnershipTransferred(address from,address to) public {
        emit VRFV2WrapperConsumerExample_OwnershipTransferred(from,to);
    }
    function emitVRFV2WrapperConsumerExample_WrappedRequestFulfilled(uint256 requestId,uint256[] memory randomWords,uint256 payment) public {
        emit VRFV2WrapperConsumerExample_WrappedRequestFulfilled(requestId,randomWords,payment);
    }
    function emitVRFV2WrapperConsumerExample_WrapperRequestMade(uint256 requestId,uint256 paid) public {
        emit VRFV2WrapperConsumerExample_WrapperRequestMade(requestId,paid);
    }
    function emitVerifier_ConfigActivated(bytes32 feedId,bytes32 configDigest) public {
        emit Verifier_ConfigActivated(feedId,configDigest);
    }
    function emitVerifier_ConfigDeactivated(bytes32 feedId,bytes32 configDigest) public {
        emit Verifier_ConfigDeactivated(feedId,configDigest);
    }
    function emitVerifier_ConfigSet(bytes32 feedId,uint32 previousConfigBlockNumber,bytes32 configDigest,uint64 configCount,address[] memory signers,bytes32[] memory offchainTransmitters,uint8 f,bytes memory onchainConfig,uint64 offchainConfigVersion,bytes memory offchainConfig) public {
        emit Verifier_ConfigSet(feedId,previousConfigBlockNumber,configDigest,configCount,signers,offchainTransmitters,f,onchainConfig,offchainConfigVersion,offchainConfig);
    }
    function emitVerifier_OwnershipTransferRequested(address from,address to) public {
        emit Verifier_OwnershipTransferRequested(from,to);
    }
    function emitVerifier_OwnershipTransferred(address from,address to) public {
        emit Verifier_OwnershipTransferred(from,to);
    }
    function emitVerifier_ReportVerified(bytes32 feedId,bytes32 reportHash,address requester) public {
        emit Verifier_ReportVerified(feedId,reportHash,requester);
    }
    function emitVerifierProxy_AccessControllerSet(address oldAccessController,address newAccessController) public {
        emit VerifierProxy_AccessControllerSet(oldAccessController,newAccessController);
    }
    function emitVerifierProxy_OwnershipTransferRequested(address from,address to) public {
        emit VerifierProxy_OwnershipTransferRequested(from,to);
    }
    function emitVerifierProxy_OwnershipTransferred(address from,address to) public {
        emit VerifierProxy_OwnershipTransferred(from,to);
    }
    function emitVerifierProxy_VerifierSet(bytes32 oldConfigDigest,bytes32 newConfigDigest,address verifierAddress) public {
        emit VerifierProxy_VerifierSet(oldConfigDigest,newConfigDigest,verifierAddress);
    }
    function emitVerifierProxy_VerifierUnset(bytes32 configDigest,address verifierAddress) public {
        emit VerifierProxy_VerifierUnset(configDigest,verifierAddress);
    }
}
