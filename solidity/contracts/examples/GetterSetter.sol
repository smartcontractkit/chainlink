pragma solidity ^0.4.18;

// GetterSetter is a contract to aid debugging and testing during development.
contract GetterSetter {
  bytes32 public getBytes32;
  uint256 public getUint256;
  uint256 public requestId;

  function setBytes32(bytes32 _value) public {
    getBytes32 = _value;
  }

  function requestedBytes32(uint256 _requestId, bytes32 _value) public {
    requestId = _requestId;
    setBytes32(_value);
  }

  function setUint256(uint256 _value) public {
    getUint256 = _value;
  }

  function requestedUint256(uint256 _requestId, uint256 _value) public {
    requestId = _requestId;
    setUint256(_value);
  }
}
