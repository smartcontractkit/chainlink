// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {USDCTokenPool} from "../../../../pools/USDC/USDCTokenPool.sol";
import {USDCTokenPoolHelper} from "../../../helpers/USDCTokenPoolHelper.sol";
import {USDCSetup} from "../USDCSetup.t.sol";

contract USDCTokenPoolSetup is USDCSetup {
  USDCTokenPoolHelper internal s_usdcTokenPool;
  USDCTokenPoolHelper internal s_usdcTokenPoolWithAllowList;
  address[] internal s_allowedList;

  function setUp() public virtual override {
    super.setUp();

    s_usdcTokenPool =
      new USDCTokenPoolHelper(s_mockUSDC, s_token, new address[](0), address(s_mockRMNRemote), address(s_router));

    s_allowedList.push(vm.randomAddress());
    s_usdcTokenPoolWithAllowList =
      new USDCTokenPoolHelper(s_mockUSDC, s_token, s_allowedList, address(s_mockRMNRemote), address(s_router));

    _poolApplyChainUpdates(address(s_usdcTokenPool));
    _poolApplyChainUpdates(address(s_usdcTokenPoolWithAllowList));

    USDCTokenPool.DomainUpdate[] memory domains = new USDCTokenPool.DomainUpdate[](1);
    domains[0] = USDCTokenPool.DomainUpdate({
      destChainSelector: DEST_CHAIN_SELECTOR,
      domainIdentifier: 9999,
      allowedCaller: keccak256("allowedCaller"),
      enabled: true
    });

    s_usdcTokenPool.setDomains(domains);
    s_usdcTokenPoolWithAllowList.setDomains(domains);
  }
}
