pragma solidity 0.8.15;

import "../interfaces/automation/AutomationCompatibleInterface.sol";
import "../dev/interfaces/automation/MercuryLookupCompatibleInterface.sol";
import {ArbSys} from "../dev/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";

//interface IVerifierProxy {
//  /**
//   * @notice Verifies that the data encoded has been signed
//   * correctly by routing to the correct verifier.
//   * @param signedReport The encoded data to be verified.
//   * @return verifierResponse The encoded response from the verifier.
//   */
//  function verify(bytes memory signedReport) external returns (bytes memory verifierResponse);
//}

contract MercuryUpkeep is AutomationCompatibleInterface, MercuryLookupCompatibleInterface {
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
  string public feedLabel;
  string public queryLabel;
  bool public immutable integrationTest;

  constructor(uint256 _testRange, uint256 _interval, bool _integrationTest) {
    testRange = _testRange;
    interval = _interval;
    previousPerformBlock = 0;
    initialBlock = 0;
    counter = 0;
    feedLabel = "feedIDStr"; // or feedIDHex
    feeds = ["ETH-USD-ARBITRUM-TESTNET", "BTC-USD-ARBITRUM-TESTNET"];
    queryLabel = "blockNumber"; // timestamp not supported yet
    integrationTest = _integrationTest;
  }

  function mercuryCallback(bytes[] memory values, bytes memory extraData) external pure returns (bool, bytes memory) {
    // do sth about the chainlinkBlob data in values and extraData
    bytes memory performData = abi.encode(values, extraData);
    return (true, performData);
  }

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    if (!eligible()) {
      return (false, data);
    }
    uint256 blockNumber;
    if (integrationTest) {
      blockNumber = block.number;
    } else {
      blockNumber = ARB_SYS.arbBlockNumber();
    }
    // encode ARB_SYS as extraData to verify that it is provided to mercuryCallback correctly.
    // in reality, this can be any data or empty
    revert MercuryLookup(feedLabel, feeds, queryLabel, blockNumber, abi.encodePacked(address(ARB_SYS)));
  }

  function performUpkeep(bytes calldata performData) external {
    uint256 blockNumber;
    if (integrationTest) {
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
    if (integrationTest) {
      blockNumber = block.number;
    } else {
      blockNumber = ARB_SYS.arbBlockNumber();
    }
    return (blockNumber - initialBlock) < testRange && (blockNumber - previousPerformBlock) >= interval;
  }
}
