pragma solidity ^0.4.18;


import "solidity-stringutils/strings.sol";

library Json {
  using strings for *;

  struct Params {
    string payload;
  }

  function add(Params self, string key, string value) internal {
    self.payload = self.payload.toSlice().concat('"'.toSlice());
    self.payload = self.payload.toSlice().concat(key.toSlice());
    self.payload = self.payload.toSlice().concat('":"'.toSlice());
    self.payload = self.payload.toSlice().concat(value.toSlice());
    self.payload = self.payload.toSlice().concat('",'.toSlice());
  }

  function close(Params self) internal returns (string) {
    var slice = self.payload.toSlice();
    slice = "{".toSlice().concat(slice).toSlice();
    slice._len -= 1;
    return slice.concat("}".toSlice());
  }

}


contract ChainLinked {
  using Json for Json.Params;
}
