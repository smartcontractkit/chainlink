// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./ConfirmedOwner.sol";
import "../interfaces/AggregatorValidatorInterface.sol";

contract ValidatorProxy is AggregatorValidatorInterface, ConfirmedOwner {

  /// @notice Uses a single storage slot to store the current address
  struct AggregatorConfiguration {
    address target;
    bool hasNewProposal;
  }

  struct ValidatorConfiguration {
    AggregatorValidatorInterface target;
    bool hasNewProposal;
  }

  // Configuration for the current aggregator
  AggregatorConfiguration private s_currentAggregator;
  // Proposed aggregator address
  address private s_proposedAggregator;

  // Configuration for the current validator
  ValidatorConfiguration private s_currentValidator;
  // Proposed validator address
  AggregatorValidatorInterface private s_proposedValidator;

  event AggregatorProposed(
    address indexed aggregator
  );
  event AggregatorUpgraded(
    address indexed previous,
    address indexed current
  );
  event ValidatorProposed(
    AggregatorValidatorInterface indexed validator
  );
  event ValidatorUpgraded(
    AggregatorValidatorInterface indexed previous,
    AggregatorValidatorInterface indexed current
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
    AggregatorValidatorInterface validator
  )
    ConfirmedOwner(msg.sender)
  {
    s_currentAggregator = AggregatorConfiguration({
      target: aggregator,
      hasNewProposal: false
    });
    s_currentValidator = ValidatorConfiguration({
      target: validator,
      hasNewProposal: false
    });
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
    ValidatorConfiguration memory currentValidator = s_currentValidator;
    require(address(s_currentValidator.target) != address(0), "No validator set");
    currentValidator.target.validate(
      previousRoundId,
      previousAnswer,
      currentRoundId,
      currentAnswer
    );
    // If there is a new proposed validator, send the validate call to that validator also
    if (currentValidator.hasNewProposal) {
      s_proposedValidator.validate(
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
    require(s_proposedAggregator != proposed, "No change");
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
    AggregatorConfiguration memory current = s_currentAggregator;
    address previous = current.target;
    address proposed = s_proposedAggregator;

    // Perform the upgrade
    require(current.hasNewProposal == true, "No proposal");
    s_currentAggregator = AggregatorConfiguration({
      target: proposed,
      hasNewProposal: false
    });
    delete s_proposedAggregator;

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
    AggregatorValidatorInterface proposed
  )
    external
    onlyOwner()
  {
    require(s_proposedValidator != proposed, "No change");
    s_proposedValidator = proposed;
    // If proposed is zero address, hasNewProposal = false
    s_currentValidator.hasNewProposal = (address(proposed) != address(0));
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
    ValidatorConfiguration memory current = s_currentValidator;
    AggregatorValidatorInterface previous = current.target;
    AggregatorValidatorInterface proposed = s_proposedValidator;

    // Perform the upgrade
    require(current.hasNewProposal == true, "No proposal");
    s_currentValidator = ValidatorConfiguration({
      target: proposed,
      hasNewProposal: false
    });
    delete s_proposedValidator;

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
      AggregatorValidatorInterface current,
      bool hasProposal,
      AggregatorValidatorInterface proposed
    )
  {
    current = s_currentValidator.target;
    hasProposal = s_currentValidator.hasNewProposal;
    proposed = s_proposedValidator;
  }

}