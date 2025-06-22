// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./IReserveManager.sol";

contract ReserveManager is IReserveManager {
    uint256 public lastTotalMinted;
    uint256 public lastTotalReserve;
    uint256 private requestIdCounter;

    function updateReserves(UpdateReserves memory updateReserves) external override {
        lastTotalMinted = updateReserves.totalMinted;
        lastTotalReserve = updateReserves.totalReserve;
        
        requestIdCounter++;
        emit RequestReserveUpdate(requestIdCounter);
    }
}