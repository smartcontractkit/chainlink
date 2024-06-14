// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {OCR2Abstract} from "../../ocr/OCR2Abstract.sol";
import {OCR2Base} from "../../ocr/OCR2Base.sol";
import {OCR2Helper} from "../helpers/OCR2Helper.sol";
import {OCR2Setup} from "./OCR2Setup.t.sol";

contract OCR2BaseSetup is OCR2Setup {
  OCR2Helper internal s_OCR2Base;

  bytes32[] internal s_rs;
  bytes32[] internal s_ss;
  bytes32 internal s_rawVs;

  uint40 internal s_latestEpochAndRound;

  function setUp() public virtual override {
    OCR2Setup.setUp();
    s_OCR2Base = new OCR2Helper();

    bytes32 testReportDigest = getTestReportDigest();

    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](2);
    uint8[] memory vs = new uint8[](2);

    // Calculate signatures
    (vs[0], rs[0], ss[0]) = vm.sign(PRIVATE0, testReportDigest);
    (vs[1], rs[1], ss[1]) = vm.sign(PRIVATE1, testReportDigest);

    s_rs = rs;
    s_ss = ss;
    s_rawVs = bytes32(bytes1(vs[0] - 27)) | (bytes32(bytes1(vs[1] - 27)) >> 8);
  }

  function getBasicConfigDigest(uint8 f, uint64 currentConfigCount) internal view returns (bytes32) {
    bytes memory configBytes = abi.encode("");
    return s_OCR2Base.configDigestFromConfigData(
      block.chainid,
      address(s_OCR2Base),
      currentConfigCount + 1,
      s_valid_signers,
      s_valid_transmitters,
      f,
      configBytes,
      s_offchainConfigVersion,
      configBytes
    );
  }

  function getTestReportDigest() internal view returns (bytes32) {
    bytes32 configDigest = getBasicConfigDigest(s_f, 0);
    bytes32[3] memory reportContext = [configDigest, configDigest, configDigest];
    return keccak256(abi.encodePacked(keccak256(REPORT), reportContext));
  }

  function getBasicConfigDigest(
    address contractAddress,
    uint8 f,
    uint64 currentConfigCount,
    bytes memory onchainConfig
  ) internal view returns (bytes32) {
    return s_OCR2Base.configDigestFromConfigData(
      block.chainid,
      contractAddress,
      currentConfigCount + 1,
      s_valid_signers,
      s_valid_transmitters,
      f,
      onchainConfig,
      s_offchainConfigVersion,
      abi.encode("")
    );
  }
}

contract OCR2Base_transmit is OCR2BaseSetup {
  bytes32 internal s_configDigest;

  function setUp() public virtual override {
    OCR2BaseSetup.setUp();
    bytes memory configBytes = abi.encode("");

    s_configDigest = getBasicConfigDigest(s_f, 0);
    s_OCR2Base.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, configBytes, s_offchainConfigVersion, configBytes
    );
  }

  function test_Transmit2SignersSuccess_gas() public {
    vm.pauseGasMetering();
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];

    vm.startPrank(s_valid_transmitters[0]);
    vm.resumeGasMetering();
    s_OCR2Base.transmit(reportContext, REPORT, s_rs, s_ss, s_rawVs);
  }

  // Reverts

  function test_ForkedChain_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];

    uint256 chain1 = block.chainid;
    uint256 chain2 = chain1 + 1;
    vm.chainId(chain2);
    vm.expectRevert(abi.encodeWithSelector(OCR2Base.ForkedChain.selector, chain1, chain2));
    vm.startPrank(s_valid_transmitters[0]);
    s_OCR2Base.transmit(reportContext, REPORT, s_rs, s_ss, s_rawVs);
  }

  function test_WrongNumberOfSignatures_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];

    vm.expectRevert(OCR2Base.WrongNumberOfSignatures.selector);
    s_OCR2Base.transmit(reportContext, REPORT, new bytes32[](0), new bytes32[](0), s_rawVs);
  }

  function test_ConfigDigestMismatch_Revert() public {
    bytes32 configDigest;
    bytes32[3] memory reportContext = [configDigest, configDigest, configDigest];

    vm.expectRevert(abi.encodeWithSelector(OCR2Base.ConfigDigestMismatch.selector, s_configDigest, configDigest));
    s_OCR2Base.transmit(reportContext, REPORT, new bytes32[](0), new bytes32[](0), s_rawVs);
  }

  function test_SignatureOutOfRegistration_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];

    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](1);

    vm.expectRevert(OCR2Base.SignaturesOutOfRegistration.selector);
    s_OCR2Base.transmit(reportContext, REPORT, rs, ss, s_rawVs);
  }

  function test_UnAuthorizedTransmitter_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](2);

    vm.expectRevert(OCR2Base.UnauthorizedTransmitter.selector);
    s_OCR2Base.transmit(reportContext, REPORT, rs, ss, s_rawVs);
  }

  function test_NonUniqueSignature_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];
    bytes32[] memory rs = s_rs;
    bytes32[] memory ss = s_ss;

    rs[1] = rs[0];
    ss[1] = ss[0];
    // Need to reset the rawVs to be valid
    bytes32 rawVs = bytes32(bytes1(uint8(28) - 27)) | (bytes32(bytes1(uint8(28) - 27)) >> 8);

    vm.startPrank(s_valid_transmitters[0]);
    vm.expectRevert(OCR2Base.NonUniqueSignatures.selector);
    s_OCR2Base.transmit(reportContext, REPORT, rs, ss, rawVs);
  }

  function test_UnauthorizedSigner_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];
    bytes32[] memory rs = new bytes32[](2);
    rs[0] = s_configDigest;
    bytes32[] memory ss = rs;

    vm.startPrank(s_valid_transmitters[0]);
    vm.expectRevert(OCR2Base.UnauthorizedSigner.selector);
    s_OCR2Base.transmit(reportContext, REPORT, rs, ss, s_rawVs);
  }
}

contract OCR2Base_setOCR2Config is OCR2BaseSetup {
  function test_SetConfigSuccess_gas() public {
    vm.pauseGasMetering();
    bytes memory configBytes = abi.encode("");
    uint32 configCount = 0;

    bytes32 configDigest = getBasicConfigDigest(s_f, configCount++);

    address[] memory transmitters = s_OCR2Base.getTransmitters();
    assertEq(0, transmitters.length);

    vm.expectEmit();
    emit OCR2Abstract.ConfigSet(
      0,
      configDigest,
      configCount,
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      configBytes,
      s_offchainConfigVersion,
      configBytes
    );

    s_OCR2Base.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, configBytes, s_offchainConfigVersion, configBytes
    );

    transmitters = s_OCR2Base.getTransmitters();
    assertEq(s_valid_transmitters, transmitters);

    configDigest = getBasicConfigDigest(s_f, configCount++);

    vm.expectEmit();
    emit OCR2Abstract.ConfigSet(
      uint32(block.number),
      configDigest,
      configCount,
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      configBytes,
      s_offchainConfigVersion,
      configBytes
    );
    vm.resumeGasMetering();
    s_OCR2Base.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, configBytes, s_offchainConfigVersion, configBytes
    );
  }

  // Reverts
  function test_RepeatAddress_Revert() public {
    address[] memory signers = new address[](10);
    signers[0] = address(1245678);
    address[] memory transmitters = new address[](10);
    transmitters[0] = signers[0];

    vm.expectRevert(
      abi.encodeWithSelector(OCR2Base.InvalidConfig.selector, OCR2Base.InvalidConfigErrorType.REPEATED_ORACLE_ADDRESS)
    );
    s_OCR2Base.setOCR2Config(signers, transmitters, 2, abi.encode(""), 100, abi.encode(""));
  }

  function test_SingerCannotBeZeroAddress_Revert() public {
    uint256 f = 1;
    address[] memory signers = new address[](3 * f + 1);
    address[] memory transmitters = new address[](3 * f + 1);
    for (uint160 i = 0; i < 3 * f + 1; ++i) {
      signers[i] = address(i + 1);
      transmitters[i] = address(i + 1000);
    }

    signers[0] = address(0);

    vm.expectRevert(OCR2Base.OracleCannotBeZeroAddress.selector);
    s_OCR2Base.setOCR2Config(signers, transmitters, uint8(f), abi.encode(""), 100, abi.encode(""));
  }

  function test_TransmitterCannotBeZeroAddress_Revert() public {
    uint256 f = 1;
    address[] memory signers = new address[](3 * f + 1);
    address[] memory transmitters = new address[](3 * f + 1);
    for (uint160 i = 0; i < 3 * f + 1; ++i) {
      signers[i] = address(i + 1);
      transmitters[i] = address(i + 1000);
    }

    transmitters[0] = address(0);

    vm.expectRevert(OCR2Base.OracleCannotBeZeroAddress.selector);
    s_OCR2Base.setOCR2Config(signers, transmitters, uint8(f), abi.encode(""), 100, abi.encode(""));
  }

  function test_OracleOutOfRegister_Revert() public {
    address[] memory signers = new address[](10);
    address[] memory transmitters = new address[](0);

    vm.expectRevert(
      abi.encodeWithSelector(
        OCR2Base.InvalidConfig.selector, OCR2Base.InvalidConfigErrorType.NUM_SIGNERS_NOT_NUM_TRANSMITTERS
      )
    );
    s_OCR2Base.setOCR2Config(signers, transmitters, 2, abi.encode(""), 100, abi.encode(""));
  }

  function test_FTooHigh_Revert() public {
    address[] memory signers = new address[](0);
    uint8 f = 1;

    vm.expectRevert(abi.encodeWithSelector(OCR2Base.InvalidConfig.selector, OCR2Base.InvalidConfigErrorType.F_TOO_HIGH));
    s_OCR2Base.setOCR2Config(signers, new address[](0), f, abi.encode(""), 100, abi.encode(""));
  }

  function test_FMustBePositive_Revert() public {
    uint8 f = 0;

    vm.expectRevert(
      abi.encodeWithSelector(OCR2Base.InvalidConfig.selector, OCR2Base.InvalidConfigErrorType.F_MUST_BE_POSITIVE)
    );
    s_OCR2Base.setOCR2Config(new address[](0), new address[](0), f, abi.encode(""), 100, abi.encode(""));
  }

  function test_TooManySigners_Revert() public {
    address[] memory signers = new address[](32);

    vm.expectRevert(
      abi.encodeWithSelector(OCR2Base.InvalidConfig.selector, OCR2Base.InvalidConfigErrorType.TOO_MANY_SIGNERS)
    );
    s_OCR2Base.setOCR2Config(signers, new address[](0), 0, abi.encode(""), 100, abi.encode(""));
  }
}
