pragma solidity 0.8.6;

contract UpkeepEIP3668 {
  event PerformingUpkeep(
    address indexed from,
    uint256 initialBlock,
    uint256 lastBlock,
    uint256 previousBlock,
    uint256 counter,
    bytes resp
  );

  error OffchainLookup(address sender, string[] urls, bytes callData, bytes4 callbackFunction, bytes extraData);

  bytes4 internal constant CALLBACK_SELECTOR = UpkeepEIP3668.callback.selector;

  uint256 public testRange;
  uint256 public interval;
  uint256 public lastBlock;
  uint256 public previousPerformBlock;
  uint256 public initialBlock;
  uint256 public counter;
  string[] public urls;

  constructor(uint256 _testRange, uint256 _interval) {
    testRange = _testRange;
    interval = _interval;
    previousPerformBlock = 0;
    lastBlock = block.number;
    initialBlock = 0;
    counter = 0;
    urls = ["https://catfact.ninja/fact?q={sender}+{data}", "https://www.google.com/search?q={sender}+{data}"];
  }

  function callback(bytes calldata resp, bytes calldata extra) external view returns (bool, bytes memory) {
    bool extraTrue = abi.decode(extra, (bool));
    return (extraTrue, resp);
  }

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    revert OffchainLookup(address(this), urls, abi.encode(data), CALLBACK_SELECTOR, abi.encode(true));
  }

  function performUpkeep(bytes calldata performData) external {
    if (initialBlock == 0) {
      initialBlock = block.number;
    }
    counter = counter + 1;
    emit PerformingUpkeep(tx.origin, initialBlock, lastBlock, previousPerformBlock, counter, performData);
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

  function setURLs(string[] memory input) external {
    urls = input;
  }
}
