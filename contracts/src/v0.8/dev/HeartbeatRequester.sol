// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./../interfaces/TypeAndVersionInterface.sol";
import "./interfaces/IAggregatorProxy.sol";
import "./interfaces/IAuthorizedForwarder.sol";
import "./../ConfirmedOwner.sol";

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
   * @dev adds a permitted caller and proxy combination.
   * @param permittedCaller the permitted caller
   * @param proxy the proxy corresponding to this caller
   */
  function addHeartbeat(address permittedCaller, address proxy) external onlyOwner {
    s_heartbeatList[permittedCaller] = IAggregatorProxy(proxy);
    emit HeartbeatSet(permittedCaller, proxy);
  }

  /**
   * @dev removes a permitted caller and proxy combination.
   * @param permittedCaller the permitted caller
   */
  function removeHeartbeat(address permittedCaller) external onlyOwner {
    delete s_heartbeatList[permittedCaller];
    emit HeartbeatRemoved(permittedCaller);
  }

  /**
   * @dev updates the author forwarder address.
   * @param newAuthForwarder the new authorized forwarder address
   */
  function setAuthForwarder(address newAuthForwarder) external onlyOwner {
    i_authForwarder = IAuthorizedForwarder(newAuthForwarder);
    emit AuthorizedForwarderSet(newAuthForwarder);
  }

  /**
   * @dev returns the authorized forwarder.
   */
  function getAuthForwarder() external view returns (address) {
    return address(i_authForwarder);
  }

  /**
   * @dev fetches aggregator address from proxy and forward function call to the aggregator via authorized forwarder.
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
