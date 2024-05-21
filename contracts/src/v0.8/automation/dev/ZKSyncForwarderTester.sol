// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

contract ZKSyncForwarderTester {
  event NumberSet(uint256 newNum, uint256 oldNum);
  event BytesSet(bytes newBytes, bytes oldBytes);
  event BoolSet(bool newBool, bool oldBool);
  event IndexedFields1(uint256 indexed id, bool b);
  event IndexedFields2(bool b, uint256 price);

  uint256 public a;
  bytes public b;
  bool public c;

  function getBlock() external view returns (uint256) {
    return block.number;
  }

  function getBlockHash(uint256 n) external view returns (bytes32 hash) {
    if (n >= block.number || block.number - n > 256) {
      return "";
    }
    return blockhash(n);
  }

  function setNumber(uint256 _a) external {
    emit NumberSet(_a, a);
    a = _a;
  }

  function setBytes(bytes calldata _b) external {
    emit BytesSet(_b, b);
    b = _b;
  }

  function setBool(bool _c) external {
    emit BoolSet(_c, c);
    c = _c;
  }

  function emitEvents(bool b) external {
    emit IndexedFields1(1, b);
    emit IndexedFields2(b, 100);
  }
}
