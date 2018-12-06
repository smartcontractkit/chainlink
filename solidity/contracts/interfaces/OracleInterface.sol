pragma solidity 0.4.24;

interface OracleInterface {
  function cancel(bytes32 externalId) external;
  function fulfillData(uint256 internalId, bytes32 data) external returns (bool);
  function getAuthorizationStatus(address node) external view returns (bool);
  function requestData(
    address sender,
    uint256 amount,
    uint256 version,
    bytes32 specId,
    address callbackAddress,
    bytes4 callbackFunctionId,
    bytes32 externalId,
    bytes data
  ) external;
  function setFulfillmentPermission(address node, bool allowed) external;
  function withdraw(address recipient, uint256 amount) external;
  function withdrawable() external view returns (uint256);
}
