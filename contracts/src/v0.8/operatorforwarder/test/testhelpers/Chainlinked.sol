pragma solidity ^0.8.0;

import {ChainlinkClient, Chainlink} from "../../ChainlinkClient.sol";

/**
 * @title The Chainlinked contract
 * @notice Contract writers can inherit this contract in order to create requests for the
 * Chainlink network. ChainlinkClient is an alias of the Chainlinked contract.
 */
// solhint-disable
contract Chainlinked is ChainlinkClient {
  /**
   * @notice Creates a request that can hold additional parameters
   * @param _specId The Job Specification ID that the request will be created for
   * @param _callbackAddress The callback address that the response will be sent to
   * @param _callbackFunctionSignature The callback function signature to use for the callback address
   * @return A Chainlink Request struct in memory
   */
  function newRequest(
    bytes32 _specId,
    address _callbackAddress,
    bytes4 _callbackFunctionSignature
  ) internal pure returns (Chainlink.Request memory) {
    return _buildChainlinkRequest(_specId, _callbackAddress, _callbackFunctionSignature);
  }

  /**
   * @notice Creates a Chainlink request to the stored oracle address
   * @dev Calls `sendChainlinkRequestTo` with the stored oracle address
   * @param _req The initialized Chainlink Request
   * @param _payment The amount of LINK to send for the request
   * @return The request ID
   */
  function chainlinkRequest(Chainlink.Request memory _req, uint256 _payment) internal returns (bytes32) {
    return _sendChainlinkRequest(_req, _payment);
  }

  /**
   * @notice Creates a Chainlink request to the specified oracle address
   * @dev Generates and stores a request ID, increments the local nonce, and uses `transferAndCall` to
   * send LINK which creates a request on the target oracle contract.
   * Emits ChainlinkRequested event.
   * @param _oracle The address of the oracle for the request
   * @param _req The initialized Chainlink Request
   * @param _payment The amount of LINK to send for the request
   * @return requestId The request ID
   */
  function chainlinkRequestTo(
    address _oracle,
    Chainlink.Request memory _req,
    uint256 _payment
  ) internal returns (bytes32 requestId) {
    return _sendChainlinkRequestTo(_oracle, _req, _payment);
  }

  /**
   * @notice Sets the stored oracle address
   * @param _oracle The address of the oracle contract
   */
  function setOracle(address _oracle) internal {
    _setChainlinkOracle(_oracle);
  }

  /**
   * @notice Sets the LINK token address
   * @param _link The address of the LINK token contract
   */
  function setLinkToken(address _link) internal {
    _setChainlinkToken(_link);
  }

  /**
   * @notice Retrieves the stored address of the LINK token
   * @return The address of the LINK token
   */
  function chainlinkToken() internal view returns (address) {
    return _chainlinkTokenAddress();
  }

  /**
   * @notice Retrieves the stored address of the oracle contract
   * @return The address of the oracle contract
   */
  function oracleAddress() internal view returns (address) {
    return _chainlinkOracleAddress();
  }

  /**
   * @notice Ensures that the fulfillment is valid for this contract
   * @dev Use if the contract developer prefers methods instead of modifiers for validation
   * @param _requestId The request ID for fulfillment
   */
  function fulfillChainlinkRequest(
    bytes32 _requestId
  )
    internal
    recordChainlinkFulfillment(_requestId) // solhint-disable-next-line no-empty-blocks
  {}

  /**
   * @notice Sets the stored oracle and LINK token contracts with the addresses resolved by ENS
   * @dev Accounts for subnodes having different resolvers
   * @param _ens The address of the ENS contract
   * @param _node The ENS node hash
   */
  function setChainlinkWithENS(address _ens, bytes32 _node) internal {
    _useChainlinkWithENS(_ens, _node);
  }

  /**
   * @notice Sets the stored oracle contract with the address resolved by ENS
   * @dev This may be called on its own as long as `setChainlinkWithENS` has been called previously
   */
  function setOracleWithENS() internal {
    _updateChainlinkOracleWithENS();
  }

  /**
   * @notice Allows for a request which was created on another contract to be fulfilled
   * on this contract
   * @param _oracle The address of the oracle contract that will fulfill the request
   * @param _requestId The request ID used for the response
   */
  function addExternalRequest(address _oracle, bytes32 _requestId) internal {
    _addChainlinkExternalRequest(_oracle, _requestId);
  }
}
