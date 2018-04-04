pragma solidity ^0.4.11;


import '../token/ERC20.sol';


contract LinkReceiver {

  bool public fallbackCalled;
  bool public callDataCalled;
  uint public tokensReceived;


  function onTokenTransfer(address _from, uint _amount, bytes _data)
  public returns (bool success) {
    fallbackCalled = true;
    if (_data.length > 0) {
      require(address(this).delegatecall(_data, msg.sender, _from, _amount));
    }
    return true;
  }

  function callbackWithoutWithdrawl() {
    callDataCalled = true;
  }

  function callbackWithWithdrawl(uint _value, address _from, address _token) {
    callDataCalled = true;
    ERC20 token = ERC20(_token);
    token.transferFrom(_from, this, _value);
    tokensReceived = _value;
  }

}
