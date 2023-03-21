// SPDX-License-Identifier: MIT
pragma solidity >=0.4.22 <0.9.0;

library SCALibrary {
  // keccak256("EIP712Domain(uint256 chainId, address verifyingContract)");
  bytes32 internal constant DOMAIN_SEPARATOR = hex"1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea61";

  // keccak256("executeTransactionFromEntryPoint(address to, uint256 value, bytes calldata data)");
  bytes32 internal constant TYPEHASH = hex"4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa0";

  enum LinkPaymentType {
    DIRECT_FUNDING,
    SUBSCRIPTION
  }

  struct DirectFundingData {
    address recipient; // recipient of the top-up
    uint256 topupThreshold; // set to zero to disable auto-topup
    uint256 topupAmount;
  }

  function getUserOpFullHash(bytes32 userOpHash) internal pure returns (bytes32 fullHash) {
    bytes32 hashOfEncoding = keccak256(abi.encode(SCALibrary.TYPEHASH, userOpHash));
    fullHash = keccak256(abi.encodePacked(bytes1(0x19), bytes1(0x01), SCALibrary.DOMAIN_SEPARATOR, hashOfEncoding));
  }
}
