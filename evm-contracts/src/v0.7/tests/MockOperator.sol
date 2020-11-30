pragma solidity 0.7.0;

import "../interfaces/OperatorInterface.sol";
import "../dev/OperatorProxy.sol";

contract MockOperator is OperatorInterface {
  address public immutable proxy;
  bool private s_isAuthorized;

  constructor(address link) {
    proxy = address(new OperatorProxy(link));
  }

  function setIsAuthorized(bool authorized) public {
    s_isAuthorized = authorized;
  }

  function isAuthorizedSender(address) external override view returns (bool) {
    return s_isAuthorized;
  }
  function fulfillOracleRequest(
    bytes32,
    uint256,
    address,
    bytes4,
    uint256,
    bytes32
  ) external override returns (bool) {
    return false;
  }
  function setAuthorizedSender(address, bool) external override {}
  function withdraw(address, uint256) external override {}
  function withdrawable() external override view returns (uint256) {
    return 0;
  }
}