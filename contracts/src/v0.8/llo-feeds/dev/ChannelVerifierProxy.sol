// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {IChannelVerifierProxy} from "./interfaces/IChannelVerifierProxy.sol";

contract ChannelVerifierProxy is IChannelVerifierProxy, ConfirmedOwner, TypeAndVersionInterface {

    error AccessForbidden();

    constructor() ConfirmedOwner(msg.sender) {
    }

    /// @inheritdoc TypeAndVersionInterface
    function typeAndVersion() external pure override returns (string memory) {
        return "ChannelVerifierProxy 0.0.0";
    }

    /// @inheritdoc IChannelVerifierProxy
    function verify(
        bytes calldata payload,
        bytes calldata parameterPayload
    ) external payable returns (bytes memory) {
        _verify(payload, parameterPayload);

        return new bytes(0);
    }

    /// @inheritdoc IChannelVerifierProxy
    function verifyBulk(
        bytes[] calldata payloads,
        bytes calldata parameterPayload
    ) external payable returns (bytes[] memory verifiedReports) {
        verifiedReports = new bytes[](payloads.length);
        for (uint256 i; i < payloads.length; ++i) {
            _verify(payloads[i], parameterPayload);
        }

        return new bytes[](0);
    }

    function _verify(bytes calldata payload, bytes calldata parameterPayload) internal pure returns (bytes memory, bytes memory) {
        return (payload, parameterPayload);
    }

    function supportsInterface(bytes4 interfaceId) external pure returns (bool) {
        return interfaceId == type(IChannelVerifierProxy).interfaceId;
    }
}