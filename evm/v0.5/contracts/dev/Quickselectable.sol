pragma solidity 0.5.0;

import "../vendor/SafeMath.sol";

contract Quickselectable {
  using SafeMath for uint256;

  function quickselect(int256[] memory _a, uint256 _k)
    internal
    pure
    returns (int256)
  {
    int256[] memory a = _a;
    uint256 k = _k;
    uint256 aLen = a.length;
    int256[] memory a1 = new int256[](aLen);
    int256[] memory a2 = new int256[](aLen);
    uint256 a1Len;
    uint256 a2Len;
    int256 pivot;
    uint256 i;

    while (true) {
      pivot = a[aLen.div(2)];
      a1Len = 0;
      a2Len = 0;
      for (i = 0; i < aLen; i++) {
        if (a[i] < pivot) {
          a1[a1Len] = a[i];
          a1Len++;
        } else if (a[i] > pivot) {
          a2[a2Len] = a[i];
          a2Len++;
        }
      }
      if (k <= a1Len) {
        aLen = a1Len;
        (a, a1) = swap(a, a1);
      } else if (k > (aLen.sub(a2Len))) {
        k = k.sub(aLen.sub(a2Len));
        aLen = a2Len;
        (a, a2) = swap(a, a2);
      } else {
        return pivot;
      }
    }
  }

  function swap(int256[] memory _a, int256[] memory _b)
    private
    pure
    returns(int256[] memory, int256[] memory)
  {
    return (_b, _a);
  }

}
