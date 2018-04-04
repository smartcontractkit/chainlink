pragma solidity ^0.4.11;


import "../ERC677Token.sol";
import "../token/StandardToken.sol";


contract Token677 is StandardToken, ERC677Token {
    string public constant name = "Example ERC677 Token";
    string public constant symbol = "ERC677";
    uint8 public constant decimals = 18;
    uint256 public totalSupply;

    function Token677(uint _initialBalance)
    {
        balances[msg.sender] = _initialBalance;
        totalSupply = _initialBalance;
    }
}
