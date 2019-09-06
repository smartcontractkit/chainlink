pragma solidity 0.5.0;
pragma experimental ABIEncoderV2;

contract CoordinatorInterface {

  struct ServiceAgreement {
    uint256 payment;
    uint256 expiration;
    uint256 endAt;
    address[] oracles;
    // This effectively functions as an ID tag for the service agreement. It is
    // calculated as the keccak256 hash of the normalized JSON request to create
    // the ServiceAgreement, but that identity is unused, and its value is
    // essentially arbitrary.
    bytes32 requestDigest;

    // Specification of aggregator interface. See ../tests/MeanAggregator.sol
    // for example
    address aggregator;
    // Selectors for the interface methods must be specified, because their
    // arguments can vary from aggregator to aggregator.
    //
    // Function selector for aggregator initiateJob method
    bytes4 aggInitiateJobSelector;
    // Function selector for aggregator fulfill method
    bytes4 aggFulfillSelector;
  }

  struct OracleSignatures {
    uint8[] vs;
    bytes32[] rs;
    bytes32[] ss;
  }

  function initiateServiceAgreement(
    ServiceAgreement memory _agreement,
    OracleSignatures memory _signatures)
    public returns (bytes32);

  function fulfillOracleRequest(
    bytes32 _requestId,
    bytes32 _aggregatorArgs)
    external returns (bool);
}
