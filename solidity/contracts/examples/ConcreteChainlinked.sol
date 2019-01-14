pragma solidity 0.4.24;

import "../Chainlinked.sol";


contract ConcreteChainlinked is Chainlinked {

  constructor(address _link, address _oracle)
    public
  {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  event Run(
    bytes32 id,
    address callbackAddress,
    bytes4 callbackfunctionSelector,
    bytes data
  );

  function publicNewRun(
    bytes32 _id,
    address _address,
    bytes _fulfillmentSignature
  )
    public
  {
    ChainlinkLib.Run memory run = newRun(
      _id, _address, bytes4(keccak256(_fulfillmentSignature)));
    emit Run(
      run.id,
      run.callbackAddress,
      run.callbackFunctionId,
      run.buf.buf
    );
  }

  function publicRequestRun(
    bytes32 _id,
    address _address,
    bytes _fulfillmentSignature,
    uint256 _wei
  )
    public
  {
    ChainlinkLib.Run memory run = newRun(
      _id, _address, bytes4(keccak256(_fulfillmentSignature)));
    chainlinkRequest(run, _wei);
  }

  function publicRequestRunFrom(
    address _oracle,
    bytes32 _id,
    address _address,
    bytes _fulfillmentSignature,
    uint256 _wei
  )
    public
  {
    ChainlinkLib.Run memory run = newRun(_id, _address, bytes4(keccak256(_fulfillmentSignature)));
    chainlinkRequestFrom(_oracle, run, _wei);
  }

  function publicCancelRequest(bytes32 _requestId) public {
    cancelChainlinkRequest(_requestId);
  }

  function publicChainlinkToken() public view returns (address) {
    return chainlinkToken();
  }

  function fulfillRequest(bytes32 _requestId, bytes32)
    public
    checkChainlinkFulfillment(_requestId)
  {}

  function publicCompleteChainlinkFulfillment(bytes32 _requestId, bytes32) public {
    completeChainlinkFulfillment(_requestId);
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
