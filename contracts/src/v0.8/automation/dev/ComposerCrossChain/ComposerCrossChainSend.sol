pragma solidity 0.8.16;

import "../../../shared/access/ConfirmedOwner.sol";
import "../../interfaces/AutomationCompatibleInterface.sol";
import "../interfaces/ComposerCompatibleInterfaceV1.sol";
import "../../../ChainSpecificUtil.sol";
import "../../../vendor/openzeppelin-contracts/contracts/utils/Strings.sol";
import "../../../vendor/Strings.sol";
import "./CCIPDeps/TokenTransfer.sol";

contract ComposerCrossChainSend is
    ConfirmedOwner,
    TokenTransfer,
    AutomationCompatibleInterface,
    ComposerCompatibleInterfaceV1
{
    using strings for strings.slice;

    event TopUpSent(address receiver, uint256 amount, uint256 currentNonce);

    uint32 private constant MIN_GAS_FOR_PERFORM = 200_000;
    string private s_scriptHash;

    string private s_receivers;
    string private s_rpcUrl;

    mapping(address => uint256) s_nonces;

    struct Tuple {
        address addr;
        uint256 nonce;
    }

    constructor(
        string memory scriptHash,
        address _router,
        uint64 _destinationChainSelector,
        address _token,
        string memory _receivers,
        string memory _rpcUrl
    ) ConfirmedOwner(msg.sender) TokenTransfer(_router, _destinationChainSelector, _token) {
        s_scriptHash = scriptHash;
        s_receivers = _receivers;
        s_rpcUrl = _rpcUrl;
    }

    function checkUpkeep(bytes calldata /* data */ ) external view override returns (bool, bytes memory) {
        return revertForFeedLookup();
    }

    // Pass the addresses to watch and top-up, and the public RPC URL to be invoked from Functions.
    function revertForFeedLookup() public view returns (bool, bytes memory) {
        uint256 blockNumber = ChainSpecificUtil._getBlockNumber();
        string[] memory functionsArguments = new string[](2);

        functionsArguments[0] = s_receivers;
        functionsArguments[1] = s_rpcUrl;

        // Emit Composer request revert.
        revert ComposerRequestV1(s_scriptHash, functionsArguments, false, "", new string[](0), "", 0, "");
    }

    // Modified checkCallback function that matches the StreamsLookupCompatibleInterface, but
    // accepts the result of a functions call to correctly ABI-encode it. This is a stopgap
    // function only intended to exist while Functions does not yet implement more sophisticated
    // ABI-encoding.
    function checkCallback(bytes[] memory data, bytes memory lookupData)
        external
        view
        override
        returns (bool, bytes memory)
    {
        require(data.length == 1, "should only have one item for abi-decoding");
        string memory values = abi.decode(data[0], (string));

        // Parse the comma separated string Functions result.
        strings.slice memory s = strings.toSlice(values);
        strings.slice memory delim = strings.toSlice(",");
        string[] memory tuples = new string[](s.count(delim) + 1);
        for (uint256 i = 0; i < tuples.length; i++) {
            tuples[i] = s.split(delim).toString();
        }

        // Convert the strings to Tuples.
        // Format: 0x0000000000000000000000000000000000000000-0-0
        //               receiver address^        needs_update^ ^nonce
        Tuple[] memory tupleStructs = new Tuple[](tuples.length);
        uint256 counter = 0;
        for (uint256 i = 0; i < tuples.length; i++) {
            Tuple memory tuple;
            strings.slice memory tupleString = strings.toSlice(tuples[i]);
            strings.slice memory dashDelim = strings.toSlice("-");

            string memory addressString = tupleString.split(dashDelim).toString();
            tuple.addr = toAddress(addressString);

            string memory needsUpdateString = tupleString.split(dashDelim).toString();
            uint256 needsUpdate = stringToUint(needsUpdateString);

            string memory nonceString = tupleString.split(dashDelim).toString();
            tuple.nonce = stringToUint(nonceString);

            if (needsUpdate == 1) {
                tupleStructs[counter] = tuple;
                counter++;
            }
        }

        // Adjusts the length of tuples to not include empty values.
        assembly {
            mstore(tupleStructs, counter)
        }

        // Return the well-formatted performData.
        bytes memory performData = abi.encode(tupleStructs, lookupData);
        return (tupleStructs.length > 0, performData);
    }

    // Invoke CCIP sends of BnM token to receivers that need a top-up and do not already have a CCIP send in-flight.
    function performUpkeep(bytes calldata performData) external override {
        (Tuple[] memory tuples, /* bytes memory lookupData */ ) = abi.decode(performData, (Tuple[], bytes));

        for (uint256 i = 0; i < tuples.length; i++) {
            Tuple memory tuple = tuples[i];
            if (s_nonces[tuple.addr] <= tuple.nonce) {
                transferTokensPayNative(1, tuple.addr); // amount to send hardcoded to 1.
                emit TopUpSent(tuple.addr, 1, s_nonces[tuple.addr]);
                s_nonces[tuple.addr]++;
            }
            if (gasleft() < MIN_GAS_FOR_PERFORM) {
                return;
            }
        }
    }

    // HELPER FUNCTIONS
    function stringToUint(string memory s) public pure returns (uint256) {
        bytes memory b = bytes(s);
        uint256 result = 0;
        for (uint256 i = 0; i < b.length; i++) {
            uint256 c = uint256(uint8(b[i]));
            if (c >= 48 && c <= 57) {
                result = result * 10 + (c - 48);
            }
        }
        return result;
    }

    function fromHexChar(uint8 c) public pure returns (uint8) {
        if (bytes1(c) >= bytes1("0") && bytes1(c) <= bytes1("9")) {
            return c - uint8(bytes1("0"));
        }
        if (bytes1(c) >= bytes1("a") && bytes1(c) <= bytes1("f")) {
            return 10 + c - uint8(bytes1("a"));
        }
        if (bytes1(c) >= bytes1("A") && bytes1(c) <= bytes1("F")) {
            return 10 + c - uint8(bytes1("A"));
        }
        return 0;
    }

    function hexStringToAddress(string memory s) public pure returns (bytes memory) {
        bytes memory ss = bytes(s);
        require(ss.length % 2 == 0); // length must be even
        bytes memory r = new bytes(ss.length/2);
        for (uint256 i = 0; i < ss.length / 2; ++i) {
            r[i] = bytes1(fromHexChar(uint8(ss[2 * i])) * 16 + fromHexChar(uint8(ss[2 * i + 1])));
        }

        return r;
    }

    function toAddress(string memory s) public pure returns (address) {
        bytes memory _bytes = hexStringToAddress(s);
        require(_bytes.length >= 1 + 20, "toAddress_outOfBounds");
        address tempAddress;

        assembly {
            tempAddress := div(mload(add(add(_bytes, 0x20), 1)), 0x1000000000000000000000000)
        }

        return tempAddress;
    }
}
