pragma solidity ^0.4.23;

contract GetMoney {
  address[] public payees;
  event LogMoney(uint256 indexed amount);

  function receive() public payable {
    payees.push(msg.sender);
    LogMoney(msg.value);
  }
}
