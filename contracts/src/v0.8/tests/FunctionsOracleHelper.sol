// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../dev/functions/FunctionsOracle.sol";

contract FunctionsOracleHelper is FunctionsOracle {
  function callValidateReport(bytes calldata report) external pure returns (bool isValid) {
    bytes32 configDigest;
    uint40 epochAndRound;
    isValid = _validateReport(configDigest, epochAndRound, report);
  }

  function callReport(bytes calldata report) external {
    address[maxNumOracles] memory signers;
    signers[0] = msg.sender;
    _report(gasleft(), msg.sender, 1, signers, report);
  }

  function callReportMultipleSigners(bytes calldata report, address secondSigner) external {
    address[maxNumOracles] memory signers;
    signers[0] = msg.sender;
    signers[1] = secondSigner;
    _report(gasleft(), msg.sender, 2, signers, report);
  }
}
