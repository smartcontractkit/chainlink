pragma solidity ^0.4.23;

import "../../../solidity/contracts/Chainlinked.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract SpecAndRunLog is Chainlinked, Ownable {
  bytes32 internal requestId;
  bytes32 public currentPrice;

  event RequestFulfilled(
    bytes32 indexed requestId,
    bytes32 indexed price
  );

  constructor(address _link, address _oracle) Ownable() public {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  function request() public {
    string[] memory tasks = new string[](1);
    tasks[0] = "httppost";

    ChainlinkLib.Spec memory spec = newSpec(tasks, this, "fulfill(bytes32,bytes32)");
    spec.add("msg", "hello_chainlink");
    spec.add("url", "http://localhost:6690");
    requestId = chainlinkRequest(spec, LINK(1));
  }

  function cancelRequest() public onlyOwner {
    oracle.cancel(requestId);
  }

  function fulfill(bytes32 _requestId, bytes32 _price)
    public
    checkChainlinkFulfillment(_requestId)
  {
  }
}
