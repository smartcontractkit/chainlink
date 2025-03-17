// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IReceiver} from "../keystone/interfaces/IReceiver.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";
import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";
import {IDataFeedsCache} from "./interfaces/IDataFeedsCache.sol";
import {ITokenRecover} from "./interfaces/ITokenRecover.sol";

import {IERC165} from "../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC165.sol";
import {IERC20} from "../vendor/openzeppelin-solidity/v5.0.2/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../vendor/openzeppelin-solidity/v5.0.2/contracts/token/ERC20/utils/SafeERC20.sol";

contract DataFeedsCache is IDataFeedsCache, IReceiver, ITokenRecover, ITypeAndVersion, OwnerIsCreator {
  using SafeERC20 for IERC20;

  string public constant override typeAndVersion = "DataFeedsCache 1.0.0";

  // solhint-disable-next-line
  uint256 public constant override version = 7;

  /// Cache State

  struct WorkflowMetadata {
    address allowedSender; // Address of the sender allowed to send new reports
    address allowedWorkflowOwner; // ─╮ Address of the workflow owner
    bytes10 allowedWorkflowName; // ──╯ Name of the workflow
  }

  struct FeedConfig {
    uint8[] bundleDecimals; // Only appliciable to Bundle reports - Decimal reports have decimals encoded into the DataId.
    string description; // Description of the feed (e.g. "LINK / USD")
    WorkflowMetadata[] workflowMetadata; // Metadata for the feed
  }

  struct ReceivedBundleReport {
    bytes32 dataId; // Data ID of the feed from the received report
    uint32 timestamp; // Timestamp of the feed from the received report
    bytes bundle; // Report data in raw bytes
  }

  struct ReceivedDecimalReport {
    bytes32 dataId; // Data ID of the feed from the received report
    uint32 timestamp; // ─╮ Timestamp of the feed from the received report
    uint224 answer; // ───╯  Report data in uint224
  }

  struct StoredBundleReport {
    bytes bundle; // The latest bundle report stored for a feed
    uint32 timestamp; // The timestamp of the latest bundle report
  }

  struct StoredDecimalReport {
    uint224 answer; // ───╮ The latest decimal report stored for a feed
    uint32 timestamp; // ─╯ The timestamp of the latest decimal report
  }

  /// The message sender determines which feed is being requested, as each proxy has a single associated feed
  mapping(address aggProxy => bytes16 dataId) private s_aggregatorProxyToDataId;

  /// The latest decimal reports for each decimal feed. This will always equal s_decimalReports[s_dataIdToRoundId[dataId]][dataId]
  mapping(bytes16 dataId => StoredDecimalReport) private s_latestDecimalReports;

  /// Decimal reports for each feed, per round
  mapping(uint256 roundId => mapping(bytes16 dataId => StoredDecimalReport)) private s_decimalReports;

  /// The latest bundle reports for each bundle feed
  mapping(bytes16 dataId => StoredBundleReport) private s_latestBundleReports;

  /// The latest round id for each feed
  mapping(bytes16 dataId => uint256 roundId) private s_dataIdToRoundId;

  /// Addresses that are permitted to configure all feeds
  mapping(address feedAdmin => bool isFeedAdmin) private s_feedAdmins;

  mapping(bytes16 dataId => FeedConfig) private s_feedConfigs;

  /// Whether a given Sender and Workflow have permission to write feed updates.
  /// reportHash is the keccak256 hash of the abi.encoded(dataId, sender, workflowOwner and workflowName)
  mapping(bytes32 reportHash => bool) private s_writePermissions;

  event BundleReportUpdated(bytes16 indexed dataId, uint256 indexed timestamp, bytes bundle);
  event DecimalReportUpdated(
    bytes16 indexed dataId, uint256 indexed roundId, uint256 indexed timestamp, uint224 answer
  );
  event DecimalFeedConfigSet(
    bytes16 indexed dataId, uint8 decimals, string description, WorkflowMetadata[] workflowMetadata
  );
  event BundleFeedConfigSet(
    bytes16 indexed dataId, uint8[] decimals, string description, WorkflowMetadata[] workflowMetadata
  );
  event FeedConfigRemoved(bytes16 indexed dataId);
  event TokenRecovered(address indexed token, address indexed to, uint256 amount);

  event FeedAdminSet(address indexed feedAdmin, bool indexed isAdmin);

  event ProxyDataIdRemoved(address indexed proxy, bytes16 indexed dataId);
  event ProxyDataIdUpdated(address indexed proxy, bytes16 indexed dataId);

  event InvalidUpdatePermission(bytes16 indexed dataId, address sender, address workflowOwner, bytes10 workflowName);
  event StaleDecimalReport(bytes16 indexed dataId, uint256 reportTimestamp, uint256 latestTimestamp);
  event StaleBundleReport(bytes16 indexed dataId, uint256 reportTimestamp, uint256 latestTimestamp);

  error ArrayLengthMismatch();
  error EmptyConfig();
  error ErrorSendingNative(address to, uint256 amount, bytes data);
  error FeedNotConfigured(bytes16 dataId);
  error InsufficientBalance(uint256 balance, uint256 requiredBalance);
  error InvalidAddress(address addr);
  error InvalidDataId();
  error InvalidWorkflowName(bytes10 workflowName);
  error UnauthorizedCaller(address caller);
  error NoMappingForSender(address proxy);

  modifier onlyFeedAdmin() {
    if (!s_feedAdmins[msg.sender]) revert UnauthorizedCaller(msg.sender);
    _;
  }

  /// @inheritdoc IERC165
  function supportsInterface(
    bytes4 interfaceId
  ) public pure returns (bool) {
    return (
      interfaceId == type(IDataFeedsCache).interfaceId || interfaceId == type(IERC165).interfaceId
        || interfaceId == type(IReceiver).interfaceId || interfaceId == type(ITokenRecover).interfaceId
        || interfaceId == type(ITypeAndVersion).interfaceId
    );
  }

  /// @notice Get the workflow metadata of a feed
  /// @param dataId data ID of the feed
  /// @param startIndex The cursor to start fetching the metadata from
  /// @param maxCount The number of metadata to fetch
  /// @return workflowMetadata The metadata of the feed
  function getFeedMetadata(
    bytes16 dataId,
    uint256 startIndex,
    uint256 maxCount
  ) external view returns (WorkflowMetadata[] memory workflowMetadata) {
    FeedConfig storage feedConfig = s_feedConfigs[dataId];

    uint256 workflowMetadataLength = feedConfig.workflowMetadata.length;

    if (workflowMetadataLength == 0) {
      revert FeedNotConfigured(dataId);
    }

    if (startIndex >= workflowMetadataLength) return new WorkflowMetadata[](0);
    uint256 endIndex = startIndex + maxCount;
    endIndex = endIndex > workflowMetadataLength || maxCount == 0 ? workflowMetadataLength : endIndex;

    workflowMetadata = new WorkflowMetadata[](endIndex - startIndex);
    for (uint256 idx; idx < workflowMetadata.length; idx++) {
      workflowMetadata[idx] = feedConfig.workflowMetadata[idx + startIndex];
    }

    return workflowMetadata;
  }

  /// @notice Checks to see if this data ID, msg.sender, workflow owner, and workflow name are permissioned
  /// @param dataId The data ID for the feed
  /// @param workflowMetadata workflow metadata
  function checkFeedPermission(
    bytes16 dataId,
    WorkflowMetadata memory workflowMetadata
  ) external view returns (bool hasPermission) {
    bytes32 permission = _createReportHash(
      dataId,
      workflowMetadata.allowedSender,
      workflowMetadata.allowedWorkflowOwner,
      workflowMetadata.allowedWorkflowName
    );
    return s_writePermissions[permission];
  }

  // ================================================================
  // │                  Contract Config Interface                   │
  // ================================================================

  /// @notice Initializes the config for a decimal feed
  /// @param dataIds The data IDs of the feeds to configure
  /// @param descriptions The descriptions of the feeds
  /// @param workflowMetadata List of workflow metadata (owners, senders, and names) for every feed
  function setDecimalFeedConfigs(
    bytes16[] calldata dataIds,
    string[] calldata descriptions,
    WorkflowMetadata[] calldata workflowMetadata
  ) external onlyFeedAdmin {
    if (workflowMetadata.length == 0 || dataIds.length == 0) {
      revert EmptyConfig();
    }

    if (dataIds.length != descriptions.length) {
      revert ArrayLengthMismatch();
    }

    for (uint256 i; i < dataIds.length; ++i) {
      bytes16 dataId = dataIds[i];
      if (dataId == bytes16(0)) revert InvalidDataId();
      FeedConfig storage feedConfig = s_feedConfigs[dataId];

      if (feedConfig.workflowMetadata.length > 0) {
        // Feed is already configured, remove the previous config
        for (uint256 j; j < feedConfig.workflowMetadata.length; ++j) {
          WorkflowMetadata memory feedCurrentWorkflowMetadata = feedConfig.workflowMetadata[j];
          bytes32 reportHash = _createReportHash(
            dataId,
            feedCurrentWorkflowMetadata.allowedSender,
            feedCurrentWorkflowMetadata.allowedWorkflowOwner,
            feedCurrentWorkflowMetadata.allowedWorkflowName
          );
          delete s_writePermissions[reportHash];
        }

        delete s_feedConfigs[dataId];

        emit FeedConfigRemoved(dataId);
      }

      for (uint256 j; j < workflowMetadata.length; ++j) {
        WorkflowMetadata memory feedWorkflowMetadata = workflowMetadata[j];
        // Do those checks only once for the first data id
        if (i == 0) {
          if (feedWorkflowMetadata.allowedSender == address(0)) {
            revert InvalidAddress(feedWorkflowMetadata.allowedSender);
          }
          if (feedWorkflowMetadata.allowedWorkflowOwner == address(0)) {
            revert InvalidAddress(feedWorkflowMetadata.allowedWorkflowOwner);
          }
          if (feedWorkflowMetadata.allowedWorkflowName == bytes10(0)) {
            revert InvalidWorkflowName(feedWorkflowMetadata.allowedWorkflowName);
          }
        }

        bytes32 reportHash = _createReportHash(
          dataId,
          feedWorkflowMetadata.allowedSender,
          feedWorkflowMetadata.allowedWorkflowOwner,
          feedWorkflowMetadata.allowedWorkflowName
        );
        s_writePermissions[reportHash] = true;
        feedConfig.workflowMetadata.push(feedWorkflowMetadata);
      }

      feedConfig.description = descriptions[i];

      emit DecimalFeedConfigSet({
        dataId: dataId,
        decimals: _getDecimals(dataId),
        description: descriptions[i],
        workflowMetadata: workflowMetadata
      });
    }
  }

  /// @notice Initializes the config for a bundle feed
  /// @param dataIds The data IDs of the feeds to configure
  /// @param descriptions The descriptions of the feeds
  /// @param decimalsMatrix The number of decimals for each data point in the bundle for the feed
  /// @param workflowMetadata List of workflow metadata (owners, senders, and names) for every feed
  function setBundleFeedConfigs(
    bytes16[] calldata dataIds,
    string[] calldata descriptions,
    uint8[][] calldata decimalsMatrix,
    WorkflowMetadata[] calldata workflowMetadata
  ) external onlyFeedAdmin {
    if (workflowMetadata.length == 0 || dataIds.length == 0) {
      revert EmptyConfig();
    }

    if (dataIds.length != descriptions.length || dataIds.length != decimalsMatrix.length) {
      revert ArrayLengthMismatch();
    }

    for (uint256 i; i < dataIds.length; ++i) {
      bytes16 dataId = dataIds[i];
      if (dataId == bytes16(0)) revert InvalidDataId();
      FeedConfig storage feedConfig = s_feedConfigs[dataId];

      if (feedConfig.workflowMetadata.length > 0) {
        // Feed is already configured, remove the previous config
        for (uint256 j; j < feedConfig.workflowMetadata.length; ++j) {
          WorkflowMetadata memory feedCurrentWorkflowMetadata = feedConfig.workflowMetadata[j];
          bytes32 reportHash = _createReportHash(
            dataId,
            feedCurrentWorkflowMetadata.allowedSender,
            feedCurrentWorkflowMetadata.allowedWorkflowOwner,
            feedCurrentWorkflowMetadata.allowedWorkflowName
          );
          delete s_writePermissions[reportHash];
        }

        delete s_feedConfigs[dataId];

        emit FeedConfigRemoved(dataId);
      }

      for (uint256 j; j < workflowMetadata.length; ++j) {
        WorkflowMetadata memory feedWorkflowMetadata = workflowMetadata[j];
        // Do those checks only once for the first data id
        if (i == 0) {
          if (feedWorkflowMetadata.allowedSender == address(0)) {
            revert InvalidAddress(feedWorkflowMetadata.allowedSender);
          }
          if (feedWorkflowMetadata.allowedWorkflowOwner == address(0)) {
            revert InvalidAddress(feedWorkflowMetadata.allowedWorkflowOwner);
          }
          if (feedWorkflowMetadata.allowedWorkflowName == bytes10(0)) {
            revert InvalidWorkflowName(feedWorkflowMetadata.allowedWorkflowName);
          }
        }

        bytes32 reportHash = _createReportHash(
          dataId,
          feedWorkflowMetadata.allowedSender,
          feedWorkflowMetadata.allowedWorkflowOwner,
          feedWorkflowMetadata.allowedWorkflowName
        );
        s_writePermissions[reportHash] = true;
        feedConfig.workflowMetadata.push(feedWorkflowMetadata);
      }

      feedConfig.bundleDecimals = decimalsMatrix[i];
      feedConfig.description = descriptions[i];

      emit BundleFeedConfigSet({
        dataId: dataId,
        decimals: decimalsMatrix[i],
        description: descriptions[i],
        workflowMetadata: workflowMetadata
      });
    }
  }

  /// @notice Removes feeds and all associated data, for a set of feeds
  /// @param dataIds And array of data IDs to delete the data and configs of
  function removeFeedConfigs(
    bytes16[] calldata dataIds
  ) external onlyFeedAdmin {
    for (uint256 i; i < dataIds.length; ++i) {
      bytes16 dataId = dataIds[i];
      if (s_feedConfigs[dataId].workflowMetadata.length == 0) revert FeedNotConfigured(dataId);

      for (uint256 j; j < s_feedConfigs[dataId].workflowMetadata.length; ++j) {
        WorkflowMetadata memory feedWorkflowMetadata = s_feedConfigs[dataId].workflowMetadata[j];
        bytes32 reportHash = _createReportHash(
          dataId,
          feedWorkflowMetadata.allowedSender,
          feedWorkflowMetadata.allowedWorkflowOwner,
          feedWorkflowMetadata.allowedWorkflowName
        );
        delete s_writePermissions[reportHash];
      }

      delete s_feedConfigs[dataId];

      emit FeedConfigRemoved(dataId);
    }
  }

  /// @notice Sets a feed admin for all feeds, only callable by the Owner
  /// @param feedAdmin The feed admin
  function setFeedAdmin(address feedAdmin, bool isAdmin) external onlyOwner {
    if (feedAdmin == address(0)) revert InvalidAddress(feedAdmin);

    s_feedAdmins[feedAdmin] = isAdmin;
    emit FeedAdminSet(feedAdmin, isAdmin);
  }

  /// @notice Returns a bool is an address has feed admin permission for all feeds
  /// @param feedAdmin The feed admin
  /// @return isFeedAdmin bool if the address is the feed admin for all feeds
  function isFeedAdmin(
    address feedAdmin
  ) external view returns (bool) {
    return s_feedAdmins[feedAdmin];
  }

  /// @inheritdoc IDataFeedsCache
  function updateDataIdMappingsForProxies(
    address[] calldata proxies,
    bytes16[] calldata dataIds
  ) external onlyFeedAdmin {
    uint256 numberOfProxies = proxies.length;
    if (numberOfProxies != dataIds.length) revert ArrayLengthMismatch();

    for (uint256 i; i < numberOfProxies; i++) {
      s_aggregatorProxyToDataId[proxies[i]] = dataIds[i];

      emit ProxyDataIdUpdated(proxies[i], dataIds[i]);
    }
  }

  /// @inheritdoc IDataFeedsCache
  function getDataIdForProxy(
    address proxy
  ) external view returns (bytes16 dataId) {
    return s_aggregatorProxyToDataId[proxy];
  }

  /// @inheritdoc IDataFeedsCache
  function removeDataIdMappingsForProxies(
    address[] calldata proxies
  ) external onlyFeedAdmin {
    uint256 numberOfProxies = proxies.length;

    for (uint256 i; i < numberOfProxies; i++) {
      address proxy = proxies[i];
      bytes16 dataId = s_aggregatorProxyToDataId[proxy];
      delete s_aggregatorProxyToDataId[proxy];
      emit ProxyDataIdRemoved(proxy, dataId);
    }
  }

  // ================================================================
  // │                   Token Transfer Interface                   │
  // ================================================================

  /// @inheritdoc ITokenRecover
  function recoverTokens(IERC20 token, address to, uint256 amount) external onlyOwner {
    if (address(token) == address(0)) {
      if (amount > address(this).balance) {
        revert InsufficientBalance(address(this).balance, amount);
      }
      (bool success, bytes memory data) = to.call{value: amount}("");
      if (!success) revert ErrorSendingNative(to, amount, data);
    } else {
      if (amount > token.balanceOf(address(this))) {
        revert InsufficientBalance(token.balanceOf(address(this)), amount);
      }
      token.safeTransfer(to, amount);
    }
    emit TokenRecovered(address(token), to, amount);
  }

  // ================================================================
  // │                    Cache Update Interface                    │
  // ================================================================

  /// @inheritdoc IReceiver
  function onReport(bytes calldata metadata, bytes calldata report) external {
    (address workflowOwner, bytes10 workflowName) = _getWorkflowMetaData(metadata);

    // The first 32 bytes is the offset to the array
    // The second 32 bytes is the length of the array
    uint256 numReports = uint256(bytes32(report[32:64]));

    // Decimal reports contain 96 bytes per report
    // The total length should equal to the sum of:
    // 32 bytes for the offset
    // 32 bytes for the number of reports
    // the number of reports times 96
    if (report.length == numReports * 96 + 64) {
      ReceivedDecimalReport[] memory decodedDecimalReports = abi.decode(report, (ReceivedDecimalReport[]));
      for (uint256 i; i < numReports; ++i) {
        ReceivedDecimalReport memory decodedDecimalReport = decodedDecimalReports[i];
        // single dataId can have multiple permissions, to be updated by multiple Workflows
        bytes16 dataId = bytes16(decodedDecimalReport.dataId);
        bytes32 permission = _createReportHash(dataId, msg.sender, workflowOwner, workflowName);
        if (!s_writePermissions[permission]) {
          emit InvalidUpdatePermission(dataId, msg.sender, workflowOwner, workflowName);
          continue;
        }

        if (decodedDecimalReport.timestamp <= s_latestDecimalReports[dataId].timestamp) {
          emit StaleDecimalReport(dataId, decodedDecimalReport.timestamp, s_latestDecimalReports[dataId].timestamp);
          continue;
        }

        StoredDecimalReport memory decimalReport =
          StoredDecimalReport({answer: decodedDecimalReport.answer, timestamp: decodedDecimalReport.timestamp});

        uint256 roundId = ++s_dataIdToRoundId[dataId];

        s_latestDecimalReports[dataId] = decimalReport;
        s_decimalReports[roundId][dataId] = decimalReport;

        emit DecimalReportUpdated(dataId, roundId, decimalReport.timestamp, decimalReport.answer);

        // Needed for DF1 backward compatibility
        emit NewRound(roundId, address(0), decodedDecimalReport.timestamp);
        emit AnswerUpdated(int256(uint256(decodedDecimalReport.answer)), roundId, block.timestamp);
      }
    }
    // Bundle reports contain more bytes for the offsets
    // The total length should equal to the sum of:
    // 32 bytes for the offset
    // 32 bytes for the number of reports
    // the number of reports times 224
    else {
      //For byte reports decode using ReceivedFeedReportBundle struct
      ReceivedBundleReport[] memory decodedBundleReports = abi.decode(report, (ReceivedBundleReport[]));
      for (uint256 i; i < decodedBundleReports.length; ++i) {
        ReceivedBundleReport memory decodedBundleReport = decodedBundleReports[i];
        bytes16 dataId = bytes16(decodedBundleReport.dataId);
        // same dataId can have multiple permissions
        bytes32 permission = _createReportHash(dataId, msg.sender, workflowOwner, workflowName);
        if (!s_writePermissions[permission]) {
          emit InvalidUpdatePermission(dataId, msg.sender, workflowOwner, workflowName);
          continue;
        }

        if (decodedBundleReport.timestamp <= s_latestBundleReports[dataId].timestamp) {
          emit StaleBundleReport(dataId, decodedBundleReport.timestamp, s_latestBundleReports[dataId].timestamp);
          continue;
        }

        StoredBundleReport memory bundleReport =
          StoredBundleReport({bundle: decodedBundleReport.bundle, timestamp: decodedBundleReport.timestamp});

        s_latestBundleReports[dataId] = bundleReport;

        emit BundleReportUpdated(dataId, bundleReport.timestamp, bundleReport.bundle);
      }
    }
  }
  // ================================================================
  // │                        Helper Methods                        │
  // ================================================================

  /// @notice Gets the Decimals of the feed from the data Id
  /// @param dataId The data ID for the feed
  /// @return feedDecimals The number of decimals the feed has
  function _getDecimals(
    bytes16 dataId
  ) internal pure returns (uint8 feedDecimals) {
    // Get the report type from data id. Report type has index of 7
    bytes1 reportType = _getDataType(dataId, 7);

    // For decimal reports convert to uint8, then shift
    if (reportType >= hex"20" && reportType <= hex"60") {
      return uint8(reportType) - 32;
    }

    // If not decimal type, return 0
    return 0;
  }

  /// @notice Extracts the workflow name and the workflow owner from the metadata parameter of onReport
  /// @param metadata The metadata in bytes format
  /// @return workflowOwner The owner of the workflow
  /// @return workflowName  The name of the workflow
  function _getWorkflowMetaData(
    bytes memory metadata
  ) internal pure returns (address, bytes10) {
    address workflowOwner;
    bytes10 workflowName;
    // (first 32 bytes contain length of the byte array)
    // workflow_cid             // offset 32, size 32
    // workflow_name            // offset 64, size 10
    // workflow_owner           // offset 74, size 20
    // report_name              // offset 94, size  2
    assembly {
      // no shifting needed for bytes10 type
      workflowName := mload(add(metadata, 64))
      // shift right by 12 bytes to get the actual value
      workflowOwner := shr(mul(12, 8), mload(add(metadata, 74)))
    }
    return (workflowOwner, workflowName);
  }

  /// @notice Extracts a byte from the data ID, to check data types
  /// @param dataId The data ID for the feed
  /// @param index The index of the byte to extract from the data Id
  /// @return dataType result The keccak256 hash of the abi.encoded inputs
  function _getDataType(bytes16 dataId, uint256 index) internal pure returns (bytes1 dataType) {
    // Convert bytes16 to bytes
    return abi.encodePacked(dataId)[index];
  }

  /// @notice Creates a report hash used to permission write access
  /// @param dataId The data ID for the feed
  /// @param sender The msg.sender of the transaction calling into onReport
  /// @param workflowOwner The owner of the workflow
  /// @param workflowName The name of the workflow
  /// @return reportHash The keccak256 hash of the abi.encoded inputs
  function _createReportHash(
    bytes16 dataId,
    address sender,
    address workflowOwner,
    bytes10 workflowName
  ) internal pure returns (bytes32) {
    return keccak256(abi.encode(dataId, sender, workflowOwner, workflowName));
  }

  // ================================================================
  // │                    Data Access Interface                     │
  // ================================================================

  /// Bundle Feed Interface

  function latestBundle() external view returns (bytes memory bundle) {
    bytes16 dataId = s_aggregatorProxyToDataId[msg.sender];
    if (dataId == bytes16(0)) revert NoMappingForSender(msg.sender);

    return (s_latestBundleReports[dataId].bundle);
  }

  function bundleDecimals() external view returns (uint8[] memory bundleFeedDecimals) {
    bytes16 dataId = s_aggregatorProxyToDataId[msg.sender];
    if (dataId == bytes16(0)) revert NoMappingForSender(msg.sender);

    return s_feedConfigs[dataId].bundleDecimals;
  }

  function latestBundleTimestamp() external view returns (uint256 timestamp) {
    bytes16 dataId = s_aggregatorProxyToDataId[msg.sender];
    if (dataId == bytes16(0)) revert NoMappingForSender(msg.sender);

    return s_latestBundleReports[dataId].timestamp;
  }

  /// AggregatorInterface

  function latestAnswer() external view returns (int256 answer) {
    bytes16 dataId = s_aggregatorProxyToDataId[msg.sender];
    if (dataId == bytes16(0)) revert NoMappingForSender(msg.sender);

    return int256(uint256(s_latestDecimalReports[dataId].answer));
  }

  function latestTimestamp() external view returns (uint256 timestamp) {
    bytes16 dataId = s_aggregatorProxyToDataId[msg.sender];
    if (dataId == bytes16(0)) revert NoMappingForSender(msg.sender);

    return s_latestDecimalReports[dataId].timestamp;
  }

  function latestRound() external view returns (uint256 round) {
    bytes16 dataId = s_aggregatorProxyToDataId[msg.sender];
    if (dataId == bytes16(0)) revert NoMappingForSender(msg.sender);

    return s_dataIdToRoundId[dataId];
  }

  function getAnswer(
    uint256 roundId
  ) external view returns (int256 answer) {
    bytes16 dataId = s_aggregatorProxyToDataId[msg.sender];
    if (dataId == bytes16(0)) revert NoMappingForSender(msg.sender);

    return int256(uint256(s_decimalReports[roundId][dataId].answer));
  }

  function getTimestamp(
    uint256 roundId
  ) external view returns (uint256 timestamp) {
    bytes16 dataId = s_aggregatorProxyToDataId[msg.sender];
    if (dataId == bytes16(0)) revert NoMappingForSender(msg.sender);

    return s_decimalReports[roundId][dataId].timestamp;
  }

  /// AggregatorV3Interface

  function decimals() external view returns (uint8 feedDecimals) {
    bytes16 dataId = s_aggregatorProxyToDataId[msg.sender];
    if (dataId == bytes16(0)) revert NoMappingForSender(msg.sender);
    return _getDecimals(dataId);
  }

  function description() external view returns (string memory feedDescription) {
    bytes16 dataId = s_aggregatorProxyToDataId[msg.sender];
    if (dataId == bytes16(0)) revert NoMappingForSender(msg.sender);

    return s_feedConfigs[dataId].description;
  }

  function getRoundData(
    uint80 roundId
  ) external view returns (uint80 id, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) {
    bytes16 dataId = s_aggregatorProxyToDataId[msg.sender];
    if (dataId == bytes16(0)) revert NoMappingForSender(msg.sender);

    uint256 timestamp = s_decimalReports[uint256(roundId)][dataId].timestamp;

    return (roundId, int256(uint256(s_decimalReports[uint256(roundId)][dataId].answer)), timestamp, timestamp, roundId);
  }

  function latestRoundData()
    external
    view
    returns (uint80 id, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    bytes16 dataId = s_aggregatorProxyToDataId[msg.sender];
    if (dataId == bytes16(0)) revert NoMappingForSender(msg.sender);

    uint80 roundId = uint80(s_dataIdToRoundId[dataId]);
    uint256 timestamp = s_latestDecimalReports[dataId].timestamp;

    return (roundId, int256(uint256(s_latestDecimalReports[dataId].answer)), timestamp, timestamp, roundId);
  }

  /// Direct access
  function getLatestBundle(
    bytes16 dataId
  ) external view returns (bytes memory bundle) {
    if (dataId == bytes16(0)) revert InvalidDataId();
    return (s_latestBundleReports[dataId].bundle);
  }

  function getBundleDecimals(
    bytes16 dataId
  ) external view returns (uint8[] memory bundleFeedDecimals) {
    if (dataId == bytes16(0)) revert InvalidDataId();
    return s_feedConfigs[dataId].bundleDecimals;
  }

  function getLatestBundleTimestamp(
    bytes16 dataId
  ) external view returns (uint256 timestamp) {
    if (dataId == bytes16(0)) revert InvalidDataId();
    return s_latestBundleReports[dataId].timestamp;
  }

  function getLatestAnswer(
    bytes16 dataId
  ) external view returns (int256 answer) {
    if (dataId == bytes16(0)) revert InvalidDataId();
    return int256(uint256(s_latestDecimalReports[dataId].answer));
  }

  function getLatestTimestamp(
    bytes16 dataId
  ) external view returns (uint256 timestamp) {
    if (dataId == bytes16(0)) revert InvalidDataId();
    return s_latestDecimalReports[dataId].timestamp;
  }

  function getLatestRoundData(
    bytes16 dataId
  ) external view returns (uint80 id, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) {
    if (dataId == bytes16(0)) revert InvalidDataId();

    uint80 roundId = uint80(s_dataIdToRoundId[dataId]);
    uint256 timestamp = s_latestDecimalReports[dataId].timestamp;

    return (roundId, int256(uint256(s_latestDecimalReports[dataId].answer)), timestamp, timestamp, roundId);
  }

  function getDecimals(
    bytes16 dataId
  ) external pure returns (uint8 feedDecimals) {
    if (dataId == bytes16(0)) revert InvalidDataId();
    return _getDecimals(dataId);
  }

  function getDescription(
    bytes16 dataId
  ) external view returns (string memory feedDescription) {
    if (dataId == bytes16(0)) revert InvalidDataId();
    return s_feedConfigs[dataId].description;
  }
}
