pragma solidity ^0.4.23;

import "./Chainlinked.sol";
import "./Ownable.sol";

contract Consumer is Chainlinked, Ownable {
  bytes32 internal requestId;
  bytes32 internal jobId;
  bytes32 public currentPrice;

  function Consumer(address _link, address _oracle, bytes32 _jobId) public {
    setLinkToken(_link);
    setOracle(_oracle);
    jobId = _jobId;
  }

  function requestEthereumPrice(string _currency) public {
    ChainlinkLib.Run memory run = newRun(jobId, this, "fulfill(bytes32,bytes32)");
    run.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    run.addStringArray("path", path);
    requestId = chainlinkRequest(run, 1 szabo);
  }

  function stringToBytes32(string memory source)
    internal
    returns (bytes32 result) {
      bytes memory tempEmptyStringTest = bytes(source);
      if (tempEmptyStringTest.length == 0) {
          return 0x0;
    }

    assembly {
        result := mload(add(source, 32))
    }
  }

  function cancelRequest(uint256 _requestId) 
    public 
    onlyOwner 
  {
    oracle.cancel(_requestId);
  }

  function fulfill(bytes32 _requestId, bytes32 _data)
    public
    onlyOracle
    checkRequestId(_requestId)
  {
    currentPrice = _data;
  }

  modifier checkRequestId(bytes32 _requestId) {
    require(requestId == _requestId);
    _;
  }

}