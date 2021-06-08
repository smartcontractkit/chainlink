pragma solidity 0.5.0;

contract ServiceAgreementDecoder {

  struct ServiceAgreement {
    uint256 payment;
    uint256 expiration;
    uint256 endAt;
    address[] oracles;
    // This effectively functions as an ID tag for the off-chain job of the
    // service agreement. It is calculated as the keccak256 hash of the
    // normalized JSON request to create the ServiceAgreement, but that identity
    // is unused, and its value is essentially arbitrary.
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

  function decodeServiceAgreement(
    bytes memory _serviceAgreementData
  )
    internal
    pure
    returns(ServiceAgreement memory)
  {
    // solhint-disable indent
    ServiceAgreement memory agreement;

    ( agreement.payment,
      agreement.expiration,
      agreement.endAt,
      agreement.oracles,
      agreement.requestDigest,
      agreement.aggregator,
      agreement.aggInitiateJobSelector,
      agreement.aggFulfillSelector) =
      abi.decode(
        _serviceAgreementData,
        ( uint256,
        uint256,
        uint256,
        address[],
        bytes32,
        address,
        bytes4,
        bytes4 )
      );

    return agreement;
  }
}
