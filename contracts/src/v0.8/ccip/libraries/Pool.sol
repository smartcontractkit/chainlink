// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library Pool {
  // bytes4(keccak256("POOL_RETURN_DATA_V1_TAG"))
  bytes4 public constant POOL_RETURN_DATA_V1_TAG = 0x179fa694;

  struct PoolReturnDataV1 {
    bytes destPoolAddress;
    bytes destPoolData;
  }

  function _generatePoolReturnDataV1(
    bytes memory remotePoolAddress,
    bytes memory destPoolData
  ) internal pure returns (bytes memory) {
    return abi.encode(PoolReturnDataV1({destPoolAddress: remotePoolAddress, destPoolData: destPoolData}));

    // TODO next PR: actually use the tag
    return abi.encodeWithSelector(
      POOL_RETURN_DATA_V1_TAG, PoolReturnDataV1({destPoolAddress: remotePoolAddress, destPoolData: destPoolData})
    );
  }
}
