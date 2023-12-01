// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {TypeAndVersionInterface} from "../../../interfaces/TypeAndVersionInterface.sol";
// solhint-disable-next-line no-unused-import
import {ForwarderInterface} from "../interfaces/ForwarderInterface.sol";

// ./dev dependencies - to be moved from ./dev after audit
import {CrossDomainForwarder} from "../CrossDomainForwarder.sol";
import {CrossDomainOwnable} from "../CrossDomainOwnable.sol";

import {IScrollMessenger} from "@scroll-tech/contracts/libraries/IScrollMessenger.sol";
import {Address} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";

///
/// @title ScrollCrossDomainForwarder - L1 xDomain account representation
/// @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
/// @dev Any other L2 contract which uses this contract's address as a privileged position,
///  can be considered to be owned by the `l1Owner`
///
contract ScrollCrossDomainForwarder is TypeAndVersionInterface, CrossDomainForwarder {
  IScrollMessenger private immutable i_SCROLL_CROSS_DOMAIN_MESSENGER;

  ///
  /// @notice creates a new Scroll xDomain Forwarder contract
  /// @param crossDomainMessengerAddr the xDomain bridge messenger (Scroll bridge L2) contract address
  /// @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
  ///
  constructor(IScrollMessenger crossDomainMessengerAddr, address l1OwnerAddr) CrossDomainOwnable(l1OwnerAddr) {
    // solhint-disable-next-line custom-errors
    require(address(crossDomainMessengerAddr) != address(0), "Invalid xDomain Messenger address");
    i_SCROLL_CROSS_DOMAIN_MESSENGER = crossDomainMessengerAddr;
  }

  ///
  /// @notice versions:
  ///
  /// - ScrollCrossDomainForwarder 1.0.0: initial release
  ///
  /// @inheritdoc TypeAndVersionInterface
  ///
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "ScrollCrossDomainForwarder 1.0.0";
  }

  ///
  /// @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` being the L1 owner address
  /// @inheritdoc ForwarderInterface
  ///
  function forward(address target, bytes memory data) external virtual override onlyL1Owner {
    Address.functionCall(target, data, "Forwarder call reverted");
  }

  ///
  /// @notice This is always the address of the Scroll Cross Domain Messenger contract
  ///
  function crossDomainMessenger() public view returns (address) {
    return address(i_SCROLL_CROSS_DOMAIN_MESSENGER);
  }

  ///
  /// @notice The call MUST come from the L1 owner (via cross-chain message.) Reverts otherwise.
  ///
  modifier onlyL1Owner() override {
    // solhint-disable-next-line custom-errors
    require(msg.sender == crossDomainMessenger(), "Sender is not the L2 messenger");
    // solhint-disable-next-line custom-errors
    require(
      IScrollMessenger(crossDomainMessenger()).xDomainMessageSender() == l1Owner(),
      "xDomain sender is not the L1 owner"
    );
    _;
  }

  ///
  /// @notice The call MUST come from the proposed L1 owner (via cross-chain message.) Reverts otherwise.
  ///
  modifier onlyProposedL1Owner() override {
    address messenger = crossDomainMessenger();
    // solhint-disable-next-line custom-errors
    require(msg.sender == messenger, "Sender is not the L2 messenger");
    // solhint-disable-next-line custom-errors
    require(IScrollMessenger(messenger).xDomainMessageSender() == s_l1PendingOwner, "Must be proposed L1 owner");
    _;
  }
}
