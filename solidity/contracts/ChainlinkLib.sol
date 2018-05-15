pragma solidity ^0.4.23;

import "solidity-cborutils/contracts/CBOR.sol";

library ChainlinkLib {
  bytes4 internal constant oracleRequestDataFid = bytes4(keccak256("requestData(uint256,bytes32,address,bytes4,bytes32,bytes)"));
  bytes4 internal constant oracleSpecAndRunFid = bytes4(keccak256("specAndRun(uint256,address,bytes4,bytes32,bytes)"));

  using CBOR for Buffer.buffer;

  struct Run {
    bytes32 jobId;
    address callbackAddress;
    bytes4 callbackFunctionId;
    bytes32 requestId;
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
    uint256 _clArgsVersion
  ) internal pure returns (bytes memory) {
    return abi.encodeWithSelector(
      oracleRequestDataFid,
      _clArgsVersion,
      self.jobId,
      self.callbackAddress,
      self.callbackFunctionId,
      self.requestId,
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

  function close(Run memory self) internal pure {
    self.buf.endSequence();
  }

  struct Job {
    address callbackAddress;
    bytes4 callbackFunctionId;
    bytes32 requestId;
    Buffer.buffer buf;
  }

  function initialize(
    Job memory self,
    string[] _tasks,
    address _callbackAddress,
    string _callbackFunctionSignature
  ) internal pure returns (ChainlinkLib.Job memory) {
    Buffer.init(self.buf, 128);
    self.callbackAddress = _callbackAddress;
    self.callbackFunctionId = bytes4(keccak256(_callbackFunctionSignature));

    self.buf.startMap();
    self.buf.encodeString("tasks");
    self.buf.startArray();
    for (uint256 i = 0; i < _tasks.length; i++) {
      self.buf.encodeString(_tasks[i]);
    }
    self.buf.endSequence();
    self.buf.encodeString("params");
    self.buf.startMap();
    return self;
  }

  function encodeForOracle(
    Job memory self,
    uint256 _clArgsVersion
  ) internal pure returns (bytes memory) {
    return abi.encodeWithSelector(
      oracleSpecAndRunFid,
      _clArgsVersion,
      self.callbackAddress,
      self.callbackFunctionId,
      self.requestId,
      self.buf.buf);
  }

  function add(Job memory self, string _key, string _value)
    internal pure
  {
    self.buf.encodeString(_key);
    self.buf.encodeString(_value);
  }

  function addStringArray(Job memory self, string _key, string[] memory _values)
    internal pure
  {
    self.buf.encodeString(_key);
    self.buf.startArray();
    for (uint256 i = 0; i < _values.length; i++) {
      self.buf.encodeString(_values[i]);
    }
    self.buf.endSequence();
  }

  function close(Job memory self) internal pure {
    self.buf.endSequence();
    self.buf.endSequence();
  }
}
