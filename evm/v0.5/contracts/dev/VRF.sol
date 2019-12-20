pragma solidity 0.5.0;

////////////////////////////////////////////////////////////////////////////////
//       XXX: Do not use in production until this code has been audited.
////////////////////////////////////////////////////////////////////////////////

/** ****************************************************************************
  * @notice On-chain verification of verifiable-random-function (VRF) proofs as
  * @notice described in
  * @notice https://tools.ietf.org/html/draft-goldbe-vrf-01#section-5.3 and
  * @notice https://eprint.iacr.org/2017/099.pdf (security proofs)

  * @dev Bibliographic references:

  * @dev Goldberg, et al., "Verifiable Random Functions (VRFs)", Internet Draft
  * @dev draft-irtf-cfrg-vrf-05, IETF, Aug 11 2019,
  * @dev https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-vrf-05

  * @dev Papadopoulos, et al., "Making NSEC5 Practical for DNSSEC", Cryptology
  * @dev ePrint Archive, Report 2017/099, 2017
  * ****************************************************************************
  * @dev USAGE

  * @dev The main entry point is randomValueFromVRFProof. See its docstring.
  * ****************************************************************************
  * @dev PURPOSE

  * @dev Reggie the Random Oracle (not his real job) wants to provide randomness
  * @dev to Vera the verifier in such a way that Vera can be sure he's not
  * @dev making his output up to suit himself. Reggie provides Vera a public key
  * @dev to which he knows the secret key. Each time Vera provides a seed to
  * @dev Reggie, he gives back a value which is computed completely
  * @dev deterministically from the seed and the secret key.

  * @dev Reggie provides a proof by which Vera can verify that the output was
  * @dev correctly computed once Reggie tells it to her, but without that proof,
  * @dev the output is indistinguishable to her from a uniform random sample
  * @dev from the output space.

  * @dev The purpose of this contract is to perform that verification.
  * ****************************************************************************
  * @dev DESIGN NOTES

  * @dev The VRF algorithm verified here satisfies the full unqiqueness, full
  * @dev collision resistance, and full pseudorandomness security properties.
  * @dev See "SECURITY PROPERTIES" below, and
  * @dev https://tools.ietf.org/html/draft-goldbe-vrf-01#section-3

  * @dev An elliptic curve point is generally represented in the solidity code
  * @dev as a uint256[2], corresponding to its affine coordinates in
  * @dev GF(FIELD_SIZE).

  * @dev For the sake of efficiency, this implementation deviates from the spec
  * @dev in some minor ways:

  * @dev - Keccak hash rather than the SHA256 hash recommended in
  * @dev   https://tools.ietf.org/html/draft-goldbe-vrf-01#section-5.5 . This is
  * @dev   because keccak costs much less gas on the EVM. The impact onsecurity
  * @dev   should be minor.

  * @dev - Secp256k1 curve instead of the P-256 or ED25519 curves recommended in
  * @dev   https://tools.ietf.org/html/draft-goldbe-vrf-01#section-5.5 . This is
  * @dev   because it's much cheaper to abuse ECRECOVER for the most expensive
  * @dev   ECC arithmetic, when computing in the EVM.

  * @dev - hashToCurve recursively hashes until it finds a curve
  * @dev   x-ordinate. On the EVM, this is slightly more efficient than the
  * @dev   recommendation in
  * @dev   https://tools.ietf.org/html/draft-goldbe-vrf-01#section-5.4.1.1 step
  * @dev   4 to concatenate with a nonce then hash, and rehash with the nonce
  * @dev   updated until a valid x-ordinate is found.

  * @dev - In the calculation of the challenge value "c", the "u" value
  * @dev   (i.e. the value computed by Reggie as the nonce times the secp256k1
  * @dev   generator point, see steps 4 and 7 of
  * @dev   https://tools.ietf.org/html/draft-goldbe-vrf-01#section-5.3) is
  * @dev   replaced by its ethereum address of the point, which is the lower 160
  * @dev   bits of the keccak hash of the original u. This is because we only
  * @dev   verify the calculation of u up to its address, by abusing ECRECOVER.
  * ****************************************************************************
  * @dev SECURITY PROPERTIES

  * @dev Here are the security properties for this VRF:

  * @dev Full uniqueness: For any seed and valid VRF public key, there is
  * @dev   exactly one VRF output which can be proved to come from that seed, in
  * @dev   the sense that the proof will pass verifyVRFProof.

  * @dev Full collision resistance: It's cryptographically infeasible to find
  * @dev   two seeds with same VRF output from a fixed, valid VRF key

  * @dev Full pseudorandomness: Absent the proofs that the VRF outputs are
  * @dev   derived from a given seed, the outputs are computationally
  * @dev   indistinguishable from randomness.

  * @dev https://eprint.iacr.org/2017/099.pdf, Appendix B contains the proofs
  * @dev for these properties. The introduction to
  * @dev https://tools.ietf.org/html/draft-goldbe-vrf-01#section-5 is very
  * @dev conservative about the security properties it claims, but is implicitly
  * @dev stronger in its claims for the specific cipher suites described in
  * @dev section 5.5. The reason for this is given in the "NOTE" at the bottom
  * @dev of section 5.5, namely, to quote Appendix B:

  * @dev    If the group E is fixed and trusted to have been correctly
  * @dev    generated (i.e., E is known to have a subgroup of prime order q),
  * @dev    and the generator g is known to be in G − {1}, then the verifier
  * @dev    just needs to check that [VRF public key] PK ∈ E. (This is the only
  * @dev    requirement on PK in the proof [of trusted uniqueness] above.)

  * @dev A similar note is on the proof for trusted collision-resistance:

  * @dev     **Collision resistance without trusting the key**. Similarly
  * @dev     to the case with uniqueness, our VRF can be modified the same way
  * @dev     to attain collision resistance without needing to trust the key
  * @dev     generation. The modifications are the same as in the case of
  * @dev     uniqueness (to ensure that F_{SK} is uniquely defined), with the
  * @dev     additional check that PK^f≠1 to ensure that x is not divisible by q

  * @dev (For secp256k1, f, the cofactor of the subgroup, is 1)

  * @dev Thus, here we rely on the fact that the secp256k1 parameters are
  * @dev correct, and we can check that the VRF public key lies on secp256k1 and
  * @dev is not the generator or the zero point, so we do not have to trust in
  * @dev correct key generation.
*/
contract VRF {

  // See https://www.secg.org/sec2-v2.pdf, section 2.4.1, for these constants.
  uint256 constant private GROUP_ORDER = // Number of points in Secp256k1
    // solium-disable-next-line indentation
    0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141;
  // Prime characteristic of the galois field over which Secp256k1 is defined
  uint256 constant private FIELD_SIZE =
    // solium-disable-next-line indentation
    0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F;
  uint256 constant private WORD_LENGTH_BYTES = 0x20;

  // (base^exponent) % FIELD_SIZE
  // Cribbed from https://medium.com/@rbkhmrcr/precompiles-solidity-e5d29bd428c4
  function bigModExp(uint256 base, uint256 exponent)
    internal view returns (uint256 exponentiation) {
      uint256 callResult;
      uint256[6] memory bigModExpContractInputs;
      bigModExpContractInputs[0] = WORD_LENGTH_BYTES;  // Length of base
      bigModExpContractInputs[1] = WORD_LENGTH_BYTES;  // Length of exponent
      bigModExpContractInputs[2] = WORD_LENGTH_BYTES;  // Length of modulus
      bigModExpContractInputs[3] = base;
      bigModExpContractInputs[4] = exponent;
      bigModExpContractInputs[5] = FIELD_SIZE;
      uint256[1] memory output;
      assembly { // solhint-disable-line no-inline-assembly
      callResult := staticcall(
        not(0),                   // Gas cost: no limit
        0x05,                     // Bigmodexp contract address
        bigModExpContractInputs,
        0xc0,                     // Length of input segment
        output,
        0x20                      // Length of output segment
      )
      }
      if (callResult == 0) {revert("bigModExp failure!");}
      return output[0];
    }

  // Let q=FIELD_SIZE. q % 4 = 3, ∴ p≡r^2 mod q ⇒ p^SQRT_POWER≡m±r mod q.  See
  // https://en.wikipedia.org/wiki/Modular_square_root#Prime_or_prime_power_modulus
  uint256 constant private SQRT_POWER = (FIELD_SIZE + 1) >> 2;

  // Computes a s.t. a^2 = x in the field. Assumes a exists
  function squareRoot(uint256 x) internal view returns (uint256) {
    return bigModExp(x, SQRT_POWER);
  }

  function ySquared(uint256 x) internal pure returns (uint256) {
    // Curve is y^2=x^3+7. See section 2.4.1 of https://www.secg.org/sec2-v2.pdf
    uint256 xCubed = mulmod(x, mulmod(x, x, FIELD_SIZE), FIELD_SIZE);
    return addmod(xCubed, 7, FIELD_SIZE);
  }

  function isOnCurve(uint256[2] memory p) internal pure returns (bool) {
    return ySquared(p[0]) == mulmod(p[1], p[1], FIELD_SIZE);
  }

  // Hash x uniformly into {0, ..., FIELD_SIZE-1}.
  function zqHash(uint256 x) internal pure returns (uint256 x_) {
    x_ = x;
    // Rejecting if x >= q corresponds to step 1 in section 2.3.6 of
    // http://www.secg.org/sec1-v2.pdf , which is part of the definition of
    // RS2ECP via section 2.3.4 via OS2ECP via
    // https://tools.ietf.org/html/rfc8032#section-5.1.3
    while (x_ >= FIELD_SIZE) {
      x_ = uint256(keccak256(abi.encodePacked(x_)));
    }
  }

  // One-way hash function onto the curve.
  function hashToCurve(uint256[2] memory pk, uint256 input)
    internal view returns (uint256[2] memory rv) {
      rv[0] = zqHash(uint256(keccak256(abi.encodePacked(pk, input))));
      rv[1] = squareRoot(ySquared(rv[0]));
      // Keep re-hashing until rv[1]^2 = rv[0]^3 + 7 mod P
      while (mulmod(rv[1], rv[1], FIELD_SIZE) != ySquared(rv[0])) {
        rv[0] = zqHash(uint256(keccak256(abi.encodePacked(rv[0]))));
        rv[1] = squareRoot(ySquared(rv[0]));
      }
      // See
      // https://tools.ietf.org/html/draft-goldbe-vrf-01#section-5.4.1.1
      // step 4.D, referencing RS2ECP,
      // https://tools.ietf.org/html/draft-goldbe-vrf-01#section-5.5 , for
      // definition of RS2ECP, and http://www.secg.org/sec1-v2.pdf#page=17
      // , steps 2.3-2.4.1 of section 2.3.4 for relevant part of OS2ECP
      // definition. Together, these specify that the y ordinate must be
      // even.
      if (rv[1] % 2 == 1) {
        rv[1] = FIELD_SIZE - rv[1];
      }
    }

  // Returns true iff q==scalar*x, with cryptographically high probability.
  // Based on Vitalik Buterin's idea in ethresear.ch post mentioned below.
  //
  // scalar must be non-zero
  function ecmulVerify(uint256[2] memory x, uint256 scalar, uint256[2] memory q)
    internal pure returns(bool) {
      require(scalar != 0); // Rules out an ecrecover failure case
      // This ecrecover returns the address associated with c*R. See
      // https://ethresear.ch/t/you-can-kinda-abuse-ecrecover-to-do-ecmul-in-secp256k1-today/2384/9
      // The point corresponding to the address returned by ecrecover(0,v,r,s=c*r)
      // is (r⁻¹ mod Q) * (c*r * R - 0 * g) = c * R, where R is the point
      // specified by (v, r). See https://crypto.stackexchange.com/a/18106
      bytes32 scalarTimesX0 = bytes32(mulmod(scalar, x[0], GROUP_ORDER));
      uint8 parity = x[1] % 2 != 0 ? 28 : 27;
      return ecrecover(bytes32(0), parity, bytes32(x[0]), scalarTimesX0) ==
        address(uint256(keccak256(abi.encodePacked(q)))); // Takes bottom 160 bits
    }

  // Returns x1/z1-x2/z2=(x1z2+x2z1)/(z1z2) in projective coordinates on P¹(𝔽ₙ)
  function projectiveSub(uint256 x1, uint256 z1, uint256 x2, uint256 z2)
    internal pure returns(uint256 x3, uint256 z3) {
      uint256 num1 = mulmod(z2, x1, FIELD_SIZE);
      uint256 num2 = mulmod(FIELD_SIZE - x2, z1, FIELD_SIZE);
      (x3, z3) = (addmod(num1, num2, FIELD_SIZE), mulmod(z1, z2, FIELD_SIZE));
    }

  // Returns x1/z1*x2/z2=(x1x2)/(z1z2), in projective coordinates on P¹(𝔽ₙ)
  function projectiveMul(uint256 x1, uint256 z1, uint256 x2, uint256 z2)
    internal pure returns(uint256 x3, uint256 z3) {
      (x3, z3) = (mulmod(x1, x2, FIELD_SIZE), mulmod(z1, z2, FIELD_SIZE));
    }

  /** *********************************************************************
      @notice Computes elliptic-curve sum, in projective co-ordinates

      @dev Using projective coordinates avoids costly divisions

      @dev To use this with p and p in affine coordinates, call
      @dev projectiveECAdd(px, py, qx, qy). This will return
      @dev the addition of (px, py, 1) and (qx, qy, 1), in the
      @dev secp256k1 group.

      @dev This can be used to calculate the z which is the inverse to zInv
      @dev in isValidVRFOutput. But consider using a faster
      @dev re-implementation.

      @dev This function assumes [px,py,1],[qx,qy,1] are valid projective
           coordinates of secp256k1 points. That is safe in this contract,
           because this method is only used by linearCombination, which checks
           points are on the curve via ecrecover.
      **************************************************************************
      @param px The first affine coordinate of the first summand
      @param py The second affine coordinate of the first summand
      @param qx The first affine coordinate of the second summand
      @param qy The second affine coordinate of the second summand

      (px,py) and (qx,qy) must be distinct, valid secp256k1 points.
      **************************************************************************
      @return [px,py,1]+[qx,qy,1] as points on secp256k1, in P²(𝔽ₙ)
  */
  function projectiveECAdd(uint256 px, uint256 py, uint256 qx, uint256 qy)
    internal pure returns(uint256 sx, uint256 sy, uint256 sz) {
      // See "Group law for E/K : y^2 = x^3 + ax + b", in section 3.1.2, p. 80,
      // "Guide to Elliptic Curve Cryptography" by Hankerson, Menezes and Vanstone
      // We take the equations there for (sx,sy), and homogenize them to
      // projective coordinates. That way, no inverses are required, here, and we
      // only need the one inverse in affineECAdd.

      // We only need the "point addition" equations from Hankerson et al. Can
      // skip the "point doubling" equations because p1 == p2 is cryptographically
      // impossible, and require'd not to be the case in linearCombination.

      // Add extra "projective coordinate" to the two points
      (uint256 z1, uint256 z2) = (1, 1);

      // (lx, lz) = (qy-py)/(qx-px), i.e., gradient of secant line.
      uint256 lx = addmod(qy, FIELD_SIZE - py, FIELD_SIZE);
      uint256 lz = addmod(qx, FIELD_SIZE - px, FIELD_SIZE);

      uint256 dx; // Accumulates denominator from sx calculation
      // sx=((qy-py)/(qx-px))^2-px-qx
      (sx, dx) = projectiveMul(lx, lz, lx, lz); // ((qy-py)/(qx-px))^2
      (sx, dx) = projectiveSub(sx, dx, px, z1); // ((qy-py)/(qx-px))^2-px
      (sx, dx) = projectiveSub(sx, dx, qx, z2); // ((qy-py)/(qx-px))^2-px-qx

      uint256 dy; // Accumulates denominator from sy calculation
      // sy=((qy-py)/(qx-px))(px-sx)-py
      (sy, dy) = projectiveSub(px, z1, sx, dx); // px-sx
      (sy, dy) = projectiveMul(sy, dy, lx, lz); // ((qy-py)/(qx-px))(px-sx)
      (sy, dy) = projectiveSub(sy, dy, py, z1); // ((qy-py)/(qx-px))(px-sx)-py

      if (dx != dy) { // Cross-multiply to put everything over a common denominator
        sx = mulmod(sx, dy, FIELD_SIZE);
        sy = mulmod(sy, dx, FIELD_SIZE);
        sz = mulmod(dx, dy, FIELD_SIZE);
      } else { // Already over a common denominator, use that for z ordinate
        sz = dx;
      }
    }

  // Returns p1+p2, as affine points on secp256k1.
  //
  // invZ must be the inverse of the z returned by projectiveECAdd(p1, p2).
  // It is computed off-chain to save gas.
  //
  // p1 and p2 must be distinct, because projectiveECAdd doesn't handle
  // point doubling.
  function affineECAdd(
    uint256[2] memory p1, uint256[2] memory p2,
    uint256 invZ) internal pure returns (uint256[2] memory) {
    uint256 x;
    uint256 y;
    uint256 z;
    (x, y, z) = projectiveECAdd(p1[0], p1[1], p2[0], p2[1]);
    require(mulmod(z, invZ, FIELD_SIZE) == 1, "invZ must be inverse of z");
    // Clear the z ordinate of the projective representation by dividing through
    // by it, to obtain the affine representation
    return [mulmod(x, invZ, FIELD_SIZE), mulmod(y, invZ, FIELD_SIZE)];
  }

  // Returns true iff address(c*p+s*g) == lcWitness, where g is generator.
  function verifyLinearCombinationWithGenerator(
    uint256 c, uint256[2] memory p, uint256 s, address lcWitness)
    internal pure returns (bool) {
      // Rule out ecrecover failure modes which return address 0.
      require(lcWitness != address(0), "bad witness");
      // https://ethresear.ch/t/you-can-kinda-abuse-ecrecover-to-do-ecmul-in-secp256k1-today/2384/9
      // The point corresponding to the address returned by
      // ecrecover(-s*p[0],v,p[0],c*p[0]) is
      // (p[0]⁻¹ mod GROUP_ORDER)*(c*p[0]-(-s)*p[0]*g)=c*p+s*g, where v
      // is the parity of p[1]. See https://crypto.stackexchange.com/a/18106
      bytes32 pseudoHash = bytes32(GROUP_ORDER - mulmod(p[0], s, GROUP_ORDER));
      // https://bitcoin.stackexchange.com/questions/38351/ecdsa-v-r-s-what-is-v
      uint8 v = (p[1] % 2 == 0) ? 27 : 28;
      bytes32 pseudoSignature = bytes32(mulmod(c, p[0], GROUP_ORDER));
      address computed = ecrecover(pseudoHash, v, bytes32(p[0]), pseudoSignature);
      return computed == lcWitness;
    }

  // c*p1 + s*p2. Requires cp1Witness=c*p1 and sp2Witness=s*p2. Also
  // requires cp1Witness != sp2Witness (which is fine for this application,
  // since it is cryptographically impossible for them to be equal. In the
  // (cryptographically impossible) case that a prover accidentally derives
  // a proof with equal c*p1 and s*p2, they should retry with a different
  // proof nonce.) Assumes that all points are on secp256k1.
  function linearCombination(
    uint256 c, uint256[2] memory p1, uint256[2] memory cp1Witness,
    uint256 s, uint256[2] memory p2, uint256[2] memory sp2Witness,
    uint256 zInv)
    internal pure returns (uint256[2] memory) {
      require((cp1Witness[0] - sp2Witness[0]) % FIELD_SIZE != 0,
              "points in sum must be distinct");
      require(ecmulVerify(p1, c, cp1Witness), "First multiplication check failed");
      require(ecmulVerify(p2, s, sp2Witness), "Second multiplication check failed");
      return affineECAdd(cp1Witness, sp2Witness, zInv);
    }

  // Pseudo-random number from inputs. Matches vrf.go/scalarFromCurve, and
  // https://tools.ietf.org/html/draft-goldbe-vrf-01#section-5.4.2 . It
  // draft calls (in section 5.3 step 5 and section 5.4.2 steps 3-5) for
  // taking the first hash without checking that it corresponds to a number
  // less than the group order, which will lead to a slight bias in the
  // sample.
  function scalarFromCurve(
    uint256[2] memory hash, uint256[2] memory pk, uint256[2] memory gamma,
    address uWitness, uint256[2] memory v)
    internal pure returns (uint256 s) {
      return uint256(
        keccak256(abi.encodePacked(hash, pk, gamma, v, uWitness)));
    }

  // True if (gamma, c, s) is a correctly constructed randomness proof from pk
  // and seed. zInv must be the inverse of the third ordinate from
  // projectiveECAdd applied to cGammaWitness and sHashWitness. Corresponds to
  // section 5.3 of the IETF draft.
  //
  // TODO(alx): Since I'm only using pk in the ecrecover call, I could only pass
  // the x ordinate, and the parity of the y ordinate in the top bit of uWitness
  // (which I could make a uint256 without using any extra space.) Would save
  // about 2000 gas. (Roughly 2.5%.)
  //
  // TODO(alx): It would probably also be fine to only pass the address
  // witnesses for cGammaWitness and sHashWitness, which would save another 4000
  // gas.
  function verifyVRFProof(
    uint256[2] memory pk, uint256[2] memory gamma, uint256 c, uint256 s,
    uint256 seed, address uWitness, uint256[2] memory cGammaWitness,
    uint256[2] memory sHashWitness, uint256 zInv)
    internal view {
      require(isOnCurve(pk), "public key is not on curve");
      require(isOnCurve(gamma), "gamma is not on curve");
      require(isOnCurve(cGammaWitness), "cGammaWitness is not on curve");
      require(isOnCurve(sHashWitness), "sHashWitness is not on curve");
      // Step 4. of IETF draft section 5.3 (pk corresponds to 5.3's y, and here
      // we use the hash of u instead of u itself.)
      require(
        verifyLinearCombinationWithGenerator(c, pk, s, uWitness),
        "addr(c*pk+s*g)≠_uWitness"
      );
      // Step 5. of IETF draft section 5.3 (pk corresponds to y, seed to beta)
      uint256[2] memory hash = hashToCurve(pk, seed);
      // Step 6. of IETF draft section 5.3
      uint256[2] memory v = linearCombination(
        c, gamma, cGammaWitness, s, hash, sHashWitness, zInv);
      // Steps 7. and 8. of IETF draft section 5.3
      require(c == scalarFromCurve(hash, pk, gamma, uWitness, v), "invalid proof");
    }

  /** **************************************************************************
      @notice isValidVRFOutput returns true iff the proof can be verified as
      showing that output was generated as mandated.

      @dev See the invocation of verifyVRFProof in VRF.js, for an example.
      **************************************************************************
      @dev Let x be the secret key associated with the public key pk (which is
           known as y in section 5.3 of the IETF draft.)

      @param pk Affine coordinates of the secp256k1 public key for this VRF.
      @param gamma Intermediate output of the VRF as an affine secp256k1 point
      @param c The challenge value for proof that gamma = x*hashToCurve(seed)
              See the variable c on  p. oeuta28 of
              https://www.cs.bu.edu/~gold-be/papers/VRF_ietf99_print.pdf
      @param s The response value for the proof. See s on p. 28
      @param seed The input seed from which the VRF output is computed. Also
             known as alpha in the IETF draft, section 5.3
      @param uWitness The ethereum address of c*pk + s*<generator>, in
             elliptic-curve arithmetic. This corresponds to u in section 5.3 of
             the IETF draft, but there it as an elliptic curve point, not an
             address.
      @param cGammaWitness c*gamma on the elliptic-curve
      @param sHashWitness s*hashToCurve(seed) on the elliptic-curve
      @param zInv Inverse of the third ordinate of the return value from
             projectiveECAdd(c*gamma, s*hashToCurve(seed)). Passed in here
             to save gas, because computing modular inverses is expensive in the
             EVM.
      @param output The actual output of the VRF. Known as beta in the
             IETF standard, section 5.3
      **************************************************************************
      @return True iff all the above parameters are correct
  */
  function isValidVRFOutput(
    uint256[2] calldata pk, uint256[2] calldata gamma, uint256 c, uint256 s,
    uint256 seed, address uWitness, uint256[2] calldata cGammaWitness,
    uint256[2] calldata sHashWitness, uint256 zInv, uint256 output)
    external view returns (bool) {
      return verifyVRFProof(
        pk, gamma, c, s, seed, uWitness, cGammaWitness, sHashWitness,
        zInv) &&
        (uint256(keccak256(abi.encodePacked(gamma))) == output);
    }
}
