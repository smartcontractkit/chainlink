// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Test} from "forge-std/Test.sol";

import {PrimaryAggregator} from "../PrimaryAggregator.sol";

import {LinkTokenInterface} from "../../shared/interfaces/LinkTokenInterface.sol";
import {AccessControllerInterface} from "../../shared/interfaces/AccessControllerInterface.sol";
import {AggregatorValidatorInterface} from "../../shared/interfaces/AggregatorValidatorInterface.sol";
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

  function exposed_totalLinkDue() external view returns (uint256 linkDue) {
    return _totalLinkDue();
  }
}

contract PrimaryAggregatorBaseTest is Test {
  uint256 constant MAX_NUM_ORACLES = 31;

  address constant BILLING_ACCESS_CONTROLLER_ADDRESS = address(100);
  address constant REQUESTER_ACCESS_CONTROLLER_ADDRESS = address(101);

  int192 constant MIN_ANSWER = 0;
  int192 constant MAX_ANSWER = 100;

  LinkToken s_link;
  LinkTokenInterface linkTokenInterface;

  PrimaryAggregator aggregator;
  PrimaryAggregatorHarness harness;

  function setUp() public virtual {
    s_link = new LinkToken();

    linkTokenInterface = LinkTokenInterface(address(s_link));
    AccessControllerInterface _billingAccessController = AccessControllerInterface(BILLING_ACCESS_CONTROLLER_ADDRESS);
    AccessControllerInterface _requesterAccessController = AccessControllerInterface(REQUESTER_ACCESS_CONTROLLER_ADDRESS);

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
  address[] signers = new address[](MAX_NUM_ORACLES);
  address[] transmitters = new address[](MAX_NUM_ORACLES);
  uint8 f = 1;
  bytes onchainConfig = abi.encodePacked(uint8(1), MIN_ANSWER, MAX_ANSWER);
  uint64 offchainConfigVersion = 1;
  bytes offchainConfig = "1";

  function setUp() public override virtual {
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

contract SetValidatorConfig is PrimaryAggregatorBaseTest {
  event ValidatorConfigSet(
    AggregatorValidatorInterface indexed previousValidator,
    uint32 previousGasLimit,
    AggregatorValidatorInterface indexed currentValidator,
    uint32 currentGasLimit
  );

  AggregatorValidatorInterface oldValidator = AggregatorValidatorInterface(address(0x0));
  AggregatorValidatorInterface newValidator = AggregatorValidatorInterface(address(42));


  function test_EmitsValidatorConfigSet() public {
    vm.expectEmit();
    emit ValidatorConfigSet(oldValidator, 0, newValidator, 1);

    aggregator.setValidatorConfig(
      newValidator,
      1
    );
  }

}

contract GetValidatorConfig is PrimaryAggregatorBaseTest {
  AggregatorValidatorInterface newValidator = AggregatorValidatorInterface(address(42));
  uint32 newGasLimit = 1;

  function setUp() public override {
    super.setUp();

    aggregator.setValidatorConfig(
      newValidator,
      newGasLimit
    );
  }

  function test_ReturnsValidatorConfig() public view {
    (AggregatorValidatorInterface returnedValidator, uint32 returnedGasLimit) = aggregator.getValidatorConfig();
    assertEq(address(returnedValidator), address(newValidator), "did not return the right validator");
    assertEq(returnedGasLimit, newGasLimit, "did not return the right gas limit");
  }
}

contract SetRequesterAccessController is PrimaryAggregatorBaseTest {
  event RequesterAccessControllerSet(AccessControllerInterface old, AccessControllerInterface current);

  AccessControllerInterface oldAccessControllerInterface = AccessControllerInterface(REQUESTER_ACCESS_CONTROLLER_ADDRESS);
  AccessControllerInterface newAccessControllerInterface = AccessControllerInterface(address(42));

  function test_EmitsRequesterAccessControllerSet() public {
    vm.expectEmit();
    emit RequesterAccessControllerSet(oldAccessControllerInterface, newAccessControllerInterface);

    aggregator.setRequesterAccessController(newAccessControllerInterface);
  }
}

contract GetRequesterAccessController is PrimaryAggregatorBaseTest {
  AccessControllerInterface newAccessControllerInterface = AccessControllerInterface(address(42));
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

// TODO: this is a big one, come back to it
contract Trasmit is ConfiguredPrimaryAggregatorBaseTest {}

// TODO: once transmit logic is updated we can test these better
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

contract SetLinkToken is PrimaryAggregatorBaseTest {
  event LinkTokenSet(
    LinkTokenInterface indexed oldLinkToken,
    LinkTokenInterface indexed newLinkToken
  );

  LinkToken n_linkToken;
  LinkTokenInterface newLinkToken;

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

  AccessControllerInterface oldBillingAccessController = AccessControllerInterface(BILLING_ACCESS_CONTROLLER_ADDRESS);
  AccessControllerInterface newBillingAccessController = AccessControllerInterface(address(42));

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

  address constant USER = address(42);

  function test_RevertIf_NotOwner() public {
    vm.mockCall(
      BILLING_ACCESS_CONTROLLER_ADDRESS,
      abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
      abi.encode(false)
    );
    vm.startPrank(USER);
    vm.expectRevert("Only owner&billingAdmin can call");

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
    vm.expectRevert("Only payee can withdraw");

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
  address constant USER = address(42);

  function test_RevertIf_NotOwner() public {
    vm.mockCall(
      BILLING_ACCESS_CONTROLLER_ADDRESS,
      abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
      abi.encode(false)
    );
    vm.startPrank(USER);
    vm.expectRevert("Only owner&billingAdmin can call");

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
 
    vm.expectRevert("insufficient funds");

    aggregator.withdrawFunds(USER, 1e9);
  }
}

contract LinkAvailableForPayment is PrimaryAggregatorBaseTest {
  uint256 LINK_AMOUNT = 1e9;

  function setUp() public override {
    super.setUp();

    deal(address(s_link), address(aggregator), LINK_AMOUNT);
  }

  function test_ReturnsBalanceWhenNothingDue() public view {
    assertEq(
      aggregator.linkAvailableForPayment(),
      int256(LINK_AMOUNT),
      "did not return the correct balance"
    );
  }

  function test_ReturnsRemainingBalanceWhenHasDues() public view {
    // TODO: run a transmit so that there is an amount that is due
    // then test that LINK_AMOUNT - AMOUNT_DUE is what gets returned
  }
}

contract OracleObservationCount is ConfiguredPrimaryAggregatorBaseTest {
  function test_ReturnsZeroWhenNoObservations() public view {
    assertEq(
      aggregator.oracleObservationCount(transmitters[0]),
      0,
      "did not return 0 for observation count"
    );
  }
  
  function test_ReturnsCorrectObservationCount() public view {
    // TODO: run a transmit then write this test
  }
}

contract SetPayees is ConfiguredPrimaryAggregatorBaseTest {
  event PayeeshipTransferred(
    address indexed transmitter,
    address indexed previous,
    address indexed current
  );

  address[] payees = transmitters;

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
  event PayeeshipTransferRequested(
    address indexed transmitter,
    address indexed current,
    address indexed proposed
  );

  address[] payees = new address[](transmitters.length);

  address constant PROPOSED = address(43);

  function setUp() public override {
    super.setUp();

    for (uint256 index = 0; index < transmitters.length; index++) {
      payees[index] = address(uint160(1000+index));
    }

    aggregator.setPayees(transmitters, payees);
  }

  function test_RevertIf_SenderNotCurrentPayee() public {
    vm.expectRevert("only current payee can update");

    aggregator.transferPayeeship(address(42), address(43));
  }

  function test_RevertIf_SenderIsProposed() public {
    vm.startPrank(payees[0]);
    vm.expectRevert("cannot transfer to self");

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
  event PayeeshipTransferred(
    address indexed transmitter,
    address indexed previous,
    address indexed current
  );

  address[] payees = new address[](transmitters.length);
  address constant PROPOSED = address(42);

  function setUp() public override {
    super.setUp();

    for (uint256 index = 0; index < transmitters.length; index++) {
      payees[index] = address(uint160(1000+index));
    }

    aggregator.setPayees(transmitters, payees);

    vm.startPrank(payees[0]);
    aggregator.transferPayeeship(transmitters[0], PROPOSED);
    vm.stopPrank();
  }

  function test_RevertIf_SenderIsNotProposed() public {
    vm.startPrank(address(43));
    vm.expectRevert("only proposed payees can accept");

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
    assertEq(
      aggregator.typeAndVersion(),
      "PrimaryAggregator 1.0.0",
      "did not return the right type and version"
    );
  }
}

