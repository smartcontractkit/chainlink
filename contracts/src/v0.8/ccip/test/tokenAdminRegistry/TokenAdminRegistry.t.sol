// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {IPoolV1} from "../../interfaces/IPool.sol";

import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {TokenSetup} from "../TokenSetup.t.sol";

contract TokenAdminRegistrySetup is TokenSetup {
  address internal s_registryModule = makeAddr("registryModule");

  function setUp() public virtual override {
    TokenSetup.setUp();

    s_tokenAdminRegistry.addRegistryModule(s_registryModule);
  }
}

contract TokenAdminRegistry_getPools is TokenAdminRegistrySetup {
  function test_getPools_Success() public {
    address[] memory tokens = new address[](1);
    tokens[0] = s_sourceTokens[0];

    address[] memory got = s_tokenAdminRegistry.getPools(tokens);
    assertEq(got.length, 1);
    assertEq(got[0], s_sourcePoolByToken[tokens[0]]);

    got = s_tokenAdminRegistry.getPools(s_sourceTokens);
    assertEq(got.length, s_sourceTokens.length);
    for (uint256 i = 0; i < s_sourceTokens.length; i++) {
      assertEq(got[i], s_sourcePoolByToken[s_sourceTokens[i]]);
    }

    address doesNotExist = makeAddr("doesNotExist");
    tokens[0] = doesNotExist;
    got = s_tokenAdminRegistry.getPools(tokens);
    assertEq(got.length, 1);
    assertEq(got[0], address(0));
  }
}

contract TokenAdminRegistry_getPool is TokenAdminRegistrySetup {
  function test_getPool_Success() public view {
    address got = s_tokenAdminRegistry.getPool(s_sourceTokens[0]);
    assertEq(got, s_sourcePoolByToken[s_sourceTokens[0]]);
  }
}

contract TokenAdminRegistry_setPool is TokenAdminRegistrySetup {
  function test_setPool_Success() public {
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

  function test_setPool_ZeroAddressRemovesPool_Success() public {
    address pool = makeAddr("pool");
    vm.mockCall(pool, abi.encodeWithSelector(IPoolV1.isSupportedToken.selector), abi.encode(true));
    s_tokenAdminRegistry.setPool(s_sourceTokens[0], pool);

    assertEq(s_tokenAdminRegistry.getPool(s_sourceTokens[0]), pool);

    vm.expectEmit();
    emit TokenAdminRegistry.PoolSet(s_sourceTokens[0], pool, address(0));

    s_tokenAdminRegistry.setPool(s_sourceTokens[0], address(0));

    assertEq(s_tokenAdminRegistry.getPool(s_sourceTokens[0]), address(0));
  }

  function test_setPool_InvalidTokenPoolToken_Revert() public {
    address pool = makeAddr("pool");
    vm.mockCall(pool, abi.encodeWithSelector(IPoolV1.isSupportedToken.selector), abi.encode(false));

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.InvalidTokenPoolToken.selector, s_sourceTokens[0]));
    s_tokenAdminRegistry.setPool(s_sourceTokens[0], pool);
  }

  function test_setPool_OnlyAdministrator_Revert() public {
    vm.stopPrank();

    vm.expectRevert(
      abi.encodeWithSelector(TokenAdminRegistry.OnlyAdministrator.selector, address(this), s_sourceTokens[0])
    );
    s_tokenAdminRegistry.setPool(s_sourceTokens[0], makeAddr("pool"));
  }
}

contract TokenAdminRegistry_getAllConfiguredTokens is TokenAdminRegistrySetup {
  function test_Fuzz_getAllConfiguredTokens_Success(uint8 numberOfTokens) public {
    TokenAdminRegistry cleanTokenAdminRegistry = new TokenAdminRegistry();
    for (uint160 i = 0; i < numberOfTokens; ++i) {
      cleanTokenAdminRegistry.proposeAdministrator(address(i), address(i + 1000));
    }

    uint160 count = 0;
    for (uint160 start = 0; start < numberOfTokens; start += count++) {
      address[] memory got = cleanTokenAdminRegistry.getAllConfiguredTokens(uint64(start), uint64(count));
      if (start + count > numberOfTokens) {
        assertEq(got.length, numberOfTokens - start);
      } else {
        assertEq(got.length, count);
      }

      for (uint160 j = 0; j < got.length; ++j) {
        assertEq(got[j], address(j + start));
      }
    }
  }

  function test_getAllConfiguredTokens_outOfBounds_Success() public view {
    address[] memory tokens = s_tokenAdminRegistry.getAllConfiguredTokens(type(uint64).max, 10);
    assertEq(tokens.length, 0);
  }
}

contract TokenAdminRegistry_transferAdminRole is TokenAdminRegistrySetup {
  function test_transferAdminRole_Success() public {
    address token = s_sourceTokens[0];

    address currentAdmin = s_tokenAdminRegistry.getTokenConfig(token).administrator;
    address newAdmin = makeAddr("newAdmin");

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferRequested(token, currentAdmin, newAdmin);

    s_tokenAdminRegistry.transferAdminRole(token, newAdmin);

    TokenAdminRegistry.TokenConfig memory config = s_tokenAdminRegistry.getTokenConfig(token);

    // Assert only the pending admin updates, without affecting the pending admin.
    assertEq(config.pendingAdministrator, newAdmin);
    assertEq(config.administrator, currentAdmin);
  }

  function test_transferAdminRole_OnlyAdministrator_Revert() public {
    vm.stopPrank();

    vm.expectRevert(
      abi.encodeWithSelector(TokenAdminRegistry.OnlyAdministrator.selector, address(this), s_sourceTokens[0])
    );
    s_tokenAdminRegistry.transferAdminRole(s_sourceTokens[0], makeAddr("newAdmin"));
  }
}

contract TokenAdminRegistry_acceptAdminRole is TokenAdminRegistrySetup {
  function test_acceptAdminRole_Success() public {
    address token = s_sourceTokens[0];

    address currentAdmin = s_tokenAdminRegistry.getTokenConfig(token).administrator;
    address newAdmin = makeAddr("newAdmin");

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferRequested(token, currentAdmin, newAdmin);

    s_tokenAdminRegistry.transferAdminRole(token, newAdmin);

    TokenAdminRegistry.TokenConfig memory config = s_tokenAdminRegistry.getTokenConfig(token);

    // Assert only the pending admin updates, without affecting the pending admin.
    assertEq(config.pendingAdministrator, newAdmin);
    assertEq(config.administrator, currentAdmin);

    vm.startPrank(newAdmin);

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferred(token, newAdmin);

    s_tokenAdminRegistry.acceptAdminRole(token);

    config = s_tokenAdminRegistry.getTokenConfig(token);

    // Assert only the pending admin updates, without affecting the pending admin.
    assertEq(config.pendingAdministrator, address(0));
    assertEq(config.administrator, newAdmin);
  }

  function test_acceptAdminRole_OnlyPendingAdministrator_Revert() public {
    address token = s_sourceTokens[0];
    address currentAdmin = s_tokenAdminRegistry.getTokenConfig(token).administrator;
    address newAdmin = makeAddr("newAdmin");

    s_tokenAdminRegistry.transferAdminRole(token, newAdmin);

    TokenAdminRegistry.TokenConfig memory config = s_tokenAdminRegistry.getTokenConfig(token);

    // Assert only the pending admin updates, without affecting the pending admin.
    assertEq(config.pendingAdministrator, newAdmin);
    assertEq(config.administrator, currentAdmin);

    address notNewAdmin = makeAddr("notNewAdmin");
    vm.startPrank(notNewAdmin);

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.OnlyPendingAdministrator.selector, notNewAdmin, token));
    s_tokenAdminRegistry.acceptAdminRole(token);
  }
}

contract TokenAdminRegistry_isAdministrator is TokenAdminRegistrySetup {
  function test_isAdministrator_Success() public {
    address newAdmin = makeAddr("newAdmin");
    address newToken = makeAddr("newToken");
    assertFalse(s_tokenAdminRegistry.isAdministrator(newToken, newAdmin));
    assertFalse(s_tokenAdminRegistry.isAdministrator(newToken, OWNER));

    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);
    changePrank(newAdmin);
    s_tokenAdminRegistry.acceptAdminRole(newToken);

    assertTrue(s_tokenAdminRegistry.isAdministrator(newToken, newAdmin));
    assertFalse(s_tokenAdminRegistry.isAdministrator(newToken, OWNER));
  }
}

contract TokenAdminRegistry_proposeAdministrator is TokenAdminRegistrySetup {
  function test_proposeAdministrator_module_Success() public {
    vm.startPrank(s_registryModule);
    address newAdmin = makeAddr("newAdmin");
    address newToken = makeAddr("newToken");

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferRequested(newToken, address(0), newAdmin);

    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);

    assertEq(s_tokenAdminRegistry.getTokenConfig(newToken).pendingAdministrator, newAdmin);
    assertEq(s_tokenAdminRegistry.getTokenConfig(newToken).administrator, address(0));
    assertEq(s_tokenAdminRegistry.getTokenConfig(newToken).tokenPool, address(0));

    changePrank(newAdmin);
    s_tokenAdminRegistry.acceptAdminRole(newToken);

    assertTrue(s_tokenAdminRegistry.isAdministrator(newToken, newAdmin));
  }

  function test_proposeAdministrator_owner_Success() public {
    address newAdmin = makeAddr("newAdmin");
    address newToken = makeAddr("newToken");

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferRequested(newToken, address(0), newAdmin);

    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);

    assertEq(s_tokenAdminRegistry.getTokenConfig(newToken).pendingAdministrator, newAdmin);

    changePrank(newAdmin);
    s_tokenAdminRegistry.acceptAdminRole(newToken);

    assertTrue(s_tokenAdminRegistry.isAdministrator(newToken, newAdmin));
  }

  function test_proposeAdministrator_reRegisterWhileUnclaimed_Success() public {
    address newAdmin = makeAddr("wrongAddress");
    address newToken = makeAddr("newToken");

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferRequested(newToken, address(0), newAdmin);

    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);

    assertEq(s_tokenAdminRegistry.getTokenConfig(newToken).pendingAdministrator, newAdmin);

    newAdmin = makeAddr("correctAddress");

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferRequested(newToken, address(0), newAdmin);

    // Ensure we can still register the correct admin while the previous admin is unclaimed.
    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);

    changePrank(newAdmin);
    s_tokenAdminRegistry.acceptAdminRole(newToken);

    assertTrue(s_tokenAdminRegistry.isAdministrator(newToken, newAdmin));
  }

  mapping(address token => address admin) internal s_AdminByToken;

  function test_Fuzz_proposeAdministrator_Success(address[50] memory tokens, address[50] memory admins) public {
    TokenAdminRegistry cleanTokenAdminRegistry = new TokenAdminRegistry();
    for (uint256 i = 0; i < tokens.length; i++) {
      if (admins[i] == address(0)) {
        continue;
      }
      if (cleanTokenAdminRegistry.getTokenConfig(tokens[i]).administrator != address(0)) {
        continue;
      }
      cleanTokenAdminRegistry.proposeAdministrator(tokens[i], admins[i]);
      s_AdminByToken[tokens[i]] = admins[i];
    }

    for (uint256 i = 0; i < tokens.length; i++) {
      assertEq(cleanTokenAdminRegistry.getTokenConfig(tokens[i]).pendingAdministrator, s_AdminByToken[tokens[i]]);
    }
  }

  function test_proposeAdministrator_OnlyRegistryModule_Revert() public {
    address newToken = makeAddr("newToken");
    vm.stopPrank();

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.OnlyRegistryModuleOrOwner.selector, address(this)));
    s_tokenAdminRegistry.proposeAdministrator(newToken, OWNER);
  }

  function test_proposeAdministrator_ZeroAddress_Revert() public {
    address newToken = makeAddr("newToken");

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.ZeroAddress.selector));
    s_tokenAdminRegistry.proposeAdministrator(newToken, address(0));
  }

  function test_proposeAdministrator_AlreadyRegistered_Revert() public {
    address newAdmin = makeAddr("newAdmin");
    address newToken = makeAddr("newToken");

    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);
    changePrank(newAdmin);
    s_tokenAdminRegistry.acceptAdminRole(newToken);

    changePrank(OWNER);

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.AlreadyRegistered.selector, newToken));
    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);
  }
}

contract TokenAdminRegistry_addRegistryModule is TokenAdminRegistrySetup {
  function test_addRegistryModule_Success() public {
    address newModule = makeAddr("newModule");

    s_tokenAdminRegistry.addRegistryModule(newModule);

    assertTrue(s_tokenAdminRegistry.isRegistryModule(newModule));

    // Assert the event is not emitted if the module is already added.
    vm.recordLogs();
    s_tokenAdminRegistry.addRegistryModule(newModule);

    vm.assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_addRegistryModule_OnlyOwner_Revert() public {
    address newModule = makeAddr("newModule");
    vm.stopPrank();

    vm.expectRevert("Only callable by owner");
    s_tokenAdminRegistry.addRegistryModule(newModule);
  }
}

contract TokenAdminRegistry_removeRegistryModule is TokenAdminRegistrySetup {
  function test_removeRegistryModule_Success() public {
    address newModule = makeAddr("newModule");

    s_tokenAdminRegistry.addRegistryModule(newModule);

    assertTrue(s_tokenAdminRegistry.isRegistryModule(newModule));

    vm.expectEmit();
    emit TokenAdminRegistry.RegistryModuleRemoved(newModule);

    s_tokenAdminRegistry.removeRegistryModule(newModule);

    assertFalse(s_tokenAdminRegistry.isRegistryModule(newModule));

    // Assert the event is not emitted if the module is already removed.
    vm.recordLogs();
    s_tokenAdminRegistry.removeRegistryModule(newModule);

    vm.assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_removeRegistryModule_OnlyOwner_Revert() public {
    address newModule = makeAddr("newModule");
    vm.stopPrank();

    vm.expectRevert("Only callable by owner");
    s_tokenAdminRegistry.removeRegistryModule(newModule);
  }
}
