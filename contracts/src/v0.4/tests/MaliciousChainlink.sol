pragma solidity 0.4.24;

import { CBOR as CBOR_Chainlink } from "../vendor/CBOR.sol";
import { Buffer as Buffer_Chainlink } from "../vendor/Buffer.sol";

library MaliciousChainlink {
  using CBOR_Chainlink for Buffer_Chainlink.buffer;

  struct Request {
    bytes32 specId;
    address callbackAddress;
    bytes4 callbackFunctionId;
    uint256 nonce;
    Buffer_Chainlink.buffer buf;
  }

  struct WithdrawRequest {
    bytes32 specId;
    address callbackAddress;
    bytes4 callbackFunctionId;
    uint256 nonce;
    Buffer_Chainlink.buffer buf;
  }

  function initializeWithdraw(
    WithdrawRequest memory self,
    bytes32 _specId,
    address _callbackAddress,
    bytes4 _callbackFunction
  ) internal pure returns (MaliciousChainlink.WithdrawRequest memory) {
    Buffer_Chainlink.init(self.buf, 128);
    self.specId = _specId;
    self.callbackAddress = _callbackAddress;
    self.callbackFunctionId = _callbackFunction;
    return self;
  }

  function add(Request memory self, string _key, string _value)
    internal pure
  {
    self.buf.encodeString(_key);
    self.buf.encodeString(_value);
  }

  function addBytes(Request memory self, string _key, bytes _value)
    internal pure
  {
    self.buf.encodeString(_key);
    self.buf.encodeBytes(_value);
  }

  function addInt(Request memory self, string _key, int256 _value)
    internal pure
  {
    self.buf.encodeString(_key);
    self.buf.encodeInt(_value);
  }

  function addUint(Request memory self, string _key, uint256 _value)
    internal pure
  {
    self.buf.encodeString(_key);
    self.buf.encodeUInt(_value);
  }

  function addStringArray(Request memory self, string _key, string[] memory _values)
    internal pure
  {
    self.buf.encodeString(_key);
    self.buf.startArray();
    for (uint256 i = 0; i < _values.length; i++) {
      self.buf.encodeString(_values[i]);
    }
    self.buf.endSequence();
  }
}
