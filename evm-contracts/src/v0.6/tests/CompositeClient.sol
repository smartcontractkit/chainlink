// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "../ChainlinkClient.sol";
import "../VRFConsumerBase.sol";

/**
 * @notice CompositeClient to test that both ChainlinkClient and 
 * VRFConsumerBase can be inherited by the same contract 
 */
contract CompositeClient is ChainlinkClient, VRFConsumerBase {
  
  constructor(address _vrfCoordinator, address _link) public
    // solhint-disable-next-line no-empty-blocks
    VRFConsumerBase(_vrfCoordinator, _link) { /* empty */ }
  
  function fulfillRandomness(bytes32, uint256)
    internal override
  {
    return;
  }
}