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
    uint8[] _vs,
    bytes32[] _rs,
    bytes32[] _ss,
    bytes32 _requestDigest
  ) public
  {
    require(_oracles.length == _vs.length && _vs.length == _rs.length && _rs.length == _ss.length, "Must pass in as many signatures as oracles");

    bytes32 serviceAgreementID = getId(_payment, _expiration, _oracles, _requestDigest);

    for (uint i = 0; i < _oracles.length; i++) {
      address signer = getOracleAddressFromSASignature(serviceAgreementID, _vs[i], _rs[i], _ss[i]);
      require(_oracles[i] == signer, "Invalid oracle signature specified in SA");
    }

    serviceAgreements[serviceAgreementID] = ServiceAgreement(
      _payment,
      _expiration,
      _oracles,
      _requestDigest
    );
  }

  function getOracleAddressFromSASignature(
    bytes32 _serviceAgreementID,
    uint8 _v,
    bytes32 _r,
    bytes32 _s
  )
    private pure returns (address)
  {
    bytes32 prefixedHash = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", _serviceAgreementID));
    return ecrecover(prefixedHash, _v, _r, _s);
  }
}
