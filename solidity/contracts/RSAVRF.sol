pragma solidity 0.4.24;

contract RSAVRF {

  uint256 constant keySizeBits = 2048;  // Size of RSA modulus
  // solium-disable-next-line zeppelin/no-arithmetic-operations
  uint256 constant keySizeWords = keySizeBits / 256;
  // solium-disable-next-line zeppelin/no-arithmetic-operations
  uint256 constant keySizeBytes = keySizeBits / 8;

  // _RSAModulus must be in big-endian format
  constructor(uint256[keySizeWords] memory _RSAModulus) public {
    require(inputSize == 640, "Changed key size: adjust unrolled copy in bigModExp");
    require(publicExponent == 3, "Changed exponent sie; adjust gasCost calculation");
    require(keySizeBits % 256 == 0, "Key size must be multiple of uint256's.");
    modulus = _RSAModulus;
  }

  // The traditional public key
  //
  // A publicExponent of 3 drastically reduces gas costs. See gasCost.
  uint256 constant publicExponent = 3;
  uint256[keySizeWords] public modulus;  // Set during initialization

  uint256 constant wordSizeBytes = 32;

  // solium-disable-next-line zeppelin/no-arithmetic-operations
  uint256 constant inputSize = 3 * wordSizeBytes + keySizeBytes + wordSizeBytes + keySizeBytes;
  //                           Argument lengths  + base length  +exponent length+ modulus length
  uint256 constant inputSizeBytes = inputSize / wordSizeBytes;

  // Gas cost for bigModExp precompiled contract call.
  // https://github.com/ethereum/EIPs/blob/master/EIPS/eip-198.md#specification
  //
  // exponent of 3 results in ADJUSTED_EXPONENT_LENGTH of 1, so it can be
  // excluded from the calculation.
  // solium-disable-next-line zeppelin/no-arithmetic-operations
  uint256 constant gasCost = (keySizeBytes ** 2) / 4 + 96 * keySizeBytes - 3072;

  // (_base**exponent) % modulus. _base must be in big-endian format
  function bigModExp(uint256[keySizeWords] memory _base)
    public view returns(uint256[keySizeWords] memory result) {

    // Lay out the arguments in memory as bigModExp expects them:
    //
    // base-length||exponent-length||modulus-length||   base     ||publicExponent||modulus
    //   1 word   ||    1 word     ||   1 word     ||keySizeBytes||    1 word    ||keySizeBytes
    //
    // Length units are bytes, each number is in big-endian format,
    uint256[inputSize] memory inputs;
    inputs[0] = keySizeBytes;  // _base length
    inputs[1] = wordSizeBytes; // exponent length
    inputs[2] = keySizeBytes;  // modulus length
    inputs[3] = _base[0];
    inputs[4] = _base[1];  // It's just a little more efficient, to unroll these
    inputs[5] = _base[2];
    inputs[6] = _base[3];
    inputs[7] = _base[4];
    inputs[8] = _base[5];
    inputs[9] = _base[6];
    inputs[10] = _base[7];

    inputs[11] = publicExponent;

    inputs[12] = modulus[0];
    inputs[13] = modulus[1];
    inputs[14] = modulus[2];
    inputs[15] = modulus[3];
    inputs[16] = modulus[4];
    inputs[17] = modulus[5];
    inputs[18] = modulus[6];
    inputs[19] = modulus[7];
    assert(7 == keySizeWords-1);
    assert(19 == inputSizeBytes-1); // That unrolling changes if key size changes

    // Now, do the bigModExp
    int success;
    uint256 inputsLength = inputSize; // No constants in assembly
    uint256 resultsLength = keySizeBytes;
    uint256 gcost = gasCost;
    assembly{
      result := mload(0x40)  // Store result at start of free memory
      mstore(0x40, add(result, resultsLength))  // Move freemem ptr past result
      success := staticcall(
        gcost,
        0x05,                  // BigModExp contract address
        inputs, inputsLength,  // Input segment
        result, resultsLength) // Output segment
    }
    if (success == 0) {revert("bigModExp call failed");}
  }
}
