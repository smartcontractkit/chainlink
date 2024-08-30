// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {MockZKSyncL1Bridge} from "../../mocks/zksync/MockZKSyncL1Bridge.sol";
import {ZKSyncSequencerUptimeFeedInterface} from "../../../dev/interfaces/ZKSyncSequencerUptimeFeedInterface.sol";
import {ZKSyncValidator} from "../../../dev/zksync/ZKSyncValidator.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ZKSyncValidatorTest is L2EPTest {
  address internal constant L2_SEQ_STATUS_RECORDER_ADDRESS = 0x491B1dDA0A8fa069bbC1125133A975BF4e85a91b;
  uint32 internal constant INIT_GAS_LIMIT = 1900000;
  uint32 internal constant MAIN_NET_CHAIN_ID = 300;
  uint32 internal constant INIT_GAS_PER_PUBDATA_BYTE_LIMIT = 800;

  MockZKSyncL1Bridge internal s_mockZKSyncL1Bridge;

  /// Setup
  function setUp() public {
    s_mockZKSyncL1Bridge = new MockZKSyncL1Bridge();
  }

  function getValidValidator() internal returns (ZKSyncValidator) {
    return
      buildValidator(
        address(s_mockZKSyncL1Bridge),
        address(0xa04Fc18f012B1a5A8231c7Ee4b916Dd6dbd271b6),
        INIT_GAS_LIMIT,
        MAIN_NET_CHAIN_ID,
        INIT_GAS_PER_PUBDATA_BYTE_LIMIT
      );
  }

  function buildValidator(
    address l1CrossDomainMessengerAddress,
    address l2UptimeFeedAddr,
    uint32 gasLimit,
    uint32 chainId,
    uint32 l2GasPerPubdataByteLimit
  ) internal returns (ZKSyncValidator) {
    return
      new ZKSyncValidator(l1CrossDomainMessengerAddress, l2UptimeFeedAddr, gasLimit, chainId, l2GasPerPubdataByteLimit);
  }
}

contract ZKSyncValidator_Constructor is ZKSyncValidatorTest {
  /// @notice it correctly validates that the chain id is valid
  function test_ConstructingRevertedWithInvalidChainId() public {
    vm.expectRevert(ZKSyncValidator.InvalidChainID.selector);
    buildValidator(
      address(0xFe31891940A2e5f04B76eD8bD1038E44127d1512),
      address(0xa04Fc18f012B1a5A8231c7Ee4b916Dd6dbd271b6),
      INIT_GAS_LIMIT,
      0,
      INIT_GAS_PER_PUBDATA_BYTE_LIMIT
    );
  }

  /// @notice it correctly validates that the L1 bridge address is not zero
  function test_ConstructingRevertedWithZeroL1BridgeAddress() public {
    vm.expectRevert(
      abi.encodeWithSelector(ZKSyncValidator.ZeroAddressNotAllowed.selector, "Invalid xDomain Messenger address")
    );
    buildValidator(
      address(0),
      address(0xa04Fc18f012B1a5A8231c7Ee4b916Dd6dbd271b6),
      INIT_GAS_LIMIT,
      MAIN_NET_CHAIN_ID,
      INIT_GAS_PER_PUBDATA_BYTE_LIMIT
    );
  }

  /// @notice it correctly validates that the L2 Uptime feed address is not zero
  function test_ConstructingRevertedWithZeroL2UpdateFeedAddress() public {
    vm.expectRevert(
      abi.encodeWithSelector(ZKSyncValidator.ZeroAddressNotAllowed.selector, "Invalid ZKSyncSequencerUptimeFeedInterface contract address")
    );
    buildValidator(
      address(0xa04Fc18f012B1a5A8231c7Ee4b916Dd6dbd271b6),
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
    ZKSyncValidator validator = getValidValidator();
    assertEq(validator.getL2GasPerPubdataByteLimit(), INIT_GAS_PER_PUBDATA_BYTE_LIMIT);

    uint32 newGasPerPubDataByteLimit = 2000000;
    validator.setL2GasPerPubdataByteLimit(newGasPerPubDataByteLimit);
    assertEq(validator.getL2GasPerPubdataByteLimit(), newGasPerPubDataByteLimit);
  }
}

contract ZKSyncValidator_GetChainId is ZKSyncValidatorTest {
  /// @notice it correctly gets the chain id
  function test_CorrectlyGetsTheChainId() public {
    ZKSyncValidator validator = getValidValidator();
    assertEq(validator.getChainId(), MAIN_NET_CHAIN_ID);
  }
}

contract ZKSyncValidator_Validate is ZKSyncValidatorTest {
  /// @notice it reverts if called by account with no access
  function test_RevertsIfCalledByAnAccountWithNoAccess() public {
    ZKSyncValidator validator = getValidValidator();

    vm.startPrank(s_strangerAddr);
    vm.expectRevert("No access");
    validator.validate(0, 0, 1, 1);
  }

  /// @notice it posts sequencer status when there is not status change
  function test_PostSequencerStatusWhenThereIsNotStatusChange() public {
    ZKSyncValidator validator = getValidValidator();

    // Gives access to the s_eoaValidator
    validator.addAccess(s_eoaValidator);

    // Sets block.timestamp to a later date
    uint256 futureTimestampInSeconds = block.timestamp + 5000;
    vm.startPrank(s_eoaValidator);
    vm.warp(futureTimestampInSeconds);

    // Sets up the expected event data
    bytes memory message = abi.encodeWithSelector(ZKSyncSequencerUptimeFeedInterface.updateStatus.selector, false, futureTimestampInSeconds);
    
    vm.expectEmit(false, false, false, true);
    // emit SentMessage(
    //   address(validator), // sender
    //   L2_SEQ_STATUS_RECORDER_ADDRESS, // target
    //   0, // value
    //   0, // nonce
    //   INIT_GAS_LIMIT, // gas limit
    //   message
    // );

    // Runs the function (which produces the event to test)
    validator.validate(0, 0, 0, 0);
  }

  // /// @notice it post sequencer offline
  // function test_PostSequencerOffline() public {
  //   // Gives access to the s_eoaValidator
  //   s_zksyncValidator.addAccess(s_eoaValidator);

  //   // Sets block.timestamp to a later date
  //   uint256 futureTimestampInSeconds = block.timestamp + 10000;
  //   vm.startPrank(s_eoaValidator);
  //   vm.warp(futureTimestampInSeconds);

  //   // Sets up the expected event data
  //   vm.expectEmit(false, false, false, true);
  //   emit SentMessage(
  //     address(s_zksyncValidator), // sender
  //     L2_SEQ_STATUS_RECORDER_ADDRESS, // target
  //     0, // value
  //     0, // nonce
  //     INIT_GAS_LIMIT, // gas limit
  //     abi.encodeWithSelector(ScrollSequencerUptimeFeed.updateStatus.selector, true, futureTimestampInSeconds) // message
  //   );

  //   // Runs the function (which produces the event to test)
  //   s_zksyncValidator.validate(0, 0, 1, 1);
  // }
}
