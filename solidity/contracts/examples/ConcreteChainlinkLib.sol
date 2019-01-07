pragma solidity 0.4.24;

import "../ChainlinkLib.sol";

contract ConcreteChainlinkLib {
  using ChainlinkLib for ChainlinkLib.Run;
  using CBOR for Buffer.buffer;

  ChainlinkLib.Run run;

  event RunData(bytes payload);

  constructor() public {
    ChainlinkLib.Run memory r2 = run;
    Buffer.init(r2.buf, 128);
    run = r2;
  }

  function closeEvent() public {
    ChainlinkLib.Run memory r2 = run;
    r2.close();
    run = r2;
    emit RunData(run.buf.buf);
  }

  function add(string _key, string _value) public {
    ChainlinkLib.Run memory r2 = run;
    r2.add(_key, _value);
    run = r2;
  }

  function addBytes(string _key, bytes _value) public {
    ChainlinkLib.Run memory r2 = run;
    r2.addBytes(_key, _value);
    run = r2;
  }

  function addInt(string _key, int256 _value) public {
    ChainlinkLib.Run memory r2 = run;
    r2.addInt(_key, _value);
    run = r2;
  }

  function addUint(string _key, uint256 _value) public {
    ChainlinkLib.Run memory r2 = run;
    r2.addUint(_key, _value);
    run = r2;
  }

  // Temporarily have method receive bytes32[] memory until experimental
  // string[] memory can be invoked from truffle tests.
  function addStringArray(string _key, bytes32[] memory _values) public {
    string[] memory strings = new string[](_values.length);
    for (uint256 i = 0; i < _values.length; i++) {
      strings[i] = bytes32ToString(_values[i]);
    }
    ChainlinkLib.Run memory r2 = run;
    r2.addStringArray(_key, strings);
    run = r2;
  }

  function bytes32ToString(bytes32 x) private pure returns (string) {
    bytes memory bytesString = new bytes(32);
    uint charCount = 0;
    for (uint j = 0; j < 32; j++) {
      byte char = byte(bytes32(uint(x) * 2 ** (8 * j))); //solium-disable-line
      if (char != 0) {
        bytesString[charCount] = char;
        charCount++;
      }
    }
    bytes memory bytesStringTrimmed = new bytes(charCount);
    for (j = 0; j < charCount; j++) {
      bytesStringTrimmed[j] = bytesString[j];
    }
    return string(bytesStringTrimmed);
  }
}
