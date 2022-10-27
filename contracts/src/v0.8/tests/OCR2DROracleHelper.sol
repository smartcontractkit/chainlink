// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../dev/ocr2dr/OCR2DROracle.sol";

contract OCR2DROracleHelper is OCR2DROracle {
    function callReport(bytes calldata report) external {
        bytes32 configDigest;
        uint40 epochAndRound;
        _report(configDigest, epochAndRound, report);
    }
}
