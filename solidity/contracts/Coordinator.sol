pragma solidity ^0.4.24;

// Coordinator handles oracle service aggreements between one or more oracles.
contract Coordinator {

  function getId(uint256 _payment, uint256 _expiration, bytes32 _jobSpecId)
    public pure returns (bytes32)
  {
    return keccak256(abi.encodePacked(_payment, _expiration, _jobSpecId));
  }

}
