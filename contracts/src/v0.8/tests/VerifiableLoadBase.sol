// SPDX-License-Identifier: MIT
pragma solidity ^0.8.16;

import "../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import {IKeeperRegistryMaster, IAutomationV21PlusCommon} from "../automation/interfaces/v2_1/IKeeperRegistryMaster.sol";
import {ArbSys} from "../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";
import "../automation/v2_1/AutomationRegistrar2_1.sol";
import {LogTriggerConfig} from "../automation/v2_1/AutomationUtils2_1.sol";

abstract contract VerifiableLoadBase is ConfirmedOwner {
  error IndexOutOfRange();
  event LogEmitted(uint256 indexed upkeepId, uint256 indexed blockNum, address indexed addr);
  event LogEmittedAgain(uint256 indexed upkeepId, uint256 indexed blockNum, address indexed addr);
  event UpkeepTopUp(uint256 upkeepId, uint96 amount, uint256 blockNum);

  using EnumerableSet for EnumerableSet.UintSet;
  ArbSys internal constant ARB_SYS = ArbSys(0x0000000000000000000000000000000000000064);
  //bytes32 public constant emittedSig = 0x97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf08; //keccak256(LogEmitted(uint256,uint256,address))
  bytes32 public immutable emittedSig = LogEmitted.selector;
  // bytes32 public constant emittedAgainSig = 0xc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d; //keccak256(LogEmittedAgain(uint256,uint256,address))
  bytes32 public immutable emittedAgainSig = LogEmittedAgain.selector;

  mapping(uint256 => uint256) public lastTopUpBlocks;
  mapping(uint256 => uint256) public intervals;
  mapping(uint256 => uint256) public previousPerformBlocks;
  mapping(uint256 => uint256) public firstPerformBlocks;
  mapping(uint256 => uint256) public counters;
  mapping(uint256 => uint256) public performGasToBurns;
  mapping(uint256 => uint256) public checkGasToBurns;
  mapping(uint256 => uint256) public performDataSizes;
  mapping(uint256 => uint256) public gasLimits;
  mapping(bytes32 => bool) public dummyMap; // used to force storage lookup
  mapping(uint256 => uint256[]) public delays; // how to query for delays for a certain past period: calendar day and/or past 24 hours

  mapping(uint256 => mapping(uint16 => uint256[])) public bucketedDelays;
  mapping(uint256 => uint16) public buckets;
  EnumerableSet.UintSet internal s_upkeepIDs;
  AutomationRegistrar2_1 public registrar;
  LinkTokenInterface public linkToken;
  IKeeperRegistryMaster public registry;
  // check if an upkeep is eligible for adding funds at this interval
  uint256 public upkeepTopUpCheckInterval = 5;
  // an upkeep will get this amount of LINK for every top up
  uint96 public addLinkAmount = 200000000000000000; // 0.2 LINK
  // if an upkeep's balance is less than this threshold * min balance, this upkeep is eligible for adding funds
  uint8 public minBalanceThresholdMultiplier = 20;
  // if this contract is using arbitrum block number
  bool public immutable useArbitrumBlockNum;

  // the following fields are immutable bc if they are adjusted, the existing upkeeps' delays will be stored in
  // different sizes of buckets. it's better to redeploy this contract with new values.
  uint16 public immutable BUCKET_SIZE = 100;

  // search feeds in notion: "Schema and Feed ID Registry"
  string[] public feedsHex = [
    "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
    "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"
  ];
  string public feedParamKey = "feedIdHex"; // feedIDs for v0.3
  string public timeParamKey = "blockNumber"; // timestamp for v0.3

  /**
   * @param _registrar a automation registrar 2.1 address
   * @param _useArb if this contract will use arbitrum block number
   */
  constructor(AutomationRegistrar2_1 _registrar, bool _useArb) ConfirmedOwner(msg.sender) {
    registrar = _registrar;
    (address registryAddress, ) = registrar.getConfig();
    registry = IKeeperRegistryMaster(payable(address(registryAddress)));
    linkToken = registrar.LINK();
    useArbitrumBlockNum = _useArb;
  }

  receive() external payable {}

  function setParamKeys(string memory _feedParamKey, string memory _timeParamKey) external {
    feedParamKey = _feedParamKey;
    timeParamKey = _timeParamKey;
  }

  function setFeeds(string[] memory _feeds) external {
    feedsHex = _feeds;
  }

  /**
   * @notice withdraws LINKs from this contract to msg sender when testing is finished.
   */
  function withdrawLinks() external onlyOwner {
    uint256 balance = linkToken.balanceOf(address(this));
    linkToken.transfer(msg.sender, balance);
  }

  function getBlockNumber() internal view returns (uint256) {
    if (useArbitrumBlockNum) {
      return ARB_SYS.arbBlockNumber();
    } else {
      return block.number;
    }
  }

  /**
   * @notice sets registrar, registry, and link token address.
   * @param newRegistrar the new registrar address
   */
  function setConfig(AutomationRegistrar2_1 newRegistrar) external {
    registrar = newRegistrar;
    (address registryAddress, ) = registrar.getConfig();
    registry = IKeeperRegistryMaster(payable(address(registryAddress)));
    linkToken = registrar.LINK();
  }

  /**
   * @notice gets an array of active upkeep IDs.
   * @param startIndex the start index of upkeep IDs
   * @param maxCount the max number of upkeep IDs requested
   * @return an array of active upkeep IDs
   */
  function getActiveUpkeepIDsDeployedByThisContract(
    uint256 startIndex,
    uint256 maxCount
  ) external view returns (uint256[] memory) {
    uint256 maxIdx = s_upkeepIDs.length();
    if (startIndex >= maxIdx) revert IndexOutOfRange();
    if (maxCount == 0) {
      maxCount = maxIdx - startIndex;
    }
    uint256[] memory ids = new uint256[](maxCount);
    for (uint256 idx = 0; idx < maxCount; idx++) {
      ids[idx] = s_upkeepIDs.at(startIndex + idx);
    }
    return ids;
  }

  /**
   * @notice gets an array of active upkeep IDs.
   * @param startIndex the start index of upkeep IDs
   * @param maxCount the max number of upkeep IDs requested
   * @return an array of active upkeep IDs
   */
  function getAllActiveUpkeepIDsOnRegistry(
    uint256 startIndex,
    uint256 maxCount
  ) external view returns (uint256[] memory) {
    return registry.getActiveUpkeepIDs(startIndex, maxCount);
  }

  /**
   * @notice register an upkeep via the registrar.
   * @param params a registration params struct
   * @return an upkeep ID
   */
  function _registerUpkeep(AutomationRegistrar2_1.RegistrationParams memory params) private returns (uint256) {
    uint256 upkeepId = registrar.registerUpkeep(params);
    s_upkeepIDs.add(upkeepId);
    gasLimits[upkeepId] = params.gasLimit;
    return upkeepId;
  }

  /**
   * @notice returns a log trigger config
   */
  function getLogTriggerConfig(
    address addr,
    uint8 selector,
    bytes32 topic0,
    bytes32 topic1,
    bytes32 topic2,
    bytes32 topic3
  ) external view returns (bytes memory logTrigger) {
    LogTriggerConfig memory cfg = LogTriggerConfig({
      contractAddress: addr,
      filterSelector: selector,
      topic0: topic0,
      topic1: topic1,
      topic2: topic2,
      topic3: topic3
    });
    return abi.encode(cfg);
  }

  // this function sets pipeline data and trigger config for log trigger upkeeps
  function batchPreparingUpkeeps(
    uint256[] calldata upkeepIds,
    uint8 selector,
    bytes32 topic0,
    bytes32 topic1,
    bytes32 topic2,
    bytes32 topic3
  ) external {
    uint256 len = upkeepIds.length;
    for (uint256 i = 0; i < len; i++) {
      uint256 upkeepId = upkeepIds[i];

      this.updateUpkeepPipelineData(upkeepId, abi.encode(upkeepId));

      uint8 triggerType = registry.getTriggerType(upkeepId);
      if (triggerType == 1) {
        // currently no using a filter selector
        bytes memory triggerCfg = this.getLogTriggerConfig(address(this), selector, topic0, topic2, topic2, topic3);
        registry.setUpkeepTriggerConfig(upkeepId, triggerCfg);
      }
    }
  }

  // this function sets pipeline data and trigger config for log trigger upkeeps
  function batchPreparingUpkeepsSimple(uint256[] calldata upkeepIds, uint8 log, uint8 selector) external {
    uint256 len = upkeepIds.length;
    for (uint256 i = 0; i < len; i++) {
      uint256 upkeepId = upkeepIds[i];

      this.updateUpkeepPipelineData(upkeepId, abi.encode(upkeepId));

      uint8 triggerType = registry.getTriggerType(upkeepId);
      if (triggerType == 1) {
        // currently no using a filter selector
        bytes32 sig = emittedSig;
        if (log != 0) {
          sig = emittedAgainSig;
        }
        bytes memory triggerCfg = this.getLogTriggerConfig(
          address(this),
          selector,
          sig,
          bytes32(abi.encode(upkeepId)),
          bytes32(0),
          bytes32(0)
        );
        registry.setUpkeepTriggerConfig(upkeepId, triggerCfg);
      }
    }
  }

  function updateLogTriggerConfig1(
    uint256 upkeepId,
    address addr,
    uint8 selector,
    bytes32 topic0,
    bytes32 topic1,
    bytes32 topic2,
    bytes32 topic3
  ) external {
    registry.setUpkeepTriggerConfig(upkeepId, this.getLogTriggerConfig(addr, selector, topic0, topic1, topic2, topic3));
  }

  function updateLogTriggerConfig2(uint256 upkeepId, bytes calldata cfg) external {
    registry.setUpkeepTriggerConfig(upkeepId, cfg);
  }

  /**
   * @notice batch registering upkeeps.
   * @param number the number of upkeeps to be registered
   * @param gasLimit the gas limit of each upkeep
   * @param triggerType the trigger type of this upkeep, 0 for conditional, 1 for log trigger
   * @param triggerConfig the trigger config of this upkeep
   * @param amount the amount of LINK to fund each upkeep
   * @param checkGasToBurn the amount of check gas to burn
   * @param performGasToBurn the amount of perform gas to burn
   */
  function batchRegisterUpkeeps(
    uint8 number,
    uint32 gasLimit,
    uint8 triggerType,
    bytes memory triggerConfig,
    uint96 amount,
    uint256 checkGasToBurn,
    uint256 performGasToBurn
  ) external {
    AutomationRegistrar2_1.RegistrationParams memory params = AutomationRegistrar2_1.RegistrationParams({
      name: "test",
      encryptedEmail: bytes(""),
      upkeepContract: address(this),
      gasLimit: gasLimit,
      adminAddress: address(this), // use address of this contract as the admin
      triggerType: triggerType,
      checkData: bytes(""), // update pipeline data later bc upkeep id is not available now
      triggerConfig: triggerConfig,
      offchainConfig: bytes(""),
      amount: amount
    });

    linkToken.approve(address(registrar), amount * number);

    uint256[] memory upkeepIds = new uint256[](number);
    for (uint8 i = 0; i < number; i++) {
      uint256 upkeepId = _registerUpkeep(params);
      upkeepIds[i] = upkeepId;
      checkGasToBurns[upkeepId] = checkGasToBurn;
      performGasToBurns[upkeepId] = performGasToBurn;
    }
  }

  function topUpFund(uint256 upkeepId, uint256 blockNum) public {
    if (blockNum - lastTopUpBlocks[upkeepId] > upkeepTopUpCheckInterval) {
      IAutomationV21PlusCommon.UpkeepInfoLegacy memory info = registry.getUpkeep(upkeepId);
      uint96 minBalance = registry.getMinBalanceForUpkeep(upkeepId);
      if (info.balance < minBalanceThresholdMultiplier * minBalance) {
        addFunds(upkeepId, addLinkAmount);
        lastTopUpBlocks[upkeepId] = blockNum;
        emit UpkeepTopUp(upkeepId, addLinkAmount, blockNum);
      }
    }
  }

  function getMinBalanceForUpkeep(uint256 upkeepId) external view returns (uint96) {
    return registry.getMinBalanceForUpkeep(upkeepId);
  }

  function getForwarder(uint256 upkeepID) external view returns (address) {
    return registry.getForwarder(upkeepID);
  }

  function getBalance(uint256 id) external view returns (uint96 balance) {
    return registry.getBalance(id);
  }

  function getTriggerType(uint256 upkeepId) external view returns (uint8) {
    return registry.getTriggerType(upkeepId);
  }

  function burnPerformGas(uint256 upkeepId, uint256 startGas, uint256 blockNum) public {
    uint256 performGasToBurn = performGasToBurns[upkeepId];
    while (startGas - gasleft() + 10000 < performGasToBurn) {
      dummyMap[blockhash(blockNum)] = false;
    }
  }

  /**
   * @notice adds fund for an upkeep.
   * @param upkeepId the upkeep ID
   * @param amount the amount of LINK to be funded for the upkeep
   */
  function addFunds(uint256 upkeepId, uint96 amount) public {
    linkToken.approve(address(registry), amount);
    registry.addFunds(upkeepId, amount);
  }

  /**
   * @notice updates pipeline data for an upkeep. In order for the upkeep to be performed, the pipeline data must be the abi encoded upkeep ID.
   * @param upkeepId the upkeep ID
   * @param pipelineData the new pipeline data for the upkeep
   */
  function updateUpkeepPipelineData(uint256 upkeepId, bytes calldata pipelineData) external {
    registry.setUpkeepCheckData(upkeepId, pipelineData);
  }

  function withdrawLinks(uint256 upkeepId) external {
    registry.withdrawFunds(upkeepId, address(this));
  }

  function batchWithdrawLinks(uint256[] calldata upkeepIds) external {
    uint256 len = upkeepIds.length;
    for (uint32 i = 0; i < len; i++) {
      this.withdrawLinks(upkeepIds[i]);
    }
  }

  /**
   * @notice batch canceling upkeeps.
   * @param upkeepIds an array of upkeep IDs
   */
  function batchCancelUpkeeps(uint256[] calldata upkeepIds) external {
    uint256 len = upkeepIds.length;
    for (uint8 i = 0; i < len; i++) {
      registry.cancelUpkeep(upkeepIds[i]);
      s_upkeepIDs.remove(upkeepIds[i]);
    }
  }

  function eligible(uint256 upkeepId) public view returns (bool) {
    if (firstPerformBlocks[upkeepId] == 0) {
      return true;
    }
    return (getBlockNumber() - previousPerformBlocks[upkeepId]) >= intervals[upkeepId];
  }

  //  /**
  //   * @notice set a new add LINK amount.
  //   * @param amount the new value
  //   */
  //  function setAddLinkAmount(uint96 amount) external {
  //    addLinkAmount = amount;
  //  }
  //
  //  function setUpkeepTopUpCheckInterval(uint256 newInterval) external {
  //    upkeepTopUpCheckInterval = newInterval;
  //  }
  //
  //  function setMinBalanceThresholdMultiplier(uint8 newMinBalanceThresholdMultiplier) external {
  //    minBalanceThresholdMultiplier = newMinBalanceThresholdMultiplier;
  //  }

  //  function setPerformGasToBurn(uint256 upkeepId, uint256 value) public {
  //    performGasToBurns[upkeepId] = value;
  //  }
  //
  //  function setCheckGasToBurn(uint256 upkeepId, uint256 value) public {
  //    checkGasToBurns[upkeepId] = value;
  //  }

  function setPerformDataSize(uint256 upkeepId, uint256 value) public {
    performDataSizes[upkeepId] = value;
  }

  function setUpkeepGasLimit(uint256 upkeepId, uint32 gasLimit) public {
    registry.setUpkeepGasLimit(upkeepId, gasLimit);
    gasLimits[upkeepId] = gasLimit;
  }

  function setInterval(uint256 upkeepId, uint256 _interval) external {
    intervals[upkeepId] = _interval;
    firstPerformBlocks[upkeepId] = 0;
    counters[upkeepId] = 0;

    delete delays[upkeepId];
    uint16 currentBucket = buckets[upkeepId];
    for (uint16 i = 0; i <= currentBucket; i++) {
      delete bucketedDelays[upkeepId][i];
    }
    delete buckets[upkeepId];
  }

  /**
   * @notice batch setting intervals for an array of upkeeps.
   * @param upkeepIds an array of upkeep IDs
   * @param interval a new interval
   */
  function batchSetIntervals(uint256[] calldata upkeepIds, uint32 interval) external {
    uint256 len = upkeepIds.length;
    for (uint256 i = 0; i < len; i++) {
      this.setInterval(upkeepIds[i], interval);
    }
  }

  /**
   * @notice batch updating pipeline data for all upkeeps.
   * @param upkeepIds an array of upkeep IDs
   */
  function batchUpdatePipelineData(uint256[] calldata upkeepIds) external {
    uint256 len = upkeepIds.length;
    for (uint256 i = 0; i < len; i++) {
      uint256 upkeepId = upkeepIds[i];
      this.updateUpkeepPipelineData(upkeepId, abi.encode(upkeepId));
    }
  }

  /**
   * @notice finds all log trigger upkeeps and emits logs to serve as the initial trigger for upkeeps
   */
  function batchSendLogs(uint8 log) external {
    uint256[] memory upkeepIds = this.getActiveUpkeepIDsDeployedByThisContract(0, 0);
    uint256 len = upkeepIds.length;
    uint256 blockNum = getBlockNumber();
    for (uint256 i = 0; i < len; i++) {
      uint256 upkeepId = upkeepIds[i];
      uint8 triggerType = registry.getTriggerType(upkeepId);
      if (triggerType == 1) {
        if (log == 0) {
          emit LogEmitted(upkeepId, blockNum, address(this));
        } else {
          emit LogEmittedAgain(upkeepId, blockNum, address(this));
        }
      }
    }
  }

  function getUpkeepInfo(uint256 upkeepId) public view returns (IAutomationV21PlusCommon.UpkeepInfoLegacy memory) {
    return registry.getUpkeep(upkeepId);
  }

  function getUpkeepTriggerConfig(uint256 upkeepId) public view returns (bytes memory) {
    return registry.getUpkeepTriggerConfig(upkeepId);
  }

  function getUpkeepPrivilegeConfig(uint256 upkeepId) public view returns (bytes memory) {
    return registry.getUpkeepPrivilegeConfig(upkeepId);
  }

  function setUpkeepPrivilegeConfig(uint256 upkeepId, bytes memory cfg) external {
    registry.setUpkeepPrivilegeConfig(upkeepId, cfg);
  }

  function sendLog(uint256 upkeepId, uint8 log) external {
    uint256 blockNum = getBlockNumber();
    if (log == 0) {
      emit LogEmitted(upkeepId, blockNum, address(this));
    } else {
      emit LogEmittedAgain(upkeepId, blockNum, address(this));
    }
  }

  function getDelaysLength(uint256 upkeepId) public view returns (uint256) {
    return delays[upkeepId].length;
  }

  function getBucketedDelaysLength(uint256 upkeepId) public view returns (uint256) {
    uint16 currentBucket = buckets[upkeepId];
    uint256 len = 0;
    for (uint16 i = 0; i <= currentBucket; i++) {
      len += bucketedDelays[upkeepId][i].length;
    }
    return len;
  }

  function getDelays(uint256 upkeepId) public view returns (uint256[] memory) {
    return delays[upkeepId];
  }

  function getBucketedDelays(uint256 upkeepId, uint16 bucket) public view returns (uint256[] memory) {
    return bucketedDelays[upkeepId][bucket];
  }

  function getSumDelayLastNPerforms(uint256 upkeepId, uint256 n) public view returns (uint256, uint256) {
    uint256[] memory delays = delays[upkeepId];
    return getSumDelayLastNPerforms(delays, n);
  }

  function getSumDelayInBucket(uint256 upkeepId, uint16 bucket) public view returns (uint256, uint256) {
    uint256[] memory delays = bucketedDelays[upkeepId][bucket];
    return getSumDelayLastNPerforms(delays, delays.length);
  }

  function getSumDelayLastNPerforms(uint256[] memory delays, uint256 n) internal view returns (uint256, uint256) {
    uint256 i;
    uint256 len = delays.length;
    if (n == 0 || n >= len) {
      n = len;
    }
    uint256 sum = 0;

    for (i = 0; i < n; i++) sum = sum + delays[len - i - 1];
    return (sum, n);
  }

  function getPxDelayLastNPerforms(uint256 upkeepId, uint256 p, uint256 n) public view returns (uint256) {
    return getPxDelayLastNPerforms(delays[upkeepId], p, n);
  }

  function getPxDelayLastNPerforms(uint256[] memory delays, uint256 p, uint256 n) internal view returns (uint256) {
    uint256 i;
    uint256 len = delays.length;
    if (n == 0 || n >= len) {
      n = len;
    }
    uint256[] memory subArr = new uint256[](n);

    for (i = 0; i < n; i++) subArr[i] = (delays[len - i - 1]);
    quickSort(subArr, int256(0), int256(subArr.length - 1));

    if (p == 100) {
      return subArr[subArr.length - 1];
    }
    return subArr[(p * subArr.length) / 100];
  }

  function quickSort(uint256[] memory arr, int256 left, int256 right) private pure {
    int256 i = left;
    int256 j = right;
    if (i == j) return;
    uint256 pivot = arr[uint256(left + (right - left) / 2)];
    while (i <= j) {
      while (arr[uint256(i)] < pivot) i++;
      while (pivot < arr[uint256(j)]) j--;
      if (i <= j) {
        (arr[uint256(i)], arr[uint256(j)]) = (arr[uint256(j)], arr[uint256(i)]);
        i++;
        j--;
      }
    }
    if (left < j) quickSort(arr, left, j);
    if (i < right) quickSort(arr, i, right);
  }
}
