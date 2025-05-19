// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";
import {ITokenMessenger} from "../../pools/USDC/ITokenMessenger.sol";

// This contract mocks both the ITokenMessenger and IMessageTransmitter
// contracts involved with the Cross Chain Token Protocol.
contract MockUSDCTokenMessenger is ITokenMessenger {
  uint32 private immutable i_messageBodyVersion;
  address private immutable i_transmitter;

  bytes32 public constant DESTINATION_TOKEN_MESSENGER = keccak256("i_destinationTokenMessenger");

  uint64 public s_nonce;

  constructor(uint32 version, address transmitter) {
    i_messageBodyVersion = version;
    s_nonce = 1;
    i_transmitter = transmitter;
  }

  function depositForBurnWithCaller(
    uint256 amount,
    uint32 destinationDomain,
    bytes32 mintRecipient,
    address burnToken,
    bytes32 destinationCaller
  ) external returns (uint64) {
    IBurnMintERC20(burnToken).transferFrom(msg.sender, address(this), amount);
    IBurnMintERC20(burnToken).burn(amount);
    emit DepositForBurn(
      s_nonce,
      burnToken,
      amount,
      msg.sender,
      mintRecipient,
      destinationDomain,
      DESTINATION_TOKEN_MESSENGER,
      destinationCaller
    );
    return s_nonce++;
  }

  function messageBodyVersion() external view returns (uint32) {
    return i_messageBodyVersion;
  }

  function localMessageTransmitter() external view returns (address) {
    return i_transmitter;
  }
}
