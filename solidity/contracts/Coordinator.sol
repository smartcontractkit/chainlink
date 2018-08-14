pragma solidity ^0.4.24;

// Coordinator handles oracle service aggreements between one or more oracles.
contract Coordinator {

  struct ServiceAgreement {
    uint256 payment;
    uint256 expiration;
    bytes32 jobSpecID;
  }

  mapping(bytes32 => ServiceAgreement) internal serviceAgreements;

  function getId(uint256 _payment, uint256 _expiration, bytes32 _jobSpecId)
    public pure returns (bytes32)
  {
    return keccak256(abi.encodePacked(_payment, _expiration, _jobSpecId));
  }

  function initiateServiceAgreement(uint256 _payment, uint256 _expiration, bytes32 _jobSpecId)
    public
  {
    bytes32 id = getId(_payment, _expiration, _jobSpecId);

    serviceAgreements[id] = ServiceAgreement(
      _payment,
      _expiration,
      _jobSpecId
    );
  }

  function getServiceAgreement(bytes32 _id)
    public returns (bytes memory)
  {
    ServiceAgreement sa = serviceAgreements[_id];
    return abi.encode(sa.payment, sa.expiration, sa.jobSpecID);
  }
}
