pragma solidity 0.7.0;

import "./Operator.sol";

contract OperatorGenerator {

    address public link;

    event OperatorCreated(address indexed operator, address indexed owner);

    constructor(address linkAddress) public {
        link = linkAddress;
    }

    function createOperator() external returns (address operatorAddress){
        Operator operator = new Operator(link, msg.sender);
        operatorAddress = address(operator);
        emit OperatorCreated(operatorAddress, msg.sender);
    }
}