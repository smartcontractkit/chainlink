// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./OCR2Aggregator.sol";
import "./SimpleReadAccessController.sol";

/**
 * @notice Wrapper of OCR2Aggregator which checks read access on Aggregator-interface methods
 */
contract AccessControlledOCR2Aggregator is OCR2Aggregator, SimpleReadAccessController {

  constructor(
    LinkTokenInterface _link,
    int192 _minAnswer,
    int192 _maxAnswer,
    AccessControllerInterface _billingAccessController,
    AccessControllerInterface _requesterAccessController,
    uint8 _decimals,
    string memory description
  )
    OCR2Aggregator(
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
    return "AccessControlledOCR2Aggregator 1.0.0-alpha";
  }


  /*
   * v2 Aggregator interface
   */

  /// @inheritdoc OCR2Aggregator
  function latestAnswer()
    public
    override
    view
    checkAccess()
    returns (int256)
  {
    return super.latestAnswer();
  }

  /// @inheritdoc OCR2Aggregator
  function latestTimestamp()
    public
    override
    view
    checkAccess()
    returns (uint256)
  {
    return super.latestTimestamp();
  }

  /// @inheritdoc OCR2Aggregator
  function latestRound()
    public
    override
    view
    checkAccess()
    returns (uint256)
  {
    return super.latestRound();
  }

  /// @inheritdoc OCR2Aggregator
  function getAnswer(uint256 _roundId)
    public
    override
    view
    checkAccess()
    returns (int256)
  {
    return super.getAnswer(_roundId);
  }

  /// @inheritdoc OCR2Aggregator
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

  /// @inheritdoc OCR2Aggregator
  function description()
    public
    override
    view
    checkAccess()
    returns (string memory)
  {
    return super.description();
  }

  /// @inheritdoc OCR2Aggregator
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

  /// @inheritdoc OCR2Aggregator
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
