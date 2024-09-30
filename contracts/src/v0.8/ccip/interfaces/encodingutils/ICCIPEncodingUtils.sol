// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {CCIPHome} from "../../capability/CCIPHome.sol";
import {RMNRemote} from "../../rmn/RMNRemote.sol";

/// @dev this file exposes structs that are otherwise internal to the CCIP codebase
/// doing this allows those structs to be encoded and decoded with type safety in offchain code
/// and tests because generated wrappers are available
interface ICCIPEncodingUtils {
  /// @dev the RMN Report struct is used in integration / E2E tests
  function exposeRmnReport(bytes32 rmnReportVersion, RMNRemote.Report memory rmnReport) external;

  /// @dev the OCR3Config Config struct is used in integration / E2E tests
  function exposeOCR3Config(CCIPHome.OCR3Config[] calldata config) external view returns (bytes memory);
}
