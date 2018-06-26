pragma solidity ^0.4.24;

contract GetMoney {
  address[] public payees;
  event LogMoney(uint256 indexed amount);

  function receive() public payable {
    payees.push(msg.sender);
    emit LogMoney(msg.value);
  }
}
