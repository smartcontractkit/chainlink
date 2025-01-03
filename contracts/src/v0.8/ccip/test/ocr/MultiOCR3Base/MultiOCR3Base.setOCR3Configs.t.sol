// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {MultiOCR3Base} from "../../../ocr/MultiOCR3Base.sol";
import {MultiOCR3Helper} from "../../helpers/MultiOCR3Helper.sol";
import {MultiOCR3BaseSetup} from "./MultiOCR3BaseSetup.t.sol";

import {Vm} from "forge-std/Vm.sol";

contract MultiOCR3Base_setOCR3Configs is MultiOCR3BaseSetup {
  function test_SetConfigsZeroInput() public {
    vm.recordLogs();
    s_multiOCR3.setOCR3Configs(new MultiOCR3Base.OCRConfigArgs[](0));

    // No logs emitted
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);
  }

  function test_SetConfigWithSigners() public {
    uint8 F = 2;

    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(0));

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(F, s_validSigners, s_validTransmitters),
      F: F,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });

    vm.expectEmit();
    emit MultiOCR3Base.ConfigSet(
      ocrConfigs[0].ocrPluginType,
      ocrConfigs[0].configDigest,
      ocrConfigs[0].signers,
      ocrConfigs[0].transmitters,
      ocrConfigs[0].F
    );

    vm.expectEmit();
    emit MultiOCR3Helper.AfterConfigSet(ocrConfigs[0].ocrPluginType);

    s_multiOCR3.setOCR3Configs(ocrConfigs);

    MultiOCR3Base.OCRConfig memory expectedConfig = MultiOCR3Base.OCRConfig({
      configInfo: MultiOCR3Base.ConfigInfo({
        configDigest: ocrConfigs[0].configDigest,
        F: ocrConfigs[0].F,
        n: uint8(ocrConfigs[0].signers.length),
        isSignatureVerificationEnabled: ocrConfigs[0].isSignatureVerificationEnabled
      }),
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    _assertOCRConfigEquality(s_multiOCR3.latestConfigDetails(0), expectedConfig);
  }

  function test_SetConfigWithSignersMismatchingTransmitters() public {
    uint8 F = 2;

    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(0));

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(F, s_validSigners, s_partialTransmitters),
      F: F,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_partialTransmitters
    });

    vm.expectEmit();
    emit MultiOCR3Base.ConfigSet(
      ocrConfigs[0].ocrPluginType,
      ocrConfigs[0].configDigest,
      ocrConfigs[0].signers,
      ocrConfigs[0].transmitters,
      ocrConfigs[0].F
    );

    vm.expectEmit();
    emit MultiOCR3Helper.AfterConfigSet(ocrConfigs[0].ocrPluginType);

    s_multiOCR3.setOCR3Configs(ocrConfigs);

    MultiOCR3Base.OCRConfig memory expectedConfig = MultiOCR3Base.OCRConfig({
      configInfo: MultiOCR3Base.ConfigInfo({
        configDigest: ocrConfigs[0].configDigest,
        F: ocrConfigs[0].F,
        n: uint8(ocrConfigs[0].signers.length),
        isSignatureVerificationEnabled: ocrConfigs[0].isSignatureVerificationEnabled
      }),
      signers: s_validSigners,
      transmitters: s_partialTransmitters
    });
    _assertOCRConfigEquality(s_multiOCR3.latestConfigDetails(0), expectedConfig);
  }

  function test_SetConfigWithoutSigners() public {
    uint8 F = 1;
    address[] memory signers = new address[](0);

    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(0));

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(F, signers, s_validTransmitters),
      F: F,
      isSignatureVerificationEnabled: false,
      signers: signers,
      transmitters: s_validTransmitters
    });

    vm.expectEmit();
    emit MultiOCR3Base.ConfigSet(
      ocrConfigs[0].ocrPluginType,
      ocrConfigs[0].configDigest,
      ocrConfigs[0].signers,
      ocrConfigs[0].transmitters,
      ocrConfigs[0].F
    );

    vm.expectEmit();
    emit MultiOCR3Helper.AfterConfigSet(ocrConfigs[0].ocrPluginType);

    s_multiOCR3.setOCR3Configs(ocrConfigs);

    MultiOCR3Base.OCRConfig memory expectedConfig = MultiOCR3Base.OCRConfig({
      configInfo: MultiOCR3Base.ConfigInfo({
        configDigest: ocrConfigs[0].configDigest,
        F: ocrConfigs[0].F,
        n: uint8(ocrConfigs[0].signers.length),
        isSignatureVerificationEnabled: ocrConfigs[0].isSignatureVerificationEnabled
      }),
      signers: signers,
      transmitters: s_validTransmitters
    });
    _assertOCRConfigEquality(s_multiOCR3.latestConfigDetails(0), expectedConfig);
  }

  function test_SetConfigIgnoreSigners() public {
    uint8 F = 1;

    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(0));

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(F, new address[](0), s_validTransmitters),
      F: F,
      isSignatureVerificationEnabled: false,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });

    vm.expectEmit();
    emit MultiOCR3Base.ConfigSet(
      ocrConfigs[0].ocrPluginType,
      ocrConfigs[0].configDigest,
      s_emptySigners,
      ocrConfigs[0].transmitters,
      ocrConfigs[0].F
    );

    vm.expectEmit();
    emit MultiOCR3Helper.AfterConfigSet(ocrConfigs[0].ocrPluginType);

    s_multiOCR3.setOCR3Configs(ocrConfigs);

    MultiOCR3Base.OCRConfig memory expectedConfig = MultiOCR3Base.OCRConfig({
      configInfo: MultiOCR3Base.ConfigInfo({
        configDigest: ocrConfigs[0].configDigest,
        F: ocrConfigs[0].F,
        n: 0,
        isSignatureVerificationEnabled: ocrConfigs[0].isSignatureVerificationEnabled
      }),
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });
    _assertOCRConfigEquality(s_multiOCR3.latestConfigDetails(0), expectedConfig);

    // Verify no signer role is set
    for (uint256 i = 0; i < s_validSigners.length; ++i) {
      MultiOCR3Base.Oracle memory signerOracle = s_multiOCR3.getOracle(0, s_validSigners[i]);
      assertEq(uint8(signerOracle.role), uint8(MultiOCR3Base.Role.Unset));
    }
  }

  function test_SetMultipleConfigs() public {
    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(0));
    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(1));
    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(2));

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](3);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(2, s_validSigners, s_validTransmitters),
      F: 2,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    ocrConfigs[1] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 1,
      configDigest: _getBasicConfigDigest(1, s_validSigners, s_validTransmitters),
      F: 1,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    ocrConfigs[2] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 2,
      configDigest: _getBasicConfigDigest(1, s_partialSigners, s_partialTransmitters),
      F: 1,
      isSignatureVerificationEnabled: true,
      signers: s_partialSigners,
      transmitters: s_partialTransmitters
    });

    for (uint256 i; i < ocrConfigs.length; ++i) {
      vm.expectEmit();
      emit MultiOCR3Base.ConfigSet(
        ocrConfigs[i].ocrPluginType,
        ocrConfigs[i].configDigest,
        ocrConfigs[i].signers,
        ocrConfigs[i].transmitters,
        ocrConfigs[i].F
      );

      vm.expectEmit();
      emit MultiOCR3Helper.AfterConfigSet(ocrConfigs[i].ocrPluginType);
    }
    s_multiOCR3.setOCR3Configs(ocrConfigs);

    for (uint256 i; i < ocrConfigs.length; ++i) {
      MultiOCR3Base.OCRConfig memory expectedConfig = MultiOCR3Base.OCRConfig({
        configInfo: MultiOCR3Base.ConfigInfo({
          configDigest: ocrConfigs[i].configDigest,
          F: ocrConfigs[i].F,
          n: uint8(ocrConfigs[i].signers.length),
          isSignatureVerificationEnabled: ocrConfigs[i].isSignatureVerificationEnabled
        }),
        signers: ocrConfigs[i].signers,
        transmitters: ocrConfigs[i].transmitters
      });
      _assertOCRConfigEquality(s_multiOCR3.latestConfigDetails(ocrConfigs[i].ocrPluginType), expectedConfig);
    }

    // pluginType 3 remains unconfigured
    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(3));
  }

  function testFuzz_SetConfig_Success(MultiOCR3Base.OCRConfigArgs memory ocrConfig, uint64 randomAddressOffset) public {
    // condition: cannot assume max oracle count
    vm.assume(ocrConfig.transmitters.length <= 255);
    vm.assume(ocrConfig.signers.length <= 255);
    // condition: at least one transmitter
    vm.assume(ocrConfig.transmitters.length > 0);
    // condition: number of transmitters does not exceed signers
    vm.assume(ocrConfig.signers.length == 0 || ocrConfig.transmitters.length <= ocrConfig.signers.length);

    // condition: F > 0
    ocrConfig.F = uint8(bound(ocrConfig.F, 1, 3));

    uint256 transmittersLength = ocrConfig.transmitters.length;

    // Force addresses to be unique (with a random offset for broader testing)
    for (uint160 i = 0; i < transmittersLength; ++i) {
      ocrConfig.transmitters[i] = vm.addr(PRIVATE0 + randomAddressOffset + i);
      // condition: non-zero oracle address
      vm.assume(ocrConfig.transmitters[i] != address(0));
    }

    if (ocrConfig.signers.length == 0) {
      ocrConfig.isSignatureVerificationEnabled = false;
    } else {
      ocrConfig.isSignatureVerificationEnabled = true;

      // condition: number of signers > 3F
      vm.assume(ocrConfig.signers.length > 3 * ocrConfig.F);

      uint256 signersLength = ocrConfig.signers.length;

      // Force addresses to be unique - continuing generation with an offset after the transmitter addresses
      for (uint160 i = 0; i < signersLength; ++i) {
        ocrConfig.signers[i] = vm.addr(PRIVATE0 + randomAddressOffset + i + transmittersLength);
        // condition: non-zero oracle address
        vm.assume(ocrConfig.signers[i] != address(0));
      }
    }

    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(ocrConfig.ocrPluginType));

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = ocrConfig;

    vm.expectEmit();
    emit MultiOCR3Base.ConfigSet(
      ocrConfig.ocrPluginType, ocrConfig.configDigest, ocrConfig.signers, ocrConfig.transmitters, ocrConfig.F
    );
    vm.expectEmit();
    emit MultiOCR3Helper.AfterConfigSet(ocrConfig.ocrPluginType);
    s_multiOCR3.setOCR3Configs(ocrConfigs);

    MultiOCR3Base.OCRConfig memory expectedConfig = MultiOCR3Base.OCRConfig({
      configInfo: MultiOCR3Base.ConfigInfo({
        configDigest: ocrConfig.configDigest,
        F: ocrConfig.F,
        n: ocrConfig.isSignatureVerificationEnabled ? uint8(ocrConfig.signers.length) : 0,
        isSignatureVerificationEnabled: ocrConfig.isSignatureVerificationEnabled
      }),
      signers: ocrConfig.signers,
      transmitters: ocrConfig.transmitters
    });
    _assertOCRConfigEquality(s_multiOCR3.latestConfigDetails(ocrConfig.ocrPluginType), expectedConfig);
  }

  function test_UpdateConfigTransmittersWithoutSigners() public {
    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(0));

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(1, s_emptySigners, s_validTransmitters),
      F: 1,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });
    s_multiOCR3.setOCR3Configs(ocrConfigs);

    address[] memory newTransmitters = s_partialSigners;

    ocrConfigs[0].F = 2;
    ocrConfigs[0].configDigest = _getBasicConfigDigest(2, s_emptySigners, newTransmitters);
    ocrConfigs[0].transmitters = newTransmitters;

    vm.expectEmit();
    emit MultiOCR3Base.ConfigSet(
      ocrConfigs[0].ocrPluginType,
      ocrConfigs[0].configDigest,
      ocrConfigs[0].signers,
      ocrConfigs[0].transmitters,
      ocrConfigs[0].F
    );
    vm.expectEmit();
    emit MultiOCR3Helper.AfterConfigSet(ocrConfigs[0].ocrPluginType);

    s_multiOCR3.setOCR3Configs(ocrConfigs);

    MultiOCR3Base.OCRConfig memory expectedConfig = MultiOCR3Base.OCRConfig({
      configInfo: MultiOCR3Base.ConfigInfo({
        configDigest: ocrConfigs[0].configDigest,
        F: ocrConfigs[0].F,
        n: uint8(ocrConfigs[0].signers.length),
        isSignatureVerificationEnabled: ocrConfigs[0].isSignatureVerificationEnabled
      }),
      signers: s_emptySigners,
      transmitters: newTransmitters
    });
    _assertOCRConfigEquality(s_multiOCR3.latestConfigDetails(0), expectedConfig);

    // Verify oracle roles get correctly re-assigned
    for (uint256 i; i < newTransmitters.length; ++i) {
      MultiOCR3Base.Oracle memory transmitterOracle = s_multiOCR3.getOracle(0, newTransmitters[i]);
      assertEq(transmitterOracle.index, i);
      assertEq(uint8(transmitterOracle.role), uint8(MultiOCR3Base.Role.Transmitter));
    }

    // Verify old transmitters get correctly unset
    for (uint256 i = newTransmitters.length; i < s_validTransmitters.length; ++i) {
      MultiOCR3Base.Oracle memory transmitterOracle = s_multiOCR3.getOracle(0, s_validTransmitters[i]);
      assertEq(uint8(transmitterOracle.role), uint8(MultiOCR3Base.Role.Unset));
    }
  }

  function test_UpdateConfigSigners() public {
    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(0));

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(2, s_validSigners, s_validTransmitters),
      F: 2,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    s_multiOCR3.setOCR3Configs(ocrConfigs);

    address[] memory newSigners = s_partialTransmitters;
    address[] memory newTransmitters = s_partialSigners;

    ocrConfigs[0].F = 1;
    ocrConfigs[0].configDigest = _getBasicConfigDigest(1, newSigners, newTransmitters);
    ocrConfigs[0].signers = newSigners;
    ocrConfigs[0].transmitters = newTransmitters;

    vm.expectEmit();
    emit MultiOCR3Base.ConfigSet(
      ocrConfigs[0].ocrPluginType,
      ocrConfigs[0].configDigest,
      ocrConfigs[0].signers,
      ocrConfigs[0].transmitters,
      ocrConfigs[0].F
    );
    vm.expectEmit();
    emit MultiOCR3Helper.AfterConfigSet(ocrConfigs[0].ocrPluginType);

    s_multiOCR3.setOCR3Configs(ocrConfigs);

    MultiOCR3Base.OCRConfig memory expectedConfig = MultiOCR3Base.OCRConfig({
      configInfo: MultiOCR3Base.ConfigInfo({
        configDigest: ocrConfigs[0].configDigest,
        F: ocrConfigs[0].F,
        n: uint8(ocrConfigs[0].signers.length),
        isSignatureVerificationEnabled: ocrConfigs[0].isSignatureVerificationEnabled
      }),
      signers: newSigners,
      transmitters: newTransmitters
    });
    _assertOCRConfigEquality(s_multiOCR3.latestConfigDetails(0), expectedConfig);

    // Verify oracle roles get correctly re-assigned
    for (uint256 i; i < newSigners.length; ++i) {
      MultiOCR3Base.Oracle memory signerOracle = s_multiOCR3.getOracle(0, newSigners[i]);
      assertEq(signerOracle.index, i);
      assertEq(uint8(signerOracle.role), uint8(MultiOCR3Base.Role.Signer));

      MultiOCR3Base.Oracle memory transmitterOracle = s_multiOCR3.getOracle(0, newTransmitters[i]);
      assertEq(transmitterOracle.index, i);
      assertEq(uint8(transmitterOracle.role), uint8(MultiOCR3Base.Role.Transmitter));
    }

    // Verify old signers / transmitters get correctly unset
    for (uint256 i = newSigners.length; i < s_validSigners.length; ++i) {
      MultiOCR3Base.Oracle memory signerOracle = s_multiOCR3.getOracle(0, s_validSigners[i]);
      assertEq(uint8(signerOracle.role), uint8(MultiOCR3Base.Role.Unset));

      MultiOCR3Base.Oracle memory transmitterOracle = s_multiOCR3.getOracle(0, s_validTransmitters[i]);
      assertEq(uint8(transmitterOracle.role), uint8(MultiOCR3Base.Role.Unset));
    }
  }

  // Reverts

  function test_RevertWhen_RepeatTransmitterAddress() public {
    address[] memory signers = s_validSigners;
    address[] memory transmitters = s_validTransmitters;
    transmitters[0] = signers[0];

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(1, signers, transmitters),
      F: 1,
      isSignatureVerificationEnabled: true,
      signers: signers,
      transmitters: transmitters
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        MultiOCR3Base.InvalidConfig.selector, MultiOCR3Base.InvalidConfigErrorType.REPEATED_ORACLE_ADDRESS
      )
    );
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_RevertWhen_RepeatSignerAddress() public {
    address[] memory signers = s_validSigners;
    address[] memory transmitters = s_validTransmitters;
    signers[1] = signers[0];

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(1, signers, transmitters),
      F: 1,
      isSignatureVerificationEnabled: true,
      signers: signers,
      transmitters: transmitters
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        MultiOCR3Base.InvalidConfig.selector, MultiOCR3Base.InvalidConfigErrorType.REPEATED_ORACLE_ADDRESS
      )
    );
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_RevertWhen_SignerCannotBeZeroAddress() public {
    uint8 F = 1;
    address[] memory signers = new address[](3 * F + 1);
    address[] memory transmitters = new address[](3 * F + 1);
    for (uint160 i = 0; i < 3 * F + 1; ++i) {
      signers[i] = address(i + 1);
      transmitters[i] = address(i + 1000);
    }

    signers[0] = address(0);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(F, signers, transmitters),
      F: F,
      isSignatureVerificationEnabled: true,
      signers: signers,
      transmitters: transmitters
    });

    vm.expectRevert(MultiOCR3Base.OracleCannotBeZeroAddress.selector);
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_RevertWhen_TransmitterCannotBeZeroAddress() public {
    uint8 F = 1;
    address[] memory signers = new address[](3 * F + 1);
    address[] memory transmitters = new address[](3 * F + 1);
    for (uint160 i = 0; i < 3 * F + 1; ++i) {
      signers[i] = address(i + 1);
      transmitters[i] = address(i + 1000);
    }

    transmitters[0] = address(0);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(F, signers, transmitters),
      F: F,
      isSignatureVerificationEnabled: true,
      signers: signers,
      transmitters: transmitters
    });

    vm.expectRevert(MultiOCR3Base.OracleCannotBeZeroAddress.selector);
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_RevertWhen_StaticConfigChange() public {
    uint8 F = 1;

    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(0));

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(F, s_validSigners, s_validTransmitters),
      F: F,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });

    s_multiOCR3.setOCR3Configs(ocrConfigs);

    // signature verification cannot change
    ocrConfigs[0].isSignatureVerificationEnabled = false;
    vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.StaticConfigCannotBeChanged.selector, 0));
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_RevertWhen_FTooHigh() public {
    address[] memory signers = new address[](0);
    address[] memory transmitters = new address[](1);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(1, signers, transmitters),
      F: 1,
      isSignatureVerificationEnabled: true,
      signers: signers,
      transmitters: transmitters
    });

    vm.expectRevert(
      abi.encodeWithSelector(MultiOCR3Base.InvalidConfig.selector, MultiOCR3Base.InvalidConfigErrorType.F_TOO_HIGH)
    );
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_RevertWhen_FMustBePositive() public {
    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(0, s_validSigners, s_validTransmitters),
      F: 0,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        MultiOCR3Base.InvalidConfig.selector, MultiOCR3Base.InvalidConfigErrorType.F_MUST_BE_POSITIVE
      )
    );
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_RevertWhen_NoTransmitters() public {
    address[] memory signers = new address[](0);
    address[] memory transmitters = new address[](0);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(10, signers, transmitters),
      F: 1,
      isSignatureVerificationEnabled: false,
      signers: signers,
      transmitters: transmitters
    });

    vm.expectRevert(
      abi.encodeWithSelector(MultiOCR3Base.InvalidConfig.selector, MultiOCR3Base.InvalidConfigErrorType.NO_TRANSMITTERS)
    );
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_RevertWhen_TooManyTransmitters() public {
    address[] memory signers = new address[](0);
    address[] memory transmitters = new address[](257);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(10, signers, transmitters),
      F: 10,
      isSignatureVerificationEnabled: false,
      signers: signers,
      transmitters: transmitters
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        MultiOCR3Base.InvalidConfig.selector, MultiOCR3Base.InvalidConfigErrorType.TOO_MANY_TRANSMITTERS
      )
    );
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_RevertWhen_TooManySigners() public {
    address[] memory signers = new address[](257);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(1, signers, s_validTransmitters),
      F: 1,
      isSignatureVerificationEnabled: true,
      signers: signers,
      transmitters: s_validTransmitters
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        MultiOCR3Base.InvalidConfig.selector, MultiOCR3Base.InvalidConfigErrorType.TOO_MANY_SIGNERS
      )
    );
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_RevertWhen_MoreTransmittersThanSigners() public {
    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(1, s_validSigners, s_partialTransmitters),
      F: 1,
      isSignatureVerificationEnabled: true,
      signers: s_partialSigners,
      transmitters: s_validTransmitters
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        MultiOCR3Base.InvalidConfig.selector, MultiOCR3Base.InvalidConfigErrorType.TOO_MANY_TRANSMITTERS
      )
    );
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }
}
