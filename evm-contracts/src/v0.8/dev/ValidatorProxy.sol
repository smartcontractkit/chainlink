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

  /// @notice Configuration of the current and proposed aggregator
  ProxyConfiguration private s_currentAggregator;
  address private s_proposedAggregator;

  /// @notice Configuration of the current and proposed validator
  ProxyConfiguration private s_currentValidator;
  address private s_proposedValidator;

  event NewAggregatorProposed(
    address indexed aggregator
  );
  event ProposedAggregatorRetracted(
    address indexed aggregator
  );
  event AggregatorUpgraded(
    address indexed previous,
    address indexed current
  );
  event NewValidatorProposed(
    address indexed validator
  );
  event ProposedValidatorRetracted(
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

  constructor(
    address aggregator,
    address validator
  )
    ConfirmedOwner(msg.sender)
  {
    s_currentAggregator.target = aggregator;
    s_currentValidator.target = validator;
  }

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
      return false;
    }

    // Send the validate call to the current validator
    ProxyConfiguration memory currentValidator = s_currentValidator;
    bool success = AggregatorValidatorInterface(currentValidator.target).validate(
      previousRoundId,
      previousAnswer,
      currentRoundId,
      currentAnswer
    );
    // If there is a new proposed validator, send the validate call to that validator also
    if (currentValidator.hasNewProposal) {
      bool proposedSuccess = AggregatorValidatorInterface(s_proposedValidator).validate(
        previousRoundId,
        previousAnswer,
        currentRoundId,
        currentAnswer
      );
      success = success && proposedSuccess;
    }
    return success;
  }

  /** AGGREGATOR CONFIGURATION FUNCTIONS **/

  function proposeNewAggregator(
    address proposed
  )
    external
    onlyOwner()
  {
    s_proposedAggregator = proposed;
    s_currentAggregator.hasNewProposal = true;
    emit NewAggregatorProposed(proposed);
  }

  function retractProposedAggregator()
    external
    onlyOwner()
  {
    address proposed = s_proposedAggregator;
    s_proposedAggregator = address(0);
    s_currentAggregator.hasNewProposal = false;
    emit ProposedAggregatorRetracted(proposed);
  }

  function upgradeAggregator()
    external
    onlyOwner()
  {
    // Get configuration in memory
    ProxyConfiguration memory current = s_currentAggregator;
    address previous = current.target;
    address proposed = s_proposedAggregator;

    // Perform the upgrade
    require(current.hasNewProposal == true && proposed != address(0), "No proposal");
    current.target = proposed;
    current.hasNewProposal = false;

    s_proposedAggregator = address(0);

    emit AggregatorUpgraded(previous, proposed);
  }

  /** VALIDATOR CONFIGURATION FUNCTIONS **/

  function proposeNewValidator(
    address proposed
  )
    external
    onlyOwner()
  {
    s_proposedValidator = proposed;
    s_currentValidator.hasNewProposal = true;
    emit NewValidatorProposed(proposed);
  }

  function retractProposedValidator()
    external
    onlyOwner()
  {
    address proposed = s_proposedValidator;
    s_proposedValidator = address(0);
    s_currentValidator.hasNewProposal = false;
    emit ProposedValidatorRetracted(proposed);
  }

  function upgradeValidator()
    external
    onlyOwner()
  {
    // Get configuration in memory
    ProxyConfiguration memory current = s_currentValidator;
    address previous = current.target;
    address proposed = s_proposedValidator;

    // Perform the upgrade
    require(current.hasNewProposal == true && proposed != address(0), "No proposal");
    current.target = proposed;
    current.hasNewProposal = false;

    s_proposedValidator = address(0);

    emit ValidatorUpgraded(previous, proposed);
  }

}