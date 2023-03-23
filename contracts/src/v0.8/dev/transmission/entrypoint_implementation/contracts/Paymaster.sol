// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "../../../vendor/entrypoint/interfaces/IPaymaster.sol";
import "./SCALibrary.sol";
import "../../../vendor/entrypoint/core/Helpers.sol";
import "../../../../interfaces/LinkTokenInterface.sol";
import "./SCALibrary.sol";

/// @dev LINK token paymaster implementation.
contract Paymaster is IPaymaster {
  error OnlyCallableFromLink();
  error InvalidCalldata();

  LinkTokenInterface public immutable LINK;
  uint256 weiPerUnitLink = 5000000000000000;

  mapping(bytes32 => bool) userOpHashMapping;
  mapping(address => uint256) subscriptions;

  constructor(LinkTokenInterface linkToken) {
    LINK = linkToken;
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
      require(false, "already tried");
      // return ("", _packValidationData(true, 0, 0)); // already tried
    }

    uint256 extraCostJuels = extractExtraCostJuels(userOp);
    uint256 costJuels = (1e18 * maxCost) / weiPerUnitLink + extraCostJuels;
    if (subscriptions[userOp.sender] < costJuels) {
      require(false, "insufficient funds");
      // return ("", _packValidationData(true, 0, 0)); // insufficient funds
    }

    userOpHashMapping[userOpHash] = true;
    return (abi.encode(userOp.sender, extraCostJuels), _packValidationData(false, 0, 0)); // already tried
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
    uint256 costJuels = (1e18 * actualGasCost) / weiPerUnitLink;
    subscriptions[sender] -= (costJuels + extraCostJuels);
  }
}
