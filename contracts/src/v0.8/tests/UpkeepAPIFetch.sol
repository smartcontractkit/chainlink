pragma solidity 0.8.6;

contract UpkeepAPIFetch {
  event PerformingUpkeep(
    address indexed from,
    uint256 initialBlock,
    uint256 lastBlock,
    uint256 previousBlock,
    uint256 counter,
    string gender,
    string name
  );

  error ChainlinkAPIFetch(string url, bytes extraData, string[] jsonFields, bytes4 callbackSelector);

  bytes4 internal constant CALLBACK_SELECTOR = UpkeepAPIFetch.callback.selector;

  uint256 public testRange;
  uint256 public interval;
  uint256 public lastBlock;
  uint256 public previousPerformBlock;
  uint256 public initialBlock;
  uint256 public counter;
  string public url;
  string[] public fields;

  constructor(uint256 _testRange, uint256 _interval) {
    testRange = _testRange;
    interval = _interval;
    previousPerformBlock = 0;
    lastBlock = block.number;
    initialBlock = 0;
    counter = 0;
    fields = ["gender", "name"];
    url = "https://api.genderize.io/?name=chris";
  }

  function callback(
    bytes calldata extraData,
    string[] calldata values,
    uint256 statusCode
  ) external view returns (bool, bytes memory) {
    //    if (statusCode > 299) {
    //      // could also pass true here with statusCode so performUpkeep could trigger changes when a url sees an error
    //      return (false, abi.encode("error", statusCode));
    //    }
    string memory gender = values[0];
    string memory name = values[1];
    return (true, abi.encode(gender, name));
  }

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    revert ChainlinkAPIFetch(url, abi.encode(data), fields, CALLBACK_SELECTOR);
  }

  function performUpkeep(bytes calldata performData) external {
    if (initialBlock == 0) {
      initialBlock = block.number;
    }
    counter = counter + 1;
    (string memory gender, string memory name) = abi.decode(performData, (string, string));
    emit PerformingUpkeep(tx.origin, initialBlock, lastBlock, previousPerformBlock, counter, gender, name);
    previousPerformBlock = lastBlock;
  }

  function eligible() public view returns (bool) {
    if (initialBlock == 0) {
      return true;
    }

    return (block.number - initialBlock) < testRange && (block.number - lastBlock) >= interval;
  }

  function setConfig(uint256 _testRange, uint256 _interval) external {
    testRange = _testRange;
    interval = _interval;
    initialBlock = 0;
    counter = 0;
  }

  function setURLs(string memory input) external {
    url = input;
  }

  function stringToUint(string memory s) public pure returns (uint256) {
    bytes memory b = bytes(s);
    uint256 result = 0;
    for (uint256 i = 0; i < b.length; i++) {
      uint256 c = uint256(uint8(b[i]));
      if (c >= 48 && c <= 57) {
        result = result * 10 + (c - 48);
      }
    }
    return result;
  }
}
