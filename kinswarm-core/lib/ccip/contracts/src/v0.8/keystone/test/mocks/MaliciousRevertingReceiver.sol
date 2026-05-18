// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {IReceiver} from "../../interfaces/IReceiver.sol";
import {Test} from "forge-std/Test.sol";

/// A malicious receiver that uses max allowed for ERC165 checks and consumes all gas in `onReport()`
/// Causes parent Forwarder contract to revert if it doesn't handle gas tracking accurately
contract MaliciousRevertingReceiver is IReceiver, IERC165, Test {
  function onReport(bytes calldata, bytes calldata) external view override {
    // consumes about 63/64 of all gas available
    uint256 targetGasRemaining = 200;
    for (uint256 i = 0; gasleft() > targetGasRemaining; i++) {}
  }

  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    // Consume up to the maximum amount of gas that can be consumed in this check
    // This loop consumes roughly 29_000 gas
    for (uint256 i = 0; i < 670; i++) {}

    return interfaceId == type(IReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
  }
}
