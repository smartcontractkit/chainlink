// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./FunctionsDON.sol";
import "../../interfaces/TypeAndVersionInterface.sol";

/**
 * @title The Functions Decentralized Oracle Network (DON) Factory
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 * @notice Creates FunctionsDON contracts of a specific version
 */
contract FunctionsDONFactory is TypeAndVersionInterface {
  using EnumerableSet for EnumerableSet.AddressSet;

  EnumerableSet.AddressSet private s_created;

  event DONCreated(address indexed don, address indexed owner, address indexed sender);

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion() external pure override returns (string memory) {
    return "FunctionsDONFactory 0.0.0";
  }

  /**
   * @notice creates a new DON contract with the msg.sender as the proposed owner
   * @notice msg.sender will still need to call DON.acceptOwnership()
   * @return address Address of a newly deployed DON
   */
  function deployNewDON() external returns (address) {
    FunctionsDON don = new FunctionsDON();
    don.transferOwnership(msg.sender);
    s_created.add(address(don));
    emit DONCreated(address(don), msg.sender, msg.sender);
    return address(don);
  }

  /**
   * @notice Verifies whether this factory deployed an address
   * @param DONAddress The DON address in question
   * @return bool True if an DON has been created at that address
   */
  function created(address DONAddress) external view returns (bool) {
    return s_created.contains(DONAddress);
  }
}
