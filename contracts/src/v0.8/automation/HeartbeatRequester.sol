// SPDX-License-Identifier: MIT
// solhint-disable-next-line one-contract-per-file
pragma solidity 0.8.6;

import {TypeAndVersionInterface} from "./../interfaces/TypeAndVersionInterface.sol";
import {ConfirmedOwner} from "../shared/access/ConfirmedOwner.sol";

// defines some interfaces for type safety and reduces encoding/decoding
// does not use the full interfaces intentionally because the requester only uses a fraction of them
interface IAggregatorProxy {
  function aggregator() external view returns (address);
}

interface IOffchainAggregator {
  function requestNewRound() external returns (uint80);
}

/**
 * @notice The heartbeat requester will maintain a mapping from allowed callers to corresponding proxies. When requested
 *         by eligible caller, it will call a proxy for an aggregator address and request a new round. The aggregator
 *         is gated by permissions and this requester address needs to be whitelisted.
 */
contract HeartbeatRequester is TypeAndVersionInterface, ConfirmedOwner {
  event HeartbeatPermitted(address indexed permittedCaller, address newProxy, address oldProxy);
  event HeartbeatRemoved(address indexed permittedCaller, address removedProxy);

  error HeartbeatNotPermitted();

  mapping(address => IAggregatorProxy) internal s_heartbeatList;

  /**
   * @notice versions:
   * - HeartbeatRequester 1.0.0: The requester fetches the latest aggregator address from proxy, and request a new round
   *                             using the aggregator address.
   */
  string public constant override typeAndVersion = "HeartbeatRequester 1.0.0";

  constructor() ConfirmedOwner(msg.sender) {}

  /**
   * @notice adds a permitted caller and proxy combination.
   * @param permittedCaller the permitted caller
   * @param proxy the proxy corresponding to this caller
   */
  function permitHeartbeat(address permittedCaller, IAggregatorProxy proxy) external onlyOwner {
    address oldProxy = address(s_heartbeatList[permittedCaller]);
    s_heartbeatList[permittedCaller] = proxy;
    emit HeartbeatPermitted(permittedCaller, address(proxy), oldProxy);
  }

  /**
   * @notice removes a permitted caller and proxy combination.
   * @param permittedCaller the permitted caller
   */
  function removeHeartbeat(address permittedCaller) external onlyOwner {
    address removedProxy = address(s_heartbeatList[permittedCaller]);
    delete s_heartbeatList[permittedCaller];
    emit HeartbeatRemoved(permittedCaller, removedProxy);
  }

  /**
   * @notice fetches aggregator address from proxy and requests a new round.
   * @param proxy the proxy address
   */
  function getAggregatorAndRequestHeartbeat(address proxy) external {
    IAggregatorProxy proxyInterface = s_heartbeatList[msg.sender];
    if (address(proxyInterface) != proxy) revert HeartbeatNotPermitted();

    IOffchainAggregator aggregator = IOffchainAggregator(proxyInterface.aggregator());
    aggregator.requestNewRound();
  }
}
