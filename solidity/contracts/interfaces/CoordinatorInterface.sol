pragma solidity 0.4.24;

interface CoordinatorInterface {
  function initiateServiceAgreement(
    uint256 payment,
    uint256 expiration,
    uint256 endAt,
    address[] oracles,
    uint8[] vs,
    bytes32[] rs,
    bytes32[] ss,
    bytes32 requestDigest
  ) external returns (bytes32);
  function fulfillOracleRequest(bytes32 requestId, bytes32 data) external returns (bool);
}
