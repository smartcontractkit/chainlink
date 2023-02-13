// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;
import "../interfaces/KeeperCompatibleInterface.sol";
import "../vendor/GovernorBravoDelegate.sol";

/// @title Chainlink Keepers Compatible GovernorBravoDelegate Automator
contract GovernorBravoAutomator is KeeperCompatibleInterface {
  /// @notice Possible actions that can be taken in the performUpkeep function.
  /// QUEUE => calls 'queue(id)' on the governance contract
  /// EXECUTE => calls 'execute(id)' on the governance contract
  /// CANCEL => calls 'cancel(id)' on the governance contract
  /// UPDATE_INDEX => updates the starting proposal index within the
  ///                 upkeep contract to reduce the amount of proposals
  ///                 the need to be checked
  enum Action {
    QUEUE,
    EXECUTE,
    CANCEL,
    UPDATE_INDEX
  }

  GovernorBravoDelegate public immutable s_governanceContract;
  IGovernorBravoToken public immutable s_governanceTokenContract;
  uint256 public s_proposalStartingIndex;
  Action public action;

  constructor(
    GovernorBravoDelegate _governanceContract,
    uint256 _proposalStartingIndex,
    IGovernorBravoToken _tokenContract
  ) {
    s_governanceContract = _governanceContract;
    s_proposalStartingIndex = (_proposalStartingIndex > 0) ? _proposalStartingIndex : 1; // proposals start at 1.
    s_governanceTokenContract = _tokenContract;
  }

  ///@notice Simulated at each block by the Chainlink Keepers network. Checks if there are any actions (queue() or execute()) required on a governance contract. Also tracks a 'starting index'.
  ///@return upkeepNeeded return true if performUpkeep should be called
  ///@return performData bytes encoded: (governance action required, index of proposal)
  function checkUpkeep(
    bytes calldata /* checkData */
  ) external view override returns (bool upkeepNeeded, bytes memory performData) {
    //Get number of proposal
    uint256 proposalCount = s_governanceContract.proposalCount();

    //Find starting index
    uint256 newStartingIndex = findStartingIndex();

    //If new starting index found, update in performUpkeep
    if (newStartingIndex > s_proposalStartingIndex) {
      performData = abi.encode(Action.UPDATE_INDEX, newStartingIndex);
      return (true, performData);
    }

    //Go through each proposal and check the current state
    for (uint256 i = newStartingIndex; i <= proposalCount; i++) {
      GovernorBravoDelegateStorageV1.ProposalState state = s_governanceContract.state(i);
      (, address proposer, , , , , , , , ) = s_governanceContract.proposals(i);

      if (state == GovernorBravoDelegateStorageV1.ProposalState.Succeeded) {
        //If the state is 'Succeeded' then call 'queue' with the Proposal ID
        performData = abi.encode(Action.QUEUE, i);
        return (true, performData);
      } else if (state == GovernorBravoDelegateStorageV1.ProposalState.Queued) {
        //If the state is 'Queued' then call 'execute' with the Proposal ID
        performData = abi.encode(Action.EXECUTE, i);
        return (true, performData);
      } else if (
        s_governanceTokenContract.getPriorVotes(proposer, sub256(block.number, 1)) <
        s_governanceContract.proposalThreshold()
      ) {
        performData = abi.encode(Action.CANCEL, i);
        return (true, performData);
      }
    }

    revert("no action needed");
  }

  ///@notice Chainlink Keepers will execute when checkUpkeep returns 'true'. Decodes the 'performData' passed in from checkUpkeep and performs an action as needed.
  ///@param performData bytes encoded: (governance action required, index of proposal)
  ///@dev The governance contract has action validation built-in
  function performUpkeep(bytes calldata performData) external override {
    //Decode performData
    (Action performAction, uint256 proposalIndex) = abi.decode(performData, (Action, uint256));

    //Check state of proposal at index
    GovernorBravoDelegateStorageV1.ProposalState state = s_governanceContract.state(proposalIndex);

    //Revalidate state and action of provided index
    if (performAction == Action.QUEUE && state == GovernorBravoDelegateStorageV1.ProposalState.Succeeded) {
      s_governanceContract.queue(proposalIndex);
    } else if (performAction == Action.EXECUTE && state == GovernorBravoDelegateStorageV1.ProposalState.Queued) {
      s_governanceContract.execute(proposalIndex);
    } else if (performAction == Action.CANCEL && state != GovernorBravoDelegateStorageV1.ProposalState.Executed) {
      s_governanceContract.cancel(proposalIndex);
    } else if (performAction == Action.UPDATE_INDEX) {
      uint256 newStartingIndex = findStartingIndex();
      require(newStartingIndex > s_proposalStartingIndex, "No update required");
      s_proposalStartingIndex = newStartingIndex;
    }
  }

  ///@notice Goes through each proposal, if any proposal is in a state that needs to be checked OR its the last proposal, will set it as the starting proposal index for future checks
  ///@return index The proposal index to start checking from

  function findStartingIndex() public view returns (uint256 index) {
    // Set current starting index
    uint256 pendings_proposalStartIndex = s_proposalStartingIndex;
    // Get current proposal count
    uint256 proposalCount = s_governanceContract.proposalCount();

    for (uint256 i = pendings_proposalStartIndex; i <= proposalCount; i++) {
      GovernorBravoDelegateStorageV1.ProposalState state = s_governanceContract.state(i);
      if (
        state == GovernorBravoDelegateStorageV1.ProposalState.Pending ||
        state == GovernorBravoDelegateStorageV1.ProposalState.Active ||
        state == GovernorBravoDelegateStorageV1.ProposalState.Succeeded ||
        state == GovernorBravoDelegateStorageV1.ProposalState.Queued ||
        i == proposalCount
      ) {
        pendings_proposalStartIndex = i;
        break;
      }
    }
    return (pendings_proposalStartIndex);
  }

  function sub256(uint256 a, uint256 b) internal pure returns (uint256) {
    require(b <= a, "subtraction underflow");
    return a - b;
  }
}

interface IGovernorBravoToken {
  function getPriorVotes(address account, uint256 blockNumber) external view returns (uint96);
}
