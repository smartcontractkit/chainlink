// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {IReceiver} from "../../interfaces/IReceiver.sol";

contract MaliciousInterfaceReceiver is IReceiver, IERC165 {
  error InsufficientGasProvided();
  event GasProvided(uint256 gasProvided, uint256 gasExpected, bool sufficient);
  event MessageReceived(bytes metadata, bytes[] mercuryReports);
  bytes public latestReport;
  uint256 internal s_expectedGasLimit;

  constructor(uint256 expectedGasLimit) {
    s_expectedGasLimit = expectedGasLimit;
  }

  function onReport(bytes calldata, bytes calldata) external {
    uint256 providedGas = gasleft();
    emit GasProvided(providedGas, s_expectedGasLimit, providedGas >= s_expectedGasLimit);

    if (providedGas < s_expectedGasLimit) {
      revert InsufficientGasProvided();
    }
  }

  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    // Consume up to the maximum amount of gas that can be consumed in this check
    // This loop consumes roughly 29_000 gas
    for (uint256 i = 0; i < 670; i++) {}

    return interfaceId == type(IReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
  }
}
