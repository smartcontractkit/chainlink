pragma solidity ^0.4.21;

import "../../../solidity/contracts/Chainlinked.sol";

contract RunLog is Chainlinked {
  uint256 private requestId;

  function RunLog(address _oracle) public {
    setOracle(_oracle);
  }

  function request() public {
    Chainlink.Run memory run = newRun("MY_JOB_ID", this, "fulfill(uint256,bytes32)");
    run.add("msg", "hello_chainlink");
    requestId = chainlinkRequest(run);
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
