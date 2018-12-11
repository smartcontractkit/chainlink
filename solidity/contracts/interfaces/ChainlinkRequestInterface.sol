pragma solidity 0.4.24;

interface ChainlinkRequestInterface {
  function cancel(bytes32 externalId) external;
  function requestData(
    address sender,
    uint256 amount,
    uint256 version,
    bytes32 id,
    address callbackAddress,
    bytes4 callbackFunctionId,
    bytes32 externalId,
    bytes data
  ) external;
}