// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./OCR2DROracle.sol";
import "../../interfaces/TypeAndVersionInterface.sol";

/**
 * @title OCR2DROracle Factory
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 * @notice Creates OCR2DROracle contracts for node operators
 */
contract OCR2DROracleFactory is TypeAndVersionInterface {
  using EnumerableSet for EnumerableSet.AddressSet;

  EnumerableSet.AddressSet private s_created;

  event OracleCreated(address indexed oracle, address indexed owner, address indexed sender);

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion() external pure override returns (string memory) {
    return "OCR2DROracleFactory 0.0.0";
  }

  /**
   * @notice creates a new Oracle contract with the msg.sender as the proposed owner
   * @notice msg.sender will still need to call oracle.acceptOwnership()
   * @return address Address of a newly deployed oracle
   */
  function deployNewOracle() external returns (address) {
    OCR2DROracle oracle = new OCR2DROracle();
    oracle.transferOwnership(msg.sender);
    s_created.add(address(oracle));
    emit OracleCreated(address(oracle), msg.sender, msg.sender);
    return address(oracle);
  }

  /**
   * @notice Verifies whether this factory deployed an address
   * @param oracleAddress The oracle address in question
   * @return bool True if an oracle has been created at that address
   */
  function created(address oracleAddress) external view returns (bool) {
    return s_created.contains(oracleAddress);
  }
}
