// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {DualAggregator} from "../../DualAggregator.sol";

import {Test} from "forge-std/Test.sol";

contract ReportGenerator is Test {
  int192 public constant JUELS_PER_FEECOIN = 1000000000000000000;

  struct SignedReport {
    bytes32[3] reportContext;
    bytes report;
    // ECDSA signatures
    bytes32[] rs;
    bytes32[] ss;
    bytes32 rawVs;
  }

  uint8 public s_f;
  uint32 public s_currentEpoch;
  uint256[] public s_oraclePrivateKeys;
  bytes32 public s_configDigest;
  DualAggregator public s_aggregator;

  constructor(DualAggregator aggregator, uint256[] memory privateKeys, bytes32 configDigest, uint8 f) Test() {
    s_aggregator = aggregator;
    s_f = f;
    s_configDigest = configDigest;

    for (uint256 i; i < privateKeys.length; ++i) {
      s_oraclePrivateKeys.push(privateKeys[i]);
    }
  }

  function generateSignedReport(int192 price, uint32 observationTimestamp) public returns (SignedReport memory) {
    // build report
    bytes memory report = _buildReport(price, observationTimestamp);

    // build report context
    bytes32[3] memory reportContext = _buildReportContext();

    // sign report
    (bytes32[] memory rs, bytes32[] memory ss, bytes32 rawVs) = _signReport(report, reportContext);

    ++s_currentEpoch;

    return SignedReport({reportContext: reportContext, report: report, rs: rs, ss: ss, rawVs: rawVs});
  }

  function _buildReport(int192 price, uint32 observationTimestamp) internal view returns (bytes memory) {
    // build the observers bytes (string of uint8 id numbers for each observer) and observations array
    bytes memory rawObservers = new bytes(s_oraclePrivateKeys.length);
    int192[] memory observations = new int192[](s_oraclePrivateKeys.length);
    for (uint256 i; i < s_oraclePrivateKeys.length; ++i) {
      rawObservers[i] = bytes1(uint8(i));
      observations[i] = price;
    }
    // We can convert because we know the length of the array is less than 32
    bytes32 observers = bytes32(rawObservers);

    return abi.encode(observationTimestamp, observers, observations, JUELS_PER_FEECOIN);
  }

  function _buildReportContext() internal view returns (bytes32[3] memory) {
    // reportContext[0]: ConfigDigest
    // reportContext[1]: 27 byte padding, 4-byte epoch and 1-byte round
    // reportContext[2]: ExtraHash

    // 27 bytes of padding, 4 bytes for epoch, 1 byte for round
    bytes memory rawEpochAndRound = abi.encodePacked(bytes27(0), s_currentEpoch, uint8(1));
    bytes32 epochAndRound = bytes32(rawEpochAndRound);

    return [s_configDigest, epochAndRound, bytes32(0)];
  }

  function _signReport(
    bytes memory report,
    bytes32[3] memory reportContext
  ) internal view returns (bytes32[] memory, bytes32[] memory, bytes32) {
    bytes32 hashedReport = keccak256(abi.encodePacked(keccak256(report), reportContext));

    uint256 numSigners = s_f + 1;
    bytes32[] memory rs = new bytes32[](numSigners);
    bytes32[] memory ss = new bytes32[](numSigners);
    bytes memory vs = new bytes(numSigners);
    for (uint256 i; i < numSigners; ++i) {
      (uint8 v, bytes32 r, bytes32 s) = vm.sign(s_oraclePrivateKeys[i], hashedReport);
      rs[i] = r;
      ss[i] = s;
      vs[i] = bytes1(v - 27);
    }

    // vs can be converted to bytes32 because we know the length of the array is less than 32
    bytes32 rawVs = bytes32(vs);

    return (rs, ss, rawVs);
  }
}
