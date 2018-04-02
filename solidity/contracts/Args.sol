pragma solidity ^0.4.18;

contract Args {
  bytes constant stringType = "string,";
  bytes constant bytes32Type = "bytes32,";
  bytes constant bytes32ArrayType = "bytes32[],";

  bytes types;
  bytes names;
  uint16[] lengths;
  bytes values; // all implied to be [disk] storage.

  // do we need to do lengths for every bytes array?
  event Data(
    bytes types,
    uint16[] lengths,
    bytes names,
    bytes values
  );

  function fireEvent()
    public
  {
    Data(types, lengths, names, values);
  }

  function add(string _key, string _value)
    public
  {
    types = concat(types, stringType);
    bytes memory value = bytes(_value);
    lengths.push(uint16(value.length));
    names = concat(names, bytes(_key));
    names = concat(names, ",");
    values = concat(values, value);
  }

  function addBytes32(string _key, bytes32 _value)
    public
  {
    types = concat(types, bytes32Type);
    lengths.push(32);
    names = concat(names, concat(bytes(_key), ","));
    values = concat(values, toBytes(_value));
  }

  function addBytes32Array(string _key, bytes32[] memory _values)
    public
  {
    types = concat(types, bytes32Type);
    lengths.push(uint16(_values.length));
    names = concat(names, concat(bytes(_key), ","));
    values = concat(values, toBytes(_values.length));
    for (uint256 i = 0; i < _values.length; i++) {
      values = concat(values, toBytes(_values[i]));
    }
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
    returns (bytes memory c)
  {
      // Store the length of the first array
      uint alen = a.length;
      // Store the length of BOTH arrays
      uint totallen = alen + b.length;
      // Count the loops required for array a (sets of 32 bytes)
      uint loopsa = (a.length + 31) / 32;
      // Count the loops required for array b (sets of 32 bytes)
      uint loopsb = (b.length + 31) / 32;
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
}
