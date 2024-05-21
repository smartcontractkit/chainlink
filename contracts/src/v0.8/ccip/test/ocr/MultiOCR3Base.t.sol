// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {MultiOCR3Base} from "../../ocr/MultiOCR3Base.sol";
import {BaseTest} from "../BaseTest.t.sol";
import {MultiOCR3Helper} from "../helpers/MultiOCR3Helper.sol";

import {Vm} from "forge-std/Vm.sol";

// TODO: revisit pulling more tests from OCR2BaseNoChecks
contract MultiOCR3BaseSetup is BaseTest {
  // Signer private keys used for these test
  uint256 internal constant PRIVATE0 = 0x7b2e97fe057e6de99d6872a2ef2abf52c9b4469bc848c2465ac3fcd8d336e81d;
  uint256 internal constant PRIVATE1 = 0xab56160806b05ef1796789248e1d7f34a6465c5280899159d645218cd216cee6;
  uint256 internal constant PRIVATE2 = 0x6ec7caa8406a49b76736602810e0a2871959fbbb675e23a8590839e4717f1f7f;
  uint256 internal constant PRIVATE3 = 0x80f14b11da94ae7f29d9a7713ea13dc838e31960a5c0f2baf45ed458947b730a;

  address[] internal s_validSigners;
  address[] internal s_validTransmitters;
  uint256[] internal s_validSignerKeys;

  bytes internal constant REPORT = abi.encode("testReport");
  MultiOCR3Helper internal s_multiOCR3;

  function setUp() public virtual override {
    BaseTest.setUp();

    s_validTransmitters = new address[](4);
    for (uint160 i = 0; i < 4; ++i) {
      s_validTransmitters[i] = address(4 + i);
    }

    s_validSignerKeys = new uint256[](4);
    s_validSignerKeys[0] = PRIVATE0;
    s_validSignerKeys[1] = PRIVATE1;
    s_validSignerKeys[2] = PRIVATE2;
    s_validSignerKeys[3] = PRIVATE3;

    //0xc110458BE52CaA6bB68E66969C3218A4D9Db0211
    //0xc110a19c08f1da7F5FfB281dc93630923F8E3719
    //0xc110fdF6e8fD679C7Cc11602d1cd829211A18e9b
    //0xc11028017c9b445B6bF8aE7da951B5cC28B326C0
    s_validSigners = new address[](s_validSignerKeys.length);
    for (uint256 i; i < s_validSignerKeys.length; ++i) {
      s_validSigners[i] = vm.addr(s_validSignerKeys[i]);
    }

    s_multiOCR3 = new MultiOCR3Helper();
  }

  /// @dev returns a mock config digest with config digest computation logic similar to OCR2Base
  function _getBasicConfigDigest(
    uint8 F,
    address[] memory signers,
    address[] memory transmitters
  ) internal view returns (bytes32) {
    bytes memory configBytes = abi.encode("");
    uint256 configVersion = 1;

    uint256 h = uint256(
      keccak256(
        abi.encode(
          block.chainid, address(s_multiOCR3), signers, transmitters, F, configBytes, configVersion, configBytes
        )
      )
    );
    uint256 prefixMask = type(uint256).max << (256 - 16); // 0xFFFF00..00
    uint256 prefix = 0x0001 << (256 - 16); // 0x000100..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }

  /// @dev returns a hash value in the same format as the h value on which the signature verified
  ///      in the _transmit function
  function _getTestReportDigest(bytes32 configDigest) internal pure returns (bytes32) {
    bytes32[3] memory reportContext = [configDigest, configDigest, configDigest];
    return keccak256(abi.encodePacked(keccak256(REPORT), reportContext));
  }

  function _assertOCRConfigEquality(
    MultiOCR3Base.OCRConfig memory configA,
    MultiOCR3Base.OCRConfig memory configB
  ) internal pure {
    vm.assertEq(configA.configInfo.configDigest, configB.configInfo.configDigest);
    vm.assertEq(configA.configInfo.F, configB.configInfo.F);
    vm.assertEq(configA.configInfo.n, configB.configInfo.n);
    vm.assertEq(configA.configInfo.uniqueReports, configB.configInfo.uniqueReports);
    vm.assertEq(configA.configInfo.isSignatureVerificationEnabled, configB.configInfo.isSignatureVerificationEnabled);

    vm.assertEq(configA.signers, configB.signers);
    vm.assertEq(configA.transmitters, configB.transmitters);
  }

  function _assertOCRConfigUnconfigured(MultiOCR3Base.OCRConfig memory config) internal pure {
    assertEq(config.configInfo.configDigest, bytes32(""));
    assertEq(config.signers.length, 0);
    assertEq(config.transmitters.length, 0);
  }

  function _getSignaturesForDigest(
    uint256[] memory signerPrivateKeys,
    bytes32 configDigest
  ) internal pure returns (bytes32[] memory rs, bytes32[] memory ss, uint8[] memory vs, bytes32 rawVs) {
    rs = new bytes32[](signerPrivateKeys.length);
    ss = new bytes32[](signerPrivateKeys.length);
    vs = new uint8[](signerPrivateKeys.length);

    // Calculate signatures
    for (uint256 i; i < signerPrivateKeys.length; ++i) {
      (vs[i], rs[i], ss[i]) = vm.sign(signerPrivateKeys[i], _getTestReportDigest(configDigest));
      rawVs = rawVs | (bytes32(bytes1(vs[i] - 27)) >> (8 * i));
    }

    return (rs, ss, vs, rawVs);
  }
}

contract MultiOCR3Base_transmit is MultiOCR3BaseSetup {
  address[] internal s_emptySigners;
  address[] internal s_partialSigners;
  uint256[] internal s_partialSignerKeys;

  bytes32 internal s_configDigest1;
  bytes32 internal s_configDigest2;
  bytes32 internal s_configDigest3;

  function setUp() public virtual override {
    super.setUp();

    s_emptySigners = new address[](0);

    s_partialSigners = new address[](2);
    s_partialSigners[0] = s_validSigners[0];
    s_partialSigners[1] = s_validSigners[1];

    s_partialSignerKeys = new uint256[](2);
    s_partialSignerKeys[0] = s_validSignerKeys[0];
    s_partialSignerKeys[1] = s_validSignerKeys[1];

    s_configDigest1 = _getBasicConfigDigest(1, s_validSigners, s_validTransmitters);
    s_configDigest2 = _getBasicConfigDigest(1, s_validSigners, s_validTransmitters);
    s_configDigest3 = _getBasicConfigDigest(2, s_emptySigners, s_validTransmitters);

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](3);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: s_configDigest1,
      F: 1,
      uniqueReports: false,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    ocrConfigs[1] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 1,
      configDigest: s_configDigest2,
      F: 1,
      uniqueReports: true,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    ocrConfigs[2] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 2,
      configDigest: s_configDigest3,
      F: 2,
      uniqueReports: false,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });

    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_TransmitSignersNonUniqueReports_gas_Success() public {
    vm.pauseGasMetering();
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_partialSignerKeys, s_configDigest1);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(0, s_configDigest1, uint32(uint256(s_configDigest1) >> 8));

    vm.startPrank(s_validTransmitters[1]);
    vm.resumeGasMetering();
    s_multiOCR3.transmit(0, reportContext, REPORT, rs, ss, rawVs);
  }

  function test_TransmitUniqueReportSigners_gas_Success() public {
    vm.pauseGasMetering();
    bytes32[3] memory reportContext = [s_configDigest2, s_configDigest2, s_configDigest2];

    uint256[] memory signerKeys = new uint256[](3);
    signerKeys[0] = s_validSignerKeys[0];
    signerKeys[1] = s_validSignerKeys[1];
    signerKeys[2] = s_validSignerKeys[2];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) = _getSignaturesForDigest(signerKeys, s_configDigest2);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(1, s_configDigest2, uint32(uint256(s_configDigest2) >> 8));

    vm.startPrank(s_validTransmitters[2]);
    vm.resumeGasMetering();
    s_multiOCR3.transmit(1, reportContext, REPORT, rs, ss, rawVs);
  }

  function test_TransmitWithoutSignatureVerification_gas_Success() public {
    vm.pauseGasMetering();
    bytes32[3] memory reportContext = [s_configDigest3, s_configDigest3, s_configDigest3];

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(2, s_configDigest3, uint32(uint256(s_configDigest3) >> 8));

    vm.startPrank(s_validTransmitters[0]);
    vm.resumeGasMetering();
    s_multiOCR3.transmit(2, reportContext, REPORT, new bytes32[](0), new bytes32[](0), bytes32(""));
  }

  // TODO: revisit the below 2 test-cases
  // function test_TransmitWithLessCalldataArgs_Success() public {
  //   vm.pauseGasMetering();
  //   bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

  //   (bytes32[] memory rs, bytes32[] memory ss, , bytes32 rawVs) =
  //     _getSignaturesForDigest(s_partialSignerKeys, s_configDigest1);

  //   vm.expectEmit();
  //   emit MultiOCR3Base.Transmitted(0, s_configDigest1, uint32(uint256(s_configDigest1) >> 8));

  //   // The transmit should succeed, even though there are more args to the external transmit function
  //   vm.startPrank(s_validTransmitters[1]);
  //   vm.resumeGasMetering();
  //   s_multiOCR3.transmit2(42, 0, reportContext, REPORT, rs, ss, rawVs);
  // }

  // function test_TransmitWithExtraCalldataArgs_Success() public {
  //   vm.pauseGasMetering();
  //   bytes32[3] memory reportContext = [s_configDigest3, s_configDigest3, s_configDigest3];

  //   vm.expectEmit();
  //   emit MultiOCR3Base.Transmitted(2, s_configDigest3, uint32(uint256(s_configDigest3) >> 8));

  //   vm.startPrank(s_validTransmitters[0]);
  //   vm.resumeGasMetering();
  //   s_multiOCR3.transmit3(2, reportContext, REPORT);
  // }

  // Reverts
  function test_ForkedChain_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_partialSignerKeys, s_configDigest1);

    uint256 chain1 = block.chainid;
    uint256 chain2 = chain1 + 1;
    vm.chainId(chain2);
    vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.ForkedChain.selector, chain1, chain2));
    vm.startPrank(s_validTransmitters[0]);
    s_multiOCR3.transmit(0, reportContext, REPORT, rs, ss, rawVs);
  }

  function test_ZeroSignatures_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.WrongNumberOfSignatures.selector);
    s_multiOCR3.transmit(0, reportContext, REPORT, new bytes32[](0), new bytes32[](0), bytes32(""));
  }

  function test_TooManySignatures_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    // 1 signature too many
    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, s_configDigest2);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.WrongNumberOfSignatures.selector);
    s_multiOCR3.transmit(1, reportContext, REPORT, rs, ss, rawVs);
  }

  function test_InsufficientSignatures_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    // Missing 1 signature for unique report
    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_partialSignerKeys, s_configDigest2);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.WrongNumberOfSignatures.selector);
    s_multiOCR3.transmit(1, reportContext, REPORT, rs, ss, rawVs);
  }

  function test_ConfigDigestMismatch_Revert() public {
    bytes32 configDigest;
    bytes32[3] memory reportContext = [configDigest, configDigest, configDigest];

    (,,, bytes32 rawVs) = _getSignaturesForDigest(s_partialSignerKeys, s_configDigest1);

    vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.ConfigDigestMismatch.selector, s_configDigest1, configDigest));
    s_multiOCR3.transmit(0, reportContext, REPORT, new bytes32[](0), new bytes32[](0), rawVs);
  }

  function test_SignatureOutOfRegistration_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](1);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.SignaturesOutOfRegistration.selector);
    s_multiOCR3.transmit(0, reportContext, REPORT, rs, ss, bytes32(""));
  }

  function test_UnAuthorizedTransmitter_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](2);

    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_multiOCR3.transmit(0, reportContext, REPORT, rs, ss, bytes32(""));
  }

  function test_NonUniqueSignature_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    (bytes32[] memory rs, bytes32[] memory ss, uint8[] memory vs, bytes32 rawVs) =
      _getSignaturesForDigest(s_partialSignerKeys, s_configDigest1);

    rs[1] = rs[0];
    ss[1] = ss[0];
    // Need to reset the rawVs to be valid
    rawVs = bytes32(bytes1(vs[0] - 27)) | (bytes32(bytes1(vs[0] - 27)) >> 8);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.NonUniqueSignatures.selector);
    s_multiOCR3.transmit(0, reportContext, REPORT, rs, ss, rawVs);
  }

  function test_UnauthorizedSigner_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest1, s_configDigest1, s_configDigest1];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_partialSignerKeys, s_configDigest1);

    rs[0] = s_configDigest1;
    ss = rs;

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.UnauthorizedSigner.selector);
    s_multiOCR3.transmit(0, reportContext, REPORT, rs, ss, rawVs);
  }

  function test_UnconfiguredPlugin_Revert() public {
    bytes32 configDigest;
    bytes32[3] memory reportContext = [configDigest, configDigest, configDigest];
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](2);

    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_multiOCR3.transmit(42, reportContext, REPORT, rs, ss, bytes32(""));
  }
}

contract MultiOCR3Base_setOCR3Configs is MultiOCR3BaseSetup {
  // TODO: fuzz test for setOCR3Configs (single config)

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
      uniqueReports: false,
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
    s_multiOCR3.setOCR3Configs(ocrConfigs);

    MultiOCR3Base.OCRConfig memory expectedConfig = MultiOCR3Base.OCRConfig({
      configInfo: MultiOCR3Base.ConfigInfo({
        configDigest: ocrConfigs[0].configDigest,
        F: ocrConfigs[0].F,
        n: uint8(s_validSigners.length),
        uniqueReports: ocrConfigs[0].uniqueReports,
        isSignatureVerificationEnabled: ocrConfigs[0].isSignatureVerificationEnabled
      }),
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    _assertOCRConfigEquality(s_multiOCR3.latestConfigDetails(0), expectedConfig);
  }

  function test_SetConfigWithoutSigners_Success() public {
    uint8 F = 2;
    address[] memory signers = new address[](0);

    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(0));

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(F, signers, s_validTransmitters),
      F: F,
      uniqueReports: false,
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
    s_multiOCR3.setOCR3Configs(ocrConfigs);

    MultiOCR3Base.OCRConfig memory expectedConfig = MultiOCR3Base.OCRConfig({
      configInfo: MultiOCR3Base.ConfigInfo({
        configDigest: ocrConfigs[0].configDigest,
        F: ocrConfigs[0].F,
        n: uint8(s_validTransmitters.length),
        uniqueReports: ocrConfigs[0].uniqueReports,
        isSignatureVerificationEnabled: ocrConfigs[0].isSignatureVerificationEnabled
      }),
      signers: signers,
      transmitters: s_validTransmitters
    });
    _assertOCRConfigEquality(s_multiOCR3.latestConfigDetails(0), expectedConfig);
  }

  function test_SetMultipleConfigs_Success() public {
    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(0));
    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(1));
    _assertOCRConfigUnconfigured(s_multiOCR3.latestConfigDetails(2));

    address[] memory partialSigners = new address[](2);
    partialSigners[0] = s_validSigners[0];
    partialSigners[1] = s_validSigners[1];

    address[] memory partialTransmitters = new address[](2);
    partialTransmitters[0] = s_validTransmitters[0];
    partialTransmitters[1] = s_validTransmitters[1];

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](3);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(2, s_validSigners, s_validTransmitters),
      F: 2,
      uniqueReports: false,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    ocrConfigs[1] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 1,
      configDigest: _getBasicConfigDigest(1, s_validSigners, s_validTransmitters),
      F: 1,
      uniqueReports: true,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    ocrConfigs[2] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 2,
      configDigest: _getBasicConfigDigest(1, partialSigners, partialSigners),
      F: 1,
      uniqueReports: true,
      isSignatureVerificationEnabled: true,
      signers: partialSigners,
      transmitters: partialTransmitters
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
    }
    s_multiOCR3.setOCR3Configs(ocrConfigs);

    for (uint256 i; i < ocrConfigs.length; ++i) {
      MultiOCR3Base.OCRConfig memory expectedConfig = MultiOCR3Base.OCRConfig({
        configInfo: MultiOCR3Base.ConfigInfo({
          configDigest: ocrConfigs[i].configDigest,
          F: ocrConfigs[i].F,
          n: uint8(ocrConfigs[i].signers.length),
          uniqueReports: ocrConfigs[i].uniqueReports,
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

  // Reverts
  // TODO: implement revert tests after re-introducing validations

  function test_RepeatAddress_Revert() public {
    address[] memory signers = new address[](1);
    signers[0] = address(1245678);
    address[] memory transmitters = new address[](1);
    transmitters[0] = signers[0];

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: 0,
      configDigest: _getBasicConfigDigest(1, signers, transmitters),
      F: 1,
      uniqueReports: false,
      isSignatureVerificationEnabled: true,
      signers: signers,
      transmitters: transmitters
    });

    vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.InvalidConfig.selector, "repeated transmitter address"));
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  function test_SingerCannotBeZeroAddress_Revert() public {
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
      uniqueReports: false,
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
      uniqueReports: false,
      isSignatureVerificationEnabled: true,
      signers: signers,
      transmitters: transmitters
    });

    vm.expectRevert(MultiOCR3Base.OracleCannotBeZeroAddress.selector);
    s_multiOCR3.setOCR3Configs(ocrConfigs);
  }

  //   function test_OracleOutOfRegister_Revert() public {
  //     address[] memory signers = new address[](10);
  //     address[] memory transmitters = new address[](0);

  //     vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.InvalidConfig.selector, "oracle addresses out of registration"));
  //     s_multiOCR3.setOCR2Config(signers, transmitters, 2, abi.encode(""), 100, abi.encode(""));
  //   }

  //   function test_FTooHigh_Revert() public {
  //     address[] memory signers = new address[](0);
  //     uint8 f = 1;

  //     vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.InvalidConfig.selector, "faulty-oracle f too high"));
  //     s_multiOCR3.setOCR2Config(signers, new address[](0), f, abi.encode(""), 100, abi.encode(""));
  //   }

  //   function test_FMustBePositive_Revert() public {
  //     uint8 f = 0;

  //     vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.InvalidConfig.selector, "f must be positive"));
  //     s_multiOCR3.setOCR2Config(new address[](0), new address[](0), f, abi.encode(""), 100, abi.encode(""));
  //   }

  //   function test_TooManySigners_Revert() public {
  //     address[] memory signers = new address[](32);

  //     vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.InvalidConfig.selector, "too many signers"));
  //     s_multiOCR3.setOCR2Config(signers, new address[](0), 0, abi.encode(""), 100, abi.encode(""));
  //   }
}
