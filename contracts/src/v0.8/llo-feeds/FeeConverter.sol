// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {IFeeConverter} from "./interfaces/IFeeConverter.sol";
import {ConfirmedOwner} from "../ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../shared/vendor/IERC165.sol";

/*
 * @title FeeManager
 * @author Michael Fletcher
 * @notice This contract is used for the handling the conversion of fees from one token to another.
 */
contract FeeConverter is IFeeConverter, ConfirmedOwner, TypeAndVersionInterface {
  //the native token address
  address private immutable NATIVE_ADDRESS;

  /**
   * @notice Construct the FeeManager contract
   * @param nativeAddress The address of the NATIVE token
   */
  constructor(address nativeAddress) ConfirmedOwner(msg.sender) {
    //set the native address
    NATIVE_ADDRESS = nativeAddress;
  }

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "FeeConverter 0.0.1";
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) external pure override returns (bool) {
    return false;
  }
}
