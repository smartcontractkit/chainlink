// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Common} from "../../libraries/Common.sol";

interface IDestinationVerifier is IERC165 {
  /**
   * @notice Verifies that the data encoded has been signed correctly using the signatures included within the payload.
   * @param signedReport The encoded data to be verified.
   * @param parameterPayload The encoded parameters to be used in the verification and billing process.
   * @param sender The address that requested to verify the contract.Used for logging and applying the fee.
   * @dev Verification is typically only done through the proxy contract so we can't just use msg.sender.
   * @return verifierResponse The encoded verified response.
   */
  function verify(bytes calldata signedReport, bytes calldata parameterPayload, address sender) external payable returns (bytes memory verifierResponse);

  /**
   * @notice Bulk verifies that the data encoded has been signed correctly using the signatures included within the payload.
   * @param signedReports The encoded data to be verified.
   * @param parameterPayload The encoded parameters to be used in the verification and billing process.
   * @param sender The address that requested to verify the contract. Used for logging and applying the fee.
   * @dev Verification is typically only done through the proxy contract so we can't just use msg.sender.
   * @return verifiedReports The encoded verified responses.
   */
  function verifyBulk(bytes[] calldata signedReports, bytes calldata parameterPayload, address sender) external payable returns (bytes[] memory verifiedReports);

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
    * @param DONConfigID The new config active status
    * @param isActive The new config active status
    */
  function setConfigActive(bytes24 DONConfigID, bool isActive) external;


  /**
   * @notice Returns the current access controller
   * @return The address of the access controller
   */
  function getAccessController() external view returns (address);

  /**
    * @notice Returns the current fee manager
    * @return The address of the fee manager
    */
  function getFeeManager() external view returns (address);
}
