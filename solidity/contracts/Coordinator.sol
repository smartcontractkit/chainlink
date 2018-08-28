pragma solidity ^0.4.24;
pragma experimental ABIEncoderV2; //solium-disable-line

// Coordinator handles oracle service aggreements between one or more oracles.
contract Coordinator {
  struct ServiceAgreement {
    uint256 payment;
    uint256 expiration;
    address[] oracles;
    string requestDigest;
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
    string _msg
  ) public
  {
    //require(_oracles.length == _signatures.length);

    bytes32 id = getId(_payment, _expiration, _oracles, _msg);
    emit EmitString("ID");
    emit EmitID(id);

    for (uint i = 0; i < _oracles.length; i++) {
      emit EmitString("!!! SHOULD verify each participant");
      emit EmitString(_msg);

      //bytes[] signature = _signatures[i];

      uint8 v = _vs[i];
      bytes32 r = _rs[i];
      bytes32 s = _ss[i];

      emit EmitV(v);
      emit EmitR(r);
      emit EmitS(s);

      address signer = getOracleAddressFromSASignature(_msg, v, r, s);// signature);
      emit EmitAddress(signer);

      // memory said = _payment + _expiration + _oracles + keccack256(_normalizedJSON)
      // signature = sign(said)

      address oracle = _oracles[i];
      emit EmitAddress(oracle);
      // require(
      //   oracle == signer,
      //   "!!! oracle is not the signer: TODO: can it do string interpolation of the addresses???"
      // );
    }

    // bytes32 id = getId(_payment, _expiration, _oracles, _msg);

    serviceAgreements[id] = ServiceAgreement(
      _payment,
      _expiration,
      _oracles,
      _msg
    );
  }

  //function getOracleAddressFromSASignature(bytes32 _hash, bytes32 _sig) returns (address) {
  //function getOracleAddressFromSASignature(bytes32 _hash, bytes[] _sig) returns (address) {
  function getOracleAddressFromSASignature(string _msg, uint8 _v, bytes32 _r, bytes32 _s) private pure returns (address) {
    //bytes32 r;
    //bytes32 s;
    //uint8 v;

    //if (sig.length != 65) {
      //return 0;
    //}

    //assembly {
      //v := byte(0, sig)
      //r := and(sig, 0x00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff)
      ////s := mload(add(sig, 33))
    //}

    // https://github.com/ethereum/go-ethereum/issues/2053
    // if (_v < 27) {
      // _v += 27;
    // }

    //if (v != 27 && v != 28) {
      //return 0;
    //}

    ///* prefix might be needed for geth only
     //* https://github.com/ethereum/go-ethereum/issues/3731
     //*/
    bytes memory prefix = "\x19Ethereum Signed Message:\n11";
    bytes32 prefixedHash = keccak256(abi.encodePacked(prefix, _msg));

    return ecrecover(prefixedHash, _v, _r, _s);
  }
}
