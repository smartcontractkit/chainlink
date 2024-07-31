// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {OCR2BaseNoChecks} from "../../ocr/OCR2BaseNoChecks.sol";
import {OCR2NoChecksHelper} from "../helpers/OCR2NoChecksHelper.sol";
import {OCR2Setup} from "./OCR2Setup.t.sol";

contract OCR2BaseNoChecksSetup is OCR2Setup {
  OCR2NoChecksHelper internal s_OCR2Base;

  bytes32[] internal s_rs;
  bytes32[] internal s_ss;
  bytes32 internal s_rawVs;

  function setUp() public virtual override {
    OCR2Setup.setUp();
    s_OCR2Base = new OCR2NoChecksHelper();
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
}

contract OCR2BaseNoChecks_transmit is OCR2BaseNoChecksSetup {
  bytes32 internal s_configDigest;

  function setUp() public virtual override {
    OCR2BaseNoChecksSetup.setUp();
    bytes memory configBytes = abi.encode("");

    s_configDigest = getBasicConfigDigest(s_f, 0);
    s_OCR2Base.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, configBytes, s_offchainConfigVersion, configBytes
    );
  }

  function test_TransmitSuccess_gas() public {
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
    vm.expectRevert(abi.encodeWithSelector(OCR2BaseNoChecks.ForkedChain.selector, chain1, chain2));
    vm.startPrank(s_valid_transmitters[0]);
    s_OCR2Base.transmit(reportContext, REPORT, s_rs, s_ss, s_rawVs);
  }

  function test_ConfigDigestMismatch_Revert() public {
    bytes32 configDigest;

    bytes32[3] memory reportContext = [configDigest, configDigest, configDigest];

    vm.expectRevert(
      abi.encodeWithSelector(OCR2BaseNoChecks.ConfigDigestMismatch.selector, s_configDigest, configDigest)
    );
    s_OCR2Base.transmit(reportContext, REPORT, new bytes32[](0), new bytes32[](0), s_rawVs);
  }

  function test_UnAuthorizedTransmitter_Revert() public {
    bytes32[3] memory reportContext = [s_configDigest, s_configDigest, s_configDigest];
    bytes32[] memory rs = new bytes32[](3);
    bytes32[] memory ss = new bytes32[](3);

    vm.expectRevert(OCR2BaseNoChecks.UnauthorizedTransmitter.selector);
    s_OCR2Base.transmit(reportContext, REPORT, rs, ss, s_rawVs);
  }
}

contract OCR2BaseNoChecks_setOCR2Config is OCR2BaseNoChecksSetup {
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

  function test_SetConfigSuccess_gas() public {
    vm.pauseGasMetering();
    bytes memory configBytes = abi.encode("");
    uint32 configCount = 0;

    bytes32 configDigest = getBasicConfigDigest(s_f, configCount++);

    address[] memory transmitters = s_OCR2Base.getTransmitters();
    assertEq(0, transmitters.length);

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

    s_OCR2Base.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, configBytes, s_offchainConfigVersion, configBytes
    );

    transmitters = s_OCR2Base.getTransmitters();
    assertEq(s_valid_transmitters, transmitters);

    configDigest = getBasicConfigDigest(s_f, configCount++);

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
    s_OCR2Base.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, configBytes, s_offchainConfigVersion, configBytes
    );
  }

  // Reverts
  function test_RepeatAddress_Revert() public {
    address[] memory signers = new address[](4);
    address[] memory transmitters = new address[](4);
    transmitters[0] = address(1245678);
    transmitters[1] = address(1245678);
    transmitters[2] = address(1245678);
    transmitters[3] = address(1245678);

    vm.expectRevert(
      abi.encodeWithSelector(
        OCR2BaseNoChecks.InvalidConfig.selector, OCR2BaseNoChecks.InvalidConfigErrorType.REPEATED_ORACLE_ADDRESS
      )
    );
    s_OCR2Base.setOCR2Config(signers, transmitters, 1, abi.encode(""), 100, abi.encode(""));
  }

  function test_FMustBePositive_Revert() public {
    uint8 f = 0;

    vm.expectRevert(
      abi.encodeWithSelector(
        OCR2BaseNoChecks.InvalidConfig.selector, OCR2BaseNoChecks.InvalidConfigErrorType.F_MUST_BE_POSITIVE
      )
    );
    s_OCR2Base.setOCR2Config(new address[](0), new address[](0), f, abi.encode(""), 100, abi.encode(""));
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

    vm.expectRevert(OCR2BaseNoChecks.OracleCannotBeZeroAddress.selector);
    s_OCR2Base.setOCR2Config(signers, transmitters, uint8(f), abi.encode(""), 100, abi.encode(""));
  }

  function test_TooManyTransmitter_Revert() public {
    address[] memory transmitters = new address[](100);

    vm.expectRevert(
      abi.encodeWithSelector(
        OCR2BaseNoChecks.InvalidConfig.selector, OCR2BaseNoChecks.InvalidConfigErrorType.TOO_MANY_TRANSMITTERS
      )
    );
    s_OCR2Base.setOCR2Config(new address[](0), transmitters, 0, abi.encode(""), 100, abi.encode(""));
  }
}
