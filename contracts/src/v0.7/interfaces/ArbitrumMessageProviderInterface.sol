pragma solidity ^0.7.0;

interface ArbitrumMessageProviderInterface {
    event InboxMessageDelivered(uint256 indexed messageNum, bytes data);

    event InboxMessageDeliveredFromOrigin(uint256 indexed messageNum);
}
