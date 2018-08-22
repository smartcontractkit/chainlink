pragma solidity ^0.4.24;
pragma experimental ABIEncoderV2; //solium-disable-line

// Coordinator handles oracle service aggreements between one or more oracles.
contract Coordinator {

  struct ServiceAgreement {
    uint256 payment;
    uint256 expiration;
    bytes32 jobSpecID;
  }

  mapping(bytes32 => ServiceAgreement) public serviceAgreements;

  function getId(uint256 _payment, uint256 _expiration, bytes32 _jobSpecID)
    public pure returns (bytes32)
  {
    return keccak256(abi.encodePacked(_payment, _expiration, _jobSpecID));
  }

  function initiateServiceAgreement(uint256 _payment, uint256 _expiration, bytes32 _jobSpecID)
    public
  {
    bytes32 id = getId(_payment, _expiration, _jobSpecID);

    serviceAgreements[id] = ServiceAgreement(
      _payment,
      _expiration,
      _jobSpecID
    );
  }
}
