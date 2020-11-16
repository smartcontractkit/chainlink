pragma solidity ^0.7.0;

import "./Operator.sol";

/**
 * @title Operator Generator
 * @notice Generates Operator contracts for node operators
 */
contract OperatorGenerator {

    address public link;

    event OperatorCreated(address indexed operator, address indexed owner);

    /**
     * @param linkAddress address
     */
    constructor(address linkAddress) public {
        link = linkAddress;
    }

    /**
     * @notice Create a new Operator contract with the msg.sender as owner
     * @return operatorAddress
     */
    function createOperator() external returns (address operatorAddress){
        Operator operator = new Operator(link, msg.sender);
        operatorAddress = address(operator);
        emit OperatorCreated(operatorAddress, msg.sender);
    }
}