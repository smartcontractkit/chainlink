// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/interfaces/IERC165.sol";
import {Common} from "../../../libraries/Common.sol";
import {IVerifierFeeManager} from "../../interfaces/IVerifierFeeManager.sol";

interface IFeeManager is IERC165, IVerifierFeeManager {
  struct Quote {
    address quoteAddress;
  }

  /**
   * @notice Processes the fee for a report, billing the subscriber and paying the reward manager
   * @param payload report and quote data to process the fee for
   * @param subscriber address of the user to process fee for
   */
  function processFee(bytes calldata payload, address subscriber) external payable;

  /**
   * @notice Processes the fees for each report in the payload, billing the subscriber and paying the reward manager
   * @param payloads reports and quotes to process
   * @param subscriber address of the user to process fee for
   */
  function processFeeBulk(bytes[] calldata payloads, address subscriber) external payable;

  /**
   * @notice Calculate the applied fee and the reward from a report. If the sender is a subscriber, they will receive a discount.
   * @param subscriber address trying to verify
   * @param report report to calculate the fee for
   * @param quote any metadata required to fetch the fee
   * @return (fee, reward) fee and the reward data
   */
  function getFeeAndReward(
    address subscriber,
    bytes memory report,
    Quote memory quote
  ) external returns (Common.Asset memory, Common.Asset memory);

  /**
   * @notice Sets the fee recipients within the reward manager
   * @param configDigest digest of the configuration
   * @param rewardRecipientAndWeights the address and weights of all the recipients to receive rewards
   */
  function setFeeRecipients(
    bytes32 configDigest,
    Common.AddressAndWeight[] calldata rewardRecipientAndWeights
  ) external;

  /**
   * @notice Sets the native surcharge
   * @param surcharge surcharge to be paid if paying in native
   */
  function setNativeSurcharge(uint64 surcharge) external;

  /**
   * @notice Adds a subscriber to the fee manager
   * @param subscriber address of the subscriber
   * @param feedId feed id to apply the discount to
   * @param token token to apply the discount to
   * @param discount discount to be applied to the fee
   */
  function updateSubscriberDiscount(address subscriber, bytes32 feedId, address token, uint64 discount) external;

  /**
   * @notice Withdraws any native or LINK rewards to the owner address
   * @param quantity quantity of tokens to withdraw, address(0) is native
   * @param quantity quantity to withdraw
   */
  function withdraw(address assetAddress, uint192 quantity) external;

  /**
   * @notice Returns the link balance of the fee manager
   * @return link balance of the fee manager
   */
  function linkAvailableForPayment() external returns (uint256);

  /**
   * @notice Admin function to pay the LINK deficit for a given config digest
   * @param configDigest the config digest to pay the deficit for
   */
  function payLinkDeficit(bytes32 configDigest) external;

  /**
   * @notice The structure to hold a fee and reward to verify a report
   * @param digest the digest linked to the fee and reward
   * @param fee the fee paid to verify the report
   * @param reward the reward paid upon verification
   */
  struct FeeAndReward {
    bytes32 configDigest;
    Common.Asset fee;
    Common.Asset reward;
  }
}
