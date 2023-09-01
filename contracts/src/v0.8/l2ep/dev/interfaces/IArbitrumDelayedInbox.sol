// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../../vendor/arb-bridge-eth/v0.8.0-custom/contracts/bridge/interfaces/IInbox.sol";

/**
 * @notice This interface extends Arbitrum's IInbox interface to include
 * the calculateRetryableSubmissionFee.  This new function was added as part
 * of Arbitrum's Nitro migration but was excluded from the IInbox interface.  This setup
 * works for us as the team has added it as a public function to the IInbox proxy
 * contract's implementation
 */
interface IArbitrumDelayedInbox is IInbox {
  function calculateRetryableSubmissionFee(uint256 dataLength, uint256 baseFee) external view returns (uint256);
}
