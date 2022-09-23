// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../dev/ocr2dr/OCR2DR.sol";

contract OCR2DRTestHelper {
  using OCR2DR for OCR2DR.Request;
  using OCR2DR for OCR2DR.HttpQuery;
  using OCR2DR for OCR2DR.HttpHeader;

  OCR2DR.Request private req;

  event RequestData(bytes data);

  function closeEvent() public {
    emit RequestData(req.encodeCBOR());
  }

  function initializeRequestForInlineJavaScript(string memory sourceCode) public {
    OCR2DR.Request memory r;
    r.initializeRequestForInlineJavaScript(sourceCode);
    storeRequest(r);
  }

  function addSecrets(bytes memory secrets) public {
    OCR2DR.Request memory r = req;
    r.addSecrets(secrets);
    storeRequest(r);
  }

  function addTwoArgs(string memory arg1, string memory arg2) public {
    string[] memory args = new string[](2);
    args[0] = arg1;
    args[1] = arg2;
    OCR2DR.Request memory r = req;
    r.addArgs(args);
    storeRequest(r);
  }

  function addQuery(string memory url) public {
    OCR2DR.HttpQuery memory q;
    q.initializeHttpQuery(OCR2DR.HttpVerb.Get, url);
    OCR2DR.Request memory r = req;
    r.addHttpQuery(q);
    storeRequest(r);
  }

  function addQueryWithTwoHeaders(
    string memory url,
    string memory h1k,
    string memory h1v,
    string memory h2k,
    string memory h2v
  ) public {
    OCR2DR.HttpQuery memory q;
    q.initializeHttpQuery(OCR2DR.HttpVerb.Get, url);
    q.addHttpHeader(h1k, h1v);
    q.addHttpHeader(h2k, h2v);
    OCR2DR.Request memory r = req;
    r.addHttpQuery(q);
    storeRequest(r);
  }

  function storeRequest(OCR2DR.Request memory r) private {
    req.location = r.location;
    req.language = r.language;
    req.source = r.source;
    req.args = r.args;
    req.secrets = r.secrets;
    delete req.queries;
    for (uint256 i = 0; i < r.queries.length; i++) {
      OCR2DR.HttpQuery storage nq = req.queries.push();
      nq.verb = r.queries[i].verb;
      nq.url = r.queries[i].url;
      for (uint256 j = 0; j < r.queries[i].headers.length; j++) {
        nq.headers.push(r.queries[i].headers[j]);
      }
    }
  }
}
