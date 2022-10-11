// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../dev/ocr2dr/OCR2DR.sol";

contract OCR2DRTestHelper {
  using OCR2DR for OCR2DR.Request;
  using OCR2DR for OCR2DR.HttpQuery;
  using OCR2DR for OCR2DR.HttpHeader;

  OCR2DR.Request private s_req;

  event RequestData(bytes data);

  function closeEvent() public {
    emit RequestData(s_req.encodeCBOR());
  }

  function initializeRequestForInlineJavaScript(string memory sourceCode) public {
    OCR2DR.Request memory r;
    r.initializeRequestForInlineJavaScript(sourceCode);
    storeRequest(r);
  }

  function addSecrets(bytes memory secrets) public {
    OCR2DR.Request memory r = s_req;
    r.addInlineSecrets(secrets);
    storeRequest(r);
  }

  function addEmptyArgs() public pure {
    OCR2DR.Request memory r;
    string[] memory args;
    r.addArgs(args);
  }

  function addTwoArgs(string memory arg1, string memory arg2) public {
    string[] memory args = new string[](2);
    args[0] = arg1;
    args[1] = arg2;
    OCR2DR.Request memory r = s_req;
    r.addArgs(args);
    storeRequest(r);
  }

  function addQuery(string memory url) public {
    OCR2DR.HttpQuery memory q;
    q.initializeHttpQuery(OCR2DR.HttpVerb.Get, url);
    OCR2DR.Request memory r = s_req;
    r.addHttpQuery(q);
    storeRequest(r);
  }

  function setEmptyQueries() public pure {
    OCR2DR.Request memory r;
    OCR2DR.HttpQuery[] memory q;
    r.setHttpQueries(q);
  }

  function setEmptyHeaders() public pure {
    OCR2DR.HttpQuery memory q;
    OCR2DR.HttpHeader[] memory h;
    q.setHttpHeaders(h);
  }

  function setTwoQueries(string memory url1, string memory url2) public {
    OCR2DR.Request memory r = s_req;
    OCR2DR.HttpQuery memory q1;
    OCR2DR.HttpQuery memory q2;
    q1.initializeHttpQuery(OCR2DR.HttpVerb.Get, url1);
    q2.initializeHttpQuery(OCR2DR.HttpVerb.Get, url2);
    r.addHttpQuery(q1);
    r.addHttpQuery(q2);
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
    OCR2DR.HttpHeader[] memory headers = new OCR2DR.HttpHeader[](1);
    headers[0].key = h1k;
    headers[0].value = h1v;
    q.setHttpHeaders(headers);
    q.addHttpHeader(h2k, h2v);
    OCR2DR.Request memory r = s_req;
    r.addHttpQuery(q);
    storeRequest(r);
  }

  function storeRequest(OCR2DR.Request memory r) private {
    s_req.codeLocation = r.codeLocation;
    s_req.language = r.language;
    s_req.source = r.source;
    s_req.args = r.args;
    s_req.secretsLocation = r.secretsLocation;
    s_req.secrets = r.secrets;
    delete s_req.queries;
    for (uint256 i = 0; i < r.queries.length; i++) {
      OCR2DR.HttpQuery storage nq = s_req.queries.push();
      nq.verb = r.queries[i].verb;
      nq.url = r.queries[i].url;
      for (uint256 j = 0; j < r.queries[i].headers.length; j++) {
        nq.headers.push(r.queries[i].headers[j]);
      }
    }
  }
}
