pragma solidity ^0.6.0;

import './AnswerValidatorInterface.sol';
import '../interfaces/FlagsInterface.sol';
import '../Owned.sol';
import '../vendor/SafeMath.sol';
import '../SignedSafeMath.sol';

contract HistoricDeviationValidator is Owned, AnswerValidatorInterface {
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
    address newFlags,
    uint24 newFT
  )
    public
  {
    setFlagsAddress(newFlags);
    setFlaggingThreshold(newFT);
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

  function setFlaggingThreshold(uint24 newFT)
    public
    onlyOwner()
  {
    uint24 previousFT = uint24(flaggingThreshold);

    if (previousFT != newFT) {
      flaggingThreshold = newFT;

      emit FlaggingThresholdUpdated(previousFT, newFT);
    }
  }

  function setFlagsAddress(address newFlags)
    public
    onlyOwner()
  {
    address previous = address(flags);

    if (previous != newFlags) {
      flags = FlagsInterface(newFlags);

      emit FlagsAddressUpdated(previous, newFlags);
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
    uint256 difference = abs(current - previous);
    return difference.mul(THRESHOLD_MULTIPLIER).div(abs(previous));
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

