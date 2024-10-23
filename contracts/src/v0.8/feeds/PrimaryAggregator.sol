// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {AggregatorV2V3Interface} from "../shared/interfaces/AggregatorV2V3Interface.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";
import {OCR2Abstract} from "../shared/ocr2/OCR2Abstract.sol";

contract PrimaryAggregator is OCR2Abstract, OwnerIsCreator, AggregatorV2V3Interface {
  struct Transmitter {
    bool active;
    uint8 index;
    uint96 paymentJuels;
  }
  mapping (address => Transmitter) internal s_transmitters;

  struct Signer {
    bool active;
    uint8 index;
  }
  mapping (address => Signer) internal s_signers;

  address[] internal s_signersList;
  address[] internal s_transmittersList;
  uint32[MAX_NUM_ORACLES] internal s_rewardFromAggregatorRoundId;
  bytes32 s_latestConfigDigest;

  struct HotVars {
    uint8 f;
    uint40 latestEpochAndRound;
    uint32 latestAggregatorRoundId;
    uint32 maximumGasPriceGwei;
    uint32 reasonableGasPriceGwei;
    uint32 observationPaymentGjuels;
    uint32 transmissionPaymentGjuels;
    uint24 accountingGas;
  }
  HotVars internal s_hotVars;

  struct Transmission {
    int192 answer;
    uint32 observationsTimestamp;
    uint32 transmissionTimestamp;
  }
  mapping(uint32 => Transmission) internal s_transmissions;

  int192 immutable public minAnswer;
  int192 immutable public maxAnswer;

  uint32 internal s_configCount;
  uint32 internal s_latestConfigBlockNumber;

  struct SetConfigArgs {
    address[] signers;
    address[] transmitters;
    uint8 f;
    bytes onchainConfig;
    uint64 offchainConfigVersion;
    bytes offchainConfig;
  }

  AggregatorV2V3Interface internal fallbackFeed;

  uint8 immutable public override decimals;

  constructor (uint8 decimals_) {
    decimals = decimals_;
  }

  function typeAndVersion() external pure returns (string memory) {
    return "PrimaryAggregator v0";
  }

  function description() external view returns (string memory) {
    return "PrimaryAggregator";
  }

  function version() external view returns (uint256) {
    return 0;
  }

  function _requirePositiveF (
    uint256 f
  )
    internal
    pure
    virtual
  {
    require(0 < f, "f must be positive");
  }

  function _payOracles()
    internal
  {
    unchecked {
      LinkTokenInterface linkToken = s_linkToken;
      uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;
      uint32[maxNumOracles] memory rewardFromAggregatorRoundId = s_rewardFromAggregatorRoundId;
      address[] memory transmitters = s_transmittersList;
      for (uint transmitteridx = 0; transmitteridx < transmitters.length; transmitteridx++) {
        uint256 reimbursementAmountJuels = s_transmitters[transmitters[transmitteridx]].paymentJuels;
        s_transmitters[transmitters[transmitteridx]].paymentJuels = 0;
        uint256 obsCount = latestAggregatorRoundId - rewardFromAggregatorRoundId[transmitteridx];
        uint256 juelsAmount =
          obsCount * uint256(s_hotVars.observationPaymentGjuels) * (1 gwei) + reimbursementAmountJuels;
        if (juelsAmount > 0) {
            address payee = s_payees[transmitters[transmitteridx]];
            // Poses no re-entrancy issues, because LINK.transfer does not yield
            // control flow.
            require(linkToken.transfer(payee, juelsAmount), "insufficient funds");
            rewardFromAggregatorRoundId[transmitteridx] = latestAggregatorRoundId;
            emit OraclePaid(transmitters[transmitteridx], payee, juelsAmount, linkToken);
          }
      }
      // "Zero" the accounting storage variables
      s_rewardFromAggregatorRoundId = rewardFromAggregatorRoundId;
    }
  }

  /// @inheritdoc OCR2Abstract
  function setConfig(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  )
    external
    override
    onlyOwner()
  {
    require(signers.length <= MAX_NUM_ORACLES, "too many oracles");
    require(signers.length == transmitters.length, "oracle length mismatch");
    require(3*f < signers.length, "faulty-oracle f too high");
    _requirePositiveF(f);
    require(keccak256(onchainConfig) == keccak256(abi.encodePacked(uint8(1) /*version*/, minAnswer, maxAnswer)), "invalid onchainConfig");

    SetConfigArgs memory args = SetConfigArgs({
      signers: signers,
      transmitters: transmitters,
      f: f,
      onchainConfig: onchainConfig,
      offchainConfigVersion: offchainConfigVersion,
      offchainConfig: offchainConfig
    });

    s_hotVars.latestEpochAndRound = 0;
    _payOracles();

    // remove any old signer/transmitter addresses
    uint256 oldLength = s_signersList.length;
    for (uint256 i = 0; i < oldLength; i++) {
      address signer = s_signersList[i];
      address transmitter = s_transmittersList[i];
      delete s_signers[signer];
      delete s_transmitters[transmitter];
    }
    delete s_signersList;
    delete s_transmittersList;

    // add new signer/transmitter addresses
    for (uint i = 0; i < args.signers.length; i++) {
      require(
        !s_signers[args.signers[i]].active,
        "repeated signer address"
      );
      s_signers[args.signers[i]] = Signer({
        active: true,
        index: uint8(i)
      });
      require(
        !s_transmitters[args.transmitters[i]].active,
        "repeated transmitter address"
      );
      s_transmitters[args.transmitters[i]] = Transmitter({
        active: true,
        index: uint8(i),
        paymentJuels: 0
      });
    }
    s_signersList = args.signers;
    s_transmittersList = args.transmitters;

    s_hotVars.f = args.f;
    uint32 previousConfigBlockNumber = s_latestConfigBlockNumber;
    s_latestConfigBlockNumber = uint32(block.number);
    s_configCount += 1;
    s_latestConfigDigest = _configDigestFromConfigData(
      block.chainid,
      address(this),
      s_configCount,
      args.signers,
      args.transmitters,
      args.f,
      args.onchainConfig,
      args.offchainConfigVersion,
      args.offchainConfig
    );

    emit ConfigSet(
      previousConfigBlockNumber,
      s_latestConfigDigest,
      s_configCount,
      args.signers,
      args.transmitters,
      args.f,
      args.onchainConfig,
      args.offchainConfigVersion,
      args.offchainConfig
    );

    uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;
    for (uint256 i = 0; i < args.signers.length; i++) {
      s_rewardFromAggregatorRoundId[i] = latestAggregatorRoundId;
    }
  }

  /// @inheritdoc OCR2Abstract
  function latestConfigDetails()
    external
    override
    view
    returns (
      uint32 configCount,
      uint32 blockNumber,
      bytes32 configDigest
    )
  {
    return (s_configCount, s_latestConfigBlockNumber, s_latestConfigDigest);
  }

  /**
   * @return list of addresses permitted to transmit reports to this contract

   * @dev The list will match the order used to specify the transmitter during setConfig
   */
  function getTransmitters()
    external
    view
    returns(address[] memory)
  {
    return s_transmittersList;
  }

  function getRoundData(
    uint80 _roundId
  ) external view returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) {
    return (0, 0, 0, 0 , 0);
  }

  function latestRoundData()
    external
    view
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) {
      return (0, 0, 0, 0 , 0);
    }

  function latestAnswer() external view returns (int256) { return 0; }

  function latestTimestamp() external view returns (uint256) { return 0; }

  function latestRound() external view returns (uint256) { return 0; }

  function getAnswer(uint256 roundId) external view returns (int256) { return 0; }

  function getTimestamp(uint256 roundId) external view returns (uint256) { return 0; }

  function latestConfigDetails()
    override
    external
    view
    returns (uint32 configCount, uint32 blockNumber, bytes32 configDigest) {
      return (0, 0, 0);
    }

  /// @notice optionally returns the latest configDigest and epoch for which a
  /// report was successfully transmitted. Alternatively, the contract may return
  /// scanLogs set to true and use Transmitted events to provide this information
  /// to offchain watchers.
  /// @return scanLogs indicates whether to rely on the configDigest and epoch
  /// returned or whether to scan logs for the Transmitted event instead.
  /// @return configDigest
  /// @return epoch
  function latestConfigDigestAndEpoch()
    override
    external
    view
    returns (bool scanLogs, bytes32 configDigest, uint32 epoch) {
      return (false, 0, 0);
    }

  /// @notice transmit is called to post a new report to the contract
  /// @param reportContext [0]: ConfigDigest, [1]: 27 byte padding, 4-byte epoch and 1-byte round, [2]: ExtraHash
  /// @param report serialized report, which the signatures are signing.
  /// @param rs ith element is the R components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries
  /// @param ss ith element is the S components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries
  /// @param rawVs ith element is the the V component of the ith signature
  function transmit(
    // NOTE: If these parameters are changed, expectedMsgDataLength and/or
    // TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT need to be changed accordingly
    bytes32[3] calldata reportContext,
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs // signatures
  ) override external {}

}


