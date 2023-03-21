// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.12;

import "../interfaces/IPaymaster.sol";
import "./SCALibrary.sol";
import "../contracts/Helpers.sol";
import "../../../../interfaces/LinkTokenInterface.sol";
import "./SCALibrary.sol";

/**
 * the interface exposed by a paymaster contract, who agrees to pay the gas for user's operations.
 * a paymaster must hold a stake to cover the required entrypoint stake and also the gas for the transaction.
 */
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
      return ("", _packValidationData(true, 0, 0)); // already tried
    }

    uint256 costJuels = (1e18 * maxCost) / weiPerUnitLink + extractExtraCostJuels(userOp);
    if (subscriptions[userOp.sender] < costJuels) {
      return ("", _packValidationData(true, 0, 0)); // insufficient funds
    }

    userOpHashMapping[userOpHash] = true;
    return (abi.encode(userOp.sender), _packValidationData(false, 0, 0)); // already tried
  }

  function extractExtraCostJuels(UserOperation calldata userOp) internal returns (uint256 extraCost) {
    if (userOp.paymasterAndData.length == 20) {
      return 0; // no extra data.
    }

    uint8 paymentType = uint8(userOp.paymasterAndData[20]);
    if (paymentType == uint8(SCALibrary.LinkPaymentType.DIRECT_FUNDING)) {
      SCALibrary.DirectFundingData memory directFundingData = abi.decode(userOp.paymasterAndData[21:], (SCALibrary.DirectFundingData));
      if (
        directFundingData.topupThreshold != 0 &&
        LINK.balanceOf(directFundingData.recipient) < directFundingData.topupThreshold
      ) {
        LINK.transfer(directFundingData.recipient, directFundingData.topupAmount);
        extraCost = directFundingData.topupAmount;
      }
    }
  }

  function postOp(PostOpMode mode, bytes calldata context, uint256 actualGasCost) external {
    address sender = abi.decode(context, (address));
    uint256 costJuels = (1e18 * actualGasCost) / weiPerUnitLink;
    subscriptions[sender] -= costJuels;
  }
}
