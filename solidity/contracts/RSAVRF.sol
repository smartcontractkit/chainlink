pragma solidity 0.4.24;

contract RSAVRF {

  uint256 constant keySizeBits = 4096;  // Size of RSA modulus
  // solium-disable-next-line zeppelin/no-arithmetic-operations
  uint256 constant keySizeWords = keySizeBits / 256;
  // solium-disable-next-line zeppelin/no-arithmetic-operations
  uint256 constant keySizeBytes = keySizeBits / 8;

  constructor(uint256[keySizeWords] memory _RSAModulus) public {
    require(inputSize == 1152, "Changed key size: adjust unrolled copy in bigModExp");
    require(keySizeBits % 256 == 0, "Key size must be multiple of uint256's.");
    modulus = _RSAModulus;
  }

  // The traditional public key
  uint256 constant publicExponent = 3;
  uint256[keySizeWords] public modulus;  // Set during initialization

  // solium-disable-next-line zeppelin/no-arithmetic-operations
  uint256 constant inputSize = 3*32 + keySizeBytes + 32 + keySizeBytes;
  // solium-disable-next-line zeppelin/no-arithmetic-operations
  uint256 constant gasCost = (keySizeBytes**2)/4 + 96 * keySizeBytes - 3072;

  // (_base**exponent) % modulus.
  function bigModExp(uint256[keySizeWords] memory _base)
    public view returns(uint256[keySizeWords] memory result) {

    // Lay out the arguments in memory as bigModExp expects them:
    // base-length||exponent-length||modulus-length||base||exponent||length
    //   1 word         1 word          1 word
    // Here base-length units are bytes, each number is in little-endian format,
    // solium-disable-next-line zeppelin/no-arithmetic-operations
    uint256[inputSize] memory inputs;
    inputs[0] = keySizeBytes;  // _base length
    inputs[1] = 32;            // exponent length
    inputs[2] = keySizeBytes;  // modulus length
    inputs[3] = _base[0];
    inputs[4] = _base[1];  // It's just a little more efficient, to unroll these
    inputs[5] = _base[2];
    inputs[6] = _base[3];
    inputs[7] = _base[4];
    inputs[8] = _base[5];
    inputs[9] = _base[6];
    inputs[10] = _base[7];
    inputs[11] = _base[8];
    inputs[12] = _base[9];
    inputs[13] = _base[10];
    inputs[14] = _base[11];
    inputs[15] = _base[12];
    inputs[16] = _base[13];
    inputs[17] = _base[14];
    inputs[18] = _base[15];

    inputs[19] = publicExponent;

    inputs[20] = modulus[0];
    inputs[21] = modulus[1];
    inputs[22] = modulus[2];
    inputs[23] = modulus[3];
    inputs[24] = modulus[4];
    inputs[25] = modulus[5];
    inputs[26] = modulus[6];
    inputs[27] = modulus[7];
    inputs[28] = modulus[8];
    inputs[29] = modulus[9];
    inputs[30] = modulus[10];
    inputs[31] = modulus[11];
    inputs[32] = modulus[12];
    inputs[33] = modulus[13];
    inputs[34] = modulus[14];
    inputs[35] = modulus[15];

    // Now, do the bigModExp
    int success;
    uint256 inputsLength = inputSize; // No constants in assembly
    // solium-disable-next-line zeppelin/no-arithmetic-operations
    uint256 olen = 32 * modulus.length;
    uint256 gcost = gasCost;
    assembly{
      result := mload(0x40)  // Store result at start of free memory
      mstore(0x40, add(result, olen))  // Move free-memory pointer beyond result
      success := staticcall(
        gcost,
        0x05,                  // BigModExp contract address
        inputs, inputsLength,  // Input segment
        result, olen)          // Output segment
    }
    if (success == 0) {revert("bigModExp call failed");}
  }
}
