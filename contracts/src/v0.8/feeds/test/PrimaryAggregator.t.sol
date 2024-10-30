// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";

import {PrimaryAggregator} from "../PrimaryAggregator.sol";

import {AccessControllerInterface} from "../../shared/interfaces/AccessControllerInterface.sol";
import {AggregatorValidatorInterface} from "../../shared/interfaces/AggregatorValidatorInterface.sol";
import {LinkTokenInterface} from "../../shared/interfaces/LinkTokenInterface.sol";
import {LinkToken} from "../../shared/token/ERC677/LinkToken.sol";

contract PrimaryAggregatorHarness is PrimaryAggregator {
  constructor(
    LinkTokenInterface link,
    int192 minAnswer_,
    int192 maxAnswer_,
    AccessControllerInterface billingAccessController,
    AccessControllerInterface requesterAccessController,
    uint8 decimals_,
    string memory description_
  )
    PrimaryAggregator(
      link,
      minAnswer_,
      maxAnswer_,
      billingAccessController,
      requesterAccessController,
      decimals_,
      description_
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
    return
      _configDigestFromConfigData(
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
}

contract PrimaryAggregatorBaseTest is Test {
  uint256 internal constant MAX_NUM_ORACLES = 31;

  address internal constant BILLING_ACCESS_CONTROLLER_ADDRESS = address(100);
  address internal constant REQUESTER_ACCESS_CONTROLLER_ADDRESS = address(101);

  int192 internal constant MIN_ANSWER = 0;
  int192 internal constant MAX_ANSWER = 100;

  LinkToken internal s_link;
  LinkTokenInterface internal linkTokenInterface;

  PrimaryAggregator internal aggregator;
  PrimaryAggregatorHarness internal harness;

  function setUp() public virtual {
    s_link = new LinkToken();

    linkTokenInterface = LinkTokenInterface(address(s_link));
    AccessControllerInterface _billingAccessController = AccessControllerInterface(BILLING_ACCESS_CONTROLLER_ADDRESS);
    AccessControllerInterface _requesterAccessController = AccessControllerInterface(
      REQUESTER_ACCESS_CONTROLLER_ADDRESS
    );

    aggregator = new PrimaryAggregator(
      linkTokenInterface,
      MIN_ANSWER,
      MAX_ANSWER,
      _billingAccessController,
      _requesterAccessController,
      18,
      "TEST"
    );
    harness = new PrimaryAggregatorHarness(
      linkTokenInterface,
      MIN_ANSWER,
      MAX_ANSWER,
      _billingAccessController,
      _requesterAccessController,
      18,
      "TEST"
    );
  }
}

contract ConfiguredPrimaryAggregatorBaseTest is PrimaryAggregatorBaseTest {
  address[] internal signers = new address[](MAX_NUM_ORACLES);
  address[] internal transmitters = new address[](MAX_NUM_ORACLES);
  uint8 internal f = 1;
  bytes internal onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
  uint64 internal offchainConfigVersion = 1;
  bytes internal offchainConfig = "1";
  bytes32 internal configDigest;

  function setUp() public virtual override {
    super.setUp();

    for (uint256 i = 0; i < MAX_NUM_ORACLES; i++) {
      signers[i] = vm.addr(uint160(1000 + i));
      transmitters[i] = vm.addr(uint160(2000 + i));
    }

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
  }
}

contract Constructor is PrimaryAggregatorBaseTest {
  function test_constructor() public view {
    // TODO: add more checks here if we want
    assertEq(aggregator.i_minAnswer(), MIN_ANSWER, "minAnswer not set correctly");
    assertEq(aggregator.i_maxAnswer(), MAX_ANSWER, "maxAnswer not set correctly");
    assertEq(aggregator.decimals(), 18, "decimals not set correctly");
  }
}

contract SetConfig is PrimaryAggregatorBaseTest {
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

    vm.expectRevert(PrimaryAggregator.TooManyOracles.selector);

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_OracleLengthMismatch() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES - 1);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(PrimaryAggregator.OracleLengthMismatch.selector);

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_fTooHigh() public {
    address[] memory signers = new address[](1);
    address[] memory transmitters = new address[](1);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(PrimaryAggregator.FaultyOracleFTooHigh.selector);

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_fNotPositive() public {
    address[] memory signers = new address[](1);
    address[] memory transmitters = new address[](1);
    uint8 f = 0;
    bytes memory onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(PrimaryAggregator.FMustBePositive.selector);

    aggregator.setConfig(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig);
  }

  function test_RevertIf_onchainConfigInvalid() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert(PrimaryAggregator.InvalidOnChainConfig.selector);

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

    vm.expectRevert(PrimaryAggregator.RepeatedSignerAddress.selector);

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

    vm.expectRevert(PrimaryAggregator.RepeatedTransmitterAddress.selector);

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

contract LatestConfigDetails is PrimaryAggregatorBaseTest {
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

contract GetTransmitters is ConfiguredPrimaryAggregatorBaseTest {
  function test_ReturnsTransmittersList() public view {
    assertEq(aggregator.getTransmitters(), transmitters, "transmiters list is not the same");
  }
}

contract SetValidatorConfig is PrimaryAggregatorBaseTest {
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

contract GetValidatorConfig is PrimaryAggregatorBaseTest {
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

contract SetRequesterAccessController is PrimaryAggregatorBaseTest {
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

contract GetRequesterAccessController is PrimaryAggregatorBaseTest {
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
contract RequestNewRound is ConfiguredPrimaryAggregatorBaseTest {}

contract Trasmit is ConfiguredPrimaryAggregatorBaseTest {
  uint32 epoch = 0;
  uint32 round = 0;

  function setUp() public override {
    super.setUp();
  }

  function test_RevertIf_UnauthorizedTransmitter() public {
    vm.expectRevert(PrimaryAggregator.UnauthorizedTransmitter.selector);
    bytes32[3] memory reportContext = [
      bytes32(abi.encodePacked("1")),
      bytes32(abi.encodePacked("2")),
      bytes32(abi.encodePacked("3"))
    ];
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
    vm.expectRevert(PrimaryAggregator.ConfigDigestMismatch.selector);

    bytes32[3] memory reportContext = [
      bytes32(abi.encodePacked("1")),
      bytes32(abi.encodePacked("2")),
      bytes32(abi.encodePacked("3"))
    ];
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
    vm.expectRevert(PrimaryAggregator.CalldataLengthMismatch.selector);

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
    vm.expectRevert(PrimaryAggregator.WrongNumberOfSignatures.selector);

    bytes memory epochAndRound = abi.encodePacked(bytes27(0), epoch, round);

    bytes32[3] memory reportContext = [configDigest, bytes32(epochAndRound), bytes32(abi.encodePacked("1"))];
    bytes memory report = new bytes(0);
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](1);
    bytes32[] memory ss = new bytes32[](1);

    aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  function test_RevertIf_SignaturesOutOfRegistration() public {
    vm.startPrank(transmitters[0]);
    vm.expectRevert(PrimaryAggregator.SignaturesOutOfRegistration.selector);

    bytes memory epochAndRound = abi.encodePacked(bytes27(0), uint32(epoch), uint32(round));
    bytes32[3] memory reportContext = [configDigest, bytes32(epochAndRound), bytes32(abi.encodePacked("1"))];
    bytes memory report = new bytes(0);
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](1);

    aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  function test_RevertIf_SignatureError() public {
    vm.startPrank(transmitters[0]);
    vm.expectRevert(PrimaryAggregator.SignatureError.selector);

    bytes memory epochAndRound = abi.encodePacked(bytes27(0), uint32(epoch), uint32(round));
    bytes32[3] memory reportContext = [configDigest, bytes32(epochAndRound), bytes32(abi.encodePacked("1"))];
    bytes memory report = new bytes(0);
    bytes32 rawVs = bytes32(abi.encodePacked("1"));
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](2);

    rs[0] = bytes32(abi.encodePacked("1"));
    rs[1] = bytes32(abi.encodePacked("1"));
    ss[0] = bytes32(abi.encodePacked("1"));
    ss[1] = bytes32(abi.encodePacked("1"));

    aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }

  // TODO: figure this test case out better
  // for some reason it thinks the signers aren't active
  function test_RevertIf_DuplicateSigner() public {
    vm.startPrank(transmitters[0]);
    vm.expectRevert(PrimaryAggregator.DuplicateSigner.selector);

    bytes memory epochAndRound = abi.encodePacked(bytes27(0), uint32(epoch), uint32(round));
    bytes32[3] memory reportContext = [configDigest, bytes32(epochAndRound), bytes32(abi.encodePacked("1"))];
    bytes memory report = new bytes(0);

    bytes32 h = keccak256(abi.encode(keccak256(report), reportContext));
    (uint8 v1, bytes32 r1, bytes32 s1) = vm.sign(1, h);
    bytes32[] memory rs = new bytes32[](2);
    bytes32[] memory ss = new bytes32[](2);

    rs[0] = r1;
    rs[1] = r1;
    ss[0] = s1;
    ss[1] = s1;

    bytes32 rawVs = bytes32(uint256(v1));

    aggregator.transmit(reportContext, report, rs, ss, rawVs);
  }
}

contract TransmittedPrimaryAggregatorBaseTest is ConfiguredPrimaryAggregatorBaseTest {
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

contract LatestTransmissionDetails is TransmittedPrimaryAggregatorBaseTest {
  function test_RevertIf_NotEOA() public {
    vm.expectRevert(PrimaryAggregator.OnlyCallableByEOA.selector);
    aggregator.latestTransmissionDetails();
  }

  function test_ReturnsLatestTransmissionDetails() public {
    (bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer, uint64 latestTimestamp) = aggregator
      .latestTransmissionDetails();

    assertEq(configDigest, bytes32(abi.encodePacked("1")));
    assertEq(epoch, 1);
    assertEq(round, 1);
    assertEq(latestAnswer, 1);
    assertEq(latestTimestamp, 1);
  }
}

// TODO: once transmit logic is updated we can test these better
contract LatestConfigDigestAndEpoch is TransmittedPrimaryAggregatorBaseTest {
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
contract LatestAnswer is TransmittedPrimaryAggregatorBaseTest {
  function test_ReturnsLatestAnswer() public view {
    assertEq(aggregator.latestAnswer(), 1);
  }
}
contract LatestTimestamp is TransmittedPrimaryAggregatorBaseTest {
  function test_ReturnsLatestTimestamp() public view {
    assertEq(aggregator.latestTimestamp(), 1);
  }
}
contract LatestRound is TransmittedPrimaryAggregatorBaseTest {
  function test_ReturnsLatestRound() public view {
    assertEq(aggregator.latestRound(), 1);
  }
}
contract GetAnswer is TransmittedPrimaryAggregatorBaseTest {
  function test_ReturnsCorrectAnswer() public view {
    assertEq(aggregator.getAnswer(1), 1);
  }
}
contract GetTimestamp is TransmittedPrimaryAggregatorBaseTest {
  function test_ReturnsCorrectTimestamp() public view {
    assertEq(aggregator.getTimestamp(1), 1);
  }
}
contract Description is TransmittedPrimaryAggregatorBaseTest {
  function test_ReturnsCorrectDescription() public view {
    assertEq(aggregator.description(), "TEST");
  }
}
contract GetRoundData is TransmittedPrimaryAggregatorBaseTest {
  function test_ReturnsCorrectRoundData() public view {
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) = aggregator
      .getRoundData(1);

    assertEq(roundId, 1);
    assertEq(answer, 1);
    assertEq(startedAt, 1);
    assertEq(updatedAt, 1);
    assertEq(answeredInRound, 1);
  }
}
contract LatestRoundData is TransmittedPrimaryAggregatorBaseTest {
  function test_ReturnsLatestRoundData() public view {
    (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) = aggregator
      .latestRoundData();

    assertEq(roundId, 1);
    assertEq(answer, 1);
    assertEq(startedAt, 1);
    assertEq(updatedAt, 1);
    assertEq(answeredInRound, 1);
  }
}

contract SetLinkToken is PrimaryAggregatorBaseTest {
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

contract GetLinkToken is PrimaryAggregatorBaseTest {
  function test_ReturnsLinkToken() public view {
    assertEq(
      address(aggregator.getLinkToken()),
      address(linkTokenInterface),
      "did not return the right link token interface"
    );
  }
}

contract SetBillingAccessController is PrimaryAggregatorBaseTest {
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

contract GetBillingAccessController is PrimaryAggregatorBaseTest {
  function test_ReturnsBillingAccessController() public view {
    assertEq(
      address(aggregator.getBillingAccessController()),
      BILLING_ACCESS_CONTROLLER_ADDRESS,
      "did not return the right billing access controller"
    );
  }
}

contract SetBilling is PrimaryAggregatorBaseTest {
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
    vm.expectRevert(PrimaryAggregator.OnlyOwnerAndBillingAdminCanCall.selector);

    aggregator.setBilling(0, 0, 0, 0, 0);
  }

  function test_EmitsBillingSet() public {
    vm.expectEmit();
    emit BillingSet(0, 0, 0, 0, 0);

    aggregator.setBilling(0, 0, 0, 0, 0);
  }
}

contract GetBilling is PrimaryAggregatorBaseTest {
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

contract WithdrawPayment is ConfiguredPrimaryAggregatorBaseTest {
  function test_RevertIf_NotPayee() public {
    vm.expectRevert(PrimaryAggregator.OnlyPayeeCanWithdraw.selector);

    aggregator.withdrawPayment(address(42));
  }

  function test_PaysOracles() public {
    // TODO: mock and except the call to the mock
  }
}

contract OwedPayment is ConfiguredPrimaryAggregatorBaseTest {
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

contract WithdrawFunds is ConfiguredPrimaryAggregatorBaseTest {
  address internal constant USER = address(42);

  function test_RevertIf_NotOwner() public {
    vm.mockCall(
      BILLING_ACCESS_CONTROLLER_ADDRESS,
      abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
      abi.encode(false)
    );
    vm.startPrank(USER);
    vm.expectRevert(PrimaryAggregator.OnlyOwnerAndBillingAdminCanCall.selector);

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
      address(s_link),
      abi.encodeWithSelector(LinkTokenInterface.transfer.selector, USER, 0),
      abi.encode(false)
    );

    vm.expectRevert(PrimaryAggregator.InsufficientFunds.selector);

    aggregator.withdrawFunds(USER, 1e9);
  }
}

contract LinkAvailableForPayment is PrimaryAggregatorBaseTest {
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

contract OracleObservationCount is ConfiguredPrimaryAggregatorBaseTest {
  function test_ReturnsZeroWhenNoObservations() public view {
    assertEq(aggregator.oracleObservationCount(transmitters[0]), 0, "did not return 0 for observation count");
  }

  function test_ReturnsCorrectObservationCount() public view {
    // TODO: run a transmit then write this test
  }
}

contract SetPayees is ConfiguredPrimaryAggregatorBaseTest {
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

contract TransferPayeeship is ConfiguredPrimaryAggregatorBaseTest {
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
    vm.expectRevert(PrimaryAggregator.OnlyCurrentPayeeCanUpdate.selector);

    aggregator.transferPayeeship(address(42), address(43));
  }

  function test_RevertIf_SenderIsProposed() public {
    vm.startPrank(payees[0]);
    vm.expectRevert(PrimaryAggregator.CannotTransferToSelf.selector);

    aggregator.transferPayeeship(transmitters[0], payees[0]);
  }

  function test_EmitsPayeeshipTransferredRequested() public {
    vm.startPrank(payees[0]);
    vm.expectEmit();
    emit PayeeshipTransferRequested(transmitters[0], payees[0], PROPOSED);

    aggregator.transferPayeeship(transmitters[0], PROPOSED);
  }
}

contract AcceptPayeeship is ConfiguredPrimaryAggregatorBaseTest {
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
    vm.expectRevert(PrimaryAggregator.OnlyProposedPayeesCanAccept.selector);

    aggregator.acceptPayeeship(transmitters[0]);
  }

  function test_EmitsPayeeshipTransferred() public {
    vm.startPrank(PROPOSED);
    vm.expectEmit();
    emit PayeeshipTransferred(transmitters[0], payees[0], PROPOSED);

    aggregator.acceptPayeeship(transmitters[0]);
  }
}

contract TypeAndVersion is PrimaryAggregatorBaseTest {
  function test_IsCorrect() public view {
    assertEq(aggregator.typeAndVersion(), "PrimaryAggregator 1.0.0", "did not return the right type and version");
  }
}
