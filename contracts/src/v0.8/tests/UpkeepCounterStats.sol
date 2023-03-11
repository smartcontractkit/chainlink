pragma solidity ^0.8.6;

import "../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import "../automation/2_0/KeeperRegistrar2_0.sol";

contract UpkeepCounterStats {
  using EnumerableSet for EnumerableSet.UintSet;
  event PerformingUpkeep(
    uint256 initialBlock,
    uint256 lastBlock,
    uint256 previousBlock,
    uint256 counter,
    bytes performData
  );

  mapping(uint256 => uint256) public upkeepIdsToIntervals;
  mapping(uint256 => uint256) public upkeepIdsToPreviousPerformBlock;
  mapping(uint256 => uint256) public upkeepIdsToInitialBlock;
  mapping(uint256 => uint256) public upkeepIdsToCounter;
  mapping(uint256 => uint256) public upkeepIdsToPerformGasToBurn;
  mapping(uint256 => uint256) public upkeepIdsToCheckGasToBurn;
  mapping(uint256 => uint256) public upkeepIdsToPerformDataSize;
  mapping(bytes32 => bool) public dummyMap; // used to force storage lookup
  mapping(uint256 => uint256[]) private upkeepIdsToDelay;
  EnumerableSet.UintSet internal s_upkeepIDs;
  KeeperRegistrar2_0 public registrar;
  AutomationRegistryBaseInterface public registry;
  LinkTokenInterface public linkToken;

  constructor(address registrarAddress) {
    registrar = KeeperRegistrar2_0(registrarAddress);
    (,,, address registryAddress,) = registrar.getRegistrationConfig();
    registry = AutomationRegistryBaseInterface(registryAddress);
    linkToken = registrar.LINK();
  }

  function registerUpkeep(RegistrationParams memory params) external returns (uint256) {
    uint256 upkeepId = registrar.registerUpkeep(params);
    s_upkeepIDs.add(upkeepId);
    return upkeepId;
  }

  function batchRegisterUpkeeps(uint8 number, uint32 gasLimit, uint96 amount) external returns (uint256[] memory) {
    RegistrationParams memory params = RegistrationParams({
      name: "test",
      encryptedEmail: '0x',
      upkeepContract: address(this),
      gasLimit: gasLimit,
      adminAddress: address(this), // cannot use msg.sender otherwise updateCheckData won't work
      checkData: '0x', // update check data later bc upkeep id is not available now
      offchainConfig: '0x',
      amount: amount
    });

    uint256[] memory upkeepIds = new uint256[](number);
    for (uint8 i = 0; i < number; i++) {
      uint256 upkeepId = this.registerUpkeep(params);
      upkeepIds[i] = upkeepId;
    }
    return upkeepIds;
  }

  function addFunds(uint256 upkeepId, uint96 amount) external {
    linkToken.approve(address(registry), amount);
    registry.addFunds(upkeepId, amount);
  }

  function updateCheckData(uint256 upkeepId, bytes calldata checkData) external {
    registry.updateCheckData(upkeepId, checkData);
  }

  function cancelUpkeep(uint256 upkeepId) external {
    registry.cancelUpkeep(upkeepId);
    s_upkeepIDs.remove(upkeepId);
    // keep data in mappings in case needed afterwards?
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
    (uint256 upkeepId, bytes memory performDataPlaceHolder) = abi.decode(
      performData,
      (uint256, bytes)
    );
    uint256 performGasToBurn = upkeepIdsToPerformGasToBurn[upkeepId];
    uint256 initialBlock = upkeepIdsToInitialBlock[upkeepId];
    uint256 blockNum = block.number;
    uint256 interval = upkeepIdsToIntervals[upkeepId];
    if (initialBlock == 0) {
      upkeepIdsToInitialBlock[upkeepId] = blockNum;
      initialBlock = blockNum;
    } else {
      // Calculate and append delay
      uint256 delay = blockNum - upkeepIdsToPreviousPerformBlock[upkeepId] - interval;
      upkeepIdsToDelay[upkeepId].push(delay);
    }

    //upkeepIdsToLastBlock[upkeepId] = blockNum;
    uint256 counter = upkeepIdsToCounter[upkeepId] + 1;
    upkeepIdsToCounter[upkeepId] = counter;
    emit PerformingUpkeep(initialBlock, blockNum, upkeepIdsToPreviousPerformBlock[upkeepId], counter, performData);
    upkeepIdsToPreviousPerformBlock[upkeepId] = blockNum;

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

  function setPerformGasToBurn(uint256 upkeepId, uint256 value) public {
    upkeepIdsToPerformGasToBurn[upkeepId] = value;
  }

  function setCheckGasToBurn(uint256 upkeepId, uint256 value) public {
    upkeepIdsToCheckGasToBurn[upkeepId] = value;
  }

  function setPerformDataSize(uint256 upkeepId, uint256 value) public {
    upkeepIdsToPerformDataSize[upkeepId] = value;
  }

  function setSpread(uint256 upkeepId, uint256 _interval) external {
    upkeepIdsToIntervals[upkeepId] = _interval;
    upkeepIdsToInitialBlock[upkeepId] = 0;
    upkeepIdsToCounter[upkeepId] = 0;

    delete upkeepIdsToDelay[upkeepId];
  }

  function getDelaysLength(uint256 upkeepId) public view returns (uint256) {
    return upkeepIdsToDelay[upkeepId].length;
  }

  function getDelays(uint256 upkeepId) public view returns (uint256[] memory) {
    return upkeepIdsToDelay[upkeepId];
  }

  function getSumDelayLastNPerforms(uint256 upkeepId, uint256 n) public view returns (uint256, uint256) {
    uint256[] memory delays = upkeepIdsToDelay[upkeepId];
    uint256 i;
    uint256 len = delays.length;
    if (n == 0 || n >= len) {
      n = len;
    }
    uint256 sum = 0;

    for (i = 0; i < n; i++) sum = sum + delays[len - i - 1];
    return (sum, n);
  }

  function getPxDelayForAllUpkeeps(uint256 p) public view returns (uint256[] memory upkeepIds, uint256[] memory pxDelays) {
    uint256 len = s_upkeepIDs.length();
    uint256[] memory upkeepIds = new uint256[](len);
    uint256[] memory pxDelays = new uint256[](len);

    for (uint256 idx = 0; idx < len; idx++) {
      uint256 upkeepId = s_upkeepIDs.at(idx);
      uint256[] memory delays = upkeepIdsToDelay[upkeepId];
      upkeepIds[idx] = upkeepId;
      pxDelays[idx] = this.getPxDelayLastNPerforms(upkeepId, p, delays.length);
    }

    return (upkeepIds, pxDelays);
  }

  function getPxDelayLastNPerforms(uint256 upkeepId, uint256 p, uint256 n) public view returns (uint256) {
    uint256[] memory delays = upkeepIdsToDelay[upkeepId];
    uint256 i;
    uint256 len = delays.length;
    if (n == 0 || n >= len) {
      n = len;
    }
    uint256[] memory subArr = new uint256[](n);

    for (i = 0; i < n; i++) subArr[i] = (delays[len - i - 1]);
    quickSort(subArr, int256(0), int256(subArr.length - 1));

    uint256 index = (p * subArr.length) / 100;
    return subArr[index];
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
