pragma solidity 0.4.24;
pragma experimental ABIEncoderV2; // solium-disable-line no-experimental 

contract ChainlinkRequestInterface {
  struct Request {
    bytes32 sAId;
    address callbackAddress;
    bytes4 callbackFunctionId;
    uint256 nonce;
    uint256 dataVersion;
    bytes data;
  }

  function oracleRequest(
    address sender,
    uint256 payment,
    Request memory request
  ) public;

  function cancelOracleRequest(
    bytes32 requestId,
    uint256 payment,
    bytes4 callbackFunctionId,
    uint256 expiration
  ) external;
}
