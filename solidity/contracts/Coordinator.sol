pragma solidity ^0.4.24;
pragma experimental ABIEncoderV2; //solium-disable-line

// Coordinator handles oracle service aggreements between one or more oracles.
contract Coordinator {
  struct ServiceAgreement {
    uint256 payment;
    uint256 expiration;
    address[] oracles;
    string strRequestDigest;
    bytes32 requestDigest;
  }

  mapping(bytes32 => ServiceAgreement) public serviceAgreements;

  event EmitString(string msg);
  event EmitAddress(address msg);
  event EmitV(uint8 v);
  event EmitR(bytes32 r);
  event EmitS(bytes32 s);
  event EmitID(bytes32 id);

  function getId(
    uint256 _payment,
    uint256 _expiration,
    address[] _oracles,
    string _msg
  )
    public pure returns (bytes32)
  {
    return keccak256(abi.encodePacked(_payment, _expiration, _oracles, _msg));
  }

  // XXX: No nested structs in web3
  //struct Signature {
    //uint8 v;
    //bytes32 r;
    //bytes32 s;
  //}

  function initiateServiceAgreement(
    uint256 _payment,
    uint256 _expiration,
    address[] _oracles,
    uint8[] _vs,
    bytes32[] _rs,
    bytes32[] _ss,
    // XXX: no nested structs in web3
    // bytes[][] _signatures,
    // Signature[] _signatures,
    string _msg,
    bytes32 _requestDigest
  ) public
  {
    bytes32 id = getId(_payment, _expiration, _oracles, _msg);

    for (uint i = 0; i < _oracles.length; i++) {
      emit EmitString("Oracle at index");
      emit EmitAddress(_oracles[i]);
      emit EmitString("Signers...");
      address strSigner = getOracleAddressFromSASignatureStr(_msg, _vs[i], _rs[i], _ss[i]);
      emit EmitAddress(strSigner);
      address signer = getOracleAddressFromSASignature(_requestDigest, _vs[i], _rs[i], _ss[i]);
      emit EmitAddress(signer);

      // require(
      //   // _oracles[i] == strSigner,
      //   _oracles[i] == signer,
      //   "!!! oracle is not the signer: TODO: can it do string interpolation of the addresses???"
      // );
    }

    serviceAgreements[id] = ServiceAgreement(
      _payment,
      _expiration,
      _oracles,
      _msg,
      _requestDigest
    );
  }

  function getOracleAddressFromSASignature(bytes32 _requestDigest, uint8 _v, bytes32 _r, bytes32 _s) private pure returns (address) {
    bytes memory prefix = "\x19Ethereum Signed Message:\n11";
    bytes32 prefixedHash = keccak256(abi.encodePacked(prefix, _requestDigest));

    return ecrecover(prefixedHash, _v, _r, _s);
  }

  function getOracleAddressFromSASignatureStr(string _msg, uint8 _v, bytes32 _r, bytes32 _s) private pure returns (address) {
    bytes memory prefix = "\x19Ethereum Signed Message:\n11";
    bytes32 prefixedHash = keccak256(abi.encodePacked(prefix, _msg));

    return ecrecover(prefixedHash, _v, _r, _s);
  }
}
