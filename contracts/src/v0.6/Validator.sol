// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "./interfaces/ArbitrumInboxInterface.sol";
import "./SimpleWriteAccessController.sol";

contract ArbitrumValidator is SimpleWriteAccessController {
    IInbox arbitrumInbox;
    address flagsAddress;
    // address could use a more robust design: https://eips.ethereum.org/EIPS/eip-1967
    address arbitrumFlag = address(2);

    constructor(address _inboxAddress, address _flagAddress) {
        arbitrumInbox = IInbox(_inboxAddress);
        flagsAddress = _flagAddress;
    }

    function validate(
        uint256 previousRoundId,
        int256 previousAnswer,
        uint256 currentRoundId,
        int256 currentAnswer
    ) external hasAccess() returns (bool) {
        bool raise = currentAnswer == 1;
        if (raise) {
            arbitrumInbox.sendL1FundedContractTransaction(
                30000000,
                400000000,
                flagsAddress,
                abi.encodeWithSignature("raiseFlag(address)", arbitrumFlag)
            );
        } else {
            arbitrumInbox.sendL1FundedContractTransaction(
                30000000,
                400000000,
                flagsAddress,
                abi.encodeWithSignature("lowerFlags(address[])", [arbitrumFlag])
            );
        }
        return true;
    }
}
