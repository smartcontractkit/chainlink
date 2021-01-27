pragma solidity ^0.6.0;

/**
 * @title The Owned contract
 * @notice A contract with helpers for basic contract ownership.
 */
contract Owned {

  address private __owner;
  address private __pendingOwner;

  event OwnershipTransferRequested(
    address indexed from,
    address indexed to
  );
  event OwnershipTransferred(
    address indexed from,
    address indexed to
  );

  constructor() public {
    _setOwner(msg.sender);
  }

  /**
   * @return The address of the owner.
   */
  function owner() public view returns (address) {
    return _owner();
  }

  /**
   * @return The owner slot.
   */
  function _owner() internal virtual view returns (address) {
    return __owner;
  }

  /**
   * @dev Sets the address of the owner.
   * @param newOwner Address of the new owner.
   */
  function _setOwner(address newOwner) internal virtual {
    __owner = newOwner;
  }

  /**
   * @return The pending owner slot.
   */
  function _pendingOwner() internal virtual view returns (address) {
    return __pendingOwner;
  }

  /**
   * @dev Sets the address of the pending owner.
   * @param newPendingOwner Address of the new pending owner.
   */
  function _setPendingOwner(address newPendingOwner) internal virtual {
    __pendingOwner = newPendingOwner;
  }

  /**
   * @dev Allows an owner to begin transferring ownership to a new address,
   * pending.
   */
  function transferOwnership(address _to) external onlyOwner() {
    _setPendingOwner(_to);

    emit OwnershipTransferRequested(_owner(), _to);
  }

  /**
   * @dev Allows an ownership transfer to be completed by the recipient.
   */
  function acceptOwnership() external {
    require(msg.sender == _pendingOwner(), "Must be proposed owner");

    address oldOwner = _owner();
    _setOwner(msg.sender);
    _setPendingOwner(address(0));

    emit OwnershipTransferred(oldOwner, msg.sender);
  }

  /**
   * @dev Reverts if called by anyone other than the contract owner.
   */
  modifier onlyOwner() {
    require(msg.sender == _owner(), "Only callable by owner");
    _;
  }
}
