pragma solidity 0.4.24;

// GetterSetter is a contract to aid debugging and testing during development.
contract GetterSetter {
  bytes32 public getBytes32;
  uint256 public getUint256;
  bytes32 public requestId;
  bytes public getBytes;

  event SetBytes32(address indexed from, bytes32 indexed value);
  event SetUint256(address indexed from, uint256 indexed value);
  event SetBytes(address indexed from, bytes value);

  event Output(bytes32 b32, uint256 u256, bytes32 b322);

  function setBytes32(bytes32 _value) public {
    getBytes32 = _value;
    emit SetBytes32(msg.sender, _value);
  }

  function requestedBytes32(bytes32 _requestId, bytes32 _value) public {
    requestId = _requestId;
    setBytes32(_value);
  }

  function setBytes(bytes _value) public {
    getBytes = _value;
    emit SetBytes(msg.sender, _value);
  }

  function requestedBytes(bytes32 _requestId, bytes _value) public {
    requestId = _requestId;
    setBytes(_value);
  }

  function setUint256(uint256 _value) public {
    getUint256 = _value;
    emit SetUint256(msg.sender, _value);
  }

  function requestedUint256(bytes32 _requestId, uint256 _value) public {
    requestId = _requestId;
    setUint256(_value);
  }
}
