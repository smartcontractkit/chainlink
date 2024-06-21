// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice: IMPORTANT NOTICE for anyone who wants to use this contract
/// @notice Source: https://github.com/ethereum-optimism/optimism/blob/71b93116738ee98c9f8713b1a5dfe626ce06c1b2/packages/contracts-bedrock/src/libraries/GasPayingToken.sol
/// @notice The original code was trimmed down to include only the necessary interface elements required to interact with GasPriceOracle
/// @notice We need this file so that Solidity compiler will not complain because some functions don't exist
/// @notice In reality, we don't embed this code into our own contracts, instead we make cross-contract calls on predeployed GasPriceOracle contract

import {Storage} from "./Storage.sol";
import {Constants} from "./Constants.sol";
import {LibString} from "../deps/LibString.sol";

/// @title IGasToken
/// @notice Implemented by contracts that are aware of the custom gas token used
///         by the L2 network.
interface IGasToken {
  /// @notice Getter for the ERC20 token address that is used to pay for gas and its decimals.
  function gasPayingToken() external view returns (address, uint8);
  /// @notice Returns the gas token name.
  function gasPayingTokenName() external view returns (string memory);
  /// @notice Returns the gas token symbol.
  function gasPayingTokenSymbol() external view returns (string memory);
  /// @notice Returns true if the network uses a custom gas token.
  function isCustomGasToken() external view returns (bool);
}

/// @title GasPayingToken
/// @notice Handles reading and writing the custom gas token to storage.
///         To be used in any place where gas token information is read or
///         written to state. If multiple contracts use this library, the
///         values in storage should be kept in sync between them.
library GasPayingToken {
  /// @notice The storage slot that contains the address and decimals of the gas paying token
  bytes32 internal constant GAS_PAYING_TOKEN_SLOT = bytes32(uint256(keccak256("opstack.gaspayingtoken")) - 1);

  /// @notice The storage slot that contains the ERC20 `name()` of the gas paying token
  bytes32 internal constant GAS_PAYING_TOKEN_NAME_SLOT = bytes32(uint256(keccak256("opstack.gaspayingtokenname")) - 1);

  /// @notice the storage slot that contains the ERC20 `symbol()` of the gas paying token
  bytes32 internal constant GAS_PAYING_TOKEN_SYMBOL_SLOT =
    bytes32(uint256(keccak256("opstack.gaspayingtokensymbol")) - 1);

  /// @notice Reads the gas paying token and its decimals from the magic
  ///         storage slot. If nothing is set in storage, then the ether
  ///         address is returned instead.
  function getToken() internal view returns (address addr_, uint8 decimals_) {
    bytes32 slot = Storage.getBytes32(GAS_PAYING_TOKEN_SLOT);
    addr_ = address(uint160(uint256(slot) & uint256(type(uint160).max)));
    if (addr_ == address(0)) {
      addr_ = Constants.ETHER;
      decimals_ = 18;
    } else {
      decimals_ = uint8(uint256(slot) >> 160);
    }
  }

  /// @notice Reads the gas paying token's name from the magic storage slot.
  ///         If nothing is set in storage, then the ether name, 'Ether', is returned instead.
  function getName() internal view returns (string memory name_) {
    (address addr, ) = getToken();
    if (addr == Constants.ETHER) {
      name_ = "Ether";
    } else {
      name_ = LibString.fromSmallString(Storage.getBytes32(GAS_PAYING_TOKEN_NAME_SLOT));
    }
  }

  /// @notice Reads the gas paying token's symbol from the magic storage slot.
  ///         If nothing is set in storage, then the ether symbol, 'ETH', is returned instead.
  function getSymbol() internal view returns (string memory symbol_) {
    (address addr, ) = getToken();
    if (addr == Constants.ETHER) {
      symbol_ = "ETH";
    } else {
      symbol_ = LibString.fromSmallString(Storage.getBytes32(GAS_PAYING_TOKEN_SYMBOL_SLOT));
    }
  }

  /// @notice Writes the gas paying token, its decimals, name and symbol to the magic storage slot.
  function set(address _token, uint8 _decimals, bytes32 _name, bytes32 _symbol) internal {
    Storage.setBytes32(GAS_PAYING_TOKEN_SLOT, bytes32((uint256(_decimals) << 160) | uint256(uint160(_token))));
    Storage.setBytes32(GAS_PAYING_TOKEN_NAME_SLOT, _name);
    Storage.setBytes32(GAS_PAYING_TOKEN_SYMBOL_SLOT, _symbol);
  }

  /// @notice Maps a string to a normalized null-terminated small string.
  function sanitize(string memory _str) internal pure returns (bytes32) {
    require(bytes(_str).length <= 32, "GasPayingToken: string cannot be greater than 32 bytes");

    return LibString.toSmallString(_str);
  }
}
