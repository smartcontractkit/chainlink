pragma solidity 0.8.16;

import "../automation/interfaces/AutomationCompatibleInterface.sol";
import "../automation/interfaces/StreamsLookupCompatibleInterface.sol";
import {ArbSys} from "../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";

interface IVerifierProxy {
  /**
   * @notice Verifies that the data encoded has been signed
   * correctly by routing to the correct verifier.
   * @param signedReport The encoded data to be verified.
   * @return verifierResponse The encoded response from the verifier.
   */
  function verify(bytes memory signedReport) external returns (bytes memory verifierResponse);
}

contract StreamsLookupUpkeep is AutomationCompatibleInterface, StreamsLookupCompatibleInterface {
  event MercuryPerformEvent(address indexed sender, uint256 indexed blockNumber, bytes v0, bytes verifiedV0, bytes ed);

  ArbSys internal constant ARB_SYS = ArbSys(0x0000000000000000000000000000000000000064);
  // keep these in sync with verifier proxy in RDD
  IVerifierProxy public production_testnet_verifier_proxy = IVerifierProxy(0x2ff010DEbC1297f19579B4246cad07bd24F2488A);
  IVerifierProxy public staging_testnet_verifier_proxy = IVerifierProxy(0x2ff010DEbC1297f19579B4246cad07bd24F2488A);

  uint256 public testRange;
  uint256 public interval;
  uint256 public previousPerformBlock;
  uint256 public initialBlock;
  uint256 public counter;
  string[] public feeds;
  string public feedParamKey;
  string public timeParamKey;
  bool public immutable useArbBlock;
  bool public staging;
  bool public verify;
  bool public shouldRevertCallback;
  bool public callbackReturnBool;
  uint256 public verifyNthReport;

  // find related info here: https://docs.chain.link/data-streams/stream-ids?network=arbitrum&page=1
  constructor(uint256 _testRange, uint256 _interval, bool _useArbBlock, bool _staging, bool _verify) {
    testRange = _testRange;
    interval = _interval;
    previousPerformBlock = 0;
    initialBlock = 0;
    counter = 0;
    useArbBlock = _useArbBlock;
    feedParamKey = "feedIDs"; // feedIDs for v0.3
    timeParamKey = "timestamp"; // timestamp
    // search feeds in notion: "Schema and Feed ID Registry"
    feeds = [
      "0x00037da06d56d083fe599397a4769a042d63aa73dc4ef57709d31e9971a5b439", // BTC / USD in testnet v0.3
      "0x000359843a543ee2fe414dc14c7e7920ef10f4372990b79d6361cdc0dd1ba782" // ETH / USD
    ];
    staging = _staging;
    verify = _verify;
    callbackReturnBool = true;
    verifyNthReport = 0;
  }

  function setParamKeys(string memory _feedParamKey, string memory _timeParamKey) external {
    feedParamKey = _feedParamKey;
    timeParamKey = _timeParamKey;
  }

  function setFeeds(string[] memory _feeds) external {
    feeds = _feeds;
  }

  function setShouldRevertCallback(bool value) public {
    shouldRevertCallback = value;
  }

  function setCallbackReturnBool(bool value) public {
    callbackReturnBool = value;
  }

  function reset() external {
    previousPerformBlock = 0;
    initialBlock = 0;
    counter = 0;
  }

  function checkCallback(bytes[] memory values, bytes memory extraData) external view returns (bool, bytes memory) {
    require(!shouldRevertCallback, "shouldRevertCallback is true");
    // do sth about the chainlinkBlob data in values and extraData
    bytes memory performData = abi.encode(values, extraData);
    return (callbackReturnBool, performData);
  }

  function checkErrorHandler(
    uint256 errCode,
    bytes memory extraData
  ) external view override returns (bool upkeepNeeded, bytes memory performData) {
    // dummy function with default values
    return (false, new bytes(0));
  }

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    if (!eligible()) {
      return (false, data);
    }
    uint256 timeParam = block.timestamp;

    // encode ARB_SYS as extraData to verify that it is provided to checkCallback correctly.
    // in reality, this can be any data or empty
    revert StreamsLookup(feedParamKey, feeds, timeParamKey, timeParam, abi.encodePacked(address(ARB_SYS)));
  }

  function performUpkeep(bytes calldata performData) external {
    uint256 blockNumber;
    if (useArbBlock) {
      blockNumber = ARB_SYS.arbBlockNumber();
    } else {
      blockNumber = block.number;
    }
    if (initialBlock == 0) {
      initialBlock = blockNumber;
    }
    (bytes[] memory values, bytes memory extraData) = abi.decode(performData, (bytes[], bytes));
    previousPerformBlock = blockNumber;
    counter = counter + 1;

    bytes memory v0 = "";
    bytes memory v1 = "";
    if (verify) {
      if (staging) {
        v0 = staging_testnet_verifier_proxy.verify(values[verifyNthReport]);
      } else {
        v0 = production_testnet_verifier_proxy.verify(values[verifyNthReport]);
      }
    }
    emit MercuryPerformEvent(msg.sender, blockNumber, values[verifyNthReport], v0, extraData);
  }

  function setProductionTestnetVerifierProxy(IVerifierProxy proxy) external {
    production_testnet_verifier_proxy = proxy;
  }

  function setStagingTestnetVerifierProxy(IVerifierProxy proxy) external {
    staging_testnet_verifier_proxy = proxy;
  }

  function setStaging(bool _staging) external {
    staging = _staging;
  }

  function setVerifyNthReport(uint256 _n) external {
    verifyNthReport = _n;
  }

  function eligible() public view returns (bool) {
    if (initialBlock == 0) {
      return true;
    }

    uint256 blockNumber;
    if (useArbBlock) {
      blockNumber = ARB_SYS.arbBlockNumber();
    } else {
      blockNumber = block.number;
    }
    return (blockNumber - initialBlock) < testRange && (blockNumber - previousPerformBlock) >= interval;
  }
}
