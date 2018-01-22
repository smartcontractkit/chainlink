pragma solidity ^0.4.18;

import "solidity-stringutils/strings.sol";
import "../../contracts/ChainLink.sol";

contract Consumer {
  using strings for *;

  ChainLink private chainLink;
  uint256 private nonce;
  bytes32 public currentPrice;

  function Consumer(address _chainLink) public {
    chainLink = ChainLink(_chainLink);
  }

  function requestEthereumPrice() public {
    bytes4 fid = bytes4(keccak256("fulfill(uint256,bytes32)"));
    string memory payload = "{";
    payload = addParameter(payload, "url", "https://etherprice.com/api");
    payload = addParameter(payload, "path", "recent,usd");
    nonce = chainLink.requestData(this, fid, closeJSON(payload));
  }

  function addParameter(
    string data,
    string key,
    string value
  ) returns(string) {
    data = data.toSlice().concat(key.toSlice());
    data = data.toSlice().concat(":\"".toSlice());
    data = data.toSlice().concat(value.toSlice());
    return data.toSlice().concat("\",".toSlice());
  }

  function closeJSON(string data) returns (string) {
    var slice = data.toSlice();
    slice._len -= 1;
    return slice.concat("}".toSlice());
  }

  function fulfill(uint256 _nonce, bytes32 _data)
    public
    checkSender
    checkNonce(_nonce)
  {
    currentPrice = _data;
  }

  modifier checkSender() {
    require(msg.sender == address(chainLink));
    _;
  }

  modifier checkNonce(uint256 _nonce) {
    require(nonce == _nonce);
    _;
  }
}
