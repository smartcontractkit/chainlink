// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

contract KeeperRegistrar1_2Mock {
  event AutoApproveAllowedSenderSet(address indexed senderAddress, bool allowed);
  event ConfigChanged(
    uint8 autoApproveConfigType,
    uint32 autoApproveMaxAllowed,
    address keeperRegistry,
    uint96 minLINKJuels
  );
  event OwnershipTransferRequested(address indexed from, address indexed to);
  event OwnershipTransferred(address indexed from, address indexed to);
  event RegistrationApproved(bytes32 indexed hash, string displayName, uint256 indexed upkeepId);
  event RegistrationRejected(bytes32 indexed hash);
  event RegistrationRequested(
    bytes32 indexed hash,
    string name,
    bytes encryptedEmail,
    address indexed upkeepContract,
    uint32 gasLimit,
    address adminAddress,
    bytes checkData,
    uint96 amount,
    uint8 indexed source
  );

  function emitAutoApproveAllowedSenderSet(address senderAddress, bool allowed) public {
    emit AutoApproveAllowedSenderSet(senderAddress, allowed);
  }

  function emitConfigChanged(
    uint8 autoApproveConfigType,
    uint32 autoApproveMaxAllowed,
    address keeperRegistry,
    uint96 minLINKJuels
  ) public {
    emit ConfigChanged(autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels);
  }

  function emitOwnershipTransferRequested(address from, address to) public {
    emit OwnershipTransferRequested(from, to);
  }

  function emitOwnershipTransferred(address from, address to) public {
    emit OwnershipTransferred(from, to);
  }

  function emitRegistrationApproved(bytes32 hash, string memory displayName, uint256 upkeepId) public {
    emit RegistrationApproved(hash, displayName, upkeepId);
  }

  function emitRegistrationRejected(bytes32 hash) public {
    emit RegistrationRejected(hash);
  }

  function emitRegistrationRequested(
    bytes32 hash,
    string memory name,
    bytes memory encryptedEmail,
    address upkeepContract,
    uint32 gasLimit,
    address adminAddress,
    bytes memory checkData,
    uint96 amount,
    uint8 source
  ) public {
    emit RegistrationRequested(
      hash,
      name,
      encryptedEmail,
      upkeepContract,
      gasLimit,
      adminAddress,
      checkData,
      amount,
      source
    );
  }

  enum AutoApproveType {
    DISABLED,
    ENABLED_SENDER_ALLOWLIST,
    ENABLED_ALL
  }

  AutoApproveType public s_autoApproveConfigType;
  uint32 public s_autoApproveMaxAllowed;
  uint32 public s_approvedCount;
  address public s_keeperRegistry;
  uint256 public s_minLINKJuels;

  // Function to set mock return data for the getRegistrationConfig function
  function setRegistrationConfig(
    AutoApproveType _autoApproveConfigType,
    uint32 _autoApproveMaxAllowed,
    uint32 _approvedCount,
    address _keeperRegistry,
    uint256 _minLINKJuels
  ) external {
    s_autoApproveConfigType = _autoApproveConfigType;
    s_autoApproveMaxAllowed = _autoApproveMaxAllowed;
    s_approvedCount = _approvedCount;
    s_keeperRegistry = _keeperRegistry;
    s_minLINKJuels = _minLINKJuels;
  }

  // Mock getRegistrationConfig function
  function getRegistrationConfig()
    external
    view
    returns (
      AutoApproveType autoApproveConfigType,
      uint32 autoApproveMaxAllowed,
      uint32 approvedCount,
      address keeperRegistry,
      uint256 minLINKJuels
    )
  {
    return (s_autoApproveConfigType, s_autoApproveMaxAllowed, s_approvedCount, s_keeperRegistry, s_minLINKJuels);
  }
}
