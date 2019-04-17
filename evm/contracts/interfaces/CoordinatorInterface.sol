pragma solidity 0.4.24;
pragma experimental ABIEncoderV2;

contract CoordinatorInterface {

  struct ServiceAgreement {
    uint256 payment;
    uint256 expiration;
    uint256 endAt;
    address[] oracles;
    bytes32 requestDigest;
  }

  struct OracleSignatures {
    uint8[] vs;
    bytes32[] rs;
    bytes32[] ss;
  }

  function initiateServiceAgreement(
    ServiceAgreement memory _agreement, 
    OracleSignatures memory _signatures
  ) public returns (bytes32);

  function fulfillOracleRequest(bytes32 requestId, bytes32 data) external returns (bool);
}
