pragma solidity ^0.4.18;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";

contract ChainLink is Ownable {

  struct Callback {
    address addr;
    bytes4 fid;
  }

  uint256 private nonce;
  mapping(uint256 => Callback) private callbacks;

  event Request(
    uint256 indexed nonce,
    address indexed to,
    bytes4 indexed fid,
    string data
  );

  function requestData(
    address _callbackAddress,
    bytes4 _callbackFID,
    string _data
  )
    public
    returns (uint256)
  {
    nonce += 1;
    Callback memory cb = Callback(_callbackAddress, _callbackFID);
    callbacks[nonce] = cb;
    Request(nonce, cb.addr, cb.fid, _data);
    return nonce;
  }

  function fulfillData(uint256 _nonce, bytes32 _data)
    public
    onlyOwner
    hasNonce(_nonce)
  {
    Callback memory cb = callbacks[_nonce];
    require(cb.addr.call(cb.fid, _nonce, _data));
    delete callbacks[_nonce];
  }

  modifier hasNonce(uint256 _nonce) {
    require(callbacks[_nonce].addr != address(0));
    _;
  }
}
