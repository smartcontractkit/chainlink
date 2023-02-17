// SPDX-License-Identifier: MIT

/**
  The Cron contract is a chainlink keepers-powered cron job runner for smart contracts.
  The contract enables developers to trigger actions on various targets using cron
  strings to specify the cadence. For example, a user may have 3 tasks that require
  regular service in their dapp ecosystem:
    1) 0xAB..CD, update(1), "0 0 * * *"     --> runs update(1) on 0xAB..CD daily at midnight
    2) 0xAB..CD, update(2), "30 12 * * 0-4" --> runs update(2) on 0xAB..CD weekdays at 12:30
    3) 0x12..34, trigger(), "0 * * * *"     --> runs trigger() on 0x12..34 hourly

  To use this contract, a user first deploys this contract and registers it on the chainlink
  keeper registry. Then the user adds cron jobs by following these steps:
    1) Convert a cron string to an encoded cron spec by calling encodeCronString()
    2) Take the encoding, target, and handler, and create a job by sending a tx to createCronJob()
    3) Cron job is running :)
*/

pragma solidity 0.8.6;

import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/proxy/Proxy.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "../../ConfirmedOwner.sol";
import "../KeeperBase.sol";
import "../../interfaces/automation/KeeperCompatibleInterface.sol";
import {Cron as CronInternal, Spec} from "../../libraries/internal/Cron.sol";
import {Cron as CronExternal} from "../../libraries/external/Cron.sol";
import {getRevertMsg} from "../../utils/utils.sol";

/**
 * @title The CronUpkeep contract
 * @notice A keeper-compatible contract that runs various tasks on cron schedules.
 * Users must use the encodeCronString() function to encode their cron jobs before
 * setting them. This keeps all the string manipulation off chain and reduces gas costs.
 */
contract CronUpkeep is KeeperCompatibleInterface, KeeperBase, ConfirmedOwner, Pausable, Proxy {
  using EnumerableSet for EnumerableSet.UintSet;

  event CronJobExecuted(uint256 indexed id, uint256 timestamp);
  event CronJobCreated(uint256 indexed id, address target, bytes handler);
  event CronJobUpdated(uint256 indexed id, address target, bytes handler);
  event CronJobDeleted(uint256 indexed id);

  error CallFailed(uint256 id, string reason);
  error CronJobIDNotFound(uint256 id);
  error ExceedsMaxJobs();
  error InvalidHandler();
  error TickInFuture();
  error TickTooOld();
  error TickDoesntMatchSpec();

  address immutable s_delegate;
  uint256 public immutable s_maxJobs;
  uint256 private s_nextCronJobID = 1;
  EnumerableSet.UintSet private s_activeCronJobIDs;

  mapping(uint256 => uint256) private s_lastRuns;
  mapping(uint256 => Spec) private s_specs;
  mapping(uint256 => address) private s_targets;
  mapping(uint256 => bytes) private s_handlers;
  mapping(uint256 => bytes32) private s_handlerSignatures;

  /**
   * @param owner the initial owner of the contract
   * @param delegate the contract to delegate checkUpkeep calls to
   * @param maxJobs the max number of cron jobs this contract will support
   * @param firstJob an optional encoding of the first cron job
   */
  constructor(
    address owner,
    address delegate,
    uint256 maxJobs,
    bytes memory firstJob
  ) ConfirmedOwner(owner) {
    s_delegate = delegate;
    s_maxJobs = maxJobs;
    if (firstJob.length > 0) {
      (address target, bytes memory handler, Spec memory spec) = abi.decode(firstJob, (address, bytes, Spec));
      createCronJobFromSpec(target, handler, spec);
    }
  }

  /**
   * @notice Executes the cron job with id encoded in performData
   * @param performData abi encoding of cron job ID and the cron job's next run-at datetime
   */
  function performUpkeep(bytes calldata performData) external override whenNotPaused {
    (uint256 id, uint256 tickTime, address target, bytes memory handler) = abi.decode(
      performData,
      (uint256, uint256, address, bytes)
    );
    validate(id, tickTime, target, handler);
    s_lastRuns[id] = block.timestamp;
    (bool success, bytes memory payload) = target.call(handler);
    if (!success) {
      revert CallFailed(id, getRevertMsg(payload));
    }
    emit CronJobExecuted(id, block.timestamp);
  }

  /**
   * @notice Creates a cron job from the given encoded spec
   * @param target the destination contract of a cron job
   * @param handler the function signature on the target contract to call
   * @param encodedCronSpec abi encoding of a cron spec
   */
  function createCronJobFromEncodedSpec(
    address target,
    bytes memory handler,
    bytes memory encodedCronSpec
  ) external onlyOwner {
    if (s_activeCronJobIDs.length() >= s_maxJobs) {
      revert ExceedsMaxJobs();
    }
    Spec memory spec = abi.decode(encodedCronSpec, (Spec));
    createCronJobFromSpec(target, handler, spec);
  }

  /**
   * @notice Updates a cron job from the given encoded spec
   * @param id the id of the cron job to update
   * @param newTarget the destination contract of a cron job
   * @param newHandler the function signature on the target contract to call
   * @param newEncodedCronSpec abi encoding of a cron spec
   */
  function updateCronJob(
    uint256 id,
    address newTarget,
    bytes memory newHandler,
    bytes memory newEncodedCronSpec
  ) external onlyOwner onlyValidCronID(id) {
    Spec memory newSpec = abi.decode(newEncodedCronSpec, (Spec));
    s_targets[id] = newTarget;
    s_handlers[id] = newHandler;
    s_specs[id] = newSpec;
    s_handlerSignatures[id] = handlerSig(newTarget, newHandler);
    emit CronJobUpdated(id, newTarget, newHandler);
  }

  /**
   * @notice Deletes the cron job matching the provided id. Reverts if
   * the id is not found.
   * @param id the id of the cron job to delete
   */
  function deleteCronJob(uint256 id) external onlyOwner onlyValidCronID(id) {
    delete s_lastRuns[id];
    delete s_specs[id];
    delete s_targets[id];
    delete s_handlers[id];
    delete s_handlerSignatures[id];
    s_activeCronJobIDs.remove(id);
    emit CronJobDeleted(id);
  }

  /**
   * @notice Pauses the contract, which prevents executing performUpkeep
   */
  function pause() external onlyOwner {
    _pause();
  }

  /**
   * @notice Unpauses the contract
   */
  function unpause() external onlyOwner {
    _unpause();
  }

  /**
   * @notice Get the id of an eligible cron job
   * @return upkeepNeeded signals if upkeep is needed, performData is an abi encoding
   * of the id and "next tick" of the elligible cron job
   */
  function checkUpkeep(bytes calldata) external override whenNotPaused cannotExecute returns (bool, bytes memory) {
    _delegate(s_delegate);
  }

  /**
   * @notice gets a list of active cron job IDs
   * @return list of active cron job IDs
   */
  function getActiveCronJobIDs() external view returns (uint256[] memory) {
    uint256 length = s_activeCronJobIDs.length();
    uint256[] memory jobIDs = new uint256[](length);
    for (uint256 idx = 0; idx < length; idx++) {
      jobIDs[idx] = s_activeCronJobIDs.at(idx);
    }
    return jobIDs;
  }

  /**
   * @notice gets a cron job
   * @param id the cron job ID
   * @return target - the address a cron job forwards the eth tx to
             handler - the encoded function sig to execute when forwarding a tx
             cronString - the string representing the cron job
             nextTick - the timestamp of the next time the cron job will run
   */
  function getCronJob(uint256 id)
    external
    view
    onlyValidCronID(id)
    returns (
      address target,
      bytes memory handler,
      string memory cronString,
      uint256 nextTick
    )
  {
    Spec memory spec = s_specs[id];
    return (s_targets[id], s_handlers[id], CronExternal.toCronString(spec), CronExternal.nextTick(spec));
  }

  /**
   * @notice Adds a cron spec to storage and the ID to the list of jobs
   * @param target the destination contract of a cron job
   * @param handler the function signature on the target contract to call
   * @param spec the cron spec to create
   */
  function createCronJobFromSpec(
    address target,
    bytes memory handler,
    Spec memory spec
  ) internal {
    uint256 newID = s_nextCronJobID;
    s_activeCronJobIDs.add(newID);
    s_targets[newID] = target;
    s_handlers[newID] = handler;
    s_specs[newID] = spec;
    s_lastRuns[newID] = block.timestamp;
    s_handlerSignatures[newID] = handlerSig(target, handler);
    s_nextCronJobID++;
    emit CronJobCreated(newID, target, handler);
  }

  function _implementation() internal view override returns (address) {
    return s_delegate;
  }

  /**
   * @notice validates the input to performUpkeep
   * @param id the id of the cron job
   * @param tickTime the observed tick time
   * @param target the contract to forward the tx to
   * @param handler the handler of the contract receiving the forwarded tx
   */
  function validate(
    uint256 id,
    uint256 tickTime,
    address target,
    bytes memory handler
  ) private {
    tickTime = tickTime - (tickTime % 60); // remove seconds from tick time
    if (block.timestamp < tickTime) {
      revert TickInFuture();
    }
    if (tickTime <= s_lastRuns[id]) {
      revert TickTooOld();
    }
    if (!CronInternal.matches(s_specs[id], tickTime)) {
      revert TickDoesntMatchSpec();
    }
    if (handlerSig(target, handler) != s_handlerSignatures[id]) {
      revert InvalidHandler();
    }
  }

  /**
   * @notice returns a unique identifier for target/handler pairs
   * @param target the contract to forward the tx to
   * @param handler the handler of the contract receiving the forwarded tx
   * @return a hash of the inputs
   */
  function handlerSig(address target, bytes memory handler) private pure returns (bytes32) {
    return keccak256(abi.encodePacked(target, handler));
  }

  modifier onlyValidCronID(uint256 id) {
    if (!s_activeCronJobIDs.contains(id)) {
      revert CronJobIDNotFound(id);
    }
    _;
  }
}
