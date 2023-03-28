// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "../../../vendor/entrypoint/interfaces/IPaymaster.sol";
import "./SCALibrary.sol";
import "../../../vendor/entrypoint/core/Helpers.sol";
import "../../../../interfaces/LinkTokenInterface.sol";
import "../../../../interfaces/AggregatorV3Interface.sol";
import "./SCALibrary.sol";

/// @dev LINK token paymaster implementation.
/// TODO: more documentation.
contract Paymaster is IPaymaster {
  error OnlyCallableFromLink();
  error InvalidCalldata();
  error UserOperationAlreadyTried(bytes32 userOpHash);
  error InsufficientFunds(uint256 juelsNeeded, uint256 subscriptionBalance);

  LinkTokenInterface public immutable LINK;
  AggregatorV3Interface public immutable LINK_ETH_FEED;

  struct Config {
    uint32 stalenessSeconds;
    int256 fallbackWeiPerUnitLink;
  }
  Config public s_config;

  mapping(bytes32 => bool) userOpHashMapping;
  mapping(address => uint256) subscriptions;

  constructor(LinkTokenInterface linkToken, AggregatorV3Interface linkEthFeed) {
    LINK = linkToken;
    LINK_ETH_FEED = linkEthFeed;
  }

  function setConfig(uint32 stalenessSeconds, int256 fallbackWeiPerUnitLink) external {
    s_config = Config({stalenessSeconds: stalenessSeconds, fallbackWeiPerUnitLink: fallbackWeiPerUnitLink});
  }

  function onTokenTransfer(address /* _sender */, uint256 _amount, bytes calldata _data) external {
    if (msg.sender != address(LINK)) {
      revert OnlyCallableFromLink();
    }
    if (_data.length != 32) {
      revert InvalidCalldata();
    }

    address subscription = abi.decode(_data, (address));
    subscriptions[subscription] += _amount;
  }

  function validatePaymasterUserOp(
    UserOperation calldata userOp,
    bytes32 userOpHash,
    uint256 maxCost
  ) external returns (bytes memory context, uint256 validationData) {
    if (userOpHashMapping[userOpHash]) {
      revert UserOperationAlreadyTried(userOpHash);
    }

    uint256 extraCostJuels = extractExtraCostJuels(userOp);
    uint256 costJuels = getCostJuels(maxCost) + extraCostJuels;
    if (subscriptions[userOp.sender] < costJuels) {
      revert InsufficientFunds(costJuels, subscriptions[userOp.sender]);
    }

    userOpHashMapping[userOpHash] = true;
    return (abi.encode(userOp.sender, extraCostJuels), _packValidationData(false, 0, 0)); // success
  }

  /// @dev Calculates any extra LINK cost for the user operation, based on the funding type passed to the
  /// @dev paymaster.
  function extractExtraCostJuels(UserOperation calldata userOp) internal returns (uint256 extraCost) {
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
        LINK.balanceOf(directFundingData.recipient) < directFundingData.topupThreshold
      ) {
        LINK.transfer(directFundingData.recipient, directFundingData.topupAmount);
        extraCost = directFundingData.topupAmount;
      }
    }
  }

  /// @dev Deducts user subscription balance after execution.
  function postOp(PostOpMode /* mode */, bytes calldata context, uint256 actualGasCost) external {
    (address sender, uint256 extraCostJuels) = abi.decode(context, (address, uint256));
    subscriptions[sender] -= (getCostJuels(actualGasCost) + extraCostJuels);
  }

  function getCostJuels(uint256 costWei) internal view returns (uint256 costJuels) {
    uint256 costJuels = (1e18 * costWei) / uint256(getFeedData());
  }

  function getFeedData() internal view returns (int256) {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    uint256 timestamp;
    int256 weiPerUnitLink;
    (, weiPerUnitLink, , timestamp, ) = LINK_ETH_FEED.latestRoundData();
    if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
      weiPerUnitLink = s_config.fallbackWeiPerUnitLink;
    }
    return weiPerUnitLink;
  }
}
