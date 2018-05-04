pragma solidity ^0.4.23;
pragma experimental ABIEncoderV2;

import "./Buffer.sol";
import "./CBOR.sol";

library ChainlinkLib {
  using CBOR for Buffer.buffer;

  struct Run {
    bytes32 id;
    bytes32 jobId;
    address callbackAddress;
    bytes4 callbackFunctionId;
    Buffer.buffer buf;
  }

  function add(Run memory self, string _key, string _value)
    internal
  {
    self.buf.encodeString(_key);
    self.buf.encodeString(_value);
  }

  function addStringArray(Run memory self, string _key, string[] memory _values)
    internal
  {
    self.buf.encodeString(_key);
    self.buf.startArray();
    for (uint256 i = 0; i < _values.length; i++) {
      self.buf.encodeString(_values[i]);
    }
    self.buf.endSequence();
  }

  function close(Run memory self)
    internal
    returns (bytes)
  {
     self.buf.endSequence();
     return self.buf.buf;
  }
}
