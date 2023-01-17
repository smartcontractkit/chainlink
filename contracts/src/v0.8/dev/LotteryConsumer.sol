// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../VRFConsumerBaseV2.sol";
import "../ConfirmedOwner.sol";
import "../interfaces/VRFCoordinatorV2Interface.sol";

contract LotteryConsumer is ConfirmedOwner, VRFConsumerBaseV2 {
  struct RequestConfig {
    // VRF v2 keyhash to use.
    bytes32 keyHash;
    // VRF v2 subscription ID.
    uint64 subscriptionId;
    // Minimum number of confirmations for the request.
    uint16 minRequestConfirmations;
    // Callback gas limit.
    uint32 callbackGasLimit;
    // How many random words to get back.
    uint32 numWords;
  }

  RequestConfig internal s_requestConfig;

  struct LotteryRequest {
    // clientRequestId is provided by the client application that is calling
    // into the VRF web2 service.
    bytes32 clientRequestId;
    // lotteryType is provided by the client application that is calling
    // into the VRF web2 service.
    // Type 1 - 5 of 35
    // Type 2 - 6 of 49
    uint8 lotteryType;
    // vrfExternalRequestId is provided by the VRF web2 service.
    // it is a UUID in integer form.
    uint128 vrfExternalRequestId;
  }

  struct LotteryOutcome {
    // clientRequestId is provided by the client application that is calling
    // into the VRF web2 service.
    bytes32 clientRequestId;
    // lotteryType is provided by the client application that is calling
    // into the VRF web2 service.
    // Type 1 - 5 of 35
    // Type 2 - 6 of 49
    uint8 lotteryType;
    // vrfExternalRequestId is provided by the VRF web2 service.
    // it is a UUID in integer form.
    uint128 vrfExternalRequestId;
    // Winning numbers for this lottery.
    uint8[] winningNumbers;
  }

  event LotteryStarted(uint256 indexed vrfRequestId, LotteryRequest request);
  event LotterySettled(uint256 indexed vrfRequestId, LotteryRequest request, LotteryOutcome outcome);
  event AllowedCallerAdded(address caller);
  event AllowedCallerRemoved(address caller);

  mapping(uint256 => LotteryRequest) /* VRF request ID */ /* lottery request object */
    internal s_requests;

  mapping(uint256 => LotteryOutcome) /* VRF request ID */ /* lottery outcome */
    internal s_outcomes;

  mapping(address => bool) /* EOA caller address */ /* is allowed caller */
    internal s_allowedCallers;

  uint256 internal s_mostRecentVrfRequestId;

  error UnallowedCaller(address caller);

  VRFCoordinatorV2Interface internal s_vrfCoordinator;

  constructor(address _vrfCoordinator) ConfirmedOwner(msg.sender) VRFConsumerBaseV2(_vrfCoordinator) {
    require(_vrfCoordinator != address(0), "vrf coordinator must be non-zero");
    s_vrfCoordinator = VRFCoordinatorV2Interface(_vrfCoordinator);
  }

  function requestRandomness(bytes32 clientRequestId, uint8 lotteryType, uint128 vrfExternalRequestId) external {
    // TODO: should we disallow making requests w/ the same client request Id?
    RequestConfig memory config = s_requestConfig;
    uint256 requestId = s_vrfCoordinator.requestRandomWords(
      config.keyHash,
      config.subscriptionId,
      config.minRequestConfirmations,
      config.callbackGasLimit,
      config.numWords);
    LotteryRequest memory request;
    request.clientRequestId = clientRequestId;
    request.lotteryType = lotteryType;
    request.vrfExternalRequestId = vrfExternalRequestId;
    s_requests[requestId] = request;
    emit LotteryStarted(requestId, request);
  }

  function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal override {
    LotteryRequest memory lotteryRequest = s_requests[requestId];
    require(lotteryRequest.lotteryType > 0, "request unrecognized");
    // Winning numbers are the first 5 elements of the shuffled array (in the event of a 5 of 35)
    uint8[] memory shuffle = shuffle35(randomWords);
    uint8[] memory winningNumbers = new uint8[](5); // only 5 of 35 for now
    for (uint256 i = 0; i < winningNumbers.length; i++) {
      winningNumbers[i] = shuffle[i];
    }
    LotteryOutcome memory outcome;
    outcome.clientRequestId = lotteryRequest.clientRequestId;
    outcome.lotteryType = lotteryRequest.lotteryType;
    outcome.vrfExternalRequestId = lotteryRequest.vrfExternalRequestId;
    outcome.winningNumbers = winningNumbers;
    s_outcomes[requestId] = outcome;
    delete s_requests[requestId]; // prevent re-fulfillment
    emit LotterySettled(requestId, lotteryRequest, outcome);
  }

  function getMostRecentVrfRequestId() public view returns (uint256) {
    return s_mostRecentVrfRequestId;
  }

  // ---------- Configuration CRUD  ----------
  function setRequestConfig(
    bytes32 keyHash,
    uint64 subscriptionId,
    uint16 minRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords
  ) external onlyOwner {
    s_requestConfig.keyHash = keyHash;
    s_requestConfig.subscriptionId = subscriptionId;
    s_requestConfig.minRequestConfirmations = minRequestConfirmations;
    s_requestConfig.callbackGasLimit = callbackGasLimit;
    s_requestConfig.numWords = numWords;
  }

  function getRequestConfig() public view returns (
    bytes32 keyHash,
    uint64 subscriptionId,
    uint16 minRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords
  ) {
    RequestConfig memory cfg = s_requestConfig;
    return (
      cfg.keyHash,
      cfg.subscriptionId,
      cfg.minRequestConfirmations,
      cfg.callbackGasLimit,
      cfg.numWords
    );
  }

  // ---------- Allowed caller CRUD ----------

  // function addAllowedCaller(address caller) external onlyOwner {
  //   s_allowedCallers[caller] = true;
  //   emit AllowedCallerAdded(caller);
  // }

  // function removeAllowedCaller(address caller) external onlyOwner {
  //   delete s_allowedCallers[caller];
  //   emit AllowedCallerRemoved(caller);
  // }

  // function isCallerAllowed(address caller) public view returns (bool) {
  //   return s_allowedCallers[caller];
  // }

  // modifier onlyAllowedCallers(address caller) {
  //   if (!s_allowedCallers[caller]) {
  //     revert UnallowedCaller(caller);
  //   }
  //   _;
  // }

  // ---------- Helpers --------------
  function shuffle35(uint256[] memory randomness) public pure returns (uint8[] memory) {
    require(randomness.length == 35, "must be of length 35");
    // Fill array from 1 to 35 inclusive.
    uint8[] memory ret = new uint8[](35);
    for (uint256 i = 0; i < ret.length; i++) {
      ret[i] = uint8(i+1);
    }
    // Shuffle the array using the provided randomness.
    for (uint256 i = 0; i < ret.length; i++) {
      uint256 n = i + (randomness[i] % (35 - i));
      uint8 temp = ret[n];
      ret[n] = ret[i];
      ret[i] = temp;
    }
    return ret;
  }
}
