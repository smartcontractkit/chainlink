// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {LinkTokenInterface} from "../../shared/interfaces/LinkTokenInterface.sol";
import {VRFV2WrapperInterface} from "../interfaces/VRFV2WrapperInterface.sol";

contract VRFV2WrapperUnderFundingConsumer is ConfirmedOwner {
  LinkTokenInterface internal immutable LINK;
  VRFV2WrapperInterface internal immutable VRF_V2_WRAPPER;

  constructor(address _link, address _vrfV2Wrapper) ConfirmedOwner(msg.sender) {
    LINK = LinkTokenInterface(_link);
    VRF_V2_WRAPPER = VRFV2WrapperInterface(_vrfV2Wrapper);
  }

  function makeRequest(uint32 _callbackGasLimit, uint16 _requestConfirmations, uint32 _numWords) external onlyOwner {
    LINK.transferAndCall(
      address(VRF_V2_WRAPPER),
      // Pay less than the needed amount
      VRF_V2_WRAPPER.calculateRequestPrice(_callbackGasLimit) - 1,
      abi.encode(_callbackGasLimit, _requestConfirmations, _numWords)
    );
  }
}
