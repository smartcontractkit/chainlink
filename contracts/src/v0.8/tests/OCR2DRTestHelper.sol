// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../dev/ocr2dr/OCR2DR.sol";

contract OCR2DRTestHelper {
  using OCR2DR for OCR2DR.Request;

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

  function storeRequest(OCR2DR.Request memory r) private {
    s_req.codeLocation = r.codeLocation;
    s_req.language = r.language;
    s_req.source = r.source;
    s_req.args = r.args;
    s_req.secretsLocation = r.secretsLocation;
    s_req.secrets = r.secrets;
  }
}
