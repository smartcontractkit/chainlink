// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Internal} from "../../../libraries/Internal.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

contract OffRamp_getExecutionState is OffRampSetup {
  mapping(uint64 sourceChainSelector => mapping(uint64 seqNum => Internal.MessageExecutionState state)) internal
    s_differentialExecutionState;

  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 32
  function testFuzz_Differential_Success(
    uint64 sourceChainSelector,
    uint16[500] memory seqNums,
    uint8[500] memory values
  ) public {
    for (uint256 i = 0; i < seqNums.length; ++i) {
      // Only use the first three slots. This makes sure existing slots get overwritten
      // as the tests uses 500 sequence numbers.
      uint16 seqNum = seqNums[i] % 386;
      Internal.MessageExecutionState state = Internal.MessageExecutionState(values[i] % 4);
      s_differentialExecutionState[sourceChainSelector][seqNum] = state;
      s_offRamp.setExecutionStateHelper(sourceChainSelector, seqNum, state);
      assertEq(uint256(state), uint256(s_offRamp.getExecutionState(sourceChainSelector, seqNum)));
    }

    for (uint256 i = 0; i < seqNums.length; ++i) {
      uint16 seqNum = seqNums[i] % 386;
      Internal.MessageExecutionState expectedState = s_differentialExecutionState[sourceChainSelector][seqNum];
      assertEq(uint256(expectedState), uint256(s_offRamp.getExecutionState(sourceChainSelector, seqNum)));
    }
  }

  function test_GetExecutionState() public {
    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 0, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3);

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 1, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (3 << 2));

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 1, Internal.MessageExecutionState.IN_PROGRESS);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 2));

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 2, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 2) + (3 << 4));

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 127, Internal.MessageExecutionState.IN_PROGRESS);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 2) + (3 << 4) + (1 << 254));

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 128, Internal.MessageExecutionState.SUCCESS);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 2) + (3 << 4) + (1 << 254));
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 1), 2);

    assertEq(
      uint256(Internal.MessageExecutionState.FAILURE), uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 0))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.IN_PROGRESS),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 1))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.FAILURE), uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 2))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.IN_PROGRESS),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 127))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.SUCCESS),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 128))
    );
  }

  function test_GetDifferentChainExecutionState() public {
    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 0, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1 + 1, 0), 0);

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 127, Internal.MessageExecutionState.IN_PROGRESS);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 254));
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1 + 1, 0), 0);

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, 128, Internal.MessageExecutionState.SUCCESS);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 254));
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 1), 2);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1 + 1, 0), 0);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1 + 1, 1), 0);

    s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1 + 1, 127, Internal.MessageExecutionState.FAILURE);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 0), 3 + (1 << 254));
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, 1), 2);
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1 + 1, 0), (3 << 254));
    assertEq(s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1 + 1, 1), 0);

    assertEq(
      uint256(Internal.MessageExecutionState.FAILURE), uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 0))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.IN_PROGRESS),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 127))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.SUCCESS),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, 128))
    );

    assertEq(
      uint256(Internal.MessageExecutionState.UNTOUCHED),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1 + 1, 0))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.FAILURE),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1 + 1, 127))
    );
    assertEq(
      uint256(Internal.MessageExecutionState.UNTOUCHED),
      uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1 + 1, 128))
    );
  }

  function test_FillExecutionState() public {
    for (uint64 i = 0; i < 384; ++i) {
      s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, i, Internal.MessageExecutionState.FAILURE);
    }

    for (uint64 i = 0; i < 384; ++i) {
      assertEq(
        uint256(Internal.MessageExecutionState.FAILURE),
        uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, i))
      );
    }

    for (uint64 i = 0; i < 3; ++i) {
      assertEq(type(uint256).max, s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, i));
    }

    for (uint64 i = 0; i < 384; ++i) {
      s_offRamp.setExecutionStateHelper(SOURCE_CHAIN_SELECTOR_1, i, Internal.MessageExecutionState.IN_PROGRESS);
    }

    for (uint64 i = 0; i < 384; ++i) {
      assertEq(
        uint256(Internal.MessageExecutionState.IN_PROGRESS),
        uint256(s_offRamp.getExecutionState(SOURCE_CHAIN_SELECTOR_1, i))
      );
    }

    for (uint64 i = 0; i < 3; ++i) {
      // 0x555... == 0b101010101010.....
      assertEq(
        0x5555555555555555555555555555555555555555555555555555555555555555,
        s_offRamp.getExecutionStateBitMap(SOURCE_CHAIN_SELECTOR_1, i)
      );
    }
  }
}
