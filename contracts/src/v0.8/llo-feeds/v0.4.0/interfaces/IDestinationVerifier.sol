// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Common} from "../../libraries/Common.sol";
import {IDestinationVerifierProxyInterface} from "./IDestinationVerifierProxyInterface.sol";

interface IDestinationVerifier is IDestinationVerifierProxyInterface {

  /**
   * @notice sets off-chain reporting protocol configuration incl. participating oracles
   * @param signers addresses with which oracles sign the reports
   * @param f number of faulty oracles the system can tolerate
   * @param recipientAddressesAndWeights the addresses and weights of all the recipients to receive rewards
   */
  function setConfig(
    address[] memory signers,
    uint8 f,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) external;

  /**
    * @notice Sets the fee manager address
    * @param _feeManager The address of the fee manager
    */
  function setFeeManager(address _feeManager) external;

  /**
    * @notice Sets the access controller address
    * @param _accessController The address of the access controller
    */
  function setAccessController(address _accessController) external;

  /**
    * @notice Updates the config active status
    * @param DONConfigID The ID of the config to update
    * @param isActive The new config active status
    */
  function setConfigActive(uint256 DONConfigID, bool isActive) external;
}
