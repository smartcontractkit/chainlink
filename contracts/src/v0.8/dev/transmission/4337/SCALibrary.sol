// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

library SCALibrary {
  // keccak256("EIP712Domain(uint256 chainId, address verifyingContract)");
  bytes32 internal constant DOMAIN_SEPARATOR = hex"1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea61";

  // keccak256("executeTransactionFromEntryPoint(address to, uint256 value, bytes calldata data)");
  bytes32 internal constant TYPEHASH = hex"4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa0";

  enum LinkPaymentType {
    DIRECT_FUNDING,
    SUBSCRIPTION // TODO: implement
  }

  struct DirectFundingData {
    address recipient; // recipient of the top-up
    uint256 topupThreshold; // set to zero to disable auto-topup
    uint256 topupAmount;
  }

  function getUserOpFullHash(bytes32 userOpHash, address scaAddress) internal view returns (bytes32 fullHash) {
    bytes32 hashOfEncoding = keccak256(abi.encode(SCALibrary.TYPEHASH, userOpHash));
    fullHash = keccak256(
      abi.encodePacked(
        bytes1(0x19),
        bytes1(0x01),
        SCALibrary.DOMAIN_SEPARATOR,
        block.chainid,
        scaAddress,
        hashOfEncoding
      )
    );
  }

  function recoverSignature(bytes memory signature, bytes32 fullHash) internal pure returns (address) {
    bytes32 r;
    bytes32 s;
    assembly {
      r := mload(add(signature, 0x20))
      s := mload(add(signature, 0x40))
    }
    uint8 v = uint8(signature[64]);

    return ecrecover(fullHash, v + 27, r, s);
  }
}
