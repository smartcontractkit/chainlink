pragma solidity 0.8.15;

import "@openzeppelin/contracts/utils/Strings.sol";

contract UpkeepAPIFetch {
  event PerformingUpkeep(address indexed from, uint256 lastBlock, uint256 counter, string id, string name);

  error ChainlinkAPIFetch(string url, bytes extraData, string[] jsonFields, bytes4 callbackSelector);

  bytes4 internal constant CALLBACK_SELECTOR = UpkeepAPIFetch.callback.selector;

  uint256 public testRange;
  uint256 public interval;
  uint256 public lastBlock;
  uint256 public previousPerformBlock;
  uint256 public initialBlock;
  uint256 public counter;
  string public url;
  string public id;
  string public pokemon;
  string[] public fields;

  constructor(uint256 _testRange, uint256 _interval) {
    testRange = _testRange;
    interval = _interval;
    previousPerformBlock = 0;
    lastBlock = block.number;
    initialBlock = 0;
    counter = 0;
    fields = ["id", "name"];
    url = "https://pokeapi.co/api/v2/pokemon/";
  }

  function callback(
    bytes calldata extraData,
    string[] calldata values,
    uint256 statusCode
  ) external view returns (bool, bytes memory) {
    if (statusCode > 299) {
      // pass true here with msg to perform to trigger changes when a url sees an error
      return (true, abi.encode("error", "error"));
    }
    string memory pid = values[0];
    string memory name = values[1];
    return (true, abi.encode(pid, name));
  }

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    if (!eligible()) {
      return (false, data);
    }
    string memory pid = Strings.toString(counter + 1);
    string memory urlWithId = string(abi.encodePacked(url, pid));
    revert ChainlinkAPIFetch(urlWithId, abi.encode(data), fields, CALLBACK_SELECTOR);
  }

  function performUpkeep(bytes calldata performData) external {
    if (initialBlock == 0) {
      initialBlock = block.number;
    }
    lastBlock = block.number;
    counter = counter + 1;
    (string memory pid, string memory name) = abi.decode(performData, (string, string));
    if (keccak256(abi.encodePacked(pid)) == keccak256(abi.encodePacked("error"))) {
      counter = 0;
    } else {
      id = pid;
      pokemon = name;
    }
    emit PerformingUpkeep(tx.origin, lastBlock, counter, pid, name);
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
}
