// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {MockBridgehub} from "../../mocks/zksync/MockZKSyncL1Bridge.sol";
import {ISequencerUptimeFeed} from "../../../dev/interfaces/ISequencerUptimeFeed.sol";
import {ZKSyncValidator} from "../../../dev/zksync/ZKSyncValidator.sol";
import {BaseValidator} from "../../../dev/shared/BaseValidator.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ZKSyncValidatorTest is L2EPTest {
  address internal constant L2_SEQ_STATUS_RECORDER_ADDRESS = address(0x491B1dDA0A8fa069bbC1125133A975BF4e85a91b);
  address internal constant DUMMY_L1_XDOMAIN_MSNGR_ADDR = address(0xa04Fc18f012B1a5A8231c7Ee4b916Dd6dbd271b6);
  address internal constant DUMMY_L2_UPTIME_FEED_ADDR = address(0xFe31891940A2e5f04B76eD8bD1038E44127d1512);
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

contract ZKSyncValidator_Constructor is ZKSyncValidatorTest {
  /// @notice it correctly validates that the chain id is valid
  function test_ConstructingRevertedWithInvalidChainId() public {
    vm.expectRevert(ZKSyncValidator.InvalidChainID.selector);
    new ZKSyncValidator(
      DUMMY_L1_XDOMAIN_MSNGR_ADDR,
      DUMMY_L2_UPTIME_FEED_ADDR,
      INIT_GAS_LIMIT,
      BAD_CHAIN_ID,
      INIT_GAS_PER_PUBDATA_BYTE_LIMIT
    );
  }

  /// @notice it correctly validates that the L1 bridge address is not zero
  function test_ConstructingRevertedWithZeroL1BridgeAddress() public {
    vm.expectRevert(BaseValidator.L1CrossDomainMessengerAddressZero.selector);
    new ZKSyncValidator(
      address(0),
      DUMMY_L2_UPTIME_FEED_ADDR,
      INIT_GAS_LIMIT,
      MAIN_NET_CHAIN_ID,
      INIT_GAS_PER_PUBDATA_BYTE_LIMIT
    );
  }

  /// @notice it correctly validates that the L2 Uptime feed address is not zero
  function test_ConstructingRevertedWithZeroL2UpdateFeedAddress() public {
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

contract ZKSyncValidator_GetSetL2GasPerPubdataByteLimit is ZKSyncValidatorTest {
  /// @notice it correctly updates the gas limit per pubdata byte
  function test_CorrectlyGetsAndUpdatesTheGasPerPubdataByteLimit() public {
    assertEq(s_zksyncValidator.getL2GasPerPubdataByteLimit(), INIT_GAS_PER_PUBDATA_BYTE_LIMIT);

    uint32 newGasPerPubDataByteLimit = 2000000;
    s_zksyncValidator.setL2GasPerPubdataByteLimit(newGasPerPubDataByteLimit);
    assertEq(s_zksyncValidator.getL2GasPerPubdataByteLimit(), newGasPerPubDataByteLimit);
  }
}

contract ZKSyncValidator_GetChainId is ZKSyncValidatorTest {
  /// @notice it correctly gets the chain id
  function test_CorrectlyGetsTheChainId() public {
    assertEq(s_zksyncValidator.getChainId(), MAIN_NET_CHAIN_ID);
  }
}

contract ZKSyncValidator_Validate is ZKSyncValidatorTest {
  /// @notice it reverts if called by account with no access
  function test_RevertsIfCalledByAnAccountWithNoAccess() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("No access");
    s_zksyncValidator.validate(0, 0, 1, 1);
  }

  /// @notice it posts sequencer status when there is not status change
  function test_PostSequencerStatusWhenThereIsNotStatusChange() public {
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

  /// @notice it post sequencer offline
  function test_PostSequencerOffline() public {
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
