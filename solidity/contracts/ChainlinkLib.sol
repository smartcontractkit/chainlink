pragma solidity ^0.4.18;

library ChainlinkLib {
  bytes constant stringType = "string,";
  bytes constant bytes32Type = "bytes32,";
  bytes constant bytes32ArrayType = "bytes32[],";

  struct Run {
    bytes32 id;
    bytes32 jobId;
    address callbackAddress;
    bytes4 callbackFunctionId;
    bytes names;
    bytes types;
    bytes values;
  }

  function add(Run memory self, string _key, string _value)
    internal
  {
    self.types = concat(self.types, stringType);
    self.names = concat(concat(self.names, bytes(_key)), ",");
    self.values = append(self.values, bytes(_value));
  }

  function addBytes32(Run self, string _key, bytes32 _value)
    internal
  {
    self.names = concat(self.names, concat(bytes(_key), ","));
    self.types = concat(self.types, bytes32Type);
    self.values = concat(self.values, toBytes(_value));
  }

  function addBytes32Array(Run self, string _key, bytes32[] memory _values)
    internal
  {
    self.names = concat(self.names, concat(bytes(_key), ","));
    self.types = concat(self.types, bytes32ArrayType);
    self.values = concat(self.values, toBytes(_values.length));
    for (uint256 i = 0; i < _values.length; i++) {
      self.values = concat(self.values, toBytes(_values[i]));
    }
  }

  function payload(Run self)
    internal
    returns (bytes)
  {
    bytes memory result = addLengthPrefix(self.names);
    return addLengthPrefix(append(append(result, self.types), self.values));
  }

  function toBytes(bytes32 _b)
    internal
    pure
    returns (bytes memory)
  {
    bytes memory c = new bytes(32);
    uint charCount = 0;
    for (uint i = 0; i < 32; i++) {
        c[i] = byte(bytes32(uint(_b) * 2 ** (8 * i)));
    }
    return c;
  }

  function toBytes(bytes4 _b)
    internal
    pure
    returns (bytes memory)
  {
    bytes memory c = new bytes(4);
    uint charCount = 0;
    for (uint i = 0; i < 4; i++) {
        c[i] = byte(bytes32(uint(_b) * 2 ** (8 * i)));
    }
    return c;
  }

  function toBytes(uint256 _val)
    internal
    pure
    returns (bytes memory)
  {
    bytes memory val = new bytes(32);
    assembly { mstore(add(val, 32), _val) }
    return val;
  }

  // https://ethereum.stackexchange.com/a/40456/24978
  function concat(bytes memory a, bytes memory b)
    internal
    pure
    returns (bytes c)
  {
      // Store the length of the first array
      uint256 alen = a.length;
      // Store the length of BOTH arrays
      uint256 totallen = alen + b.length;
      // Count the loops required for array a (sets of 32 bytes)
      uint256 loopsa = (a.length + 31) / 32;
      // Count the loops required for array b (sets of 32 bytes)
      uint256 loopsb = (b.length + 31) / 32;
      assembly {
          let m := mload(0x40)
          // Load the length of both arrays to the head of the new bytes array
          mstore(m, totallen)
          // Add the contents of a to the array
          for {  let i := 0 } lt(i, loopsa) { i := add(1, i) } { mstore(add(m, mul(32, add(1, i))), mload(add(a, mul(32, add(1, i))))) }
          // Add the contents of b to the array
          for {  let i := 0 } lt(i, loopsb) { i := add(1, i) } { mstore(add(m, add(mul(32, add(1, i)), alen)), mload(add(b, mul(32, add(1, i))))) }
          mstore(0x40, add(m, add(32, totallen)))
          c := m
      }
  }

  function append(bytes memory _a, bytes memory _b)
    internal
    pure
    returns (bytes memory c)
  {
    if (_a.length == 0) {
      return addLengthPrefix(_b);
    } else {
      return nestBytes(_a, _b);
    }
  }

  function nestBytes(bytes memory _a, bytes memory _b)
    internal
    pure
    returns (bytes memory c)
  {
      uint256 loopsA = ((_a.length + 31) / 32) * 32;
      uint256 loopsB = (((_b.length + 31) / 32) + 1) * 32;
      uint256 totalLen = loopsA + loopsB;
      assembly {
          let mem := mload(0x40)
          mstore(mem, totalLen)
          mem := add(32, mem)
          for {  let i := 0 } lt(i, loopsA) { i := add(32, i) } {
            mstore(add(mem, i), mload(add(_a, add(i, 32))))
          }
          for {  let i := 0 } lt(i, loopsB) { i := add(32, i) } {
            mstore(add(mem, add(i, loopsA)), mload(add(_b, i)))
          }
          mstore(0x40, add(mem, totalLen))
          c := sub(mem, 32)
      }
  }

  function addLengthPrefix(bytes memory _in)
    internal
    pure
    returns (bytes memory c)
  {
      uint256 totalLen = (((_in.length + 31) / 32) + 1) * 32;
      assembly {
          let mem := mload(0x40)
          mstore(mem, totalLen)
          mem := add(32, mem)
          for {  let i := 0 } lt(i, totalLen) { i := add(32, i) } {
            mstore(add(mem, i), mload(add(_in, i)))
          }
          mstore(0x40, add(mem, totalLen))
          c := sub(mem, 32)
      }
  }

}
