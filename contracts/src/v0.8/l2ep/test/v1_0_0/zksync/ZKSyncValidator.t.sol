// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {MockBridgehub} from "../../mocks/zksync/MockZKSyncL1Bridge.sol";
import {ISequencerUptimeFeed} from "../../../interfaces/ISequencerUptimeFeed.sol";
import {ZKSyncValidator} from "../../../zksync/ZKSyncValidator.sol";
import {BaseValidator} from "../../../base/BaseValidator.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ZKSyncValidator_Setup is L2EPTest {
  address internal immutable L2_SEQ_STATUS_RECORDER_ADDRESS = makeAddr("L2_SEQ_STATUS_RECORDER_ADDRESS");
  address internal immutable DUMMY_L1_XDOMAIN_MSNGR_ADDR = makeAddr("DUMMY_L1_XDOMAIN_MSNGR_ADDR");
  address internal immutable DUMMY_L2_UPTIME_FEED_ADDR = makeAddr("DUMMY_L2_UPTIME_FEED_ADDR");
  uint32 internal constant INIT_GAS_PER_PUBDATA_BYTE_LIMIT = 800;
  uint32 internal constant INIT_GAS_LIMIT = 1900000;
  uint32 internal constant MAIN_NET_CHAIN_ID = 300;
  uint32 internal constant BAD_CHAIN_ID = 0;

  ISequencerUptimeFeed internal s_zksyncSequencerUptimeFeed;
  MockBridgehub internal s_mockZKSyncL1Bridge;
  ZKSyncValidator internal s_zksyncValidator;

  /// Fake event that will get emitted when `requestL2TransactionDirect` is called
  /// Definition is taken from MockZKSyncL1Bridge
  event SentMessage(address indexed sender, bytes message);

  /// Setup
  function setUp() public {
    s_mockZKSyncL1Bridge = new MockBridgehub();

    s_zksyncValidator = new ZKSyncValidator(
      address(s_mockZKSyncL1Bridge),
      DUMMY_L2_UPTIME_FEED_ADDR,
      INIT_GAS_LIMIT,
      MAIN_NET_CHAIN_ID,
      INIT_GAS_PER_PUBDATA_BYTE_LIMIT
    );
  }
}

contract ZKSyncValidator_Constructor is ZKSyncValidator_Setup {
  /// @notice Reverts when chain ID is invalid
  function test_Constructor_RevertWhen_ChainIdIsInvalid() public {
    vm.expectRevert(ZKSyncValidator.InvalidChainID.selector);
    new ZKSyncValidator(
      DUMMY_L1_XDOMAIN_MSNGR_ADDR,
      DUMMY_L2_UPTIME_FEED_ADDR,
      INIT_GAS_LIMIT,
      BAD_CHAIN_ID,
      INIT_GAS_PER_PUBDATA_BYTE_LIMIT
    );
  }

  /// @notice Reverts when L1 bridge address is zero
  function test_Constructor_RevertWhen_L1BridgeAddressIsZero() public {
    vm.expectRevert(BaseValidator.L1CrossDomainMessengerAddressZero.selector);
    new ZKSyncValidator(
      address(0),
      DUMMY_L2_UPTIME_FEED_ADDR,
      INIT_GAS_LIMIT,
      MAIN_NET_CHAIN_ID,
      INIT_GAS_PER_PUBDATA_BYTE_LIMIT
    );
  }

  /// @notice Reverts when L2 update feed address is zero
  function test_Constructor_RevertWhen_L2UpdateFeedAddressIsZero() public {
    vm.expectRevert(BaseValidator.L2UptimeFeedAddrZero.selector);
    new ZKSyncValidator(
      DUMMY_L1_XDOMAIN_MSNGR_ADDR,
      address(0),
      INIT_GAS_LIMIT,
      MAIN_NET_CHAIN_ID,
      INIT_GAS_PER_PUBDATA_BYTE_LIMIT
    );
  }
}

contract ZKSyncValidator_GetSetL2GasPerPubdataByteLimit is ZKSyncValidator_Setup {
  /// @notice Correctly gets and updates the gas per pubdata byte limit
  function test_GetSetL2GasPerPubdataByteLimit_CorrectlyHandlesGasPerPubdataByteLimit() public {
    assertEq(s_zksyncValidator.getL2GasPerPubdataByteLimit(), INIT_GAS_PER_PUBDATA_BYTE_LIMIT);

    uint32 newGasPerPubDataByteLimit = 2000000;
    s_zksyncValidator.setL2GasPerPubdataByteLimit(newGasPerPubDataByteLimit);
    assertEq(s_zksyncValidator.getL2GasPerPubdataByteLimit(), newGasPerPubDataByteLimit);
  }
}

contract ZKSyncValidator_GetChainId is ZKSyncValidator_Setup {
  /// @notice Correctly gets the chain ID
  function test_GetChainId_CorrectlyGetsTheChainId() public view {
    assertEq(s_zksyncValidator.getChainId(), MAIN_NET_CHAIN_ID);
  }
}

contract ZKSyncValidator_Validate is ZKSyncValidator_Setup {
  /// @notice Reverts if called by an account with no access
  function test_Validate_RevertWhen_CalledByAccountWithNoAccess() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("No access");
    s_zksyncValidator.validate(0, 0, 1, 1);
  }

  /// @notice Posts sequencer status when there is no status change
  function test_Validate_PostSequencerStatus_NoStatusChange() public {
    // Gives access to the s_eoaValidator
    s_zksyncValidator.addAccess(s_eoaValidator);

    // Sets block.timestamp to a later date
    uint256 futureTimestampInSeconds = block.timestamp + 5000;
    vm.startPrank(s_eoaValidator);
    vm.warp(futureTimestampInSeconds);

    // Sets up the expected event data
    bytes memory message = abi.encodeWithSelector(
      ISequencerUptimeFeed.updateStatus.selector,
      false,
      futureTimestampInSeconds
    );

    vm.expectEmit(false, false, false, true);
    emit SentMessage(address(s_zksyncValidator), message);

    // Runs the function (which produces the event to test)
    s_zksyncValidator.validate(0, 0, 0, 0);
  }

  /// @notice Posts sequencer offline status
  function test_Validate_PostSequencerOffline() public {
    // Gives access to the s_eoaValidator
    s_zksyncValidator.addAccess(s_eoaValidator);

    // Sets block.timestamp to a later date
    uint256 futureTimestampInSeconds = block.timestamp + 10000;
    vm.startPrank(s_eoaValidator);
    vm.warp(futureTimestampInSeconds);

    // Sets up the expected event data
    vm.expectEmit(false, false, false, true);
    emit SentMessage(
      address(s_zksyncValidator),
      abi.encodeWithSelector(ISequencerUptimeFeed.updateStatus.selector, true, futureTimestampInSeconds)
    );

    // Runs the function (which produces the event to test)
    s_zksyncValidator.validate(0, 0, 1, 1);
  }
}
