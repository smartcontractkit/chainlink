pragma solidity ^0.4.18;

contract BytesUtils {

  function bytes4toBytes(bytes4 _b)
    internal
    pure
    returns (bytes memory)
  {
    bytes memory c = new bytes(4);
    uint charCount = 0;
    for (uint i = 0; i < 4; i++) {
        c[i] = byte(bytes32(uint(_b) * 2 ** (8 * (28 + i))));
    }
    return c;
  }

  function bytes32toBytes(bytes32 _b)
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

  function uint256toBytes(uint256 _val)
    internal
    pure
    returns (bytes memory)
  {
    bytes memory val = new bytes(32);
    assembly { mstore(add(val, 32), _val) }
    return val;
  }

  function addressToBytes(address a)
    internal
    pure
    returns (bytes memory)
  {
    bytes memory m = new bytes(32);
    assembly { mstore(add(m, 32), a) }
    return m;
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

  function append(bytes memory _a, bytes memory _b)
    internal
    pure
    returns (bytes memory c)
  {
    uint256 lenA = _a.length;
    uint256 loopsA = ((_a.length + 31) / 32) * 32;
    uint256 loopsB = ((_b.length + 31) / 32) * 32;
    uint256 totalLen = _a.length + loopsB;
    assembly {
      let mem := mload(0x40)
      mstore(mem, totalLen)
      mem := add(32, mem)
      for {  let i := 0 } lt(i, loopsA) { i := add(32, i) } {
        mstore(add(mem, i), mload(add(_a, add(i, 32))))
      }
      for {  let i := 0 } lt(i, loopsB) { i := add(32, i) } {
        mstore(add(mem, add(i, lenA)), mload(add(add(_b, i), 32)))
      }
      mstore(0x40, add(mem, totalLen))
      c := sub(mem, 32)
    }
  }
}
