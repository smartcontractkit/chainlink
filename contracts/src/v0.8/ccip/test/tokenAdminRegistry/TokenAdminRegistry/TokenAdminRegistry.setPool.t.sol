// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {IPoolV1} from "../../../interfaces/IPool.sol";
import {TokenAdminRegistry} from "../../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {TokenAdminRegistrySetup} from "./TokenAdminRegistrySetup.t.sol";

contract TokenAdminRegistry_setPool is TokenAdminRegistrySetup {
  function test_setPool() public {
    address pool = makeAddr("pool");
    vm.mockCall(pool, abi.encodeWithSelector(IPoolV1.isSupportedToken.selector), abi.encode(true));

    vm.expectEmit();
    emit TokenAdminRegistry.PoolSet(s_sourceTokens[0], s_sourcePoolByToken[s_sourceTokens[0]], pool);

    s_tokenAdminRegistry.setPool(s_sourceTokens[0], pool);

    assertEq(s_tokenAdminRegistry.getPool(s_sourceTokens[0]), pool);

    // Assert the event is not emitted if the pool is the same as the current pool.
    vm.recordLogs();
    s_tokenAdminRegistry.setPool(s_sourceTokens[0], pool);

    vm.assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_setPool_ZeroAddressRemovesPool() public {
    address pool = makeAddr("pool");
    vm.mockCall(pool, abi.encodeWithSelector(IPoolV1.isSupportedToken.selector), abi.encode(true));
    s_tokenAdminRegistry.setPool(s_sourceTokens[0], pool);

    assertEq(s_tokenAdminRegistry.getPool(s_sourceTokens[0]), pool);

    vm.expectEmit();
    emit TokenAdminRegistry.PoolSet(s_sourceTokens[0], pool, address(0));

    s_tokenAdminRegistry.setPool(s_sourceTokens[0], address(0));

    assertEq(s_tokenAdminRegistry.getPool(s_sourceTokens[0]), address(0));
  }

  function test_RevertWhen_setPool_InvalidTokenPoolToken() public {
    address pool = makeAddr("pool");
    vm.mockCall(pool, abi.encodeWithSelector(IPoolV1.isSupportedToken.selector), abi.encode(false));

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.InvalidTokenPoolToken.selector, s_sourceTokens[0]));
    s_tokenAdminRegistry.setPool(s_sourceTokens[0], pool);
  }

  function test_RevertWhen_setPool_OnlyAdministrator() public {
    vm.stopPrank();

    vm.expectRevert(
      abi.encodeWithSelector(TokenAdminRegistry.OnlyAdministrator.selector, address(this), s_sourceTokens[0])
    );
    s_tokenAdminRegistry.setPool(s_sourceTokens[0], makeAddr("pool"));
  }
}
