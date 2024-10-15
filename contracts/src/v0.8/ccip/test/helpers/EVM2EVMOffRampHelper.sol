// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

contract EVM2EVMOffRampHelper {
  uint64 public s_nonce;
  mapping(address sender => uint64 nonce) public s_nonces;

  function execute(
    address[] memory senders
  ) external {
    for (uint256 i; i < senders.length; i++) {
      s_nonces[senders[i]]++;
    }
  }

  function metadataHash() external pure returns (bytes32) {
    return 0x0;
  }

  function getSenderNonce(
    address sender
  ) external view returns (uint64 nonce) {
    return s_nonces[sender];
  }
}
