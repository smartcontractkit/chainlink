// SPDX-License-Identifier: MIT
pragma solidity 0.6.6;

import "../BlockhashStore.sol";

contract BlockhashStoreTestHelper is BlockhashStore {
  function godmodeSetHash(uint256 n, bytes32 h) public {
    s_blockhashes[n] = h;
  }
}
