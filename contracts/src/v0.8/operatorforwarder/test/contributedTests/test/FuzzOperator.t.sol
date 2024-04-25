pragma solidity ^0.8.19;

import "../helpers/Deployer.sol";
import "../../../AuthorizedReceiver.sol";
import "../../../Operator.sol";

contract FuzzTokenReceiver is Test {

    bytes4 private constant ORACLE_REQUEST_SELECTOR = Operator.oracleRequest.selector;
    bytes4 private constant OPERATOR_REQUEST_SELECTOR = Operator.operatorRequest.selector;
    uint256 private constant SELECTOR_LENGTH = 4;
    uint256 private constant EXPECTED_REQUEST_WORDS = 2;
    uint256 private constant MINIMUM_REQUEST_LENGTH = SELECTOR_LENGTH + (32 * EXPECTED_REQUEST_WORDS);
    bytes receivedData;

    address _spoofedLinkToken = address(0xCCBBAA);
    
    function testFuzzOnTokenTransferFromLink(
        address _actualSender,
        uint256 _actualAmount,
        address _sender,
        uint256 _amount,
        bytes32 _specId,
        address _callbackAddress,
        bytes4 _callbackFunc,
        uint256 _nonce,
        uint256 _dataVersion,
        bytes calldata _dataBytes
        ) public {
            

        bytes memory tokenTransferData = abi.encodeWithSignature(
            "oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)",
            _sender,
            _amount,
            _specId,
            _callbackAddress,
            _callbackFunc,
            _nonce,
            _dataVersion,
            _dataBytes
        );

        bytes memory expectedData = abi.encodeWithSignature(
            "oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)",
            _actualSender,
            _actualAmount,
            _specId,
            _callbackAddress,
            _callbackFunc,
            _nonce,
            _dataVersion,
            _dataBytes
        );
        
        _onTokenTransfer(_actualSender, _actualAmount, tokenTransferData);

        assertEq(receivedData, expectedData, "data does not match expected");

    }

    function _onTokenTransfer(
    address sender,
    uint256 amount,
    bytes memory data
    ) private permittedFunctionsForLINK(data) {
        assembly {
            // solhint-disable-next-line avoid-low-level-calls
            mstore(add(data, 36), sender) // ensure correct sender is passed
            // solhint-disable-next-line avoid-low-level-calls
            mstore(add(data, 68), amount) // ensure correct amount is passed0.8.19
        }
        // solhint-disable-next-line avoid-low-level-calls
        (bool success, ) = address(this).delegatecall(data); // calls oracleRequest
        require(success, "Unable to create request");
    }

    // @dev Reverts if the given data does not begin with the `oracleRequest` function selector
    // @param data The data payload of the request
    modifier permittedFunctionsForLINK(bytes memory data) {
        bytes4 funcSelector;
        assembly {
        // solhint-disable-next-line avoid-low-level-calls
        funcSelector := mload(add(data, 32))
        }
        _validateTokenTransferAction(funcSelector, data);
        _;
    }

    // @notice Require that the token transfer action is valid
    // @dev OPERATOR_REQUEST_SELECTOR = multiword, ORACLE_REQUEST_SELECTOR = singleword
    function _validateTokenTransferAction(bytes4 funcSelector, bytes memory data) internal pure {
        require(data.length >= MINIMUM_REQUEST_LENGTH, "Invalid request length");
        require(
        funcSelector == OPERATOR_REQUEST_SELECTOR || funcSelector == ORACLE_REQUEST_SELECTOR,
        "Must use whitelisted functions"
        );
    }

    // Callback function for oracle request fulfillment
    fallback(bytes calldata _receivedData) external returns (bytes memory) {
        receivedData = _receivedData;
    }
}