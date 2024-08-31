// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IAny2EVMMessageReceiver} from "../../../interfaces/IAny2EVMMessageReceiver.sol";
import {Client} from "../../../libraries/Client.sol";

import {IERC165} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";

contract MaybeRevertMessageReceiver is IAny2EVMMessageReceiver, IERC165 {
  error ReceiveRevert();
  error CustomError(bytes err);

  event ValueReceived(uint256 amount);
  event MessageReceived();

  address private s_manager;
  bool public s_toRevert;
  bytes private s_err;

  constructor(bool toRevert) {
    s_manager = msg.sender;
    s_toRevert = toRevert;
  }

  function setRevert(bool toRevert) external {
    s_toRevert = toRevert;
  }

  function setErr(bytes memory err) external {
    s_err = err;
  }

  /// @notice IERC165 supports an interfaceId
  /// @param interfaceId The interfaceId to check
  /// @return true if the interfaceId is supported
  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == type(IAny2EVMMessageReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
  }

  function ccipReceive(Client.Any2EVMMessage calldata) external override {
    if (s_toRevert) {
      revert CustomError(s_err);
    }
    emit MessageReceived();
  }

  receive() external payable {
    if (s_toRevert) {
      revert ReceiveRevert();
    }

    emit ValueReceived(msg.value);
  }
}
