// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {ILogAutomation, Log} from "../interfaces/ILogAutomation.sol";
import "../interfaces/StreamsLookupCompatibleInterface.sol";
import {ArbSys} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";

interface IVerifierProxy {
  /**
   * @notice Verifies that the data encoded has been signed
   * correctly by routing to the correct verifier.
   * @param signedReport The encoded data to be verified.
   * @return verifierResponse The encoded response from the verifier.
   */
  function verify(bytes memory signedReport) external returns (bytes memory verifierResponse);
}

contract LogTriggeredStreamsLookup is ILogAutomation, StreamsLookupCompatibleInterface {
  event PerformingLogTriggerUpkeep(
    address indexed from,
    uint256 orderId,
    uint256 amount,
    address exchange,
    uint256 blockNumber,
    bytes blob,
    bytes verified
  );
  event LimitOrderExecuted(uint256 indexed orderId, uint256 indexed amount, address indexed exchange); // keccak(LimitOrderExecuted(uint256,uint256,address)) => 0xd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd
  event IgnoringErrorHandlerData();

  ArbSys internal constant ARB_SYS = ArbSys(0x0000000000000000000000000000000000000064);
  IVerifierProxy internal constant VERIFIER = IVerifierProxy(0x09DFf56A4fF44e0f4436260A04F5CFa65636A481);

  // for log trigger
  bytes32 constant sentSig = 0x3e9c37b3143f2eb7e9a2a0f8091b6de097b62efcfe48e1f68847a832e521750a;
  bytes32 constant withdrawnSig = 0x0a71b8ed921ff64d49e4d39449f8a21094f38a0aeae489c3051aedd63f2c229f;
  bytes32 constant executedSig = 0xd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd;

  // for mercury config
  bool public useArbitrumBlockNum;
  bool public verify;
  string[] public feedsHex = ["0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"];
  string public feedParamKey = "feedIdHex";
  string public timeParamKey = "blockNumber";
  uint256 public counter;
  bool public checkErrReturnBool;

  constructor(bool _useArbitrumBlockNum, bool _verify, bool _checkErrReturnBool) {
    useArbitrumBlockNum = _useArbitrumBlockNum;
    verify = _verify;
    checkErrReturnBool = _checkErrReturnBool;
    counter = 0;
  }

  function start() public {
    // need an initial event to begin the cycle
    emit LimitOrderExecuted(1, 100, address(0x0));
  }

  function setTimeParamKey(string memory timeParam) external {
    timeParamKey = timeParam;
  }

  function setFeedParamKey(string memory feedParam) external {
    feedParamKey = feedParam;
  }

  function setFeedsHex(string[] memory newFeeds) external {
    feedsHex = newFeeds;
  }

  function checkLog(
    Log calldata log,
    bytes memory
  ) external override returns (bool upkeepNeeded, bytes memory performData) {
    uint256 blockNum = getBlockNumber();

    // filter by event signature
    if (log.topics[0] == executedSig) {
      // filter by indexed parameters
      bytes memory t1 = abi.encodePacked(log.topics[1]); // bytes32 to bytes
      uint256 orderId = abi.decode(t1, (uint256));
      bytes memory t2 = abi.encodePacked(log.topics[2]);
      uint256 amount = abi.decode(t2, (uint256));
      bytes memory t3 = abi.encodePacked(log.topics[3]);
      address exchange = abi.decode(t3, (address));

      revert StreamsLookup(
        feedParamKey,
        feedsHex,
        timeParamKey,
        blockNum,
        abi.encode(orderId, amount, exchange, executedSig)
      );
    }
    revert("could not find matching event sig");
  }

  function performUpkeep(bytes calldata performData) external override {
    if (performData.length == 0) {
      emit IgnoringErrorHandlerData();
      return;
    }
    (bytes[] memory values, bytes memory extraData) = abi.decode(performData, (bytes[], bytes));
    (uint256 orderId, uint256 amount, address exchange, bytes32 logTopic0) = abi.decode(
      extraData,
      (uint256, uint256, address, bytes32)
    );

    bytes memory verifiedResponse = "";
    if (verify) {
      verifiedResponse = VERIFIER.verify(values[0]);
    }

    counter = counter + 1;
    if (logTopic0 == executedSig) {
      emit LimitOrderExecuted(1, 100, address(0x0));
    }

    emit PerformingLogTriggerUpkeep(
      tx.origin,
      orderId,
      amount,
      exchange,
      getBlockNumber(),
      values[0],
      verifiedResponse
    );
  }

  function checkCallback(
    bytes[] memory values,
    bytes memory extraData
  ) external view override returns (bool, bytes memory) {
    // do sth about the chainlinkBlob data in values and extraData
    bytes memory performData = abi.encode(values, extraData);
    return (true, performData);
  }

  function checkErrorHandler(
    uint256 errCode,
    bytes memory extraData
  ) external view override returns (bool upkeepNeeded, bytes memory performData) {
    // dummy function with default values
    return (checkErrReturnBool, new bytes(0));
  }

  function getBlockNumber() internal view returns (uint256) {
    if (useArbitrumBlockNum) {
      return ARB_SYS.arbBlockNumber();
    } else {
      return block.number;
    }
  }
}
