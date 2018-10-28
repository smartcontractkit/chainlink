pragma solidity ^0.4.24;

import "../Chainlinked.sol";

contract Consumer is Chainlinked {
  bytes32 internal specId;
  bytes32 public currentPrice;

  constructor(address _link, address _oracle, bytes32 _specId) public {
    setLinkToken(_link);
    setOracle(_oracle);
    specId = _specId;
  }

  function requestEthereumPrice() public {
    chainlinkLib.Run memory run = newRun(specId, this, "reportPrice(bytes32,uint256)");
    chainlinkRequest(run, LINK(1));
  }

  function reportPrice(bytes32 _requestId, uint256 _price)
    public
    checkChainlinkFulfillment(_requestId)
  {
    //...
  }

}
