// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";

import {DualAggregator} from "../DualAggregator.sol";

import {AccessControllerInterface} from "../../shared/interfaces/AccessControllerInterface.sol";
import {AggregatorValidatorInterface} from "../../shared/interfaces/AggregatorValidatorInterface.sol";
import {LinkTokenInterface} from "../../shared/interfaces/LinkTokenInterface.sol";
import {LinkToken} from "../../shared/token/ERC677/LinkToken.sol";
import {ReportGenerator} from "./testhelpers/ReportGenerator.t.sol";

contract DualAggregatorHarness is DualAggregator {
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
    DualAggregator(
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

  function exposed_configDigestFromConfigData(
    uint256 chainId,
    address contractAddress,
    uint64 configCount,
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external pure returns (bytes32) {
    return _configDigestFromConfigData(
      chainId,
      contractAddress,
      configCount,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  function exposed_totalLinkDue() external view returns (uint256 linkDue) {
    return _totalLinkDue();
  }

  function exposed_getSyncPrimaryRound() external view returns (uint80 roundId) {
    return _getSyncPrimaryRound();
  }

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

  // helper function to inject transmissions
  function injectTransmissions(
    int192[] memory answers,
    uint32[] memory observationsTimestamps,
    uint32[] memory recordedTimestamps
  ) public {
    for (uint32 i = 0; i < answers.length; i++) {
      setTransmission(i + 1, answers[i], observationsTimestamps[i], recordedTimestamps[i]);
    }
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
  LinkTokenInterface internal linkTokenInterface;

  DualAggregator internal aggregator;
  DualAggregatorHarness internal harness;

  function setUp() public virtual {
    s_link = new LinkToken();

    linkTokenInterface = LinkTokenInterface(address(s_link));
    AccessControllerInterface _billingAccessController = AccessControllerInterface(BILLING_ACCESS_CONTROLLER_ADDRESS);
    AccessControllerInterface _requesterAccessController =
      AccessControllerInterface(REQUESTER_ACCESS_CONTROLLER_ADDRESS);

    aggregator = new DualAggregator(
      linkTokenInterface,
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
    harness = new DualAggregatorHarness(
      linkTokenInterface,
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
  address[] internal signers = new address[](MAX_NUM_ORACLES);
  address[] internal transmitters = new address[](MAX_NUM_ORACLES);
  uint8 internal f = 1;
  bytes internal onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
  uint64 internal offchainConfigVersion = 1;
  bytes internal offchainConfig = "1";
  bytes32 internal configDigest;
  ReportGenerator internal s_reportGenerator;

  function setUp() public virtual override {
    super.setUp();

    uint256[] memory privateKeys = new uint256[](MAX_NUM_ORACLES);
    for (uint256 i = 0; i < MAX_NUM_ORACLES - 1; i++) {
      uint256 privateKey = uint256(keccak256(abi.encodePacked(i, "oracle-generator-seed")));
      address publicKey = vm.addr(privateKey);
      privateKeys[i] = privateKey;
      transmitters[i] = publicKey;
      signers[i] = publicKey;
    }

    // add the secondary proxy as an approved transmitter/signer
    uint256 secondaryProxyPrivateKey = uint256(keccak256(abi.encodePacked(uint256(102), "oracle-generator-seed")));
    privateKeys[MAX_NUM_ORACLES - 1] = secondaryProxyPrivateKey;
    transmitters[MAX_NUM_ORACLES - 1] = SECONDARY_PROXY;
    signers[MAX_NUM_ORACLES - 1] = SECONDARY_PROXY;

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
    configDigest = harness.exposed_configDigestFromConfigData(
      block.chainid,
      address(aggregator),
      1,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
    s_reportGenerator = new ReportGenerator(aggregator, privateKeys, configDigest, f);
  }
}

contract Constructor is DualAggregatorBaseTest {
  function test_constructor() public view {
    // TODO: add more checks here if we want
    assertEq(aggregator.i_minAnswer(), MIN_ANSWER, "minAnswer not set correctly");
    assertEq(aggregator.i_maxAnswer(), MAX_ANSWER, "maxAnswer not set correctly");
    assertEq(aggregator.decimals(), 18, "decimals not set correctly");
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

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_OracleLengthMismatch() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES - 1);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(DualAggregator.OracleLengthMismatch.selector);

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_fTooHigh() public {
    address[] memory signers = new address[](1);
    address[] memory transmitters = new address[](1);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(DualAggregator.FaultyOracleFTooHigh.selector);

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_fNotPositive() public {
    address[] memory signers = new address[](1);
    address[] memory transmitters = new address[](1);
    uint8 f = 0;
    bytes memory onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(DualAggregator.FMustBePositive.selector);

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_onchainConfigInvalid() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(DualAggregator.InvalidOnChainConfig.selector);

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
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

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
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

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_HappyPath() public {
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

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);

    assertEq(true, true, "the setConfig transaction rolled back");
  }
}

contract LatestConfigDetails is DualAggregatorBaseTest {
  address[] internal signers = new address[](MAX_NUM_ORACLES);
  address[] internal transmitters = new address[](MAX_NUM_ORACLES);
  uint8 internal f = 1;
  bytes internal onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
  uint64 internal offchainConfigVersion = 1;
  bytes internal offchainConfig = "1";

  function setUp() public override {
    super.setUp();

    for (uint256 i = 0; i < MAX_NUM_ORACLES; i++) {
      signers[i] = vm.addr(uint160(1000 + i));
      transmitters[i] = vm.addr(uint160(2000 + i));
    }

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_ReturnsConfigDetails() public view {
    (uint32 configCount, uint32 blockNumber, bytes32 configDigest) = aggregator.latestConfigDetails();

    assertEq(configCount, 1, "config count not incremented");
    assertEq(blockNumber, block.number, "block number is wrong");
    assertEq(
      configDigest,
      harness.exposed_configDigestFromConfigData(
        block.chainid,
        address(aggregator),
        configCount,
        signers,
        transmitters,
        f,
        onchainConfig,
        offchainConfigVersion,
        offchainConfig
      ),
      "configDigest is not correct"
    );
  }
}

contract GetTransmitters is ConfiguredDualAggregatorBaseTest {
  function test_ReturnsTransmittersList() public view {
    assertEq(aggregator.getTransmitters(), transmitters, "transmiters list is not the same");
  }
}

contract SetValidatorConfig is DualAggregatorBaseTest {
  event ValidatorConfigSet(
    AggregatorValidatorInterface indexed previousValidator,
    uint32 previousGasLimit,
    AggregatorValidatorInterface indexed currentValidator,
    uint32 currentGasLimit
  );

  AggregatorValidatorInterface internal oldValidator = AggregatorValidatorInterface(address(0x0));
  AggregatorValidatorInterface internal newValidator = AggregatorValidatorInterface(address(42));

  function test_EmitsValidatorConfigSet() public {
    vm.expectEmit();
    emit ValidatorConfigSet(oldValidator, 0, newValidator, 1);

    aggregator.setValidatorConfig(newValidator, 1);
  }
}

contract GetValidatorConfig is DualAggregatorBaseTest {
  AggregatorValidatorInterface internal newValidator = AggregatorValidatorInterface(address(42));
  uint32 internal newGasLimit = 1;

  function setUp() public override {
    super.setUp();

    aggregator.setValidatorConfig(newValidator, newGasLimit);
  }

  function test_ReturnsValidatorConfig() public view {
    (AggregatorValidatorInterface returnedValidator, uint32 returnedGasLimit) = aggregator.getValidatorConfig();
    assertEq(address(returnedValidator), address(newValidator), "did not return the right validator");
    assertEq(returnedGasLimit, newGasLimit, "did not return the right gas limit");
  }
}

contract SetRequesterAccessController is DualAggregatorBaseTest {
  event RequesterAccessControllerSet(AccessControllerInterface old, AccessControllerInterface current);

  AccessControllerInterface internal oldAccessControllerInterface =
    AccessControllerInterface(REQUESTER_ACCESS_CONTROLLER_ADDRESS);
  AccessControllerInterface internal newAccessControllerInterface = AccessControllerInterface(address(42));

  function test_EmitsRequesterAccessControllerSet() public {
    vm.expectEmit();
    emit RequesterAccessControllerSet(oldAccessControllerInterface, newAccessControllerInterface);

    aggregator.setRequesterAccessController(newAccessControllerInterface);
  }
}

contract GetRequesterAccessController is DualAggregatorBaseTest {
  AccessControllerInterface internal newAccessControllerInterface = AccessControllerInterface(address(42));

  function setUp() public override {
    super.setUp();

    aggregator.setRequesterAccessController(newAccessControllerInterface);
  }

  function test_ReturnsRequesterAccessController() public view {
    assertEq(
      address(aggregator.getRequesterAccessController()),
      address(newAccessControllerInterface),
      "did not return the right access controller interface"
    );
  }
}

// TODO: determine if we need this method still
contract RequestNewRound is ConfiguredDualAggregatorBaseTest {}

contract Transmit is ConfiguredDualAggregatorBaseTest {
  uint32 constant CUTOFF_TIME = 40;

  struct Report {
    int192 price;
    uint32 timestamp;
  }

  function setUp() public override {
    super.setUp();

    aggregator.setCutoffTime(CUTOFF_TIME);
  }

  function test_RevertIf_UnauthorizedTransmitter() public {
    vm.expectRevert(DualAggregator.UnauthorizedTransmitter.selector);
    bytes32[3] memory reportContext =
      [bytes32(abi.encodePacked("1")), bytes32(abi.encodePacked("2")), bytes32(abi.encodePacked("3"))];
    bytes memory report = abi.encodePacked("1");
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](1);
    bytes32[] memory ss = new bytes32[](1);

    rs[0] = bytes32(abi.encodePacked("1"));
    ss[0] = bytes32(abi.encodePacked("1"));

    aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  function test_RevertIf_ConfigDigestMismatch() public {
    vm.startPrank(transmitters[0]);
    vm.expectRevert(DualAggregator.ConfigDigestMismatch.selector);

    bytes32[3] memory reportContext =
      [bytes32(abi.encodePacked("1")), bytes32(abi.encodePacked("2")), bytes32(abi.encodePacked("3"))];
    bytes memory report = abi.encodePacked("1");
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](1);
    bytes32[] memory ss = new bytes32[](1);

    rs[0] = bytes32(abi.encodePacked("1"));
    ss[0] = bytes32(abi.encodePacked("1"));

    aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  function test_RevertIf_CalldataLengthMismatch() public {
    vm.startPrank(transmitters[0]);
    vm.expectRevert(DualAggregator.CalldataLengthMismatch.selector);

    bytes32[3] memory reportContext = [configDigest, bytes32(abi.encodePacked("2")), bytes32(abi.encodePacked("3"))];
    bytes memory report = abi.encodePacked("1");
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](1);
    bytes32[] memory ss = new bytes32[](1);

    rs[0] = bytes32(abi.encodePacked("1"));
    ss[0] = bytes32(abi.encodePacked("1"));

    aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  function test_RevertIf_WrongNumberOfSignatures() public {
    vm.startPrank(transmitters[0]);
    vm.expectRevert(DualAggregator.WrongNumberOfSignatures.selector);

    bytes memory epochAndRound = abi.encodePacked(bytes27(0), uint32(0), uint8(0));

    bytes32[3] memory reportContext = [configDigest, bytes32(epochAndRound), bytes32(abi.encodePacked("1"))];
    bytes memory report = new bytes(0);
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](1);
    bytes32[] memory ss = new bytes32[](1);

    aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  function test_RevertIf_SignaturesOutOfRegistration() public {
    vm.startPrank(transmitters[0]);
    vm.expectRevert(DualAggregator.SignaturesOutOfRegistration.selector);

    bytes memory epochAndRound = abi.encodePacked(bytes27(0), uint32(0), uint8(0));
    bytes32[3] memory reportContext = [configDigest, bytes32(epochAndRound), bytes32(abi.encodePacked("1"))];
    bytes memory report = new bytes(0);
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](1);

    aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  function test_RevertIf_SignatureError() public {
    vm.startPrank(transmitters[0]);

    ReportGenerator.SignedReport memory signedReport =
      s_reportGenerator.generateSignedReport(0, uint32(block.timestamp));

    signedReport.rs[0] = bytes32(abi.encodePacked("1"));
    signedReport.rs[1] = bytes32(abi.encodePacked("1"));
    signedReport.ss[0] = bytes32(abi.encodePacked("1"));
    signedReport.ss[1] = bytes32(abi.encodePacked("1"));

    vm.expectRevert(DualAggregator.SignatureError.selector);
    aggregator.transmit(
      signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
    );
  }

  function test_RevertIf_DuplicateSigner() public {
    vm.startPrank(transmitters[0]);

    ReportGenerator.SignedReport memory signedReport =
      s_reportGenerator.generateSignedReport(0, uint32(block.timestamp));

    signedReport.rs[1] = signedReport.rs[0];
    signedReport.ss[1] = signedReport.ss[0];

    vm.expectRevert(DualAggregator.DuplicateSigner.selector);
    aggregator.transmit(
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
    aggregator.transmit(
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
    aggregator.transmitSecondary(
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
    aggregator.transmit(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );

    vm.expectRevert(DualAggregator.StaleReport.selector);
    aggregator.transmitSecondary(
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
    aggregator.transmitSecondary(
      signedReport1.reportContext, signedReport1.report, signedReport1.rs, signedReport1.ss, signedReport1.rawVs
    );

    vm.expectRevert(DualAggregator.StaleReport.selector);
    aggregator.transmit(
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
    aggregator.transmitSecondary(
      signedReport2.reportContext, signedReport2.report, signedReport2.rs, signedReport2.ss, signedReport2.rawVs
    );

    // check the standard feed
    _changePrank(aggregator.getTransmitters()[0]);
    (, int256 standardAnswer,,,) = aggregator.latestRoundData();
    assertEq(report3.price, standardAnswer, "standard feed answer is not correct");

    // check the secondary feed
    _changePrank(SECONDARY_PROXY);
    (, int256 secondaryAnswer,,,) = aggregator.latestRoundData();
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
    aggregator.transmitSecondary(
      signedReport2.reportContext, signedReport2.report, signedReport2.rs, signedReport2.ss, signedReport2.rawVs
    );

    // check the standard feed
    _changePrank(aggregator.getTransmitters()[0]);
    (, int256 standardAnswer,,,) = aggregator.latestRoundData();
    assertEq(report3.price, standardAnswer, "standard feed answer is not correct");

    // check the secondary feed
    _changePrank(SECONDARY_PROXY);
    (, int256 secondaryAnswer,,,) = aggregator.latestRoundData();
    assertEq(report1.price, secondaryAnswer, "secondary feed answer is not correct");

    _checkRounds(2, 1);
  }

  function test_IncomingSecondaryReportRevertsDueToMaxIterations() public {
    // define 4 as the new max sync iterations
    aggregator.setMaxSyncIterations(4);
    aggregator.setCutoffTime(60);

    // sign the report 1
    ReportGenerator.SignedReport memory signedReport1 =
      s_reportGenerator.generateSignedReport(1, uint32(block.timestamp));
    _mineBlock();

    // transmit the signed report from the primary feed
    _changePrank(aggregator.getTransmitters()[0]);
    aggregator.transmit(
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
    aggregator.transmitSecondary(
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
    _changePrank(aggregator.getTransmitters()[0]);

    if (transmitSecondary) {
      aggregator.transmitSecondary(
        signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
      );
    }

    if (transmitPrimary) {
      aggregator.transmit(
        signedReport.reportContext, signedReport.report, signedReport.rs, signedReport.ss, signedReport.rawVs
      );
    }

    // check the standard feed
    (, int256 standardAnswer,,,) = aggregator.latestRoundData();
    assertEq(int256(expectedStandardFeedAnswer), standardAnswer, "standard feed answer is not correct");

    // check the secondary feed
    _changePrank(SECONDARY_PROXY);
    (, int256 secondaryAnswer,,,) = aggregator.latestRoundData();
    assertEq(int256(expectedSecondaryFeedAnswer), secondaryAnswer, "secondary feed answer is not correct");
  }

  function _checkRounds(uint256 expectedPrimaryRound, uint256 expectedSecondaryRound) internal {
    _changePrank(aggregator.getTransmitters()[0]);
    assertEq(expectedPrimaryRound, aggregator.latestRound(), "standard feed round is not correct");

    _changePrank(SECONDARY_PROXY);
    assertEq(expectedSecondaryRound, aggregator.latestRound(), "secondary feed round is not correct");
  }

  function _mineBlock() internal {
    skip(12);
  }
}

contract TransmittedDualAggregatorBaseTest is ConfiguredDualAggregatorBaseTest {
  bytes32[] internal rs;
  bytes32[] internal ss;
  uint32 internal epoch = 0;
  uint32 internal round = 0;

  // TODO: fix the CalldataLengthMismatch issue
  function setUp() public override {
    super.setUp();

    // vm.startPrank(transmitters[0]);
    // bytes memory epochAndRound = abi.encodePacked(
    //   bytes27(0),
    //   uint32(epoch),
    //   uint32(round)
    // );
    // bytes32[3] memory reportContext = [configDigest, bytes32(epochAndRound), bytes32(abi.encodePacked("1"))];
    // bytes memory report = new bytes(0);
    // bytes32 rawVs = bytes32(abi.encode(uint32(1)));
    //
    // rs.push(bytes32(uint256(uint160(signers[0]))));
    // rs.push(bytes32(uint256(uint160(signers[0]))));
    // ss.push(bytes32(uint256(uint160(signers[0]))));
    // ss.push(bytes32(uint256(uint160(signers[0]))));
    //
    // aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }
}

contract LatestTransmissionDetails is TransmittedDualAggregatorBaseTest {
  function test_RevertIf_NotEOA() public {
    vm.expectRevert(DualAggregator.OnlyCallableByEOA.selector);
    aggregator.latestTransmissionDetails();
  }

  function test_ReturnsLatestTransmissionDetails() public view {
    (bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer, uint64 latestTimestamp) =
      aggregator.latestTransmissionDetails();

    assertEq(configDigest, bytes32(abi.encodePacked("1")));
    assertEq(epoch, 1);
    assertEq(round, 1);
    assertEq(latestAnswer, 1);
    assertEq(latestTimestamp, 1);
  }
}

// TODO: once transmit logic is updated we can test these better
contract LatestConfigDigestAndEpoch is TransmittedDualAggregatorBaseTest {
  function test_ReturnsLatestConfigDigestAndEpoch() public view {
    (bool scanLogs, bytes32 configDigest, uint32 epoch) = aggregator.latestConfigDigestAndEpoch();

    assertEq(scanLogs, false, "scanLogs was not correct");
    assertEq(
      configDigest,
      harness.exposed_configDigestFromConfigData(
        block.chainid,
        address(aggregator),
        1,
        signers,
        transmitters,
        f,
        onchainConfig,
        offchainConfigVersion,
        offchainConfig
      ),
      "configDigest incorrect"
    );
    assertEq(epoch, 1, "epoch not correct");
  }
}

contract RoundDataDualAggregatorBaseTest is ConfiguredDualAggregatorBaseTest {
  int192[] internal answers = [int192(10), int192(11), int192(12), int192(13), int192(14), int192(15)];
  uint32[] internal observationsTimestamps = [uint32(1), uint32(6), uint32(11), uint32(16), uint32(21), uint32(26)];
  uint32[] internal recordedTimestamps = [uint32(5), uint32(10), uint32(15), uint32(20), uint32(25), uint32(30)];

  function setUp() public virtual override {
    super.setUp();
  }

  function setDualAggregatorBase(
    uint256 startingTime,
    uint32 cutoffTime,
    uint32 latestPrimaryRound,
    uint32 latestSecondaryRound
  ) public {
    harness.injectTransmissions(answers, observationsTimestamps, recordedTimestamps);
    harness.setLatestRoundIds(latestPrimaryRound, latestSecondaryRound);
    harness.setCutoffTime(cutoffTime);
    vm.warp(startingTime);
  }
}

// latestAnswer(): test primary and secondary caller
contract LatestAnswer is RoundDataDualAggregatorBaseTest {
  function setUp() public override {
    super.setUp();
    setDualAggregatorBase(31, 21, 6, 2);
  }

  // test return the latest primary answer
  function test_ReturnsLatestPrimaryAnswer() public view {
    assertEq(harness.latestAnswer(), 15);
  }

  // test return the latest secondary answer
  function test_ReturnsLatestSecondaryAnswer() public {
    vm.startPrank(address(102));
    assertEq(harness.latestAnswer(), 11);
  }
}

// latestTimestamp(): test primary and secondary caller
contract LatestTimestamp is RoundDataDualAggregatorBaseTest {
  function setUp() public override {
    super.setUp();
    setDualAggregatorBase(31, 21, 6, 2);
  }

  // test return the latest primary timestamp
  function test_ReturnsLatestPrimaryTimestamp() public view {
    assertEq(harness.latestTimestamp(), 30);
  }

  // test return the latest secondary timestamp
  function test_ReturnsLatestSecondaryTimestamp() public {
    vm.startPrank(address(102));
    assertEq(harness.latestTimestamp(), 10);
  }
}

// latestRound(): test all the paths
contract LatestRound is RoundDataDualAggregatorBaseTest {
  function setUp() public override {
    super.setUp();
    setDualAggregatorBase(30, 20, 6, 2);
  }

  // test return the latest primary round id
  function test_ReturnsLatestPrimaryRoundId() public {
    vm.warp(31);
    assertEq(harness.latestRound(), 6);
  }

  // test latest primary round in the same block, return previous one
  function test_ReturnsPreviousPrimaryRoundId() public view {
    assertEq(harness.latestRound(), 5);
  }

  // test return the latest secondary round id
  function test_ReturnsLatestSecondaryRoundId() public {
    vm.startPrank(address(102));
    assertEq(harness.latestRound(), 2);
  }

  // test return the secondary round id synced with primary rounds
  function test_ReturnsSyncedRoundId() public {
    harness.setCutoffTime(9);

    vm.startPrank(address(102));
    assertEq(harness.latestRound(), 4);
  }
}

// getAnswer(): test primary and secondary caller
contract GetAnswer is RoundDataDualAggregatorBaseTest {
  function setUp() public override {
    super.setUp();
    setDualAggregatorBase(30, 0, 6, 2);
  }

  // test return primary answer
  function test_ReturnsPrimaryGetAnswer() public {
    vm.warp(31);
    assertEq(harness.getAnswer(6), 15);
  }

  // test return primary answer, not allowed
  function test_ReturnsPrimaryGetAnswerNotAllowed() public view {
    assertEq(harness.getAnswer(6), 0);
  }

  // test return secondary answer
  function test_ReturnsSecondaryGetAnswer() public {
    vm.startPrank(address(102));
    assertEq(harness.getAnswer(2), 11);
  }

  // test return secondary answer, not allowed
  function test_ReturnsSecondaryGetAnswerNotAllowed() public {
    harness.setCutoffTime(20);

    vm.startPrank(address(102));
    assertEq(harness.getAnswer(3), 0);
  }
}

// getTimestamp(): test primary and secondary caller
contract GetTimestamp is RoundDataDualAggregatorBaseTest {
  function setUp() public override {
    super.setUp();
    setDualAggregatorBase(30, 0, 6, 2);
  }

  // test return primary timestamp
  function test_ReturnsPrimaryGetTimestamp() public {
    vm.warp(31);
    assertEq(harness.getTimestamp(6), 30);
  }

  // test return primary timestamp, not allowed
  function test_ReturnsPrimaryGetTimestampNotAllowed() public view {
    assertEq(harness.getTimestamp(6), 0);
  }

  // test return secondary timestamp
  function test_ReturnsSecondaryGetTimestamp() public {
    vm.startPrank(address(102));
    assertEq(harness.getTimestamp(2), 10);
  }

  // test return secondary timestamp, not allowed
  function test_ReturnsSecondaryGetTimestampNotAllowed() public {
    harness.setCutoffTime(20);

    vm.startPrank(address(102));
    assertEq(harness.getTimestamp(3), 0);
  }
}

contract Description is TransmittedDualAggregatorBaseTest {
  function test_ReturnsCorrectDescription() public view {
    assertEq(aggregator.description(), "TEST");
  }
}

// getRoundData(): test primary and secondary caller
contract GetRoundData is RoundDataDualAggregatorBaseTest {
  function setUp() public override {
    super.setUp();
    setDualAggregatorBase(30, 0, 6, 2);
  }

  // test return primary round data
  function test_ReturnsPrimaryGetRoundData() public {
    vm.warp(31);
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      harness.getRoundData(6);

    assertEq(roundId, 6);
    assertEq(answer, 15);
    assertEq(startedAt, 26);
    assertEq(updatedAt, 30);
    assertEq(answeredInRound, 6);
  }

  // test return primary round data, not allowed
  function test_ReturnsPrimaryGetRoundDataNotAllowed() public view {
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      harness.getRoundData(6);

    assertEq(roundId, 0);
    assertEq(answer, 0);
    assertEq(startedAt, 0);
    assertEq(updatedAt, 0);
    assertEq(answeredInRound, 0);
  }

  // test return secondary round data
  function test_ReturnsSecondaryGetRoundData() public {
    vm.startPrank(address(102));
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      harness.getRoundData(2);

    assertEq(roundId, 2);
    assertEq(answer, 11);
    assertEq(startedAt, 6);
    assertEq(updatedAt, 10);
    assertEq(answeredInRound, 2);
  }

  // test return secondary round data, not allowed
  function test_ReturnsSecondaryGetRoundDataNotAllowed() public {
    harness.setCutoffTime(20);

    vm.startPrank(address(102));
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      harness.getRoundData(3);

    assertEq(roundId, 0);
    assertEq(answer, 0);
    assertEq(startedAt, 0);
    assertEq(updatedAt, 0);
    assertEq(answeredInRound, 0);
  }
}

// latestRoundData(): test primary and secondary caller
contract LatestRoundData is RoundDataDualAggregatorBaseTest {
  function setUp() public override {
    super.setUp();
    setDualAggregatorBase(31, 21, 6, 2);
  }

  // test return the latest primary round data
  function test_ReturnsLatestPrimaryRoundData() public view {
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      harness.latestRoundData();

    assertEq(roundId, 6);
    assertEq(answer, 15);
    assertEq(startedAt, 26);
    assertEq(updatedAt, 30);
    assertEq(answeredInRound, 6);
  }

  // test return the latest secondary round data
  function test_ReturnsLatestSecondaryRoundData() public {
    vm.startPrank(address(102));
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) =
      harness.latestRoundData();

    assertEq(roundId, 2);
    assertEq(answer, 11);
    assertEq(startedAt, 6);
    assertEq(updatedAt, 10);
    assertEq(answeredInRound, 2);
  }
}

contract SetLinkToken is DualAggregatorBaseTest {
  event LinkTokenSet(LinkTokenInterface indexed oldLinkToken, LinkTokenInterface indexed newLinkToken);

  LinkToken internal n_linkToken;
  LinkTokenInterface internal newLinkToken;

  function setUp() public override {
    super.setUp();
    n_linkToken = new LinkToken();
    newLinkToken = LinkTokenInterface(address(n_linkToken));
  }

  // TODO: determine the right way to make this `transfer` call fail
  // function test_RevertIf_TransferFundsFailed() public {
  //   vm.expectRevert("transfer remaining funds failed");
  //   aggregator.setLinkToken(newLinkToken, address(43));
  // }

  function test_EmitsLinkTokenSet() public {
    deal(address(n_linkToken), address(aggregator), 1e5);
    vm.expectEmit();
    emit LinkTokenSet(linkTokenInterface, newLinkToken);

    aggregator.setLinkToken(newLinkToken, address(43));
  }
}

contract GetLinkToken is DualAggregatorBaseTest {
  function test_ReturnsLinkToken() public view {
    assertEq(
      address(aggregator.getLinkToken()), address(linkTokenInterface), "did not return the right link token interface"
    );
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

    aggregator.setBillingAccessController(newBillingAccessController);
  }
}

contract GetBillingAccessController is DualAggregatorBaseTest {
  function test_ReturnsBillingAccessController() public view {
    assertEq(
      address(aggregator.getBillingAccessController()),
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

    aggregator.setBilling(0, 0, 0, 0, 0);
  }

  function test_EmitsBillingSet() public {
    vm.expectEmit();
    emit BillingSet(0, 0, 0, 0, 0);

    aggregator.setBilling(0, 0, 0, 0, 0);
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
    ) = aggregator.getBilling();

    assertEq(returnedMaxGasPriceGwei, 0, "maxGasPriceGwei incorrect");
    assertEq(returnedReasonableGasPriceGwei, 0, "reasonableGasPriceGwei incorrect");
    assertEq(returnedObservationPaymentGjuels, 0, "observationPaymentGjuels incorrect");
    assertEq(returnedTransmissionPaymentGjuels, 0, "transmissionPaymentGjuels incorrect");
    assertEq(returnedAccountingGas, 0, "accountingGas incorrect");
  }
}

contract WithdrawPayment is ConfiguredDualAggregatorBaseTest {
  function test_RevertIf_NotPayee() public {
    vm.expectRevert(DualAggregator.OnlyPayeeCanWithdraw.selector);

    aggregator.withdrawPayment(address(42));
  }

  function test_PaysOracles() public {
    // TODO: mock and except the call to the mock
  }
}

contract OwedPayment is ConfiguredDualAggregatorBaseTest {
  // TODO: need to figure out a way to toggle the `active` bit on a transmitter
  // right now this is just
  function test_ReturnZeroIfTransmitterNotActive() public view {
    uint256 returnedValue = aggregator.owedPayment(transmitters[0]);

    assertEq(returnedValue, 0, "did not return 0 when transmitter inactive");
  }

  function test_ReturnOwedAmount() public view {
    // TODO: will need to run a transmit here to increase the amount the transmitter is owed
    uint256 returnedValue = aggregator.owedPayment(transmitters[0]);

    assertEq(returnedValue, 0, "did not return the correct owed amount");
  }
}

contract WithdrawFunds is ConfiguredDualAggregatorBaseTest {
  address internal constant USER = address(42);

  function test_RevertIf_NotOwner() public {
    vm.mockCall(
      BILLING_ACCESS_CONTROLLER_ADDRESS,
      abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
      abi.encode(false)
    );
    vm.startPrank(USER);
    vm.expectRevert(DualAggregator.OnlyOwnerAndBillingAdminCanCall.selector);

    aggregator.withdrawFunds(USER, 42);
  }

  // TODO: need to run a transmit to ensure the user has a lot to withdraw
  // function test_RevertIf_InsufficientBalance() public {
  //   vm.expectRevert("insufficient balance");
  //
  //   aggregator.withdrawFunds(USER, 1e9);
  // }

  function test_RevertIf_InsufficientFunds() public {
    vm.mockCall(
      address(s_link), abi.encodeWithSelector(LinkTokenInterface.transfer.selector, USER, 0), abi.encode(false)
    );

    vm.expectRevert(DualAggregator.InsufficientFunds.selector);

    aggregator.withdrawFunds(USER, 1e9);
  }
}

contract LinkAvailableForPayment is DualAggregatorBaseTest {
  uint256 internal LINK_AMOUNT = 1e9;

  function setUp() public override {
    super.setUp();

    deal(address(s_link), address(aggregator), LINK_AMOUNT);
  }

  function test_ReturnsBalanceWhenNothingDue() public view {
    assertEq(aggregator.linkAvailableForPayment(), int256(LINK_AMOUNT), "did not return the correct balance");
  }

  function test_ReturnsRemainingBalanceWhenHasDues() public view {
    // TODO: run a transmit so that there is an amount that is due
    // then test that LINK_AMOUNT - AMOUNT_DUE is what gets returned
  }
}

contract OracleObservationCount is ConfiguredDualAggregatorBaseTest {
  function test_ReturnsZeroWhenNoObservations() public view {
    assertEq(aggregator.oracleObservationCount(transmitters[0]), 0, "did not return 0 for observation count");
  }

  function test_ReturnsCorrectObservationCount() public view {
    // TODO: run a transmit then write this test
  }
}

contract SetPayees is ConfiguredDualAggregatorBaseTest {
  event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current);

  address[] internal payees = transmitters;

  function test_EmitsPayeeshipTransferred() public {
    vm.expectEmit();
    for (uint256 index = 0; index < transmitters.length; index++) {
      address transmitter = transmitters[0];
      address payee = payees[0];
      address currentPayee = address(0);
      emit PayeeshipTransferred(transmitter, currentPayee, payee);
    }

    aggregator.setPayees(transmitters, payees);
  }
}

contract TransferPayeeship is ConfiguredDualAggregatorBaseTest {
  event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed);

  address[] internal payees = new address[](transmitters.length);
  address internal constant PROPOSED = address(43);

  function setUp() public override {
    super.setUp();

    for (uint256 index = 0; index < transmitters.length; index++) {
      payees[index] = address(uint160(1000 + index));
    }

    aggregator.setPayees(transmitters, payees);
  }

  function test_RevertIf_SenderNotCurrentPayee() public {
    vm.expectRevert(DualAggregator.OnlyCurrentPayeeCanUpdate.selector);

    aggregator.transferPayeeship(address(42), address(43));
  }

  function test_RevertIf_SenderIsProposed() public {
    vm.startPrank(payees[0]);
    vm.expectRevert(DualAggregator.CannotTransferToSelf.selector);

    aggregator.transferPayeeship(transmitters[0], payees[0]);
  }

  function test_EmitsPayeeshipTransferredRequested() public {
    vm.startPrank(payees[0]);
    vm.expectEmit();
    emit PayeeshipTransferRequested(transmitters[0], payees[0], PROPOSED);

    aggregator.transferPayeeship(transmitters[0], PROPOSED);
  }
}

contract AcceptPayeeship is ConfiguredDualAggregatorBaseTest {
  event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current);

  address[] internal payees = new address[](transmitters.length);
  address internal constant PROPOSED = address(42);

  function setUp() public override {
    super.setUp();

    for (uint256 index = 0; index < transmitters.length; index++) {
      payees[index] = address(uint160(1000 + index));
    }

    aggregator.setPayees(transmitters, payees);

    vm.startPrank(payees[0]);
    aggregator.transferPayeeship(transmitters[0], PROPOSED);
    vm.stopPrank();
  }

  function test_RevertIf_SenderIsNotProposed() public {
    vm.startPrank(address(43));
    vm.expectRevert(DualAggregator.OnlyProposedPayeesCanAccept.selector);

    aggregator.acceptPayeeship(transmitters[0]);
  }

  function test_EmitsPayeeshipTransferred() public {
    vm.startPrank(PROPOSED);
    vm.expectEmit();
    emit PayeeshipTransferred(transmitters[0], payees[0], PROPOSED);

    aggregator.acceptPayeeship(transmitters[0]);
  }
}

contract TypeAndVersion is DualAggregatorBaseTest {
  function test_IsCorrect() public view {
    assertEq(aggregator.typeAndVersion(), "DualAggregator 1.0.0", "did not return the right type and version");
  }
}

// _getSyncPrimaryRound(): test all the paths
contract GetSyncPrimaryRound is RoundDataDualAggregatorBaseTest {
  function setUp() public override {
    super.setUp();
  }

  // test with 0 reports transmitted
  function test_zeroTransmissions() public view {
    assertEq(harness.exposed_getSyncPrimaryRound(), 0);
  }

  // test with cutoff time reaching the secondary round id
  function test_returnSecondaryRoundId() public {
    setDualAggregatorBase(30, 20, 6, 2);
    assertEq(harness.exposed_getSyncPrimaryRound(), 2);
  }

  // test with cutoff time condition matching in round id 4
  function test_returnSyncFourthRoundId() public {
    setDualAggregatorBase(30, 9, 6, 2);

    assertEq(harness.exposed_getSyncPrimaryRound(), 4);
  }

  // test with cutoff time condition matching in the latest round id
  function test_returnSyncLatestRoundId() public {
    setDualAggregatorBase(50, 10, 6, 2);

    vm.warp(50);
    assertEq(harness.exposed_getSyncPrimaryRound(), 6);
  }
}
