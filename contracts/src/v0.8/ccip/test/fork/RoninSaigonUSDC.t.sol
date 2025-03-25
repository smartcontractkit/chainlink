pragma solidity ^0.8.24;

import {IOwner} from "../../interfaces/IOwner.sol";
import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {Internal} from "../../libraries/Internal.sol";
import {OffRamp} from "../../offRamp/OffRamp.sol";
import {BurnMintERC20} from "../../../shared/token/ERC20/BurnMintERC20.sol";
import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {BurnMintWithLockReleaseFlagTokenPool} from "../../pools/USDC/BurnMintWithLockReleaseFlagTokenPool.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {AccessControl} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/access/AccessControl.sol";

import {Test} from "forge-std/Test.sol";
import {console2 as console} from "forge-std/Console2.sol";
import {stdStorage, StdStorage} from "forge-std/Test.sol";              

interface RoninUSDC {
  function addMinters(address[] memory _addedMinters) external;
  function admin() external returns (address);
  function MINTER_ROLE() external returns (bytes32);
  function grantRole(bytes32 role, address account) external;
}

contract RoninSaigonUSDCTest is Test {
    using stdStorage for StdStorage;

    event ExecutionStateChanged(
    uint64 indexed sequenceNumber, bytes32 indexed messageId, Internal.MessageExecutionState state, bytes returnData
  );

    BurnMintWithLockReleaseFlagTokenPool public s_burnMintTokenPool;
    address public constant RMNProxy = 0xf206c6D3f3810eBbD75e7B4684291b5e51023D2f;

    TokenAdminRegistry public constant tokenAdminRegistry = TokenAdminRegistry(0x057879f376041D527a98327DE2Ec00F201c9cA25);

    IBurnMintERC20 public constant USDC = IBurnMintERC20(0x067FBFf8990c58Ab90BaE3c97241C5d736053F77);

    address public constant OFFRAMP = 0x77008Fbd8Ae8f395beF9c6a55905896f3Ead75e9;

    address public constant ROUTER = 0x0aCAe4e51D3DA12Dd3F45A66e8b660f740e6b820;
    address public CURRENT_USDC_ADMIN = 0xCee681C9108c42C710c6A8A949307D5F13C9F3ca;

    function setUp() public  {
    }

    function test_deployAndConfigureNewTokenPool() public {
        vm.createSelectFork("https://saigon-testnet.roninchain.com/rpc");

        s_burnMintTokenPool = new BurnMintWithLockReleaseFlagTokenPool(
            IBurnMintERC20(USDC),
            6,
            new address[](0),
            RMNProxy,
            ROUTER
        );

        bytes[] memory remotePoolAddresses = new bytes[](1);
        remotePoolAddresses[0] = abi.encode(0xAff3fE524ea94118EF09DaDBE3c77ba6AA0005EC);

        TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);

        chainUpdates[0] = TokenPool.ChainUpdate({
            remoteChainSelector: 16015286601757825753,
            remotePoolAddresses: remotePoolAddresses,
            remoteTokenAddress: abi.encode(0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238),
            outboundRateLimiterConfig: RateLimiter.Config({
                isEnabled: false,
                capacity: 0,
                rate: 0
            }),
            inboundRateLimiterConfig: RateLimiter.Config({
                isEnabled: false,
                capacity: 0,
                rate: 0
            })
        });

        s_burnMintTokenPool.applyChainUpdates(new uint64[](0), chainUpdates);
    
        // Get the admin of the current pool, should be MCMS
        address currentAdmin = tokenAdminRegistry.getTokenConfig(address(USDC)).administrator;

        // Impersonate so that the new pool can be set
        vm.startPrank(currentAdmin);

        // Set the new pool
        tokenAdminRegistry.setPool(address(USDC), address(s_burnMintTokenPool));

        // Check the new pool is set
        assertEq(tokenAdminRegistry.getTokenConfig(address(USDC)).tokenPool, address(s_burnMintTokenPool));

        // Have USDC set the new pool as a minter and burner
        vm.startPrank(CURRENT_USDC_ADMIN);
        bytes32 MINTER_ROLE = RoninUSDC(address(USDC)).MINTER_ROLE();

        // Write to storage and grant minter role since an admin address is not known
        stdstore
        .enable_packed_slots()
        .target(address(USDC))
        .sig("hasRole(bytes32,address)")
        .with_key(MINTER_ROLE)
        .with_key(address(s_burnMintTokenPool))
        .checked_write(true);

        assertTrue(AccessControl(address(USDC)).hasRole(MINTER_ROLE, address(s_burnMintTokenPool)));

bytes memory manualExecutionData = hex"8926c4ee000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000005a0000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000004a00000000000000000000000000000000000000000000000000000000000000540000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000de41ba4fc9d91ad9000000000000000000000000968d0cd7343f711216817e617d3f92a23dc91c07000000000000000000000000968d0cd7343f711216817e617d3f92a23dc91c07000000000000000000000000000000000000000000000000000000000000053b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000097d90c9d3e0b50ca60e1ae45f6a81010f9fb5340000000000000000000000000000000000000000000000000000775e6621072100000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000220a77de758956ac4ff26bac468d4c8a6da5039aef4179106aa62ec2f7d2c0f14d5000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000001c7d4b196cb0c7b01d743fbc6116a902379c723800000000000000000000000000000000000000000000000000000000000f42400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000002bf200000000000000000000000000000000000000000000000000000000000000020000000000000000000000000aff3fe524ea94118ef09dadbe3c77ba6aa0005ec0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000067fbff8990c58ab90bae3c97241c5d736053f770000000000000000000000000000000000000000000000000000000000000020fa7c07de000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000f42400000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000f4240";

        address receiver = 0x968D0Cd7343f711216817E617d3f92a23dC91c07;
        uint256 balanceBefore = USDC.balanceOf(receiver);
        uint256 expectedBalance = 1e6;
        bytes32 messageId = 0xa77de758956ac4ff26bac468d4c8a6da5039aef4179106aa62ec2f7d2c0f14d5;

        // Check the event logs
        vm.expectEmit(true, true, true, true);
        emit ExecutionStateChanged(
            1339,
            messageId,
            Internal.MessageExecutionState.SUCCESS,
            ""
        );

        // Manually execute the message
        (bool success, ) = OFFRAMP.call(manualExecutionData);
        assertTrue(success, "Manual Execution Failed when it should not have");

        uint256 balanceAfter = USDC.balanceOf(receiver);
        assertEq(balanceAfter, balanceBefore + expectedBalance, "Balance did not increase as expected");

        
    }
 
}