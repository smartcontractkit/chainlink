// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./FunctionsOracle.sol";
import "../../interfaces/TypeAndVersionInterface.sol";

/**
 * @title The Functions Decentralized Oracle Network (Oracle) Factory
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 * @notice Creates FunctionsOracle contracts of a specific version
 */
contract FunctionsOracleFactory is TypeAndVersionInterface {
  using EnumerableSet for EnumerableSet.AddressSet;

  EnumerableSet.AddressSet private s_created;

  event OracleCreated(address indexed don, address indexed owner, address indexed sender);

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion() external pure override returns (string memory) {
    return "FunctionsOracleFactory 0.0.0";
  }

  /**
   * @notice creates a new Oracle contract with the msg.sender as the proposed owner
   * @notice msg.sender will still need to call Oracle.acceptOwnership()
   * @return address Address of a newly deployed Oracle
   */
  function deployNewOracle() external returns (address) {
    FunctionsOracle don = new FunctionsOracle();
    don.transferOwnership(msg.sender);
    s_created.add(address(don));
    emit OracleCreated(address(don), msg.sender, msg.sender);
    return address(don);
  }

  /**
   * @notice Verifies whether this factory deployed an address
   * @param OracleAddress The Oracle address in question
   * @return bool True if an Oracle has been created at that address
   */
  function created(address OracleAddress) external view returns (bool) {
    return s_created.contains(OracleAddress);
  }
}
