// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import {IVersioned} from "./interfaces/IVersioned.sol";

abstract contract Versioned is IVersioned {
  string internal s_id;
  uint16 internal s_version;

  constructor(string memory id, uint16 version) {
    s_id = id;
    s_version = version;
  }

  /**
   * @inheritdoc IVersioned
   */
  function idAndVersion() public view override returns (string memory id, uint16 version) {
    return (s_id, s_version);
  }
}
