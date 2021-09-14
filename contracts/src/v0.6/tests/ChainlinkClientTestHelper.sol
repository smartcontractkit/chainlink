// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "../ChainlinkClient.sol";
import "../vendor/SafeMathChainlink.sol";

contract ChainlinkClientTestHelper is ChainlinkClient {
  using SafeMathChainlink for uint256;

  constructor(
    address _link,
    address _oracle
  ) public {
    setChainlinkToken(_link);
    setChainlinkOracle(_oracle);
  }

  event Request(
    bytes32 id,
    address callbackAddress,
    bytes4 callbackfunctionSelector,
    bytes data
  );
  event LinkAmount(
    uint256 amount
  );

  function publicNewRequest(
    bytes32 _id,
    address _address,
    bytes memory _fulfillmentSignature
  )
    public
  {
    Chainlink.Request memory req = buildChainlinkRequest(
      _id, _address, bytes4(keccak256(_fulfillmentSignature)));
    emit Request(
      req.id,
      req.callbackAddress,
      req.callbackFunctionId,
      req.buf.buf
    );
  }

  function publicRequest(
    bytes32 _id,
    address _address,
    bytes memory _fulfillmentSignature,
    uint256 _wei
  )
    public
  {
    Chainlink.Request memory req = buildChainlinkRequest(
      _id, _address, bytes4(keccak256(_fulfillmentSignature)));
    sendChainlinkRequest(req, _wei);
  }

  function publicRequestRunTo(
    address _oracle,
    bytes32 _id,
    address _address,
    bytes memory _fulfillmentSignature,
    uint256 _wei
  )
    public
  {
    Chainlink.Request memory run = buildChainlinkRequest(_id, _address, bytes4(keccak256(_fulfillmentSignature)));
    sendChainlinkRequestTo(_oracle, run, _wei);
  }

  function publicCancelRequest(
    bytes32 _requestId,
    uint256 _payment,
    bytes4 _callbackFunctionId,
    uint256 _expiration
  ) public {
    cancelChainlinkRequest(_requestId, _payment, _callbackFunctionId, _expiration);
  }

  function publicChainlinkToken() public view returns (address) {
    return chainlinkTokenAddress();
  }

  function publicFulfillChainlinkRequest(bytes32 _requestId, bytes32) public {
    fulfillRequest(_requestId, bytes32(0));
  }

  function fulfillRequest(bytes32 _requestId, bytes32)
    public
  {
    validateChainlinkCallback(_requestId);
  }

  function publicLINK(
    uint256 _amount
  )
    public
  {
    emit LinkAmount(LINK.mul(_amount));
  }

  function publicOracleAddress()
    public
    view
    returns (
      address
    )
  {
    return chainlinkOracleAddress();
  }

  function publicAddExternalRequest(
    address _oracle,
    bytes32 _requestId
  )
    public
  {
    addChainlinkExternalRequest(_oracle, _requestId);
  }
}
