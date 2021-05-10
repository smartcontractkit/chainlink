// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./ConfirmedOwner.sol";
import "../interfaces/AggregatorValidatorInterface.sol";

contract ValidatorProxy is AggregatorValidatorInterface, ConfirmedOwner {

  /// @notice Uses a single storage slot to store the current address
  struct ProxyConfiguration {
    address target;
    bool hasNewProposal;
  }

  // Configuration for the current aggregator
  ProxyConfiguration private s_currentAggregator;
  // Proposed aggregator address
  address private s_proposedAggregator;

  // Configuration for the current validator
  ProxyConfiguration private s_currentValidator;
  // Proposed validator address
  address private s_proposedValidator;

  event AggregatorProposed(
    address indexed aggregator
  );
  event AggregatorUpgraded(
    address indexed previous,
    address indexed current
  );
  event ValidatorProposed(
    address indexed validator
  );
  event ValidatorUpgraded(
    address indexed previous,
    address indexed current
  );
  /// @notice The proposed aggregator called validate, but the call was not passed on to any validators
  event ProposedAggregatorValidateCall(
    address indexed proposed,
    uint256 previousRoundId,
    int256 previousAnswer,
    uint256 currentRoundId,
    int256 currentAnswer
  );

  /**
   * @notice Construct the ValidatorProxy with an aggregator and a validator
   * @param aggregator address
   * @param validator address
   */
  constructor(
    address aggregator,
    address validator
  )
    ConfirmedOwner(msg.sender)
  {
    s_currentAggregator.target = aggregator;
    s_currentValidator.target = validator;
  }

  /**
   * @notice Validate a transmission
   * @dev Must be called by either the `s_currentAggregator.target`, or the `s_proposedAggregator`.
   * If called by the `s_currentAggregator.target` this function passes the call on to the `s_currentValidator.target`
   * and the `s_proposedValidator`, if it is set.
   * If called by the `s_proposedAggregator` this function emits a `ProposedAggregatorValidateCall` to signal that
   * the call was received.
   * @param previousRoundId uint256
   * @param previousAnswer int256
   * @param currentRoundId uint256
   * @param currentAnswer int256
   * @return bool
   */
  function validate(
    uint256 previousRoundId,
    int256 previousAnswer,
    uint256 currentRoundId,
    int256 currentAnswer
  )
    external
    override
    returns (
      bool
    )
  {
    address currentAggregator = s_currentAggregator.target;
    address proposedAggregator = s_proposedAggregator;
    require(msg.sender == currentAggregator || msg.sender == proposedAggregator, "Not a configured aggregator");
    // If the aggregator is still in proposed state, emit an event and don't push to any validator.
    // This is to confirm that `validate` is being called prior to upgrade.
    if (msg.sender == proposedAggregator) {
      emit ProposedAggregatorValidateCall(
        proposedAggregator,
        previousRoundId,
        previousAnswer,
        currentRoundId,
        currentAnswer
      );
      return true;
    }

    // Send the validate call to the current validator
    ProxyConfiguration memory currentValidator = s_currentValidator;
    require(s_currentValidator.target != address(0), "No validator set");
    AggregatorValidatorInterface(currentValidator.target).validate(
      previousRoundId,
      previousAnswer,
      currentRoundId,
      currentAnswer
    );
    // If there is a new proposed validator, send the validate call to that validator also
    if (currentValidator.hasNewProposal) {
      AggregatorValidatorInterface(s_proposedValidator).validate(
        previousRoundId,
        previousAnswer,
        currentRoundId,
        currentAnswer
      );
    }
    return true;
  }

  /** AGGREGATOR CONFIGURATION FUNCTIONS **/

  /**
   * @notice Propose an aggregator
   * @dev A zero address can be used to unset the proposed aggregator. Only owner can call.
   * @param proposed address
   */
  function proposeNewAggregator(
    address proposed
  )
    external
    onlyOwner()
  {
    s_proposedAggregator = proposed;
    // If proposed is zero address, hasNewProposal = false
    s_currentAggregator.hasNewProposal = (proposed != address(0));
    emit AggregatorProposed(proposed);
  }

  /**
   * @notice Upgrade the aggregator by setting the current aggregator as the proposed aggregator.
   * @dev Must have a proposed aggregator. Only owner can call.
   */
  function upgradeAggregator()
    external
    onlyOwner()
  {
    // Get configuration in memory
    ProxyConfiguration memory current = s_currentAggregator;
    address previous = current.target;
    address proposed = s_proposedAggregator;

    // Perform the upgrade
    require(current.hasNewProposal == true, "No proposal");
    current.target = proposed;
    current.hasNewProposal = false;

    s_currentAggregator = current;
    s_proposedAggregator = address(0);

    emit AggregatorUpgraded(previous, proposed);
  }

  /**
   * @notice Get aggregator details
   * @return current address
   * @return hasProposal bool
   * @return proposed address
   */
  function getAggregators()
    external
    view
    returns(
      address current,
      bool hasProposal,
      address proposed
    )
  {
    current = s_currentAggregator.target;
    hasProposal = s_currentAggregator.hasNewProposal;
    proposed = s_proposedAggregator;
  }

  /** VALIDATOR CONFIGURATION FUNCTIONS **/

  /**
   * @notice Propose an validator
   * @dev A zero address can be used to unset the proposed validator. Only owner can call.
   * @param proposed address
   */
  function proposeNewValidator(
    address proposed
  )
    external
    onlyOwner()
  {
    s_proposedValidator = proposed;
    // If proposed is zero address, hasNewProposal = false
    s_currentValidator.hasNewProposal = (proposed != address(0));
    emit ValidatorProposed(proposed);
  }

  /**
   * @notice Upgrade the validator by setting the current validator as the proposed validator.
   * @dev Must have a proposed validator. Only owner can call.
   */
  function upgradeValidator()
    external
    onlyOwner()
  {
    // Get configuration in memory
    ProxyConfiguration memory current = s_currentValidator;
    address previous = current.target;
    address proposed = s_proposedValidator;

    // Perform the upgrade
    require(current.hasNewProposal == true, "No proposal");
    current.target = proposed;
    current.hasNewProposal = false;

    s_currentValidator = current;
    s_proposedValidator = address(0);

    emit ValidatorUpgraded(previous, proposed);
  }

  /**
   * @notice Get validator details
   * @return current address
   * @return hasProposal bool
   * @return proposed address
   */
  function getValidators()
    external
    view
    returns(
      address current,
      bool hasProposal,
      address proposed
    )
  {
    current = s_currentValidator.target;
    hasProposal = s_currentValidator.hasNewProposal;
    proposed = s_proposedValidator;
  }

}