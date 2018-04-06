pragma solidity ^0.4.18;

import "../../../solidity/contracts/Chainlinked.sol";

contract RunLog is Chainlinked {
  uint256 private requestId;

  function RunLog(address _oracle) public {
    setOracle(_oracle);
  }

  function request() public {
    var fid = bytes4(keccak256("fulfill(uint256,bytes32)"));
    var data = '{"msg":"hello_chainlink"}';
    requestId = oracle.requestData("MY_JOB_ID", this, fid, data);
  }

  function fulfill(uint256 _requestId, bytes32 _data)
    public
    onlyOracle
    checkRequestId(_requestId)
  {
  }

  modifier checkRequestId(uint256 _requestId) {
    require(requestId == _requestId);
    _;
  }
}
