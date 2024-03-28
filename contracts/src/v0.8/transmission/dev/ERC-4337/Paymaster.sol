// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IPaymaster} from "../../../vendor/entrypoint/interfaces/IPaymaster.sol";
import {SCALibrary} from "./SCALibrary.sol";
import {LinkTokenInterface} from "../../../shared/interfaces/LinkTokenInterface.sol";
import {AggregatorV3Interface} from "../../../shared/interfaces/AggregatorV3Interface.sol";
import {ConfirmedOwner} from "../../../shared/access/ConfirmedOwner.sol";
import {UserOperation} from "../../../vendor/entrypoint/interfaces/UserOperation.sol";
import {_packValidationData} from "../../../vendor/entrypoint/core/Helpers.sol";

/// @dev LINK token paymaster implementation.
/// TODO: more documentation.
contract Paymaster is IPaymaster, ConfirmedOwner {
  error OnlyCallableFromLink();
  error InvalidCalldata();
  error Unauthorized(address sender, address validator);
  error UserOperationAlreadyTried(bytes32 userOpHash);
  error InsufficientFunds(uint256 juelsNeeded, uint256 subscriptionBalance);

  LinkTokenInterface public immutable i_linkToken;
  AggregatorV3Interface public immutable i_linkEthFeed;
  address public immutable i_entryPoint;

  struct Config {
    uint32 stalenessSeconds;
    int256 fallbackWeiPerUnitLink;
  }
  Config public s_config;

  mapping(bytes32 => bool) internal s_userOpHashMapping;
  mapping(address => uint256) internal s_subscriptions;

  constructor(
    LinkTokenInterface linkToken,
    AggregatorV3Interface linkEthFeed,
    address entryPoint
  ) ConfirmedOwner(msg.sender) {
    i_linkToken = linkToken;
    i_linkEthFeed = linkEthFeed;
    i_entryPoint = entryPoint;
  }

  function setConfig(uint32 stalenessSeconds, int256 fallbackWeiPerUnitLink) external onlyOwner {
    s_config = Config({stalenessSeconds: stalenessSeconds, fallbackWeiPerUnitLink: fallbackWeiPerUnitLink});
  }

  function onTokenTransfer(address /* _sender */, uint256 _amount, bytes calldata _data) external {
    if (msg.sender != address(i_linkToken)) {
      revert OnlyCallableFromLink();
    }
    if (_data.length != 32) {
      revert InvalidCalldata();
    }

    address subscription = abi.decode(_data, (address));
    s_subscriptions[subscription] += _amount;
  }

  function validatePaymasterUserOp(
    UserOperation calldata userOp,
    bytes32 userOpHash,
    uint256 maxCost
  ) external returns (bytes memory context, uint256 validationData) {
    if (msg.sender != i_entryPoint) {
      revert Unauthorized(msg.sender, i_entryPoint);
    }
    if (s_userOpHashMapping[userOpHash]) {
      revert UserOperationAlreadyTried(userOpHash);
    }

    uint256 extraCostJuels = _handleExtraCostJuels(userOp);
    uint256 costJuels = _getCostJuels(maxCost) + extraCostJuels;
    if (s_subscriptions[userOp.sender] < costJuels) {
      revert InsufficientFunds(costJuels, s_subscriptions[userOp.sender]);
    }

    s_userOpHashMapping[userOpHash] = true;
    return (abi.encode(userOp.sender, extraCostJuels), _packValidationData(false, 0, 0)); // success
  }

  /// @dev Calculates any extra LINK cost for the user operation, based on the funding type passed to the
  /// @dev paymaster. Handles funding the LINK token funding described in the user operation.
  /// TODO: add logic for subscription top-up.
  function _handleExtraCostJuels(UserOperation calldata userOp) internal returns (uint256 extraCost) {
    if (userOp.paymasterAndData.length == 20) {
      return 0; // no extra data, stop here
    }

    uint8 paymentType = uint8(userOp.paymasterAndData[20]);

    // For direct funding, use top-up logic.
    if (paymentType == uint8(SCALibrary.LinkPaymentType.DIRECT_FUNDING)) {
      SCALibrary.DirectFundingData memory directFundingData = abi.decode(
        userOp.paymasterAndData[21:],
        (SCALibrary.DirectFundingData)
      );
      if (
        directFundingData.topupThreshold != 0 &&
        i_linkToken.balanceOf(directFundingData.recipient) < directFundingData.topupThreshold
      ) {
        i_linkToken.transfer(directFundingData.recipient, directFundingData.topupAmount);
        extraCost = directFundingData.topupAmount;
      }
    }
    return extraCost;
  }

  /// @dev Deducts user subscription balance after execution.
  function postOp(PostOpMode /* mode */, bytes calldata context, uint256 actualGasCost) external {
    if (msg.sender != i_entryPoint) {
      revert Unauthorized(msg.sender, i_entryPoint);
    }
    (address sender, uint256 extraCostJuels) = abi.decode(context, (address, uint256));
    s_subscriptions[sender] -= (_getCostJuels(actualGasCost) + extraCostJuels);
  }

  function _getCostJuels(uint256 costWei) internal view returns (uint256 costJuels) {
    costJuels = (1e18 * costWei) / uint256(_getFeedData());
    return costJuels;
  }

  function _getFeedData() internal view returns (int256) {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    uint256 timestamp;
    int256 weiPerUnitLink;
    (, weiPerUnitLink, , timestamp, ) = i_linkEthFeed.latestRoundData();
    if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
      weiPerUnitLink = s_config.fallbackWeiPerUnitLink;
    }
    return weiPerUnitLink;
  }
}
