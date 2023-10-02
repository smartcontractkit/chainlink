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
  IVerifierProxy internal constant PRODUCTION_TESTNET_VERIFIER_PROXY =
    IVerifierProxy(0x09DFf56A4fF44e0f4436260A04F5CFa65636A481);
  IVerifierProxy internal constant STAGING_TESTNET_VERIFIER_PROXY =
    IVerifierProxy(0x60448B880c9f3B501af3f343DA9284148BD7D77C);

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
      //"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", // ETH / USD in production testnet v0.2
      //"0x4254432d5553442d415242495452554d2d544553544e45540000000000000000" // BTC / USD in production testnet v0.2
      "0x00028c915d6af0fd66bba2d0fc9405226bca8d6806333121a7d9832103d1563c" // ETH / USD in staging testnet v0.3
    ];
    staging = _staging;
    verify = _verify;
    callbackReturnBool = true;
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

  function reset() public {
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

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    if (!eligible()) {
      return (false, data);
    }
    uint256 timeParam;
    if (keccak256(abi.encodePacked(feedParamKey)) == keccak256(abi.encodePacked("feedIdHex"))) {
      if (useArbBlock) {
        timeParam = ARB_SYS.arbBlockNumber();
      } else {
        timeParam = block.number;
      }
    } else {
      // assume this will be feedIDs for v0.3
      timeParam = block.timestamp;
    }

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
        v0 = STAGING_TESTNET_VERIFIER_PROXY.verify(values[0]);
      } else {
        v0 = PRODUCTION_TESTNET_VERIFIER_PROXY.verify(values[0]);
      }
    }
    emit MercuryPerformEvent(msg.sender, blockNumber, values[0], v0, extraData);
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
