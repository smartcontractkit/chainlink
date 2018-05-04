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
    bytes32 externalId;
    Buffer.buffer buf;
  }

  function initialize(
    Run memory self,
    bytes32 _jobId,
    address _callbackAddress,
    string _callbackFunctionSignature
  ) internal pure returns (ChainlinkLib.Run memory) {
    Buffer.init(self.buf, 128);
    self.jobId = _jobId;
    self.callbackAddress = _callbackAddress;
    self.callbackFunctionId = bytes4(keccak256(_callbackFunctionSignature));
    self.buf.startMap();
    return self;
  }

  function encodeForOracle(
    Run memory self,
    bytes4 _oracleFid,
    uint256 _clArgsVersion
  ) internal pure returns (bytes memory) {
    return abi.encodeWithSelector(
      _oracleFid,
      _clArgsVersion,
      self.jobId,
      self.callbackAddress,
      self.callbackFunctionId,
      self.externalId,
      self.buf.buf);
  }

  function add(Run memory self, string _key, string _value)
    internal pure
  {
    self.buf.encodeString(_key);
    self.buf.encodeString(_value);
  }

  function addStringArray(Run memory self, string _key, string[] memory _values)
    internal pure
  {
    self.buf.encodeString(_key);
    self.buf.startArray();
    for (uint256 i = 0; i < _values.length; i++) {
      self.buf.encodeString(_values[i]);
    }
    self.buf.endSequence();
  }

  function close(Run memory self)
    internal pure
  {
     self.buf.endSequence();
  }
}
