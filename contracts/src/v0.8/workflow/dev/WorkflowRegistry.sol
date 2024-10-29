// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";

import {Strings} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/Strings.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

contract WorkflowRegistry is OwnerIsCreator, ITypeAndVersion {
  // Bindings
  using EnumerableSet for EnumerableSet.Bytes32Set;
  using EnumerableSet for EnumerableSet.AddressSet;
  using EnumerableSet for EnumerableSet.UintSet;

  // Constants
  string public constant override typeAndVersion = "WorkflowRegistry 1.0.0";
  uint8 private constant MAX_WORKFLOW_NAME_LENGTH = 64;
  uint8 private constant MAX_URL_LENGTH = 200;
  uint8 private constant MAX_PAGINATION_LIMIT = 100;

  // Enums
  enum WorkflowStatus {
    ACTIVE,
    PAUSED
  }

  // Structs
  struct WorkflowMetadata {
    bytes32 workflowID; //     Unique identifier from hash of owner address, WASM binary content, config content and secrets URL
    address owner; // ─────────────────╮ Workflow owner
    uint32 donID; //                   │ Unique identifier for the Workflow DON
    WorkflowStatus status; // ─────────╯ Current status of the workflow (active, paused)
    string workflowName; //    Human readable string capped at 64 characters length
    string binaryURL; //       URL to the WASM binary
    string configURL; //       URL to the config
    string secretsURL; //      URL to the encrypted secrets. Workflow DON applies a default refresh period (e.g. daily)
  }

  // Mappings
  /// @dev Maps an owner address to a set of their workflow (name + owner) hashess.
  mapping(address owner => EnumerableSet.Bytes32Set workflowKeys) private s_ownerWorkflowKeys;
  /// @dev Maps a DON ID to a set of workflow IDs.
  mapping(uint32 donID => EnumerableSet.Bytes32Set workflowKeys) private s_donWorkflowKeys;
  /// @dev Maps a workflow (name + owner) hash to its corresponding workflow metadata.
  mapping(bytes32 workflowKey => WorkflowMetadata workflowMetadata) private s_workflows;
  /// @dev Mapping to track workflows by secretsURL hash (owner + secretsURL).
  /// This is used to find all workflows that have the same secretsURL when a force secrets update event is requested.
  mapping(bytes32 secretsURLHash => EnumerableSet.Bytes32Set workflowKeys) private s_secretsHashToWorkflows;

  /// @dev List of all authorized EOAs/contracts allowed to access this contract's state functions. All view functions are open access.
  EnumerableSet.AddressSet private s_authorizedAddresses;
  /// @dev List of all authorized DON IDs.
  EnumerableSet.UintSet private s_allowedDONs;

  // Events
  event AllowedDONsUpdatedV1(uint32[] donIDs, bool allowed);
  event AuthorizedAddressesUpdatedV1(address[] addresses, bool allowed);
  event WorkflowRegisteredV1(
    bytes32 indexed workflowID,
    address indexed workflowOwner,
    uint32 indexed donID,
    WorkflowStatus status,
    string workflowName,
    string binaryURL,
    string configURL,
    string secretsURL
  );
  event WorkflowUpdatedV1(
    bytes32 indexed oldWorkflowID,
    address indexed workflowOwner,
    uint32 indexed donID,
    bytes32 newWorkflowID,
    string workflowName,
    string binaryURL,
    string configURL,
    string secretsURL
  );
  event WorkflowPausedV1(
    bytes32 indexed workflowID,
    address indexed workflowOwner,
    uint32 indexed donID,
    string workflowName
  );
  event WorkflowActivatedV1(
    bytes32 indexed workflowID,
    address indexed workflowOwner,
    uint32 indexed donID,
    string workflowName
  );
  event WorkflowDeletedV1(
    bytes32 indexed workflowID,
    address indexed workflowOwner,
    uint32 indexed donID,
    string workflowName
  );
  event WorkflowForceUpdateSecretsRequestedV1(string indexed secretsURL, address indexed owner, string[] workflowNames);

  // Errors
  error OnlyAuthorizedAddress();
  error OnlyAllowedDONID();
  error InvalidWorkflowID();
  error WorkflowAlreadyInDesiredStatus();
  error WorkflowDoesNotExist();
  error WorkflowIDAlreadyExists();
  error WorkflowIDNotUpdated();
  error WorkflowContentNotUpdated();
  error WorkflowAlreadyRegistered();
  error WorkflowNameTooLong(uint256 length, uint8 maxAllowedLength);
  error URLTooLong(uint256 length, uint8 maxAllowedLength);

  // Modifiers
  // Check if the caller is an authorized address
  modifier onlyAuthorizedAddresses() {
    if (!s_authorizedAddresses.contains(msg.sender)) revert OnlyAuthorizedAddress();
    _;
  }

  // External state functions
  // ================================================================
  // |                            ADMIN                             |
  // ================================================================
  /**
   * @notice Updates the list of allowed DON IDs.
   * @dev If a DON ID with associated workflows is removed from the allowed DONs list,
   *      the only allowed actions on workflows for that DON are to pause or delete them.
   *      It will no longer be possible to update, activate, or register new workflows for a removed DON.
   * @param donIDs The list of unique identifiers for Workflow DONs.
   * @param allowed True if they should be added to the allowlist, false to remove them.
   */
  function updateAllowedDONs(uint32[] calldata donIDs, bool allowed) external onlyOwner {
    uint256 length = donIDs.length;
    for (uint256 i = 0; i < length; ++i) {
      if (allowed) {
        s_allowedDONs.add(donIDs[i]);
      } else {
        s_allowedDONs.remove(donIDs[i]);
      }
    }

    emit AllowedDONsUpdatedV1(donIDs, allowed);
  }

  /// @notice Updates a list of authorized addresses that can register workflows.
  /// @param addresses The list of addresses.
  /// @param allowed True if they should be added to whitelist, false to remove them.
  /// @dev We don't check if an existing authorized address will be set to false, please take extra caution.
  function updateAuthorizedAddresses(address[] calldata addresses, bool allowed) external onlyOwner {
    uint256 length = addresses.length;
    for (uint256 i = 0; i < length; ++i) {
      if (allowed) {
        s_authorizedAddresses.add(addresses[i]);
      } else {
        s_authorizedAddresses.remove(addresses[i]);
      }
    }

    emit AuthorizedAddressesUpdatedV1(addresses, allowed);
  }

  // ================================================================
  // |                       Workflow Management                    |
  // ================================================================
  /**
   * @notice Registers a new workflow.
   * @dev Registers a new workflow after validating the caller, DON ID, workflow name, and workflow metadata.
   *      This function performs the following steps:
   *      - Validates the caller is authorized and the DON ID is allowed.
   *      - Validates the workflow metadata (workflowID, binaryURL, configURL, secretsURL) lengths.
   *      - Checks if the workflow with the given name already exists for the owner.
   *      - Stores the workflow metadata in the appropriate mappings for the owner and DON.
   *      - Adds the secretsURL to the hash mapping if present.
   *
   * Requirements:
   * - Caller must be an authorized address.
   * - The provided DON ID must be allowed.
   * - The workflow name must not exceed `MAX_WORKFLOW_NAME_LENGTH`.
   * - Workflow metadata must be valid and adhere to set requirements.
   * - Workflow with the given name must not already exist for the owner.
   *
   * Emits:
   * - `WorkflowRegisteredV1` event upon successful registration.
   *
   * @param workflowName The human-readable name for the workflow. Must not exceed 64 characters.
   * @param workflowID The unique identifier for the workflow based on the WASM binary content, config content and secrets URL.
   * @param donID The unique identifier of the Workflow DON that this workflow is associated with.
   * @param status Initial status for this workflow after registration (e.g., Active or Paused).
   * @param binaryURL The URL pointing to the WASM binary for the workflow.
   * @param configURL The URL pointing to the configuration file for the workflow.
   * @param secretsURL The URL pointing to the secrets file for the workflow. Can be empty if there are no secrets.
   */
  function registerWorkflow(
    string calldata workflowName,
    bytes32 workflowID,
    uint32 donID,
    WorkflowStatus status,
    string calldata binaryURL,
    string calldata configURL,
    string calldata secretsURL
  ) external onlyAuthorizedAddresses {
    address sender = msg.sender;
    if (!s_allowedDONs.contains(donID)) revert OnlyAllowedDONID();

    if (bytes(workflowName).length > MAX_WORKFLOW_NAME_LENGTH) {
      revert WorkflowNameTooLong(bytes(workflowName).length, MAX_WORKFLOW_NAME_LENGTH);
    }

    _validateWorkflowMetadata(workflowID, binaryURL, configURL, secretsURL);

    bytes32 workflowKey = _computeOwnerAndStringFieldHashKey(sender, workflowName);
    if (s_workflows[workflowKey].owner != address(0)) {
      revert WorkflowAlreadyRegistered();
    }

    // Create new workflow entry
    s_workflows[workflowKey] = WorkflowMetadata({
      workflowID: workflowID,
      owner: sender,
      donID: donID,
      status: status,
      workflowName: workflowName,
      binaryURL: binaryURL,
      configURL: configURL,
      secretsURL: secretsURL
    });

    s_ownerWorkflowKeys[sender].add(workflowKey);
    s_donWorkflowKeys[donID].add(workflowKey);

    // Hash the secretsURL and add the workflow to the secrets hash mapping
    if (bytes(secretsURL).length > 0) {
      bytes32 secretsHash = _computeOwnerAndStringFieldHashKey(sender, secretsURL);
      s_secretsHashToWorkflows[secretsHash].add(workflowKey);
    }

    emit WorkflowRegisteredV1(workflowID, sender, donID, status, workflowName, binaryURL, configURL, secretsURL);
  }

  /**
   * @notice Updates the workflow metadata for a given workflow.
   * @dev Updates the workflow metadata based on the provided parameters.
   *
   * - If a field needs to be updated, the new value should be provided.
   * - If the value should remain unchanged, provide the same value as before.
   * - To remove an optional field (such as `configURL` or `secretsURL`), pass an empty string ("").
   *
   * This function performs the following steps:
   * - Validates the provided workflow metadata.
   * - Retrieves the workflow by the caller's address and `workflowName`.
   * - Updates only the fields that have changed.
   * - Ensures that the workflow ID (`newWorkflowID`) must change and at least one of the URLs must also change.
   * - Updates the `secretsURL` hash mappings if the `secretsURL` changes.
   *
   * Requirements:
   * - `binaryURL` must always be provided, as it is required.
   * - If only one field is being updated, the other fields must be provided with their current values to keep them unchanged,
   *   otherwise they will be treated as empty strings.
   * - The DON ID must be in the allowed list to perform updates.
   * - The caller must be an authorized address. This means that even if the caller is the owner of the workflow, if they were
   *   later removed from the authorized addresses list, they will not be able to perform updates.
   *
   * Emits:
   * - `WorkflowUpdatedV1` event indicating the workflow has been successfully updated.
   *
   * @param workflowName The human-readable name for the workflow.
   * @param newWorkflowID The rehashed unique identifier for the workflow.
   * @param binaryURL The URL pointing to the WASM binary. Must always be provided.
   * @param configURL The URL pointing to the configuration file. Provide an empty string ("") to remove it.
   * @param secretsURL The URL pointing to the secrets file. Provide an empty string ("") to remove it.
   */
  function updateWorkflow(
    string calldata workflowName,
    bytes32 newWorkflowID,
    string calldata binaryURL,
    string calldata configURL,
    string calldata secretsURL
  ) external onlyAuthorizedAddresses {
    _validateWorkflowMetadata(newWorkflowID, binaryURL, configURL, secretsURL);

    address sender = msg.sender;
    (, WorkflowMetadata storage workflow) = _getWorkflowFromStorageByName(sender, workflowName);

    // Check if the DON ID is allowed
    if (!s_allowedDONs.contains(workflow.donID)) revert OnlyAllowedDONID();

    // Read current values from storage into local variables
    bytes32 currentWorkflowID = workflow.workflowID;
    string memory currentBinaryURL = workflow.binaryURL;
    string memory currentConfigURL = workflow.configURL;
    string memory currentSecretsURL = workflow.secretsURL;

    // Condition to revert: WorkflowID must change, and at least one URL must change
    if (currentWorkflowID == newWorkflowID) {
      revert WorkflowIDNotUpdated();
    }

    // Determine which URLs have changed
    bool sameBinaryURL = Strings.equal(currentBinaryURL, binaryURL);
    bool sameConfigURL = Strings.equal(currentConfigURL, configURL);
    bool sameSecretsURL = Strings.equal(currentSecretsURL, secretsURL);
    if (sameBinaryURL && sameConfigURL && sameSecretsURL) {
      revert WorkflowContentNotUpdated();
    }

    // Update all fields that have changed and the relevant sets
    workflow.workflowID = newWorkflowID;
    if (!sameBinaryURL) {
      workflow.binaryURL = binaryURL;
    }
    if (!sameConfigURL) {
      workflow.configURL = configURL;
    }
    if (!sameSecretsURL) {
      // Remove the old secrets hash if secretsURL is not empty
      if (bytes(currentSecretsURL).length > 0) {
        // Using keccak256 instead of _computeOwnerAndStringFieldHashKey as currentSecretsURL is memory
        bytes32 oldSecretsHash = keccak256(abi.encodePacked(sender, currentSecretsURL));
        s_secretsHashToWorkflows[oldSecretsHash].remove(currentWorkflowID);
      }

      workflow.secretsURL = secretsURL;

      // Add the new secrets hash if secretsURL is not empty
      if (bytes(secretsURL).length > 0) {
        bytes32 newSecretsHash = _computeOwnerAndStringFieldHashKey(sender, secretsURL);
        s_secretsHashToWorkflows[newSecretsHash].add(newWorkflowID);
      }
    }

    // Emit an event after updating the workflow
    emit WorkflowUpdatedV1(
      currentWorkflowID,
      sender,
      workflow.donID,
      newWorkflowID,
      workflow.workflowName,
      binaryURL,
      configURL,
      secretsURL
    );
  }

  /**
   * @notice Pauses an existing workflow.
   * @dev Workflows with any DON ID can be paused.
   *      If a caller was later removed from the authorized addresses list, they will still be able to pause the workflow.
   * @param workflowName The human-readable name for the workflow. It should be unique per owner.
   */
  function pauseWorkflow(string calldata workflowName) external {
    _updateWorkflowStatus(workflowName, WorkflowStatus.PAUSED);
  }

  /**
   * @notice Activates an existing workflow.
   * @dev The DON ID for the workflow must be in the allowed list to perform this action.
   *      The caller must also be an authorized address. This means that even if the caller is the owner of the workflow,
   *      if they were later removed from the authorized addresses list, they will not be able to activate the workflow.
   * @param workflowName The human-readable name for the workflow. It should be unique per owner.
   */
  function activateWorkflow(string calldata workflowName) external onlyAuthorizedAddresses {
    _updateWorkflowStatus(workflowName, WorkflowStatus.ACTIVE);
  }

  /**
   * @notice Deletes an existing workflow, removing it from the contract storage.
   * @dev This function permanently removes a workflow associated with the caller's address.
   *      Workflows with any DON ID can be deleted.
   *      The caller must also be an authorized address. This means that even if the caller is the owner of the workflow,
   *      if they were later removed from the authorized addresses list, they will not be able to delete the workflow.
   *
   * The function performs the following operations:
   * - Retrieves the workflow metadata using the workflow name and owner address.
   * - Ensures that only the owner of the workflow can perform this operation.
   * - Deletes the workflow from the `s_workflows` mapping.
   * - Removes the workflow from associated sets (`s_ownerWorkflowKeys`, `s_donWorkflowKeys`, and `s_secretsHashToWorkflows`).
   *
   * Requirements:
   * - The caller must be the owner of the workflow and an authorized address.
   *
   * Emits:
   * - `WorkflowDeletedV1` event indicating that the workflow has been deleted successfully.
   *
   * @param workflowName The human-readable name of the workflow to delete.
   */
  function deleteWorkflow(string calldata workflowName) external onlyAuthorizedAddresses {
    address sender = msg.sender;

    // Retrieve workflow metadata from storage
    (bytes32 workflowKey, WorkflowMetadata storage workflow) = _getWorkflowFromStorageByName(sender, workflowName);

    // Remove the workflow from the owner and DON mappings
    s_ownerWorkflowKeys[sender].remove(workflowKey);
    s_donWorkflowKeys[workflow.donID].remove(workflowKey);

    // Remove the workflow from the secrets hash set if secretsURL is not empty
    if (bytes(workflow.secretsURL).length > 0) {
      // Using keccak256 instead of _computeOwnerAndStringFieldHashKey as secretsURL is storage ref
      bytes32 secretsHash = keccak256(abi.encodePacked(sender, workflow.secretsURL));
      s_secretsHashToWorkflows[secretsHash].remove(workflowKey);
    }

    // Delete the workflow metadata from storage
    delete s_workflows[workflowKey];

    // Emit an event indicating the workflow has been deleted
    emit WorkflowDeletedV1(workflow.workflowID, sender, workflow.donID, workflowName);
  }

  /**
   * @notice Requests a force update for workflows that share the same secrets URL.
   * @dev This function allows an owner to request a force update for all workflows that share a given `secretsURL`.
   *      The `secretsURL` can be shared between multiple workflows, but they must all belong to the same owner.
   *      This function ensures that the caller owns all workflows associated with the given `secretsURL`.
   *
   * The function performs the following steps:
   * - Hashes the provided `secretsURL` and `msg.sender` to generate a unique mapping key.
   * - Retrieves all workflows associated with the given secrets hash.
   * - Collects the names of all matching workflows and emits an event indicating a force update request.
   *
   * Requirements:
   * - The caller must be the owner of all workflows that share the given `secretsURL`.
   *
   * Emits:
   * - `WorkflowForceUpdateSecretsRequestedV1` event indicating that a force update for workflows using this `secretsURL` has been requested.
   * @param secretsURL The URL pointing to the updated secrets file. This can be shared among multiple workflows.
   */
  function requestForceUpdateSecrets(string calldata secretsURL) external {
    address sender = msg.sender;

    // Use secretsURL and sender hash key to get the mapping key
    bytes32 secretsHash = _computeOwnerAndStringFieldHashKey(sender, secretsURL);

    // Retrieve all workflow keys associated with the given secrets hash
    EnumerableSet.Bytes32Set storage workflowKeys = s_secretsHashToWorkflows[secretsHash];
    uint256 matchCount = workflowKeys.length();

    // No workflows found with the provided secretsURL
    if (matchCount == 0) {
      revert WorkflowDoesNotExist();
    }

    // Create an array for matching workflow names
    string[] memory matchingWorkflowNames = new string[](matchCount);

    // Iterate through matched workflows and gather workflow names
    for (uint256 i = 0; i < matchCount; ++i) {
      bytes32 workflowKey = workflowKeys.at(i);
      WorkflowMetadata storage workflow = s_workflows[workflowKey];
      matchingWorkflowNames[i] = workflow.workflowName;
    }

    // Emit a single event for all matching workflows
    emit WorkflowForceUpdateSecretsRequestedV1(secretsURL, sender, matchingWorkflowNames);
  }

  // External view functions
  // ================================================================
  // |                        Workflow Queries                      |
  // ================================================================
  /// @notice Returns workflow metadata.
  /// @param workflowOwner Address that owns this workflow.
  /// @param workflowName The human-readable name for the workflow.
  /// @return WorkflowMetadata The metadata of the workflow.
  function getWorkflowMetadata(
    address workflowOwner,
    string calldata workflowName
  ) external view returns (WorkflowMetadata memory) {
    bytes32 workflowKey = _computeOwnerAndStringFieldHashKey(workflowOwner, workflowName);
    WorkflowMetadata storage workflow = s_workflows[workflowKey];

    if (workflow.owner == address(0)) revert WorkflowDoesNotExist();

    return workflow;
  }

  /**
   * @notice Retrieves a list of workflow metadata for a specific owner.
   * @dev This function allows paginated retrieval of workflows owned by a particular address.
   *      If the `limit` is set to 0 or exceeds the MAX_PAGINATION_LIMIT, the MAX_PAGINATION_LIMIT will be used instead in both cases.
   * @param workflowOwner The address of the workflow owner for whom the workflow metadata is being retrieved.
   * @param start The index at which to start retrieving workflows (zero-based index).
   *      If the start index is greater than or equal to the total number of workflows, an empty array is returned.
   * @param limit The maximum number of workflow metadata entries to retrieve.
   *      If the limit exceeds the available number of workflows from the start index, only the available entries are returned.
   * @return workflowMetadataList An array of WorkflowMetadata structs containing metadata of workflows owned by the specified owner.
   */
  function getWorkflowMetadataListByOwner(
    address workflowOwner,
    uint256 start,
    uint256 limit
  ) external view returns (WorkflowMetadata[] memory workflowMetadataList) {
    uint256 totalWorkflows = s_ownerWorkflowKeys[workflowOwner].length();
    if (start >= totalWorkflows) {
      return new WorkflowMetadata[](0);
    }

    if (limit > MAX_PAGINATION_LIMIT || limit == 0) {
      limit = MAX_PAGINATION_LIMIT;
    }

    uint256 end = (start + limit > totalWorkflows) ? totalWorkflows : start + limit;

    uint256 resultLength = end - start;
    workflowMetadataList = new WorkflowMetadata[](resultLength);

    for (uint256 i = 0; i < resultLength; ++i) {
      bytes32 workflowKey = s_ownerWorkflowKeys[workflowOwner].at(start + i);
      workflowMetadataList[i] = s_workflows[workflowKey];
    }

    return workflowMetadataList;
  }

  /**
   * @notice Retrieves a list of workflow metadata for a specific DON ID.
   * @dev This function allows paginated retrieval of workflows associated with a particular DON.
   *      If the `limit` is set to 0 or exceeds the MAX_PAGINATION_LIMIT, the MAX_PAGINATION_LIMIT will be used instead in both cases.
   * @param donID The unique identifier of the DON whose associated workflows are being retrieved.
   * @param start The index at which to start retrieving workflows (zero-based index).
   *      If the start index is greater than or equal to the total number of workflows, an empty array is returned.
   * @param limit The maximum number of workflow metadata entries to retrieve.
   *      If the limit exceeds the available number of workflows from the start index, only the available entries are returned.
   * @return workflowMetadataList An array of WorkflowMetadata structs containing metadata of workflows associated with the specified DON ID.
   */
  function getWorkflowMetadataListByDON(
    uint32 donID,
    uint256 start,
    uint256 limit
  ) external view returns (WorkflowMetadata[] memory workflowMetadataList) {
    uint256 totalWorkflows = s_donWorkflowKeys[donID].length();
    if (start >= totalWorkflows) {
      return new WorkflowMetadata[](0);
    }

    if (limit > MAX_PAGINATION_LIMIT || limit == 0) {
      limit = MAX_PAGINATION_LIMIT;
    }

    uint256 end = (start + limit > totalWorkflows) ? totalWorkflows : start + limit;

    uint256 resultLength = end - start;
    workflowMetadataList = new WorkflowMetadata[](resultLength);

    for (uint256 i = 0; i < resultLength; ++i) {
      bytes32 workflowKey = s_donWorkflowKeys[donID].at(start + i);
      workflowMetadataList[i] = s_workflows[workflowKey];
    }

    return workflowMetadataList;
  }

  /**
   * @notice Fetch all allowed DON IDs
   * @return allowedDONs List of all allowed DON IDs
   */
  function getAllAllowedDONs() external view returns (uint32[] memory allowedDONs) {
    uint256 len = s_allowedDONs.length();
    allowedDONs = new uint32[](len);
    for (uint256 i = 0; i < len; ++i) {
      allowedDONs[i] = uint32(s_allowedDONs.at(i));
    }

    return allowedDONs;
  }

  /**
   * @notice Fetch all authorized addresses
   * @return authorizedAddresses List of all authorized addresses
   */
  function getAllAuthorizedAddresses() external view returns (address[] memory authorizedAddresses) {
    uint256 len = s_authorizedAddresses.length();
    authorizedAddresses = new address[](len);
    for (uint256 i = 0; i < len; ++i) {
      authorizedAddresses[i] = s_authorizedAddresses.at(i);
    }

    return authorizedAddresses;
  }

  // ================================================================
  // |                        Internal Helpers                      |
  // ================================================================
  /**
   * @dev Internal function to update the workflow status.
   *
   * This function is used to change the status of an existing workflow, either to "Paused" or "Active".
   *
   * The function performs the following operations:
   * - Retrieves the workflow metadata from storage based on the workflow name.
   * - Only the owner of the workflow can update the status.
   * - Checks if the workflow is already in the desired status, and reverts if no change is necessary to avoid unnecessary
   *   storage writes.
   * - Updates the status of the workflow and emits the appropriate event (`WorkflowPausedV1` or `WorkflowActivatedV1`).
   *
   * Emits:
   * - `WorkflowPausedV1` or `WorkflowActivatedV1` event indicating that the relevant workflow status has been updated.
   *
   * @param workflowName The human-readable name of the workflow.
   * @param newStatus The new status to set for the workflow (either `Paused` or `Active`).
   */
  function _updateWorkflowStatus(string calldata workflowName, WorkflowStatus newStatus) internal {
    address sender = msg.sender;

    // Retrieve workflow metadata once
    (, WorkflowMetadata storage workflow) = _getWorkflowFromStorageByName(sender, workflowName);

    // Avoid unnecessary storage writes if already in the desired status
    if (workflow.status == newStatus) {
      revert WorkflowAlreadyInDesiredStatus();
    }

    // Check if the DON ID is allowed when activating a workflow
    if (newStatus == WorkflowStatus.ACTIVE && !s_allowedDONs.contains(workflow.donID)) {
      revert OnlyAllowedDONID();
    }

    // Update the workflow status
    workflow.status = newStatus;

    // Emit the appropriate event based on newStatus
    if (newStatus == WorkflowStatus.PAUSED) {
      emit WorkflowPausedV1(workflow.workflowID, sender, workflow.donID, workflowName);
    } else if (newStatus == WorkflowStatus.ACTIVE) {
      emit WorkflowActivatedV1(workflow.workflowID, sender, workflow.donID, workflowName);
    }
  }

  /**
   * @dev Internal function to retrieve a workflow by the owner and name.
   *
   * Passing in `msg.sender` as the owner effectively ensures that the workflow key is uniquely tied to the caller's
   * address and workflow name, thus belonging to the caller.
   *
   * The resulting key is used to uniquely identify the workflow for that specific owner.
   *
   * Note: Although a hash collision is theoretically possible, the likelihood is so astronomically low with `keccak256`
   * (which produces a 256-bit hash) that it can be disregarded for all practical purposes.
   *
   * This function is used in place of a modifier in update functions in ensuring workflow ownership and also returns
   * the workflow key and workflow storage.
   *
   * However, if an address besides the msg.sender is passed in, this makes no guarantee on ownership or permissioning
   * and calling functions should handle those separately accordingly.
   *
   * @param sender The address of the owner of the workflow.
   * @param workflowName The human-readable name of the workflow.
   * @return workflowKey The unique key for the workflow.
   * @return workflow The metadata of the workflow.
   */
  function _getWorkflowFromStorageByName(
    address sender,
    string calldata workflowName
  ) internal view returns (bytes32 workflowKey, WorkflowMetadata storage workflow) {
    workflowKey = _computeOwnerAndStringFieldHashKey(sender, workflowName);
    workflow = s_workflows[workflowKey];

    if (workflow.owner == address(0)) revert WorkflowDoesNotExist();

    return (workflowKey, workflow);
  }

  /// @dev Internal function to validate the metadata for a workflow.
  /// @param workflowID The unique identifier for the workflow.
  function _validateWorkflowMetadata(
    bytes32 workflowID,
    string calldata binaryURL,
    string calldata configURL,
    string calldata secretsURL
  ) internal pure {
    if (workflowID == bytes32(0)) revert InvalidWorkflowID();

    if (bytes(binaryURL).length > MAX_URL_LENGTH) {
      revert URLTooLong(bytes(binaryURL).length, MAX_URL_LENGTH);
    }

    if (bytes(configURL).length > MAX_URL_LENGTH) {
      revert URLTooLong(bytes(configURL).length, MAX_URL_LENGTH);
    }

    if (bytes(secretsURL).length > MAX_URL_LENGTH) {
      revert URLTooLong(bytes(secretsURL).length, MAX_URL_LENGTH);
    }
  }

  /**
   * @dev Internal function to compute a unique hash from the owner's address and a given field.
   *
   * This function is used to generate a unique identifier by combining an owner's address with a specific field,
   * ensuring uniqueness for operations like workflow management or secrets handling.
   *
   * The `field` parameter here is of type `calldata string`, which may not work for all use cases.
   *
   * @param owner The address of the owner. Typically used to uniquely associate the field with the owner.
   * @param field A string field, such as the workflow name or secrets URL, that is used to generate the unique hash.
   * @return A unique bytes32 hash computed from the combination of the owner's address and the given field.
   */
  function _computeOwnerAndStringFieldHashKey(address owner, string calldata field) internal pure returns (bytes32) {
    return keccak256(abi.encodePacked(owner, field));
  }
}
