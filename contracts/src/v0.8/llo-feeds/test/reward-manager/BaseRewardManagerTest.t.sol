// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {ERC20Mock} from "../../../shared/vendor/ERC20Mock.sol";
import "../../RewardManager.sol";
import {Common} from "../../../libraries/internal/Common.sol";
import "forge-std/console.sol";

/**
  * @title BaseRewardManagerTest
  * @author Michael Fletcher
  * @notice Base class for all reward manager tests
  * @dev This contract is intended to be inherited from and not used directly. It contains functionality to setup a primary and secondary
  */
contract BaseRewardManagerTest is Test {

    //contracts
    ERC20Mock internal asset;
    ERC20Mock internal unsupported;
    RewardManager internal rewardManager;

    //contract owner
    address internal constant ADMIN = address(1);
    //user to represent verifier contract
    address internal constant USER = address(2);

    //default recipients configured in reward manager
    address internal constant DEFAULT_RECIPIENT_1 = address(3);
    address internal constant DEFAULT_RECIPIENT_2 = address(4);
    address internal constant DEFAULT_RECIPIENT_3 = address(5);
    address internal constant DEFAULT_RECIPIENT_4 = address(6);
    address internal constant DEFAULT_RECIPIENT_5 = address(7);
    address internal constant DEFAULT_RECIPIENT_6 = address(8);
    address internal constant DEFAULT_RECIPIENT_7 = address(9);

    //default address for unregistered recipient
    address internal constant INVALID_RECIPIENT = address(0);

    //two pools should be enough to test all edge cases
    bytes32 internal constant PRIMARY_POOL_ID = keccak256("primary_pool");
    bytes32 internal constant SECONDARY_POOL_ID = keccak256("secondary_pool");

    //convenience arrays of all pool combinations used for testing
    bytes32[] internal PRIMARY_POOL_ARRAY = [PRIMARY_POOL_ID];
    bytes32[] internal SECONDARY_POOL_ARRAY = [SECONDARY_POOL_ID];
    bytes32[] internal ALL_POOLS = [PRIMARY_POOL_ID, SECONDARY_POOL_ID];

    //erc20 config
    uint256 internal constant DEFAULT_MINT_QUANTITY = 100 ether;

    //reward scalar (this should match the const in the contract)
    uint256 internal constant POOL_SCALAR = 10000;

    //the selector for the Unauthorized error
    bytes4 internal constant UNAUTHORIZED_ERROR_SELECTOR = bytes4(keccak256("Unauthorized()"));

    function setUp() public virtual {
        //change to admin user
        vm.startPrank(ADMIN);

        //init required contracts
        _initializeERC20Contracts();
        _initializeRewardManager();
    }

    function _initializeERC20Contracts() internal {
        //create the contracts
        asset = new ERC20Mock();
        unsupported = new ERC20Mock();

        //mint some tokens to the admin
        asset.mint(ADMIN, DEFAULT_MINT_QUANTITY);
        unsupported.mint(ADMIN, DEFAULT_MINT_QUANTITY);

        //mint some tokens to the user
        asset.mint(USER, DEFAULT_MINT_QUANTITY);
        unsupported.mint(USER, DEFAULT_MINT_QUANTITY);
    }

    function _initializeRewardManager() internal {
        //create the contract
        rewardManager = new RewardManager(address(asset));
    }

    function createPrimaryPool() public {
        rewardManager.setRewardRecipients(PRIMARY_POOL_ID, getPrimaryRecipients());
    }

    function createSecondaryPool() public {
        rewardManager.setRewardRecipients(SECONDARY_POOL_ID, getSecondaryRecipients());
    }

    //override this to test variations of different recipients. changing this function will require existing tests to be updated as constants are hardcoded to be explicit
    function getPrimaryRecipients() public virtual returns (Common.AddressAndWeight[] memory) {
        //array of recipients
        Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);

        //init each recipient with even weights. 2500 = 25% of pool
        recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 2500);
        recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 2500);
        recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, 2500);
        recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, 2500);

        return recipients;
    }

    function getPrimaryRecipientAddresses() public returns (address[] memory) {
        //array of recipients
        address[] memory recipients = new address[](4);

        //init each recipient with even weights. 2500 = 25% of pool
        recipients[0] = DEFAULT_RECIPIENT_1;
        recipients[1] = DEFAULT_RECIPIENT_2;
        recipients[2] = DEFAULT_RECIPIENT_3;
        recipients[3] = DEFAULT_RECIPIENT_4;

        return recipients;
    }

    //override this to test variations of different recipients. changing this function will require existing tests to be updated as constants are hardcoded to be explicit
    function getSecondaryRecipients() public virtual returns (Common.AddressAndWeight[] memory) {
        //array of recipients
        Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);

        //init each recipient with even weights
        recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 2500);
        recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, 2500);
        recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, 2500);
        recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_7, 2500);

        return recipients;
    }

    function getSecondaryRecipientAddresses() public returns (address[] memory) {
        //array of recipients
        address[] memory recipients = new address[](4);

        //init each recipient with even weights
        recipients[0] = DEFAULT_RECIPIENT_1;
        recipients[1] = DEFAULT_RECIPIENT_5;
        recipients[2] = DEFAULT_RECIPIENT_6;
        recipients[3] = DEFAULT_RECIPIENT_7;

        return recipients;
    }

    function addFundsToPool(bytes32 poolId, address sender, Common.Asset memory amount) public {
        //record the current address and switch to the sender
        address originalAddr = msg.sender;
        changePrank(sender);

        //approve the amount we're paying into the pool
        ERC20Mock(amount.assetAddress).approve(address(rewardManager), amount.amount);

        //this represents the verifier adding some funds to the pool
        rewardManager.onFeePaid(poolId, sender, amount);

        //change back to the original address
        changePrank(originalAddr);
    }

    function getAsset(uint256 quantity) public returns (Common.Asset memory) {
        return Common.Asset(address(asset), quantity);
    }

    function getUnsupportedAsset(uint256 quantity) public returns (Common.Asset memory) {
        return Common.Asset(address(unsupported), quantity);
    }

    function getAssetBalance(address addr) public returns (uint256) {
        return asset.balanceOf(addr);
    }

    function getUnsupportedBalance(address addr) public returns (uint256) {
        return unsupported.balanceOf(addr);
    }

    function claimRewards(bytes32[] memory poolIds, address sender) public {
        //record the current address and switch to the recipient
        address originalAddr = msg.sender;
        changePrank(sender);

        //claim the rewards under this recipient address
        rewardManager.claimRewards(poolIds);

        //change back to the original address
        changePrank(originalAddr);
    }

    function payRecipients(bytes32 poolId, address[] memory recipients, address sender) public {
        //record the current address and switch to the recipient
        address originalAddr = msg.sender;
        changePrank(sender);

        //pay the recipients
        rewardManager.payRecipients(poolId, recipients);

        //change back to the original address
        changePrank(originalAddr);
    }
}