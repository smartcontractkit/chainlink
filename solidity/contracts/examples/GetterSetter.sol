pragma solidity ^0.4.18;

// GetterSetter is a contract to aid debugging and testing during development.
contract GetterSetter {
  bytes32 public getBytes32;
  uint256 public getUint256;
  uint256 public requestId;

  event SetBytes32(address indexed from, bytes32 indexed value);
  event SetUint256(address indexed from, uint256 indexed value);

  function setBytes32(bytes32 _value) public {
    getBytes32 = _value;
    SetBytes32(msg.sender, _value);
  }

  function requestedBytes32(uint256 _requestId, bytes32 _value) public {
    requestId = _requestId;
    setBytes32(_value);
  }

  function setUint256(uint256 _value) public {
    getUint256 = _value;
    SetUint256(msg.sender, _value);
  }

  function requestedUint256(uint256 _requestId, uint256 _value) public {
    requestId = _requestId;
    setUint256(_value);
  }
}
