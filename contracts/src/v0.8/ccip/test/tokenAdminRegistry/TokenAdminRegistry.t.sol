// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {IPool} from "../../interfaces/IPool.sol";

import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {TokenSetup} from "../TokenSetup.t.sol";

contract TokenAdminRegistrySetup is TokenSetup {
  event AdministratorRegistered(address indexed token, address indexed administrator);
  event PoolSet(address indexed token, address indexed previousPool, address indexed newPool);
  event AdministratorTransferRequested(address indexed token, address indexed currentAdmin, address indexed newAdmin);
  event AdministratorTransferred(address indexed token, address indexed newAdmin);
  event DisableReRegistrationSet(address indexed token, bool disabled);
  event RegistryModuleAdded(address indexed module);
  event RegistryModuleRemoved(address indexed module);

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

contract TokenAdminRegistry_isTokenSupportedOnRemoteChain is TokenAdminRegistrySetup {
  function test_isTokenSupportedOnRemoteChain_Success() public {
    uint64 nonExistentChainSelector = 2523356;

    assertTrue(s_tokenAdminRegistry.isTokenSupportedOnRemoteChain(s_sourceTokens[0], DEST_CHAIN_SELECTOR));
    assertFalse(s_tokenAdminRegistry.isTokenSupportedOnRemoteChain(s_sourceTokens[0], nonExistentChainSelector));

    address nonExistentToken = makeAddr("nonExistentToken");

    assertFalse(s_tokenAdminRegistry.isTokenSupportedOnRemoteChain(nonExistentToken, DEST_CHAIN_SELECTOR));
  }
}

contract TokenAdminRegistry_setPool is TokenAdminRegistrySetup {
  function test_setPool_Success() public {
    address pool = makeAddr("pool");
    vm.mockCall(pool, abi.encodeWithSelector(IPool.isSupportedToken.selector), abi.encode(true));

    vm.expectEmit();
    emit PoolSet(s_sourceTokens[0], s_sourcePoolByToken[s_sourceTokens[0]], pool);

    s_tokenAdminRegistry.setPool(s_sourceTokens[0], pool);

    assertEq(s_tokenAdminRegistry.getPool(s_sourceTokens[0]), pool);

    // Assert the event is not emitted if the pool is the same as the current pool.
    vm.recordLogs();
    s_tokenAdminRegistry.setPool(s_sourceTokens[0], pool);

    vm.assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_setPool_ZeroAddressRemovesPool_Success() public {
    address pool = makeAddr("pool");
    vm.mockCall(pool, abi.encodeWithSelector(IPool.isSupportedToken.selector), abi.encode(true));
    s_tokenAdminRegistry.setPool(s_sourceTokens[0], pool);

    assertEq(s_tokenAdminRegistry.getPool(s_sourceTokens[0]), pool);

    vm.expectEmit();
    emit PoolSet(s_sourceTokens[0], pool, address(0));

    s_tokenAdminRegistry.setPool(s_sourceTokens[0], address(0));

    assertEq(s_tokenAdminRegistry.getPool(s_sourceTokens[0]), address(0));
  }

  function test_setPool_InvalidTokenPoolToken_Revert() public {
    address pool = makeAddr("pool");
    vm.mockCall(pool, abi.encodeWithSelector(IPool.isSupportedToken.selector), abi.encode(false));

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
      cleanTokenAdminRegistry.registerAdministratorPermissioned(address(i), address(i + 1000));
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
    emit AdministratorTransferRequested(token, currentAdmin, newAdmin);

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
    emit AdministratorTransferRequested(token, currentAdmin, newAdmin);

    s_tokenAdminRegistry.transferAdminRole(token, newAdmin);

    TokenAdminRegistry.TokenConfig memory config = s_tokenAdminRegistry.getTokenConfig(token);

    // Assert only the pending admin updates, without affecting the pending admin.
    assertEq(config.pendingAdministrator, newAdmin);
    assertEq(config.administrator, currentAdmin);

    vm.startPrank(newAdmin);

    vm.expectEmit();
    emit AdministratorTransferred(token, newAdmin);

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

contract TokenAdminRegistry_setDisableReRegistration is TokenAdminRegistrySetup {
  function test_setDisableReRegistration_Success() public {
    vm.expectEmit();
    emit DisableReRegistrationSet(s_sourceTokens[0], true);

    s_tokenAdminRegistry.setDisableReRegistration(s_sourceTokens[0], true);

    assertTrue(s_tokenAdminRegistry.getTokenConfig(s_sourceTokens[0]).disableReRegistration);

    vm.expectEmit();
    emit DisableReRegistrationSet(s_sourceTokens[0], false);

    s_tokenAdminRegistry.setDisableReRegistration(s_sourceTokens[0], false);

    assertFalse(s_tokenAdminRegistry.getTokenConfig(s_sourceTokens[0]).disableReRegistration);
  }
}

contract TokenAdminRegistry_isAdministrator is TokenAdminRegistrySetup {
  function test_isAdministrator_Success() public {
    address newOwner = makeAddr("newOwner");
    address newToken = makeAddr("newToken");
    assertFalse(s_tokenAdminRegistry.isAdministrator(newToken, newOwner));
    assertFalse(s_tokenAdminRegistry.isAdministrator(newToken, OWNER));

    s_tokenAdminRegistry.registerAdministratorPermissioned(newToken, newOwner);

    assertTrue(s_tokenAdminRegistry.isAdministrator(newToken, newOwner));
    assertFalse(s_tokenAdminRegistry.isAdministrator(newToken, OWNER));
  }
}

contract TokenAdminRegistry_registerAdministrator is TokenAdminRegistrySetup {
  function test_registerAdministrator_Success() public {
    vm.startPrank(s_registryModule);
    address newOwner = makeAddr("newOwner");
    address newToken = makeAddr("newToken");

    vm.expectEmit();
    emit AdministratorRegistered(newToken, newOwner);

    s_tokenAdminRegistry.registerAdministrator(newToken, newOwner);

    assertTrue(s_tokenAdminRegistry.isAdministrator(newToken, newOwner));
  }

  function test_registerAdministrator__disableReRegistration_Revert() public {
    vm.startPrank(s_registryModule);
    address newOwner = makeAddr("newOwner");
    address newToken = makeAddr("newToken");

    s_tokenAdminRegistry.registerAdministrator(newToken, newOwner);

    vm.startPrank(newOwner);

    s_tokenAdminRegistry.setDisableReRegistration(newToken, true);

    vm.startPrank(s_registryModule);
    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.AlreadyRegistered.selector, newToken));

    s_tokenAdminRegistry.registerAdministrator(newToken, newOwner);
  }

  function test_registerAdministrator_OnlyRegistryModule_Revert() public {
    address newToken = makeAddr("newToken");
    vm.stopPrank();

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.OnlyRegistryModule.selector, address(this)));
    s_tokenAdminRegistry.registerAdministrator(newToken, OWNER);
  }
}

contract TokenAdminRegistry_registerAdministratorPermissioned is TokenAdminRegistrySetup {
  function test_registerAdministratorPermissioned_Success() public {
    address newAdmin = makeAddr("newAdmin");
    address newToken = makeAddr("newToken");

    vm.expectEmit();
    emit AdministratorRegistered(newToken, newAdmin);

    s_tokenAdminRegistry.registerAdministratorPermissioned(newToken, newAdmin);

    assertTrue(s_tokenAdminRegistry.isAdministrator(newToken, newAdmin));
    assertEq(s_tokenAdminRegistry.getTokenConfig(newToken).isRegistered, true);
  }

  mapping(address token => address admin) internal s_AdminByToken;

  function test_Fuzz_registerAdministratorPermissioned_Success(
    address[50] memory tokens,
    address[50] memory admins
  ) public {
    TokenAdminRegistry cleanTokenAdminRegistry = new TokenAdminRegistry();
    for (uint256 i = 0; i < tokens.length; i++) {
      if (admins[i] == address(0)) {
        continue;
      }
      if (cleanTokenAdminRegistry.getTokenConfig(tokens[i]).isRegistered) {
        continue;
      }
      cleanTokenAdminRegistry.registerAdministratorPermissioned(tokens[i], admins[i]);
      s_AdminByToken[tokens[i]] = admins[i];
    }

    for (uint256 i = 0; i < tokens.length; i++) {
      assertTrue(cleanTokenAdminRegistry.isAdministrator(tokens[i], s_AdminByToken[tokens[i]]));
    }
  }

  function test_registerAdministratorPermissioned_ZeroAddress_Revert() public {
    address newToken = makeAddr("newToken");

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.ZeroAddress.selector));
    s_tokenAdminRegistry.registerAdministratorPermissioned(newToken, address(0));
  }

  function test_registerAdministratorPermissioned_AlreadyRegistered_Revert() public {
    address newAdmin = makeAddr("newAdmin");
    address newToken = makeAddr("newToken");

    s_tokenAdminRegistry.registerAdministratorPermissioned(newToken, newAdmin);

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.AlreadyRegistered.selector, newToken));
    s_tokenAdminRegistry.registerAdministratorPermissioned(newToken, newAdmin);
  }

  function test_registerAdministratorPermissioned_OnlyOwner_Revert() public {
    address newOwner = makeAddr("newOwner");
    address newToken = makeAddr("newToken");
    vm.stopPrank();

    vm.expectRevert("Only callable by owner");
    s_tokenAdminRegistry.registerAdministratorPermissioned(newToken, newOwner);
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
    emit RegistryModuleRemoved(newModule);

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
