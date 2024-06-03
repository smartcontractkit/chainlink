// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Consumer} from "./Consumer.sol";
import {Chainlink} from "../../../Chainlink.sol";

contract GasGuzzlingConsumer is Consumer {
  using Chainlink for Chainlink.Request;

  constructor(address _link, address _oracle, bytes32 _specId) {
    _setChainlinkToken(_link);
    _setChainlinkOracle(_oracle);
    s_specId = _specId;
  }

  function gassyRequestEthereumPrice(uint256 _payment) public {
    Chainlink.Request memory req = _buildChainlinkRequest(s_specId, address(this), this.gassyFulfill.selector);
    req._add("get", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = "USD";
    req._addStringArray("path", path);
    _sendChainlinkRequest(req, _payment);
  }

  function gassyFulfill(bytes32 _requestId, bytes32) public recordChainlinkFulfillment(_requestId) {
    while (true) {}
  }

  function gassyMultiWordRequest(uint256 _payment) public {
    Chainlink.Request memory req = _buildChainlinkRequest(s_specId, address(this), this.gassyMultiWordFulfill.selector);
    req._add("get", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = "USD";
    req._addStringArray("path", path);
    _sendChainlinkRequest(req, _payment);
  }

  function gassyMultiWordFulfill(bytes32 _requestId, bytes memory) public recordChainlinkFulfillment(_requestId) {
    while (true) {}
  }
}
