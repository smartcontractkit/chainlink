// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Consumer} from "./Consumer.sol";

contract BasicConsumer is Consumer {
  constructor(address _link, address _oracle, bytes32 _specId) {
    _setChainlinkToken(_link);
    _setChainlinkOracle(_oracle);
    s_specId = _specId;
  }
}
