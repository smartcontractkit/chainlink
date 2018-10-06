pragma solidity ^0.4.24;

import "../Chainlinked.sol";

contract PaymentForDelivery is Chainlinked {
  bytes32 internal specId;

  constructor(address _link, address _oracle, bytes32 _specId) public {
    setLinkToken(_link);
    setOracle(_oracle);
    specId = _specId;
  }

  function trackShipment(
    uint256 _shipmentId
  ) public {
    ChainlinkLib.Run memory run = newRun(specId, this, "arrivalPrice(bytes32,uint256)");
    run.addString("carrier", "FedEx");
    run.addUint("shipmentId", _shipmentId);
    string[] memory sources = new string[](3);
    sources[0] = "BraveNewCoin";
    sources[1] = "CryptoCompare";
    sources[2] = "CoinMarketCap";
    run.addStringArray("sources", sources);
    chainlinkRequest(run, LINK(5));
  }

  function arrivalPrice(bytes32 _requestId, uint256 _price)
    public
    checkChainlinkFulfillment(_requestId)
  {
    //...
  }

}
