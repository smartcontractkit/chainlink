// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

import "../interfaces/LinkTokenInterface.sol";
import "../ConfirmedOwner.sol";
import "../interfaces/VRFCoordinatorV2Interface.sol";
import "../interfaces/ERC677ReceiverInterface.sol";
import "../VRFConsumerBaseV2.sol";

contract VRFWebCoordinator is
  ConfirmedOwner,
  AccessControl,
  ERC677ReceiverInterface,
  ReentrancyGuard,
  VRFConsumerBaseV2
{
  bytes32 public constant REQUESTER_ROLE = keccak256("REQUESTER_ROLE");
  bytes32 public constant REGISTER_ROLE = keccak256("REGISTER_ROLE");

  VRFCoordinatorV2 public immutable COORDINATOR;
  LinkTokenInterface public immutable LINK;

  struct Config {
    uint64 subscriptionId;
    bytes32 keyHash;
    uint16 requestConfirmations;
    uint32 callbackGasLimit;
  }
  Config public s_config;

  struct APISubscription {
    bytes32 apiKeyHash;
    uint256 requestCap;
    uint256 requestCount;
  }
  mapping(bytes32 => APISubscription) public s_apiSubscriptions; /* web consumer api key hash */ /* api subscription */

  struct APIRequest {
    bytes32 apiKeyHash;
    bytes32 vrfWebRequestId;
  }
  mapping(uint256 => APIRequest) public s_apiRequests; /* vrf request id */ /* api key hash */

  error InvalidAPISubscription(bytes32 apiKeyHash);
  error NotEnoughRequestAllowance();
  error OnlyCallableFromLink();
  error InvalidCalldata();
  error TopUpTooSmall();
  error TransferAndCallFailed();
  error ApiKeyAlreadyRegistered();
  error InvalidVRFSubscription(uint64 subId);
  error InvalidMaxGasLimit();
  error InvalidRequestConfirmations();
  error InvalidVRFRequestId(uint256 requestId);

  // TODO: what fields in these events should be indexed?
  event VRFWebRandomnessRequested(bytes32 vrfWebRequestId, bytes32 indexed apiKeyHash, uint256 vrfRequestId);
  event VRFWebRandomnessFulfilled(
    bytes32 vrfWebRequestId,
    bytes32 indexed apiKeyHash,
    uint256 vrfRequestId,
    uint256[] randomWords
  );
  event ConfigSet(uint64 subscriptionId, bytes32 keyHash, uint16 requestConfirmations, uint32 callbackGasLimit);
  event ApiKeyHashRegistered(bytes32 indexed apiKeyHash, uint256 requestCap);
  event RequestAllowanceUpdated(bytes32 indexed apiKeyHash, uint256 oldCap, uint256 newCap, uint256 topUpAmountJuels);

  constructor(address vrfCoordinatorV2, address linkToken)
    ConfirmedOwner(msg.sender)
    VRFConsumerBaseV2(vrfCoordinatorV2)
  {
    COORDINATOR = VRFCoordinatorV2(vrfCoordinatorV2);
    LINK = LinkTokenInterface(linkToken);

    // Grant the contract deployer the default admin role: it will be able
    // to grant and revoke any roles.
    _setupRole(DEFAULT_ADMIN_ROLE, msg.sender);
  }

  function requestRandomWords(
    bytes32 apiKeyHash,
    bytes32 vrfWebRequestId,
    uint32 numWords
  ) external returns (uint256) {
    // check if api key exists
    APISubscription memory apiSub = s_apiSubscriptions[apiKeyHash];
    if (apiSub.apiKeyHash == 0x0) {
      revert InvalidAPISubscription(apiKeyHash);
    }

    // check if enough request allowance
    if (apiSub.requestCount >= apiSub.requestCap) {
      revert NotEnoughRequestAllowance();
    }

    Config memory cfg = s_config;

    uint256 vrfRequestId = COORDINATOR.requestRandomWords(
      cfg.keyHash,
      cfg.subscriptionId,
      cfg.requestConfirmations,
      cfg.callbackGasLimit,
      numWords
    );

    s_apiSubscriptions[apiKeyHash].requestCount++;
    s_apiRequests[vrfRequestId] = APIRequest({apiKeyHash: apiKeyHash, vrfWebRequestId: vrfWebRequestId});

    emit VRFWebRandomnessRequested(vrfWebRequestId, apiKeyHash, vrfRequestId);

    return vrfRequestId;
  }

  function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal override {
    APIRequest memory request = s_apiRequests[requestId];
    if (request.apiKeyHash == 0x0) {
      revert InvalidVRFRequestId(requestId);
    }

    delete s_apiRequests[requestId];

    emit VRFWebRandomnessFulfilled(request.apiKeyHash, request.vrfWebRequestId, requestId, randomWords);
  }

  function setConfig(
    uint64 subscriptionId,
    bytes32 keyHash,
    uint16 requestConfirmations,
    uint32 callbackGasLimit
  ) external onlyOwner {
    // validate config against config in the coordinator,
    // so that we don't set invalid parameters.
    (uint16 minimumRequestConfirmations, uint32 maxGasLimit, uint32 _stalenessSeconds, uint32 _gasAfterPaymentCalculation) = COORDINATOR.getConfig();
    if (
      requestConfirmations < minimumRequestConfirmations
    ) {
      revert InvalidRequestConfirmations();
    }

    if (callbackGasLimit > maxGasLimit) {
      revert InvalidMaxGasLimit();
    }

    // check if the given subscription exists in the VRF coordinator.
    try COORDINATOR.getSubscription(subscriptionId) returns (
      uint96, /* balance */
      uint64, /* reqCount */
      address, /* owner */
      address[] memory /* consumers */
    ) {
      // nothing to do, all good.
    } catch {
      // subscription does not exist, have to revert, otherwise all randomness
      // requests will just revert.
      revert InvalidVRFSubscription(subscriptionId);
    }

    s_config = Config({
      subscriptionId: subscriptionId,
      keyHash: keyHash,
      requestConfirmations: requestConfirmations,
      callbackGasLimit: callbackGasLimit
    });

    emit ConfigSet(subscriptionId, keyHash, requestConfirmations, callbackGasLimit);
  }

  function registerApiKey(bytes32 apiKeyHash, uint256 requestCap) external {
    require(hasRole(REGISTER_ROLE, msg.sender), "Caller is not a registerer");
    if (s_apiSubscriptions[apiKeyHash].apiKeyHash != 0x0) {
      revert ApiKeyAlreadyRegistered();
    }

    s_apiSubscriptions[apiKeyHash] = APISubscription({apiKeyHash: apiKeyHash, requestCap: requestCap, requestCount: 0});

    emit ApiKeyHashRegistered(apiKeyHash, requestCap);
  }

  function onTokenTransfer(
    address sender,
    uint256 amountJuels,
    bytes calldata data
  ) external nonReentrant {
    if (msg.sender != address(LINK)) {
      revert OnlyCallableFromLink();
    }
    if (data.length != 256) {
      revert InvalidCalldata();
    }

    bytes32 apiKeyHash = abi.decode(data, (bytes32));
    if (s_apiSubscriptions[apiKeyHash].apiKeyHash == 0x0) {
      revert InvalidAPISubscription(apiKeyHash);
    }

    // As with regular VRF, anyone can fund a particular API key's request allowance.
    uint256 allowance = requestAllowanceFromJuels(amountJuels);
    if (allowance == 0) {
      revert TopUpTooSmall();
    }

    uint256 oldCap = s_apiSubscriptions[apiKeyHash].requestCap;
    s_apiSubscriptions[apiKeyHash].requestCap += allowance;

    // forward the link that was transferred to us to the vrf coordinator to top up the subscription
    bool success = LINK.transferAndCall(address(COORDINATOR), amountJuels, abi.encode(s_config.subscriptionId));
    if (!success) {
      revert TransferAndCallFailed();
    }

    emit RequestAllowanceUpdated(apiKeyHash, oldCap, oldCap + allowance, amountJuels);
  }

  function requestAllowanceFromJuels(uint256 amountJuels) public returns (uint256) {
    // TODO: figure out a good conversion. Will probably need product input.
    return 100;
  }

  // --------- role management for requesters and registerers ---------
  function addRandomnessRequesters(address[] randomnessRequesters) external onlyOwner {
    grantMany(REQUESTER_ROLE, randomnessRequesters);
  }

  function removeRandomnessRequesters(address[] randomnessRequesters) external onlyOwner {
    revokeMany(REQUESTER_ROLE, randomnessRequesters);
  }

  function addApiKeyRegisterers(address[] apiKeyRegisterers) external onlyOwner {
    grantMany(REGISTER_ROLE, apiKeyRegisterers);
  }

  function removeApiKeyRegisterers(address[] apiKeyRegisterers) external onlyOwner {
    revokeMany(REGISTER_ROLE, apiKeyRegisterers);
  }

  function grantMany(bytes32 role, address[] addresses) private {
    for (uint256 i = 0; i < addresses.length; i++) {
      grantRole(role, addresses[i]);
    }
  }

  function revokeMany(bytes32 role, address[] addresses) private {
    for (uint256 i = 0; i < addresses.length; i++) {
      revokeRole(role, addresses[i]);
    }
  }
}

interface VRFCoordinatorV2 is VRFCoordinatorV2Interface {
  function getConfig()
    external
    view
    returns (
      uint16 minimumRequestConfirmations,
      uint32 maxGasLimit,
      uint32 stalenessSeconds,
      uint32 gasAfterPaymentCalculation
    );
}
