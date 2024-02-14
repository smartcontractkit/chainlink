// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IMessageTransmitter} from "../../pools/USDC/IMessageTransmitter.sol";

contract MockUSDCTransmitter is IMessageTransmitter {
  // Indicated whether the receiveMessage() call should succeed.
  bool public s_shouldSucceed;
  uint32 private immutable i_version;
  uint32 private immutable i_localDomain;

  constructor(uint32 _version, uint32 _localDomain) {
    i_version = _version;
    i_localDomain = _localDomain;
    s_shouldSucceed = true;
  }

  function receiveMessage(bytes calldata, bytes calldata) external view returns (bool success) {
    return s_shouldSucceed;
  }

  function setShouldSucceed(bool shouldSucceed) external {
    s_shouldSucceed = shouldSucceed;
  }

  function version() external view returns (uint32) {
    return i_version;
  }

  function localDomain() external view returns (uint32) {
    return i_localDomain;
  }
}
