// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../Chainlink.sol";
import "../vendor/CBORChainlink.sol";
import "../vendor/BufferChainlink.sol";

contract ChainlinkTestHelper {
  using Chainlink for Chainlink.Request;
  using CBORChainlink for BufferChainlink.buffer;

  Chainlink.Request private req;

  event RequestData(bytes payload);

  function closeEvent() public {
    emit RequestData(req.buf.buf);
  }

  function setBuffer(bytes memory data) public {
    Chainlink.Request memory r2 = req;
    r2.setBuffer(data);
    req = r2;
  }

  function add(string memory _key, string memory _value) public {
    Chainlink.Request memory r2 = req;
    r2.add(_key, _value);
    req = r2;
  }

  function addBytes(string memory _key, bytes memory _value) public {
    Chainlink.Request memory r2 = req;
    r2.addBytes(_key, _value);
    req = r2;
  }

  function addInt(string memory _key, int256 _value) public {
    Chainlink.Request memory r2 = req;
    r2.addInt(_key, _value);
    req = r2;
  }

  function addUint(string memory _key, uint256 _value) public {
    Chainlink.Request memory r2 = req;
    r2.addUint(_key, _value);
    req = r2;
  }

  // Temporarily have method receive bytes32[] memory until experimental
  // string[] memory can be invoked from truffle tests.
  function addStringArray(string memory _key, string[] memory _values) public {
    Chainlink.Request memory r2 = req;
    r2.addStringArray(_key, _values);
    req = r2;
  }
}
