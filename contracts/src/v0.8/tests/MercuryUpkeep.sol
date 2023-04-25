pragma solidity 0.8.15;

import "../interfaces/automation/AutomationCompatibleInterface.sol";
import "../dev/interfaces/automation/MercuryLookupCompatibleInterface.sol";

contract MercuryUpkeep is AutomationCompatibleInterface, MercuryLookupCompatibleInterface {
  event MercuryEvent(address indexed from, bytes data);

  error MercuryLookup(string feedLabel, string[] feedList, string queryLabel, uint256 query, bytes extraData);

  uint256 public testRange;
  uint256 public interval;
  uint256 public lastBlock;
  uint256 public previousPerformBlock;
  uint256 public initialBlock;
  uint256 public counter;
  string[] public feeds;
  string public feedLabel;
  string public queryLabel;

  constructor(uint256 _testRange, uint256 _interval) {
    testRange = _testRange;
    interval = _interval;
    previousPerformBlock = 0;
    lastBlock = block.number;
    initialBlock = 0;
    counter = 0;
    feedLabel = "feedIDStr"; // or feedIDHex
    feeds = ["ETH-USD-ARBITRUM-TESTNET", "BTC-USD-ARBITRUM-TESTNET"]; // or ["0x4554482d5553442d415242495452554d2d544553544e45540000000000000000","0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"]
    // blockNumber is used for special customers, timestamp is an option for all customers but it's under development
    // the advantage to use timestamp is that this contract can be deployed at a different chain than where Mercury server lives
    // when using blockNumber, the server probably will fail the request bc blockNumber and timestamp are too far apart
    queryLabel = "blockNumber";
  }

  // make a mercury compatible interface + automation compatible interface => generate ABI
  // values => chainlinkblobs
  function mercuryCallback(bytes[] memory values, bytes memory extraData) external view returns (bool, bytes memory) {
    //    this is where they do something with the chainlinkBlobHex
    return (true, extraData);
  }

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    if (!eligible()) {
      return (false, data);
    }
    revert MercuryLookup(feedLabel, feeds, queryLabel, block.number, data);
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
