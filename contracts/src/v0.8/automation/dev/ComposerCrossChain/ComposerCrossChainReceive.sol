pragma solidity 0.8.16;

import "../../../shared/access/ConfirmedOwner.sol";
import "./CCIPDeps/CCIPReceiver.sol";
import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/token/ERC20/IERC20.sol";

contract ComposerCrossChainReceive is ConfirmedOwner, CCIPReceiver {
  IERC20 s_token;
  uint256 s_nonce;
  uint256 s_minBalance;

  constructor(address _router, address _token) ConfirmedOwner(msg.sender) CCIPReceiver(_router) {
    s_token = IERC20(_token);
    s_minBalance = 1; // for now, hardcoded to 1 wei-worth of a token
  }

  // If balance falls below min balance, request top-up. Include current nonce in return tuple.
  function getStatus() external view returns (bool updateNeeded, uint256 nonce) {
    return (s_token.balanceOf(address(this)) < s_minBalance, s_nonce);
  }

  function _ccipReceive(Client.Any2EVMMessage memory message) internal override {
    s_nonce++; // prevent double spend from sender
  }

  // admin
  function withdraw(address to) external onlyOwner {
    s_token.transfer(to, s_token.balanceOf(address(this)));
  }

  function resetNonce() external onlyOwner {
    s_nonce = 0;
  }
}
