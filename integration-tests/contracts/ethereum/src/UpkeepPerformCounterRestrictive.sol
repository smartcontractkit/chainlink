pragma solidity 0.7.6;

contract UpkeepPerformCounterRestrictive {
  event PerformingUpkeep(bool eligible, address from, uint256 initialCall, uint256 nextEligible, uint256 blockNumber);

  uint256 public initialCall = 0;
  uint256 public nextEligible = 0;
  uint256 public testRange;
  uint256 public averageEligibilityCadence;
  uint256 public checkGasToBurn;
  uint256 public performGasToBurn;
  mapping(bytes32 => bool) public dummyMap; // used to force storage lookup

  uint256 private count = 0;

  constructor(uint256 _testRange, uint256 _averageEligibilityCadence) {
    testRange = _testRange;
    averageEligibilityCadence = _averageEligibilityCadence;
  }

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    uint256 startGas = gasleft();
    uint256 blockNum = block.number - 1;
    bool dummy;
    // burn gas
    while (startGas - gasleft() < checkGasToBurn) {
      dummy = dummy && dummyMap[blockhash(blockNum)]; // arbitrary storage reads
      blockNum--;
    }
    return (eligible(), abi.encode(dummy));
  }

  function performUpkeep(bytes calldata) external {
    uint256 startGas = gasleft();
    bool eligible = eligible();
    uint256 blockNum = block.number;
    emit PerformingUpkeep(eligible, tx.origin, initialCall, nextEligible, blockNum);
    require(eligible);
    if (initialCall == 0) {
      initialCall = blockNum;
    }
    nextEligible = (blockNum + (rand() % (averageEligibilityCadence * 2))) + 1;
    count++;
    // burn gas
    blockNum--;
    while (startGas - gasleft() < performGasToBurn) {
      dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
      blockNum--;
    }
  }

  function setCheckGasToBurn(uint256 value) public {
    checkGasToBurn = value;
  }

  function setPerformGasToBurn(uint256 value) public {
    performGasToBurn = value;
  }

  function getCountPerforms() public view returns (uint256) {
    return count;
  }

  function eligible() internal view returns (bool) {
    return initialCall == 0 || (block.number - initialCall < testRange && block.number > nextEligible);
  }

  function checkEligible() public view returns (bool) {
    return eligible();
  }

  function reset() external {
    initialCall = 0;
    count = 0;
  }

  function setSpread(uint256 _newTestRange, uint256 _newAverageEligibilityCadence) external {
    testRange = _newTestRange;
    averageEligibilityCadence = _newAverageEligibilityCadence;
  }

  function rand() private view returns (uint256) {
    return uint256(keccak256(abi.encode(blockhash(block.number - 1), address(this))));
  }
}
