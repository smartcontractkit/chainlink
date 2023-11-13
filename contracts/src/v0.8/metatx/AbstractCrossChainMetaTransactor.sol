// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import {IERC2771Recipient} from "../vendor/IERC2771Recipient.sol";
import {IRouterClient} from "../ccip/interfaces/IRouterClient.sol";
import {Client} from "../ccip/libraries/Client.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";

/// @dev AbstractCrossChainMetaTransactor extends ERC20 token to add cross chain transfer functionality
/// @dev Also, it trusts ERC2771 forwarder to forward meta-transactions
abstract contract AbstractCrossChainMetaTransactor is OwnerIsCreator, IERC2771Recipient {
  /// @dev forwarder verifies signatures for meta transactions and forwards the
  /// @dev request to this contract
  address private immutable i_forwarder;
  IRouterClient private immutable i_ccipRouter;
  /// @dev address of account privileged to fund/withdraw native token for CCIP fee
  address private immutable i_ccipFeeProviderAddress;
  /// @dev CCIP chain ID
  uint64 private immutable i_ccipChainId;

  /// @notice This error is thrown whenever a zero address is passed
  error ZeroAddress();

  constructor(address forwarder, address ccipRouter, address ccipFeeProviderAddress, uint64 ccipChainId) {
    if (forwarder == address(0) || ccipRouter == address(0) || ccipFeeProviderAddress == address(0)) {
      revert ZeroAddress();
    }
    i_forwarder = forwarder;
    i_ccipRouter = IRouterClient(ccipRouter);
    i_ccipFeeProviderAddress = ccipFeeProviderAddress;
    i_ccipChainId = ccipChainId;
  }

  // ================================================================
  // |                        Meta-Transaction                      |
  // ================================================================

  /// @dev Transfers "amount" of this token to receiver address in destination chain.
  /// @dev contract needs to be funded with native token, which is used for CCIP fee
  /// @param receiver This is the address that tokens are transferred to
  /// @param amount Total token amount to be transferred
  /// @param destinationChainId Destination chain ID
  function metaTransfer(
    address receiver,
    uint256 amount,
    uint64 destinationChainId
  ) external virtual validateTrustedForwarder returns (bytes32) {
    if (!_isCrossChainTransfer(destinationChainId)) {
      _transfer(_msgSender(), receiver, amount);
      return ""; // return empty bytes32 because there is no ccip message ID
    }
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({token: address(this), amount: amount});
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(receiver),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(0), // use native token instead of ERC20 tokens
      extraArgs: ""
    });
    _transfer(_msgSender(), address(this), amount);
    _approve(address(this), address(i_ccipRouter), amount);

    uint256 fee = i_ccipRouter.getFee(destinationChainId, message);

    return i_ccipRouter.ccipSend{value: fee}(destinationChainId, message);
  }

  /// @dev Sets `amount` as the allowance of `spender` over the `owner` s tokens.
  /// @dev Sample implementation: https://github.com/OpenZeppelin/openzeppelin-contracts/blob/d59306bd06a241083841c2e4a39db08e1f3722cc/contracts/token/ERC20/ERC20.sol#L308-L314
  /// @param owner Token owner approving allowance
  /// @param spender Approved token spender
  /// @param amount Total token amount to be approved
  function _approve(address owner, address spender, uint256 amount) internal virtual;

  /// @dev Moves `amount` of tokens from `sender` to `recipient`.
  /// @dev Sample implementation: https://github.com/OpenZeppelin/openzeppelin-contracts/blob/d59306bd06a241083841c2e4a39db08e1f3722cc/contracts/token/ERC20/ERC20.sol#L222-L240
  /// @param sender Token sender
  /// @param recipient Token recipient
  /// @param amount Total token amount to be approved
  function _transfer(address sender, address recipient, uint256 amount) internal virtual;

  function _isCrossChainTransfer(uint64 chainId) private view returns (bool) {
    return i_ccipChainId != chainId;
  }

  function getCCIPChainId() public view returns (uint64) {
    return i_ccipChainId;
  }

  function getCCIPRouter() public view returns (IRouterClient) {
    return i_ccipRouter;
  }

  // ================================================================
  // |                        Contract Funding                      |
  // ================================================================

  /// @dev For cross-chain transfers, this contract needs to be funded with native token
  receive() external payable {}

  error WithdrawFailure();

  /// @dev withdraws all native tokens from this contract to CCIP Fee Provider address
  /// @dev Only callable by CCIP Fee Provider address
  function withdrawNative() external validateCCIPFeeProvider {
    uint256 amount = address(this).balance;
    // Owner can receive Ether since the address of owner is payable
    (bool success, ) = i_ccipFeeProviderAddress.call{value: amount}("");
    if (!success) {
      revert WithdrawFailure();
    }
  }

  function getCCIPFeeProvider() public view returns (address) {
    return i_ccipFeeProviderAddress;
  }

  // ================================================================
  // |                        Forwarder                             |
  // ================================================================

  /// @notice Address of the trusted forwarder
  /// @return forwarder The address of the Forwarder contract that is being used.
  function getTrustedForwarder() public view returns (address forwarder) {
    return i_forwarder;
  }

  /// @inheritdoc IERC2771Recipient
  function isTrustedForwarder(address forwarder) public view override returns (bool) {
    return forwarder == i_forwarder;
  }

  /// @inheritdoc IERC2771Recipient
  function _msgSender() internal view override returns (address msgSender) {
    if (msg.data.length >= 20 && isTrustedForwarder(msg.sender)) {
      // At this point we know that the sender is a trusted forwarder,
      // so we trust that the last bytes of msg.data are the verified sender address.
      // extract sender address from the end of msg.data
      assembly {
        msgSender := shr(96, calldataload(sub(calldatasize(), 20)))
      }
    } else {
      msgSender = msg.sender;
    }
  }

  /// @inheritdoc IERC2771Recipient
  function _msgData() internal view override returns (bytes calldata msgData) {
    if (msg.data.length >= 20 && isTrustedForwarder(msg.sender)) {
      return msg.data[0:msg.data.length - 20];
    } else {
      return msg.data;
    }
  }

  function getForwarder() public view returns (address) {
    return i_forwarder;
  }

  // ================================================================
  // |                      Access Control                           |
  // ================================================================

  error MustBeTrustedForwarder(address sender);
  error MustBeCCIPFeeProvider(address sender);

  modifier validateTrustedForwarder() {
    if (!isTrustedForwarder(msg.sender)) {
      revert MustBeTrustedForwarder(msg.sender);
    }
    _;
  }

  modifier validateCCIPFeeProvider() {
    if (msg.sender != i_ccipFeeProviderAddress) {
      revert MustBeCCIPFeeProvider(msg.sender);
    }
    _;
  }
}
