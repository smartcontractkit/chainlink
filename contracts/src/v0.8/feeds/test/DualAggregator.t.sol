// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";

import {DualAggregator} from "../DualAggregator.sol";

import {AccessControllerInterface} from "../../shared/interfaces/AccessControllerInterface.sol";
import {AggregatorValidatorInterface} from "../../shared/interfaces/AggregatorValidatorInterface.sol";
import {LinkTokenInterface} from "../../shared/interfaces/LinkTokenInterface.sol";
import {LinkToken} from "../../shared/token/ERC677/LinkToken.sol";

import {AggregatorValidator} from "./testhelpers/AggregatorValidator.t.sol";
import {DualAggregatorHelper} from "./testhelpers/DualAggregatorHelper.sol";
import {ReportGenerator} from "./testhelpers/ReportGenerator.t.sol";

contract DualAggregatorHarness is DualAggregatorHelper {
  constructor(
    LinkTokenInterface link,
    int192 minAnswer_,
    int192 maxAnswer_,
    AccessControllerInterface billingAccessController,
    AccessControllerInterface requesterAccessController,
    uint8 decimals_,
    string memory description_,
    address secondaryProxy_,
    uint32 cutoffTime_,
    uint32 maxSyncIterations_
  )
    DualAggregatorHelper(
      link,
      minAnswer_,
      maxAnswer_,
      billingAccessController,
      requesterAccessController,
      decimals_,
      description_,
      secondaryProxy_,
      cutoffTime_,
      maxSyncIterations_
    )
  {}

  // helper function to define the latest round ids
  function setLatestRoundIds(uint32 _latestAggregatorRoundId, uint32 _latestSecondaryRoundId) public {
    s_hotVars.latestAggregatorRoundId = _latestAggregatorRoundId;
    s_hotVars.latestSecondaryRoundId = _latestSecondaryRoundId;
  }

  // helper function to add a transmission without depending on transmit()
  function setTransmission(
    uint32 _roundId,
    int192 _answer,
    uint32 _observationsTimestamp,
    uint32 _recordedTimestamp
  ) public {
    s_transmissions[_roundId] = Transmission({
      answer: _answer,
      observationsTimestamp: _observationsTimestamp,
      recordedTimestamp: _recordedTimestamp
    });
  }

  // helper function to inject transmissions, answer starts at 10 and increases by 1 every round,
  // the observations timestamp starts at 1 and increases by 5 every round, the recorded timestamp
  // starts at 5 and increases by 5 every round
  function injectTransmissions(
    uint32 amount
  ) public {
    for (uint32 i = 0; i < amount; i++) {
      uint32 roundId = i + 1; // starts at one, increases by 1 every round [1, 2, 3, 4, 5, ...]
      int192 answer = int32(i) + 10; // starts at 10, increases by 1 every round [10, 11, 12, 13, 14, ...]
      uint32 observationsTimestamp = i * 5 + 1; // starts at 1, increases by 5 every round [1, 6, 11, 16, 21, ...]
      uint32 recordedTimestamp = i * 5 + 5; // starts at 5, increases by 5 every round [5, 10, 15, 20, 25, ...]

      setTransmission(roundId, answer, observationsTimestamp, recordedTimestamp);
    }
  }

  // helper function to define the latestSecondary
  function isLatestSecondary(
    bool isSecondary
  ) public {
    s_hotVars.isLatestSecondary = isSecondary;
  }
}

contract DualAggregatorBaseTest is Test {
  uint256 internal constant MAX_NUM_ORACLES = 31;

  address internal constant BILLING_ACCESS_CONTROLLER_ADDRESS = address(100);
  address internal constant REQUESTER_ACCESS_CONTROLLER_ADDRESS = address(101);
  address internal constant SECONDARY_PROXY = address(102);

  int192 internal constant MIN_ANSWER = 0;
  int192 internal constant MAX_ANSWER = 100;

  LinkToken internal s_link;
  DualAggregatorHarness internal s_aggregator;

  function setUp() public virtual {
    s_link = new LinkToken();

    AccessControllerInterface _billingAccessController = AccessControllerInterface(BILLING_ACCESS_CONTROLLER_ADDRESS);
    AccessControllerInterface _requesterAccessController =
      AccessControllerInterface(REQUESTER_ACCESS_CONTROLLER_ADDRESS);

    s_aggregator = new DualAggregatorHarness(
      LinkTokenInterface(address(s_link)),
      MIN_ANSWER,
      MAX_ANSWER,
      _billingAccessController,
      _requesterAccessController,
      18,
      "TEST",
      SECONDARY_PROXY,
      0,
      10
    );
  }

  function _changePrank(
    address _prank
  ) internal {
    vm.stopPrank();
    vm.startPrank(_prank);
  }
}

contract ConfiguredDualAggregatorBaseTest is DualAggregatorBaseTest {
  address[] internal s_signers;
  address[] internal s_transmitters;
  uint8 internal s_f;
  bytes internal s_onchainConfig;
  uint64 internal s_offchainConfigVersion;
  bytes internal s_offchainConfig;
  bytes32 internal s_configDigest;
  ReportGenerator internal s_reportGenerator;

  function setUp() public virtual override {
    super.setUp();

    s_f = 1;
    s_onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
    s_offchainConfigVersion = 1;
    s_offchainConfig = "1";

    uint256[] memory privateKeys = new uint256[](MAX_NUM_ORACLES);
    for (uint256 i = 0; i < MAX_NUM_ORACLES; i++) {
      uint256 privateKey = uint256(keccak256(abi.encodePacked(i, "oracle-generator-seed")));
      address publicKey = vm.addr(privateKey);
      privateKeys[i] = privateKey;
      s_transmitters.push(publicKey);
      s_signers.push(publicKey);
    }

    s_aggregator.setConfig(s_signers, s_transmitters, s_f, s_onchainConfig, s_offchainConfigVersion, s_offchainConfig);
    s_configDigest = s_aggregator.configDigestFromConfigData(
      block.chainid,
      address(s_aggregator),
      1,
      s_signers,
      s_transmitters,
      s_f,
      s_onchainConfig,
      s_offchainConfigVersion,
      s_offchainConfig
    );
    s_reportGenerator = new ReportGenerator(s_aggregator, privateKeys, s_configDigest, s_f);
  }
}

contract RoundDataDualAggregatorBaseTest is ConfiguredDualAggregatorBaseTest {
  function setDualAggregatorBase(
    uint256 startingTime,
    uint32 cutoffTime,
    uint32 latestPrimaryRound,
    uint32 latestSecondaryRound,
    bool isLatestSecondary
  ) internal {
    s_aggregator.injectTransmissions(6);
    s_aggregator.setLatestRoundIds(latestPrimaryRound, latestSecondaryRound);
    s_aggregator.setCutoffTime(cutoffTime);
    s_aggregator.isLatestSecondary(isLatestSecondary);
    vm.warp(startingTime);
  }
}

contract Constructor is DualAggregatorBaseTest {
  function test_constructor() public view {
    assertEq(s_aggregator.i_minAnswer(), MIN_ANSWER, "minAnswer not set correctly");
    assertEq(s_aggregator.i_maxAnswer(), MAX_ANSWER, "maxAnswer not set correctly");
    assertEq(s_aggregator.decimals(), 18, "decimals not set correctly");
    assertEq(s_aggregator.description(), "TEST", "description not set correctly");
    assertEq(
      address(s_aggregator.getBillingAccessController()),
      BILLING_ACCESS_CONTROLLER_ADDRESS,
      "billing access controller not set correctly"
    );
    assertEq(
      address(s_aggregator.getRequesterAccessController()),
      REQUESTER_ACCESS_CONTROLLER_ADDRESS,
      "requester access controller not set correctly"
    );
    assertEq(address(s_aggregator.getLinkToken()), address(s_link), "link token not set correctly");
  }
}

contract SetConfig is DualAggregatorBaseTest {
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

  function test_RevertIf_SignersTooLong() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES + 1);
    address[] memory transmitters = new address[](31);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(DualAggregator.TooManyOracles.selector);

    s_aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_OracleLengthMismatch() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES - 1);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(DualAggregator.OracleLengthMismatch.selector);

    s_aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_fTooHigh() public {
    address[] memory signers = new address[](1);
    address[] memory transmitters = new address[](1);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(DualAggregator.FaultyOracleFTooHigh.selector);

    s_aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_fNotPositive() public {
    address[] memory signers = new address[](1);
    address[] memory transmitters = new address[](1);
    uint8 f = 0;
    bytes memory onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(DualAggregator.FMustBePositive.selector);

    s_aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_onchainConfigInvalid() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(DualAggregator.InvalidOnChainConfig.selector);

    s_aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_RepeatedSigner() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES);
    uint8 f = 1;
    bytes memory onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    for (uint256 i = 0; i < MAX_NUM_ORACLES; i++) {
      transmitters[i] = address(uint160(2000 + i));
    }

    vm.expectRevert(DualAggregator.RepeatedSignerAddress.selector);

    s_aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_RepeatedTransmitter() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES);
    uint8 f = 1;
    bytes memory onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    for (uint256 i = 0; i < MAX_NUM_ORACLES; i++) {
      signers[i] = address(uint160(1000 + i));
    }

    vm.expectRevert(DualAggregator.RepeatedTransmitterAddress.selector);

    s_aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_SetsTheConfigAndEmitsEvent() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES);
    uint8 f = 1;
    bytes memory onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    for (uint256 i = 0; i < MAX_NUM_ORACLES; i++) {
      signers[i] = vm.addr(uint160(1000 + i));
      transmitters[i] = vm.addr(uint160(2000 + i));
    }

    bytes32 configDigest = s_aggregator.configDigestFromConfigData(
      block.chainid,
      address(s_aggregator),
      1,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    vm.expectEmit();
    emit ConfigSet(0, configDigest, 1, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
    s_aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);

    (uint32 configCount, uint32 blockNumber, bytes32 digest) = s_aggregator.latestConfigDetails();

    assertEq(configCount, 1, "config count not incremented");
    assertEq(blockNumber, block.number, "block number is wrong");
    assertEq(digest, configDigest, "configDigest is not correct");
  }

  function test_ShouldBeAbleToRemoveSigners() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES);
    uint8 f = 1;
    bytes memory onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    for (uint256 i = 0; i < MAX_NUM_ORACLES; i++) {
      signers[i] = vm.addr(uint160(1000 + i));
      transmitters[i] = vm.addr(uint160(2000 + i));
    }

    bytes32 configDigest = s_aggregator.configDigestFromConfigData(
      block.chainid,
      address(s_aggregator),
      1,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    uint32 blockNumber = uint32(block.number);

    vm.expectEmit();
    emit ConfigSet(0, configDigest, 1, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
    s_aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);

    // remove half the signers
    address[] memory newSigners = new address[](MAX_NUM_ORACLES / 2);
    address[] memory newTransmitters = new address[](MAX_NUM_ORACLES / 2);
    for (uint256 i = 0; i < MAX_NUM_ORACLES / 2; i++) {
      newSigners[i] = signers[i];
      newTransmitters[i] = transmitters[i];
    }

    bytes32 newConfigDigest = s_aggregator.configDigestFromConfigData(
      block.chainid,
      address(s_aggregator),
      2,
      newSigners,
      newTransmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    vm.expectEmit();
    emit ConfigSet(
      blockNumber,
      newConfigDigest,
      2,
      newSigners,
      newTransmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
    s_aggregator.setConfig(newSigners, newTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }
}

contract LatestConfigDetails is ConfiguredDualAggregatorBaseTest {
  function test_ReturnsConfigDetails() public view {
    (uint32 configCount, uint32 blockNumber, bytes32 configDigest) = s_aggregator.latestConfigDetails();

    assertEq(configCount, 1, "config count not incremented");
    assertEq(blockNumber, block.number, "block number is wrong");
    assertEq(configDigest, s_configDigest, "configDigest is not correct");
  }
}

contract GetTransmitters is ConfiguredDualAggregatorBaseTest {
  function test_ReturnsTransmittersList() public view {
    assertEq(s_aggregator.getTransmitters(), s_transmitters, "transmiters list is not the same");
  }
}

contract SetValidatorConfig is DualAggregatorBaseTest {
  event ValidatorConfigSet(
    AggregatorValidatorInterface indexed previousValidator,
    uint32 previousGasLimit,
    AggregatorValidatorInterface indexed currentValidator,
    uint32 currentGasLimit
  );

  function test_EmitsValidatorConfigSet() public {
    AggregatorValidatorInterface oldValidator = AggregatorValidatorInterface(address(0x0));
    AggregatorValidatorInterface newValidator = AggregatorValidatorInterface(address(42));

    vm.expectEmit();
    emit ValidatorConfigSet(oldValidator, 0, newValidator, 1);

    s_aggregator.setValidatorConfig(newValidator, 1);
  }
}

contract GetValidatorConfig is DualAggregatorBaseTest {
  function test_ReturnsValidatorConfig() public {
    AggregatorValidatorInterface newValidator = AggregatorValidatorInterface(address(42));
    uint32 newGasLimit = 1;

    s_aggregator.setValidatorConfig(newValidator, newGasLimit);
    (AggregatorValidatorInterface returnedValidator, uint32 returnedGasLimit) = s_aggregator.getValidatorConfig();
    assertEq(address(returnedValidator), address(newValidator), "did not return the right validator");
    assertEq(returnedGasLimit, newGasLimit, "did not return the right gas limit");
  }
}

contract ValidateAnswer is ConfiguredDualAggregatorBaseTest {
  function test_RevertIf_InsufficientGas() public {
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(1, uint32(block.timestamp));
    AggregatorValidator newValidator = new AggregatorValidator();

    uint256 gas = gasleft();
    s_aggregator.setValidatorConfig(newValidator, uint32(gas));

    _changePrank(s_aggregator.getTransmitters()[0]);
    vm.expectRevert(DualAggregator.InsufficientGas.selector);
    s_aggregator.transmit(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );

    (uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) =
      newValidator.getLatestValidatedValues();

    vm.assertEq(previousRoundId, 0, "previous round id not zero");
    vm.assertEq(previousAnswer, 0, "previous answer not zero");
    vm.assertEq(currentRoundId, 0, "current round id not zero");
    vm.assertEq(currentAnswer, 0, "current answer not zero");
  }

  function test_ValidateWithZeroAddressContract() public {
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(1, uint32(block.timestamp));
    AggregatorValidatorInterface newValidator = AggregatorValidatorInterface(address(0));

    uint256 gas = gasleft();
    s_aggregator.setValidatorConfig(newValidator, uint32(gas));

    _changePrank(s_aggregator.getTransmitters()[0]);
    s_aggregator.transmit(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );
  }

  function test_ValidateWithContractAddress() public {
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(1, uint32(block.timestamp));

    AggregatorValidator newValidator = new AggregatorValidator();
    s_aggregator.setValidatorConfig(newValidator, 21_000_000);

    _changePrank(s_aggregator.getTransmitters()[0]);
    s_aggregator.transmit(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );

    (uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) =
      newValidator.getLatestValidatedValues();

    vm.assertEq(previousRoundId, 0, "previous round id not the same");
    vm.assertEq(previousAnswer, 0, "previous answer not the same");
    vm.assertEq(currentRoundId, 1, "current round id not the same");
    vm.assertEq(currentAnswer, 1, "current answer not the same");
  }
}

contract SetRequesterAccessController is DualAggregatorBaseTest {
  event RequesterAccessControllerSet(AccessControllerInterface old, AccessControllerInterface current);

  function test_EmitsRequesterAccessControllerSet() public {
    AccessControllerInterface oldAccessControllerInterface =
      AccessControllerInterface(REQUESTER_ACCESS_CONTROLLER_ADDRESS);
    AccessControllerInterface newAccessControllerInterface = AccessControllerInterface(address(42));

    vm.expectEmit();
    emit RequesterAccessControllerSet(oldAccessControllerInterface, newAccessControllerInterface);

    s_aggregator.setRequesterAccessController(newAccessControllerInterface);
  }
}

contract GetRequesterAccessController is DualAggregatorBaseTest {
  function test_ReturnsRequesterAccessController() public {
    AccessControllerInterface newAccessControllerInterface = AccessControllerInterface(address(42));
    s_aggregator.setRequesterAccessController(newAccessControllerInterface);

    assertEq(
      address(s_aggregator.getRequesterAccessController()),
      address(newAccessControllerInterface),
      "did not return the right access controller interface"
    );
  }
}

contract RequestNewRound is ConfiguredDualAggregatorBaseTest {
  event RoundRequested(address indexed requester, bytes32 configDigest, uint32 epoch, uint8 round);

  address internal constant USER = address(777);

  function test_RevertIf_NotRequester() public {
    _changePrank(USER);

    vm.mockCall(
      REQUESTER_ACCESS_CONTROLLER_ADDRESS,
      abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
      abi.encode(false)
    );

    vm.expectRevert(DualAggregator.OnlyOwnerAndRequesterCanCall.selector);
    s_aggregator.requestNewRound();
  }

  function test_ShouldRequestNewRound() public {
    uint256 currentRoundId = s_aggregator.getHotVars().latestAggregatorRoundId;
    uint40 currentEpochAndRound = s_aggregator.getHotVars().latestEpochAndRound;

    vm.mockCall(
      REQUESTER_ACCESS_CONTROLLER_ADDRESS,
      abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
      abi.encode(true)
    );

    _changePrank(USER);

    vm.expectEmit();
    emit RoundRequested(USER, s_configDigest, uint32(currentEpochAndRound >> 8), uint8(currentEpochAndRound));
    uint256 newRoundId = s_aggregator.requestNewRound();

    assertEq(newRoundId, currentRoundId + 1, "round id not incremented");
  }
}

contract Transmit is ConfiguredDualAggregatorBaseTest {
  uint32 constant CUTOFF_TIME = 40;

  struct Report {
    int192 price;
    uint32 timestamp;
  }

  function setUp() public override {
    super.setUp();

    s_aggregator.setCutoffTime(CUTOFF_TIME);
  }

  function test_RevertIf_UnauthorizedTransmitter() public {
    bytes32[3] memory reportContext =
      [bytes32(abi.encodePacked("1")), bytes32(abi.encodePacked("2")), bytes32(abi.encodePacked("3"))];
    bytes memory report = abi.encodePacked("1");
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](1);
    bytes32[] memory ss = new bytes32[](1);

    rs[0] = bytes32(abi.encodePacked("1"));
    ss[0] = bytes32(abi.encodePacked("1"));

    vm.expectRevert(DualAggregator.UnauthorizedTransmitter.selector);
    s_aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  function test_RevertIf_ConfigDigestMismatch() public {
    vm.startPrank(s_transmitters[0]);

    bytes32[3] memory reportContext =
      [bytes32(abi.encodePacked("1")), bytes32(abi.encodePacked("2")), bytes32(abi.encodePacked("3"))];
    bytes memory report = abi.encodePacked("1");
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](1);
    bytes32[] memory ss = new bytes32[](1);

    rs[0] = bytes32(abi.encodePacked("1"));
    ss[0] = bytes32(abi.encodePacked("1"));

    vm.expectRevert(DualAggregator.ConfigDigestMismatch.selector);
    s_aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  function test_RevertIf_CalldataLengthMismatch() public {
    vm.startPrank(s_transmitters[0]);

    bytes32[3] memory reportContext = [s_configDigest, bytes32(abi.encodePacked("2")), bytes32(abi.encodePacked("3"))];
    bytes memory report = abi.encodePacked("1");
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](1);
    bytes32[] memory ss = new bytes32[](1);

    rs[0] = bytes32(abi.encodePacked("1"));
    ss[0] = bytes32(abi.encodePacked("1"));

    vm.expectRevert(DualAggregator.CalldataLengthMismatch.selector);
    s_aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  function test_RevertIf_WrongNumberOfSignatures() public {
    vm.startPrank(s_transmitters[0]);

    bytes memory epochAndRound = abi.encodePacked(bytes27(0), uint32(0), uint8(0));

    bytes32[3] memory reportContext = [s_configDigest, bytes32(epochAndRound), bytes32(abi.encodePacked("1"))];
    bytes memory report = new bytes(0);
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](1);
    bytes32[] memory ss = new bytes32[](1);

    vm.expectRevert(DualAggregator.WrongNumberOfSignatures.selector);
    s_aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  function test_RevertIf_SignaturesOutOfRegistration() public {
    vm.startPrank(s_transmitters[0]);

    bytes memory epochAndRound = abi.encodePacked(bytes27(0), uint32(0), uint8(0));
    bytes32[3] memory reportContext = [s_configDigest, bytes32(epochAndRound), bytes32(abi.encodePacked("1"))];
    bytes memory report = new bytes(0);
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](1);

    vm.expectRevert(DualAggregator.SignaturesOutOfRegistration.selector);
    s_aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  function test_RevertIf_SignatureError() public {
    vm.startPrank(s_transmitters[0]);

    ReportGenerator.SignedReport memory signedReport =
      s_reportGenerator.generateSignedReport(0, uint32(block.timestamp));

    signedReport.rs[0] = bytes32(abi.encodePacked("1"));
    signedReport.rs[1] = bytes32(abi.encodePacked("1"));
    signedReport.ss[0] = bytes32(abi.encodePacked("1"));
    signedReport.ss[1] = bytes32(abi.encodePacked("1"));

    vm.expectRevert(DualAggregator.SignatureError.selector);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );
  }

  function test_RevertIf_DuplicateSigner() public {
    vm.startPrank(s_transmitters[0]);

    ReportGenerator.SignedReport memory signedReport =
      s_reportGenerator.generateSignedReport(0, uint32(block.timestamp));

    signedReport.rs[1] = signedReport.rs[0];
    signedReport.ss[1] = signedReport.ss[0];

    vm.expectRevert(DualAggregator.DuplicateSigner.selector);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );
  }

  function test_RevertIf_ReportLengthMismatch() public {
    ReportGenerator.SignedReport memory signedReport =
      s_reportGenerator.generateSignedReportWithLengthMismatch(0, uint32(block.timestamp));

    _changePrank(s_transmitters[0]);

    vm.expectRevert(DualAggregator.ReportLengthMismatch.selector);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );
  }

  function test_RevertIf_NumObservationsOutOfBounds() public {
    ReportGenerator.SignedReport memory signedReport =
      s_reportGenerator.generateSignedReportWithTooManyOracles(0, uint32(block.timestamp));

    _changePrank(s_transmitters[0]);

    vm.expectRevert(DualAggregator.NumObservationsOutOfBounds.selector);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );
  }

  function test_RevertIf_TooFewValuesToTrustMedian() public {
    ReportGenerator.SignedReport memory signedReport =
      s_reportGenerator.generateSignedReportWithFOracles(0, uint32(block.timestamp), s_f);

    _changePrank(s_transmitters[0]);

    vm.expectRevert(DualAggregator.TooFewValuesToTrustMedian.selector);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );
  }

  function test_RevertIf_MedianIsOutOfMinMaxRange() public {
    ReportGenerator.SignedReport memory signedReport =
      s_reportGenerator.generateSignedReport(MAX_ANSWER + 1, uint32(block.timestamp));

    _changePrank(s_transmitters[0]);

    vm.expectRevert(DualAggregator.MedianIsOutOfMinMaxRange.selector);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );

    signedReport = s_reportGenerator.generateSignedReport(MIN_ANSWER - 1, uint32(block.timestamp));

    vm.expectRevert(DualAggregator.MedianIsOutOfMinMaxRange.selector);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );
  }

  function test_RevertIf_SendTheSameReportTwiceToThePrimaryFeed() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);

    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: 0
    });
    _checkRounds(1, 0);

    _mineBlock();

    vm.expectRevert(DualAggregator.StaleReport.selector);
    s_aggregator.transmit(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );
  }

  function test_RevertIf_SendTheSameReportTwiceToTheSecondaryFeed() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);

    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: 0,
      expectedSecondaryFeedAnswer: report1.price
    });
    _checkRounds(0, 1);

    _mineBlock();

    vm.expectRevert(DualAggregator.StaleReport.selector);
    s_aggregator.transmitSecondary(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );
  }

  function test_RevertIf_SendTheSameReportAlternateFirstPrimaryFeed() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);

    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: 0
    });
    _checkRounds(1, 0);

    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: report1.price
    });
    _checkRounds(1, 1);

    vm.expectRevert(DualAggregator.StaleReport.selector);
    s_aggregator.transmit(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );

    vm.expectRevert(DualAggregator.StaleReport.selector);
    s_aggregator.transmitSecondary(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );
  }

  function test_RevertIf_SendTheSameReportAlternateFirstSecondaryFeed() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);

    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: 0,
      expectedSecondaryFeedAnswer: report1.price
    });
    _checkRounds(0, 1);

    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: report1.price
    });
    _checkRounds(1, 1);

    vm.expectRevert(DualAggregator.StaleReport.selector);
    s_aggregator.transmitSecondary(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );

    vm.expectRevert(DualAggregator.StaleReport.selector);
    s_aggregator.transmit(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );
  }

  function test_ReadExpectedInitialState() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);

    _transmitAndCheck(signedReport1, 0, 0, false, false); // no transmission but check
  }

  function test_SyncFeedsTransmitStandardFirstNeverSameBlock() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);
    // Report 1
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: 0
    });

    _checkRounds(1, 0);
    _mineBlock();

    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: report1.price
    });

    _checkRounds(1, 1);
    _mineBlock();

    Report memory report2 = Report({price: 2, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport2 =
      s_reportGenerator.generateSignedReport(report2.price, report2.timestamp);
    // Report 2
    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report2.price,
      expectedSecondaryFeedAnswer: report1.price
    });

    _checkRounds(2, 1);
    _mineBlock();

    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: report2.price,
      expectedSecondaryFeedAnswer: report2.price
    });

    _checkRounds(2, 2);
  }

  function test_SyncFeedsTransmitSecondaryFirstNeverSameBlock() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);
    // Report 1
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: 0, // locked
      expectedSecondaryFeedAnswer: report1.price
    });
    _mineBlock();

    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: report1.price
    });
    _checkRounds(1, 1);
    _mineBlock();

    Report memory report2 = Report({price: 2, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport2 =
      s_reportGenerator.generateSignedReport(report2.price, report2.timestamp);
    // Report 2
    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: report1.price, // locked
      expectedSecondaryFeedAnswer: report2.price
    });
    _mineBlock();

    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report2.price,
      expectedSecondaryFeedAnswer: report2.price
    });

    _checkRounds(2, 2);
  }

  function test_SyncFeedsTransmitStandardFirstAlwaysSameBlock() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);
    // Report 1
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: 0
    });
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: report1.price
    });

    _mineBlock();
    _checkRounds(1, 1);

    Report memory report2 = Report({price: 2, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport2 =
      s_reportGenerator.generateSignedReport(report2.price, report2.timestamp);
    // Report 2
    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report2.price,
      expectedSecondaryFeedAnswer: report1.price
    });
    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: report2.price,
      expectedSecondaryFeedAnswer: report2.price
    });
    _checkRounds(2, 2);
  }

  function test_SyncFeedsTransmitSecondaryFirstAlwaysSameBlock() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);
    // Report 1
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: 0, // locked
      expectedSecondaryFeedAnswer: report1.price
    });
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: report1.price
    });

    _checkRounds(1, 1);
    _mineBlock();

    Report memory report2 = Report({price: 2, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport2 =
      s_reportGenerator.generateSignedReport(report2.price, report2.timestamp);
    // Report 2
    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: report1.price, // locked
      expectedSecondaryFeedAnswer: report2.price
    });
    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report2.price,
      expectedSecondaryFeedAnswer: report2.price
    });

    _checkRounds(2, 2);
  }

  function test_OutOfSyncFeedsSecondaryFeedFallbackToStandardFeed() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);
    // Report 1
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: 0
    });
    _mineBlock();
    _checkRounds(1, 0);

    Report memory report2 = Report({price: 2, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport2 =
      s_reportGenerator.generateSignedReport(report2.price, report2.timestamp);
    // Report 2
    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report2.price,
      expectedSecondaryFeedAnswer: 0
    });

    _checkRounds(2, 0);
    // skip cutoff time, actual: 40
    skip(CUTOFF_TIME + 1);
    _checkRounds(2, 2);

    Report memory report3 = Report({price: 3, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport3 =
      s_reportGenerator.generateSignedReport(report3.price, report3.timestamp);
    // Report 3
    _transmitAndCheck({
      signedReport: signedReport3,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report3.price,
      expectedSecondaryFeedAnswer: report2.price // freshest report before cutoff
    });
    // transmit old report to secondary feed
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: report3.price,
      expectedSecondaryFeedAnswer: report2.price // old report from before the cutoff time, ignore it
    });

    _checkRounds(3, 2);
    _mineBlock();

    Report memory report4 = Report({price: 4, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport4 =
      s_reportGenerator.generateSignedReport(report4.price, report4.timestamp);
    // Report 4
    _transmitAndCheck({
      signedReport: signedReport4,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report4.price,
      expectedSecondaryFeedAnswer: report2.price // secondary feed is still stale
    });
    _transmitAndCheck({
      signedReport: signedReport3,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: report4.price,
      expectedSecondaryFeedAnswer: report3.price // old report but newer than latest secondary report
    });

    _checkRounds(4, 3);
  }

  function test_OutOfSyncFeedsPrimaryIsSourcedFromSecondaryWithLockDelay() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);
    // Report 1
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: 0, // locked
      expectedSecondaryFeedAnswer: report1.price
    });
    _checkRounds(0, 1);
    _mineBlock();
    _checkRounds(1, 1);

    Report memory report2 = Report({price: 2, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport2 =
      s_reportGenerator.generateSignedReport(report2.price, report2.timestamp);
    // Report 2
    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: report1.price, // locked,
      expectedSecondaryFeedAnswer: report2.price
    });

    _checkRounds(1, 2);
  }

  function test_BothFeedsStalledIncomingSecondaryReportIsFromBeforeCutoffTime() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);
    // Report 1
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: 0
    });
    _mineBlock();
    _checkRounds(1, 0);

    // Report 2 generate from before the cutoff time
    Report memory report2 = Report({price: 2, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport2 =
      s_reportGenerator.generateSignedReport(report2.price, report2.timestamp);

    // skip cutoff time
    skip(CUTOFF_TIME + 1);
    _checkRounds(1, 1);

    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: false,
      transmitSecondary: true,
      expectedStandardFeedAnswer: report1.price, // locked
      expectedSecondaryFeedAnswer: report2.price
    });
    _checkRounds(1, 2);
    _mineBlock();

    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: false,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report2.price, // unlocked
      expectedSecondaryFeedAnswer: report2.price
    });
    _checkRounds(2, 2);
  }

  function test_BothFeedsStalledIncomingPrimaryReportIsFromBeforeCutoffTime() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);
    // Report 1
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: 0
    });
    _mineBlock();
    _checkRounds(1, 0);

    // Report 2 generate from before the cutoff time
    Report memory report2 = Report({price: 2, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport2 =
      s_reportGenerator.generateSignedReport(report2.price, report2.timestamp);
    _mineBlock();

    // skip cutoff time
    skip(CUTOFF_TIME + 1);
    _checkRounds(1, 1);

    // report 2 transmitted
    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report2.price,
      expectedSecondaryFeedAnswer: report1.price // because newest comes from before the cutoff time
    });

    _checkRounds(2, 1);
  }

  function test_BothFeedsStalledIncomingReportIsFromAfterCutoffTime() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);
    // Report 1
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: 0
    });
    _mineBlock();
    _checkRounds(1, 0);

    // skip cutoff time
    skip(CUTOFF_TIME + 1);
    _checkRounds(1, 1);

    // Report 2 generate from after the cutoff time
    Report memory report2 = Report({price: 2, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport2 =
      s_reportGenerator.generateSignedReport(report2.price, report2.timestamp);

    _transmitAndCheck({
      signedReport: signedReport2,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report2.price,
      expectedSecondaryFeedAnswer: report1.price // because newest comes from after the cutoff time
    });

    _checkRounds(2, 1);
  }

  function test_IncomingSecondaryReportHasNotBeenRecordedOlderThanLatestReportOlderThanCutoffTime() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);
    // Report 1
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: 0
    });
    _mineBlock();
    _checkRounds(1, 0);

    // generate Signed Report 2, will not reach standard feed
    Report memory report2 = Report({price: 2, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport2 =
      s_reportGenerator.generateSignedReport(report2.price, report2.timestamp);
    _mineBlock();

    // skip cutoff time
    skip(CUTOFF_TIME + 1);
    _checkRounds(1, 1);

    Report memory report3 = Report({price: 3, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport3 =
      s_reportGenerator.generateSignedReport(report3.price, report3.timestamp);
    // Report 3
    _transmitAndCheck({
      signedReport: signedReport3,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report3.price,
      expectedSecondaryFeedAnswer: report1.price
    });
    _mineBlock();
    _checkRounds(2, 1);

    // report 2 reaches the secondary feed, but it is dropped due to being an orphan
    vm.expectRevert(DualAggregator.StaleReport.selector);
    s_aggregator.transmitSecondary(
      signedReport2.reportContext, signedReport2.report, signedReport2.rs, signedReport2.ss, signedReport2.rawVs
    );

    // check the standard feed
    _changePrank(s_aggregator.getTransmitters()[0]);
    (, int256 standardAnswer,,,) = s_aggregator.latestRoundData();
    assertEq(report3.price, standardAnswer, "standard feed answer is not correct");

    // check the secondary feed
    _changePrank(SECONDARY_PROXY);
    (, int256 secondaryAnswer,,,) = s_aggregator.latestRoundData();
    assertEq(report1.price, secondaryAnswer, "secondary feed answer is not correct");

    _checkRounds(2, 1);
  }

  function test_IncomingSecondaryReportHasNotBeenRecordedOlderThanLatestReportNewerThanCutoffTime() public {
    Report memory report1 = Report({price: 1, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(report1.price, report1.timestamp);
    // Report 1
    _transmitAndCheck({
      signedReport: signedReport1,
      transmitPrimary: true,
      transmitSecondary: true,
      expectedStandardFeedAnswer: report1.price,
      expectedSecondaryFeedAnswer: report1.price
    });
    _checkRounds(1, 1);
    _mineBlock();

    // generate Signed Report 2, will not reach standard feed
    Report memory report2 = Report({price: 2, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport2 =
      s_reportGenerator.generateSignedReport(report2.price, report2.timestamp);
    _mineBlock();

    // report 3 will reach the standard feed
    Report memory report3 = Report({price: 3, timestamp: uint32(block.timestamp)});
    ReportGenerator.SignedReport memory signedReport3 =
      s_reportGenerator.generateSignedReport(report3.price, report3.timestamp);
    _transmitAndCheck({
      signedReport: signedReport3,
      transmitPrimary: true,
      transmitSecondary: false,
      expectedStandardFeedAnswer: report3.price,
      expectedSecondaryFeedAnswer: report1.price
    });
    _mineBlock();
    _checkRounds(2, 1);

    // report 2 reaches the secondary feed, but it is dropped due to being an orphan
    vm.expectRevert(DualAggregator.StaleReport.selector);
    s_aggregator.transmitSecondary(
      signedReport2.reportContext, signedReport2.report, signedReport2.rs, signedReport2.ss, signedReport2.rawVs
    );

    // check the standard feed
    _changePrank(s_aggregator.getTransmitters()[0]);
    (, int256 standardAnswer,,,) = s_aggregator.latestRoundData();
    assertEq(report3.price, standardAnswer, "standard feed answer is not correct");

    // check the secondary feed
    _changePrank(SECONDARY_PROXY);
    (, int256 secondaryAnswer,,,) = s_aggregator.latestRoundData();
    assertEq(report1.price, secondaryAnswer, "secondary feed answer is not correct");

    _checkRounds(2, 1);
  }

  function test_IncomingSecondaryReportRevertsDueToMaxIterations() public {
    // define 4 as the new max sync iterations
    s_aggregator.setMaxSyncIterations(4);
    s_aggregator.setCutoffTime(60);

    // sign the report 1
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(1, uint32(block.timestamp));
    _mineBlock();

    // transmit the signed report from the primary feed
    _changePrank(s_aggregator.getTransmitters()[0]);
    s_aggregator.transmit(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );
    _checkRounds(1, 0);

    // transmit 4 more reports
    Report memory report;
    ReportGenerator.SignedReport memory signedReport;
    for (int192 i = 2; i <= 5; ++i) {
      report = Report({price: i, timestamp: uint32(block.timestamp)});
      signedReport = s_reportGenerator.generateSignedReport(report.price, report.timestamp);
      // transmit report from the primary feed
      _transmitAndCheck({
        signedReport: signedReport,
        transmitPrimary: true,
        transmitSecondary: false,
        expectedStandardFeedAnswer: report.price,
        expectedSecondaryFeedAnswer: 0
      });
      _mineBlock();
    }

    _checkRounds(5, 0);

    // report 1 reaches the secondary feed, but it reverts due to max sync iterations
    vm.expectRevert(DualAggregator.MaxSyncIterations.selector);
    s_aggregator.transmitSecondary(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );
  }

  function _transmitAndCheck(
    ReportGenerator.SignedReport memory signedReport,
    int256 expectedStandardFeedAnswer,
    int256 expectedSecondaryFeedAnswer,
    bool transmitPrimary,
    bool transmitSecondary
  ) internal {
    _changePrank(s_aggregator.getTransmitters()[0]);

    if (transmitSecondary) {
      s_aggregator.transmitSecondary(
        signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
      );
    }

    if (transmitPrimary) {
      s_aggregator.transmit(
        signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
      );
    }

    // check the standard feed
    (, int256 standardAnswer,,,) = s_aggregator.latestRoundData();
    assertEq(int256(expectedStandardFeedAnswer), standardAnswer, "standard feed answer is not correct");

    // check the secondary feed
    _changePrank(SECONDARY_PROXY);
    (, int256 secondaryAnswer,,,) = s_aggregator.latestRoundData();
    assertEq(int256(expectedSecondaryFeedAnswer), secondaryAnswer, "secondary feed answer is not correct");
  }

  function _checkRounds(uint256 expectedPrimaryRound, uint256 expectedSecondaryRound) internal {
    _changePrank(SECONDARY_PROXY);
    assertEq(expectedSecondaryRound, s_aggregator.latestRound(), "secondary feed round is not correct");

    _changePrank(s_aggregator.getTransmitters()[0]);
    assertEq(expectedPrimaryRound, s_aggregator.latestRound(), "standard feed round is not correct");
  }

  function _mineBlock() internal {
    skip(12);
  }
}

contract TransmittedDualAggregatorBaseTest is ConfiguredDualAggregatorBaseTest {
  function setUp() public override {
    super.setUp();

    ReportGenerator.SignedReport memory signedReport = s_reportGenerator.generateSignedReport(1, 1);

    vm.warp(1);
    _changePrank(s_aggregator.getTransmitters()[0]);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );
  }
}

contract LatestTransmissionDetails is TransmittedDualAggregatorBaseTest {
  function test_RevertIf_NotEOA() public {
    vm.expectRevert(DualAggregator.OnlyCallableByEOA.selector);
    s_aggregator.latestTransmissionDetails();
  }

  function test_ReturnsLatestTransmissionDetails() public {
    vm.stopPrank();
    vm.startPrank(address(100), address(100));
    (bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer, uint64 latestTimestamp) =
      s_aggregator.latestTransmissionDetails();

    assertEq(configDigest, s_configDigest);
    assertEq(epoch, 0);
    assertEq(round, 1);
    assertEq(latestAnswer, 1);
    assertEq(latestTimestamp, 1);
  }
}

contract LatestConfigDigestAndEpoch is TransmittedDualAggregatorBaseTest {
  function test_ReturnsLatestConfigDigestAndEpoch() public view {
    (bool scanLogs, bytes32 configDigest, uint32 epoch) = s_aggregator.latestConfigDigestAndEpoch();

    assertEq(scanLogs, false, "scanLogs was not correct");
    assertEq(
      configDigest,
      s_aggregator.configDigestFromConfigData(
        block.chainid,
        address(s_aggregator),
        1,
        s_signers,
        s_transmitters,
        s_f,
        s_onchainConfig,
        s_offchainConfigVersion,
        s_offchainConfig
      ),
      "configDigest incorrect"
    );
    assertEq(epoch, 0, "epoch not correct");
  }
}

// latestAnswer(): test all paths
contract LatestAnswer is RoundDataDualAggregatorBaseTest {
  // test return the latest primary answer
  function test_ReturnsLatestPrimaryAnswer() public {
    setDualAggregatorBase(30, 20, 6, 6, false);
    assertEq(s_aggregator.latestAnswer(), 15);
  }

  // test return the latest primary answer after block
  function test_ReturnsLatestPrimaryAnswerAfterBlock() public {
    setDualAggregatorBase(30, 20, 6, 6, true);
    vm.warp(31);
    assertEq(s_aggregator.latestAnswer(), 15);
  }

  // test latest primary answer in the same block, return previous one
  function test_ReturnsPreviousPrimaryAnswer() public {
    setDualAggregatorBase(30, 20, 6, 6, true);
    assertEq(s_aggregator.latestAnswer(), 14);
  }

  // test return the latest secondary answer
  function test_ReturnsLatestSecondaryAnswer() public {
    setDualAggregatorBase(30, 20, 6, 2, false);
    vm.startPrank(address(102));
    assertEq(s_aggregator.latestAnswer(), 11);
  }

  // test return the secondary answer, synced with primary rounds
  function test_ReturnsSyncedAnswer() public {
    setDualAggregatorBase(30, 20, 6, 2, false);
    s_aggregator.setCutoffTime(9);

    vm.startPrank(address(102));
    assertEq(s_aggregator.latestAnswer(), 13);
  }
}

// latestTimestamp(): test all paths
contract LatestTimestamp is RoundDataDualAggregatorBaseTest {
  // test return the latest primary timestamp
  function test_ReturnsLatestPrimaryTimestamp() public {
    setDualAggregatorBase(30, 20, 6, 6, false);
    assertEq(s_aggregator.latestTimestamp(), 30);
  }

  // test return the latest primary timestamp after block
  function test_ReturnsLatestPrimaryTimestampAfterBlock() public {
    setDualAggregatorBase(30, 20, 6, 6, true);
    vm.warp(31);
    assertEq(s_aggregator.latestTimestamp(), 30);
  }

  // test latest primary timestamp in the same block, return previous one
  function test_ReturnsPreviousPrimaryTimestamp() public {
    setDualAggregatorBase(30, 20, 6, 6, true);
    assertEq(s_aggregator.latestTimestamp(), 25);
  }

  // test return the latest secondary timestamp
  function test_ReturnsLatestSecondaryTimestamp() public {
    setDualAggregatorBase(30, 20, 6, 2, false);
    vm.startPrank(address(102));
    assertEq(s_aggregator.latestTimestamp(), 10);
  }

  // test return the secondary timestamp, synced with primary rounds
  function test_ReturnsSyncedTimestamp() public {
    setDualAggregatorBase(30, 20, 6, 2, false);
    s_aggregator.setCutoffTime(9);

    vm.startPrank(address(102));
    assertEq(s_aggregator.latestTimestamp(), 20);
  }
}

// latestRound(): test all the paths
contract LatestRound is RoundDataDualAggregatorBaseTest {
  // test return the latest primary round id
  function test_ReturnsLatestPrimaryRoundId() public {
    setDualAggregatorBase(30, 20, 6, 6, false);
    assertEq(s_aggregator.latestRound(), 6);
  }

  // test return the latest primary round id after block
  function test_ReturnsLatestPrimaryRoundIdAfterBlock() public {
    setDualAggregatorBase(30, 20, 6, 6, true);
    vm.warp(31);
    assertEq(s_aggregator.latestRound(), 6);
  }

  // test latest primary round id in the same block, return previous one
  function test_ReturnsPreviousPrimaryRoundId() public {
    setDualAggregatorBase(30, 20, 6, 6, true);
    assertEq(s_aggregator.latestRound(), 5);
  }

  // test return the latest secondary round id
  function test_ReturnsLatestSecondaryRoundId() public {
    setDualAggregatorBase(30, 20, 6, 2, false);
    vm.startPrank(address(102));
    assertEq(s_aggregator.latestRound(), 2);
  }

  // test return the secondary round id, synced with primary rounds
  function test_ReturnsSyncedRoundId() public {
    setDualAggregatorBase(30, 20, 6, 2, false);
    s_aggregator.setCutoffTime(9);

    vm.startPrank(address(102));
    assertEq(s_aggregator.latestRound(), 4);
  }
}

// getAnswer(): test all paths
contract GetAnswer is RoundDataDualAggregatorBaseTest {
  // test return primary answer
  function test_ReturnsPrimaryGetAnswer() public {
    setDualAggregatorBase(30, 0, 6, 6, false);
    assertEq(s_aggregator.getAnswer(6), 15);
  }

  // test return primary answer after block
  function test_ReturnsPrimaryGetAnswerAfterBlock() public {
    setDualAggregatorBase(30, 0, 6, 6, true);
    vm.warp(31);
    assertEq(s_aggregator.getAnswer(6), 15);
  }

  // test return primary answer, not allowed
  function test_ReturnsPrimaryGetAnswerNotAllowed() public {
    setDualAggregatorBase(30, 0, 6, 6, true);
    assertEq(s_aggregator.getAnswer(6), 0);
  }

  // test return secondary answer
  function test_ReturnsSecondaryGetAnswer() public {
    setDualAggregatorBase(30, 0, 6, 2, false);
    vm.startPrank(address(102));
    assertEq(s_aggregator.getAnswer(2), 11);
  }

  // test return secondary answer, not allowed
  function test_ReturnsSecondaryGetAnswerNotAllowed() public {
    setDualAggregatorBase(30, 0, 6, 2, false);
    s_aggregator.setCutoffTime(20);

    vm.startPrank(address(102));
    assertEq(s_aggregator.getAnswer(3), 0);
  }

  // test return secondary answer, synced with primary feed
  function test_ReturnsSecondaryGetAnswerSynced() public {
    setDualAggregatorBase(30, 0, 6, 2, false);
    s_aggregator.setCutoffTime(14);

    vm.startPrank(address(102));
    assertEq(s_aggregator.getAnswer(3), 12);
  }
}

// getTimestamp(): test all paths
contract GetTimestamp is RoundDataDualAggregatorBaseTest {
  // test return primary timestamp
  function test_ReturnsPrimaryGetTimestamp() public {
    setDualAggregatorBase(30, 0, 6, 6, false);
    assertEq(s_aggregator.getTimestamp(6), 30);
  }

  // test return primary timestamp after block
  function test_ReturnsPrimaryGetTimestampAfterBlock() public {
    setDualAggregatorBase(30, 0, 6, 6, true);
    vm.warp(31);
    assertEq(s_aggregator.getTimestamp(6), 30);
  }

  // test return primary timestamp, not allowed
  function test_ReturnsPrimaryGetTimestampNotAllowed() public {
    setDualAggregatorBase(30, 0, 6, 6, true);
    assertEq(s_aggregator.getTimestamp(6), 0);
  }

  // test return secondary timestamp
  function test_ReturnsSecondaryGetTimestamp() public {
    setDualAggregatorBase(30, 0, 6, 2, false);
    vm.startPrank(address(102));
    assertEq(s_aggregator.getTimestamp(2), 10);
  }

  // test return secondary timestamp, not allowed
  function test_ReturnsSecondaryGetTimestampNotAllowed() public {
    setDualAggregatorBase(30, 0, 6, 2, false);
    s_aggregator.setCutoffTime(20);

    vm.startPrank(address(102));
    assertEq(s_aggregator.getTimestamp(3), 0);
  }

  // test return secondary timestamp, synced with primary feed
  function test_ReturnsSecondaryGetTimestampSynced() public {
    setDualAggregatorBase(30, 0, 6, 2, false);
    s_aggregator.setCutoffTime(14);

    vm.startPrank(address(102));
    assertEq(s_aggregator.getTimestamp(3), 15);
  }
}

contract Description is TransmittedDualAggregatorBaseTest {
  function test_ReturnsCorrectDescription() public view {
    assertEq(s_aggregator.description(), "TEST");
  }
}

// getRoundData(): test all paths
contract GetRoundData is RoundDataDualAggregatorBaseTest {
  // test return primary round data
  function test_ReturnsPrimaryGetRoundData() public {
    setDualAggregatorBase(30, 0, 6, 6, false);
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      s_aggregator.getRoundData(6);

    assertEq(roundId, 6);
    assertEq(answer, 15);
    assertEq(startedAt, 26);
    assertEq(updatedAt, 30);
    assertEq(answeredInRound, 6);
  }

  // test return primary round data after block
  function test_ReturnsPrimaryGetRoundDataAfterBlock() public {
    setDualAggregatorBase(30, 0, 6, 6, true);
    vm.warp(31);
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      s_aggregator.getRoundData(6);

    assertEq(roundId, 6);
    assertEq(answer, 15);
    assertEq(startedAt, 26);
    assertEq(updatedAt, 30);
    assertEq(answeredInRound, 6);
  }

  // test return primary round data, not allowed
  function test_ReturnsPrimaryGetRoundDataNotAllowed() public {
    setDualAggregatorBase(30, 0, 6, 6, true);
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      s_aggregator.getRoundData(6);

    assertEq(roundId, 0);
    assertEq(answer, 0);
    assertEq(startedAt, 0);
    assertEq(updatedAt, 0);
    assertEq(answeredInRound, 0);
  }

  // test return secondary round data
  function test_ReturnsSecondaryGetRoundData() public {
    setDualAggregatorBase(30, 0, 6, 2, false);
    vm.startPrank(address(102));
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      s_aggregator.getRoundData(2);

    assertEq(roundId, 2);
    assertEq(answer, 11);
    assertEq(startedAt, 6);
    assertEq(updatedAt, 10);
    assertEq(answeredInRound, 2);
  }

  // test return secondary round data, not allowed
  function test_ReturnsSecondaryGetRoundDataNotAllowed() public {
    setDualAggregatorBase(30, 0, 6, 2, false);
    s_aggregator.setCutoffTime(20);

    vm.startPrank(address(102));
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      s_aggregator.getRoundData(3);

    assertEq(roundId, 0);
    assertEq(answer, 0);
    assertEq(startedAt, 0);
    assertEq(updatedAt, 0);
    assertEq(answeredInRound, 0);
  }

  // test return secondary round data, synced with primary feed
  function test_ReturnsSecondaryGetRoundDataSynced() public {
    setDualAggregatorBase(30, 0, 6, 2, false);
    s_aggregator.setCutoffTime(14);

    vm.startPrank(address(102));
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      s_aggregator.getRoundData(3);

    assertEq(roundId, 3);
    assertEq(answer, 12);
    assertEq(startedAt, 11);
    assertEq(updatedAt, 15);
    assertEq(answeredInRound, 3);
  }
}

// latestRoundData(): test all paths
contract LatestRoundData is RoundDataDualAggregatorBaseTest {
  // test return the latest primary round data
  function test_ReturnsLatestPrimaryRoundData() public {
    setDualAggregatorBase(30, 20, 6, 6, false);
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      s_aggregator.latestRoundData();

    assertEq(roundId, 6);
    assertEq(answer, 15);
    assertEq(startedAt, 26);
    assertEq(updatedAt, 30);
    assertEq(answeredInRound, 6);
  }

  // test return the latest primary round data after block
  function test_ReturnsLatestPrimaryRoundDataAfterBlock() public {
    setDualAggregatorBase(30, 20, 6, 6, true);
    vm.warp(31);
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      s_aggregator.latestRoundData();

    assertEq(roundId, 6);
    assertEq(answer, 15);
    assertEq(startedAt, 26);
    assertEq(updatedAt, 30);
    assertEq(answeredInRound, 6);
  }

  // test latest primary round data in the same block, return previous one
  function test_ReturnsPreviousPrimaryRoundData() public {
    setDualAggregatorBase(30, 20, 6, 6, true);
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      s_aggregator.latestRoundData();

    assertEq(roundId, 5);
    assertEq(answer, 14);
    assertEq(startedAt, 21);
    assertEq(updatedAt, 25);
    assertEq(answeredInRound, 5);
  }

  // test return the latest secondary round data
  function test_ReturnsLatestSecondaryRoundData() public {
    setDualAggregatorBase(30, 20, 6, 2, false);
    vm.startPrank(address(102));
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      s_aggregator.latestRoundData();

    assertEq(roundId, 2);
    assertEq(answer, 11);
    assertEq(startedAt, 6);
    assertEq(updatedAt, 10);
    assertEq(answeredInRound, 2);
  }

  // test return the secondary round data, synced with primary rounds
  function test_ReturnsSyncedRoundData() public {
    setDualAggregatorBase(30, 20, 6, 2, false);
    s_aggregator.setCutoffTime(9);

    vm.startPrank(address(102));
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      s_aggregator.latestRoundData();

    assertEq(roundId, 4);
    assertEq(answer, 13);
    assertEq(startedAt, 16);
    assertEq(updatedAt, 20);
    assertEq(answeredInRound, 4);
  }
}

contract SetLinkToken is ConfiguredDualAggregatorBaseTest {
  event LinkTokenSet(LinkTokenInterface indexed oldLinkToken, LinkTokenInterface indexed newLinkToken);

  LinkToken internal s_newLinkToken;
  address internal s_transmitter;

  function setUp() public override {
    super.setUp();
    s_newLinkToken = new LinkToken();

    s_transmitter = s_aggregator.getTransmitters()[0];

    s_aggregator.setPayees(s_transmitters, s_transmitters);

    s_aggregator.setBilling(10, 1, 100, 100, 0);

    _changePrank(s_transmitter);

    ReportGenerator.SignedReport memory signedReport = s_reportGenerator.generateSignedReport(1, 1);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );
  }

  function test_RevertIf_TransferFundsFailed() public {
    uint256 totalDue = s_aggregator.totalLinkDue();
    deal(address(s_link), address(s_aggregator), totalDue);

    _changePrank(s_aggregator.owner());
    // we want to owe nothing so we can fail at transferring the remaining funds
    s_aggregator.setBilling(0, 0, 0, 0, 0);
    vm.mockCall(address(s_link), abi.encodeWithSelector(LinkTokenInterface.transfer.selector), abi.encode(false));
    vm.expectRevert(DualAggregator.TransferRemainingFundsFailed.selector);
    s_aggregator.setLinkToken(LinkTokenInterface(address(s_newLinkToken)), address(43));
  }

  function test_PayOutOraclesBeforeSettingNewLinkToken() public {
    // check the balance of the transmitter before the withdrawal
    uint256 balanceBefore = s_link.balanceOf(s_transmitter);
    uint256 amountDue = s_aggregator.owedPayment(s_transmitter);

    // deal the amount due to the aggregator
    uint256 totalDue = s_aggregator.totalLinkDue();
    deal(address(s_link), address(s_aggregator), totalDue);

    _changePrank(s_aggregator.owner());

    vm.expectEmit();
    emit LinkTokenSet(LinkTokenInterface(address(s_link)), LinkTokenInterface(address(s_newLinkToken)));
    s_aggregator.setLinkToken(LinkTokenInterface(address(s_newLinkToken)), s_aggregator.owner());

    // check the balance of the transmitter after the withdrawal
    uint256 balanceAfter = s_link.balanceOf(s_transmitter);

    assertEq(balanceAfter, balanceBefore + amountDue, "did not pay the transmitter");
  }
}

contract GetLinkToken is DualAggregatorBaseTest {
  function test_ReturnsLinkToken() public view {
    assertEq(address(s_aggregator.getLinkToken()), address(s_link), "did not return the right link token interface");
  }
}

contract SetBillingAccessController is DualAggregatorBaseTest {
  event BillingAccessControllerSet(AccessControllerInterface old, AccessControllerInterface current);

  AccessControllerInterface internal oldBillingAccessController =
    AccessControllerInterface(BILLING_ACCESS_CONTROLLER_ADDRESS);
  AccessControllerInterface internal newBillingAccessController = AccessControllerInterface(address(42));

  function test_EmitsBillingAccessControllerSet() public {
    vm.expectEmit();
    emit BillingAccessControllerSet(oldBillingAccessController, newBillingAccessController);

    s_aggregator.setBillingAccessController(newBillingAccessController);
  }
}

contract GetBillingAccessController is DualAggregatorBaseTest {
  function test_ReturnsBillingAccessController() public view {
    assertEq(
      address(s_aggregator.getBillingAccessController()),
      BILLING_ACCESS_CONTROLLER_ADDRESS,
      "did not return the right billing access controller"
    );
  }
}

contract SetBilling is DualAggregatorBaseTest {
  event BillingSet(
    uint32 maximumGasPriceGwei,
    uint32 reasonableGasPriceGwei,
    uint32 observationPaymentGjuels,
    uint32 transmissionPaymentGjuels,
    uint24 accountingGas
  );

  address internal constant USER = address(42);

  function test_RevertIf_NotOwner() public {
    vm.mockCall(
      BILLING_ACCESS_CONTROLLER_ADDRESS,
      abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
      abi.encode(false)
    );
    vm.startPrank(USER);
    vm.expectRevert(DualAggregator.OnlyOwnerAndBillingAdminCanCall.selector);

    s_aggregator.setBilling(0, 0, 0, 0, 0);
  }

  function test_EmitsBillingSet() public {
    vm.expectEmit();
    emit BillingSet(0, 0, 0, 0, 0);

    s_aggregator.setBilling(0, 0, 0, 0, 0);
  }
}

contract GetBilling is DualAggregatorBaseTest {
  function test_ReturnsBillingData() public view {
    (
      uint32 returnedMaxGasPriceGwei,
      uint32 returnedReasonableGasPriceGwei,
      uint32 returnedObservationPaymentGjuels,
      uint32 returnedTransmissionPaymentGjuels,
      uint32 returnedAccountingGas
    ) = s_aggregator.getBilling();

    assertEq(returnedMaxGasPriceGwei, 0, "maxGasPriceGwei incorrect");
    assertEq(returnedReasonableGasPriceGwei, 0, "reasonableGasPriceGwei incorrect");
    assertEq(returnedObservationPaymentGjuels, 0, "observationPaymentGjuels incorrect");
    assertEq(returnedTransmissionPaymentGjuels, 0, "transmissionPaymentGjuels incorrect");
    assertEq(returnedAccountingGas, 0, "accountingGas incorrect");
  }
}

contract WithdrawPayment is ConfiguredDualAggregatorBaseTest {
  address internal s_transmitter;

  function setUp() public override {
    super.setUp();

    s_transmitter = s_aggregator.getTransmitters()[0];

    s_aggregator.setPayees(s_transmitters, s_transmitters);

    s_aggregator.setBilling(10, 1, 100, 100, 0);

    _changePrank(s_transmitter);

    ReportGenerator.SignedReport memory signedReport = s_reportGenerator.generateSignedReport(1, 1);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );
  }

  function test_RevertIf_NotPayee() public {
    _changePrank(address(42));

    vm.expectRevert(DualAggregator.OnlyPayeeCanWithdraw.selector);
    s_aggregator.withdrawPayment(address(42));
  }

  function test_RevertIf_InsufficientFunds() public {
    vm.mockCall(
      address(s_link), abi.encodeWithSelector(LinkTokenInterface.transfer.selector, s_transmitter), abi.encode(false)
    );

    _changePrank(s_transmitter);

    vm.expectRevert(DualAggregator.InsufficientFunds.selector);
    s_aggregator.withdrawPayment(s_transmitter);
  }

  function test_PaysOracles() public {
    // check the balance of the transmitter before the withdrawal
    uint256 balanceBefore = s_link.balanceOf(s_transmitter);
    uint256 amountDue = s_aggregator.owedPayment(s_transmitter);

    // deal the amount due to the aggregator
    deal(address(s_link), address(s_aggregator), amountDue);

    // withdraw the payment
    s_aggregator.withdrawPayment(s_transmitter);

    // check the balance of the transmitter after the withdrawal
    uint256 balanceAfter = s_link.balanceOf(s_transmitter);

    assertEq(balanceAfter, balanceBefore + amountDue, "did not pay the transmitter");
  }
}

contract OwedPayment is ConfiguredDualAggregatorBaseTest {
  address internal s_transmitter;

  function setUp() public override {
    super.setUp();

    s_transmitter = s_aggregator.getTransmitters()[0];

    s_aggregator.setPayees(s_transmitters, s_transmitters);

    s_aggregator.setBilling(10, 1, 100, 100, 0);
  }

  function test_ReturnZeroIfTransmitterNotActive() public view {
    uint256 returnedValue = s_aggregator.owedPayment(address(42));

    assertEq(returnedValue, 0, "did not return 0 when transmitter inactive");
  }

  function test_ReturnOwedAmount() public {
    uint256 returnedValue = s_aggregator.owedPayment(s_transmitters[0]);

    assertEq(returnedValue, 0, "did not return the correct owed amount");

    _changePrank(s_transmitter);

    ReportGenerator.SignedReport memory signedReport = s_reportGenerator.generateSignedReport(1, 1);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );

    assertGt(s_aggregator.owedPayment(s_transmitter), 0, "did not return the correct owed amount");
  }
}

contract WithdrawFunds is ConfiguredDualAggregatorBaseTest {
  address internal constant USER = address(42);

  function setUp() public override {
    super.setUp();
    s_aggregator.setBilling(10, 1, 100, 100, 0);
  }

  function test_RevertIf_NotOwner() public {
    vm.mockCall(
      BILLING_ACCESS_CONTROLLER_ADDRESS,
      abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
      abi.encode(false)
    );
    vm.startPrank(USER);
    vm.expectRevert(DualAggregator.OnlyOwnerAndBillingAdminCanCall.selector);

    s_aggregator.withdrawFunds(USER, 42);
  }

  function test_RevertIf_InsufficientBalance() public {
    ReportGenerator.SignedReport memory signedReport = s_reportGenerator.generateSignedReport(1, 1);

    _changePrank(s_aggregator.getTransmitters()[0]);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );

    _changePrank(s_aggregator.owner());

    vm.expectRevert(DualAggregator.InsufficientBalance.selector);
    s_aggregator.withdrawFunds(USER, 1000 ether);
  }

  function test_RevertIf_InsufficientFunds() public {
    vm.mockCall(
      address(s_link), abi.encodeWithSelector(LinkTokenInterface.transfer.selector, USER, 0), abi.encode(false)
    );

    vm.expectRevert(DualAggregator.InsufficientFunds.selector);

    s_aggregator.withdrawFunds(USER, 1e9);
  }
}

contract LinkAvailableForPayment is ConfiguredDualAggregatorBaseTest {
  uint256 internal LINK_AMOUNT = 1e9;

  function setUp() public override {
    super.setUp();

    deal(address(s_link), address(s_aggregator), LINK_AMOUNT);
  }

  function test_ReturnsBalanceWhenNothingDue() public view {
    assertEq(s_aggregator.linkAvailableForPayment(), int256(LINK_AMOUNT), "did not return the correct balance");
  }

  function test_ReturnsRemainingBalanceWhenHasDues() public {
    _changePrank(s_aggregator.getTransmitters()[0]);

    ReportGenerator.SignedReport memory signedReport = s_reportGenerator.generateSignedReport(1, 1);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );

    // get amount due
    uint256 amountDue = s_aggregator.totalLinkDue();

    // then test that LINK_AMOUNT - AMOUNT_DUE is what gets returned
    assertEq(
      s_aggregator.linkAvailableForPayment(), int256(LINK_AMOUNT - amountDue), "did not return the correct balance"
    );
  }
}

contract OracleObservationCount is ConfiguredDualAggregatorBaseTest {
  function test_ReturnsZeroWhenNoObservations() public view {
    assertEq(s_aggregator.oracleObservationCount(s_transmitters[0]), 0, "did not return 0 for observation count");
  }

  function test_ReturnsCorrectObservationCount() public {
    address transmitter = s_aggregator.getTransmitters()[0];

    _changePrank(transmitter);

    ReportGenerator.SignedReport memory signedReport = s_reportGenerator.generateSignedReport(1, 1);
    s_aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );

    assertEq(s_aggregator.oracleObservationCount(transmitter), 1, "did not return the correct observation count");
  }
}

contract SetPayees is ConfiguredDualAggregatorBaseTest {
  event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current);

  address[] internal payees;

  function setUp() public override {
    super.setUp();
    payees = s_transmitters;
  }

  function test_EmitsPayeeshipTransferred() public {
    vm.expectEmit();
    for (uint256 index = 0; index < s_transmitters.length; index++) {
      address transmitter = s_transmitters[0];
      address payee = payees[0];
      address currentPayee = address(0);
      emit PayeeshipTransferred(transmitter, currentPayee, payee);
    }

    s_aggregator.setPayees(s_transmitters, payees);
  }
}

contract TransferPayeeship is ConfiguredDualAggregatorBaseTest {
  event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed);

  address[] internal payees;
  address internal constant PROPOSED = address(43);

  function setUp() public override {
    super.setUp();
    payees = new address[](s_transmitters.length);

    for (uint256 index = 0; index < s_transmitters.length; index++) {
      payees[index] = address(uint160(1000 + index));
    }

    s_aggregator.setPayees(s_transmitters, payees);
  }

  function test_RevertIf_SenderNotCurrentPayee() public {
    vm.expectRevert(DualAggregator.OnlyCurrentPayeeCanUpdate.selector);

    s_aggregator.transferPayeeship(address(42), address(43));
  }

  function test_RevertIf_SenderIsProposed() public {
    vm.startPrank(payees[0]);
    vm.expectRevert(DualAggregator.CannotTransferPayeeToSelf.selector);

    s_aggregator.transferPayeeship(s_transmitters[0], payees[0]);
  }

  function test_EmitsPayeeshipTransferredRequested() public {
    vm.startPrank(payees[0]);
    vm.expectEmit();
    emit PayeeshipTransferRequested(s_transmitters[0], payees[0], PROPOSED);

    s_aggregator.transferPayeeship(s_transmitters[0], PROPOSED);
  }
}

contract AcceptPayeeship is ConfiguredDualAggregatorBaseTest {
  event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current);

  address[] internal payees;
  address internal constant PROPOSED = address(42);

  function setUp() public override {
    super.setUp();
    payees = new address[](s_transmitters.length);

    for (uint256 index = 0; index < s_transmitters.length; index++) {
      payees[index] = address(uint160(1000 + index));
    }

    s_aggregator.setPayees(s_transmitters, payees);

    vm.startPrank(payees[0]);
    s_aggregator.transferPayeeship(s_transmitters[0], PROPOSED);
    vm.stopPrank();
  }

  function test_RevertIf_SenderIsNotProposed() public {
    vm.startPrank(address(43));
    vm.expectRevert(DualAggregator.OnlyProposedPayeesCanAccept.selector);

    s_aggregator.acceptPayeeship(s_transmitters[0]);
  }

  function test_EmitsPayeeshipTransferred() public {
    vm.startPrank(PROPOSED);
    vm.expectEmit();
    emit PayeeshipTransferred(s_transmitters[0], payees[0], PROPOSED);

    s_aggregator.acceptPayeeship(s_transmitters[0]);
  }
}

contract TypeAndVersion is DualAggregatorBaseTest {
  function test_IsCorrect() public view {
    assertEq(s_aggregator.typeAndVersion(), "DualAggregator 1.0.0", "did not return the right type and version");
  }
}

// _getSyncPrimaryRound(): test all the paths
contract GetSyncPrimaryRound is RoundDataDualAggregatorBaseTest {
  // test with 0 reports transmitted
  function test_zeroTransmissions() public view {
    assertEq(s_aggregator.getSyncPrimaryRound(), 0);
  }

  // test with cutoff time reaching the secondary round id
  function test_returnSecondaryRoundId() public {
    setDualAggregatorBase(30, 20, 6, 2, false);
    assertEq(s_aggregator.getSyncPrimaryRound(), 2);
  }

  // test with cutoff time condition matching in round id 4
  function test_returnSyncFourthRoundId() public {
    setDualAggregatorBase(30, 9, 6, 2, false);

    assertEq(s_aggregator.getSyncPrimaryRound(), 4);
  }

  // test with cutoff time condition matching in the latest round id
  function test_returnSyncLatestRoundId() public {
    setDualAggregatorBase(50, 10, 6, 2, false);

    vm.warp(50);
    assertEq(s_aggregator.getSyncPrimaryRound(), 6);
  }
}
