// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {MultiOCR3Base} from "../../ocr/MultiOCR3Base.sol";
import {MultiOCR3Helper} from "../helpers/MultiOCR3Helper.sol";
import {MultiOCR3BaseSetup} from "./MultiOCR3BaseSetup.t.sol";

import {Vm} from "forge-std/Vm.sol";

contract MultiOCR3Base_transmit is MultiOCR3BaseSetup {
  bytes32 internal s_configDigest1;
  bytes32 internal s_configDigest2;
  bytes32 internal s_configDigest3;

  function setUp() public virtual override {
    super.setUp();

    s_configDigest1 = _getBasicConfigDigest(1, s_validSigners, s_validTransmitters);
    s_configDigest2 = _getBasicConfigDigest(1, s_validSigners, s_validTransmitters);
    s_configDigest3 = _getBasicConfigDigest(2, s_emptySigners, s_validTransmitters);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](3);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: s_configDigest1,
      F: 1,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    ocrConfigs[1] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 1,
      configDigest: s_configDigest2,
      F: 2,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    ocrConfigs[2] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 2,
      configDigest: s_configDigest3,
      F: 1,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });

    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_TransmitSigners_gas_Success() public {
    vm.pauseGasMetering();
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    // F = 2, need 2 signatures
    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, REPORT, reportContext, 2);

    s_multiOCR3.setTransmitOcrPluginType(0);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(0, s_configDigest1, uint64(uint256(s_configDigest1)));

    vm.startPrank(s_validTransmitters[1]);
    vm.resumeGasMetering();
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, rawVs);
  }

  function test_TransmitWithoutSignatureVerification_gas_Success() public {
    vm.pauseGasMetering();
    bytes32[3] memory reportContext = [s_configDigest3, s_configDigest3, s_configDigest3];

    s_multiOCR3.setTransmitOcrPluginType(2);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(2, s_configDigest3, uint64(uint256(s_configDigest3)));

    vm.startPrank(s_validTransmitters[0]);
    vm.resumeGasMetering();
    s_multiOCR3.transmitWithoutSignatures(reportContext, REPORT);
  }

  function test_Fuzz_TransmitSignersWithSignatures_Success(uint8 F, uint64 randomAddressOffset) public {
    vm.pauseGasMetering();

    F = uint8(bound(F, 1, 3));

    // condition: signers.length > 3F
    uint8 signersLength = 3 * F + 1;
    address[] memory signers = new address[](signersLength);
    address[] memory transmitters = new address[](signersLength);
    uint256[] memory signerKeys = new uint256[](signersLength);

    // Force addresses to be unique (with a random offset for broader testing)
    for (uint160 i = 0; i < signersLength; ++i) {
      transmitters[i] = vm.addr(PRIVATE0 + randomAddressOffset + i);
      // condition: non-zero oracle address
      vm.assume(transmitters[i] != address(0));

      // condition: non-repeating addresses (no clashes with transmitters)
      signerKeys[i] = PRIVATE0 + randomAddressOffset + i + signersLength;
      signers[i] = vm.addr(signerKeys[i]);
      vm.assume(signers[i] != address(0));
    }

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 3,
      configDigest: s_configDigest1,
      F: F,
      isSignatureVerificationEnabled: true,
      signers: signers,
      transmitters: transmitters
    });
    s_multiOCR3.setOCR3Configs(ocrConfigs);
    s_multiOCR3.setTransmitOcrPluginType(3);

    // Randomise picked transmitter with random offset
    vm.startPrank(transmitters[randomAddressOffset % signersLength]);

    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    // condition: matches signature expectation for transmit
    uint8 numSignatures = F + 1;
    uint256[] memory pickedSignerKeys = new uint256[](numSignatures);

    // Randomise picked signers with random offset
    for (uint256 i; i < numSignatures; ++i) {
      pickedSignerKeys[i] = signerKeys[(i + randomAddressOffset) % numSignatures];
    }

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(pickedSignerKeys, REPORT, reportContext, numSignatures);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(3, s_configDigest1, uint64(uint256(s_configDigest1)));

    vm.resumeGasMetering();
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, rawVs);
  }

  // Reverts
  function test_ForkedChain_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, REPORT, reportContext, 2);

    s_multiOCR3.setTransmitOcrPluginType(0);

    uint256 chain1 = block.chainid;
    uint256 chain2 = chain1 + 1;
    vm.chainId(chain2);
    vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.ForkedChain.selector, chain1, chain2));

    vm.startPrank(s_validTransmitters[0]);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, rawVs);
  }

  function test_ZeroSignatures_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    s_multiOCR3.setTransmitOcrPluginType(0);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.WrongNumberOfSignatures.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, new bytes32[](0), new bytes32[](0), bytes32(""));
  }

  function test_TooManySignatures_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    // 1 signature too many
    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, REPORT, reportContext, 6);

    s_multiOCR3.setTransmitOcrPluginType(1);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.WrongNumberOfSignatures.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, rawVs);
  }

  function test_InsufficientSignatures_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    // Missing 1 signature for unique report
    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, REPORT, reportContext, 4);

    s_multiOCR3.setTransmitOcrPluginType(1);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.WrongNumberOfSignatures.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, rawVs);
  }

  function test_ConfigDigestMismatch_Revert() public {
    bytes32 configDigest;
    bytes32[3] memory reportContext = [configDigest, configDigest, configDigest];

    (,,, bytes32 rawVs) = _getSignaturesForDigest(s_validSignerKeys, REPORT, reportContext, 2);

    s_multiOCR3.setTransmitOcrPluginType(0);

    vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.ConfigDigestMismatch.selector, s_configDigest1, configDigest));
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, new bytes32[](0), new bytes32[](0), rawVs);
  }

  function test_SignatureOutOfRegistration_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](1);

    s_multiOCR3.setTransmitOcrPluginType(0);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.SignaturesOutOfRegistration.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, bytes32(""));
  }

  function test_UnAuthorizedTransmitter_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](2);

    s_multiOCR3.setTransmitOcrPluginType(0);

    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, bytes32(""));
  }

  function test_NonUniqueSignature_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    (bytes32[] memory rs, bytes32[] memory ss, uint8[] memory vs, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, REPORT, reportContext, 2);

    rs[1] = rs[0];
    ss[1] = ss[0];
    // Need to reset the rawVs to be valid
    rawVs = bytes32(bytes1(vs[0] - 27)) | (bytes32(bytes1(vs[0] - 27)) >> 8);

    s_multiOCR3.setTransmitOcrPluginType(0);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.NonUniqueSignatures.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, rawVs);
  }

  function test_UnauthorizedSigner_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, REPORT, reportContext, 2);

    rs[0] = s_configDigest1;
    ss = rs;

    s_multiOCR3.setTransmitOcrPluginType(0);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.UnauthorizedSigner.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, rawVs);
  }

  function test_UnconfiguredPlugin_Revert() public {
    bytes32 configDigest;
    bytes32[3] memory reportContext = [configDigest, configDigest, configDigest];

    s_multiOCR3.setTransmitOcrPluginType(42);

    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_multiOCR3.transmitWithoutSignatures(reportContext, REPORT);
  }

  function test_TransmitWithLessCalldataArgs_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    s_multiOCR3.setTransmitOcrPluginType(0);

    // The transmit should fail, since we are trying to transmit without signatures when signatures are enabled
    vm.startPrank(s_validTransmitters[1]);

    // report length + function selector + report length + abiencoded location of report value + report context words
    uint256 receivedLength = REPORT.length + 4 + 5 * 32;
    vm.expectRevert(
      abi.encodeWithSelector(
        MultiOCR3Base.WrongMessageLength.selector,
        // Expecting inclusion of signature constant length components
        receivedLength + 5 * 32,
        receivedLength
      )
    );
    s_multiOCR3.transmitWithoutSignatures(reportContext, REPORT);
  }

  function test_TransmitWithExtraCalldataArgs_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](2);

    s_multiOCR3.setTransmitOcrPluginType(2);

    // The transmit should fail, since we are trying to transmit with signatures when signatures are disabled
    vm.startPrank(s_validTransmitters[1]);

    // dynamic length + function selector + report length + abiencoded location of report value + report context words
    // rawVs value, lengths of rs, ss, and start locations of rs & ss -> 5 words
    uint256 receivedLength = REPORT.length + 4 + (5 * 32) + (5 * 32) + (2 * 32) + (2 * 32);
    vm.expectRevert(
      abi.encodeWithSelector(
        MultiOCR3Base.WrongMessageLength.selector,
        // Expecting exclusion of signature constant length components and rs, ss words
        receivedLength - (5 * 32) - (4 * 32),
        receivedLength
      )
    );
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, bytes32(""));
  }
}

contract MultiOCR3Base_setOCR3Configs is MultiOCR3BaseSetup {
  function test_SetConfigsZeroInput_Success() public {
    vm.recordLogs();
    s_multiOCR3.setOCR3Configs(new MultiOCR3Base.OCRConfigArgs[](0));

    // No logs emitted
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);
  }

  function test_SetConfigWithSigners_Success() public {
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

  function test_SetConfigWithSignersMismatchingTransmitters_Success() public {
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

  function test_SetConfigWithoutSigners_Success() public {
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

  function test_SetConfigIgnoreSigners_Success() public {
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

  function test_SetMultipleConfigs_Success() public {
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

  function test_Fuzz_SetConfig_Success(MultiOCR3Base.OCRConfigArgs memory ocrConfig, uint64 randomAddressOffset) public {
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

  function test_UpdateConfigTransmittersWithoutSigners_Success() public {
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

  function test_UpdateConfigSigners_Success() public {
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

  function test_RepeatTransmitterAddress_Revert() public {
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

  function test_RepeatSignerAddress_Revert() public {
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

  function test_SignerCannotBeZeroAddress_Revert() public {
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

  function test_TransmitterCannotBeZeroAddress_Revert() public {
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

  function test_StaticConfigChange_Revert() public {
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

  function test_FTooHigh_Revert() public {
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

  function test_FMustBePositive_Revert() public {
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

  function test_NoTransmitters_Revert() public {
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

  function test_TooManyTransmitters_Revert() public {
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

  function test_TooManySigners_Revert() public {
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

  function test_MoreTransmittersThanSigners_Revert() public {
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
