pragma solidity 0.4.24;

import "../Chainlinked.sol";

contract Consumer is Chainlinked {
  bytes32 internal specId;
  bytes32 public currentPrice;

  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations

  event RequestFulfilled(
    bytes32 indexed requestId,  // User-defined ID
    bytes32 indexed price
  );

  function requestEthereumPrice(string _currency) public {
    ChainlinkLib.Run memory run = newRun(specId, this, this.fulfill.selector);
    run.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    run.addStringArray("path", path);
    chainlinkRequest(run, ORACLE_PAYMENT);
  }

  function cancelRequest(bytes32 _requestId) public {
    cancelChainlinkRequest(_requestId);
  }

  function withdrawLink() public {
    LinkTokenInterface link = LinkTokenInterface(chainlinkToken());
    require(link.transfer(msg.sender, link.balanceOf(address(this))), "Unable to transfer");
  }

  function fulfill(bytes32 _requestId, bytes32 _price)
    public
    checkChainlinkFulfillment(_requestId)
  {
    emit RequestFulfilled(_requestId, _price);
    currentPrice = _price;
  }

}
