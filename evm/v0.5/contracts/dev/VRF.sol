pragma solidity 0.5.0;

////////////////////////////////////////////////////////////////////////////////
//       XXX: Do not use in production until this code has been audited.
////////////////////////////////////////////////////////////////////////////////

/** ****************************************************************************
    @notice on-chain verification of verifiable-random-function (VRF) proofs as
            described in https://eprint.iacr.org/2017/099.pdf (security proofs)
            and https://tools.ietf.org/html/draft-goldbe-vrf-01#section-5 (spec)
    ****************************************************************************
    @dev PURPOSE

    @dev Reggie the Random Oracle (not his real job) wants to provide randomness
         to Vera the verifier in such a way that Vera can be sure he's not
         making his output up to suit himself. Reggie provides Vera a public key
         to which he knows the secret key. Each time Vera provides a seed to
         Reggie, he gives back a value which is computed completely
         deterministically from the seed and the secret key, but which is
         indistinguishable from randomness to Vera. Nonetheless, Vera is able to
         verify that Reggie's output came from her seed and his secret key.

    @dev The purpose of this contract is to perform that verification.
    ****************************************************************************
    @dev USAGE

    @dev The main entry point is isValidVRFOutput. See its docstring.
    Design notes
    ------------

    An elliptic curve point is generally represented in the solidity code as a
    uint256[2], corresponding to its affine coordinates in GF(fieldSize).

    For the sake of efficiency, this implementation deviates from the spec in
    some minor ways:

    - Keccak hash rather than SHA256. This is because it's provided natively by
      the EVM, and therefore costs much less gas. The impact on security should
      be minor.

    - Secp256k1 curve instead of P-256. It abuses ECRECOVER for the most
      expensive ECC arithmetic.

    - scalarFromCurve recursively hashes and takes the relevant hash bits until
      it finds a point less than the group order. This results in uniform
      sampling over the the possible values scalarFromCurve could take. The spec
      recommends just uing the first hash output as a uint256, which is a
      slightly biased sample. See the zqHash function.

    - hashToCurve recursively hashes until it finds a curve x-ordinate. The spec
      recommends that the initial input should be concatenated with a nonce and
      then hashed, and this input should be rehashed with the nonce updated
      until an x-ordinate is found. Recursive hashing is slightly more
      efficient. The spec also recommends
      (https://tools.ietf.org/html/rfc8032#section-5.1.3 , by the specification
      of RS2ECP) that the x-ordinate should be rejected if it is greater than
      the modulus.

    - In the calculation of the challenge value "c", the "u" value (or "k*g", if
      you know the secret nonce)

    The spec also requires the y ordinate of the hashToCurve to be negated if y
    is odd. See http://www.secg.org/sec1-v2.pdf#page=17 . This sacrifices one
    bit of entropy in the random output. Instead, here y is chosen based on
    whether an extra hash of the inputs is even or odd. */

contract VRF {

  // See https://en.bitcoin.it/wiki/Secp256k1 for these constants.
  uint256 constant public GROUP_ORDER = // Number of points in Secp256k1
    // solium-disable-next-line indentation
    0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141;
  // Prime characteristic of the galois field over which Secp256k1 is defined
  // solium-disable-next-line zeppelin/no-arithmetic-operations
  uint256 constant public FIELD_SIZE =
    // solium-disable-next-line indentation
    0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F;

  // solium-disable zeppelin/no-arithmetic-operations
  uint256 constant public MINUS_ONE = FIELD_SIZE - 1;
  uint256 constant public MULTIPLICATIVE_GROUP_ORDER = FIELD_SIZE - 1;
  // pow(x, SQRT_POWER, FIELD_SIZE) == √x, since FIELD_SIZE % 4 = 3
  // https://en.wikipedia.org/wiki/Modular_square_root#Prime_or_prime_power_modulus
  uint256 constant public SQRT_POWER = (FIELD_SIZE + 1) >> 2;

  uint256 constant public WORD_LENGTH_BYTES = 0x20;

  // (base**exponent) % modulus
  // Cribbed from https://medium.com/@rbkhmrcr/precompiles-solidity-e5d29bd428c4
  function bigModExp(uint256 base, uint256 exponent, uint256 modulus)
    public view returns (uint256 exponentiation) {
      uint256 callResult;
      uint256[6] memory bigModExpContractInputs;
      bigModExpContractInputs[0] = WORD_LENGTH_BYTES;  // Length of base
      bigModExpContractInputs[1] = WORD_LENGTH_BYTES;  // Length of exponent
      bigModExpContractInputs[2] = WORD_LENGTH_BYTES;  // Length of modulus
      bigModExpContractInputs[3] = base;
      bigModExpContractInputs[4] = exponent;
      bigModExpContractInputs[5] = modulus;
      uint256[1] memory output;
      assembly { // solhint-disable-line no-inline-assembly
      callResult :=
        staticcall(13056,           // Gas cost. See EIP-198's 1st e.g.
          0x05,                     // Bigmodexp contract address
          bigModExpContractInputs,
          0xc0,                     // Length of input segment
          output,
          0x20)                     // Length of output segment
      }
      if (callResult == 0) {revert("bigModExp failure!");}
      return output[0];
    }

  // Computes a s.t. a^2 = x in the field. Assumes x is a square.
  function squareRoot(uint256 x) public view returns (uint256) {
    return bigModExp(x, SQRT_POWER, FIELD_SIZE);
  }

  function ySquared(uint256 x) public view returns (uint256) {
    // Curve equation is y^2=x^3+7. See
    return (bigModExp(x, 3, FIELD_SIZE) + 7) % FIELD_SIZE;
  }

  // Hash x uniformly into {0, ..., q-1}. Expects x to ALREADY have the
  // necessary entropy... If x < q, returns x!
  function zqHash(uint256 q, uint256 x) public pure returns (uint256 x_) {
    x_ = x;
    while (x_ >= q) {
      x_ = uint256(keccak256(abi.encodePacked(x_)));
    }
  }

  // One-way hash function onto the curve.
  function hashToCurve(uint256[2] memory k, uint256 input)
    public view returns (uint256[2] memory rv) {
      bytes32 hash = keccak256(abi.encodePacked(k, input));
      rv[0] = zqHash(FIELD_SIZE, uint256(hash));
      while (true) {
        rv[0] = zqHash(FIELD_SIZE, uint256(keccak256(abi.encodePacked(rv[0]))));
        rv[1] = squareRoot(ySquared(rv[0]));
        if (mulmod(rv[1], rv[1], FIELD_SIZE) == ySquared(rv[0])) {
          break;
        }
      }
      // Two possible y ordinates for x ordinate rv[0]; pick one "randomly"
      if (uint256(keccak256(abi.encodePacked(rv[0], input))) % 2 == 0) {
        rv[1] = -rv[1];
      }
    }

  // Bits used in Ethereum address
  uint256 constant public BOTTOM_160_BITS = 2**161 - 1;

  // Returns the ethereum address associated with point.
  function pointAddress(uint256[2] calldata point) external pure returns(address) {
    bytes memory packedPoint = abi.encodePacked(point);
    // Lower 160 bits of the keccak hash of (x,y) as 64 bytes
    return address(uint256(keccak256(packedPoint)) & BOTTOM_160_BITS);
  }

  // Returns true iff q==scalar*x, with cryptographically high probability.
  // Based on Vitalik Buterin's idea in above ethresear.ch post.
  function ecmulVerify(uint256[2] memory x, uint256 scalar, uint256[2] memory q)
    public pure returns(bool) {
      // This ecrecover returns the address associated with c*R. See
      // https://ethresear.ch/t/you-can-kinda-abuse-ecrecover-to-do-ecmul-in-secp256k1-today/2384/9
      // The point corresponding to the address returned by ecrecover(0,v,r,s=c*r)
      // is (r⁻¹ mod Q) * (c*r * R - 0 * g) = c * R, where R is the point
      // specified by (v, r). See https://crypto.stackexchange.com/a/18106
      bytes32 cTimesX0 = bytes32(mulmod(scalar, x[0], GROUP_ORDER));
      uint8 parity = x[1] % 2 != 0 ? 28 : 27;
      return ecrecover(bytes32(0), parity, bytes32(x[0]), cTimesX0) ==
        address(uint256(keccak256(abi.encodePacked(q))) & BOTTOM_160_BITS);
    }

  // Returns x1/z1+x2/z2=(x1z2+x2z1)/(z1z2) in projective coordinates on P¹(𝔽ₙ)
  function projectiveAdd(uint256 x1, uint256 z1, uint256 x2, uint256 z2)
    external pure returns(uint256 x3, uint256 z3) {
      uint256 crossMultNumerator1 = mulmod(z2, x1, FIELD_SIZE);
      uint256 crossMultNumerator2 = mulmod(z1, x2, FIELD_SIZE);
      uint256 denom = mulmod(z1, z2, FIELD_SIZE);
      uint256 numerator = addmod(crossMultNumerator1, crossMultNumerator2, FIELD_SIZE);
      return (numerator, denom);
    }

  // Returns x1/z1-x2/z2=(x1z2+x2z1)/(z1z2) in projective coordinates on P¹(𝔽ₙ)
  function projectiveSub(uint256 x1, uint256 z1, uint256 x2, uint256 z2)
    public pure returns(uint256 x3, uint256 z3) {
      uint256 num1 = mulmod(z2, x1, FIELD_SIZE);
      uint256 num2 = mulmod(FIELD_SIZE - x2, z1, FIELD_SIZE);
      (x3, z3) = (addmod(num1, num2, FIELD_SIZE), mulmod(z1, z2, FIELD_SIZE));
    }

  // Returns x1/z1*x2/z2=(x1x2)/(z1z2), in projective coordinates on P¹(𝔽ₙ)
  function projectiveMul(uint256 x1, uint256 z1, uint256 x2, uint256 z2)
    public pure returns(uint256 x3, uint256 z3) {
      (x3, z3) = (mulmod(x1, x2, FIELD_SIZE), mulmod(z1, z2, FIELD_SIZE));
    }

  // Returns x1/z1/(x2/z2)=(x1z2)/(x2z1), in projective coordinates on P¹(𝔽ₙ)
  function projectiveDiv(uint256 x1, uint256 z1, uint256 x2, uint256 z2)
    external pure returns(uint256 x3, uint256 z3) {
      (x3, z3) = (mulmod(x1, z2, FIELD_SIZE), mulmod(z1, x2, FIELD_SIZE));
    }

  /** **************************************************************************
      @notice Computes elliptic-curve sum, in projective co-ordinates

      @dev Using projective coordinates avoids costly divisions

      @dev To use this with x and y in affine coordinates, compute
           projectiveECAdd(x[0], x[1], 1, y[0], y[1], 1)

      @dev This can be used to calculate the z which is the inverse to zInv in
           isValidVRFOutput. But consider using a faster re-implementation.

      @dev This function assumes [x1,y1,z1],[x2,y2,z2] are valid projective
           coordinates of secp256k1 points. That is safe in this contract,
           because this method is only used by linearCombination, which checks
           points are on the curve via ecrecover, and ensures valid projective
           coordinates by passing z1=z2=1.
      **************************************************************************
      @param x1 The first affine coordinate of the first summand
      @param y1 The second affine coordinate of the first summand
      @param x2 The first affine coordinate of the second summand
      @param y2 The second affine coordinate of the second summand
      **************************************************************************
      @return [x1,y1,z1]+[x2,y2,z2] as points on secp256k1, in P²(𝔽ₙ)
  */
  function projectiveECAdd(uint256 x1, uint256 y1, uint256 x2, uint256 y2)
    public pure returns(uint256 x3, uint256 y3, uint256 z3) {
      // See "Group law for E/K : y^2 = x^3 + ax + b", in section 3.1.2, p. 80,
      // "Guide to Elliptic Curve Cryptography" by Hankerson, Menezes and Vanstone
      // We take the equations there for (x3,y3), and homogenize them to
      // projective coordinates. That way, no inverses are required, here, and we
      // only need the one inverse in affineECAdd.
      
      // We only need the "point addition" equations from Hankerson et al. Can
      // skip the "point doubling" equations because p1 == p2 is cryptographically
      // impossible, and require'd not to be the case in linearCombination.
      
      // Add extra "projective coordinate" to the two points
      (uint256 z1, uint256 z2) = (1, 1);
      
      // (lx, lz) = (y2-y1)/(x2-x1), i.e., gradient of secant line.
      uint256 lx = addmod(y2, FIELD_SIZE - y1, FIELD_SIZE);
      uint256 lz = addmod(x2, FIELD_SIZE - x1, FIELD_SIZE);
      
      uint256 dx; // Accumulates denominator from x3 calculation
      // x3=((y2-y1)/(x2-x1))^2-x1-x2
      (x3, dx) = projectiveMul(lx, lz, lx, lz); // ((y2-y1)/(x2-x1))^2
      (x3, dx) = projectiveSub(x3, dx, x1, z1); // ((y2-y1)/(x2-x1))^2-x1
      (x3, dx) = projectiveSub(x3, dx, x2, z2); // ((y2-y1)/(x2-x1))^2-x1-x2
      
      uint256 dy; // Accumulates denominator from y3 calculation
      // y3=((y2-y1)/(x2-x1))(x1-x3)-y1
      (y3, dy) = projectiveSub(x1, z1, x3, dx); // x1-x3
      (y3, dy) = projectiveMul(y3, dy, lx, lz); // ((y2-y1)/(x2-x1))(x1-x3)
      (y3, dy) = projectiveSub(y3, dy, y1, z1); // ((y2-y1)/(x2-x1))(x1-x3)-y1
      
      if (dx != dy) { // Cross-multiply to put everything over a common denominator
        x3 = mulmod(x3, dy, FIELD_SIZE);
        y3 = mulmod(y3, dx, FIELD_SIZE);
        z3 = mulmod(dx, dy, FIELD_SIZE);
      } else {
        z3 = dx;
      }
    }

  // Returns p1+p2, as affine points on secp256k1. invZ must be the inverse of
  // the z returned by projectiveECAdd(p1, p2). It is computed off-chain to
  // save gas.
  function affineECAdd(
    uint256[2] memory p1, uint256[2] memory p2,
    uint256 invZ) public pure returns (uint256[2] memory) {
    uint256 x;
    uint256 y;
    uint256 z;
    (x, y, z) = projectiveECAdd(p1[0], p1[1], p2[0], p2[1]);
    require(mulmod(z, invZ, FIELD_SIZE) == 1, "_invZ must be inverse of z");
    // Clear the z ordinate of the projective representation by dividing through
    // by it, to obtain the affine representation
    return [mulmod(x, invZ, FIELD_SIZE), mulmod(y, invZ, FIELD_SIZE)];
  }

  // Returns true iff address(c*p+s*g) == lcWitness, where g is generator.
  function verifyLinearCombinationWithGenerator(
    uint256 c, uint256[2] memory p, uint256 s, address lcWitness)
    public pure returns (bool) {
      // ecrecover returns 0x0 in certain failure modes. Ensure witness differs.
      require(lcWitness != address(0), "bad witness");
      // https://ethresear.ch/t/you-can-kinda-abuse-ecrecover-to-do-ecmul-in-secp256k1-today/2384/9
      // The point corresponding to the address returned by
      // ecrecover(-s*p[0],v,_p[0],_c*p[0]) is
      // (p[0]⁻¹ mod GROUP_ORDER)*(c*p[0]-(-s)*p[0]*g)=_c*p+s*g, where v
      // is the parity of p[1]. See https://crypto.stackexchange.com/a/18106
      bytes32 pseudoHash = bytes32(GROUP_ORDER - mulmod(p[0], s, GROUP_ORDER));
      // https://bitcoin.stackexchange.com/questions/38351/ecdsa-v-r-s-what-is-v
      uint8 v = (p[1] % 2 == 0) ? 27 : 28;
      bytes32 pseudoSignature = bytes32(mulmod(c, p[0], GROUP_ORDER));
      address computed = ecrecover(pseudoHash, v, bytes32(p[0]), pseudoSignature);
      return computed == lcWitness;
    }

  // c*p1 + s*p2
  function linearCombination(
    uint256 c, uint256[2] memory p1, uint256[2] memory cp1Witness,
    uint256 s, uint256[2] memory p2, uint256[2] memory sp2Witness,
    uint256 zInv)
    public pure returns (uint256[2] memory) {
      require(cp1Witness[0] != sp2Witness[0], "points must differ in sum");
      require(ecmulVerify(p1, c, cp1Witness), "First multiplication check failed");
      require(ecmulVerify(p2, s, sp2Witness), "Second multiplication check failed");
      return affineECAdd(cp1Witness, sp2Witness, zInv);
    }

  // Pseudo-random number from inputs. Corresponds to vrf.go/scalarFromCurve.
  function scalarFromCurve(
    uint256[2] memory hash, uint256[2] memory pk, uint256[2] memory gamma,
    address uWitness, uint256[2] memory v)
    public pure returns (uint256 s) {
      bytes32 iHash = keccak256(abi.encodePacked(hash, pk, gamma, v, uWitness));
      return zqHash(GROUP_ORDER, uint256(iHash));
    }

  // True if (gamma, c, s) is a correctly constructed randomness proof from pk
  // and seed. zInv must be the inverse of the third ordinate from
  // projectiveECAdd applied to cGammaWitness and sHashWitness
  function verifyVRFProof(
    uint256[2] memory pk, uint256[2] memory gamma, uint256 c, uint256 s,
    uint256 seed, address uWitness, uint256[2] memory cGammaWitness,
    uint256[2] memory sHashWitness, uint256 zInv)
    public view returns (bool) {
    // NB: Curve operations already check that (pkX, pkY), (gammaX, gammaY)
    // are valid curve points. No need to do that explicitly.
      require(
        verifyLinearCombinationWithGenerator(c, pk, s, uWitness),
        "Could not verify that address(c*pk+s*generator)=_uWitness");
      uint256[2] memory hash = hashToCurve(pk, seed);
      uint256[2] memory v = linearCombination(
        c, gamma, cGammaWitness, s, hash, sHashWitness, zInv);
      return (c == scalarFromCurve(hash, pk, gamma, uWitness, v));
    }

  /** **************************************************************************
      @notice isValidVRFOutput returns true iff the proof can be verified as
      showing that output was generated as mandated.

      @dev See the invocation of verifyVRFProof in VRF.js, for an example.
      **************************************************************************
      @dev Let x be the secret key associated with the public key pk

      @param pk Affine coordinates of the secp256k1 public key for this VRF
      @param gamma Intermediate output of the VRF as an affine secp256k1 point
      @param c The challenge value for proof that gamma = x*hashToCurve(seed)
              See the variable c on  p. 28 of
              https://www.cs.bu.edu/~goldbe/papers/VRF_ietf99_print.pdf
      @param s The response value for the proof. See s on p. 28
      @param seed The input seed from which the VRF output is computed
      @param uWitness The ethereum address of c*pk + s*<generator>, in
             elliptic-curve arithmetic
      @param cGammaWitness c*gamma on the elliptic-curve
      @param sHashWitness s*hashToCurve(seed) on the elliptic-curve
      @param zInv Inverse of the third ordinate of the return value from
             projectiveECAdd(c*gamma, s*hashToCurve(seed)). Passed in here
             to save gas, because computing modular inverses is expensive in the
             EVM.
      @param output The actual output of the VRF.
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
