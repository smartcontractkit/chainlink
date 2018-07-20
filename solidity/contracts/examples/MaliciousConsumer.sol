pragma solidity ^0.4.24;


import "../Chainlinked.sol";


contract MaliciousConsumer is Chainlinked {

  constructor(address _link, address _oracle)
    public
  {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  function requestData(string _callbackFunc)
    public
  {
    ChainlinkLib.Run memory run = newRun("specId", this, _callbackFunc);
    chainlinkRequest(run, LINK(1));
  }

  function assertFail(bytes32 _requestId, bytes32 _data)
    public
  {
    assert(1 == 2);
  }

  function cancelRequestOnFulfill(bytes32 _requestId, bytes32 _data)
    public
  {
    cancelChainlinkRequest(_requestId);
  }

  function remove() public {
    selfdestruct(address(0));
  }

  function doesNothing(bytes32 _requestId, bytes32 _data) public {}
}
