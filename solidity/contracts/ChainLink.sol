pragma solidity ^0.4.18;

import "solidity-stringutils/strings.sol";

library ChainLink {
  using strings for *;

  struct Run {
    string payload;
    bytes32 jobId;
    address receiver;
    bytes4 functionHash;
  }

  function add(Run self, string _key, string _value) internal {
    self.payload = self.payload.toSlice().concat('"'.toSlice());
    self.payload = self.payload.toSlice().concat(_key.toSlice());
    self.payload = self.payload.toSlice().concat('":"'.toSlice());
    self.payload = self.payload.toSlice().concat(_value.toSlice());
    self.payload = self.payload.toSlice().concat('",'.toSlice());
  }

  function close(Run self) internal returns (string) {
    var slice = self.payload.toSlice();
    slice = "{".toSlice().concat(slice).toSlice();
    slice._len -= 1;
    return slice.concat("}".toSlice());
  }

}
