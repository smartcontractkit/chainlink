// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import "../automation/2_0/KeeperRegistrar2_0.sol";
import "../automation/2_0/KeeperRegistry2_0.sol";

/**
 * @notice this contract must have plenty LINKs bc it will check for every active upkeeps and top up
 * addLinkAmount if their balance is lower than minBalanceThresholdMultiplier * min balance.
 * if it does not have enough LINKs, upkeeps won't perform due to low LINK balance of this contract.
 */
contract UpkeepCounterStats is ConfirmedOwner {
  error IndexOutOfRange();

  event UpkeepsRegistered(uint256[] upkeepIds);
  event UpkeepsCancelled(uint256[] upkeepIds);
  event RegistrarSet(address newRegistrar);
  event FundsAdded(uint256 upkeepId, uint96 amount);
  event UpkeepTopUp(uint256 upkeepId, uint96 amount, uint256 blockNum);
  event InsufficientFunds(uint256 balance, uint256 blockNum);
  event Received(address sender, uint256 value);
  event PerformingUpkeep(
    uint256 firstPerformBlock,
    uint256 lastBlock,
    uint256 previousBlock,
    uint256 counter
  );

  using EnumerableSet for EnumerableSet.UintSet;

  mapping(uint256 => uint256) public lastTopUpBlocks;
  mapping(uint256 => uint256) public intervals;
  mapping(uint256 => uint256) public previousPerformBlocks;
  mapping(uint256 => uint256) public firstPerformBlocks;
  mapping(uint256 => uint256) public counters;
  mapping(uint256 => uint256) public performGasToBurns;
  mapping(uint256 => uint256) public checkGasToBurns;
  mapping(uint256 => uint256) public performDataSizes;
  mapping(uint256 => uint256) public gasLimits;
  mapping(uint256 => bytes) public checkDatas;
  mapping(bytes32 => bool) public dummyMap; // used to force storage lookup
  mapping(uint256 => uint256[]) public delays;  // how to query for delays for a certain past period: calendar day and/or past 24 hours

  mapping(uint256 => mapping(uint16 => uint256[])) public bucketedDelays;
  mapping(uint256 => mapping(uint16 => uint256[])) public timestampDelays;
  mapping(uint256 => uint256[]) public timestamps;
  mapping(uint256 => uint16) public timestampBuckets;
  mapping(uint256 => uint16) public buckets;
  EnumerableSet.UintSet internal s_upkeepIDs;
  KeeperRegistrar2_0 public registrar;
  LinkTokenInterface public linkToken;
  KeeperRegistry2_0 public registry;
  uint256 public upkeepTopUpCheckInterval = 10;
  uint96 public addLinkAmount = 5000000000000000000; // 5 LINK
  uint16 public immutable BUCKET_SIZE = 100;
  uint8 public minBalanceThresholdMultiplier = 100;

  constructor(address registrarAddress) ConfirmedOwner(msg.sender) {
    registrar = KeeperRegistrar2_0(registrarAddress);
    (,,, address registryAddress,) = registrar.getRegistrationConfig();
    registry = KeeperRegistry2_0(payable(address(registryAddress)));
    linkToken = registrar.LINK();
  }

  receive() external payable {
    emit Received(msg.sender, msg.value);
  }

  /**
   * @notice withdraws LINKs from this contract to msg sender when testing is finished.
   */
  function withdrawLinks() external onlyOwner {
    uint256 balance = linkToken.balanceOf(address(this));
    linkToken.transfer(msg.sender, balance);
  }

  /**
   * @notice sets registrar, registry, and link token address.
   * @param newRegistrar the new registrar address
   */
  function setConfig(KeeperRegistrar2_0 newRegistrar) external {
    registrar = newRegistrar;
    (,,, address registryAddress,) = registrar.getRegistrationConfig();
    registry = KeeperRegistry2_0(payable(address(registryAddress)));
    linkToken = registrar.LINK();

    emit RegistrarSet(address(registrar));
  }

  /**
   * @notice gets an array of active upkeep IDs.
   * @param startIndex the start index of upkeep IDs
   * @param maxCount the max number of upkeep IDs requested
   * @return an array of active upkeep IDs
   */
  function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) external view returns (uint256[] memory) {
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
   * @notice register an upkeep via the registrar.
   * @param params a registration params struct
   * @return an upkeep ID
   */
  function _registerUpkeep(KeeperRegistrar2_0.RegistrationParams memory params) private returns (uint256) {
    uint256 upkeepId = registrar.registerUpkeep(params);
    s_upkeepIDs.add(upkeepId);
    gasLimits[upkeepId] = params.gasLimit;
    checkDatas[upkeepId] = params.checkData;
    return upkeepId;
  }

  /**
   * @notice batch registering upkeeps.
   * @param number the number of upkeeps to be registered
   * @param gasLimit the gas limit of each upkeep
   * @param amount the amount of LINK to fund each upkeep
   */
  function batchRegisterUpkeeps(uint8 number, uint32 gasLimit, uint96 amount) external {
    KeeperRegistrar2_0.RegistrationParams memory params = KeeperRegistrar2_0.RegistrationParams({
      name: "test",
      encryptedEmail: bytes(""),
      upkeepContract: address(this),
      gasLimit: gasLimit,
      adminAddress: address(this), // use address of this contract as the admin
      checkData: bytes(""), // update check data later bc upkeep id is not available now
      offchainConfig: bytes(""),
      amount: amount
    });

    linkToken.approve(address(registrar), amount * number);

    uint256[] memory upkeepIds = new uint256[](number);
    for (uint8 i = 0; i < number; i++) {
      uint256 upkeepId = _registerUpkeep(params);
      upkeepIds[i] = upkeepId;
    }
    emit UpkeepsRegistered(upkeepIds);
  }

  /**
   * @notice adds fund for an upkeep.
   * @param upkeepId the upkeep ID
   * @param amount the amount of LINK to be funded for the upkeep
   */
  function addFunds(uint256 upkeepId, uint96 amount) external {
    linkToken.approve(address(registry), amount);
    registry.addFunds(upkeepId, amount);
    emit FundsAdded(upkeepId, amount);
  }

  /**
   * @notice updates check data for an upkeep. In order for the upkeep to be performed, the check data must be the abi encoded upkeep ID.
   * @param upkeepId the upkeep ID
   * @param checkData the new check data for the upkeep
   */
  function updateCheckData(uint256 upkeepId, bytes calldata checkData) external {
    registry.updateCheckData(upkeepId, checkData);
    checkDatas[upkeepId] = checkData;
  }

  /**
   * @notice cancel an upkeep.
   * @param upkeepId the upkeep ID
   */
  function cancelUpkeep(uint256 upkeepId) external {
    registry.cancelUpkeep(upkeepId);
    s_upkeepIDs.remove(upkeepId);

    delete lastTopUpBlocks[upkeepId];
    delete intervals[upkeepId];
    delete previousPerformBlocks[upkeepId];
    delete firstPerformBlocks[upkeepId];
    delete counters[upkeepId];
    delete performGasToBurns[upkeepId];
    delete checkGasToBurns[upkeepId];
    delete performDataSizes[upkeepId];
    delete gasLimits[upkeepId];
    delete checkDatas[upkeepId];
    delete delays[upkeepId];
    uint16 currentBucket = buckets[upkeepId];
    for (uint16 i = 0; i <= currentBucket; i++) {
      delete bucketedDelays[upkeepId][i];
    }
    delete buckets[upkeepId];
  }

  function withdrawLinks(uint256 upkeepId) external {
    registry.withdrawFunds(upkeepId, address(this));
  }

  /**
   * @notice batch canceling upkeeps.
   * @param upkeepIds an array of upkeep IDs
   */
  function batchCancelUpkeeps(uint256[] calldata upkeepIds) external {
    uint256 len = upkeepIds.length;
    for (uint8 i = 0; i < len; i++) {
      this.cancelUpkeep(upkeepIds[i]);
    }
    emit UpkeepsCancelled(upkeepIds);
  }

  function checkUpkeep(bytes calldata checkData) external returns (bool, bytes memory) {
    uint256 startGas = gasleft();
    (uint256 upkeepId) = abi.decode(
      checkData,
      (uint256)
    );

    uint256 performDataSize = performDataSizes[upkeepId];
    uint256 checkGasToBurn = checkGasToBurns[upkeepId];
    bytes memory pData = abi.encode(upkeepId, new bytes(performDataSize));
    uint256 blockNum = block.number;
    bool needed = eligible(upkeepId);
    while (startGas - gasleft() + 10000 < checkGasToBurn) {
      // 10K margin over gas to burn
      // Hard coded check gas to burn
      dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
      blockNum--;
    }
    return (needed, pData);
  }

  function performUpkeep(bytes calldata performData) external {
    uint256 startGas = gasleft();
    (uint256 upkeepId, ) = abi.decode(
      performData,
      (uint256, bytes)
    );
    uint256 firstPerformBlock = firstPerformBlocks[upkeepId];
    uint256 previousPerformBlock = previousPerformBlocks[upkeepId];
    uint256 blockNum = block.number;
    if (firstPerformBlock == 0) {
      firstPerformBlocks[upkeepId] = blockNum;
      firstPerformBlock = blockNum;
      timestamps[upkeepId].push(block.timestamp);
    } else {
      // Calculate and append delay
      uint256 delay = blockNum - previousPerformBlock - intervals[upkeepId];

      uint16 timestampBucket = timestampBuckets[upkeepId];
      if (block.timestamp - 3600 > timestamps[upkeepId][timestampBucket]) {
        timestamps[upkeepId].push(block.timestamp);
        timestampBucket++;
        timestampBuckets[upkeepId] = timestampBucket;
      }

      uint16 bucket = buckets[upkeepId];
      uint256[] memory bucketDelays = bucketedDelays[upkeepId][bucket];
      if (bucketDelays.length == BUCKET_SIZE) {
        bucket++;
        buckets[upkeepId] = bucket;
      }
      bucketedDelays[upkeepId][bucket].push(delay);
      timestampDelays[upkeepId][timestampBucket].push(delay);
      delays[upkeepId].push(delay);
    }

    uint256 counter = counters[upkeepId] + 1;
    counters[upkeepId] = counter;
    emit PerformingUpkeep(firstPerformBlock, blockNum, previousPerformBlock, counter);
    previousPerformBlocks[upkeepId] = blockNum;

    // check the upkeep to see if it needs extra funds
    if (blockNum - lastTopUpBlocks[upkeepId] > upkeepTopUpCheckInterval) {
      UpkeepInfo memory info = registry.getUpkeep(upkeepId);
      uint96 minBalance = registry.getMinBalanceForUpkeep(upkeepId);
      if (info.balance < minBalanceThresholdMultiplier * minBalance) {
        this.addFunds(upkeepId, addLinkAmount);
        lastTopUpBlocks[upkeepId] = blockNum;
        emit UpkeepTopUp(upkeepId, addLinkAmount, blockNum);
      }
    }

    uint256 performGasToBurn = performGasToBurns[upkeepId];
    while (startGas - gasleft() + 10000 < performGasToBurn) {
      // 10K margin over gas to burn
      dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
      blockNum--;
    }
  }

  function eligible(uint256 upkeepId) public view returns (bool) {
    if (firstPerformBlocks[upkeepId] == 0) {
      return true;
    }
    return (block.number - previousPerformBlocks[upkeepId]) >= intervals[upkeepId];
  }

  /**
   * @notice set a new add LINK amount.
   * @param amount the new value
   */
  function setAddLinkAmount(uint96 amount) external {
    addLinkAmount = amount;
  }

  function setUpkeepTopUpCheckInterval(uint256 newInterval) external {
    upkeepTopUpCheckInterval = newInterval;
  }

  function setMinBalanceThresholdMultiplier(uint8 newMinBalanceThresholdMultiplier) external {
    minBalanceThresholdMultiplier = newMinBalanceThresholdMultiplier;
  }

  function setPerformGasToBurn(uint256 upkeepId, uint256 value) public {
    performGasToBurns[upkeepId] = value;
  }

  function setCheckGasToBurn(uint256 upkeepId, uint256 value) public {
    checkGasToBurns[upkeepId] = value;
  }

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

    currentBucket = timestampBuckets[upkeepId];
    for (uint16 i = 0; i <= currentBucket; i++) {
      delete timestampDelays[upkeepId][i];
    }
    delete timestamps[upkeepId];
    delete timestampBuckets[upkeepId];
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
   * @notice batch updating check data for all upkeeps.
   * @param upkeepIds an array of upkeep IDs
   */
  function batchUpdateCheckData(uint256[] calldata upkeepIds) external {
    uint256 len = upkeepIds.length;
    for (uint256 i = 0; i < len; i++) {
      uint256 upkeepId = upkeepIds[i];
      this.updateCheckData(upkeepId, abi.encode(upkeepId));
    }
  }

  function getDelaysLength(uint256 upkeepId) public view returns (uint256) {
    return delays[upkeepId].length;
  }

  function getDelaysLengthAtBucket(uint256 upkeepId, uint16 bucket) public view returns (uint256) {
    return bucketedDelays[upkeepId][bucket].length;
  }

  function getDelaysLengthAtTimestampBucket(uint256 upkeepId, uint16 timestampBucket) public view returns (uint256) {
    return timestampDelays[upkeepId][timestampBucket].length;
  }

  function getBucketedDelaysLength(uint256 upkeepId) public view returns (uint256) {
    uint16 currentBucket = buckets[upkeepId];
    uint256 len = 0;
    for (uint16 i = 0; i <= currentBucket; i++) {
      len += bucketedDelays[upkeepId][i].length;
    }
    return len;
  }

  function getTimestampBucketedDelaysLength(uint256 upkeepId) public view returns (uint256) {
    uint16 timestampBucket = timestampBuckets[upkeepId];
    uint256 len = 0;
    for (uint16 i = 0; i <= timestampBucket; i++) {
      len += timestampDelays[upkeepId][i].length;
    }
    return len;
  }

  function getDelays(uint256 upkeepId) public view returns (uint256[] memory) {
    return delays[upkeepId];
  }

  function getTimestampDelays(uint256 upkeepId, uint16 timestampBucket) public view returns (uint256[] memory) {
    return timestampDelays[upkeepId][timestampBucket];
  }

  function getBucketedDelays(uint256 upkeepId, uint16 bucket) public view returns (uint256[] memory) {
    return bucketedDelays[upkeepId][bucket];
  }

  function getSumDelayLastNPerforms(uint256 upkeepId, uint256 n) public view returns (uint256, uint256) {
    uint256[] memory delays = delays[upkeepId];
    return getSumDelayLastNPerforms(delays, n);
  }

  function getSumBucketedDelayLastNPerforms(uint256 upkeepId, uint256 n) public view returns (uint256, uint256) {
    uint256 len = this.getBucketedDelaysLength(upkeepId);
    if (n == 0 || n >= len) {
      n = len;
    }
    uint256 nn = n;
    uint256 sum = 0;
    uint16 currentBucket = buckets[upkeepId];
    for (uint16 i = currentBucket; i >= 0; i--) {
      uint256[] memory delays = bucketedDelays[upkeepId][i];
      (uint256 s, uint256 m) = getSumDelayLastNPerforms(delays, nn);
      sum += s;
      nn -= m;
      if (nn <= 0) {
        break;
      }
    }
    return (sum, n);
  }

  function getSumTimestampBucketedDelayLastNPerforms(uint256 upkeepId, uint256 n) public view returns (uint256, uint256) {
    uint256 len = this.getTimestampBucketedDelaysLength(upkeepId);
    if (n == 0 || n >= len) {
      n = len;
    }
    uint256 nn = n;
    uint256 sum = 0;
    uint16 timestampBucket = timestampBuckets[upkeepId];
    for (uint16 i = timestampBucket; i >= 0; i--) {
      uint256[] memory delays = timestampDelays[upkeepId][i];
      (uint256 s, uint256 m) = getSumDelayLastNPerforms(delays, nn);
      sum += s;
      nn -= m;
      if (nn <= 0) {
        break;
      }
    }
    return (sum, n);
  }

  function getSumDelayInBucket(uint256 upkeepId, uint16 bucket) public view returns (uint256, uint256) {
    uint256[] memory delays = bucketedDelays[upkeepId][bucket];
    return getSumDelayLastNPerforms(delays, delays.length);
  }

  function getSumDelayInTimestampBucket(uint256 upkeepId, uint16 timestampBucket) public view returns (uint256, uint256) {
    uint256[] memory delays = timestampDelays[upkeepId][timestampBucket];
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

  function getPxDelayForAllUpkeeps(uint256 p) public view returns (uint256[] memory, uint256[] memory) {
    uint256 len = s_upkeepIDs.length();
    uint256[] memory upkeepIds = new uint256[](len);
    uint256[] memory pxDelays = new uint256[](len);

    for (uint256 idx = 0; idx < len; idx++) {
      uint256 upkeepId = s_upkeepIDs.at(idx);
      uint256[] memory delays = delays[upkeepId];
      upkeepIds[idx] = upkeepId;
      pxDelays[idx] = getPxDelayLastNPerforms(delays, p, delays.length);
    }

    return (upkeepIds, pxDelays);
  }

  function getPxBucketedDelaysForAllUpkeeps(uint256 p) public view returns (uint256[] memory, uint256[] memory) {
    uint256 len = s_upkeepIDs.length();
    uint256[] memory upkeepIds = new uint256[](len);
    uint256[] memory pxDelays = new uint256[](len);

    for (uint256 idx = 0; idx < len; idx++) {
      uint256 upkeepId = s_upkeepIDs.at(idx);
      upkeepIds[idx] = upkeepId;
      uint16 currentBucket = buckets[upkeepId];
      uint256 delayLen = this.getBucketedDelaysLength(upkeepId);
      uint256[] memory delays = new uint256[](delayLen);
      uint256 i = 0;
      mapping(uint16 => uint256[]) storage bucketedDelays = bucketedDelays[upkeepId];
      for (uint16 j = 0; j <= currentBucket; j++) {
        uint256[] memory d = bucketedDelays[j];
        for (uint256 k = 0; k < d.length; k++) {
          delays[i++] = d[k];
        }
      }
      pxDelays[idx] = getPxDelayLastNPerforms(delays, p, delayLen);
    }

    return (upkeepIds, pxDelays);
  }

  function getPxDelayInTimestampBucket(uint256 upkeepId, uint256 p, uint16 timestampBucket) public view returns (uint256) {
    uint256[] memory delays = timestampDelays[upkeepId][timestampBucket];
    return getPxDelayLastNPerforms(delays, p, delays.length);
  }

  function getPxDelayInBucket(uint256 upkeepId, uint256 p, uint16 bucket) public view returns (uint256) {
    uint256[] memory delays = bucketedDelays[upkeepId][bucket];
    return getPxDelayLastNPerforms(delays, p, delays.length);
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
      return  subArr[subArr.length - 1];
    }
    return subArr[(p * subArr.length) / 100];
  }

  function quickSort(
    uint256[] memory arr,
    int256 left,
    int256 right
  ) private pure {
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
