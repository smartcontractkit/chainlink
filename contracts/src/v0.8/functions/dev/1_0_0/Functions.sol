// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {CBOR, Buffer} from "../../../shared/vendor/solidity-cborutils/v2.0.0/CBOR.sol";

/**
 * @title Library for Chainlink Functions
 */
library Functions {
  uint16 public constant REQUEST_DATA_VERSION = 1;
  uint256 internal constant DEFAULT_BUFFER_SIZE = 256;

  using CBOR for Buffer.buffer;

  enum Location {
    Inline,
    Remote,
    DONHosted
  }

  enum CodeLanguage {
    JavaScript
    // In future version we may add other languages
  }

  struct Request {
    Location codeLocation;
    Location secretsLocation; // Only Remote secrets are supported
    CodeLanguage language;
    string source; // Source code for Location.Inline, url for Location.Remote or slot decimal number for Location.DONHosted
    bytes encryptedSecretsReference; // Encrypted urls for Location.Remote or slot decimal number for Location.DONHosted
    string[] args;
  }

  error EmptySource();
  error EmptySecrets();
  error EmptyArgs();
  error NoInlineSecrets();

  /**
   * @notice Encodes a Request to CBOR encoded bytes
   * @param self The request to encode
   * @return CBOR encoded bytes
   */
  function encodeCBOR(Request memory self) internal pure returns (bytes memory) {
    CBOR.CBORBuffer memory buffer;
    Buffer.init(buffer.buf, DEFAULT_BUFFER_SIZE);

    CBOR.writeString(buffer, "codeLocation");
    CBOR.writeUInt256(buffer, uint256(self.codeLocation));

    CBOR.writeString(buffer, "language");
    CBOR.writeUInt256(buffer, uint256(self.language));

    CBOR.writeString(buffer, "source");
    CBOR.writeString(buffer, self.source);

    if (self.args.length > 0) {
      CBOR.writeString(buffer, "args");
      CBOR.startArray(buffer);
      for (uint256 i = 0; i < self.args.length; i++) {
        CBOR.writeString(buffer, self.args[i]);
      }
      CBOR.endSequence(buffer);
    }

    if (self.encryptedSecretsReference.length > 0) {
      if (self.secretsLocation == Location.Inline) {
        revert NoInlineSecrets();
      }
      CBOR.writeString(buffer, "secretsLocation");
      CBOR.writeUInt256(buffer, uint256(self.secretsLocation));
      CBOR.writeString(buffer, "secrets");
      CBOR.writeBytes(buffer, self.encryptedSecretsReference);
    }

    return buffer.buf.buf;
  }

  /**
   * @notice Initializes a Chainlink Functions Request
   * @dev Sets the codeLocation and code on the request
   * @param self The uninitialized request
   * @param codeLocation The user provided source code location
   * @param language The programming language of the user code
   * @param source The user provided source code or a url
   */
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

  /**
   * @notice Initializes a Chainlink Functions Request
   * @dev Simplified version of initializeRequest for PoC
   * @param self The uninitialized request
   * @param javaScriptSource The user provided JS code (must not be empty)
   */
  function initializeRequestForInlineJavaScript(Request memory self, string memory javaScriptSource) internal pure {
    initializeRequest(self, Location.Inline, CodeLanguage.JavaScript, javaScriptSource);
  }

  /**
   * @notice Adds Remote user encrypted secrets to a Request
   * @param self The initialized request
   * @param encryptedSecretsReference Encrypted comma-separated string of URLs pointing to off-chain secrets
   */
  function addSecretsReference(Request memory self, bytes memory encryptedSecretsReference) internal pure {
    if (encryptedSecretsReference.length == 0) revert EmptySecrets();

    self.secretsLocation = Location.Remote;
    self.encryptedSecretsReference = encryptedSecretsReference;
  }

  /**
   * @notice Adds DON-hosted user slot id (referencing secrets) to a Request
   * @param self The initialized request
   * @param donSlotID Slot ID of the user's secrets hosted on DON
   */
  function addDONHostedSecrets(Request memory self, bytes1 donSlotID) internal pure {
    self.secretsLocation = Location.DONHosted;
    self.encryptedSecretsReference = new bytes(1);
    self.encryptedSecretsReference[0] = bytes1(donSlotID);
  }

  /**
   * @notice Adds args for the user run function
   * @param self The initialized request
   * @param args The array of args (must not be empty)
   */
  function addArgs(Request memory self, string[] memory args) internal pure {
    if (args.length == 0) revert EmptyArgs();

    self.args = args;
  }

  /**
   * @notice Add request data version to the request CBOR
   */
  function encodeRequest(bytes memory requestCBOR) internal pure returns (bytes memory) {
    return abi.encode(REQUEST_DATA_VERSION, requestCBOR);
  }

  /**
   * @notice Retrieve the request data version from an encoded request
   */
  function decodeRequest(bytes memory requestData) internal pure returns (uint16, bytes memory) {
    (uint16 version, bytes memory requestCBOR) = abi.decode(requestData, (uint16, bytes));
    return (version, requestCBOR);
  }
}
