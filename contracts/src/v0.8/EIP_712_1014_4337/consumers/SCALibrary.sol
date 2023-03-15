// SPDX-License-Identifier: MIT
pragma solidity >=0.4.22 <0.9.0;

library SCALibrary {
  // keccak256("EIP712Domain(uint256 chainId,address verifyingContract)");
  bytes32 internal constant DOMAIN_SEPARATOR = 0x47e79534a245952e8b16893a336b85a3d9ea9fa8c573f3d803afb92a79469218;

  // keccak256("SmartContractAccountTX(address to,uint256 value,bytes data)");
  bytes32 internal constant TYPEHASH = 0x23d294a3e6e5266ba3b17997d1a601816663066087b37cad4f06b4d4e30655d9;
}
