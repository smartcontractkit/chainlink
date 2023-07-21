// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../ConfirmedOwner.sol";
import "../../interfaces/LinkTokenInterface.sol";
import "../interfaces/IVRFSubscriptionV2Plus.sol";
import "../interfaces/IVRFMigratableConsumerV2Plus.sol";

/// @notice The VRFV2SubscriptionManager contract is a contract
/// @notice that manages subscriptions to the VRF service.
/// @notice It acts as the owner of the subscription on the VRF coordinator
/// @notice and is able to perform the administrative operations pertaining
/// @notice to a subscription owner, such as adding and removing a consumer,
/// @notice and canceling the subscription and recovering it's funds.
contract VRFV2PlusSubscriptionManager is ConfirmedOwner {
  constructor() ConfirmedOwner(msg.sender) {}

  /// @notice the subscription ID that is owned by this contract.
  /// @notice this needs to be combined with the VRF coordinator address
  /// @notice in order to be used for any subscription-related operations.
  uint64 public s_subId;
  /// @notice the VRF coordinator that this the subscription ID above is for.
  /// @notice in the event a migration occurs, both s_subId and s_vrfCoordinator
  /// @notice will have to change accordingly.
  IVRFSubscriptionV2Plus public s_vrfCoordinator;
  /// @notice the LINK token contract that is used to fund the subscription.
  /// @notice it may not be available on some chains, in which case ether
  /// @notice funding is used.
  LinkTokenInterface public s_linkToken;

  function setVRFCoordinator(address vrfCoordinator) external onlyOwner {
    require(vrfCoordinator != address(0), "Invalid address");
    s_vrfCoordinator = IVRFSubscriptionV2Plus(vrfCoordinator);
  }

  /// @notice setLinkToken sets the LINK token contract that is used to fund
  /// @notice the subscription and withdraw any LINK funds.
  /// @notice no integrity checks are done on the given address other than checking
  /// @notice that it is nonzero.
  function setLinkToken(address linkToken) external onlyOwner {
    require(linkToken != address(0), "Invalid address");
    s_linkToken = LinkTokenInterface(linkToken);
  }

  function createSubscription() external onlyOwner returns (uint64 subId) {
    subId = s_vrfCoordinator.createSubscription();
    s_subId = subId;
  }

  function addConsumer(address consumer) external onlyOwner {
    s_vrfCoordinator.addConsumer(s_subId, consumer);
  }

  function removeConsumer(address consumer) external onlyOwner {
    s_vrfCoordinator.removeConsumer(s_subId, consumer);
  }

  function cancelSubscription() public onlyOwner {
    s_vrfCoordinator.cancelSubscription(s_subId, address(this));
  }

  function fundSubscriptionWithEth() external payable onlyOwner {
    s_vrfCoordinator.fundSubscriptionWithEth{value: msg.value}(s_subId);
  }

  /// @notice need to transfer link to this contract first in order for
  /// @notice this to work.
  function fundSubscriptionWithLink(uint256 amount) external onlyOwner {
    bool success = s_linkToken.transferAndCall(address(s_vrfCoordinator), amount, abi.encode(s_subId));
    require(success, "Transfer failed");
  }

  function withdrawLink() external onlyOwner {
    require(s_linkToken.transfer(msg.sender, s_linkToken.balanceOf(address(this))), "Unable to transfer");
  }

  function withdrawEth() external onlyOwner {
    (bool success, ) = payable(msg.sender).call{value: address(this).balance}("");
    require(success, "Unable to transfer");
  }

  /// @notice acceptSubscriptionOwnerTransfer accepts the transfer of
  /// @notice ownership of the subscription to this contract.
  /// @notice it then sets the s_subId to the new subscription id.
  function acceptSubscriptionOwnerTransfer(uint64 subId) external onlyOwner {
    s_vrfCoordinator.acceptSubscriptionOwnerTransfer(subId);
    s_subId = subId;
  }

  function requestSubscriptionOwnerTransfer(address newOwner) external onlyOwner {
    s_vrfCoordinator.requestSubscriptionOwnerTransfer(s_subId, newOwner);
  }

  function migrateToNewCoordinator(address newCoordinator) external onlyOwner {
    (uint96 balance, uint96 ethBalance, address owner, address[] memory consumers) = s_vrfCoordinator.getSubscription(
      s_subId
    );
    require(owner == address(this), "Not owner");
    // cancel the subscription so that we can get all of the funds here
    // note that this can only be done if there are no pending requests
    // on the subscription.
    cancelSubscription();
    // create a new subscription on the new coordinator
    IVRFSubscriptionV2Plus newCoord = IVRFSubscriptionV2Plus(newCoordinator);
    uint64 newSubId = newCoord.createSubscription();
    // at this point we should have all the funds in this contract, so
    // transfer the funds to the new subscription
    if (balance > 0) {
      s_linkToken.transferAndCall(newCoordinator, balance, abi.encode(newSubId));
    }
    if (ethBalance > 0) {
      newCoord.fundSubscriptionWithEth{value: ethBalance}(newSubId);
    }
    // add the consumers to the subscription on the new coordinator
    // and set the new coordinator on the consumers
    // note that this is bounded by MAX_CONSUMERS in the coordinator
    for (uint256 i = 0; i < consumers.length; i++) {
      newCoord.addConsumer(newSubId, consumers[i]);
      IVRFMigratableConsumerV2Plus(consumers[i]).setConfig(newCoordinator, newSubId);
    }
    // set the subscription id and the vrf coordinator in this owner contract
    s_subId = newSubId;
    s_vrfCoordinator = newCoord;
  }
}
