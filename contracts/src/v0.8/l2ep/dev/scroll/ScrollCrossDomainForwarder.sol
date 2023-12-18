// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {TypeAndVersionInterface} from "../../../interfaces/TypeAndVersionInterface.sol";
import {ForwarderInterface} from "../interfaces/ForwarderInterface.sol";

import {CrossDomainForwarder} from "../CrossDomainForwarder.sol";
import {CrossDomainOwnable} from "../CrossDomainOwnable.sol";

import {IScrollMessenger} from "@scroll-tech/contracts/libraries/IScrollMessenger.sol";
import {Address} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";

/// @title ScrollCrossDomainForwarder - L1 xDomain account representation
/// @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
/// @dev Any other L2 contract which uses this contract's address as a privileged position,
/// can be considered to be owned by the `l1Owner`
contract ScrollCrossDomainForwarder is TypeAndVersionInterface, CrossDomainForwarder {
  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "ScrollCrossDomainForwarder 1.0.0";

  address internal immutable i_scrollCrossDomainMessenger;

  /// @param crossDomainMessengerAddr the xDomain bridge messenger (Scroll bridge L2) contract address
  /// @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
  constructor(IScrollMessenger crossDomainMessengerAddr, address l1OwnerAddr) CrossDomainOwnable(l1OwnerAddr) {
    // solhint-disable-next-line custom-errors
    require(address(crossDomainMessengerAddr) != address(0), "Invalid xDomain Messenger address");
    i_scrollCrossDomainMessenger = address(crossDomainMessengerAddr);
  }

  /// @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` being the L1 owner address
  /// @inheritdoc ForwarderInterface
  function forward(address target, bytes memory data) external override onlyL1Owner {
    Address.functionCall(target, data, "Forwarder call reverted");
  }

  /// @notice This is always the address of the Scroll Cross Domain Messenger contract
  function crossDomainMessenger() external view returns (address) {
    return address(i_scrollCrossDomainMessenger);
  }

  /// @notice The call MUST come from the L1 owner (via cross-chain message.) Reverts otherwise.
  modifier onlyL1Owner() override {
    // solhint-disable-next-line custom-errors
    require(msg.sender == i_scrollCrossDomainMessenger, "Sender is not the L2 messenger");
    // solhint-disable-next-line custom-errors
    require(
      IScrollMessenger(i_scrollCrossDomainMessenger).xDomainMessageSender() == l1Owner(),
      "xDomain sender is not the L1 owner"
    );
    _;
  }

  /// @notice The call MUST come from the proposed L1 owner (via cross-chain message.) Reverts otherwise.
  modifier onlyProposedL1Owner() override {
    // solhint-disable-next-line custom-errors
    require(msg.sender == i_scrollCrossDomainMessenger, "Sender is not the L2 messenger");
    // solhint-disable-next-line custom-errors
    require(
      IScrollMessenger(i_scrollCrossDomainMessenger).xDomainMessageSender() == s_l1PendingOwner,
      "Must be proposed L1 owner"
    );
    _;
  }
}
