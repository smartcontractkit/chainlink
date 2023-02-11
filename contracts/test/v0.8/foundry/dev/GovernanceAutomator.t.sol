pragma solidity ^0.8.0;

import {BaseTest} from "../BaseTest.t.sol";
import {GovernanceAutomator, IGovernorAlphaToken} from "../../../../src/v0.8/upkeeps/GovernanceAutomator.sol";
import {MockComp} from "../../../../src/v0.8/mocks/MockComp.sol";
import {MockTimelock} from "../../../../src/v0.8/mocks/MockTimelock.sol";
import {GovernorAlpha} from "../../../../src/v0.8/vendor/GovernorAlpha.sol";

contract GovernanceAutomatorBaseTest is BaseTest {
  GovernanceAutomator s_governanceAutomator;
  GovernorAlpha s_governorAlpha;
  MockComp s_mockComp;
  MockTimelock s_mockTimelock;

  address internal constant PROPOSER_1 = 0x514910771AF9Ca656af840dff83E8264EcF986CA;
  address internal constant PROPOSER_2 = 0x326C977E6efc84E512bB9C30f76E30c160eD06FB;
  address internal constant PROPOSER_3 = 0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846;

  address internal constant POWER_VOTER = 0xb0897686c545045aFc77CF20eC7A532E3120E0F1;

  address internal constant STUBBED_ADDRESS = 0x0000000000000000000000000000000000000000;
  address internal constant GOVERNOR_ALPHA = 0x0000000000000000000000000000000000000001;
  address internal constant TOKEN_CONTRACT = 0x0000000000000000000000000000000000000002;

  function setUp() public virtual override {
    BaseTest.setUp();

    // Deploy contracts.
    s_mockComp = new MockComp();
    s_mockTimelock = new MockTimelock();
    s_governorAlpha = new GovernorAlpha(address(s_mockTimelock), address(s_mockComp), STUBBED_ADDRESS);
    s_governanceAutomator = new GovernanceAutomator(s_governorAlpha, 0, IGovernorAlphaToken(address(s_mockComp)));

    // Set up proposers and a power voter.
    s_mockComp.setPriorVotes(PROPOSER_1, 100000e18 + 1);
    s_mockComp.setPriorVotes(PROPOSER_2, 100000e18 + 1);
    s_mockComp.setPriorVotes(PROPOSER_3, 100000e18 + 1);
    s_mockComp.setPriorVotes(POWER_VOTER, uint96(s_governorAlpha.quorumVotes() + 1));
  }
}

contract GovernanceAutomatorTest is GovernanceAutomatorBaseTest {
  function testGetStartingIndex() public {
    // Assert that the starting index is 1.
    uint256 startingIndex = s_governanceAutomator.findStartingIndex();
    assertEq(startingIndex, 1);

    // Assert that an upkeep is not needed.
    vm.expectRevert("no action needed");
    s_governanceAutomator.checkUpkeep("");

    // Set the block number to one.
    vm.roll(1);
    vm.warp(1);

    // Create proposal 1.
    changePrank(PROPOSER_1);
    address[] memory targets = new address[](1);
    targets[0] = STUBBED_ADDRESS;
    uint256[] memory values = new uint256[](1);
    values[0] = 0;
    string[] memory signatures = new string[](1);
    signatures[0] = "test_sig";
    bytes[] memory calldatas = new bytes[](1);
    calldatas[0] = "";
    s_governorAlpha.propose(targets, values, signatures, calldatas, "test proposal");

    // Create proposal 2.
    changePrank(PROPOSER_2);
    targets[0] = OWNER;
    values[0] = 0;
    signatures[0] = "test_sig_2";
    calldatas[0] = "";
    s_governorAlpha.propose(targets, values, signatures, calldatas, "test proposal 2");

    // Create proposal 3.
    changePrank(PROPOSER_3);
    targets[0] = TOKEN_CONTRACT;
    values[0] = 0;
    signatures[0] = "test_sig_3";
    calldatas[0] = "";
    s_governorAlpha.propose(targets, values, signatures, calldatas, "test proposal 3");

    // Assert that the proposals are pending.
    require(s_governorAlpha.state(1) == GovernorAlpha.ProposalState.Pending, "incorrect state");
    require(s_governorAlpha.state(2) == GovernorAlpha.ProposalState.Pending, "incorrect state");
    require(s_governorAlpha.state(3) == GovernorAlpha.ProposalState.Pending, "incorrect state");

    // Assert that the starting index is still 1.
    startingIndex = s_governanceAutomator.findStartingIndex();
    assertEq(startingIndex, 1);

    // Assert that an upkeep is not needed.
    vm.expectRevert("no action needed");
    s_governanceAutomator.checkUpkeep("");

    // Enter the voting window for the proposals.
    vm.roll(3);

    // Assert that the proposals are now active.
    require(s_governorAlpha.state(1) == GovernorAlpha.ProposalState.Active, "incorrect state");
    require(s_governorAlpha.state(2) == GovernorAlpha.ProposalState.Active, "incorrect state");
    require(s_governorAlpha.state(3) == GovernorAlpha.ProposalState.Active, "incorrect state");

    // Vote for proposal 2 and against proposal 1.
    changePrank(POWER_VOTER);
    s_governorAlpha.castVote(2, true);
    s_governorAlpha.castVote(1, false);

    // End the voting window.
    vm.roll(s_governorAlpha.votingPeriod() + 3);

    // Proposal 1 has been defeated. Proposal 2 has passed. Assert that the index now needs
    // to be updated, such that proposal 1 is not revisited.
    (bool upkeepNeeded, bytes memory data) = s_governanceAutomator.checkUpkeep("");
    assertEq(upkeepNeeded, true);
    assertEq(data, abi.encode(GovernanceAutomator.Action.UPDATE_INDEX, uint256(2)));

    // Update the index via performUpkeep. Assert that the starting index is now 2.
    s_governanceAutomator.performUpkeep(data);
    startingIndex = s_governanceAutomator.findStartingIndex();
    assertEq(startingIndex, 2);

    // After running performUpkeep, the next action to be completed is to queue
    // Proposal 2. Assert that this is the case.
    (upkeepNeeded, data) = s_governanceAutomator.checkUpkeep("");
    assertEq(upkeepNeeded, true);
    assertEq(data, abi.encode(GovernanceAutomator.Action.QUEUE, uint256(2)));

    // Queue proposal 2. Assert that the queued transaction is correct.
    s_governanceAutomator.performUpkeep(data);
    MockTimelock.Transaction memory queued = s_mockTimelock.getQueuedTransaction();
    assertEq(queued.target, OWNER);
    assertEq(queued.value, 0);
    assertEq(queued.signature, "test_sig_2");

    // After queueing proposal 2, it can be executed. Assert that this is the case.
    (upkeepNeeded, data) = s_governanceAutomator.checkUpkeep("");
    assertEq(upkeepNeeded, true);
    assertEq(data, abi.encode(GovernanceAutomator.Action.EXECUTE, uint256(2)));

    // Execute proposal 2. Assert that the executed transaction is correct.
    s_governanceAutomator.performUpkeep(data);
    MockTimelock.Transaction memory executed = s_mockTimelock.getExecutedTransaction();
    assertEq(executed.target, OWNER);
    assertEq(executed.value, 0);
    assertEq(executed.signature, "test_sig_2");

    // Destroy proposer 3's voting power, making proposal 3 cancellable.
    s_mockComp.setPriorVotes(PROPOSER_3, 0);

    // Proposal 3 is cancellable. Proposal 2 has been executed. Assert that the index now needs
    // to be updated, such that proposal 2 is not revisited.
    (upkeepNeeded, data) = s_governanceAutomator.checkUpkeep("");
    assertEq(upkeepNeeded, true);
    assertEq(data, abi.encode(GovernanceAutomator.Action.UPDATE_INDEX, uint256(3)));

    // Update the index via performUpkeep. Assert that the starting index is now 3.
    s_governanceAutomator.performUpkeep(data);
    startingIndex = s_governanceAutomator.findStartingIndex();
    assertEq(startingIndex, 3);

    // Proposal 3 is cancellable. Assert that this is the case.
    (upkeepNeeded, data) = s_governanceAutomator.checkUpkeep("");
    assertEq(upkeepNeeded, true);
    assertEq(data, abi.encode(GovernanceAutomator.Action.CANCEL, uint256(3)));

    // Cancel proposal 3. Assert that the cancelled transaction is correct.
    s_governanceAutomator.performUpkeep(data);
    MockTimelock.Transaction memory cancelled = s_mockTimelock.getCancelledTransaction();
    assertEq(cancelled.target, TOKEN_CONTRACT);
    assertEq(cancelled.value, 0);
    assertEq(cancelled.signature, "test_sig_3");
  }
}
