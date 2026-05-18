// SPDX-License-Identifier: MIT

pragma solidity ^0.8.19;

contract FunctionsV1EventsMock {
  // solhint-disable-next-line gas-struct-packing
  struct Config {
    uint16 maxConsumersPerSubscription;
    uint72 adminFee;
    bytes4 handleOracleFulfillmentSelector;
    uint16 gasForCallExactCheck;
    uint32[] maxCallbackGasLimits;
  }

  event ConfigUpdated(Config param1);
  event ContractProposed(
    bytes32 proposedContractSetId,
    address proposedContractSetFromAddress,
    address proposedContractSetToAddress
  );
  event ContractUpdated(bytes32 id, address from, address to);
  event FundsRecovered(address to, uint256 amount);
  event OwnershipTransferRequested(address indexed from, address indexed to);
  event OwnershipTransferred(address indexed from, address indexed to);
  event Paused(address account);
  event RequestNotProcessed(bytes32 indexed requestId, address coordinator, address transmitter, uint8 resultCode);
  event RequestProcessed(
    bytes32 indexed requestId,
    uint64 indexed subscriptionId,
    uint96 totalCostJuels,
    address transmitter,
    uint8 resultCode,
    bytes response,
    bytes err,
    bytes callbackReturnData
  );
  event RequestStart(
    bytes32 indexed requestId,
    bytes32 indexed donId,
    uint64 indexed subscriptionId,
    address subscriptionOwner,
    address requestingContract,
    address requestInitiator,
    bytes data,
    uint16 dataVersion,
    uint32 callbackGasLimit,
    uint96 estimatedTotalCostJuels
  );
  event RequestTimedOut(bytes32 indexed requestId);
  event SubscriptionCanceled(uint64 indexed subscriptionId, address fundsRecipient, uint256 fundsAmount);
  event SubscriptionConsumerAdded(uint64 indexed subscriptionId, address consumer);
  event SubscriptionConsumerRemoved(uint64 indexed subscriptionId, address consumer);
  event SubscriptionCreated(uint64 indexed subscriptionId, address owner);
  event SubscriptionFunded(uint64 indexed subscriptionId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionOwnerTransferRequested(uint64 indexed subscriptionId, address from, address to);
  event SubscriptionOwnerTransferred(uint64 indexed subscriptionId, address from, address to);
  event Unpaused(address account);

  function emitConfigUpdated(Config memory param1) public {
    emit ConfigUpdated(param1);
  }

  function emitContractProposed(
    bytes32 proposedContractSetId,
    address proposedContractSetFromAddress,
    address proposedContractSetToAddress
  ) public {
    emit ContractProposed(proposedContractSetId, proposedContractSetFromAddress, proposedContractSetToAddress);
  }

  function emitContractUpdated(bytes32 id, address from, address to) public {
    emit ContractUpdated(id, from, to);
  }

  function emitFundsRecovered(address to, uint256 amount) public {
    emit FundsRecovered(to, amount);
  }

  function emitOwnershipTransferRequested(address from, address to) public {
    emit OwnershipTransferRequested(from, to);
  }

  function emitOwnershipTransferred(address from, address to) public {
    emit OwnershipTransferred(from, to);
  }

  function emitPaused(address account) public {
    emit Paused(account);
  }

  function emitRequestNotProcessed(
    bytes32 requestId,
    address coordinator,
    address transmitter,
    uint8 resultCode
  ) public {
    emit RequestNotProcessed(requestId, coordinator, transmitter, resultCode);
  }

  function emitRequestProcessed(
    bytes32 requestId,
    uint64 subscriptionId,
    uint96 totalCostJuels,
    address transmitter,
    uint8 resultCode,
    bytes memory response,
    bytes memory err,
    bytes memory callbackReturnData
  ) public {
    emit RequestProcessed(
      requestId,
      subscriptionId,
      totalCostJuels,
      transmitter,
      resultCode,
      response,
      err,
      callbackReturnData
    );
  }

  function emitRequestStart(
    bytes32 requestId,
    bytes32 donId,
    uint64 subscriptionId,
    address subscriptionOwner,
    address requestingContract,
    address requestInitiator,
    bytes memory data,
    uint16 dataVersion,
    uint32 callbackGasLimit,
    uint96 estimatedTotalCostJuels
  ) public {
    emit RequestStart(
      requestId,
      donId,
      subscriptionId,
      subscriptionOwner,
      requestingContract,
      requestInitiator,
      data,
      dataVersion,
      callbackGasLimit,
      estimatedTotalCostJuels
    );
  }

  function emitRequestTimedOut(bytes32 requestId) public {
    emit RequestTimedOut(requestId);
  }

  function emitSubscriptionCanceled(uint64 subscriptionId, address fundsRecipient, uint256 fundsAmount) public {
    emit SubscriptionCanceled(subscriptionId, fundsRecipient, fundsAmount);
  }

  function emitSubscriptionConsumerAdded(uint64 subscriptionId, address consumer) public {
    emit SubscriptionConsumerAdded(subscriptionId, consumer);
  }

  function emitSubscriptionConsumerRemoved(uint64 subscriptionId, address consumer) public {
    emit SubscriptionConsumerRemoved(subscriptionId, consumer);
  }

  function emitSubscriptionCreated(uint64 subscriptionId, address owner) public {
    emit SubscriptionCreated(subscriptionId, owner);
  }

  function emitSubscriptionFunded(uint64 subscriptionId, uint256 oldBalance, uint256 newBalance) public {
    emit SubscriptionFunded(subscriptionId, oldBalance, newBalance);
  }

  function emitSubscriptionOwnerTransferRequested(uint64 subscriptionId, address from, address to) public {
    emit SubscriptionOwnerTransferRequested(subscriptionId, from, to);
  }

  function emitSubscriptionOwnerTransferred(uint64 subscriptionId, address from, address to) public {
    emit SubscriptionOwnerTransferred(subscriptionId, from, to);
  }

  function emitUnpaused(address account) public {
    emit Unpaused(account);
  }
}
