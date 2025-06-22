// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;


// Data expected to update reserves
struct UpdateReserves {
    uint256 totalMinted;
    uint256 totalReserve;
}

interface IReserveManager {
    function updateReserves(UpdateReserves memory updateReserves) external;
    event RequestReserveUpdate(uint256 requestId);
}
