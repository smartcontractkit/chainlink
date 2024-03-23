// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwnerWithProposal} from "../../shared/access/ConfirmedOwnerWithProposal.sol";
import {AuthorizedReceiver} from "./AuthorizedReceiver.sol";
import {Address} from "@openzeppelin/contracts/utils/Address.sol";

// solhint-disable gas-custom-errors
contract AuthorizedForwarder is ConfirmedOwnerWithProposal, AuthorizedReceiver {
  using Address for address;

  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable linkToken;

  event OwnershipTransferRequestedWithMessage(address indexed from, address indexed to, bytes message);

  constructor(
    address link,
    address owner,
    address recipient,
    bytes memory message
  ) ConfirmedOwnerWithProposal(owner, recipient) {
    require(link != address(0), "Link token cannot be a zero address");
    linkToken = link;
    if (recipient != address(0)) {
      emit OwnershipTransferRequestedWithMessage(owner, recipient, message);
    }
  }

  string public constant typeAndVersion = "AuthorizedForwarder 1.1.0";

  // @notice Forward a call to another contract
  // @dev Only callable by an authorized sender
  // @param to address
  // @param data to forward
  function forward(address to, bytes calldata data) external validateAuthorizedSender {
    require(to != linkToken, "Cannot forward to Link token");
    _forward(to, data);
  }

  //  @notice Forward multiple calls to other contracts in a multicall style
  //  @dev Only callable by an authorized sender
  //  @param tos An array of addresses to forward the calls to
  //  @param datas An array of data to forward to each corresponding address
  function multiForward(address[] calldata tos, bytes[] calldata datas) external validateAuthorizedSender {
    require(tos.length == datas.length, "Arrays must have the same length");

    for (uint256 i = 0; i < tos.length; ++i) {
      address to = tos[i];
      require(to != linkToken, "Cannot forward to Link token");

      // Perform the forward operation
      _forward(to, datas[i]);
    }
  }

  // @notice Forward a call to another contract
  // @dev Only callable by the owner
  // @param to address
  // @param data to forward
  function ownerForward(address to, bytes calldata data) external onlyOwner {
    _forward(to, data);
  }

  // @notice Transfer ownership with instructions for recipient
  // @param to address proposed recipient of ownership
  // @param message instructions for recipient upon accepting ownership
  function transferOwnershipWithMessage(address to, bytes calldata message) external {
    transferOwnership(to);
    emit OwnershipTransferRequestedWithMessage(msg.sender, to, message);
  }

  // @notice concrete implementation of AuthorizedReceiver
  // @return bool of whether sender is authorized
  function _canSetAuthorizedSenders() internal view override returns (bool) {
    return owner() == msg.sender;
  }

  // @notice common forwarding functionality and validation
  function _forward(address to, bytes calldata data) private {
    require(to.isContract(), "Must forward to a contract");
    // solhint-disable-next-line avoid-low-level-calls
    (bool success, bytes memory result) = to.call(data);
    if (!success) {
      if (result.length == 0) revert("Forwarded call reverted without reason");
      assembly {
        revert(add(32, result), mload(result))
      }
    }
  }
}
