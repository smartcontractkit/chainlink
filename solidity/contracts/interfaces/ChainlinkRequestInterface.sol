pragma solidity 0.4.24;

interface ChainlinkRequestInterface {
  function requestData(
    address sender,
    uint256 payment,
    uint256 version,
    bytes32 id,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 expiration,
    bytes data
  ) external;

  function cancel(
    bytes32 requestId,
    uint256 payment,
    bytes4 callbackFunctionId,
    uint256 expiration
  ) external;
}