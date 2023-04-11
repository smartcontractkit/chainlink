pragma solidity 0.8.15;

contract MercuryUpkeep {
  event MercuryEvent(address indexed from, bytes data);

  error MercuryLookup(string[] feedIDStrList, uint256 instance, bytes extraData);

  uint256 public testRange;
  uint256 public interval;
  uint256 public lastBlock;
  uint256 public previousPerformBlock;
  uint256 public initialBlock;
  uint256 public counter;
  string[] public feeds;

  constructor(uint256 _testRange, uint256 _interval) {
    testRange = _testRange;
    interval = _interval;
    previousPerformBlock = 0;
    lastBlock = block.number;
    initialBlock = 0;
    counter = 0;
    feeds = ["ETH-USD", "BTC-USD"];
  }

  function mercuryCallback(bytes[] memory values, bytes memory extraData) external view returns (bool, bytes memory) {
    //    this is where they do something with the chainlinkBlobHex
    return (true, extraData);
  }

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    if (!eligible()) {
      return (false, data);
    }
    // block.number or block.timestamp depending on user
    revert MercuryLookup(feeds, block.number, data);
  }

  function performUpkeep(bytes calldata performData) external {
    if (initialBlock == 0) {
      initialBlock = block.number;
    }
    lastBlock = block.number;
    counter = counter + 1;
    emit MercuryEvent(tx.origin, performData);
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

  function setFeeds(string[] memory input) external {
    feeds = input;
  }
}
