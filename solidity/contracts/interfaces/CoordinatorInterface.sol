pragma solidity 0.4.24;

interface CoordinatorInterface {
  function executeServiceAgreement(
    address sender,
    uint256 amount,
    uint256 version,
    bytes32 sAId,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 nonce,
    bytes data
  ) external;
}
