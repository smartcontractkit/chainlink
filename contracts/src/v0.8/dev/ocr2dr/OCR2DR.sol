// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {CBORChainlink} from "../../vendor/CBORChainlink.sol";
import {BufferChainlink} from "../../vendor/BufferChainlink.sol";

/**
 * @title Library for OCR2 Direct Request functions
 */
library OCR2DR {
  uint256 internal constant DEFAULT_BUFFER_SIZE = 256;

  using CBORChainlink for BufferChainlink.buffer;

  enum Location {
    Inline
    // In future version we will add Remote location
  }

  enum CodeLanguage {
    JavaScript
    // In future version we may add other languages
  }

  enum HttpVerb {
    Get
    // In future version we may add other verbs
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
    Location codeLocation;
    Location secretsLocation;
    CodeLanguage language;
    string source; // Source code for Location.Inline or url for Location.Remote
    bytes secrets; // Encrypted secrets blob for Location.Inline or url for Location.Remote
    string[] args;
    HttpQuery[] queries;
  }

  error EmptySource();
  error EmptyUrl();
  error EmptyKey();
  error EmptyValue();
  error EmptyHeaders();
  error EmptyQueries();
  error EmptySecrets();
  error EmptyArgs();

  /**
   * @notice Encodes a Request to CBOR encoded bytes
   * @param self The request to encode
   * @return CBOR encoded bytes
   */
  function encodeCBOR(Request memory self) internal pure returns (bytes memory) {
    BufferChainlink.buffer memory buf;
    BufferChainlink.init(buf, DEFAULT_BUFFER_SIZE);

    buf.encodeString("codeLocation");
    buf.encodeUInt(uint256(self.codeLocation));

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
      buf.encodeString("secretsLocation");
      buf.encodeUInt(uint256(self.secretsLocation));
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
   */
  function initializeRequest(
    Request memory self,
    Location location,
    CodeLanguage language,
    string memory source
  ) internal pure {
    if (bytes(source).length == 0) revert EmptySource();

    self.codeLocation = location;
    self.language = language;
    self.source = source;
  }

  /**
   * @notice Initializes a OCR2DR Request
   * @dev Simplified version of initializeRequest for PoC
   * @param self The uninitialized request
   * @param javaScriptSource The user provided JS code (must not be empty)
   */
  function initializeRequestForInlineJavaScript(Request memory self, string memory javaScriptSource) internal pure {
    initializeRequest(self, Location.Inline, CodeLanguage.JavaScript, javaScriptSource);
  }

  /**
   * @notice Initializes a OCR2DR HttpQuery
   * @dev Sets the verb and url on the query
   * @param self The uninitialized query
   * @param verb The user provided HTTP verb
   * @param url The user provided HTTP/s url (must not be empty)
   */
  function initializeHttpQuery(
    HttpQuery memory self,
    HttpVerb verb,
    string memory url
  ) internal pure {
    if (bytes(url).length == 0) revert EmptyUrl();

    self.verb = verb;
    self.url = url;
  }

  /**
   * @notice Adds new HttpHeader to HttpQuery
   * @param self The initialized query
   * @param key HTTP header's key (must not be empty)
   * @param value HTTP header's value (must not be empty)
   */
  function addHttpHeader(
    HttpQuery memory self,
    string memory key,
    string memory value
  ) internal pure {
    if (bytes(key).length == 0) revert EmptyKey();
    if (bytes(value).length == 0) revert EmptyValue();

    HttpHeader[] memory headers = new HttpHeader[](self.headers.length + 1);
    for (uint256 i = 0; i < self.headers.length; i++) {
      headers[i] = self.headers[i];
    }
    headers[self.headers.length] = HttpHeader(key, value);
    self.headers = headers;
  }

  /**
   * @notice Set an array of HttpHeader to HttpQuery
   * @param self The initialized HttpQuery
   * @param headers The array of headers to be set (must not be empty)
   */
  function setHttpHeaders(HttpQuery memory self, HttpHeader[] memory headers) internal pure {
    if (headers.length == 0) revert EmptyHeaders();

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
   * @notice Set an array of HttpQuery to a Request
   * @param self The initialized request
   * @param queries The array of queries to be set (must not be empty)
   */
  function setHttpQueries(Request memory self, HttpQuery[] memory queries) internal pure {
    if (queries.length == 0) revert EmptyQueries();

    self.queries = queries;
  }

  /**
   * @notice Adds user encrypted secrets to a Request
   * @param self The initialized request
   * @param secrets The user encrypted secrets (must not be empty)
   */
  function addInlineSecrets(Request memory self, bytes memory secrets) internal pure {
    if (secrets.length == 0) revert EmptySecrets();

    self.secretsLocation = Location.Inline;
    self.secrets = secrets;
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
}
