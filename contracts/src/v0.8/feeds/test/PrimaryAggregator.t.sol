// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {PrimaryAggregator} from "../PrimaryAggregator.sol";
import {LinkTokenInterface} from "../../shared/interfaces/LinkTokenInterface.sol";
import {AccessControllerInterface} from "../../shared/interfaces/AccessControllerInterface.sol";
import {AggregatorValidatorInterface} from "../../shared/interfaces/AggregatorValidatorInterface.sol";

contract PrimaryAggregatorHarness is PrimaryAggregator {
  constructor(
    LinkTokenInterface link,
    int192 minAnswer_,
    int192 maxAnswer_,
    AccessControllerInterface billingAccessController,
    AccessControllerInterface requesterAccessController,
    uint8 decimals_,
    string memory description_
  ) PrimaryAggregator(
    link,
    minAnswer_,
    maxAnswer_,
    billingAccessController,
    requesterAccessController,
    decimals_,
    description_
  ) {}

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
}

contract PrimaryAggregatorBaseTest is Test {
  uint256 constant MAX_NUM_ORACLES = 31;

  address constant LINK_TOKEN_ADDRESS = address(1);
  address constant BILLING_ACCESS_CONTROLLER_ADDRESS = address(100);
  address constant REQUESTER_ACCESS_CONTROLLER_ADDRESS = address(101);

  int192 constant MIN_ANSWER = 0;
  int192 constant MAX_ANSWER = 100;

  PrimaryAggregator aggregator;
  PrimaryAggregatorHarness harness;

  function setUp() public virtual {
    LinkTokenInterface _link = LinkTokenInterface(LINK_TOKEN_ADDRESS);
    AccessControllerInterface _billingAccessController = AccessControllerInterface(BILLING_ACCESS_CONTROLLER_ADDRESS);
    AccessControllerInterface _requesterAccessController = AccessControllerInterface(REQUESTER_ACCESS_CONTROLLER_ADDRESS);

    aggregator = new PrimaryAggregator(
      _link,
      MIN_ANSWER,
      MAX_ANSWER,
      _billingAccessController,
      _requesterAccessController,
      18,
      "TEST"
    );
    harness = new PrimaryAggregatorHarness(
      _link,
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
  address[] signers = new address[](MAX_NUM_ORACLES);
  address[] transmitters = new address[](MAX_NUM_ORACLES);
  uint8 f = 1;
  bytes onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
  uint64 offchainConfigVersion = 1;
  bytes offchainConfig = "1";

  function setUp() public override {
    super.setUp();

    for (uint256 i = 0; i<MAX_NUM_ORACLES; i++) {
      signers[i] = address(uint160(1000+i));
      transmitters[i] = address(uint160(2000+i));
    }

    aggregator.setConfig(
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
    assertEq(aggregator.minAnswer(), MIN_ANSWER, "minAnswer not set correctly");
    assertEq(aggregator.maxAnswer(), MAX_ANSWER, "maxAnswer not set correctly");
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

    vm.expectRevert("too many oracles");

    aggregator.setConfig(
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  function test_RevertIf_OracleLengthMismatch() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES - 1);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert("oracle length mismatch");

    aggregator.setConfig(
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  function test_RevertIf_fTooHigh() public {
    address[] memory signers = new address[](1);
    address[] memory transmitters = new address[](1);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert("faulty-oracle f too high");

    aggregator.setConfig(
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  function test_RevertIf_fNotPositive() public {
    address[] memory signers = new address[](1);
    address[] memory transmitters = new address[](1);
    uint8 f = 0;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert("f must be positive");

    aggregator.setConfig(
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  function test_RevertIf_onchainConfigInvalid() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES);
    uint8 f = 1;
    bytes memory onchainConfig = "1";
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    vm.expectRevert("invalid onchainConfig");

    aggregator.setConfig(
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  function test_RevertIf_RepeatedSigner() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES);
    uint8 f = 1;
    bytes memory onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    for (uint256 i = 0; i<MAX_NUM_ORACLES; i++) {
      transmitters[i] = address(uint160(2000+i));
    }

    vm.expectRevert("repeated signer address");

    aggregator.setConfig(
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  function test_RevertIf_RepeatedTransmitter() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES);
    uint8 f = 1;
    bytes memory onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    for (uint256 i = 0; i<MAX_NUM_ORACLES; i++) {
      signers[i] = address(uint160(1000+i));
    }

    vm.expectRevert("repeated transmitter address");

    aggregator.setConfig(
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  function test_HappyPath() public {
    address[] memory signers = new address[](MAX_NUM_ORACLES);
    address[] memory transmitters = new address[](MAX_NUM_ORACLES);
    uint8 f = 1;
    bytes memory onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
    uint64 offchainConfigVersion = 1;
    bytes memory offchainConfig = "1";

    for (uint256 i = 0; i<MAX_NUM_ORACLES; i++) {
      signers[i] = address(uint160(1000+i));
      transmitters[i] = address(uint160(2000+i));
    }

    aggregator.setConfig(
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    assertEq(true, true, "the setConfig transaction rolled back");
  }
}

contract latestConfigDetails is PrimaryAggregatorBaseTest {
  address[] signers = new address[](MAX_NUM_ORACLES);
  address[] transmitters = new address[](MAX_NUM_ORACLES);
  uint8 f = 1;
  bytes onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
  uint64 offchainConfigVersion = 1;
  bytes offchainConfig = "1";

  function setUp() public override {
    super.setUp();

    for (uint256 i = 0; i<MAX_NUM_ORACLES; i++) {
      signers[i] = address(uint160(1000+i));
      transmitters[i] = address(uint160(2000+i));
    }

    aggregator.setConfig(
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  function test_ReturnsConfigDetails() public view {
    (
      uint32 configCount,
      uint32 blockNumber,
      bytes32 configDigest
    ) = aggregator.latestConfigDetails();

    assertEq(configCount, 1, "config count not incremented");
    assertEq(blockNumber, block.number, "block number is wrong");
    assertEq(configDigest, harness.exposed_configDigestFromConfigData(
      block.chainid,
      address(aggregator),
      configCount,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    ), "configDigest is not correct");
  }
}

contract GetTransmitters is ConfiguredPrimaryAggregatorBaseTest {
  function test_ReturnsTransmittersList() public view {
    assertEq(aggregator.getTransmitters(), transmitters, "transmiters list is not the same");
  }
}

contract SetValidatorConfig is ConfiguredPrimaryAggregatorBaseTest {
  event ValidatorConfigSet(
    AggregatorValidatorInterface indexed previousValidator,
    uint32 previousGasLimit,
    AggregatorValidatorInterface indexed currentValidator,
    uint32 currentGasLimit
  );

  AggregatorValidatorInterface oldValidator = AggregatorValidatorInterface(address(0x0));
  AggregatorValidatorInterface newValidator = AggregatorValidatorInterface(address(3001));


  function test_EmitsValidatorConfigSet() public {
    vm.expectEmit();
    emit ValidatorConfigSet(oldValidator, 0, newValidator, 1);

    aggregator.setValidatorConfig(
      newValidator,
      1
    );
  }

}

// TODO: remaining functions to test
// sync up with steve about how verbose we want some of these to be
contract GetValidatorConfig is ConfiguredPrimaryAggregatorBaseTest {}
contract SetRequesterAccessController is ConfiguredPrimaryAggregatorBaseTest {}
contract GetRequesterAccessController is ConfiguredPrimaryAggregatorBaseTest {}
contract RequestNewRound is ConfiguredPrimaryAggregatorBaseTest {}
contract Trasmit is ConfiguredPrimaryAggregatorBaseTest {}
contract LatestTransmissionDetails is ConfiguredPrimaryAggregatorBaseTest {}
contract LatestConfigDigestAndEpoch is ConfiguredPrimaryAggregatorBaseTest {}
contract LatestAnswer is ConfiguredPrimaryAggregatorBaseTest {}
contract LatestTimestamp is ConfiguredPrimaryAggregatorBaseTest {}
contract LatestRound is ConfiguredPrimaryAggregatorBaseTest {}
contract GetAnswer is ConfiguredPrimaryAggregatorBaseTest {}
contract GetTimestamp is ConfiguredPrimaryAggregatorBaseTest {}
contract Description is ConfiguredPrimaryAggregatorBaseTest {}
contract GetRoundData is ConfiguredPrimaryAggregatorBaseTest {}
contract LatestRoundData is ConfiguredPrimaryAggregatorBaseTest {}
contract SetLinkToken is ConfiguredPrimaryAggregatorBaseTest {}
contract GetLinkToken is ConfiguredPrimaryAggregatorBaseTest {}
contract SetBillingAccessController is ConfiguredPrimaryAggregatorBaseTest {}
contract GetBillingAccessController is ConfiguredPrimaryAggregatorBaseTest {}
contract SetBilling is ConfiguredPrimaryAggregatorBaseTest {}
contract GetBilling is ConfiguredPrimaryAggregatorBaseTest {}
contract WithdrawPayment is ConfiguredPrimaryAggregatorBaseTest {}
contract OwedPayment is ConfiguredPrimaryAggregatorBaseTest {}
contract WithdrawFunds is ConfiguredPrimaryAggregatorBaseTest {}
contract LinkAvailableForPayment is ConfiguredPrimaryAggregatorBaseTest {}
contract OracleobservationCount is ConfiguredPrimaryAggregatorBaseTest {}
contract SetPayees is ConfiguredPrimaryAggregatorBaseTest {}
contract TransferPayeeship is ConfiguredPrimaryAggregatorBaseTest {}
contract AcceptPayeeship is ConfiguredPrimaryAggregatorBaseTest {}

contract TypeAndVersion is PrimaryAggregatorBaseTest {
  function test_IsCorrect() public view {
    assertEq(aggregator.typeAndVersion(), "PrimaryAggregator 1.0.0");
  }
}

