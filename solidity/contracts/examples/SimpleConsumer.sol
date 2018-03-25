pragma solidity ^0.4.18;

import "../../contracts/Chainlinked.sol";

contract SimpleConsumer is Chainlinked {
  uint256 private requestId;
  bytes32 public currentPrice;

  function SimpleConsumer(address _oracle) public {
    setOracle(_oracle);
  }

  function requestEthereumPrice() public {
    var functionId = bytes4(keccak256("fulfill(uint256,bytes32)"));
    var data = '{"url":"https://etherprice.com/api","path":["recent","usd"]}';
    requestId = oracle.requestData("someJobId", this, functionId, data);
  }

  function fulfill(uint256 _requestId, bytes32 _data)
    public
    onlyOracle
    checkRequestId(_requestId)
  {
    currentPrice = _data;
  }

  modifier checkRequestId(uint256 _requestId) {
    require(requestId == _requestId);
    _;
  }

}
