// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "./interfaces/OperatorInterface.sol";
import "./ConfirmedOwnerWithProposal.sol";
import "./AuthorizedReceiver.sol";
import "./vendor/Address.sol";

contract AuthorizedForwarder is ConfirmedOwnerWithProposal, AuthorizedReceiver {
  using Address for address;

  address public immutable getChainlinkToken;

  event OwnershipTransferRequestedWithMessage(address indexed from, address indexed to, bytes message);

  constructor(
    address link,
    address owner,
    address recipient,
    bytes memory message
  ) ConfirmedOwnerWithProposal(owner, recipient) {
    require(link != address(0));
    getChainlinkToken = link;
    if (recipient != address(0)) {
      emit OwnershipTransferRequestedWithMessage(owner, recipient, message);
    }
  }

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion() external pure virtual returns (string memory) {
    return "AuthorizedForwarder 1.0.0";
  }

  /**
   * @notice Forward a call to another contract
   * @dev Only callable by an authorized sender
   * @param to address
   * @param data to forward
   */
  function forward(address to, bytes calldata data) external validateAuthorizedSender {
    require(to != getChainlinkToken, "Cannot forward to Link token");
    _forward(to, data);
  }

  /**
   * @notice Forward a call to another contract
   * @dev Only callable by the owner
   * @param to address
   * @param data to forward
   */
  function ownerForward(address to, bytes calldata data) external onlyOwner {
    _forward(to, data);
  }

  /**
   * @notice Transfer ownership with instructions for recipient
   * @param to address proposed recipient of ownership
   * @param message instructions for recipient upon accepting ownership
   */
  function transferOwnershipWithMessage(address to, bytes calldata message) external {
    transferOwnership(to);
    emit OwnershipTransferRequestedWithMessage(msg.sender, to, message);
  }

  /**
   * @notice concrete implementation of AuthorizedReceiver
   * @return bool of whether sender is authorized
   */
  function _canSetAuthorizedSenders() internal view override returns (bool) {
    return owner() == msg.sender;
  }

  /**
   * @notice common forwarding functionality and validation
   */
  function _forward(address to, bytes calldata data) private {
    require(to.isContract(), "Must forward to a contract");
    (bool success, bytes memory result) = to.call(data);
    if (!success) {
      if (result.length == 0) revert("Forwarded call reverted without reason");
      assembly {
        revert(add(32, result), mload(result))
      }
    }
  }
}
