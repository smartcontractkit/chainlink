// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {BurnMintERC677} from "../../../../../shared/token/ERC677/BurnMintERC677.sol";

import {HybridLockReleaseUSDCTokenPool} from "../../../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {USDCTokenPool} from "../../../../pools/USDC/USDCTokenPool.sol";
import {USDCSetup} from "../USDCSetup.t.sol";

contract HybridLockReleaseUSDCTokenPoolSetup is USDCSetup {
  HybridLockReleaseUSDCTokenPool internal s_usdcTokenPool;
  HybridLockReleaseUSDCTokenPool internal s_usdcTokenPoolTransferLiquidity;
  address[] internal s_allowedList;

  function setUp() public virtual override {
    super.setUp();

    s_usdcTokenPool = new HybridLockReleaseUSDCTokenPool(
      s_mockUSDC, s_token, new address[](0), address(s_mockRMNRemote), address(s_router)
    );

    s_usdcTokenPoolTransferLiquidity = new HybridLockReleaseUSDCTokenPool(
      s_mockUSDC, s_token, new address[](0), address(s_mockRMNRemote), address(s_router)
    );

    BurnMintERC677(address(s_token)).grantMintAndBurnRoles(address(s_usdcTokenPool));

    _poolApplyChainUpdates(address(s_usdcTokenPool));

    USDCTokenPool.DomainUpdate[] memory domains = new USDCTokenPool.DomainUpdate[](1);
    domains[0] = USDCTokenPool.DomainUpdate({
      destChainSelector: DEST_CHAIN_SELECTOR,
      domainIdentifier: 9999,
      allowedCaller: keccak256("allowedCaller"),
      enabled: true
    });

    s_usdcTokenPool.setDomains(domains);

    s_usdcTokenPool.setLiquidityProvider(DEST_CHAIN_SELECTOR, OWNER);
    s_usdcTokenPool.setLiquidityProvider(SOURCE_CHAIN_SELECTOR, OWNER);
  }
}
