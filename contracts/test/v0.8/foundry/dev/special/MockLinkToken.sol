// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MockLinkToken {
  uint public constant totalSupply = 10**27;

  mapping(address => uint256) public balances;

  constructor() {
    balances[msg.sender] = totalSupply;
  }

  /**
  * @dev transfer token for a specified address
  * @param _to The address to transfer to.
  * @param _value The amount to be transferred.
  */
  function transfer(address _to, uint256 _value) public returns (bool) {
    balances[msg.sender] = balances[msg.sender] - _value;
    balances[_to] = balances[_to] + _value;
    return true;
  }

  function transferAndCall(address _to, uint _value, bytes calldata _data) public validRecipient(_to) returns (bool success) {
    transfer(_to, _value);
    if (isContract(_to)) {
      contractFallback(_to, _value, _data);
    }
    return true;
  }

  function balanceOf(address _a) public view returns (uint256 balance) {
    return balances[_a];
  }

  modifier validRecipient(address _recipient) {
    require(_recipient != address(0) && _recipient != address(this));
    _;
  }

  function contractFallback(address _to, uint _value, bytes calldata _data)
    private
  {
    ERC677Receiver receiver = ERC677Receiver(_to);
    receiver.onTokenTransfer(msg.sender, _value, _data);
  }

  function isContract(address _addr)
    private
    returns (bool hasCode)
  {
    uint length;
    assembly { length := extcodesize(_addr) }
    return length > 0;
  }
}

interface ERC677Receiver {
  function onTokenTransfer(address _sender, uint _value, bytes calldata _data) external;
}
