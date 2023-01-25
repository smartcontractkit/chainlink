// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../dev/AuthorizedReceiver.sol";

contract AuthorizedReceiverTestHelper is AuthorizedReceiver {
  bool private s_canSetAuthorizedSenders = true;

  function changeSetAuthorizedSender(bool on) external {
    s_canSetAuthorizedSenders = on;
  }

  function verifyValidateAuthorizedSender() external view validateAuthorizedSender returns (bool) {
    return true;
  }

  function _canSetAuthorizedSenders() internal view override returns (bool) {
    return s_canSetAuthorizedSenders;
  }
}
