// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {Test} from "forge-std/Test.sol";
import {ERC20Mock} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";
import {RewardManager} from "../../RewardManager.sol";
import {Common} from "../../libraries/Common.sol";
import {IRewardManager} from "../../interfaces/IRewardManager.sol";

/**
 * @title BaseRewardManagerTest
 * @author Michael Fletcher
 * @notice Base class for all reward manager tests
 * @dev This contract is intended to be inherited from and not used directly. It contains functionality to setup a primary and secondary pool
 */
contract BaseRewardManagerTest is Test {
  //contracts
  ERC20Mock internal asset;
  ERC20Mock internal unsupported;
  RewardManager internal rewardManager;

  //default address for unregistered recipient
  address internal constant INVALID_ADDRESS = address(0);
  //contract owner
  address internal constant ADMIN = address(uint160(uint256(keccak256("ADMIN"))));
  //address to represent verifier contract
  address internal constant FEE_MANAGER = address(uint160(uint256(keccak256("FEE_MANAGER"))));
  //a general user
  address internal constant USER = address(uint160(uint256(keccak256("USER"))));

  //default recipients configured in reward manager
  address internal constant DEFAULT_RECIPIENT_1 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_1"))));
  address internal constant DEFAULT_RECIPIENT_2 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_2"))));
  address internal constant DEFAULT_RECIPIENT_3 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_3"))));
  address internal constant DEFAULT_RECIPIENT_4 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_4"))));
  address internal constant DEFAULT_RECIPIENT_5 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_5"))));
  address internal constant DEFAULT_RECIPIENT_6 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_6"))));
  address internal constant DEFAULT_RECIPIENT_7 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_7"))));

  //additional recipients not in the reward manager
  address internal constant DEFAULT_RECIPIENT_8 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_8"))));
  address internal constant DEFAULT_RECIPIENT_9 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_9"))));

  //two pools should be enough to test all edge cases
  bytes32 internal constant PRIMARY_POOL_ID = keccak256("primary_pool");
  bytes32 internal constant SECONDARY_POOL_ID = keccak256("secondary_pool");
  bytes32 internal constant INVALID_POOL_ID = keccak256("invalid_pool");
  bytes32 internal constant ZERO_POOL_ID = bytes32(0);

  //convenience arrays of all pool combinations used for testing
  bytes32[] internal PRIMARY_POOL_ARRAY = [PRIMARY_POOL_ID];
  bytes32[] internal SECONDARY_POOL_ARRAY = [SECONDARY_POOL_ID];
  bytes32[] internal ALL_POOLS = [PRIMARY_POOL_ID, SECONDARY_POOL_ID];

  //erc20 config
  uint256 internal constant DEFAULT_MINT_QUANTITY = 100 ether;

  //reward scalar (this should match the const in the contract)
  uint64 internal constant POOL_SCALAR = 1e18;
  uint64 internal constant ONE_PERCENT = POOL_SCALAR / 100;
  uint64 internal constant FIFTY_PERCENT = POOL_SCALAR / 2;
  uint64 internal constant TEN_PERCENT = POOL_SCALAR / 10;

  //the selector for each error
  bytes4 internal immutable UNAUTHORIZED_ERROR_SELECTOR = RewardManager.Unauthorized.selector;
  bytes4 internal immutable INVALID_ADDRESS_ERROR_SELECTOR = RewardManager.InvalidAddress.selector;
  bytes4 internal immutable INVALID_WEIGHT_ERROR_SELECTOR = RewardManager.InvalidWeights.selector;
  bytes4 internal immutable INVALID_POOL_ID_ERROR_SELECTOR = RewardManager.InvalidPoolId.selector;
  bytes internal constant ONLY_CALLABLE_BY_OWNER_ERROR = "Only callable by owner";
  bytes4 internal immutable INVALID_POOL_LENGTH_SELECTOR = RewardManager.InvalidPoolLength.selector;

  // Events emitted within the reward manager
  event RewardRecipientsUpdated(bytes32 indexed poolId, Common.AddressAndWeight[] newRewardRecipients);
  event RewardsClaimed(bytes32 indexed poolId, address indexed recipient, uint192 quantity);
  event FeeManagerUpdated(address newProxyAddress);
  event FeePaid(IRewardManager.FeePayment[] payments, address payee);

  function setUp() public virtual {
    //change to admin user
    vm.startPrank(ADMIN);

    //init required contracts
    _initializeERC20Contracts();
    _initializeRewardManager();
  }

  function _initializeERC20Contracts() internal {
    //create the contracts
    asset = new ERC20Mock("ASSET", "AST", ADMIN, 0);
    unsupported = new ERC20Mock("UNSUPPORTED", "UNS", ADMIN, 0);

    //mint some tokens to the admin
    asset.mint(ADMIN, DEFAULT_MINT_QUANTITY);
    unsupported.mint(ADMIN, DEFAULT_MINT_QUANTITY);

    //mint some tokens to the user
    asset.mint(FEE_MANAGER, DEFAULT_MINT_QUANTITY);
    unsupported.mint(FEE_MANAGER, DEFAULT_MINT_QUANTITY);
  }

  function _initializeRewardManager() internal {
    //create the contract
    rewardManager = new RewardManager(address(asset));

    rewardManager.setFeeManager(FEE_MANAGER);
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
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, POOL_SCALAR / 4);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, POOL_SCALAR / 4);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, POOL_SCALAR / 4);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, POOL_SCALAR / 4);

    return recipients;
  }

  function getPrimaryRecipientAddresses() public pure returns (address[] memory) {
    //array of recipients
    address[] memory recipients = new address[](4);

    recipients[0] = DEFAULT_RECIPIENT_1;
    recipients[1] = DEFAULT_RECIPIENT_2;
    recipients[2] = DEFAULT_RECIPIENT_3;
    recipients[3] = DEFAULT_RECIPIENT_4;

    return recipients;
  }

  //override this to test variations of different recipients.
  function getSecondaryRecipients() public virtual returns (Common.AddressAndWeight[] memory) {
    //array of recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);

    //init each recipient with even weights. 2500 = 25% of pool
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, POOL_SCALAR / 4);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, POOL_SCALAR / 4);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, POOL_SCALAR / 4);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_7, POOL_SCALAR / 4);

    return recipients;
  }

  function getSecondaryRecipientAddresses() public pure returns (address[] memory) {
    //array of recipients
    address[] memory recipients = new address[](4);

    recipients[0] = DEFAULT_RECIPIENT_1;
    recipients[1] = DEFAULT_RECIPIENT_5;
    recipients[2] = DEFAULT_RECIPIENT_6;
    recipients[3] = DEFAULT_RECIPIENT_7;

    return recipients;
  }

  function addFundsToPool(bytes32 poolId, Common.Asset memory amount, address sender) public {
    IRewardManager.FeePayment[] memory payments = new IRewardManager.FeePayment[](1);
    payments[0] = IRewardManager.FeePayment(poolId, uint192(amount.amount));

    addFundsToPool(payments, sender);
  }

  function addFundsToPool(IRewardManager.FeePayment[] memory payments, address sender) public {
    //record the current address and switch to the sender
    address originalAddr = msg.sender;
    changePrank(sender);

    uint256 totalPayment;
    for (uint256 i; i < payments.length; ++i) {
      totalPayment += payments[i].amount;
    }

    //approve the amount being paid into the pool
    ERC20Mock(address(asset)).approve(address(rewardManager), totalPayment);

    //this represents the verifier adding some funds to the pool
    rewardManager.onFeePaid(payments, sender);

    //change back to the original address
    changePrank(originalAddr);
  }

  function getAsset(uint256 quantity) public view returns (Common.Asset memory) {
    return Common.Asset(address(asset), quantity);
  }

  function getAssetBalance(address addr) public view returns (uint256) {
    return asset.balanceOf(addr);
  }

  function claimRewards(bytes32[] memory poolIds, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //claim the rewards
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

  function setRewardRecipients(bytes32 poolId, Common.AddressAndWeight[] memory recipients, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //pay the recipients
    rewardManager.setRewardRecipients(poolId, recipients);

    //change back to the original address
    changePrank(originalAddr);
  }

  function setFeeManager(address feeManager, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //update the proxy
    rewardManager.setFeeManager(feeManager);

    //change back to the original address
    changePrank(originalAddr);
  }

  function updateRewardRecipients(bytes32 poolId, Common.AddressAndWeight[] memory recipients, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //pay the recipients
    rewardManager.updateRewardRecipients(poolId, recipients);

    //change back to the original address
    changePrank(originalAddr);
  }
}
