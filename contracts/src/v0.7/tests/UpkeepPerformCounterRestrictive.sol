pragma solidity 0.7.6;

contract UpkeepPerformCounterRestrictive {
  event PerformingUpkeep(bool eligible, address from, uint256 initialCall, uint256 nextEligible, uint256 blockNumber);

  uint256 public initialCall = 0;
  uint256 public nextEligible = 0;
  uint256 public testRange;
  uint256 public averageEligibilityCadence;
  uint256 count = 0;

  constructor(uint256 _testRange, uint256 _averageEligibilityCadence) {
    testRange = _testRange;
    averageEligibilityCadence = _averageEligibilityCadence;
  }

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    return (eligible(), bytes(""));
  }

  function performUpkeep(bytes calldata data) external {
    bool eligible = eligible();
    uint256 blockNum = block.number;
    emit PerformingUpkeep(eligible, tx.origin, initialCall, nextEligible, blockNum);
    require(eligible);
    if (initialCall == 0) {
      initialCall = blockNum;
    }
    nextEligible = (blockNum + (rand() % (averageEligibilityCadence * 2))) + 1;
    count++;
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
