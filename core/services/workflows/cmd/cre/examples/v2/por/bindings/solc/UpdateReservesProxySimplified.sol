// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IReserveManager, UpdateReserves} from "./IReserveManager.sol";
import {IReceiverTemplate} from "./keystone/IReceiverTemplate.sol";

contract UpdateReservesProxySimplified is IReceiverTemplate {
    IReserveManager public reserveManager;

    constructor(
        address _reserveManager,
        address expectedAuthor,
        bytes10 expectedWorkflowName
    ) IReceiverTemplate(expectedAuthor, expectedWorkflowName) {
        reserveManager = IReserveManager(_reserveManager);
    }

    /// @inheritdoc IReceiverTemplate
    function _processReport(bytes calldata report) internal override {
        UpdateReserves memory updateReservesData = abi.decode(
            report,
            (UpdateReserves)
        );
        reserveManager.updateReserves(updateReservesData);
    }
}