pragma solidity 0.4.24;

// GetterSetter is a contract to aid debugging and testing during development.
contract GetterSetter {
  bytes public getBytes;

  event SetBytes(address indexed from, bytes value);

  function setBytes(bytes32 _requestID, bytes _value) public {
    getBytes = _value;
    emit SetBytes(msg.sender, _value);
  }
}
