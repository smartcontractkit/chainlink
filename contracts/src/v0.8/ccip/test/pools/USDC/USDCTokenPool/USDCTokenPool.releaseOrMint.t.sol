// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Internal} from "../../../../libraries/Internal.sol";
import {Pool} from "../../../../libraries/Pool.sol";
import {RateLimiter} from "../../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../../pools/TokenPool.sol";
import {USDCTokenPool} from "../../../../pools/USDC/USDCTokenPool.sol";
import {MockE2EUSDCTransmitter} from "../../../mocks/MockE2EUSDCTransmitter.sol";
import {USDCTokenPoolSetup} from "./USDCTokenPoolSetup.t.sol";

contract USDCTokenPool_releaseOrMint is USDCTokenPoolSetup {
  // From https://github.com/circlefin/evm-cctp-contracts/blob/377c9bd813fb86a42d900ae4003599d82aef635a/src/messages/BurnMessage.sol#L57
  function _formatMessage(
    uint32 _version,
    bytes32 _burnToken,
    bytes32 _mintRecipient,
    uint256 _amount,
    bytes32 _messageSender
  ) internal pure returns (bytes memory) {
    return abi.encodePacked(_version, _burnToken, _mintRecipient, _amount, _messageSender);
  }

  function testFuzz_ReleaseOrMint_Success(address recipient, uint256 amount) public {
    vm.assume(recipient != address(0) && recipient != address(s_token));
    amount = bound(amount, 0, _getInboundRateLimiterConfig().capacity);

    USDCMessage memory usdcMessage = USDCMessage({
      version: 0,
      sourceDomain: SOURCE_DOMAIN_IDENTIFIER,
      destinationDomain: DEST_DOMAIN_IDENTIFIER,
      nonce: 0x060606060606,
      sender: SOURCE_CHAIN_TOKEN_SENDER,
      recipient: bytes32(uint256(uint160(recipient))),
      destinationCaller: bytes32(uint256(uint160(address(s_usdcTokenPool)))),
      messageBody: _formatMessage(
        0,
        bytes32(uint256(uint160(address(s_token)))),
        bytes32(uint256(uint160(recipient))),
        amount,
        bytes32(uint256(uint160(OWNER)))
      )
    });

    bytes memory message = _generateUSDCMessage(usdcMessage);
    bytes memory attestation = bytes("attestation bytes");

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(
        USDCTokenPool.SourceTokenDataPayload({nonce: usdcMessage.nonce, sourceDomain: SOURCE_DOMAIN_IDENTIFIER})
      ),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    bytes memory offchainTokenData =
      abi.encode(USDCTokenPool.MessageAndAttestation({message: message, attestation: attestation}));

    // The mocked receiver does not release the token to the pool, so we manually do it here
    deal(address(s_token), address(s_usdcTokenPool), amount);

    vm.expectEmit();
    emit TokenPool.Minted(s_routerAllowedOffRamp, recipient, amount);

    vm.expectCall(
      address(s_mockUSDCTransmitter),
      abi.encodeWithSelector(MockE2EUSDCTransmitter.receiveMessage.selector, message, attestation)
    );

    vm.startPrank(s_routerAllowedOffRamp);
    s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: sourceTokenData.extraData,
        offchainTokenData: offchainTokenData
      })
    );
  }

  // https://etherscan.io/tx/0xac9f501fe0b76df1f07a22e1db30929fd12524bc7068d74012dff948632f0883
  function test_ReleaseOrMintRealTx() public {
    bytes memory encodedUsdcMessage =
      hex"000000000000000300000000000000000000127a00000000000000000000000019330d10d9cc8751218eaf51e8885d058642e08a000000000000000000000000bd3fa81b58ba92a82136038b25adec7066af3155000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000af88d065e77c8cc2239327c5edb3a432268e58310000000000000000000000004af08f56978be7dce2d1be3c65c005b41e79401c000000000000000000000000000000000000000000000000000000002057ff7a0000000000000000000000003a23f943181408eac424116af7b7790c94cb97a50000000000000000000000000000000000000000000000000000000000000000000000000000008274119237535fd659626b090f87e365ff89ebc7096bb32e8b0e85f155626b73ae7c4bb2485c184b7cc3cf7909045487890b104efb62ae74a73e32901bdcec91df1bb9ee08ccb014fcbcfe77b74d1263fd4e0b0e8de05d6c9a5913554364abfd5ea768b222f50c715908183905d74044bb2b97527c7e70ae7983c443a603557cac3b1c000000000000000000000000000000000000000000000000000000000000";
    bytes memory attestation = bytes("attestation bytes");

    uint32 nonce = 4730;
    uint32 sourceDomain = 3;
    uint256 amount = 100;

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(USDCTokenPool.SourceTokenDataPayload({nonce: nonce, sourceDomain: sourceDomain})),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    // The mocked receiver does not release the token to the pool, so we manually do it here
    deal(address(s_token), address(s_usdcTokenPool), amount);

    bytes memory offchainTokenData =
      abi.encode(USDCTokenPool.MessageAndAttestation({message: encodedUsdcMessage, attestation: attestation}));

    vm.expectCall(
      address(s_mockUSDCTransmitter),
      abi.encodeWithSelector(MockE2EUSDCTransmitter.receiveMessage.selector, encodedUsdcMessage, attestation)
    );

    vm.startPrank(s_routerAllowedOffRamp);
    s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: OWNER,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: sourceTokenData.extraData,
        offchainTokenData: offchainTokenData
      })
    );
  }

  // Reverts
  function test_RevertWhen_UnlockingUSDCFailed() public {
    vm.startPrank(s_routerAllowedOffRamp);
    s_mockUSDCTransmitter.setShouldSucceed(false);

    uint256 amount = 13255235235;

    USDCMessage memory usdcMessage = USDCMessage({
      version: 0,
      sourceDomain: SOURCE_DOMAIN_IDENTIFIER,
      destinationDomain: DEST_DOMAIN_IDENTIFIER,
      nonce: 0x060606060606,
      sender: SOURCE_CHAIN_TOKEN_SENDER,
      recipient: bytes32(uint256(uint160(address(s_mockUSDC)))),
      destinationCaller: bytes32(uint256(uint160(address(s_usdcTokenPool)))),
      messageBody: _formatMessage(
        0,
        bytes32(uint256(uint160(address(s_token)))),
        bytes32(uint256(uint160(OWNER))),
        amount,
        bytes32(uint256(uint160(OWNER)))
      )
    });

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(
        USDCTokenPool.SourceTokenDataPayload({nonce: usdcMessage.nonce, sourceDomain: SOURCE_DOMAIN_IDENTIFIER})
      ),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    bytes memory offchainTokenData = abi.encode(
      USDCTokenPool.MessageAndAttestation({message: _generateUSDCMessage(usdcMessage), attestation: bytes("")})
    );

    vm.expectRevert(USDCTokenPool.UnlockingUSDCFailed.selector);

    s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: OWNER,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: sourceTokenData.extraData,
        offchainTokenData: offchainTokenData
      })
    );
  }

  function test_RevertWhen_TokenMaxCapacityExceeded() public {
    uint256 capacity = _getInboundRateLimiterConfig().capacity;
    uint256 amount = 10 * capacity;
    address recipient = address(1);
    vm.startPrank(s_routerAllowedOffRamp);

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(USDCTokenPool.SourceTokenDataPayload({nonce: 1, sourceDomain: SOURCE_DOMAIN_IDENTIFIER})),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    bytes memory offchainTokenData =
      abi.encode(USDCTokenPool.MessageAndAttestation({message: bytes(""), attestation: bytes("")}));

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.TokenMaxCapacityExceeded.selector, capacity, amount, address(s_token))
    );

    s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: sourceTokenData.extraData,
        offchainTokenData: offchainTokenData
      })
    );
  }
}
