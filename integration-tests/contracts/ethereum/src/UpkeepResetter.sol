pragma solidity 0.7.6;

contract UpkeepResetter {
  function _ResetConsumerBenchmark(
    address upkeepAddress,
    uint256 testRange,
    uint256 averageEligibilityCadence,
    uint256 firstEligibleBuffer,
    uint256 checkGasToBurn,
    uint256 performGasToBurn
  ) private {
    KeeperConsumerBenchmark consumer = KeeperConsumerBenchmark(upkeepAddress);
    consumer.setFirstEligibleBuffer(firstEligibleBuffer);
    consumer.setSpread(testRange, averageEligibilityCadence);
    consumer.setCheckGasToBurn(checkGasToBurn);
    consumer.setPerformGasToBurn(performGasToBurn);
    consumer.reset();
  }

  function ResetManyConsumerBenchmark(
    address[] memory upkeepAddresses,
    uint256 testRange,
    uint256 averageEligibilityCadence,
    uint256 firstEligibleBuffer,
    uint256 checkGasToBurn,
    uint256 performGasToBurn
  ) external {
    for (uint256 i = 0; i < upkeepAddresses.length; i++) {
      _ResetConsumerBenchmark(
        upkeepAddresses[i],
        testRange,
        averageEligibilityCadence,
        firstEligibleBuffer,
        checkGasToBurn,
        performGasToBurn
      );
    }
  }
}

interface KeeperConsumerBenchmark {
  function reset() external;

  function setSpread(uint256 _newTestRange, uint256 _newAverageEligibilityCadence) external;

  function setFirstEligibleBuffer(uint256 _firstEligibleBuffer) external;

  function setCheckGasToBurn(uint256 value) external;

  function setPerformGasToBurn(uint256 value) external;
}
