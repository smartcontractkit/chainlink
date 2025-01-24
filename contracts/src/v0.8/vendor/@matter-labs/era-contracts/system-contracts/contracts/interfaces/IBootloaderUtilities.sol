// SPDX-License-Identifier: MIT
// We use a floating point pragma here so it can be used within other projects that interact with the ZKsync ecosystem without using our exact pragma version.
pragma solidity ^0.8.20;

import {Transaction} from "../libraries/TransactionHelper.sol";

interface IBootloaderUtilities {
    function getTransactionHashes(
        Transaction calldata _transaction
    ) external view returns (bytes32 txHash, bytes32 signedTxHash);
}
