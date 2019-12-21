pragma solidity 0.5.0;

/**
 * @title The Owned contract
 * @notice A contract with helpers for basic contract ownership.
 */
contract Owned {

  address public owner;
  address private pendingOwner;

  event OwnershipTransferRequested(
    address indexed from,
    address indexed to
  );
  event OwnershipTransfered(
    address indexed from,
    address indexed to
  );
  event OwnershipRenounced(
    address indexed by
  );

  constructor() public {
    owner = msg.sender;
  }

  /**
   * @dev Allows an owner to begin transferring ownership to a new address,
   * pending.
   */
  function transferOwnership(address _to)
    external
    onlyOwner()
  {
    pendingOwner = _to;

    emit OwnershipTransferRequested(owner, _to);
  }

  /**
   * @dev Allows an ownership transfer to be completed by the recipient.
   */
  function acceptOwnership()
    external
  {
    require(msg.sender == pendingOwner, "Must be requested to accept ownership");

    address oldOwner = owner;
    owner = msg.sender;
    pendingOwner = address(0);

    emit OwnershipTransfered(oldOwner, msg.sender);
  }

  /**
   * @dev Renounces ownership so that no one owns the contract.
   */
  function renounceOwnership()
    external
    onlyOwner()
  {
    address oldOwner = owner;
    owner = address(0);

    emit OwnershipTransfered(oldOwner, address(0));
  }

  /**
   * @dev Reverts if called by anyone other than the contract owner.
   */
  modifier onlyOwner() {
    require(msg.sender == owner, "Only callable by owner");
    _;
  }

}
