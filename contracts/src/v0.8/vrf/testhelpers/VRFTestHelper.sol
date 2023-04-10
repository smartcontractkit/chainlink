// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../VRF.sol";

/** ***********************************************************************
    @notice Testing harness for VRF.sol, exposing its internal methods. Not to
    @notice be used for production.
*/
contract VRFTestHelper is VRF {
  function bigModExp_(uint256 base, uint256 exponent) public view returns (uint256) {
    return super.bigModExp(base, exponent);
  }

  function squareRoot_(uint256 x) public view returns (uint256) {
    return super.squareRoot(x);
  }

  function ySquared_(uint256 x) public pure returns (uint256) {
    return super.ySquared(x);
  }

  function fieldHash_(bytes memory b) public pure returns (uint256) {
    return super.fieldHash(b);
  }

  function hashToCurve_(uint256[2] memory pk, uint256 x) public view returns (uint256[2] memory) {
    return super.hashToCurve(pk, x);
  }

  function ecmulVerify_(
    uint256[2] memory x,
    uint256 scalar,
    uint256[2] memory q
  ) public pure returns (bool) {
    return super.ecmulVerify(x, scalar, q);
  }

  function projectiveECAdd_(
    uint256 px,
    uint256 py,
    uint256 qx,
    uint256 qy
  )
    public
    pure
    returns (
      uint256,
      uint256,
      uint256
    )
  {
    return super.projectiveECAdd(px, py, qx, qy);
  }

  function affineECAdd_(
    uint256[2] memory p1,
    uint256[2] memory p2,
    uint256 invZ
  ) public pure returns (uint256[2] memory) {
    return super.affineECAdd(p1, p2, invZ);
  }

  function verifyLinearCombinationWithGenerator_(
    uint256 c,
    uint256[2] memory p,
    uint256 s,
    address lcWitness
  ) public pure returns (bool) {
    return super.verifyLinearCombinationWithGenerator(c, p, s, lcWitness);
  }

  function linearCombination_(
    uint256 c,
    uint256[2] memory p1,
    uint256[2] memory cp1Witness,
    uint256 s,
    uint256[2] memory p2,
    uint256[2] memory sp2Witness,
    uint256 zInv
  ) public pure returns (uint256[2] memory) {
    return super.linearCombination(c, p1, cp1Witness, s, p2, sp2Witness, zInv);
  }

  function scalarFromCurvePoints_(
    uint256[2] memory hash,
    uint256[2] memory pk,
    uint256[2] memory gamma,
    address uWitness,
    uint256[2] memory v
  ) public pure returns (uint256) {
    return super.scalarFromCurvePoints(hash, pk, gamma, uWitness, v);
  }

  function isOnCurve_(uint256[2] memory p) public pure returns (bool) {
    return super.isOnCurve(p);
  }

  function verifyVRFProof_(
    uint256[2] memory pk,
    uint256[2] memory gamma,
    uint256 c,
    uint256 s,
    uint256 seed,
    address uWitness,
    uint256[2] memory cGammaWitness,
    uint256[2] memory sHashWitness,
    uint256 zInv
  ) public view {
    super.verifyVRFProof(pk, gamma, c, s, seed, uWitness, cGammaWitness, sHashWitness, zInv);
  }

  function randomValueFromVRFProof_(Proof memory proof, uint256 seed) public view returns (uint256 output) {
    return super.randomValueFromVRFProof(proof, seed);
  }
}
