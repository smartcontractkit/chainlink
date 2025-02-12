// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRMNRemote} from "../../../interfaces/IRMNRemote.sol";

import {Internal} from "../../../libraries/Internal.sol";
import {MultiOCR3Base} from "../../../ocr/MultiOCR3Base.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {OffRampHelper} from "../../helpers/OffRampHelper.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

contract OffRamp_constructor is OffRampSetup {
  function test_Constructor() public {
    OffRamp.StaticConfig memory staticConfig = OffRamp.StaticConfig({
      chainSelector: DEST_CHAIN_SELECTOR,
      gasForCallExactCheck: GAS_FOR_CALL_EXACT_CHECK,
      rmnRemote: s_mockRMNRemote,
      tokenAdminRegistry: address(s_tokenAdminRegistry),
      nonceManager: address(s_inboundNonceManager)
    });
    OffRamp.DynamicConfig memory dynamicConfig = _generateDynamicOffRampConfig(address(s_feeQuoter));

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](2);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });
    sourceChainConfigs[1] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1 + 1,
      onRamp: ON_RAMP_ADDRESS_2,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    OffRamp.SourceChainConfig memory expectedSourceChainConfig1 = OffRamp.SourceChainConfig({
      router: s_destRouter,
      isEnabled: true,
      minSeqNr: 1,
      onRamp: sourceChainConfigs[0].onRamp,
      isRMNVerificationDisabled: false
    });

    OffRamp.SourceChainConfig memory expectedSourceChainConfig2 = OffRamp.SourceChainConfig({
      router: s_destRouter,
      isEnabled: true,
      minSeqNr: 1,
      onRamp: sourceChainConfigs[1].onRamp,
      isRMNVerificationDisabled: false
    });

    uint64[] memory expectedSourceChainSelectors = new uint64[](2);
    expectedSourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;
    expectedSourceChainSelectors[1] = SOURCE_CHAIN_SELECTOR_1 + 1;

    vm.expectEmit();
    emit OffRamp.StaticConfigSet(staticConfig);

    vm.expectEmit();
    emit OffRamp.DynamicConfigSet(dynamicConfig);

    vm.expectEmit();
    emit OffRamp.SourceChainSelectorAdded(SOURCE_CHAIN_SELECTOR_1);

    vm.expectEmit();
    emit OffRamp.SourceChainConfigSet(SOURCE_CHAIN_SELECTOR_1, expectedSourceChainConfig1);

    vm.expectEmit();
    emit OffRamp.SourceChainSelectorAdded(SOURCE_CHAIN_SELECTOR_1 + 1);

    vm.expectEmit();
    emit OffRamp.SourceChainConfigSet(SOURCE_CHAIN_SELECTOR_1 + 1, expectedSourceChainConfig2);

    s_offRamp = new OffRampHelper(staticConfig, dynamicConfig, sourceChainConfigs);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Execution),
      configDigest: s_configDigestExec,
      F: F,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });

    s_offRamp.setOCR3Configs(ocrConfigs);

    // Static config
    OffRamp.StaticConfig memory gotStaticConfig = s_offRamp.getStaticConfig();
    assertEq(staticConfig.chainSelector, gotStaticConfig.chainSelector);
    assertEq(address(staticConfig.rmnRemote), address(gotStaticConfig.rmnRemote));
    assertEq(staticConfig.tokenAdminRegistry, gotStaticConfig.tokenAdminRegistry);

    // Dynamic config
    OffRamp.DynamicConfig memory gotDynamicConfig = s_offRamp.getDynamicConfig();
    _assertSameConfig(dynamicConfig, gotDynamicConfig);

    // OCR Config
    MultiOCR3Base.OCRConfig memory expectedOCRConfig = MultiOCR3Base.OCRConfig({
      configInfo: MultiOCR3Base.ConfigInfo({
        configDigest: ocrConfigs[0].configDigest,
        F: ocrConfigs[0].F,
        n: 0,
        isSignatureVerificationEnabled: ocrConfigs[0].isSignatureVerificationEnabled
      }),
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });
    MultiOCR3Base.OCRConfig memory gotOCRConfig = s_offRamp.latestConfigDetails(uint8(Internal.OCRPluginType.Execution));
    _assertOCRConfigEquality(expectedOCRConfig, gotOCRConfig);

    (uint64[] memory actualSourceChainSelectors, OffRamp.SourceChainConfig[] memory actualSourceChainConfigs) =
      s_offRamp.getAllSourceChainConfigs();

    _assertSourceChainConfigEquality(actualSourceChainConfigs[0], expectedSourceChainConfig1);
    _assertSourceChainConfigEquality(actualSourceChainConfigs[1], expectedSourceChainConfig2);

    // OffRamp initial values
    assertEq("OffRamp 1.6.0", s_offRamp.typeAndVersion());
    assertEq(OWNER, s_offRamp.owner());
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());

    // assertion for source chain selector
    for (uint256 i = 0; i < expectedSourceChainSelectors.length; i++) {
      assertEq(expectedSourceChainSelectors[i], actualSourceChainSelectors[i]);
    }
  }

  // Revert
  function test_RevertWhen_ZeroOnRampAddress() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: new bytes(0),
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        gasForCallExactCheck: GAS_FOR_CALL_EXACT_CHECK,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );
  }

  function test_RevertWhen_SourceChainSelector() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: 0,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    vm.expectRevert(OffRamp.ZeroChainSelectorNotAllowed.selector);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        gasForCallExactCheck: GAS_FOR_CALL_EXACT_CHECK,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );
  }

  function test_RevertWhen_ZeroRMNRemote() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](0);

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        gasForCallExactCheck: GAS_FOR_CALL_EXACT_CHECK,
        rmnRemote: IRMNRemote(address(0)),
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );
  }

  function test_RevertWhen_ZeroChainSelector() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](0);

    vm.expectRevert(OffRamp.ZeroChainSelectorNotAllowed.selector);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: 0,
        gasForCallExactCheck: GAS_FOR_CALL_EXACT_CHECK,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );
  }

  function test_RevertWhen_ZeroTokenAdminRegistry() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](0);

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        gasForCallExactCheck: GAS_FOR_CALL_EXACT_CHECK,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: address(0),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );
  }

  function test_RevertWhen_ZeroNonceManager() public {
    uint64[] memory sourceChainSelectors = new uint64[](1);
    sourceChainSelectors[0] = SOURCE_CHAIN_SELECTOR_1;

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](0);

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        gasForCallExactCheck: GAS_FOR_CALL_EXACT_CHECK,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(0)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      sourceChainConfigs
    );
  }
}
