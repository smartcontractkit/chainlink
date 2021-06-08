pragma solidity 0.5.0;

import "../dev/CoordinatorInterface.sol";
import "../dev/ServiceAgreementDecoder.sol";

/// Computes the mean of the values the oracles pass it via fulfill method
contract MeanAggregator is ServiceAgreementDecoder {

  // Relies on Coordinator's authorization of the oracles (no need to track
  // oracle authorization in this contract.)

  mapping(bytes32 /* service agreement ID */ => uint256) payment;
  mapping(bytes32 /* service agreement ID */ => address[]) oracles;
  mapping(bytes32 /* request ID */ => uint256) numberReported;
  mapping(bytes32 /* request ID */ => mapping(address => uint256)) reportingOrder;

  // Current total for given request, divided by number of oracles reporting
  mapping(bytes32 /* request ID */ => uint256) average;
  // Remainder of total for given request from division by number of oracles.
  mapping(bytes32 /* request ID */ => uint256) remainder;

  function initiateJob(
    bytes32 _sAId, bytes memory _serviceAgreementData)
    public returns (bool success, bytes memory message) {
      ServiceAgreement memory serviceAgreement = decodeServiceAgreement(_serviceAgreementData);

      if (oracles[_sAId].length != 0) {
        return (false, bytes("job already initiated"));
      }
      if (serviceAgreement.oracles.length == 0) {
        return (false, bytes("must depend on at least one oracle"));
      }
      oracles[_sAId] = serviceAgreement.oracles;
      payment[_sAId] = serviceAgreement.payment;
      success = true;
    }

  function fulfill(bytes32 _requestId, bytes32 _sAId, address _oracle,
    bytes32 _value) public
    returns (bool success, bool complete, bytes memory response,
    int256[] memory paymentAmounts) {
      if (reportingOrder[_requestId][_oracle] != 0 ||
        numberReported[_requestId] == oracles[_sAId].length) {
        return (false, false, "oracle already reported", paymentAmounts);
      }
      uint256 oDividend = uint256(_value) / oracles[_sAId].length;
      uint256 oRemainder = uint256(_value) % oracles[_sAId].length;
      uint256 newRemainder = remainder[_requestId] + oRemainder;
      uint256 newAverage = average[_requestId] + oDividend + (newRemainder / oracles[_sAId].length);
      assert(newAverage >= average[_requestId]); // No overflow
      average[_requestId] = newAverage;
      remainder[_requestId] = newRemainder % oracles[_sAId].length;
      numberReported[_requestId] += 1;
      reportingOrder[_requestId][_oracle] = numberReported[_requestId];
      success = true;
      complete = (numberReported[_requestId] == oracles[_sAId].length);
      if (complete) {
        response = abi.encode(average[_requestId]);
        paymentAmounts = calculatePayments(_sAId, _requestId);
      }
    }

  function calculatePayments(bytes32 _sAId, bytes32 _requestId) private returns (int256[] memory paymentAmounts) {
    paymentAmounts = new int256[](oracles[_sAId].length);
    uint256 numOracles = oracles[_sAId].length;
    uint256 totalPayment = payment[_sAId];
    for (uint256 oIdx = 0; oIdx < oracles[_sAId].length; oIdx++) {
      // Linearly-decaying payout determined by each oracle's reportingIdx
      uint256 reportingIdx = reportingOrder[_requestId][oracles[_sAId][oIdx]] - 1;
      paymentAmounts[oIdx] = int256(2*(totalPayment/numOracles) - (
        (totalPayment * ((2*reportingIdx) + 1)) / (numOracles**2)));
      delete reportingOrder[_requestId][oracles[_sAId][oIdx]];
    }
  }
}
