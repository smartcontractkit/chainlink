// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {IAuthorizedReceiver} from "./interfaces/IAuthorizedReceiver.sol";

abstract contract AuthorizedReceiver is IAuthorizedReceiver {
  using EnumerableSet for EnumerableSet.AddressSet;

  event AuthorizedSendersChanged(address[] senders, address changedBy);

  error EmptySendersList();
  error UnauthorizedSender();
  error NotAllowedToSetSenders();

  EnumerableSet.AddressSet private s_authorizedSenders;
  address[] private s_authorizedSendersList;

  /**
   * @notice Sets the fulfillment permission for a given node. Use `true` to allow, `false` to disallow.
   * @param senders The addresses of the authorized Chainlink node
   */
  function setAuthorizedSenders(address[] calldata senders) external override validateAuthorizedSenderSetter {
    if (senders.length == 0) {
      revert EmptySendersList();
    }
    for (uint256 i = 0; i < s_authorizedSendersList.length; i++) {
      s_authorizedSenders.remove(s_authorizedSendersList[i]);
    }
    for (uint256 i = 0; i < senders.length; i++) {
      s_authorizedSenders.add(senders[i]);
    }
    s_authorizedSendersList = senders;
    emit AuthorizedSendersChanged(senders, msg.sender);
  }

  /**
   * @notice Retrieve a list of authorized senders
   * @return array of addresses
   */
  function getAuthorizedSenders() public view override returns (address[] memory) {
    return s_authorizedSendersList;
  }

  /**
   * @notice Use this to check if a node is authorized for fulfilling requests
   * @param sender The address of the Chainlink node
   * @return The authorization status of the node
   */
  function isAuthorizedSender(address sender) public view override returns (bool) {
    return s_authorizedSenders.contains(sender);
  }

  /**
   * @notice customizable guard of who can update the authorized sender list
   * @return bool whether sender can update authorized sender list
   */
  function _canSetAuthorizedSenders() internal virtual returns (bool);

  /**
   * @notice validates the sender is an authorized sender
   */
  function _validateIsAuthorizedSender() internal view {
    if (!isAuthorizedSender(msg.sender)) {
      revert UnauthorizedSender();
    }
  }

  /**
   * @notice prevents non-authorized addresses from calling this method
   */
  modifier validateAuthorizedSender() {
    _validateIsAuthorizedSender();
    _;
  }

  /**
   * @notice prevents non-authorized addresses from calling this method
   */
  modifier validateAuthorizedSenderSetter() {
    if (!_canSetAuthorizedSenders()) {
      revert NotAllowedToSetSenders();
    }
    _;
  }
}
