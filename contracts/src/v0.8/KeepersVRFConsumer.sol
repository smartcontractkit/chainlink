// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./automation/KeeperCompatible.sol";
import "./VRFConsumerBaseV2.sol";
import "./interfaces/VRFCoordinatorV2Interface.sol";

/**
 * @title KeepersVRFConsumer
 * @notice KeepersVRFConsumer is a Chainlink Keepers compatible contract that also acts as a
 * VRF V2 requester and consumer. In particular, a random words request is made when `performUpkeep`
 * is called in a cadence provided by the upkeep interval.
 */
contract KeepersVRFConsumer is KeeperCompatibleInterface, VRFConsumerBaseV2 {
  // Upkeep interval in seconds. This contract's performUpkeep method will
  // be called by the Keepers network roughly every UPKEEP_INTERVAL seconds.
  uint256 public immutable UPKEEP_INTERVAL;

  // VRF V2 information, provided upon contract construction.
  VRFCoordinatorV2Interface public immutable COORDINATOR;
  uint64 public immutable SUBSCRIPTION_ID;
  uint16 public immutable REQUEST_CONFIRMATIONS;
  bytes32 public immutable KEY_HASH;

  // Contract state, updated in performUpkeep and fulfillRandomWords.
  uint256 public s_lastTimeStamp;
  uint256 public s_vrfRequestCounter;
  uint256 public s_vrfResponseCounter;

  struct RequestRecord {
    uint256 requestId;
    bool fulfilled;
    uint32 callbackGasLimit;
    uint256 randomness;
  }
  mapping(uint256 => RequestRecord) public s_requests; /* request ID */ /* request record */

  constructor(
    address vrfCoordinator,
    uint64 subscriptionId,
    bytes32 keyHash,
    uint16 requestConfirmations,
    uint256 upkeepInterval
  ) VRFConsumerBaseV2(vrfCoordinator) {
    COORDINATOR = VRFCoordinatorV2Interface(vrfCoordinator);
    SUBSCRIPTION_ID = subscriptionId;
    REQUEST_CONFIRMATIONS = requestConfirmations;
    KEY_HASH = keyHash;
    UPKEEP_INTERVAL = upkeepInterval;

    s_lastTimeStamp = block.timestamp;
    s_vrfRequestCounter = 0;
    s_vrfResponseCounter = 0;
  }

  /**
   * @notice Returns true if and only if at least UPKEEP_INTERVAL seconds have elapsed
   * since the last upkeep or since construction of the contract.
   * @return upkeepNeeded true if and only if at least UPKEEP_INTERVAL seconds have elapsed since the last upkeep or since construction
   * of the contract.
   */
  function checkUpkeep(
    bytes calldata /* checkData */
  )
    external
    view
    override
    returns (
      bool upkeepNeeded,
      bytes memory /* performData */
    )
  {
    upkeepNeeded = (block.timestamp - s_lastTimeStamp) > UPKEEP_INTERVAL;
  }

  /**
   * @notice Requests random words from the VRF coordinator if UPKEEP_INTERVAL seconds have elapsed
   * since the last upkeep or since construction of the contract.
   */
  function performUpkeep(
    bytes calldata /* performData */
  ) external override {
    if ((block.timestamp - s_lastTimeStamp) > UPKEEP_INTERVAL) {
      s_lastTimeStamp = block.timestamp;

      requestRandomWords();
    }
  }

  /**
   * @notice VRF callback implementation
   * @param requestId the VRF V2 request ID, provided at request time.
   * @param randomWords the randomness provided by Chainlink VRF.
   */
  function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal override {
    // Check that the request exists. If not, revert.
    RequestRecord memory record = s_requests[requestId];
    require(record.requestId == requestId, "request ID not found in map");

    // Update the randomness in the record, and increment the response counter.
    s_requests[requestId].randomness = randomWords[0];
    s_vrfResponseCounter++;
  }

  /**
   * @notice Requests random words from Chainlink VRF.
   */
  function requestRandomWords() internal {
    uint256 requestId = COORDINATOR.requestRandomWords(
      KEY_HASH,
      SUBSCRIPTION_ID,
      REQUEST_CONFIRMATIONS,
      150000, // callback gas limit
      1 // num words
    );
    s_requests[requestId] = RequestRecord({
      requestId: requestId,
      fulfilled: false,
      callbackGasLimit: 150000,
      randomness: 0
    });
    s_vrfRequestCounter++;
  }
}
