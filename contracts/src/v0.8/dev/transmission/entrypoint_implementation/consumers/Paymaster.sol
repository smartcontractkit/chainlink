// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.12;

import "../interfaces/IPaymaster.sol";
import "./SCALibrary.sol";
import "../contracts/Helpers.sol";

/**
 * the interface exposed by a paymaster contract, who agrees to pay the gas for user's operations.
 * a paymaster must hold a stake to cover the required entrypoint stake and also the gas for the transaction.
 */
contract Paymaster is IPaymaster {
  mapping(bytes32 => bool) userOpHashMapping;
  mapping(address => uint256) subscriptions;
  uint256 weiPerUnitLink = 5000000000000000;

  function validatePaymasterUserOp(
    UserOperation calldata userOp,
    bytes32 userOpHash,
    uint256 maxCost
  ) external returns (bytes memory context, uint256 validationData) {
    if (userOpHashMapping[userOpHash]) {
      return ("", _packValidationData(true, 0, 0)); // already tried
    }

    uint256 costJuels = (1e18 * maxCost) / weiPerUnitLink;
    if (subscriptions[userOp.sender] < costJuels) {
      return ("", _packValidationData(true, 0, 0)); // insufficient funds
    }

    userOpHashMapping[userOpHash] = true;
    return (abi.encode(userOp.sender), _packValidationData(false, 0, 0)); // already tried
  }

  function postOp(PostOpMode mode, bytes calldata context, uint256 actualGasCost) external {
    address sender = abi.decode(context, (address));
    uint256 costJuels = (1e18 * actualGasCost) / weiPerUnitLink;
    subscriptions[sender] -= costJuels;
  }

  function deposit(uint256 amount, address subscription) external {
    subscriptions[subscription] += amount;
  }
}
