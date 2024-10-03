// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IERC677Receiver} from "../../../interfaces/IERC677Receiver.sol";
import {IBurnMintERC20} from "../../../token/ERC20/IBurnMintERC20.sol";
import {IERC677} from "../../../token/ERC677/IERC677.sol";

import {BurnMintERC20} from "../../../token/ERC20/BurnMintERC20.sol";
import {BurnMintERC677} from "../../../token/ERC677/BurnMintERC677.sol";
import {BaseTest} from "../../BaseTest.t.sol";
import {GenericReceiver} from "../../testhelpers/GenericReceiver.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {IERC165} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol";

contract BurnMintERC677Setup is BaseTest {
  event Transfer(address indexed from, address indexed to, uint256 value);
  event MintAccessGranted(address indexed minter);
  event BurnAccessGranted(address indexed burner);
  event MintAccessRevoked(address indexed minter);
  event BurnAccessRevoked(address indexed burner);

  BurnMintERC677 internal s_burnMintERC677;
  GenericReceiver internal s_genericReceiver;

  address internal s_mockPool = address(6243783892);
  uint256 internal s_amount = 1e18;

  function setUp() public virtual override {
    BaseTest.setUp();
    s_burnMintERC677 = new BurnMintERC677("Chainlink Token", "LINK", 18, 1e27);

    // Set s_mockPool to be a burner and minter
    s_burnMintERC677.grantMintAndBurnRoles(s_mockPool);
    deal(address(s_burnMintERC677), OWNER, s_amount);

    s_genericReceiver = new GenericReceiver(false);
  }
}

contract BurnMintERC677_constructor is BurnMintERC677Setup {
  function testConstructorSuccess() public {
    string memory name = "Chainlink token v2";
    string memory symbol = "LINK2";
    uint8 decimals = 19;
    uint256 maxSupply = 1e33;
    s_burnMintERC677 = new BurnMintERC677(name, symbol, decimals, maxSupply);

    assertEq(name, s_burnMintERC677.name());
    assertEq(symbol, s_burnMintERC677.symbol());
    assertEq(decimals, s_burnMintERC677.decimals());
    assertEq(maxSupply, s_burnMintERC677.maxSupply());
  }
}

contract BurnMintERC677_supportsInterface is BurnMintERC677Setup {
  function testConstructorSuccess() public view {
    assertTrue(s_burnMintERC677.supportsInterface(type(IERC20).interfaceId));
    assertTrue(s_burnMintERC677.supportsInterface(type(IERC677).interfaceId));
    assertTrue(s_burnMintERC677.supportsInterface(type(IBurnMintERC20).interfaceId));
    assertTrue(s_burnMintERC677.supportsInterface(type(IERC165).interfaceId));
  }
}

contract BurnMintERC677_transferAndCall is BurnMintERC677Setup {
  function testTransferAndCall() public {
    vm.startPrank(OWNER);

    vm.expectCall(
      address(s_genericReceiver),
      0,
      abi.encodeWithSelector(IERC677Receiver.onTokenTransfer.selector, OWNER, s_amount, "0x")
    );

    s_burnMintERC677.transferAndCall(address(s_genericReceiver), s_amount, "0x");
  }
}
