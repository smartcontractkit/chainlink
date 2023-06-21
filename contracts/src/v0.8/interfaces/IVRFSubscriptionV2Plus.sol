// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IVRFSubscriptionV2Plus {
  function addConsumer(uint64 subId, address consumer) external;

  function removeConsumer(uint64 subId, address consumer) external;

  function cancelSubscription(uint64 subId, address to) external;

  function acceptSubscriptionOwnerTransfer(uint64 subId) external;

  function requestSubscriptionOwnerTransfer(uint64 subId, address newOwner) external;

  function createSubscription() external returns (uint64 subId);

  function getSubscription(
    uint64 subId
  ) external view returns (uint96 balance, uint96 ethBalance, address owner, address[] memory consumers);

  function fundSubscriptionWithEth(uint64 subId) external payable;
}
