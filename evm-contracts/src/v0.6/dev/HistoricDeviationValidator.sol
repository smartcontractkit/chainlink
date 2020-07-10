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
    int256 previous,
    uint256 currentRoundId,
    int256 current
  )
    external
    override
    returns (bool)
  {
    if (previous == 0) return true;

    if (percentDiff(previous, current) > flaggingThreshold) {
      flags.raiseFlags(arrayifyMsgSender());
      return false;
    }

    return true;
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

  function percentDiff(
    int256 previous,
    int256 current
  )
    private
    returns (uint256)
  {
    int256 difference;
    if (current > previous) {
      difference = current - previous;
    } else {
      difference = previous - current;
    }
    return abs(difference).mul(THRESHOLD_MULTIPLIER).div(abs(previous));
  }

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
    returns (uint256)
  {
    return uint256(value < 0 ? value.mul(-1): value);
  }

}

