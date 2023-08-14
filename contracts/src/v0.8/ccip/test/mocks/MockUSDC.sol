// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ITokenMessenger} from "../../pools/USDC/ITokenMessenger.sol";
import {IMessageReceiver} from "../../pools/USDC/IMessageReceiver.sol";

// This contract mocks both the ITokenMessenger and IMessageReceiver
// contracts involved with the Cross Chain Token Protocol.
contract MockUSDC is ITokenMessenger, IMessageReceiver {
  uint32 private immutable i_messageBodyVersion;
  bytes32 public constant i_destinationTokenMessenger = keccak256("i_destinationTokenMessenger");

  // Indicated whether the receiveMessage() call should succeed.
  bool public s_shouldSucceed;
  uint64 public s_nonce;

  constructor(uint32 version) {
    i_messageBodyVersion = version;
    s_nonce = 1;
    s_shouldSucceed = true;
  }

  function depositForBurnWithCaller(
    uint256 amount,
    uint32 destinationDomain,
    bytes32 mintRecipient,
    address burnToken,
    bytes32 destinationCaller
  ) external returns (uint64) {
    emit DepositForBurn(
      s_nonce,
      burnToken,
      amount,
      msg.sender,
      mintRecipient,
      destinationDomain,
      i_destinationTokenMessenger,
      destinationCaller
    );
    return s_nonce++;
  }

  function receiveMessage(bytes calldata, bytes calldata) external view returns (bool success) {
    return s_shouldSucceed;
  }

  function messageBodyVersion() external view returns (uint32) {
    return i_messageBodyVersion;
  }

  function setShouldSucceed(bool shouldSucceed) external {
    s_shouldSucceed = shouldSucceed;
  }
}
