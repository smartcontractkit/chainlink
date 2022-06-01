
pragma solidity 0.8.13;
interface IGovernanceToken {
    function getPriorVotes(address account, uint blockNumber) external view returns (uint96);
}