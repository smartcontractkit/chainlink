pragma solidity ^0.8.0;

abstract contract VRFConsumerBaseV2 {
    function fulfillRandomWords(
        uint256 requestId,
        uint256[] memory randomWords
    )
    public
    virtual;
}
