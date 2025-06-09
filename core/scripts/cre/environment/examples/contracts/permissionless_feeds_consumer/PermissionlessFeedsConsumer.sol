// SPDX-License-Identifier: MIT
pragma solidity =0.8.24 ^0.8.0;

// core/scripts/cre/environment/examples/contracts/dependencies/interfaces/IOwnable.sol

interface IOwnable {
  function owner() external returns (address);

  function transferOwnership(address recipient) external;

  function acceptOwnership() external;
}

// core/scripts/cre/environment/examples/contracts/dependencies/vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol

// OpenZeppelin Contracts v4.4.1 (utils/introspection/IERC165.sol)

/**
 * @dev Interface of the ERC165 standard, as defined in the
 * https://eips.ethereum.org/EIPS/eip-165[EIP].
 *
 * Implementers can declare support of contract interfaces, which can then be
 * queried by others ({ERC165Checker}).
 *
 * For an implementation, see {ERC165}.
 */
interface IERC165 {
    /**
     * @dev Returns true if this contract implements the interface defined by
     * `interfaceId`. See the corresponding
     * https://eips.ethereum.org/EIPS/eip-165#how-interfaces-are-identified[EIP section]
     * to learn more about how these ids are created.
     *
     * This function call must use less than 30 000 gas.
     */
    function supportsInterface(bytes4 interfaceId) external view returns (bool);
}

// core/scripts/cre/environment/examples/contracts/dependencies/access/ConfirmedOwnerWithProposal.sol

/// @title The ConfirmedOwner contract
/// @notice A contract with helpers for basic contract ownership.
contract ConfirmedOwnerWithProposal is IOwnable {
  address private s_owner;
  address private s_pendingOwner;

  event OwnershipTransferRequested(address indexed from, address indexed to);
  event OwnershipTransferred(address indexed from, address indexed to);

  constructor(address newOwner, address pendingOwner) {
    // solhint-disable-next-line gas-custom-errors
    require(newOwner != address(0), "Cannot set owner to zero");

    s_owner = newOwner;
    if (pendingOwner != address(0)) {
      _transferOwnership(pendingOwner);
    }
  }

  /// @notice Allows an owner to begin transferring ownership to a new address.
  function transferOwnership(address to) public override onlyOwner {
    _transferOwnership(to);
  }

  /// @notice Allows an ownership transfer to be completed by the recipient.
  function acceptOwnership() external override {
    // solhint-disable-next-line gas-custom-errors
    require(msg.sender == s_pendingOwner, "Must be proposed owner");

    address oldOwner = s_owner;
    s_owner = msg.sender;
    s_pendingOwner = address(0);

    emit OwnershipTransferred(oldOwner, msg.sender);
  }

  /// @notice Get the current owner
  function owner() public view override returns (address) {
    return s_owner;
  }

  /// @notice validate, transfer ownership, and emit relevant events
  function _transferOwnership(address to) private {
    // solhint-disable-next-line gas-custom-errors
    require(to != msg.sender, "Cannot transfer to self");

    s_pendingOwner = to;

    emit OwnershipTransferRequested(s_owner, to);
  }

  /// @notice validate access
  function _validateOwnership() internal view {
    // solhint-disable-next-line gas-custom-errors
    require(msg.sender == s_owner, "Only callable by owner");
  }

  /// @notice Reverts if called by anyone other than the contract owner.
  modifier onlyOwner() {
    _validateOwnership();
    _;
  }
}

// core/scripts/cre/environment/examples/contracts/dependencies/interfaces/IReceiver.sol

/// @title IReceiver - receives keystone reports
/// @notice Implementations must support the IReceiver interface through ERC165.
interface IReceiver is IERC165 {
  /// @notice Handles incoming keystone reports.
  /// @dev If this function call reverts, it can be retried with a higher gas
  /// limit. The receiver is responsible for discarding stale reports.
  /// @param metadata Report's metadata.
  /// @param report Workflow report.
  function onReport(bytes calldata metadata, bytes calldata report) external;
}

// core/scripts/cre/environment/examples/contracts/dependencies/vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol

// OpenZeppelin Contracts v4.4.1 (interfaces/IERC165.sol)

// core/scripts/cre/environment/examples/contracts/dependencies/access/ConfirmedOwner.sol

/// @title The ConfirmedOwner contract
/// @notice A contract with helpers for basic contract ownership.
contract ConfirmedOwner is ConfirmedOwnerWithProposal {
  constructor(address newOwner) ConfirmedOwnerWithProposal(newOwner, address(0)) {}
}

// core/scripts/cre/environment/examples/contracts/dependencies/access/OwnerIsCreator.sol

/// @title The OwnerIsCreator contract
/// @notice A contract with helpers for basic contract ownership.
contract OwnerIsCreator is ConfirmedOwner {
  constructor() ConfirmedOwner(msg.sender) {}
}

// core/scripts/cre/environment/examples/contracts/permissionless_feeds_consumer/PermissionlessFeedsConsumer.sol

contract PermissionlessFeedsConsumer is IReceiver, OwnerIsCreator {
  event FeedReceived(bytes32 indexed feedId, uint224 price, uint32 timestamp);

  struct ReceivedFeedReport {
    bytes32 FeedId;
    uint32 Timestamp;
    uint224 Price;
  }

  struct StoredFeedReport {
    uint224 Price;
    uint32 Timestamp;
  }

  mapping(bytes32 feedId => StoredFeedReport feedReport) internal s_feedReports;

  function onReport(bytes calldata metadata, bytes calldata rawReport) external {
    ReceivedFeedReport[] memory feeds = abi.decode(rawReport, (ReceivedFeedReport[]));
    for (uint256 i = 0; i < feeds.length; ++i) {
      s_feedReports[feeds[i].FeedId] = StoredFeedReport(feeds[i].Price, feeds[i].Timestamp);
      emit FeedReceived(feeds[i].FeedId, feeds[i].Price, feeds[i].Timestamp);
    }
  }

  function getPrice(bytes32 feedId) external view returns (uint224, uint32) {
    StoredFeedReport memory report = s_feedReports[feedId];
    return (report.Price, report.Timestamp);
  }

  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == type(IReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
  }
}

