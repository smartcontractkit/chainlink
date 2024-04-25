pragma solidity ^0.8.19;

import "../helpers/Deployer.sol";
import "../../../AuthorizedReceiver.sol";

contract OperatorTest is Deployer, AuthorizedReceiver {
    address operatorAddr;
    uint256 private dataReceived;
    address operatorOwner = address(0xAABBCC);



    function setUp() public {
        _setUp();
        operator = new Operator(address(link), operatorOwner);
        operatorAddr = address(operator);
        dataReceived = 0;
    }

    // Callback function for oracle request fulfillment
    function callback(bytes32 _requestId) public {
        require(msg.sender == operatorAddr, "Only Operator can call this function");
        dataReceived += 1;
    }

    // @notice concrete implementation of AuthorizedReceiver
    // @return bool of whether sender is authorized
    function _canSetAuthorizedSenders() internal view override returns (bool) {
        return true;
    }


    function testOracleRequestFlow() public {
        // Define some mock values
        bytes32 specId = keccak256("testSpec");
        bytes4 callbackFunctionId = bytes4(keccak256("callback(bytes32)"));
        uint256 nonce = 0;
        uint256 dataVersion = 1;
        bytes memory data = "";

        uint256 initialLinkBalance = link.balanceOf(operatorAddr);
        uint256 payment = 1 ether; // Mock payment value

        uint256 withdrawableBefore = operator.withdrawable();

        // Send LINK tokens to the Operator contract using `transferAndCall`
        link.setBalance(address(alice), 1 ether);
        assertEq(link.balanceOf(address(alice)), 1 ether, "balance update failed");

        vm.prank(alice);
        link.transferAndCall(operatorAddr, payment, abi.encodeWithSignature("oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)", address(this), payment, specId, address(this), callbackFunctionId, nonce, dataVersion, data));

        // Check that the LINK tokens were transferred to the Operator contract
        assertEq(link.balanceOf(operatorAddr), initialLinkBalance + payment);
        // No withdrawable LINK as it's all locked
        assertEq(operator.withdrawable(), withdrawableBefore);
    }

    function testFulfillOracleRequest() public {
        // This test file is the callback target and actual sender contract
        // so we should enable it to set Authorised senders to itself
        address[] memory senders = new address[](2);
        senders[0] = address(this);
        senders[0] = address(bob);

        vm.prank(address(operatorOwner));
        operator.setAuthorizedSenders(senders);

        uint256 withdrawableBefore = operator.withdrawable();

        // Define mock values for creating a new oracle request
        bytes32 specId = keccak256("testSpecForFulfill");
        bytes4 callbackFunctionId = bytes4(keccak256("callback(bytes32)"));
        uint256 nonce = 1; 
        uint256 dataVersion = 1;
        bytes memory dataBytes = ""; 
        uint256 payment = 1 ether;
        uint256 expiration = block.timestamp + 5 minutes;

        // Convert bytes to bytes32
        bytes32 data = bytes32(uint256(keccak256(dataBytes)));

        // Send LINK tokens to the Operator contract using `transferAndCall`
        link.setBalance(address(bob), 1 ether);
        vm.prank(bob);
        link.transferAndCall(operatorAddr, payment, abi.encodeWithSignature("oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)", address(this), payment, specId, address(this), callbackFunctionId, nonce, dataVersion, dataBytes));


        // Fulfill the request using the operator
        bytes32 requestId = keccak256(abi.encodePacked(bob, nonce));
        vm.prank(bob);
        operator.fulfillOracleRequest(requestId, payment, address(this), callbackFunctionId, expiration, data);

        assertEq(dataReceived, 1, "Oracle request was not fulfilled");

        // Withdrawable balance
        assertEq(operator.withdrawable(), withdrawableBefore + payment, "Internal accounting not updated correctly");
    }

    function testCancelOracleRequest() public {
        // Define mock values for creating a new oracle request
        bytes32 specId = keccak256("testSpecForCancel");
        bytes4 callbackFunctionId = bytes4(keccak256("callback(bytes32)"));
        uint256 nonce = 2;
        uint256 dataVersion = 1;
        bytes memory dataBytes = "";
        uint256 payment = 1 ether;
        uint256 expiration = block.timestamp + 5 minutes; 
        
        uint256 withdrawableBefore = operator.withdrawable();

        // Convert bytes to bytes32
        bytes32 data = bytes32(uint256(keccak256(dataBytes)));

        // Send LINK tokens to the Operator contract using `transferAndCall`
        link.setBalance(address(bob), 1 ether);
        vm.prank(bob);
        link.transferAndCall(operatorAddr, payment, abi.encodeWithSignature("oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)", address(bob), payment, specId, address(bob), callbackFunctionId, nonce, dataVersion, dataBytes));

        // No withdrawable balance as it's all locked
        assertEq(operator.withdrawable(), withdrawableBefore, "Internal accounting not updated correctly");

        bytes32 requestId = keccak256(abi.encodePacked(bob, nonce));

        vm.startPrank(alice);
        vm.expectRevert(bytes("Params do not match request ID"));
        operator.cancelOracleRequest(requestId, payment, callbackFunctionId, expiration);
        vm.stopPrank();
        
        vm.startPrank(bob);
        vm.expectRevert(bytes("Request is not expired"));
        operator.cancelOracleRequest(requestId, payment, callbackFunctionId, expiration);

        vm.warp(expiration);
        operator.cancelOracleRequest(requestId, payment, callbackFunctionId, expiration);
        vm.stopPrank();

        // Check if the LINK tokens were refunded to the sender (bob in this case)
        assertEq(link.balanceOf(address(bob)), 1 ether, "Oracle request was not canceled properly");

        assertEq(operator.withdrawable(), withdrawableBefore, "Internal accounting not updated correctly");
    }

    function testUnauthorizedFulfillment() public {
        bytes32 specId = keccak256("unauthorizedFulfillSpec");
        bytes4 callbackFunctionId = bytes4(keccak256("callback(bytes32)"));
        uint256 nonce = 5;
        uint256 dataVersion = 1;
        bytes memory dataBytes = "";
        uint256 payment = 1 ether;
        uint256 expiration = block.timestamp + 5 minutes;

        link.setBalance(address(alice), payment);
        vm.prank(alice);
        link.transferAndCall(operatorAddr, payment, abi.encodeWithSignature("oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)", address(alice), payment, specId, address(this), callbackFunctionId, nonce, dataVersion, dataBytes));

        bytes32 requestId = keccak256(abi.encodePacked(alice, nonce));

        vm.prank(address(bob));
        vm.expectRevert(bytes("Not authorized sender"));
        operator.fulfillOracleRequest(requestId, payment, address(this), callbackFunctionId, expiration, bytes32(uint256(keccak256(dataBytes))));
    }
}