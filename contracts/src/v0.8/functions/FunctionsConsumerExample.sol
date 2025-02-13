// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import  { ConfirmedOwner } from "../shared/access/ConfirmedOwner.sol";
import {FunctionsClient} from "./v1_0_0/FunctionsClient.sol";
import {CBOR} from "../vendor/solidity-cborutils/v2.0.0/CBOR.sol";
//import {FunctionsRequest} from "v1_0_0/libraries/FunctionsRequest.sol";

/**
 * THIS IS AN EXAMPLE CONTRACT THAT USES HARDCODED VALUES FOR CLARITY.
 * THIS IS AN EXAMPLE CONTRACT THAT USES UN-AUDITED CODE.
 * DO NOT USE THIS CODE IN PRODUCTION.
 */
contract FunctionsConsumerExample is FunctionsClient, ConfirmedOwner {
    using CBOR for CBOR.CBORBuffer;

    uint256 internal constant DEFAULT_BUFFER_SIZE = 256;

    error EmptySource();
    error EmptySecrets();
    error EmptyArgs();
    error NoInlineSecrets();

    enum Location {
        Inline, // Provided within the Request
        Remote, // Hosted through remote location that can be accessed through a provided URL
        DONHosted // Hosted on the DON's storage
    }

    enum CodeLanguage {
        JavaScript
        // In future version we may add other languages
    }

    struct Request {
        Location codeLocation; // ════════════╸ The location of the source code that will be executed on each node in the DON
        Location secretsLocation; // ═════════╸ The location of secrets that will be passed into the source code. *Only Remote secrets are supported
        CodeLanguage language; // ════════════╸ The coding language that the source code is written in
        string source; // ════════════════════╸ Raw source code for Request.codeLocation of Location.Inline, URL for Request.codeLocation of Location.Remote, or slot decimal number for Request.codeLocation of Location.DONHosted
        bytes encryptedSecretsReference; // ══╸ Encrypted URLs for Request.secretsLocation of Location.Remote (use addSecretsReference()), or CBOR encoded slotid+version for Request.secretsLocation of Location.DONHosted (use addDONHostedSecrets())
        string[] args; // ════════════════════╸ String arguments that will be passed into the source code
        bytes[] bytesArgs; // ════════════════╸ Bytes arguments that will be passed into the source code
    }

    bytes32 public s_lastRequestId;
    bytes public s_lastResponse;
    bytes public s_lastError;

    error UnexpectedRequestID(bytes32 requestId);

    event Response(bytes32 indexed requestId, bytes response, bytes err);

//    address public router = 0x234a5fb5Bd614a7AA2FfAB244D603abFA0Ac5C5C; // arb sepolia router
    address public router = 0xb83E47C2bC239B3bf370bc41e1459A34b41238D0; // eth sepolia router

    constructor() FunctionsClient(router) ConfirmedOwner(msg.sender) { }

    /**
     * @notice Send a simple request
     * @param source JavaScript source code
     * @param encryptedSecretsUrls Encrypted URLs where to fetch user secrets
     * @param donHostedSecretsSlotID Don hosted secrets slotId
     * @param donHostedSecretsVersion Don hosted secrets version
     * @param args List of arguments accessible from within the source code
     * @param bytesArgs Array of bytes arguments, represented as hex strings
     * @param subscriptionId Billing ID
     */
    function sendRequest(
        string memory source,
        bytes memory encryptedSecretsUrls,
        uint8 donHostedSecretsSlotID,
        uint64 donHostedSecretsVersion,
        string[] memory args,
        bytes[] memory bytesArgs,
        uint64 subscriptionId,
        uint32 gasLimit,
        bytes32 donID
    ) external onlyOwner returns (bytes32 requestId) {
        Request memory req;
        initializeRequestForInlineJavaScript(req, source);
        if (encryptedSecretsUrls.length > 0)
            addSecretsReference(req, encryptedSecretsUrls);
        else if (donHostedSecretsVersion > 0) {
            addDONHostedSecrets(
                req,
                donHostedSecretsSlotID,
                donHostedSecretsVersion
            );
        }
        if (args.length > 0) setArgs(req, args);
        if (bytesArgs.length > 0) setBytesArgs(req, bytesArgs);
        s_lastRequestId = _sendRequest(
            encodeCBOR(req),
            subscriptionId,
            gasLimit,
            donID
        );
        return s_lastRequestId;
    }

    /// @notice Sets args for the user run function
    /// @param self The initialized request
    /// @param args The array of string args (must not be empty)
    function setArgs(Request memory self, string[] memory args) internal pure {
        if (args.length == 0) revert EmptyArgs();

        self.args = args;
    }

    /// @notice Sets bytes args for the user run function
    /// @param self The initialized request
    /// @param args The array of bytes args (must not be empty)
    function setBytesArgs(Request memory self, bytes[] memory args) internal pure {
        if (args.length == 0) revert EmptyArgs();

        self.bytesArgs = args;
    }

    /// @notice Adds DON-hosted secrets reference to a Request
    /// @param self The initialized request
    /// @param slotID Slot ID of the user's secrets hosted on DON
    /// @param version User data version (for the slotID)
    function addDONHostedSecrets(Request memory self, uint8 slotID, uint64 version) internal pure {
        CBOR.CBORBuffer memory buffer = CBOR.create(DEFAULT_BUFFER_SIZE);

        buffer.writeString("slotID");
        buffer.writeUInt64(slotID);
        buffer.writeString("version");
        buffer.writeUInt64(version);

        self.secretsLocation = Location.DONHosted;
        self.encryptedSecretsReference = buffer.buf.buf;
    }

    /// @notice Encodes a Request to CBOR encoded bytes
    /// @param self The request to encode
    /// @return CBOR encoded bytes
    function encodeCBOR(Request memory self) internal pure returns (bytes memory) {
        CBOR.CBORBuffer memory buffer = CBOR.create(DEFAULT_BUFFER_SIZE);

        buffer.writeString("codeLocation");
        buffer.writeUInt256(uint256(self.codeLocation));

        buffer.writeString("language");
        buffer.writeUInt256(uint256(self.language));

        buffer.writeString("source");
        buffer.writeString(self.source);

        if (self.args.length > 0) {
            buffer.writeString("args");
            buffer.startArray();
            for (uint256 i = 0; i < self.args.length; ++i) {
                buffer.writeString(self.args[i]);
            }
            buffer.endSequence();
        }

        if (self.encryptedSecretsReference.length > 0) {
            if (self.secretsLocation == Location.Inline) {
                revert NoInlineSecrets();
            }
            buffer.writeString("secretsLocation");
            buffer.writeUInt256(uint256(self.secretsLocation));
            buffer.writeString("secrets");
            buffer.writeBytes(self.encryptedSecretsReference);
        }

        if (self.bytesArgs.length > 0) {
            buffer.writeString("bytesArgs");
            buffer.startArray();
            for (uint256 i = 0; i < self.bytesArgs.length; ++i) {
                buffer.writeBytes(self.bytesArgs[i]);
            }
            buffer.endSequence();
        }

        return buffer.buf.buf;
    }

    /// @notice Initializes a Chainlink Functions Request
    /// @dev Simplified version of initializeRequest for PoC
    /// @param self The uninitialized request
    /// @param javaScriptSource The user provided JS code (must not be empty)
    function initializeRequestForInlineJavaScript(Request memory self, string memory javaScriptSource) internal pure {
        initializeRequest(self, Location.Inline, CodeLanguage.JavaScript, javaScriptSource);
    }

    function initializeRequest(
        Request memory self,
        Location codeLocation,
        CodeLanguage language,
        string memory source
    ) internal pure {
        if (bytes(source).length == 0) revert EmptySource();

        self.codeLocation = codeLocation;
        self.language = language;
        self.source = source;
    }

    /// @notice Adds Remote user encrypted secrets to a Request
    /// @param self The initialized request
    /// @param encryptedSecretsReference Encrypted comma-separated string of URLs pointing to off-chain secrets
    function addSecretsReference(Request memory self, bytes memory encryptedSecretsReference) internal pure {
        if (encryptedSecretsReference.length == 0) revert EmptySecrets();

        self.secretsLocation = Location.Remote;
        self.encryptedSecretsReference = encryptedSecretsReference;
    }

    /**
     * @notice Send a pre-encoded CBOR request
     * @param request CBOR-encoded request data
     * @param subscriptionId Billing ID
     * @param gasLimit The maximum amount of gas the request can consume
     * @param donID ID of the job to be invoked
     * @return requestId The ID of the sent request
     */
    function sendRequestCBOR(
        bytes memory request,
        uint64 subscriptionId,
        uint32 gasLimit,
        bytes32 donID
    ) external onlyOwner returns (bytes32 requestId) {
        s_lastRequestId = _sendRequest(
            request,
            subscriptionId,
            gasLimit,
            donID
        );
        return s_lastRequestId;
    }

    /**
     * @notice Store latest result/error
     * @param requestId The request ID, returned by sendRequest()
     * @param response Aggregated response from the user code
     * @param err Aggregated error from the user code or from the execution pipeline
     * Either response or error parameter will be set, but never both
     */
    function fulfillRequest(
        bytes32 requestId,
        bytes memory response,
        bytes memory err
    ) internal override {
        if (s_lastRequestId != requestId) {
            revert UnexpectedRequestID(requestId);
        }
        s_lastResponse = response;
        s_lastError = err;
        emit Response(requestId, s_lastResponse, s_lastError);
    }
}