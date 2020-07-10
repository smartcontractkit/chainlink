pragma solidity ^0.6.0;

import './AggregatorValidatorInterface.sol';
import '../interfaces/FlagsInterface.sol';
import '../Owned.sol';
import '../vendor/SafeMath.sol';
import '../SignedSafeMath.sol';

contract HistoricDeviationValidator is Owned, AggregatorValidatorInterface {
  using SafeMath for uint256;
  using SignedSafeMath for int256;

  uint32 constant public THRESHOLD_MULTIPLIER = 100000;

  uint32 public flaggingThreshold;
  FlagsInterface public flags;

  event FlaggingThresholdUpdated(
    uint24 indexed previous,
    uint24 indexed current
  );
  event FlagsAddressUpdated(
    address indexed previous,
    address indexed current
  );

  constructor(
    address _flagsAddress,
    uint24 _flaggingThreshold
  )
    public
  {
    setFlagsAddress(_flagsAddress);
    setFlaggingThreshold(_flaggingThreshold);
  }

  function validate(
    uint256 previousRoundId,
    int256 previousAnswer,
    uint256 roundId,
    int256 answer
  )
    external
    override
    returns (bool)
  {
    if (!isValid(previousRoundId, previousAnswer, roundId, answer)) {
      flags.raiseFlags(arrayifyMsgSender());
      return false;
    }

    return true;
  }

  function isValid(
    uint256 ,
    int256 previousAnswer,
    uint256 ,
    int256 answer
  )
    public
    view
    returns (bool)
  {
    if (previousAnswer == 0) return true;

    int256 change = previousAnswer.sub(answer);
    uint256 percent = abs(change.mul(THRESHOLD_MULTIPLIER).div(previousAnswer));

    return percent <= flaggingThreshold;
  }

  function setFlaggingThreshold(uint24 _flaggingThreshold)
    public
    onlyOwner()
  {
    uint24 previousFT = uint24(flaggingThreshold);

    if (previousFT != _flaggingThreshold) {
      flaggingThreshold = _flaggingThreshold;

      emit FlaggingThresholdUpdated(previousFT, _flaggingThreshold);
    }
  }

  function setFlagsAddress(address _flagsAddress)
    public
    onlyOwner()
  {
    address previous = address(flags);

    if (previous != _flagsAddress) {
      flags = FlagsInterface(_flagsAddress);

      emit FlagsAddressUpdated(previous, _flagsAddress);
    }
  }


  // PRIVATE

  function arrayifyMsgSender()
    private
    returns (address[] memory)
  {
      address[] memory addresses = new address[](1);
      addresses[0] = msg.sender;
      return addresses;
  }

  function abs(
    int256 value
  )
    private
    pure
    returns (uint256)
  {
    return uint256(value < 0 ? value.mul(-1): value);
  }

}

