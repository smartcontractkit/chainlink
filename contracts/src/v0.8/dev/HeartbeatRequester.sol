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
 *         by eligible caller, it will call a proxy for an aggregator address and call aggregator function to request a
 *         new round.
 */
contract HeartbeatRequester is TypeAndVersionInterface, ConfirmedOwner {
  event HeartbeatSet(address indexed permittedCaller, address indexed proxy);
  event HeartbeatRemoved(address indexed permittedCaller);
  event AuthorizedForwarderSet(address indexed authForwarder);

  error InvalidHeartbeatCombo();

  IAuthorizedForwarder internal i_authForwarder;
  mapping(address => IAggregatorProxy) internal s_heartbeatList;

  /**
   * @notice versions:
   * - HeartbeatRequester 1.0.0: The requester fetches the latest aggregator address from proxy, and request a new round
   *                             from authorized forwarder using the aggregator address.
   */
  string public constant override typeAndVersion = "HeartbeatRequester 1.0.0";

  /**
   * @param authForwarder the authorized forwarder address
   */
  constructor(address authForwarder) ConfirmedOwner(msg.sender) {
    i_authForwarder = IAuthorizedForwarder(authForwarder);
  }

  /**
   * @notice adds a permitted caller and proxy combination.
   * @param permittedCaller the permitted caller
   * @param proxy the proxy corresponding to this caller
   */
  function addHeartbeat(address permittedCaller, IAggregatorProxy proxy) external onlyOwner {
    s_heartbeatList[permittedCaller] = proxy;
    emit HeartbeatSet(permittedCaller, address(proxy));
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
   * @notice updates the authorized forwarder address.
   * @param newAuthForwarder the new authorized forwarder address
   */
  function setAuthForwarder(IAuthorizedForwarder newAuthForwarder) external onlyOwner {
    i_authForwarder = newAuthForwarder;
    emit AuthorizedForwarderSet(address(newAuthForwarder));
  }

  /**
   * @notice returns the authorized forwarder.
   */
  function getAuthForwarder() external view returns (IAuthorizedForwarder) {
    return i_authForwarder;
  }

  /**
   * @notice fetches aggregator address from proxy and forward function call to the aggregator via authorized forwarder.
   * @param proxy the proxy address
   * @param aggregatorFuncSig the function signature on aggregator
   */
  function requestHeartbeat(address proxy, bytes calldata aggregatorFuncSig) external {
    IAggregatorProxy proxyInterface = s_heartbeatList[msg.sender];
    if (address(proxyInterface) != proxy) revert InvalidHeartbeatCombo();

    address aggregator = proxyInterface.aggregator();
    i_authForwarder.forward(aggregator, aggregatorFuncSig);
  }
}
