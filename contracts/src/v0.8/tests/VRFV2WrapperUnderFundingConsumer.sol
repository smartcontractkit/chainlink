// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../ConfirmedOwner.sol";
import "../interfaces/ILinkToken.sol";
import "../interfaces/IVRFV2Wrapper.sol";

contract VRFV2WrapperUnderFundingConsumer is ConfirmedOwner {
  ILinkToken internal immutable LINK;
  IVRFV2Wrapper internal immutable VRF_V2_WRAPPER;

  constructor(address _link, address _vrfV2Wrapper) ConfirmedOwner(msg.sender) {
    LINK = ILinkToken(_link);
    VRF_V2_WRAPPER = IVRFV2Wrapper(_vrfV2Wrapper);
  }

  function makeRequest(
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords
  ) external onlyOwner {
    LINK.transferAndCall(
      address(VRF_V2_WRAPPER),
      // Pay less than the needed amount
      VRF_V2_WRAPPER.calculateRequestPrice(_callbackGasLimit) - 1,
      abi.encode(_callbackGasLimit, _requestConfirmations, _numWords)
    );
  }
}
