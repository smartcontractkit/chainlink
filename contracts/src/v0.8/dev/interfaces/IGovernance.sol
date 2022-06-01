pragma solidity 0.8.13;
interface IGovernance {

    //PROPOSAL STATE FLOW
    //Pending ==> Active (time-based according to votingDelay) | state driven, out of scope
    //Active ==> Defeated (if votingPeriod passed AND forVotes < quorumVotes) | state driven, out of scope
    //Active ==> Succeeded (if forVotes > quorumVotes AND forVotes > againstVotes) | state driven, out of scope
    //Succeded ==> Queued (if state is Succeeded) | **action driven, in scope**
    //Queued ==> Executed (if state is Queued AND ETA time passed) | **action driven, in scope**
    //!Executed ==> Canceled (if proposer votes < proposal threshold) | **action driven, in scope**
    
    /// @notice Possible states that a proposal may be in
    enum ProposalState {Pending, Canceled, Active, Failed, Succeeded, Queued, Expired, Executed}
    function state(uint256 proposalId) external view returns (ProposalState);
    function proposalCount() external view returns (uint256);
    function queue(uint proposalId) external;
    function cancel(uint proposalId) external;
    function execute(uint proposalId) external;
    struct Proposal {
        uint id;
        address proposer;
        uint eta;
        address[] targets;
        uint[] values;
        string[] signatures;
        bytes[] calldatas;
        uint startBlock;
        uint endBlock;
        uint forVotes;
        uint againstVotes;
        bool canceled;
        bool executed;
    }
    function proposals(uint proposalId) external view returns (Proposal memory proposal);
    function proposalThreshold() external pure returns (uint);
}