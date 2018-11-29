pragma solidity 0.4.24;

import "solidity-cborutils/contracts/CBOR.sol";

library ChainlinkLib {
  bytes internal constant reqSig = "requestData(address,uint256,uint256,bytes32,address,bytes4,bytes32,bytes)";
  bytes4 internal constant oracleRequestDataFid = bytes4(keccak256(reqSig));
  bytes internal constant coordSig = "executeServiceAgreement(address,uint256,uint256,bytes32,address,bytes4,bytes32,bytes)";
  bytes4 internal constant coordinatorRequestDataFid = bytes4(keccak256(coordSig));

  using CBOR for Buffer.buffer;

  struct Run {
    bytes32 specId;
    address callbackAddress;
    bytes4 callbackFunctionId;
    bytes32 requestId;
    Buffer.buffer buf;
  }

  function initialize(
    Run memory self,
    bytes32 _specId,
    address _callbackAddress,
    bytes4 _callbackFunction
  ) internal pure returns (ChainlinkLib.Run memory) {
    Buffer.init(self.buf, 128);
    self.specId = _specId;
    self.callbackAddress = _callbackAddress;
    self.callbackFunctionId = _callbackFunction;
    self.buf.startMap();
    return self;
  }

  function encodeForOracle(
    Run memory self,
    uint256 _clArgsVersion
  ) internal pure returns (bytes memory) {
    return abi.encodeWithSelector(
      oracleRequestDataFid,
      0, // overridden by onTokenTransfer
      0, // overridden by onTokenTransfer
      _clArgsVersion,
      self.specId,
      self.callbackAddress,
      self.callbackFunctionId,
      self.requestId,
      self.buf.buf);
  }

  function encodeForCoordinator(
    Run memory self,
    uint256 _clArgsVersion
  ) internal pure returns (bytes memory) {
    return abi.encodeWithSelector(
      coordinatorRequestDataFid, // need to change
      0, // overridden by onTokenTransfer
      0, // overridden by onTokenTransfer
      _clArgsVersion,
      self.specId,
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

  function addBytes(Run memory self, string _key, bytes _value)
    internal pure
  {
    self.buf.encodeString(_key);
    self.buf.encodeBytes(_value);
  }

  function addInt(Run memory self, string _key, int256 _value)
    internal pure
  {
    self.buf.encodeString(_key);
    self.buf.encodeInt(_value);
  }

  function addUint(Run memory self, string _key, uint256 _value)
    internal pure
  {
    self.buf.encodeString(_key);
    self.buf.encodeUInt(_value);
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
}
