import "../consumers/SCA.sol";
import "../consumers/SmartContractAccountFactory.sol";
import "../consumers/SCALibrary.sol";

pragma solidity ^0.8.6;

library SmartContractAccountHelper {
  bytes constant initailizeCode = type(SCA).creationCode;

  function getFullEndTxEncoding(
    address endContract,
    uint256 value,
    uint256 deadline,
    bytes memory data
  ) public view returns (bytes memory encoding) {
    encoding = bytes.concat(
      SCA.executeTransactionFromEntryPoint.selector,
      abi.encode(endContract, value, block.timestamp + deadline, data)
    );
  }

  function getFullHashForSigning(
    bytes memory fullEndTxEncoding,
    address owner,
    uint256 nonce
  ) public view returns (bytes32 fullHash) {
    bytes32 hashOfEncoding = keccak256(abi.encode(SCALibrary.TYPEHASH, fullEndTxEncoding, owner, nonce, block.chainid));
    fullHash = keccak256(abi.encodePacked(bytes1(0x19), bytes1(0x01), SCALibrary.DOMAIN_SEPARATOR, hashOfEncoding));
  }

  function getSCAInitCodeWithConstructor(
    address owner,
    address entryPoint
  ) public pure returns (bytes memory initCode) {
    initCode = bytes.concat(initailizeCode, abi.encode(owner, entryPoint));
  }

  function getInitCode(
    address factory,
    address owner,
    address entryPoint
  ) external pure returns (bytes memory initCode) {
    bytes32 salt = bytes32(uint256(uint160(owner)) << 96);
    bytes memory initializeCodeWithConstructor = bytes.concat(initailizeCode, abi.encode(owner, entryPoint));
    initCode = bytes.concat(
      bytes20(address(factory)),
      abi.encodeWithSelector(
        SmartContractAccountFactory.deploySmartContractAccount.selector,
        salt,
        initializeCodeWithConstructor
      )
    );
  }

  /// @dev Computes the smart contract address that results from a CREATE2 operation, per EIP-1014.
  function calculateSmartContractAccountAddress(
    address owner,
    address entryPoint,
    address factory
  ) external pure returns (address) {
    bytes32 salt = bytes32(uint256(uint160(owner)) << 96);
    bytes memory initializeCodeWithConstructor = bytes.concat(initailizeCode, abi.encode(owner, entryPoint));
    bytes32 initializeCodeHash = keccak256(initializeCodeWithConstructor);
    return address(uint160(uint256(keccak256(abi.encodePacked(hex"ff", address(factory), salt, initializeCodeHash)))));
  }
}
