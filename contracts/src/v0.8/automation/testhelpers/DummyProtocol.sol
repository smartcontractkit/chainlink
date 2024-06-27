// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

// this struct is the same as LogTriggerConfig defined in KeeperRegistryLogicA2_1 contract
struct LogTriggerConfig {
  address contractAddress;
  uint8 filterSelector; // denotes which topics apply to filter ex 000, 101, 111...only last 3 bits apply
  bytes32 topic0;
  bytes32 topic1;
  bytes32 topic2;
  bytes32 topic3;
}

contract DummyProtocol {
  event LimitOrderSent(uint256 indexed amount, uint256 indexed price, address indexed to); // keccak256(LimitOrderSent(uint256,uint256,address)) => 0x3e9c37b3143f2eb7e9a2a0f8091b6de097b62efcfe48e1f68847a832e521750a
  event LimitOrderWithdrawn(uint256 indexed amount, uint256 indexed price, address indexed from); // keccak256(LimitOrderWithdrawn(uint256,uint256,address)) => 0x0a71b8ed921ff64d49e4d39449f8a21094f38a0aeae489c3051aedd63f2c229f
  event LimitOrderExecuted(uint256 indexed orderId, uint256 indexed amount, address indexed exchange); // keccak(LimitOrderExecuted(uint256,uint256,address)) => 0xd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd

  function sendLimitedOrder(uint256 amount, uint256 price, address to) public {
    // send an order to an exchange
    emit LimitOrderSent(amount, price, to);
  }

  function withdrawLimit(uint256 amount, uint256 price, address from) public {
    // withdraw an order from an exchange
    emit LimitOrderSent(amount, price, from);
  }

  function executeLimitOrder(uint256 orderId, uint256 amount, address exchange) public {
    // execute a limit order
    emit LimitOrderExecuted(orderId, amount, exchange);
  }

  /**
   * @notice this function generates bytes for a basic log trigger config with no filter selector.
   * @param targetContract the address of contract where events will be emitted from
   * @param t0 the signature of the event to listen to
   */
  function getBasicLogTriggerConfig(
    address targetContract,
    bytes32 t0
  ) external view returns (bytes memory logTrigger) {
    LogTriggerConfig memory cfg = LogTriggerConfig({
      contractAddress: targetContract,
      filterSelector: 0,
      topic0: t0,
      topic1: 0x000000000000000000000000000000000000000000000000000000000000000,
      topic2: 0x000000000000000000000000000000000000000000000000000000000000000,
      topic3: 0x000000000000000000000000000000000000000000000000000000000000000
    });
    return abi.encode(cfg);
  }

  /**
   * @notice this function generates bytes for a customizable log trigger config.
   * @param targetContract the address of contract where events will be emitted from
   * @param selector the filter selector. this denotes which topics apply to filter ex 000, 101, 111....only last 3 bits apply
   * if 0, it won't filter based on topic 1, 2, 3.
   * if 1, it will filter based on topic 1,
   * if 2, it will filter based on topic 2,
   * if 3, it will filter based on topic 1 and topic 2,
   * if 4, it will filter based on topic 3,
   * if 5, it will filter based on topic 1 and topic 3....
   * @param t0 the signature of the event to listen to.
   * @param t1 the topic 1 of the event.
   * @param t2 the topic 2 of the event.
   * @param t3 the topic 2 of the event.
   */
  function getAdvancedLogTriggerConfig(
    address targetContract,
    uint8 selector,
    bytes32 t0,
    bytes32 t1,
    bytes32 t2,
    bytes32 t3
  ) external view returns (bytes memory logTrigger) {
    LogTriggerConfig memory cfg = LogTriggerConfig({
      contractAddress: targetContract,
      filterSelector: selector,
      topic0: t0,
      topic1: t1,
      topic2: t2,
      topic3: t3
    });
    return abi.encode(cfg);
  }
}
