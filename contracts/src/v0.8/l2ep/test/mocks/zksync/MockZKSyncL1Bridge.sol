// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IBridgehub, L2TransactionRequestDirect} from "@zksync/contracts/l1-contracts/contracts/bridgehub/IBridgehub.sol";

contract MockZKSyncL1Bridge is IBridgehub {

  function requestL2TransactionDirect(
      L2TransactionRequestDirect calldata _request
  ) external payable returns (bytes32 canonicalTxHash) {
    emit "";
  }
}
