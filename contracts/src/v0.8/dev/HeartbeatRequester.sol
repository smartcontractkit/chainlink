// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./../interfaces/TypeAndVersionInterface.sol";
import "./../ConfirmedOwner.sol";

// defines some interfaces for type safety and reduces encoding/decoding
// does not use the full interfaces intentionally because the requester only uses a fraction of them
interface IAggregatorProxy {
  function aggregator() external view returns (address);
}

interface IAuthorizedForwarder {
  function forward(address to, bytes calldata data) external;
}

interface IOffchainAggregator {
  function requestNewRound() external returns (uint80);
}

/**
 * @notice The heartbeat requester will maintain a mapping from allowed callers to corresponding proxies. When requested
 *         by eligible caller, it will call a proxy for an aggregator address and call aggregator function via an
 *         authorized forwarder to request a new round.
 */
contract HeartbeatRequester is TypeAndVersionInterface, ConfirmedOwner {
  event HeartbeatPermitted(address indexed permittedCaller, address indexed proxy);
  event HeartbeatRemoved(address indexed permittedCaller);

  error HeartbeatNotPermitted();

  bytes internal constant REQUEST_SELECTOR_CALL_DATA =
    abi.encodeWithSelector(IOffchainAggregator.requestNewRound.selector);
  mapping(address => IAggregatorProxy) internal s_heartbeatList;

  /**
   * @notice versions:
   * - HeartbeatRequester 1.0.0: The requester fetches the latest aggregator address from proxy, and request a new round
   *                             from authorized forwarder using the aggregator address.
   */
  string public constant override typeAndVersion = "HeartbeatRequester 1.0.0";

  constructor() ConfirmedOwner(msg.sender) {}

  /**
   * @notice adds a permitted caller and proxy combination.
   * @param permittedCaller the permitted caller
   * @param proxy the proxy corresponding to this caller
   */
  function permitHeartbeat(address permittedCaller, IAggregatorProxy proxy) external onlyOwner {
    s_heartbeatList[permittedCaller] = proxy;
    emit HeartbeatPermitted(permittedCaller, address(proxy));
  }

  /**
   * @notice removes a permitted caller and proxy combination.
   * @param permittedCaller the permitted caller
   */
  function removeHeartbeat(address permittedCaller) external onlyOwner {
    delete s_heartbeatList[permittedCaller];
    emit HeartbeatRemoved(permittedCaller);
  }

  /**
   * @notice fetches aggregator address from proxy and requests a new round via authorized forwarder.
   * @param proxy the proxy address
   * @param forwarder the forwarder address
   */
  function getAggregatorAndForward(address proxy, IAuthorizedForwarder forwarder) external {
    IAggregatorProxy proxyInterface = s_heartbeatList[msg.sender];
    if (address(proxyInterface) != proxy) revert HeartbeatNotPermitted();

    address aggregator = proxyInterface.aggregator();
    forwarder.forward(aggregator, REQUEST_SELECTOR_CALL_DATA);
  }
}
