// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {OCR3Setup} from "./OCR3Setup.t.sol";
import {OCR3Base} from "../../ocr/OCR3Base.sol";
import {OCR3Helper} from "../helpers/OCR3Helper.sol";

contract OCR3BaseSetup is OCR3Setup {
  event ConfigSet(
    uint32 previousConfigBlockNumber,
    bytes32 configDigest,
    uint64 configCount,
    address[] signers,
    address[] transmitters,
    uint8 f,
    bytes onchainConfig,
    uint64 offchainConfigVersion,
    bytes offchainConfig
  );

  OCR3Helper internal s_OCR3Base;

  bytes32[] internal s_rs;
  bytes32[] internal s_ss;
  bytes32 internal s_rawVs;

  uint40 internal s_latestEpochAndRound;

  function setUp() public virtual override {
    OCR3Setup.setUp();
    s_OCR3Base = new OCR3Helper();

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
    return
      s_OCR3Base.configDigestFromConfigData(
        block.chainid,
        address(s_OCR3Base),
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
    return
      s_OCR3Base.configDigestFromConfigData(
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

contract OCR3Base_transmit is OCR3BaseSetup {
  bytes32 internal s_configDigest;

  function setUp() public virtual override {
    OCR3BaseSetup.setUp();
    bytes memory configBytes = abi.encode("");

    s_configDigest = getBasicConfigDigest(s_f, 0);
    s_OCR3Base.setOCR3Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      configBytes,
      s_offchainConfigVersion,
      configBytes
    );
  }

  function testTransmit2SignersSuccess_gas() public {
    vm.pauseGasMetering();
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];

    vm.startPrank(s_valid_transmitters[0]);
    vm.resumeGasMetering();
    s_OCR3Base.transmit(reportContext, REPORT, s_rs, s_ss, s_rawVs);
  }

  // Reverts

  function testNonIncreasingSequenceNumberReverts() public {
    bytes32[3] memory reportContext = [s_configDigest, bytes32(uint256(0)) /* sequence number */, s_configDigest];

    vm.expectRevert(abi.encodeWithSelector(OCR3Base.NonIncreasingSequenceNumber.selector, 0, 0));
    s_OCR3Base.transmit(reportContext, REPORT, s_rs, s_ss, s_rawVs);
  }

  function testForkedChainReverts() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];

    uint256 chain1 = block.chainid;
    uint256 chain2 = chain1 + 1;
    vm.chainId(chain2);
    vm.expectRevert(abi.encodeWithSelector(OCR3Base.ForkedChain.selector, chain1, chain2));
    vm.startPrank(s_valid_transmitters[0]);
    s_OCR3Base.transmit(reportContext, REPORT, s_rs, s_ss, s_rawVs);
  }

  function testWrongNumberOfSignaturesReverts() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];

    vm.expectRevert(OCR3Base.WrongNumberOfSignatures.selector);
    s_OCR3Base.transmit(reportContext, REPORT, new bytes32[](0), new bytes32[](0), s_rawVs);
  }

  function testConfigDigestMismatchReverts() public {
    bytes32 configDigest;
    bytes32[3] memory reportContext = [configDigest, bytes32(uint256(1)) /* sequence number */, configDigest];

    vm.expectRevert(abi.encodeWithSelector(OCR3Base.ConfigDigestMismatch.selector, s_configDigest, configDigest));
    s_OCR3Base.transmit(reportContext, REPORT, new bytes32[](0), new bytes32[](0), s_rawVs);
  }

  function testSignatureOutOfRegistrationReverts() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];

    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](1);

    vm.expectRevert(OCR3Base.SignaturesOutOfRegistration.selector);
    s_OCR3Base.transmit(reportContext, REPORT, rs, ss, s_rawVs);
  }

  function testUnAuthorizedTransmitterReverts() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](2);

    vm.expectRevert(OCR3Base.UnauthorizedTransmitter.selector);
    s_OCR3Base.transmit(reportContext, REPORT, rs, ss, s_rawVs);
  }

  function testNonUniqueSignatureReverts() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];
    bytes32[] memory rs = s_rs;
    bytes32[] memory ss = s_ss;

    rs[1] = rs[0];
    ss[1] = ss[0];
    // Need to reset the rawVs to be valid
    bytes32 rawVs = bytes32(bytes1(uint8(28) - 27)) | (bytes32(bytes1(uint8(28) - 27)) >> 8);

    vm.startPrank(s_valid_transmitters[0]);
    vm.expectRevert(OCR3Base.NonUniqueSignatures.selector);
    s_OCR3Base.transmit(reportContext, REPORT, rs, ss, rawVs);
  }

  function testUnauthorizedSignerReverts() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];
    bytes32[] memory rs = new bytes32[](2);
    rs[0] = s_configDigest;
    bytes32[] memory ss = rs;

    vm.startPrank(s_valid_transmitters[0]);
    vm.expectRevert(OCR3Base.UnauthorizedSigner.selector);
    s_OCR3Base.transmit(reportContext, REPORT, rs, ss, s_rawVs);
  }
}

contract OCR3Base_setOCR3Config is OCR3BaseSetup {
  function testSetConfigSuccess() public {
    vm.pauseGasMetering();
    bytes memory configBytes = abi.encode("");
    uint32 configCount = 0;

    bytes32 configDigest = getBasicConfigDigest(s_f, configCount++);

    address[] memory transmitters = s_OCR3Base.getTransmitters();
    assertEq(0, transmitters.length);

    s_OCR3Base.setLatestSeqNum(3);
    uint64 seqNum = s_OCR3Base.latestSequenceNumber();
    assertEq(seqNum, 3);

    vm.expectEmit();
    emit ConfigSet(
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

    s_OCR3Base.setOCR3Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      configBytes,
      s_offchainConfigVersion,
      configBytes
    );

    transmitters = s_OCR3Base.getTransmitters();
    assertEq(s_valid_transmitters, transmitters);

    configDigest = getBasicConfigDigest(s_f, configCount++);

    seqNum = s_OCR3Base.latestSequenceNumber();
    assertEq(seqNum, 0);

    vm.expectEmit();
    emit ConfigSet(
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
    s_OCR3Base.setOCR3Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      configBytes,
      s_offchainConfigVersion,
      configBytes
    );
  }

  // Reverts
  function testRepeatAddressReverts() public {
    address[] memory signers = new address[](10);
    signers[0] = address(1245678);
    address[] memory transmitters = new address[](10);
    transmitters[0] = signers[0];

    vm.expectRevert(abi.encodeWithSelector(OCR3Base.InvalidConfig.selector, "repeated transmitter address"));
    s_OCR3Base.setOCR3Config(signers, transmitters, 2, abi.encode(""), 100, abi.encode(""));
  }

  function testSignerCannotBeZeroAddressReverts() public {
    uint256 f = 1;
    address[] memory signers = new address[](3 * f + 1);
    address[] memory transmitters = new address[](3 * f + 1);
    for (uint160 i = 0; i < 3 * f + 1; ++i) {
      signers[i] = address(i + 1);
      transmitters[i] = address(i + 1000);
    }

    signers[0] = address(0);

    vm.expectRevert(OCR3Base.OracleCannotBeZeroAddress.selector);
    s_OCR3Base.setOCR3Config(signers, transmitters, uint8(f), abi.encode(""), 100, abi.encode(""));
  }

  function testTransmitterCannotBeZeroAddressReverts() public {
    uint256 f = 1;
    address[] memory signers = new address[](3 * f + 1);
    address[] memory transmitters = new address[](3 * f + 1);
    for (uint160 i = 0; i < 3 * f + 1; ++i) {
      signers[i] = address(i + 1);
      transmitters[i] = address(i + 1000);
    }

    transmitters[0] = address(0);

    vm.expectRevert(OCR3Base.OracleCannotBeZeroAddress.selector);
    s_OCR3Base.setOCR3Config(signers, transmitters, uint8(f), abi.encode(""), 100, abi.encode(""));
  }

  function testOracleOutOfRegisterReverts() public {
    address[] memory signers = new address[](10);
    address[] memory transmitters = new address[](0);

    vm.expectRevert(abi.encodeWithSelector(OCR3Base.InvalidConfig.selector, "oracle addresses out of registration"));
    s_OCR3Base.setOCR3Config(signers, transmitters, 2, abi.encode(""), 100, abi.encode(""));
  }

  function testFTooHighReverts() public {
    address[] memory signers = new address[](0);
    uint8 f = 1;

    vm.expectRevert(abi.encodeWithSelector(OCR3Base.InvalidConfig.selector, "faulty-oracle f too high"));
    s_OCR3Base.setOCR3Config(signers, new address[](0), f, abi.encode(""), 100, abi.encode(""));
  }

  function testFMustBePositiveReverts() public {
    uint8 f = 0;

    vm.expectRevert(abi.encodeWithSelector(OCR3Base.InvalidConfig.selector, "f must be positive"));
    s_OCR3Base.setOCR3Config(new address[](0), new address[](0), f, abi.encode(""), 100, abi.encode(""));
  }

  function testTooManySignersReverts() public {
    address[] memory signers = new address[](32);

    vm.expectRevert(abi.encodeWithSelector(OCR3Base.InvalidConfig.selector, "too many signers"));
    s_OCR3Base.setOCR3Config(signers, new address[](0), 0, abi.encode(""), 100, abi.encode(""));
  }
}
