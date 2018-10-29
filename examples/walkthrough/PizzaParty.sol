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
    string[] memory sources = new string[](3);
    path[0] = "BraveNewCoin";
    path[1] = "CryptoCompare";
    path[2] = "CoinMarketCap";
    run.addStringArray("sources", sources);
    run.addString("carrier", "FedEx");
    run.addUint("shipmentId", _shipmentId);
    chainlinkRequest(run, LINK(4));
  }

  function arrivalPrice(bytes32 _requestId, uint256 _price)
    public
    checkChainlinkFulfillment(_requestId)
  {
    //...
  }

}
