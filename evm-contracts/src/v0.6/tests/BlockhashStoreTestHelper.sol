pragma solidity 0.6.2;

import "../dev/BlockhashStore.sol";

contract BlockhashStoreTestHelper is BlockhashStore {
  function godmodeSetHash(uint256 n, bytes32 h) public {
    s_blockhashes[n] = h;
  }
}
