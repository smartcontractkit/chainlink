pragma solidity ^0.4.24;
pragma experimental ABIEncoderV2; //solium-disable-line


import "../Chainlinked.sol";


contract ConcreteChainlinked is Chainlinked {

  constructor(address _link, address _oracle)
    public
  {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  event Run(
    bytes32 specId,
    address callbackAddress,
    bytes4 callbackfunctionSelector,
    bytes data
  );

  function publicNewRun(
    bytes32 _specId,
    address _address,
    string _fulfillmentSignature
  )
    public
  {
    ChainlinkLib.Run memory run = newRun(_specId, _address, _fulfillmentSignature);
    run.close();
    emit Run(
      run.specId,
      run.callbackAddress,
      run.callbackFunctionId,
      run.buf.buf
    );
  }

  function publicRequestRun(
    bytes32 _specId,
    address _address,
    string _fulfillmentSignature,
    uint256 _wei
  )
    public
  {
    ChainlinkLib.Run memory run = newRun(_specId, _address, _fulfillmentSignature);
    chainlinkRequest(run, _wei);
  }

  function publicCancelRequest(bytes32 _requestId) public {
    cancelChainlinkRequest(_requestId);
  }

  function publicChainlinkToken() public returns (address) {
    return chainlinkToken();
  }

  function fulfillRequest(bytes32 _requestId, bytes32 _)
    public
    checkChainlinkFulfillment(_requestId)
  {
  }

  event LinkAmount(uint256 amount);

  function publicLINK(uint256 _link) public {
    emit LinkAmount(LINK(_link));
  }
}
