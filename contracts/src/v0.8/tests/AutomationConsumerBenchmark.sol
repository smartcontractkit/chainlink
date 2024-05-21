pragma solidity 0.8.16;

contract AutomationConsumerBenchmark {
  event PerformingUpkeep(uint256 id, address from, uint256 initialCall, uint256 nextEligible, uint256 blockNumber);

  mapping(uint256 => uint256) public initialCall;
  mapping(uint256 => uint256) public nextEligible;
  mapping(bytes32 => bool) public dummyMap; // used to force storage lookup
  mapping(uint256 => uint256) public count;
  uint256 deployedAt;

  constructor() {
    deployedAt = block.number;
  }

  function checkUpkeep(bytes calldata checkData) external view returns (bool, bytes memory) {
    (
      uint256 id,
      uint256 interval,
      uint256 range,
      uint256 checkBurnAmount,
      uint256 performBurnAmount,
      uint256 firstEligibleBuffer
    ) = abi.decode(checkData, (uint256, uint256, uint256, uint256, uint256, uint256));
    uint256 startGas = gasleft();
    bytes32 dummyIndex = blockhash(block.number - 1);
    bool dummy;
    // burn gas
    if (checkBurnAmount > 0 && eligible(id, range, firstEligibleBuffer)) {
      while (startGas - gasleft() < checkBurnAmount) {
        dummy = dummy && dummyMap[dummyIndex]; // arbitrary storage reads
        dummyIndex = keccak256(abi.encode(dummyIndex, address(this)));
      }
    }
    return (eligible(id, range, firstEligibleBuffer), checkData);
  }

  function performUpkeep(bytes calldata performData) external {
    (
      uint256 id,
      uint256 interval,
      uint256 range,
      uint256 checkBurnAmount,
      uint256 performBurnAmount,
      uint256 firstEligibleBuffer
    ) = abi.decode(performData, (uint256, uint256, uint256, uint256, uint256, uint256));
    require(eligible(id, range, firstEligibleBuffer));
    uint256 startGas = gasleft();
    if (initialCall[id] == 0) {
      initialCall[id] = block.number;
    }
    nextEligible[id] = block.number + interval;
    count[id]++;
    emit PerformingUpkeep(id, tx.origin, initialCall[id], nextEligible[id], block.number);
    // burn gas
    bytes32 dummyIndex = blockhash(block.number - 1);
    bool dummy;
    while (startGas - gasleft() < performBurnAmount) {
      dummy = dummy && dummyMap[dummyIndex]; // arbitrary storage reads
      dummyIndex = keccak256(abi.encode(dummyIndex, address(this)));
    }
  }

  function getCountPerforms(uint256 id) public view returns (uint256) {
    return count[id];
  }

  function eligible(uint256 id, uint256 range, uint256 firstEligibleBuffer) internal view returns (bool) {
    return
      initialCall[id] == 0
        ? block.number >= firstEligibleBuffer + deployedAt
        : (block.number - initialCall[id] < range && block.number > nextEligible[id]);
  }

  function checkEligible(uint256 id, uint256 range, uint256 firstEligibleBuffer) public view returns (bool) {
    return eligible(id, range, firstEligibleBuffer);
  }

  function reset() external {
    deployedAt = block.number;
  }
}
