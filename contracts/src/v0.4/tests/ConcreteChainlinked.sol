pragma solidity 0.4.24;

import "../Chainlinked.sol";
import "../vendor/SafeMathChainlink.sol";

contract ConcreteChainlinked is Chainlinked {
  using SafeMathChainlink for uint256;

  constructor(address _link, address _oracle) public {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  event Request(
    bytes32 id,
    address callbackAddress,
    bytes4 callbackfunctionSelector,
    bytes data
  );

  function publicNewRequest(
    bytes32 _id,
    address _address,
    bytes _fulfillmentSignature
  )
    public
  {
    Chainlink.Request memory req = newRequest(
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
    bytes _fulfillmentSignature,
    uint256 _wei
  )
    public
  {
    Chainlink.Request memory req = newRequest(
      _id, _address, bytes4(keccak256(_fulfillmentSignature)));
    chainlinkRequest(req, _wei);
  }

  function publicRequestRunTo(
    address _oracle,
    bytes32 _id,
    address _address,
    bytes _fulfillmentSignature,
    uint256 _wei
  )
    public
  {
    Chainlink.Request memory run = newRequest(_id, _address, bytes4(keccak256(_fulfillmentSignature)));
    chainlinkRequestTo(_oracle, run, _wei);
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
    return chainlinkToken();
  }

  function fulfillRequest(bytes32 _requestId, bytes32)
    public
    recordChainlinkFulfillment(_requestId)
  {} // solhint-disable-line no-empty-blocks

  function publicFulfillChainlinkRequest(bytes32 _requestId, bytes32) public {
    fulfillChainlinkRequest(_requestId);
  }

  event LinkAmount(uint256 amount);

  function publicLINK(uint256 _amount) public {
    emit LinkAmount(LINK.mul(_amount));
  }

  function publicOracleAddress() public view returns (address) {
    return oracleAddress();
  }

  function publicAddExternalRequest(address _oracle, bytes32 _requestId)
    public
  {
    addExternalRequest(_oracle, _requestId);
  }
}
