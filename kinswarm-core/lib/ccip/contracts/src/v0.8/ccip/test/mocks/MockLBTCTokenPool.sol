// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../../libraries/Pool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice This mock contract facilitates testing of LBTC token transfers by burning and minting tokens.
contract MockLBTCTokenPool is TokenPool, ITypeAndVersion {
    using SafeERC20 for IERC20;

    string public constant override typeAndVersion = "MockLBTCTokenPool 1.5.1";

    // This variable i_destPoolData will have either a 32-byte or non-32-byte value, which will change the off-chain behavior.
    // If it is 32 bytes, the off-chain will consider it as attestation enabled and call the attestation API.
    // If it is non-32 bytes, the off-chain will consider it as attestation disabled.
    bytes public i_destPoolData;

    constructor(
        IERC20 token,
        address[] memory allowlist,
        address rmnProxy,
        address router,
        bytes memory destPoolData
    ) TokenPool(token, 8, allowlist, rmnProxy, router) {
        i_destPoolData = destPoolData;
    }

    function lockOrBurn(
        Pool.LockOrBurnInV1 calldata lockOrBurnIn
    ) public virtual override returns (Pool.LockOrBurnOutV1 memory) {
        IBurnMintERC20(address(i_token)).burn(lockOrBurnIn.amount);
        emit Burned(msg.sender, lockOrBurnIn.amount);

        return
            Pool.LockOrBurnOutV1({
            destTokenAddress: getRemoteToken(
                lockOrBurnIn.remoteChainSelector
            ),
            destPoolData: i_destPoolData
        });
    }

    function releaseOrMint(
        Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
    ) public virtual override returns (Pool.ReleaseOrMintOutV1 memory) {
        IBurnMintERC20(address(i_token)).mint(releaseOrMintIn.receiver, releaseOrMintIn.amount);

        emit Minted(
            msg.sender,
            releaseOrMintIn.receiver,
            releaseOrMintIn.amount
        );

        return
            Pool.ReleaseOrMintOutV1({
            destinationAmount: releaseOrMintIn.amount
        });
    }
}

