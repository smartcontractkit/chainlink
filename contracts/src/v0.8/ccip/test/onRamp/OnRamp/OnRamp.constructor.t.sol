// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRMNRemote} from "../../../interfaces/IRMNRemote.sol";
import {IRouter} from "../../../interfaces/IRouter.sol";

import {Client} from "../../../libraries/Client.sol";
import {OnRamp} from "../../../onRamp/OnRamp.sol";
import {OnRampHelper} from "../../helpers/OnRampHelper.sol";
import {OnRampSetup} from "./OnRampSetup.t.sol";

contract OnRamp_constructor is OnRampSetup {
  function test_Constructor() public {
    OnRamp.StaticConfig memory staticConfig = OnRamp.StaticConfig({
      chainSelector: SOURCE_CHAIN_SELECTOR,
      rmnRemote: s_mockRMNRemote,
      nonceManager: address(s_outboundNonceManager),
      tokenAdminRegistry: address(s_tokenAdminRegistry)
    });
    OnRamp.DynamicConfig memory dynamicConfig = _generateDynamicOnRampConfig(address(s_feeQuoter));

    vm.expectEmit();
    emit OnRamp.ConfigSet(staticConfig, dynamicConfig);
    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(DEST_CHAIN_SELECTOR, 0, s_sourceRouter, false);

    _deployOnRamp(SOURCE_CHAIN_SELECTOR, s_sourceRouter, address(s_outboundNonceManager), address(s_tokenAdminRegistry));

    OnRamp.StaticConfig memory gotStaticConfig = s_onRamp.getStaticConfig();

    assertEq(staticConfig.chainSelector, gotStaticConfig.chainSelector);
    assertEq(address(staticConfig.rmnRemote), address(gotStaticConfig.rmnRemote));
    assertEq(staticConfig.tokenAdminRegistry, gotStaticConfig.tokenAdminRegistry);

    OnRamp.DynamicConfig memory gotDynamicConfig = s_onRamp.getDynamicConfig();
    assertEq(dynamicConfig.feeQuoter, gotDynamicConfig.feeQuoter);

    // Initial values
    assertEq("OnRamp 1.6.0", s_onRamp.typeAndVersion());
    assertEq(OWNER, s_onRamp.owner());
    assertEq(1, s_onRamp.getExpectedNextSequenceNumber(DEST_CHAIN_SELECTOR));
  }

  function test_RevertWhen_Constructor_EnableAllowList_ForwardFromRouter() public {
    OnRamp.StaticConfig memory staticConfig = OnRamp.StaticConfig({
      chainSelector: SOURCE_CHAIN_SELECTOR,
      rmnRemote: s_mockRMNRemote,
      nonceManager: address(s_outboundNonceManager),
      tokenAdminRegistry: address(s_tokenAdminRegistry)
    });

    OnRamp.DynamicConfig memory dynamicConfig = _generateDynamicOnRampConfig(address(s_feeQuoter));

    // Creating a DestChainConfig and setting allowlistEnabled : true
    OnRamp.DestChainConfigArgs[] memory destChainConfigs = new OnRamp.DestChainConfigArgs[](1);
    destChainConfigs[0] = OnRamp.DestChainConfigArgs({
      destChainSelector: DEST_CHAIN_SELECTOR,
      router: s_sourceRouter,
      allowlistEnabled: true
    });

    vm.expectEmit();
    emit OnRamp.ConfigSet(staticConfig, dynamicConfig);

    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(DEST_CHAIN_SELECTOR, 0, s_sourceRouter, true);

    OnRampHelper tempOnRamp = new OnRampHelper(staticConfig, dynamicConfig, destChainConfigs);

    // Sending a message and expecting revert as allowlist is enabled with no address in allowlist
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    vm.startPrank(address(s_sourceRouter));
    vm.expectRevert(abi.encodeWithSelector(OnRamp.SenderNotAllowed.selector, OWNER));
    tempOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);
  }

  function test_RevertWhen_Constructor_InvalidConfigChainSelectorEqZero() public {
    vm.expectRevert(OnRamp.InvalidConfig.selector);
    new OnRampHelper(
      OnRamp.StaticConfig({
        chainSelector: 0,
        rmnRemote: s_mockRMNRemote,
        nonceManager: address(s_outboundNonceManager),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _generateDynamicOnRampConfig(address(s_feeQuoter)),
      _generateDestChainConfigArgs(IRouter(address(0)))
    );
  }

  function test_RevertWhen_Constructor_InvalidConfigRMNProxyEqAddressZero() public {
    vm.expectRevert(OnRamp.InvalidConfig.selector);
    s_onRamp = new OnRampHelper(
      OnRamp.StaticConfig({
        chainSelector: SOURCE_CHAIN_SELECTOR,
        rmnRemote: IRMNRemote(address(0)),
        nonceManager: address(s_outboundNonceManager),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _generateDynamicOnRampConfig(address(s_feeQuoter)),
      _generateDestChainConfigArgs(IRouter(address(0)))
    );
  }

  function test_RevertWhen_Constructor_InvalidConfigNonceManagerEqAddressZero() public {
    vm.expectRevert(OnRamp.InvalidConfig.selector);
    new OnRampHelper(
      OnRamp.StaticConfig({
        chainSelector: SOURCE_CHAIN_SELECTOR,
        rmnRemote: s_mockRMNRemote,
        nonceManager: address(0),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      _generateDynamicOnRampConfig(address(s_feeQuoter)),
      _generateDestChainConfigArgs(IRouter(address(0)))
    );
  }

  function test_RevertWhen_Constructor_InvalidConfigTokenAdminRegistryEqAddressZero() public {
    vm.expectRevert(OnRamp.InvalidConfig.selector);
    new OnRampHelper(
      OnRamp.StaticConfig({
        chainSelector: SOURCE_CHAIN_SELECTOR,
        rmnRemote: s_mockRMNRemote,
        nonceManager: address(s_outboundNonceManager),
        tokenAdminRegistry: address(0)
      }),
      _generateDynamicOnRampConfig(address(s_feeQuoter)),
      _generateDestChainConfigArgs(IRouter(address(0)))
    );
  }
}
