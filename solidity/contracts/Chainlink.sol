pragma solidity ^0.4.18;

import "solidity-stringutils/strings.sol";

library Chainlink {
  using strings for *;

  struct Run {
    string payload;
    bytes32 jobId;
    address callbackAddress;
    bytes4 callbackFunctionId;
  }

  function add(Run self, string _key, string _value) internal {
    self.payload = addKey(self, _key)
      .concat('":"'.toSlice()).toSlice()
      .concat(_value.toSlice()).toSlice()
      .concat('",'.toSlice());
  }

  function add(Run self, string _key, string[] _values) internal {
    strings.slice memory payload = addKey(self, _key)
      .concat('":["'.toSlice()).toSlice();
    for(uint256 i=0;i<_values.length-1;i++) {
      payload = payload.concat(_values[i].toSlice()).toSlice()
        .concat('","'.toSlice()).toSlice();
    }

    self.payload = payload.concat(_values[_values.length-1].toSlice())
      .toSlice().concat('"],'.toSlice());
  }

  function close(Run self) internal returns (string) {
    var slice = self.payload.toSlice();
    slice._len -= 1;
    return "{".toSlice()
      .concat(slice).toSlice()
      .concat("}".toSlice());
  }

  function addKey(Run run, string _key) private returns (strings.slice) {
    return run.payload.toSlice()
      .concat('"'.toSlice()).toSlice()
      .concat(_key.toSlice()).toSlice();
  }

}
