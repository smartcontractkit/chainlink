pragma solidity 0.4.24;

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
  uint256 constant GROUP_ORDER = // Number of points in Secp256k1
    // solium-disable-next-line indentation
    0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141;
  // Prime characteristic of the galois field over which Secp256k1 is defined
  // solium-disable-next-line zeppelin/no-arithmetic-operations
  uint256 constant FIELD_SIZE =
    // solium-disable-next-line indentation
    0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F;

  // solium-disable zeppelin/no-arithmetic-operations
  uint256 constant MINUS_ONE = FIELD_SIZE - 1;
  uint256 constant MULTIPLICATIVE_GROUP_ORDER = FIELD_SIZE - 1;
  // pow(x, SQRT_POWER, FIELD_SIZE) == √x, since FIELD_SIZE % 4 = 3
  // https://en.wikipedia.org/wiki/Modular_square_root#Prime_or_prime_power_modulus
  uint256 constant SQRT_POWER = (FIELD_SIZE + 1) >> 2;

  uint256 constant WORD_LENGTH_BYTES = 0x20;

  // (_base**_exponent) % _modulus
  // Cribbed from https://medium.com/@rbkhmrcr/precompiles-solidity-e5d29bd428c4
  function bigModExp(uint256 _base, uint256 _exponent, uint256 _modulus)
    public view returns (uint256 exponentiation) {
    uint256 callResult;
    uint256[6] memory bigModExpContractInputs;
    bigModExpContractInputs[0] = WORD_LENGTH_BYTES;  // Length of _base
    bigModExpContractInputs[1] = WORD_LENGTH_BYTES;  // Length of _exponent
    bigModExpContractInputs[2] = WORD_LENGTH_BYTES;  // Length of _modulus
    bigModExpContractInputs[3] = _base;
    bigModExpContractInputs[4] = _exponent;
    bigModExpContractInputs[5] = _modulus;
    uint256[1] memory output;
    assembly {
      callResult :=
        staticcall(13056,                    // Gas cost. See EIP-198's 1st e.g.
                   0x05,                     // Bigmodexp contract address
                   bigModExpContractInputs,
                   0xc0,                     // Length of input segment
                   output,
                   0x20)                     // Length of output segment
      }
    if (callResult == 0) {revert("bigModExp failure!");}
    return output[0];
  }

  // Computes a s.t. a^2 = _x in the field. Assumes _x is a square.
  function squareRoot(uint256 _x) public view returns (uint256) {
    return bigModExp(_x, SQRT_POWER, FIELD_SIZE);
  }

  function ySquared(uint256 _x) public view returns (uint256) {
    // Curve equation is y^2=_x^3+7. See
    return (bigModExp(_x, 3, FIELD_SIZE) + 7) % FIELD_SIZE;
  }

  // Hash _x uniformly into {0, ..., q-1}. Expects _x to ALREADY have the
  // necessary entropy... If _x < q, returns _x!
  function zqHash(uint256 q, uint256 _x) public pure returns (uint256 x) {
    x = _x;
    while (x >= q) {
      x = uint256(keccak256(abi.encodePacked(x)));
    }
  }

  // One-way hash function onto the curve.
  function hashToCurve(uint256[2] memory _k, uint256 _input)
    public view returns (uint256[2] memory rv) {
    bytes32 hash = keccak256(abi.encodePacked(_k, _input));
    rv[0] = zqHash(FIELD_SIZE, uint256(hash));
    while (true) {
      rv[0] = zqHash(FIELD_SIZE, uint256(keccak256(abi.encodePacked(rv[0]))));
      rv[1] = squareRoot(ySquared(rv[0]));
      if (mulmod(rv[1], rv[1], FIELD_SIZE) == ySquared(rv[0])) {
        break;
      }
    }
    // Two possible y ordinates for x ordinate rv[0]; pick one "randomly"
    if (uint256(keccak256(abi.encodePacked(rv[0], _input))) % 2 == 0) {
      rv[1] = -rv[1];
    }
  }

  // Bits used in Ethereum address
  uint256 constant public BOTTOM_160_BITS = 2**161 - 1;

  // Returns the ethereum address associated with point.
  function pointAddress(uint256[2] memory point) public pure returns(address) {
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
    public pure returns(uint256 x3, uint256 z3) {
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
    public pure returns(uint256 x3, uint256 z3) {
    (x3, z3) = (mulmod(x1, z2, FIELD_SIZE), mulmod(z1, x2, FIELD_SIZE));
  }

  /** **************************************************************************
      @notice Computes elliptic-curve sum, in projective co-ordinates

      @dev Using projective coordinates avoids costly divisions

      @dev To use this with x and y in affine coordinates, compute
           projectiveECAdd(x[0], x[1], 1, y[0], y[1], 1)

      @dev This can be used to calculate the z which is the inverse to _zInv in
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

  // Returns p1+p2, as affine points on secp256k1. _invZ must be the inverse of
  // the z returned by projectiveECAdd(_p1, _p2). It is computed off-chain to
  // save gas.
  function affineECAdd(
    uint256[2] memory _p1, uint256[2] memory _p2,
    uint256 _invZ) public pure returns (uint256[2] memory) {
    uint256 x;
    uint256 y;
    uint256 z;
    (x, y, z) = projectiveECAdd(_p1[0], _p1[1], _p2[0], _p2[1]);
    require(mulmod(z, _invZ, FIELD_SIZE) == 1, "_invZ must be inverse of z");
    // Clear the z ordinate of the projective representation by dividing through
    // by it, to obtain the affine representation
    return [mulmod(x, _invZ, FIELD_SIZE), mulmod(y, _invZ, FIELD_SIZE)];
  }

  // Returns true iff address(_c*_p+_s*g) == _lcWitness, where g is generator.
  function verifyLinearCombinationWithGenerator(
    uint256 _c, uint256[2] memory _p, uint256 _s, address _lcWitness)
    public pure returns (bool) {
    // ecrecover returns 0x0 in certain failure modes. Ensure witness differs.
    require(_lcWitness != address(0), "bad witness");
    // https://ethresear.ch/t/you-can-kinda-abuse-ecrecover-to-do-ecmul-in-secp256k1-today/2384/9
    // The point corresponding to the address returned by
    // ecrecover(-_s*_p[0],v,_p[0],_c*_p[0]) is
    // (_p[0]⁻¹ mod GROUP_ORDER)*(_c*_p[0]-(-_s)*_p[0]*g)=_c*_p+_s*g, where v
    // is the parity of _p[1]. See https://crypto.stackexchange.com/a/18106
    bytes32 pseudoHash = bytes32(GROUP_ORDER - mulmod(_p[0], _s, GROUP_ORDER));
    // https://bitcoin.stackexchange.com/questions/38351/ecdsa-v-r-s-what-is-v
    uint8 v = (_p[1] % 2 == 0) ? 27 : 28;
    bytes32 pseudoSignature = bytes32(mulmod(_c, _p[0], GROUP_ORDER));
    address computed = ecrecover(pseudoHash, v, bytes32(_p[0]), pseudoSignature);
    return computed == _lcWitness;
  }

  // _c*_p1 + _s*_p2
  function linearCombination(
    uint256 _c, uint256[2] memory _p1, uint256[2] memory _cp1Witness,
    uint256 _s, uint256[2] memory _p2, uint256[2] memory _sp2Witness,
    uint256 _zInv)
    public pure returns (uint256[2] memory) {
    require(_cp1Witness[0] != _sp2Witness[0], "points must differ in sum");
    require(ecmulVerify(_p1, _c, _cp1Witness), "First multiplication check failed");
    require(ecmulVerify(_p2, _s, _sp2Witness), "Second multiplication check failed");
    return affineECAdd(_cp1Witness, _sp2Witness, _zInv);
  }

  // Pseudo-random number from inputs. Corresponds to vrf.go/scalarFromCurve.
  function scalarFromCurve(
    uint256[2] memory _hash, uint256[2] memory _pk, uint256[2] memory _gamma,
    address _uWitness, uint256[2] memory _v)
    public pure returns (uint256 s) {
    bytes32 iHash = keccak256(abi.encodePacked(_hash, _pk, _gamma, _v, _uWitness));
    return zqHash(GROUP_ORDER, uint256(iHash));
  }

  // True if (gamma, c, s) is a correctly constructed randomness proof from _pk
  // and _seed. _zInv must be the inverse of the third ordinate from
  // projectiveECAdd applied to _cGammaWitness and _sHashWitness
  function verifyVRFProof(
    uint256[2] memory _pk, uint256[2] memory _gamma, uint256 _c, uint256 _s,
    uint256 _seed, address _uWitness, uint256[2] memory _cGammaWitness,
    uint256[2] memory _sHashWitness, uint256 _zInv)
    public view returns (bool) {
    // NB: Curve operations already check that (_pkX, _pkY), (_gammaX, _gammaY)
    // are valid curve points. No need to do that explicitly.
    require(
      verifyLinearCombinationWithGenerator(_c, _pk, _s, _uWitness),
      "Could not verify that address(_c*_pk+_s*generator)=_uWitness");
    uint256[2] memory hash = hashToCurve(_pk, _seed);
    uint256[2] memory v = linearCombination(
      _c, _gamma, _cGammaWitness, _s, hash, _sHashWitness, _zInv);
    return (_c == scalarFromCurve(hash, _pk, _gamma, _uWitness, v));
  }

  /** **************************************************************************
      @notice isValidVRFOutput returns true iff the proof can be verified as
      showing that _output was generated as mandated.

      @dev See the invocation of verifyVRFProof in VRF.js, for an example.
      **************************************************************************
      @dev Let x be the secret key associated with the public key _pk

      @param _pk Affine coordinates of the secp256k1 public key for this VRF
      @param _gamma Intermediate output of the VRF as an affine secp256k1 point
      @param _c The challenge value for proof that _gamma = x*hashToCurve(_seed)
              See the variable c on  p. 28 of
              https://www.cs.bu.edu/~goldbe/papers/VRF_ietf99_print.pdf
      @param _s The response value for the proof. See s on p. 28
      @param _seed The input seed from which the VRF output is computed
      @param _uWitness The ethereum address of _c*_pk + _s*<generator>, in
             elliptic-curve arithmetic
      @param _cGammaWitness _c*_gamma on the elliptic-curve
      @param _sHashWitness _s*hashToCurve(_seed) on the elliptic-curve
      @param _zInv Inverse of the third ordinate of the return value from
             projectiveECAdd(_c*_gamma, _s*hashToCurve(_seed)). Passed in here
             to save gas, because computing modular inverses is expensive in the
             EVM.
      @param _output The actual output of the VRF.
      **************************************************************************
      @return True iff all the above parameters are correct
  */
  function isValidVRFOutput(
    uint256[2] memory _pk, uint256[2] memory _gamma, uint256 _c, uint256 _s,
    uint256 _seed, address _uWitness, uint256[2] memory _cGammaWitness,
    uint256[2] memory _sHashWitness, uint256 _zInv, uint256 _output)
    public view returns (bool) {
    return verifyVRFProof(
      _pk, _gamma, _c, _s, _seed, _uWitness, _cGammaWitness, _sHashWitness,
      _zInv) &&
      (uint256(keccak256(abi.encodePacked(_gamma))) == _output);
  }
}
