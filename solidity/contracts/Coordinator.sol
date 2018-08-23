pragma solidity ^0.4.24;
pragma experimental ABIEncoderV2; //solium-disable-line

// Coordinator handles oracle service aggreements between one or more oracles.
contract Coordinator {

  struct ServiceAgreement {
    uint256 payment;
    uint256 expiration;
    address[] oracles;
    bytes32 requestDigest;
  }

  mapping(bytes32 => ServiceAgreement) public serviceAgreements;

  function getId(
    uint256 _payment,
    uint256 _expiration,
    address[] _oracles,
    bytes32 _requestDigest
  )
    public pure returns (bytes32)
  {
    return keccak256(abi.encodePacked(_payment, _expiration, _oracles, _requestDigest));
  }

  function initiateServiceAgreement(
    uint256 _payment,
    uint256 _expiration,
    address[] _oracles,
    bytes32 _requestDigest
  ) public
  {
    bytes32 id = getId(_payment, _expiration, _oracles, _requestDigest);

    serviceAgreements[id] = ServiceAgreement(
      _payment,
      _expiration,
      _oracles,
      _requestDigest
    );
  }
}
