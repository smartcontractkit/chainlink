// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import {CBORChainlink} from "../../../v0.7/vendor/CBORChainlink.sol";
import {BufferChainlink} from "../../../v0.7/vendor/BufferChainlink.sol";

/**
 * @title Library for OCR2 Direct Request functions
 */
library OCR2DR {
  uint256 internal constant defaultBufferSize = 256; // solhint-disable-line const-name-snakecase

  using CBORChainlink for BufferChainlink.buffer;

  enum CodeLocation {
    Inline
  }

  enum CodeLanguage {
    JavaScript
  }

  enum HttpVerb {
    Get
  }

  struct HttpHeader {
    string key;
    string value;
  }

  struct HttpQuery {
    HttpVerb verb;
    string url;
    HttpHeader[] headers;
  }

  struct Request {
    CodeLocation location;
    CodeLanguage language;
    string source;
    string[] args;
    bytes secrets;
    HttpQuery[] queries;
  }

  /**
   * @notice Encodes a Request to CBOR encoded bytes
   * @param self The request to encode
   * @return CBOR encoded bytes
   */
  function encodeCBOR(Request memory self) internal pure returns (bytes memory) {
    BufferChainlink.buffer memory buf;
    BufferChainlink.init(buf, defaultBufferSize);

    buf.encodeString("location");
    buf.encodeUInt(uint256(self.location));

    buf.encodeString("language");
    buf.encodeUInt(uint256(self.language));

    buf.encodeString("source");
    buf.encodeString(self.source);

    if (self.args.length > 0) {
      buf.encodeString("args");
      buf.startArray();
      for (uint256 i = 0; i < self.args.length; i++) {
        buf.encodeString(self.args[i]);
      }
      buf.endSequence();
    }

    if (self.secrets.length > 0) {
      buf.encodeString("secrets");
      buf.encodeBytes(self.secrets);
    }

    if (self.queries.length > 0) {
      buf.encodeString("queries");
      buf.startArray();
      for (uint256 i = 0; i < self.queries.length; i++) {
        buf.startMap();
        buf.encodeString("verb");
        buf.encodeUInt(uint256(self.queries[i].verb));
        buf.encodeString("url");
        buf.encodeString(self.queries[i].url);
        if (self.queries[i].headers.length > 0) {
            buf.encodeString("headers");
            buf.startMap();
            for (uint256 j = 0; j < self.queries[i].headers.length; j++) {
                buf.encodeString(self.queries[i].headers[j].key);
                buf.encodeString(self.queries[i].headers[j].value);
            }
            buf.endSequence();
        }
        buf.endSequence();
      }
      buf.endSequence();
    }

    return buf.buf;
  }

  /**
   * @notice Initializes a OCR2DR Request
   * @dev Sets the codeLocation and code on the request
   * @param self The uninitialized request
   * @param location The user provided source code location
   * @param language The programming language of the user code
   * @param source The user provided source code or a url
   * @return The initialized request
   */
  function initializeRequest(
    Request memory self,
    CodeLocation location,
    CodeLanguage language,
    string memory source
  ) internal pure returns (OCR2DR.Request memory) {
    require(bytes(source).length > 0);

    self.location = location;
    self.language = language;
    self.source = source;
    return self;
  }

  /**
   * @notice Initializes a OCR2DR Request
   * @dev Simplified version of initializeRequest for PoC
   * @param self The uninitialized request
   * @param javaScriptSource The user provided JS code
   * @return The initialized request
   */
  function initializeRequestForInlineJavaScript(Request memory self, string memory javaScriptSource)
    internal
    pure
    returns (OCR2DR.Request memory)
  {
    return initializeRequest(self, CodeLocation.Inline, CodeLanguage.JavaScript, javaScriptSource);
  }

  /**
   * @notice Initializes a OCR2DR HttpQuery
   * @dev Sets the verb and url on the query
   * @param self The uninitialized query
   * @param verb The user provided HTTP verb
   * @param url The user provided HTTP/s url
   * @return The initialized HttpQuery
   */
  function initializeHttpQuery(
    HttpQuery memory self,
    HttpVerb verb,
    string memory url
  ) internal pure returns (OCR2DR.HttpQuery memory) {
    require(bytes(url).length > 0);

    self.verb = verb;
    self.url = url;
    return self;
  }

  /**
   * @notice Adds new HttpHeader to HttpQuery
   * @param self The initialized query
   * @param key HTTP header's key
   * @param value HTTP header's value
   */
  function addHttpHeader(
    HttpQuery memory self,
    string memory key,
    string memory value
  ) internal pure {
    require(bytes(key).length > 0);
    require(bytes(value).length > 0);

    HttpHeader[] memory headers = new HttpHeader[](self.headers.length + 1);
    for (uint256 i = 0; i < self.headers.length; i++) {
      headers[i] = self.headers[i];
    }
    headers[self.headers.length] = HttpHeader(key, value);
    self.headers = headers;
  }

  /**
   * @notice Adds new HttpQuery to a Request
   * @param self The initialized request
   * @param query The initialized query to be added
   */
  function addHttpQuery(Request memory self, HttpQuery memory query) internal pure {
    HttpQuery[] memory queries = new HttpQuery[](self.queries.length + 1);
    for (uint256 i = 0; i < self.queries.length; i++) {
      queries[i] = self.queries[i];
    }
    queries[self.queries.length] = query;
    self.queries = queries;
  }

  /**
   * @notice Adds user encrypted secrets to a Request
   * @param self The initialized request
   * @param secrets The user encrypted secrets
   */
  function addSecrets(Request memory self, bytes memory secrets) internal pure {
    require(secrets.length > 0);

    self.secrets = secrets;
  }

  /**
   * @notice Adds args for the user run function
   * @param self The initialized request
   * @param args The array of args
   */
  function addArgs(Request memory self, string[] memory args) internal pure {
    require(args.length > 0);

    self.args = args;
  }
}
