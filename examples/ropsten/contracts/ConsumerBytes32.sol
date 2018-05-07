pragma solidity ^0.4.23;

import "./Chainlinked.sol";
import "github.com/OpenZeppelin/openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract ConsumerBytes32 is Chainlinked, Ownable {
  bytes32 internal requestId;
  bytes32 internal jobId;
  bytes32 public currentMarketCap;

  event RequestFulfilled(
    bytes32 indexed requestId,
    bytes32 indexed marketCap
  );

  function ConsumerBytes32(address _link, address _oracle, bytes32 _jobId) Ownable() public {
    setLinkToken(_link);
    setOracle(_oracle);
    jobId = _jobId;
  }

  function requestEthereumMarketCap(string _currency)
    public
    onlyOwner
  {
    ChainlinkLib.Run memory run = newRun(jobId, this, "fulfill(bytes32,bytes32)");
    run.add("url", "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=ETH&tsyms=USD");
    string[] memory path = new string[](4);
    path[0] = "DISPLAY";
    path[1] = "ETH";
    path[2] = _currency;
    path[3] = "MKTCAP";
    run.addStringArray("path", path);
    requestId = chainlinkRequest(run, LINK(1));
  }

  function cancelRequest(uint256 _requestId) 
    public 
    onlyOwner
  {
    oracle.cancel(_requestId);
  }

  function fulfill(bytes32 _requestId, bytes32 _marketCap)
    public
    onlyOracle
    checkRequestId(_requestId)
  {
    emit RequestFulfilled(_requestId, _marketCap);
    currentMarketCap = _marketCap;
  }

  modifier checkRequestId(bytes32 _requestId) {
    require(requestId == _requestId);
    _;
  }

}
