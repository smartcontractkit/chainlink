// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {LiquidityManagerBaseTest} from "../LiquidityManagerBaseTest.t.sol";

contract OCR3Setup is LiquidityManagerBaseTest {
  // Signer private keys used for these test
  uint256 internal constant PRIVATE0 = 0x7b2e97fe057e6de99d6872a2ef2abf52c9b4469bc848c2465ac3fcd8d336e81d;
  uint256 internal constant PRIVATE1 = 0xab56160806b05ef1796789248e1d7f34a6465c5280899159d645218cd216cee6;
  uint256 internal constant PRIVATE2 = 0x6ec7caa8406a49b76736602810e0a2871959fbbb675e23a8590839e4717f1f7f;
  uint256 internal constant PRIVATE3 = 0x80f14b11da94ae7f29d9a7713ea13dc838e31960a5c0f2baf45ed458947b730a;

  address[] internal s_valid_signers;
  address[] internal s_valid_transmitters;

  uint64 internal constant s_offchainConfigVersion = 3;
  uint8 internal constant s_f = 1;
  bytes internal constant REPORT = abi.encode("testReport");

  function setUp() public virtual override {
    LiquidityManagerBaseTest.setUp();

    s_valid_transmitters = new address[](4);
    for (uint160 i = 0; i < 4; ++i) {
      s_valid_transmitters[i] = address(4 + i);
    }

    s_valid_signers = new address[](4);
    s_valid_signers[0] = vm.addr(PRIVATE0); //0xc110458BE52CaA6bB68E66969C3218A4D9Db0211
    s_valid_signers[1] = vm.addr(PRIVATE1); //0xc110a19c08f1da7F5FfB281dc93630923F8E3719
    s_valid_signers[2] = vm.addr(PRIVATE2); //0xc110fdF6e8fD679C7Cc11602d1cd829211A18e9b
    s_valid_signers[3] = vm.addr(PRIVATE3); //0xc11028017c9b445B6bF8aE7da951B5cC28B326C0
  }
}
