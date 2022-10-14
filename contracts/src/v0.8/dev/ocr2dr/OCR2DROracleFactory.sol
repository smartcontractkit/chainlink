// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "./OCR2DROracle.sol";

/**
 * @title OCR2DROracle Factory
 * @notice Creates OCR2DROracle contracts for node operators
 */
contract OCR2DROracleFactory {
  using EnumerableSet for EnumerableSet.AddressSet;

  EnumerableSet.AddressSet private s_created;

  event OracleCreated(address indexed oracle, address indexed owner, address indexed sender);

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion() external pure virtual returns (string memory) {
    return "OCR2DROracleFactory 0.0.0";
  }

  /**
   * @notice creates a new Oracle contract with the msg.sender as owner
   * @param donPublicKey DON's public key used to encrypt user secrets
   */
  function deployNewOracle(bytes32 donPublicKey) external returns (address) {
    OCR2DROracle oracle = new OCR2DROracle(msg.sender, donPublicKey);

    s_created.add(address(oracle));
    emit OracleCreated(address(oracle), msg.sender, msg.sender);

    return address(oracle);
  }

  /**
   * @notice indicates whether this factory deployed an address
   */
  function created(address query) external view returns (bool) {
    return s_created.contains(query);
  }
}
