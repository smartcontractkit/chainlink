// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {Test} from "forge-std/Test.sol";
import {DestinationVerifierProxy} from "../../DestinationVerifierProxy.sol";
import {IERC165} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {IDestinationVerifier} from "../../interfaces/IDestinationVerifier.sol";
import {IDestinationVerifierProxy} from "../../interfaces/IDestinationVerifierProxy.sol";
import {DestinationVerifier} from "../../DestinationVerifier.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";
import {DestinationFeeManager} from "../../DestinationFeeManager.sol";
import {Common} from "../../../libraries/Common.sol";
import {ERC20Mock} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";
import {WERC20Mock} from "../../../../shared/mocks/WERC20Mock.sol";
import {DestinationRewardManager} from "../../DestinationRewardManager.sol";
import {IDestinationRewardManager} from "../../interfaces/IDestinationRewardManager.sol";

contract BaseTest is Test {
    uint64 internal constant POOL_SCALAR = 1e18;
    uint64 internal constant ONE_PERCENT = POOL_SCALAR / 100;
    uint256 internal constant MAX_ORACLES = 31;
    address internal constant ADMIN = address(1);
    address internal constant USER = address(2);
    address internal constant MOCK_VERIFIER_ADDRESS = address(100);

    uint8 internal constant FAULT_TOLERANCE = 10;

    DestinationVerifierProxy internal s_verifierProxy;
    DestinationVerifier internal s_verifier;
    DestinationFeeManager internal feeManager;
    DestinationRewardManager internal rewardManager;
    ERC20Mock internal link;
    WERC20Mock internal native;

    struct Signer {
        uint256 mockPrivateKey;
        address signerAddress;
    }

    Signer[MAX_ORACLES] internal s_signers;
    bool private s_baseTestInitialized;

    // reward manager events
    event RewardRecipientsUpdated(bytes32 indexed poolId, Common.AddressAndWeight[] newRewardRecipients);

    function setUp() public virtual {
        // BaseTest.setUp is often called multiple times from tests' setUp due to inheritance.
        if (s_baseTestInitialized) return;
        s_baseTestInitialized = true;
        vm.startPrank(ADMIN);

        s_verifierProxy = new DestinationVerifierProxy();
        s_verifier = new DestinationVerifier(address(s_verifierProxy));

        // setting up FeeManager and RewardManager
        native = new WERC20Mock();
        link = new ERC20Mock("LINK", "LINK", ADMIN, 0);
        rewardManager = new DestinationRewardManager(address(link));
        feeManager =
            new DestinationFeeManager(address(link), address(native), address(s_verifier), address(rewardManager));
        s_verifier.setFeeManager(address(feeManager));
        rewardManager.setFeeManager(address(feeManager));

        for (uint256 i; i < MAX_ORACLES; i++) {
            uint256 mockPK = i + 1;
            s_signers[i].mockPrivateKey = mockPK;
            s_signers[i].signerAddress = vm.addr(mockPK);
        }
    }

    function _getSigners(uint256 numSigners) internal view returns (Signer[] memory) {
        Signer[] memory signers = new Signer[](numSigners);
        for (uint256 i; i < numSigners; i++) {
            signers[i] = s_signers[i];
        }
        return signers;
    }

    function _getSignerAddresses(Signer[] memory signers) internal view returns (address[] memory) {
        address[] memory signerAddrs = new address[](signers.length);
        for (uint256 i = 0; i < signerAddrs.length; i++) {
            signerAddrs[i] = s_signers[i].signerAddress;
        }
        return signerAddrs;
    }

    function _signerAddressAndDonConfigKey(address signer, bytes24 DONConfigID) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(signer, DONConfigID));
    }

    function _DONConfigIdFromConfigData(address[] memory signers, uint8 f) internal pure returns (bytes24) {
        Common._quickSort(signers, 0, int256(signers.length - 1));
        bytes24 DONConfigID = bytes24(keccak256(abi.encodePacked(signers, f)));
        return DONConfigID;
    }
}
