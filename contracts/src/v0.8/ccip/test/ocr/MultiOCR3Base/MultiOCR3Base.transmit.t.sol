// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {MultiOCR3Base} from "../../../ocr/MultiOCR3Base.sol";
import {MultiOCR3BaseSetup} from "./MultiOCR3BaseSetup.t.sol";

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

  function test_TransmitSigners_gas() public {
    vm.pauseGasMetering();
    bytes32[2] memory reportContext = [s_configDigest1, s_configDigest1];

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

  function test_TransmitWithoutSignatureVerification_gas() public {
    vm.pauseGasMetering();
    bytes32[2] memory reportContext = [s_configDigest3, s_configDigest3];

    s_multiOCR3.setTransmitOcrPluginType(2);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(2, s_configDigest3, uint64(uint256(s_configDigest3)));

    vm.startPrank(s_validTransmitters[0]);
    vm.resumeGasMetering();
    s_multiOCR3.transmitWithoutSignatures(reportContext, REPORT);
  }

  function testFuzz_TransmitSignersWithSignatures_Success(uint8 F, uint64 randomAddressOffset) public {
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

    bytes32[2] memory reportContext = [s_configDigest1, s_configDigest1];

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
  function test_RevertWhen_ForkedChain() public {
    bytes32[2] memory reportContext = [s_configDigest1, s_configDigest1];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, REPORT, reportContext, 2);

    s_multiOCR3.setTransmitOcrPluginType(0);

    uint256 chain1 = uint200(block.chainid);
    uint256 chain2 = chain1 + 1;
    vm.chainId(chain2);
    vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.ForkedChain.selector, chain1, chain2));

    vm.startPrank(s_validTransmitters[0]);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, rawVs);
  }

  function test_RevertWhen_ZeroSignatures() public {
    bytes32[2] memory reportContext = [s_configDigest1, s_configDigest1];

    s_multiOCR3.setTransmitOcrPluginType(0);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.WrongNumberOfSignatures.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, new bytes32[](0), new bytes32[](0), bytes32(""));
  }

  function test_RevertWhen_TooManySignatures() public {
    bytes32[2] memory reportContext = [s_configDigest1, s_configDigest1];

    // 1 signature too many
    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, REPORT, reportContext, 6);

    s_multiOCR3.setTransmitOcrPluginType(1);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.WrongNumberOfSignatures.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, rawVs);
  }

  function test_RevertWhen_InsufficientSignatures() public {
    bytes32[2] memory reportContext = [s_configDigest1, s_configDigest1];

    // Missing 1 signature for unique report
    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, REPORT, reportContext, 4);

    s_multiOCR3.setTransmitOcrPluginType(1);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.WrongNumberOfSignatures.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, rawVs);
  }

  function test_RevertWhen_ConfigDigestMismatch() public {
    bytes32 configDigest;
    bytes32[2] memory reportContext = [configDigest, configDigest];

    (,,, bytes32 rawVs) = _getSignaturesForDigest(s_validSignerKeys, REPORT, reportContext, 2);

    s_multiOCR3.setTransmitOcrPluginType(0);

    vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.ConfigDigestMismatch.selector, s_configDigest1, configDigest));
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, new bytes32[](0), new bytes32[](0), rawVs);
  }

  function test_RevertWhen_SignatureOutOfRegistration() public {
    bytes32[2] memory reportContext = [s_configDigest1, s_configDigest1];

    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](1);

    s_multiOCR3.setTransmitOcrPluginType(0);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.SignaturesOutOfRegistration.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, bytes32(""));
  }

  function test_RevertWhen_UnAuthorizedTransmitter() public {
    bytes32[2] memory reportContext = [s_configDigest1, s_configDigest1];
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](2);

    s_multiOCR3.setTransmitOcrPluginType(0);

    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, bytes32(""));
  }

  function test_RevertWhen_NonUniqueSignature() public {
    bytes32[2] memory reportContext = [s_configDigest1, s_configDigest1];

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

  function test_RevertWhen_UnauthorizedSigner() public {
    bytes32[2] memory reportContext = [s_configDigest1, s_configDigest1];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, REPORT, reportContext, 2);

    rs[0] = s_configDigest1;
    ss = rs;

    s_multiOCR3.setTransmitOcrPluginType(0);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert(MultiOCR3Base.UnauthorizedSigner.selector);
    s_multiOCR3.transmitWithSignatures(reportContext, REPORT, rs, ss, rawVs);
  }

  function test_RevertWhen_UnconfiguredPlugin() public {
    bytes32 configDigest;
    bytes32[2] memory reportContext = [configDigest, configDigest];

    s_multiOCR3.setTransmitOcrPluginType(42);

    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_multiOCR3.transmitWithoutSignatures(reportContext, REPORT);
  }

  function test_RevertWhen_TransmitWithLessCalldataArgs() public {
    bytes32[2] memory reportContext = [s_configDigest1, s_configDigest1];

    s_multiOCR3.setTransmitOcrPluginType(0);

    // The transmit should fail, since we are trying to transmit without signatures when signatures are enabled
    vm.startPrank(s_validTransmitters[1]);

    // report length + function selector + report length + abiencoded location of report value + report context words
    uint256 receivedLength = REPORT.length + 4 + 4 * 32;
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

  function test_RevertWhen_TransmitWithExtraCalldataArgs() public {
    bytes32[2] memory reportContext = [s_configDigest1, s_configDigest1];
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](2);

    s_multiOCR3.setTransmitOcrPluginType(2);

    // The transmit should fail, since we are trying to transmit with signatures when signatures are disabled
    vm.startPrank(s_validTransmitters[1]);

    // dynamic length + function selector + report length + abiencoded location of report value + report context words
    // rawVs value, lengths of rs, ss, and start locations of rs & ss -> 5 words
    uint256 receivedLength = REPORT.length + 4 + (4 * 32) + (5 * 32) + (2 * 32) + (2 * 32);
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
