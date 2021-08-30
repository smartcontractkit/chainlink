// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "../ChainlinkClient.sol";

contract MultiWordConsumer is ChainlinkClient{
  bytes32 internal specId;
  bytes public currentPrice;

  bytes32 public first;
  bytes32 public second;

  event RequestFulfilled(
    bytes32 indexed requestId,  // User-defined ID
    bytes indexed price
  );

  event RequestMultipleFulfilled(
    bytes32 indexed requestId,
    bytes32 indexed first,
    bytes32 indexed second
  );

  constructor(address _link, address _oracle, bytes32 _specId) public {
    setChainlinkToken(_link);
    setChainlinkOracle(_oracle);
    specId = _specId;
  }

  function requestEthereumPrice(string memory _currency, uint256 _payment) public {
    requestEthereumPriceByCallback(_currency, _payment, address(this));
  }

  function requestEthereumPriceByCallback(string memory _currency, uint256 _payment, address _callback) public {
    Chainlink.Request memory req = buildChainlinkRequest(specId, _callback, this.fulfillBytes.selector);
    req.add("get", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    req.addStringArray("path", path);
    sendChainlinkRequest(req, _payment);
  }

  function requestMultipleParameters(string memory _currency, uint256 _payment) public {
    Chainlink.Request memory req = buildChainlinkRequest(specId, address(this), this.fulfillMultipleParameters.selector);
    req.add("get", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    req.addStringArray("path", path);
    sendChainlinkRequest(req, _payment);
  }

  function cancelRequest(
    address _oracle,
    bytes32 _requestId,
    uint256 _payment,
    bytes4 _callbackFunctionId,
    uint256 _expiration
  ) public {
    ChainlinkRequestInterface requested = ChainlinkRequestInterface(_oracle);
    requested.cancelOracleRequest(_requestId, _payment, _callbackFunctionId, _expiration);
  }

  function withdrawLink() public {
    LinkTokenInterface _link = LinkTokenInterface(chainlinkTokenAddress());
    require(_link.transfer(msg.sender, _link.balanceOf(address(this))), "Unable to transfer");
  }

  function addExternalRequest(address _oracle, bytes32 _requestId) external {
    addChainlinkExternalRequest(_oracle, _requestId);
  }

  function fulfillMultipleParameters(bytes32 _requestId, bytes32 _first, bytes32 _second)
    public
    recordChainlinkFulfillment(_requestId)
  {
    emit RequestMultipleFulfilled(_requestId, _first, _second);
    first = _first;
    second = _second;
  }

  function fulfillBytes(bytes32 _requestId, bytes memory _price)
    public
    recordChainlinkFulfillment(_requestId)
  {
    emit RequestFulfilled(_requestId, _price);
    currentPrice = _price;
  }
}