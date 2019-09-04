pragma solidity 0.5.0;
pragma experimental ABIEncoderV2;

import "../dev/CoordinatorInterface.sol";

/// Computes the mean of the values the oracles pass it via fulfill method
contract MeanAggregator {

  // Relies on Coordinator's authorization of the oracles (no need to track
  // oracle authorization in this contract.)

  mapping(bytes32 /* service agreement ID */ => uint256) numOracles;
  mapping(bytes32 /* request ID */ => mapping(address /* oracle */ => bool))
    reported;
  mapping(bytes32 /* request ID */ => uint256) numberReported;

  // Current total for given request, divided by number of oracles reporting
  mapping(bytes32 /* request ID */ => uint256) average;
  // Remainder of total for given request from division by number of oracles.
  mapping(bytes32 /* request ID */ => uint256) remainder;

  function initiateJob(
    bytes32 _sAId, CoordinatorInterface.ServiceAgreement memory _sa)
    public returns (bool success, bytes memory message) {
    if (numOracles[_sAId] != 0) {
      return (false, bytes("job already initiated"));
    }
    if (_sa.oracles.length == 0) {
      return (false, bytes("must depend on at least one oracle"));
    }
    numOracles[_sAId] = _sa.oracles.length;
    success = true;
  }

  function fulfill(bytes32 _requestId, bytes32 _sAId, address _oracle,
                   bytes32 _value)
    public returns (bool success, bool complete, bytes memory response) {
    if (reported[_requestId][_oracle]) {
      return (false, false, "oracle already reported");
    }
    uint256 oDividend = uint256(_value) / numOracles[_sAId];
    uint256 oRemainder = uint256(_value) % numOracles[_sAId];
    uint256 newRemainder = remainder[_requestId] + oRemainder;
    uint256 newAverage = average[_requestId] + oDividend + (newRemainder / numOracles[_sAId]);
    assert(newAverage >= average[_requestId]); // No overflow
    average[_requestId] = newAverage;
    remainder[_requestId] = newRemainder % numOracles[_sAId];
    numberReported[_requestId] += 1;
    success = true;
    reported[_requestId][_oracle] = true;
    complete = (numberReported[_requestId] == numOracles[_sAId]);
    if (complete) {
      response = abi.encode(average[_requestId]);
    }
  }
}
