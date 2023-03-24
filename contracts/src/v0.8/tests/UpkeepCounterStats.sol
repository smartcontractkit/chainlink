// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import "../automation/2_0/KeeperRegistrar2_0.sol";
import "../automation/2_0/KeeperRegistry2_0.sol";

/**
 * @notice this contract must have plenty LINKs bc it will check for every active upkeeps and top up
 * addFundsMinBalanceMultiplier * min balance if their balance is lower than minBalanceThresholdMultiplier * min.
 * if it does not have enough LINKs, upkeeps won't perform due to low LINK balance of this contract.
 * this contract also must have plenty native tokens if we want to use the topUpTransmitters function.
 */
contract UpkeepCounterStats {
  error IndexOutOfRange();

  event UpkeepsRegistered(uint256[] upkeepIds);
  event UpkeepsCancelled(uint256[] upkeepIds);
  event RegistrarSet(address newRegistrar);
  event FundsAdded(uint256 upkeepId, uint256 amount);
  event TransmitterTopUp(address transmitter, uint256 amount, uint256 blockNum);
  event UpkeepTopUp(uint256 upkeepId, uint256 amount, uint256 blockNum);
  event InsufficientFunds(uint256 balance, uint256 blockNum);
  event Received(address sender, uint256 value);
  event PerformingUpkeep(
    uint256 initialBlock,
    uint256 lastBlock,
    uint256 previousBlock,
    uint256 counter,
    bytes performData
  );

  using EnumerableSet for EnumerableSet.UintSet;

  mapping(uint256 => uint256) public upkeepIdsToLastTopUpBlock;
  mapping(uint256 => uint256) public upkeepIdsToIntervals;
  mapping(uint256 => uint256) public upkeepIdsToPreviousPerformBlock;
  mapping(uint256 => uint256) public upkeepIdsToInitialBlock;
  mapping(uint256 => uint256) public upkeepIdsToCounter;
  mapping(uint256 => uint256) public upkeepIdsToPerformGasToBurn;
  mapping(uint256 => uint256) public upkeepIdsToCheckGasToBurn;
  mapping(uint256 => uint256) public upkeepIdsToPerformDataSize;
  mapping(uint256 => uint256) public upkeepIdsToGasLimit;
  mapping(uint256 => bytes) public upkeepIdsToCheckData;
  mapping(bytes32 => bool) public dummyMap; // used to force storage lookup
  mapping(uint256 => uint256[]) public upkeepIdsToDelay;  // how to query for delays for a certain past period: calendar day and/or past 24 hours

  mapping(uint256 => mapping(uint16 => uint256[])) public upkeepIdsToBucketedDelays;
  mapping(uint256 => uint16) public upkeepIdsToBucket;
  EnumerableSet.UintSet internal s_upkeepIDs;
  KeeperRegistrar2_0 public registrar;
  LinkTokenInterface public linkToken;
  KeeperRegistry2_0 public registry;
  uint256 public lastTransmittersTopUpBlock;
  uint256 public upkeepTopUpCheckInterval = 2000;
  uint256 public transmitterTopUpCheckInterval = 2000;
  uint96 public transmitterMinBalance = 5000000000000000000;
  uint96 public transmitterAddBalance = 20000000000000000000;
  uint8 public minBalanceThresholdMultiplier = 50;
  uint8 public addFundsMinBalanceMultiplier = 100;

  constructor(address registrarAddress) {
    registrar = KeeperRegistrar2_0(registrarAddress);
    (,,, address registryAddress,) = registrar.getRegistrationConfig();
    registry = KeeperRegistry2_0(payable(address(registryAddress)));
    linkToken = registrar.LINK();
  }

  receive() external payable {
    emit Received(msg.sender, msg.value);
  }

  function fundLink(uint256 amount) external {
    linkToken.approve(msg.sender, amount);
    linkToken.transferFrom(msg.sender, address(this), amount);
  }

  function setRegistrar(KeeperRegistrar2_0 newRegistrar) external {
    registrar = newRegistrar;
    (,,, address registryAddress,) = registrar.getRegistrationConfig();
    registry = KeeperRegistry2_0(payable(address(registryAddress)));
    linkToken = registrar.LINK();

    emit RegistrarSet(address(registrar));
  }

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

  function _registerUpkeep(KeeperRegistrar2_0.RegistrationParams memory params) private returns (uint256) {
    uint256 upkeepId = registrar.registerUpkeep(params);
    s_upkeepIDs.add(upkeepId);
    upkeepIdsToGasLimit[upkeepId] = params.gasLimit;
    upkeepIdsToCheckData[upkeepId] = params.checkData;
    return upkeepId;
  }

  function batchRegisterUpkeeps(uint8 number, uint32 gasLimit, uint96 amount) external {
    KeeperRegistrar2_0.RegistrationParams memory params = KeeperRegistrar2_0.RegistrationParams({
      name: "test",
      encryptedEmail: bytes(""),
      upkeepContract: address(this),
      gasLimit: gasLimit,
      adminAddress: address(this), // cannot use msg.sender otherwise updateCheckData won't work
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

  function addFunds(uint256 upkeepId, uint96 amount) external {
    linkToken.approve(address(registry), amount);
    registry.addFunds(upkeepId, amount);
    emit FundsAdded(upkeepId, amount);
  }

  function updateCheckData(uint256 upkeepId, bytes calldata checkData) external {
    registry.updateCheckData(upkeepId, checkData);
    upkeepIdsToCheckData[upkeepId] = checkData;
  }

  function _cancelUpkeep(uint256 upkeepId) private {
    registry.cancelUpkeep(upkeepId);
    s_upkeepIDs.remove(upkeepId);
    // keep data in mappings in case needed afterwards?
  }

  function batchCancelUpkeeps(uint256[] calldata upkeepIds) external {
    uint256 len = upkeepIds.length;
    for (uint8 i = 0; i < len; i++) {
      _cancelUpkeep(upkeepIds[i]);
    }
    emit UpkeepsCancelled(upkeepIds);
  }

  function checkUpkeep(bytes calldata checkData) external returns (bool, bytes memory) {
    uint256 startGas = gasleft();
    (uint256 upkeepId) = abi.decode(
      checkData,
      (uint256)
    );

    uint256 performDataSize = upkeepIdsToPerformDataSize[upkeepId];
    uint256 checkGasToBurn = upkeepIdsToCheckGasToBurn[upkeepId];
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
    uint256 initialBlock = upkeepIdsToInitialBlock[upkeepId];
    uint256 blockNum = block.number;
    if (initialBlock == 0) {
      upkeepIdsToInitialBlock[upkeepId] = blockNum;
      initialBlock = blockNum;
    } else {
      // Calculate and append delay
      uint256 delay = blockNum - upkeepIdsToPreviousPerformBlock[upkeepId] - upkeepIdsToIntervals[upkeepId];

      uint16 bucket = upkeepIdsToBucket[upkeepId];
      uint256[] memory bucketedDelays = upkeepIdsToBucketedDelays[upkeepId][bucket];
      if (bucketedDelays.length == 100) {
        bucket++;
      }
      upkeepIdsToBucketedDelays[upkeepId][bucket].push(delay);
      upkeepIdsToDelay[upkeepId].push(delay);
    }

    uint256 counter = upkeepIdsToCounter[upkeepId] + 1;
    upkeepIdsToCounter[upkeepId] = counter;
    emit PerformingUpkeep(initialBlock, blockNum, upkeepIdsToPreviousPerformBlock[upkeepId], counter, performData);
    upkeepIdsToPreviousPerformBlock[upkeepId] = blockNum;

    // every upkeep adds funds for themselves
    if (blockNum - upkeepIdsToLastTopUpBlock[upkeepId] > upkeepTopUpCheckInterval) {
      UpkeepInfo memory info = registry.getUpkeep(upkeepId);
      uint96 minBalance = registry.getMinBalanceForUpkeep(upkeepId);
      if (info.balance < minBalanceThresholdMultiplier * minBalance) {
        this.addFunds(upkeepId, addFundsMinBalanceMultiplier * minBalance);
        upkeepIdsToLastTopUpBlock[upkeepId] = blockNum;
        emit UpkeepTopUp(upkeepId, addFundsMinBalanceMultiplier * minBalance, blockNum);
      }
    }

    uint256 performGasToBurn = upkeepIdsToPerformGasToBurn[upkeepId];
    while (startGas - gasleft() + 10000 < performGasToBurn) {
      // 10K margin over gas to burn
      dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
      blockNum--;
    }
  }

  function eligible(uint256 upkeepId) public view returns (bool) {
    if (upkeepIdsToInitialBlock[upkeepId] == 0) {
      return true;
    }
    return (block.number - upkeepIdsToPreviousPerformBlock[upkeepId]) >= upkeepIdsToIntervals[upkeepId];
  }

  function setUpkeepTopUpCheckInterval(uint256 newInterval) external {
    upkeepTopUpCheckInterval = newInterval;
  }

  function setTransmitterTopUpCheckInterval(uint256 newInterval) external {
    transmitterTopUpCheckInterval = newInterval;
  }

  function setMinBalanceMultipliers(uint8 newMinBalanceThresholdMultiplier, uint8 newAddFundsMinBalanceMultiplier) external {
    minBalanceThresholdMultiplier = newMinBalanceThresholdMultiplier;
    addFundsMinBalanceMultiplier = newAddFundsMinBalanceMultiplier;
  }

  function setTransmitterBalanceLimit(uint96 newTransmitterMinBalance, uint96 newTransmitterAddBalance) external {
    transmitterMinBalance = newTransmitterMinBalance;
    transmitterAddBalance = newTransmitterAddBalance;
  }

  function setPerformGasToBurn(uint256 upkeepId, uint256 value) public {
    upkeepIdsToPerformGasToBurn[upkeepId] = value;
  }

  function setCheckGasToBurn(uint256 upkeepId, uint256 value) public {
    upkeepIdsToCheckGasToBurn[upkeepId] = value;
  }

  function setPerformDataSize(uint256 upkeepId, uint256 value) public {
    upkeepIdsToPerformDataSize[upkeepId] = value;
  }

  function setUpkeepGasLimit(uint256 upkeepId, uint32 gasLimit) public {
    registry.setUpkeepGasLimit(upkeepId, gasLimit);
    upkeepIdsToGasLimit[upkeepId] = gasLimit;
  }

  function setInterval(uint256 upkeepId, uint256 _interval) external {
    upkeepIdsToIntervals[upkeepId] = _interval;
    upkeepIdsToInitialBlock[upkeepId] = 0;
    upkeepIdsToCounter[upkeepId] = 0;

    delete upkeepIdsToDelay[upkeepId];
    uint16 currentBucket = upkeepIdsToBucket[upkeepId];
    for (uint16 i = 0; i <= currentBucket; i++) {
      delete upkeepIdsToBucketedDelays[upkeepId][i];
    }
    upkeepIdsToBucket[upkeepId] = 0;
  }

  function batchSetIntervals(uint32 interval) external {
    uint256 len = s_upkeepIDs.length();
    for (uint256 i = 0; i < len; i++) {
      uint256 upkeepId = s_upkeepIDs.at(i);
      this.setInterval(upkeepId, interval);
    }
  }

  function batchUpdateCheckData() external {
    uint256 len = s_upkeepIDs.length();
    for (uint256 i = 0; i < len; i++) {
      uint256 upkeepId = s_upkeepIDs.at(i);
      this.updateCheckData(upkeepId, abi.encode(upkeepId));
    }
  }

  function getDelaysLength(uint256 upkeepId) public view returns (uint256) {
    return upkeepIdsToDelay[upkeepId].length;
  }

  function getDelaysLengthAtBucket(uint256 upkeepId, uint16 bucket) public view returns (uint256) {
    return upkeepIdsToBucketedDelays[upkeepId][bucket].length;
  }

  function getBucketedDelaysLength(uint256 upkeepId) public view returns (uint256) {
    uint16 currentBucket = upkeepIdsToBucket[upkeepId];
    uint256 len = 0;
    for (uint16 i = 0; i <= currentBucket; i++) {
      len += upkeepIdsToBucketedDelays[upkeepId][i].length;
    }
    return len;
  }

  function getDelays(uint256 upkeepId) public view returns (uint256[] memory) {
    return upkeepIdsToDelay[upkeepId];
  }

  function getBucketedDelays(uint256 upkeepId, uint16 bucket) public view returns (uint256[] memory) {
    return upkeepIdsToBucketedDelays[upkeepId][bucket];
  }

  function getSumDelayLastNPerforms(uint256 upkeepId, uint256 n) public view returns (uint256, uint256) {
    uint256[] memory delays = upkeepIdsToDelay[upkeepId];
    return getSumDelayLastNPerforms(delays, n);
  }

  function getSumDelayLastNPerforms1(uint256 upkeepId, uint256 n) public view returns (uint256, uint256) {
    uint256 len = this.getBucketedDelaysLength(upkeepId);
    if (n == 0 || n >= len) {
      n = len;
    }
    uint256 nn = n;
    uint256 sum = 0;
    uint16 currentBucket = upkeepIdsToBucket[upkeepId];
    for (uint16 i = currentBucket; i >= 0; i--) {
      uint256[] memory delays = upkeepIdsToBucketedDelays[upkeepId][i];
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
    uint256[] memory delays = upkeepIdsToBucketedDelays[upkeepId][bucket];
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
      uint256[] memory delays = upkeepIdsToDelay[upkeepId];
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
      uint16 currentBucket = upkeepIdsToBucket[upkeepId];
      uint256 delayLen = this.getBucketedDelaysLength(upkeepId);
      uint256[] memory delays = new uint256[](delayLen);
      uint256 i = 0;
      mapping(uint16 => uint256[]) storage bucketedDelays = upkeepIdsToBucketedDelays[upkeepId];
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

  function getPxDelayInBucket(uint256 upkeepId, uint256 p, uint16 bucket) public view returns (uint256) {
    uint256[] memory delays = upkeepIdsToBucketedDelays[upkeepId][bucket];
    return getPxDelayLastNPerforms(delays, p, delays.length);
  }

  function getPxDelayLastNPerforms(uint256 upkeepId, uint256 p, uint256 n) public view returns (uint256) {
    return getPxDelayLastNPerforms(upkeepIdsToDelay[upkeepId], p, n);
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

  function topUpTransmitters() external {
    (,,, address[] memory transmitters, ) = registry.getState();
    uint256 len = transmitters.length;
    uint256 blockNum = block.number;
    for (uint256 i = 0; i < len; i++) {
      if (transmitters[i].balance < transmitterMinBalance) {
        if (address(this).balance < transmitterAddBalance) {
          emit InsufficientFunds(address(this).balance, blockNum);
        } else {
          lastTransmittersTopUpBlock = blockNum;
          transmitters[i].call{value: transmitterAddBalance}("");
          emit TransmitterTopUp(transmitters[i], transmitterAddBalance, blockNum);
        }
      }
    }
  }
}
