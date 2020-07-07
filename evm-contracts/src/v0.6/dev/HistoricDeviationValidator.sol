pragma solidity 0.6.6;

import './AnswerValidatorInterface.sol';
import '../Flags.sol';
import '../Owned.sol';
import '../vendor/SafeMath.sol';

contract HistoricDeviationValidator is Owned, AnswerValidatorInterface {
  using SafeMath for uint256;

  uint32 constant public THRESHOLD_MULTIPLIER = 100000;

  uint32 public flaggingThreshold;
  Flags public flags;

  constructor(
    address newFlags,
    uint24 newFT
  )
    public
  {
    flags = Flags(newFlags);
    flaggingThreshold = newFT;
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


  // PRIVATE

  function percentDiff(
    int256 previous,
    int256 current
  )
    private
    returns (uint256)
  {
    uint256 difference = absolute(current - previous);
    return difference.mul(THRESHOLD_MULTIPLIER).div(absolute(previous));
  }

  function arrayifyMsgSender()
    private
    returns (address[] memory)
  {
      address[] memory addresses = new address[](1);
      addresses[0] = msg.sender;
      return addresses;
  }

  function absolute(
    int256 value
  )
    private
    returns (uint256)
  {
    if (value < 0) {
      return uint256(value * -1);
    }
    return uint256(value);
  }

}

