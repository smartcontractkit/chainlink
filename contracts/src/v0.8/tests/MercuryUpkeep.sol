pragma solidity 0.8.16;

import "../automation/interfaces/AutomationCompatibleInterface.sol";
import "../dev/automation/v2_1/interfaces/FeedLookupCompatibleInterface.sol";
import {ArbSys} from "../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";

//interface IVerifierProxy {
//  /**
//   * @notice Verifies that the data encoded has been signed
//   * correctly by routing to the correct verifier.
//   * @param signedReport The encoded data to be verified.
//   * @return verifierResponse The encoded response from the verifier.
//   */
//  function verify(bytes memory signedReport) external returns (bytes memory verifierResponse);
//}

contract MercuryUpkeep is AutomationCompatibleInterface, FeedLookupCompatibleInterface {
  event MercuryPerformEvent(
    address indexed origin,
    address indexed sender,
    uint256 indexed blockNumber,
    bytes v0,
    bytes v1,
    bytes ed
  );

  ArbSys internal constant ARB_SYS = ArbSys(0x0000000000000000000000000000000000000064);
  //  IVerifierProxy internal constant VERIFIER = IVerifierProxy(0xa4D813064dc6E2eFfaCe02a060324626d4C5667f);

  uint256 public testRange;
  uint256 public interval;
  uint256 public previousPerformBlock;
  uint256 public initialBlock;
  uint256 public counter;
  string[] public feeds;
  string public feedParamKey;
  string public timeParamKey;
  bool public immutable useL1BlockNumber;
  bool public shouldRevertCallback;
  bool public callbackReturnBool;

  constructor(uint256 _testRange, uint256 _interval, bool _useL1BlockNumber) {
    testRange = _testRange;
    interval = _interval;
    previousPerformBlock = 0;
    initialBlock = 0;
    counter = 0;
    feedParamKey = "feedIdHex"; // feedIDStr is deprecated
    feeds = [
      "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
      "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"
    ];
    timeParamKey = "blockNumber"; // timestamp not supported yet
    useL1BlockNumber = _useL1BlockNumber;
    callbackReturnBool = true;
  }

  function setShouldRevertCallback(bool value) public {
    shouldRevertCallback = value;
  }

  function setCallbackReturnBool(bool value) public {
    callbackReturnBool = value;
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
    uint256 blockNumber;
    if (useL1BlockNumber) {
      blockNumber = block.number;
    } else {
      blockNumber = ARB_SYS.arbBlockNumber();
    }
    // encode ARB_SYS as extraData to verify that it is provided to checkCallback correctly.
    // in reality, this can be any data or empty
    revert FeedLookup(feedParamKey, feeds, timeParamKey, blockNumber, abi.encodePacked(address(ARB_SYS)));
  }

  function performUpkeep(bytes calldata performData) external {
    uint256 blockNumber;
    if (useL1BlockNumber) {
      blockNumber = block.number;
    } else {
      blockNumber = ARB_SYS.arbBlockNumber();
    }
    if (initialBlock == 0) {
      initialBlock = blockNumber;
    }
    (bytes[] memory values, bytes memory extraData) = abi.decode(performData, (bytes[], bytes));
    previousPerformBlock = blockNumber;
    counter = counter + 1;
    //    bytes memory v0 = VERIFIER.verify(values[0]);
    //    bytes memory v1 = VERIFIER.verify(values[1]);
    emit MercuryPerformEvent(tx.origin, msg.sender, blockNumber, values[0], values[1], extraData);
  }

  function eligible() public view returns (bool) {
    if (initialBlock == 0) {
      return true;
    }

    uint256 blockNumber;
    if (useL1BlockNumber) {
      blockNumber = block.number;
    } else {
      blockNumber = ARB_SYS.arbBlockNumber();
    }
    return (blockNumber - initialBlock) < testRange && (blockNumber - previousPerformBlock) >= interval;
  }
}
