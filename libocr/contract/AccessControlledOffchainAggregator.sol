// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;

import "./OffchainAggregator.sol";
import "./SimpleReadAccessController.sol";

/**
 * @notice Wrapper of OffchainAggregator which checks read access on Aggregator-interface methods
 */
contract AccessControlledOffchainAggregator is OffchainAggregator, SimpleReadAccessController {

  constructor(
    uint32 _maximumGasPrice,
    uint32 _reasonableGasPrice,
    uint32 _microLinkPerEth,
    uint32 _linkGweiPerObservation,
    uint32 _linkGweiPerTransmission,
    LinkTokenInterface _link,
    int192 _minAnswer,
    int192 _maxAnswer,
    AccessControllerInterface _billingAccessController,
    AccessControllerInterface _requesterAccessController,
    uint8 _decimals,
    string memory description
  )
    OffchainAggregator(
      _maximumGasPrice,
      _reasonableGasPrice,
      _microLinkPerEth,
      _linkGweiPerObservation,
      _linkGweiPerTransmission,
      _link,
      _minAnswer,
      _maxAnswer,
      _billingAccessController,
      _requesterAccessController,
      _decimals,
      description
    ) {
    }

  /*
   * Versioning
   */

  function typeAndVersion()
    external
    override
    pure
    virtual
    returns (string memory)
  {
    return "AccessControlledOffchainAggregator 4.0.0";
  }


  /*
   * v2 Aggregator interface
   */

  /// @inheritdoc OffchainAggregator
  function latestAnswer()
    public
    override
    view
    checkAccess()
    returns (int256)
  {
    return super.latestAnswer();
  }

  /// @inheritdoc OffchainAggregator
  function latestTimestamp()
    public
    override
    view
    checkAccess()
    returns (uint256)
  {
    return super.latestTimestamp();
  }

  /// @inheritdoc OffchainAggregator
  function latestRound()
    public
    override
    view
    checkAccess()
    returns (uint256)
  {
    return super.latestRound();
  }

  /// @inheritdoc OffchainAggregator
  function getAnswer(uint256 _roundId)
    public
    override
    view
    checkAccess()
    returns (int256)
  {
    return super.getAnswer(_roundId);
  }

  /// @inheritdoc OffchainAggregator
  function getTimestamp(uint256 _roundId)
    public
    override
    view
    checkAccess()
    returns (uint256)
  {
    return super.getTimestamp(_roundId);
  }

  /*
   * v3 Aggregator interface
   */

  /// @inheritdoc OffchainAggregator
  function description()
    public
    override
    view
    checkAccess()
    returns (string memory)
  {
    return super.description();
  }

  /// @inheritdoc OffchainAggregator
  function getRoundData(uint80 _roundId)
    public
    override
    view
    checkAccess()
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    return super.getRoundData(_roundId);
  }

  /// @inheritdoc OffchainAggregator
  function latestRoundData()
    public
    override
    view
    checkAccess()
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    return super.latestRoundData();
  }

}
