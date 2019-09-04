pragma solidity 0.5.0;
pragma experimental ABIEncoderV2;

import "../dev/CoordinatorInterface.sol";

/// Used to check the basic aggregator/coordinator interactions. It does nothing
/// but emit its messages as certain types of events.
contract EmptyAggregator {

  event InitiatedJob(bytes32 said, CoordinatorInterface.ServiceAgreement sa);

  function initiateJob(
    bytes32 _saId, CoordinatorInterface.ServiceAgreement memory _sa)
    public returns (bool success, bytes memory response) {
    emit InitiatedJob(_saId, _sa);
    success = true;
  }

  event Fulfilled(
    bytes32 requestId,
    address oracle,
    bool success,
    bool complete,
    bytes fulfillment);

  function fulfill(bytes32 _requestId, bytes32 _saId, address _oracle,
                   bytes32 _fulfillment)
    public returns (bool success, bool complete, bytes memory response) {
    success = true;
    complete = true;
    response = abi.encode(_fulfillment);
    emit Fulfilled(_requestId, _oracle, success, complete, response);
  }
}
