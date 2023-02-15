// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {SafeMath} from "@openzeppelin/contracts/utils/math/SafeMath.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract MetaERC20 is IERC20 {
  using SafeMath for uint256;

  string public constant name = "BankToken";
  string public constant symbol = "BANKTOKEN";
  uint8 public constant decimals = 18;

  // NOTE: implements IERC20.totalSupply
  uint256 public override totalSupply;
  // NOTE: implements IERC20.balanceOf
  mapping(address => uint256) public override balanceOf;
  // NOTE: implements IERC20.allowance
  mapping(address => mapping(address => uint256)) public override allowance;

  bytes32 public DOMAIN_SEPARATOR;
  // keccak256("MetaTransfer(address owner, address to, uint256 amount,uint256 nonce,uint256 deadline)");
  bytes32 public constant META_TRANSFER_TYPEHASH = 0xfc3a30ed0a6a26bdf760a234a365b51a1cb10009a8fba3cb68ad3b45b789aa17;
  mapping(address => uint256) public nonces;

  constructor(uint256 _totalSupply) public {
    totalSupply = _totalSupply;
    balanceOf[msg.sender] = totalSupply;

    DOMAIN_SEPARATOR = keccak256(
      abi.encode(
        keccak256('EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)'),
        keccak256(bytes(name)),
        keccak256(bytes('1')),
        block.chainid,
        address(this)
      )
    );
  }

  /**
   * @dev Moves `amount` tokens from the caller's account to `to`.
   *
   * Returns a boolean value indicating whether the operation succeeded.
   *
   * Emits a {Transfer} event.
   */
  function transfer(address to, uint256 amount) external override returns (bool) {
    _transfer(msg.sender, to, amount);
    return true;
  }

  function _transfer(address from, address to, uint256 amount) private {
    balanceOf[from] = balanceOf[from].sub(amount);
    balanceOf[to] = balanceOf[to].add(amount);
    emit Transfer(from, to, amount);
  }

  function _approve(address owner, address spender, uint amount) private {
    allowance[owner][spender] = amount;
    emit Approval(owner, spender, amount);
  }

  /**
   * @dev Sets `amount` as the allowance of `spender` over the caller's tokens.
   *
   * Returns a boolean value indicating whether the operation succeeded.
   *
   * IMPORTANT: Beware that changing an allowance with this method brings the risk
   * that someone may use both the old and the new allowance by unfortunate
   * transaction ordering. One possible solution to mitigate this race
   * condition is to first reduce the spender's allowance to 0 and set the
   * desired value afterwards:
   * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
   *
   * Emits an {Approval} event.
   */
  function approve(address spender, uint256 amount) external override returns (bool) {
    _approve(msg.sender, spender, amount);
    return true;
  }

  /**
   * @dev Moves `amount` tokens from `from` to `to` using the
   * allowance mechanism. `amount` is then deducted from the caller's
   * allowance.
   *
   * Returns a boolean value indicating whether the operation succeeded.
   *
   * Emits a {Transfer} event.
   */
  function transferFrom(
    address from,
    address to,
    uint256 amount
  ) external override returns (bool) {
    if (allowance[from][msg.sender] != type(uint256).max) {
      allowance[from][msg.sender] = allowance[from][msg.sender].sub(amount);
    }
    _transfer(from, to, amount);
    return true;
  }

  function metaTransfer(address owner, address to, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) external {
    require(deadline >= block.timestamp, 'EXPIRED');
    bytes32 digest = keccak256(
      abi.encodePacked(
        '\x19\x01',
        DOMAIN_SEPARATOR,
        keccak256(abi.encode(META_TRANSFER_TYPEHASH, owner, to, amount, nonces[owner]++, deadline))
      )
    );
    address recoveredAddress = ecrecover(digest, v, r, s);
    require(recoveredAddress != address(0) && recoveredAddress == owner, 'INVALID_SIGNATURE');
    _transfer(owner, to, amount);
  }
}
