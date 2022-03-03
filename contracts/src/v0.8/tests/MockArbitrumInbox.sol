import {IInbox} from "../dev/vendor/arb-bridge-eth/v0.8.0-custom/contracts/bridge/interfaces/IInbox.sol";
import {IBridge} from "../dev/vendor/arb-bridge-eth/v0.8.0-custom/contracts/bridge/interfaces/IBridge.sol";

contract MockArbitrumInbox is IInbox {
  event RetryableTicketNoRefundAliasRewriteCreated(
    address destAddr,
    uint256 arbTxCallValue,
    uint256 maxSubmissionCost,
    address submissionRefundAddress,
    address valueRefundAddress,
    uint256 maxGas,
    uint256 gasPriceBid,
    bytes data
  );

  function sendL2Message(bytes calldata messageData) external override returns (uint256) {
    return 0;
  }

  function sendUnsignedTransaction(
    uint256 maxGas,
    uint256 gasPriceBid,
    uint256 nonce,
    address destAddr,
    uint256 amount,
    bytes calldata data
  ) external override returns (uint256) {
    return 0;
  }

  function sendContractTransaction(
    uint256 maxGas,
    uint256 gasPriceBid,
    address destAddr,
    uint256 amount,
    bytes calldata data
  ) external override returns (uint256) {
    return 0;
  }

  function sendL1FundedUnsignedTransaction(
    uint256 maxGas,
    uint256 gasPriceBid,
    uint256 nonce,
    address destAddr,
    bytes calldata data
  ) external payable override returns (uint256) {
    return 0;
  }

  function sendL1FundedContractTransaction(
    uint256 maxGas,
    uint256 gasPriceBid,
    address destAddr,
    bytes calldata data
  ) external payable override returns (uint256) {
    return 0;
  }

  function createRetryableTicketNoRefundAliasRewrite(
    address destAddr,
    uint256 arbTxCallValue,
    uint256 maxSubmissionCost,
    address submissionRefundAddress,
    address valueRefundAddress,
    uint256 maxGas,
    uint256 gasPriceBid,
    bytes calldata data
  ) external payable override returns (uint256) {
    emit RetryableTicketNoRefundAliasRewriteCreated(
      destAddr,
      arbTxCallValue,
      maxSubmissionCost,
      submissionRefundAddress,
      valueRefundAddress,
      maxGas,
      gasPriceBid,
      data
    );
    return 42;
  }

  function createRetryableTicket(
    address destAddr,
    uint256 arbTxCallValue,
    uint256 maxSubmissionCost,
    address submissionRefundAddress,
    address valueRefundAddress,
    uint256 maxGas,
    uint256 gasPriceBid,
    bytes calldata data
  ) external payable override returns (uint256) {
    return 0;
  }

  function depositEth(address destAddr) external payable override returns (uint256) {
    return 0;
  }

  function depositEthRetryable(
    address destAddr,
    uint256 maxSubmissionCost,
    uint256 maxGas,
    uint256 maxGasPrice
  ) external payable override returns (uint256) {
    return 0;
  }

  function bridge() external view override returns (IBridge) {
    return IBridge(address(0));
  }
}
